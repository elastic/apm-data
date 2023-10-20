// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package elasticapm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"go.elastic.co/apm/v2"
	"go.uber.org/zap"

	"github.com/elastic/apm-data/input"
	"github.com/elastic/apm-data/input/elasticapm/internal/decoder"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/rumv3"
	v2 "github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/v2"
	"github.com/elastic/apm-data/model/modelpb"
)

var (
	errUnrecognizedObject = errors.New("did not recognize object type")

	errEmptyBody = errors.New("empty body")

	// ErrQueueFull may be returned by HandleStream when the internal
	// queue is full.
	ErrQueueFull = errors.New("queue is full")

	batchPool sync.Pool
)

const (
	errorEventType            = "error"
	metricsetEventType        = "metricset"
	spanEventType             = "span"
	transactionEventType      = "transaction"
	logEventType              = "log"
	rumv3ErrorEventType       = "e"
	rumv3TransactionEventType = "x"

	v2MetadataKey    = "metadata"
	rumv3MetadataKey = "m"
)

// Processor decodes a streams and is safe for concurrent use. The processor
// accepts a channel that is used as a semaphore to control the maximum
// concurrent number of stream decode operations that can happen at any time.
// The buffered channel is meant to be shared between all the processors so
// the concurrency limit is shared between all the intake endpoints.
type Processor struct {
	streamReaderPool sync.Pool
	sem              input.Semaphore
	logger           *zap.Logger
	MaxEventSize     int
}

// Config holds configuration for Processor constructors.
type Config struct {
	// Semaphore holds a semaphore on which Processor.HandleStream will acquire a
	// token before proceeding, to limit concurrency.
	Semaphore input.Semaphore
	// Logger holds a logger for the processor. If Logger is nil,
	// then no logging will be performed.
	Logger *zap.Logger
	// MaxEventSize holds the maximum event size, in bytes.
	MaxEventSize int
}

// StreamHandler is an interface for handling an Elastic APM agent ND-JSON event
// stream, implemented by processor/stream.
type StreamHandler interface {
	HandleStream(
		ctx context.Context,
		async bool,
		baseEvent *modelpb.APMEvent,
		stream io.Reader,
		batchSize int,
		processor modelpb.BatchProcessor,
		out *Result,
	) error
}

// NewProcessor returns a new Processor for processing an event stream from
// Elastic APM agents.
func NewProcessor(cfg Config) *Processor {
	if cfg.Logger == nil {
		cfg.Logger = zap.NewNop()
	}
	return &Processor{
		MaxEventSize: cfg.MaxEventSize,
		sem:          cfg.Semaphore,
		logger:       cfg.Logger,
	}
}

func (p *Processor) readMetadata(reader *streamReader, out *modelpb.APMEvent) error {
	body, err := reader.ReadAhead()
	if err != nil {
		if err == io.EOF {
			if len(reader.LatestLine()) == 0 {
				return errEmptyBody
			}
			return &InvalidInputError{
				Message:  "EOF while reading metadata",
				Document: string(reader.LatestLine()),
			}
		}
		return reader.wrapError(err)
	}
	switch key := p.identifyEventType(body); string(key) {
	case v2MetadataKey:
		if err := v2.DecodeNestedMetadata(reader, out); err != nil {
			return reader.wrapError(err)
		}
	case rumv3MetadataKey:
		if err := rumv3.DecodeNestedMetadata(reader, out); err != nil {
			return reader.wrapError(err)
		}
	default:
		return &InvalidInputError{
			Message:  fmt.Sprintf("%q or %q required", v2MetadataKey, rumv3MetadataKey),
			Document: string(reader.LatestLine()),
		}
	}
	return nil
}

// identifyEventType takes a reader and reads ahead the first key of the
// underlying json input. This method makes some assumptions met by the
// input format:
// - the input is in JSON format
// - every valid ndjson line only has one root key
// - the bytes that we must match on are ASCII
func (p *Processor) identifyEventType(body []byte) []byte {
	// find event type, trim spaces and account for single and double quotes
	var quote byte
	var key []byte
	for i, r := range body {
		if r == '"' || r == '\'' {
			quote = r
			key = body[i+1:]
			break
		}
	}
	end := bytes.IndexByte(key, quote)
	if end == -1 {
		return nil
	}
	return key[:end]
}

// readBatch reads up to `batchSize` events from the ndjson stream into
// batch, returning the number of events read and any error encountered.
// Callers should always process the n > 0 events returned before considering
// the error err.
func (p *Processor) readBatch(
	ctx context.Context,
	baseEvent *modelpb.APMEvent,
	batchSize int,
	batch *modelpb.Batch,
	reader *streamReader,
	result *Result,
) (int, error) {

	// input events are decoded and appended to the batch
	origLen := len(*batch)
	for i := 0; i < batchSize && !reader.IsEOF(); i++ {
		body, err := reader.ReadAhead()
		if err != nil && err != io.EOF {
			err := reader.wrapError(err)
			var invalidInput *InvalidInputError
			if errors.As(err, &invalidInput) {
				result.addError(err)
				continue
			}
			// return early, we assume we can only recover from a input error types
			return len(*batch) - origLen, err
		}
		if len(body) == 0 {
			// required for backwards compatibility - sending empty lines was permitted in previous versions
			continue
		}
		// We copy the event for each iteration of the batch, as to avoid
		// shallow copies of Labels and NumericLabels.
		input := modeldecoder.Input{Base: baseEvent}
		switch eventType := p.identifyEventType(body); string(eventType) {
		case errorEventType:
			err = v2.DecodeNestedError(reader, &input, batch)
		case metricsetEventType:
			err = v2.DecodeNestedMetricset(reader, &input, batch)
		case spanEventType:
			err = v2.DecodeNestedSpan(reader, &input, batch)
		case transactionEventType:
			err = v2.DecodeNestedTransaction(reader, &input, batch)
		case logEventType:
			err = v2.DecodeNestedLog(reader, &input, batch)
		case rumv3ErrorEventType:
			err = rumv3.DecodeNestedError(reader, &input, batch)
		case rumv3TransactionEventType:
			err = rumv3.DecodeNestedTransaction(reader, &input, batch)
		default:
			err = fmt.Errorf("%w: %q", errUnrecognizedObject, eventType)
		}
		if err != nil && err != io.EOF {
			result.addError(&InvalidInputError{
				Message:  err.Error(),
				Document: string(reader.LatestLine()),
			})
		}
	}
	if reader.IsEOF() {
		return len(*batch) - origLen, io.EOF
	}
	return len(*batch) - origLen, nil
}

// HandleStream processes a stream of events in batches of batchSize at a time,
// updating result as events are accepted, or per-event errors occur.
//
// HandleStream will return an error when a terminal stream-level error occurs,
// such as the rate limit being exceeded, or due to authorization errors. In
// this case the result will only cover the subset of events accepted.
//
// Callers must not access result concurrently with HandleStream.
func (p *Processor) HandleStream(
	ctx context.Context,
	async bool,
	baseEvent *modelpb.APMEvent,
	reader io.Reader,
	batchSize int,
	processor modelpb.BatchProcessor,
	result *Result,
) error {
	sp, ctx := apm.StartSpan(ctx, "Stream", "Reporter")
	defer sp.End()
	// Limit the number of concurrent batch decodes.
	//
	// The semaphore defaults to 200 (N), only allowing N requests to read
	// an cache Y events (determined by batchSize) from the batch.
	//
	// Clients can set async to true which makes the processor process the
	// events in the background. Returns with an error `ErrQueueFull`
	// if the semaphore is full. When asynchronous processing is requested,
	// the batches are decoded synchronously, but the batch is processed
	// asynchronously.
	if err := p.semAcquire(ctx, async); err != nil {
		return err
	}
	sr := p.getStreamReader(reader)

	// Release the semaphore on early exit; this will be set to false
	// for asynchronous requests once we may no longer exit early.
	shouldReleaseSemaphore := true
	defer func() {
		sr.release()
		if shouldReleaseSemaphore {
			p.sem.Release(1)
		}
	}()

	// The first item is the metadata object.
	if err := p.readMetadata(sr, baseEvent); err != nil {
		if err == errEmptyBody {
			return nil
		}
		// no point in continuing if we couldn't read the metadata
		if _, ok := err.(*InvalidInputError); ok {
			return err
		}
		return &InvalidInputError{
			Message:  err.Error(),
			Document: string(sr.LatestLine()),
		}
	}

	if async {
		// The semaphore is released by handleStream
		shouldReleaseSemaphore = false
	}
	first := true
	for {
		err := p.handleStream(ctx, async, baseEvent, batchSize, sr, processor, result, first)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if first {
			first = false
		}
	}
}

func (p *Processor) handleStream(
	ctx context.Context,
	async bool,
	baseEvent *modelpb.APMEvent,
	batchSize int,
	sr *streamReader,
	processor modelpb.BatchProcessor,
	result *Result,
	first bool,
) (readErr error) {
	// Async requests will re-aquire the semaphore if it has more events than
	// `batchSize`. In that event, the semaphore will be acquired again. If
	// the semaphore is full, `ErrQueueFull` is returned.
	// The first iteration will not acquire the semaphore since it's already
	// acquired in the caller function.
	var n int
	if async {
		if !first {
			if err := p.semAcquire(ctx, async); err != nil {
				return err
			}
		}
		defer func() {
			// If no events have been read on an asynchronous request, release
			// the semaphore since the processing goroutine isn't scheduled.
			if n == 0 {
				p.sem.Release(1)
			}
		}()
	}
	var batch modelpb.Batch
	if b, ok := batchPool.Get().(*modelpb.Batch); ok {
		batch = (*b)[:0]
	}
	n, readErr = p.readBatch(ctx, baseEvent, batchSize, &batch, sr, result)
	if n == 0 {
		// No events to process, return the batch to the pool.
		batchPool.Put(&batch)
		return readErr
	}
	// Async requests are processed in the background and once the batch has
	// been processed, the semaphore is released.
	if async {
		go func() {
			defer p.sem.Release(1)
			if err := p.processBatch(ctx, processor, &batch); err != nil {
				p.logger.Error("failed handling async request", zap.Error(err))
			}
		}()
	} else {
		if err := p.processBatch(ctx, processor, &batch); err != nil {
			return err
		}
		result.Accepted += n
	}
	return readErr
}

// processBatch processes the batch and returns it to the pool after it's been processed.
func (p *Processor) processBatch(ctx context.Context, processor modelpb.BatchProcessor, batch *modelpb.Batch) error {
	defer func() {
		for i := range *batch {
			(*batch)[i].ReturnToVTPool()
		}
		batchPool.Put(batch)
	}()
	return processor.ProcessBatch(ctx, batch)
}

// getStreamReader returns a streamReader that reads ND-JSON lines from r.
func (p *Processor) getStreamReader(r io.Reader) *streamReader {
	if sr, ok := p.streamReaderPool.Get().(*streamReader); ok {
		sr.Reset(r)
		return sr
	}
	return &streamReader{
		processor:           p,
		NDJSONStreamDecoder: decoder.NewNDJSONStreamDecoder(r, p.MaxEventSize),
	}
}

func (p *Processor) semAcquire(ctx context.Context, async bool) error {
	if async {
		if ok := p.sem.TryAcquire(1); !ok {
			return ErrQueueFull
		}
		return nil
	}
	return p.sem.Acquire(ctx, 1)
}

// streamReader wraps NDJSONStreamReader, converting errors to stream errors.
type streamReader struct {
	processor *Processor
	*decoder.NDJSONStreamDecoder
}

// release releases the streamReader, adding it to its Processor's sync.Pool.
// The streamReader must not be used after release returns.
func (sr *streamReader) release() {
	sr.Reset(nil)
	sr.processor.streamReaderPool.Put(sr)
}

func (sr *streamReader) wrapError(err error) error {
	if _, ok := err.(decoder.JSONDecodeError); ok {
		return &InvalidInputError{
			Message:  err.Error(),
			Document: string(sr.LatestLine()),
		}
	}

	var e = err
	if err, ok := err.(modeldecoder.DecoderError); ok {
		e = err.Unwrap()
	}
	if errors.Is(e, decoder.ErrLineTooLong) {
		return &InvalidInputError{
			TooLarge: true,
			Message:  "event exceeded the permitted size",
			Document: string(sr.LatestLine()),
		}
	}
	return err
}
