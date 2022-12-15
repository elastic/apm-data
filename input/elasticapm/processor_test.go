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
	"net/netip"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/apm-data/model"
)

func TestHandleStreamReaderError(t *testing.T) {
	readErr := errors.New("read failed")
	var calls int
	var reader readerFunc = func(p []byte) (int, error) {
		calls++
		if calls > 1 {
			return 0, readErr
		}
		buf := bytes.NewBuffer(nil)
		buf.WriteString(validMetadata + "\n")
		for i := 0; i < 5; i++ {
			buf.WriteString(validTransaction + "\n")
		}
		return copy(p, buf.Bytes()), nil
	}

	sp := NewProcessor(Config{
		MaxEventSize: 100 * 1024,
		Semaphore:    make(chan struct{}, 1),
	})

	var actualResult Result
	err := sp.HandleStream(
		context.Background(), false, model.APMEvent{},
		reader, 10, nopBatchProcessor{}, &actualResult,
	)
	assert.Equal(t, readErr, err)
	assert.Equal(t, Result{Accepted: 5}, actualResult)
}

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) {
	return f(p)
}

func TestHandleStreamBatchProcessorError(t *testing.T) {
	payload := validMetadata + "\n" + validTransaction + "\n"
	for _, test := range []struct {
		name string
		err  error
	}{{
		name: "NotQueueFull",
		err:  errors.New("queue is not full, something else is wrong"),
	}, {
		name: "QueueFull",
		err:  ErrQueueFull,
	}} {
		sp := NewProcessor(Config{
			MaxEventSize: 100 * 1024,
			Semaphore:    make(chan struct{}, 1),
		})
		processor := model.ProcessBatchFunc(func(context.Context, *model.Batch) error {
			return test.err
		})

		var actualResult Result
		err := sp.HandleStream(
			context.Background(), false, model.APMEvent{},
			strings.NewReader(payload), 10, processor, &actualResult,
		)
		assert.Equal(t, test.err, err)
		assert.Zero(t, actualResult)
	}
}

func TestHandleStreamErrors(t *testing.T) {
	var (
		invalidEvent        = `{ "transaction": { "id": 12345, "trace_id": "0123456789abcdef0123456789abcdef", "parent_id": "abcdefabcdef01234567", "type": "request", "duration": 32.592981, "span_count": { "started": 21 } } }   `
		invalidJSONEvent    = `{ "invalid-json" }`
		invalidJSONMetadata = `{"metadata": {"invalid-json"}}`
		invalidMetadata     = `{"metadata": {"user": null}}`
		invalidMetadata2    = `{"not": "metadata"}`
		invalidEventType    = `{"tennis-court": {"name": "Centre Court, Wimbledon"}}`
		tooLargeEvent       = strings.Repeat("*", len(validMetadata)*2)
	)

	for _, test := range []struct {
		name     string
		payload  string
		tooLarge int
		invalid  int
		errors   []error // per-event errors
		err      error   // stream-level error
	}{{
		name:    "InvalidEvent",
		payload: validMetadata + "\n" + invalidEvent + "\n",
		invalid: 1,
		errors: []error{
			&InvalidInputError{
				Message:  `decode error: data read error: v2.transactionRoot.Transaction: v2.transaction.ID: ReadString: expects " or n,`,
				Document: invalidEvent,
			},
		},
	}, {
		name:    "InvalidJSONEvent",
		payload: validMetadata + "\n" + invalidJSONEvent + "\n",
		invalid: 1,
		errors: []error{
			&InvalidInputError{
				Message:  `did not recognize object type: "invalid-json"`,
				Document: invalidJSONEvent,
			},
		},
	}, {
		name:    "InvalidJSONMetadata",
		payload: invalidJSONMetadata + "\n",
		err: &InvalidInputError{
			Message:  "decode error: data read error: v2.metadataRoot.Metadata: v2.metadata.readFieldHash: expect :,",
			Document: invalidJSONMetadata,
		},
	}, {
		name:    "InvalidMetadata",
		payload: invalidMetadata + "\n",
		err: &InvalidInputError{
			Message:  "validation error: 'metadata' required",
			Document: invalidMetadata,
		},
	}, {
		name:    "InvalidMetadata2",
		payload: invalidMetadata2 + "\n",
		err: &InvalidInputError{
			Message:  `"metadata" or "m" required`,
			Document: invalidMetadata2,
		},
	}, {
		name:    "UnrecognizedEvent",
		payload: validMetadata + "\n" + invalidEventType + "\n",
		invalid: 1,
		errors: []error{
			&InvalidInputError{
				Message:  `did not recognize object type: "tennis-court"`,
				Document: invalidEventType,
			},
		},
	}, {
		name:     "TooLargeEvent",
		payload:  validMetadata + "\n" + tooLargeEvent + "\n",
		tooLarge: 1,
		errors: []error{
			&InvalidInputError{
				TooLarge: true,
				Message:  "event exceeded the permitted size",
				Document: tooLargeEvent[:len(validMetadata)+1],
			},
		},
	}} {
		t.Run(test.name, func(t *testing.T) {
			var actualResult Result
			p := NewProcessor(Config{
				MaxEventSize: len(validMetadata) + 1,
				Semaphore:    make(chan struct{}, 1),
			})
			err := p.HandleStream(
				context.Background(), false, model.APMEvent{},
				strings.NewReader(test.payload), 10,
				nopBatchProcessor{}, &actualResult,
			)
			assert.Equal(t, test.err, err)
			assert.Zero(t, actualResult.Accepted)
			assert.Equal(t, test.errors, actualResult.Errors)
			assert.Equal(t, test.tooLarge, actualResult.TooLarge)
			assert.Equal(t, test.invalid, actualResult.Invalid)
		})
	}
}

func TestHandleStream(t *testing.T) {
	var events []model.APMEvent
	batchProcessor := model.ProcessBatchFunc(func(ctx context.Context, batch *model.Batch) error {
		events = append(events, (*batch)...)
		return nil
	})

	payload := strings.Join([]string{
		validMetadata,
		validError,
		validMetricset,
		validSpan,
		validTransaction,
		validLog,
		"", // final newline
	}, "\n")

	p := NewProcessor(Config{
		MaxEventSize: 100 * 1024,
		Semaphore:    make(chan struct{}, 1),
	})
	err := p.HandleStream(
		context.Background(), false, model.APMEvent{},
		strings.NewReader(payload), 10, batchProcessor,
		&Result{},
	)
	require.NoError(t, err)

	processors := make([]model.Processor, len(events))
	for i, event := range events {
		processors[i] = event.Processor
	}
	assert.Equal(t, []model.Processor{
		model.ErrorProcessor,
		model.MetricsetProcessor,
		model.SpanProcessor,
		model.TransactionProcessor,
		model.LogProcessor,
	}, processors)
}

func TestHandleStreamRUMv3(t *testing.T) {
	var events []model.APMEvent
	batchProcessor := model.ProcessBatchFunc(func(ctx context.Context, batch *model.Batch) error {
		events = append(events, (*batch)...)
		return nil
	})

	payload := strings.Join([]string{
		validRUMv3Metadata,
		validRUMv3Error,
		validRUMv3Transaction,
		"", // final newline
	}, "\n")

	p := NewProcessor(Config{
		MaxEventSize: 100 * 1024,
		Semaphore:    make(chan struct{}, 1),
	})
	err := p.HandleStream(
		context.Background(), false, model.APMEvent{},
		strings.NewReader(payload), 10, batchProcessor,
		&Result{},
	)
	require.NoError(t, err)

	processors := make([]model.Processor, len(events))
	for i, event := range events {
		processors[i] = event.Processor
	}
	assert.Equal(t, []model.Processor{
		model.ErrorProcessor,
		model.TransactionProcessor,
		model.MetricsetProcessor,
		model.MetricsetProcessor,
		model.SpanProcessor,
		model.SpanProcessor,
		model.SpanProcessor,
		model.SpanProcessor,
		model.SpanProcessor,
		model.SpanProcessor,
		model.SpanProcessor,
		model.SpanProcessor,
	}, processors)
}

func TestHandleStreamBaseEvent(t *testing.T) {
	requestTimestamp := time.Date(2018, 8, 1, 10, 0, 0, 0, time.UTC)

	baseEvent := model.APMEvent{
		Timestamp: requestTimestamp,
		UserAgent: model.UserAgent{Original: "rum-2.0"},
		Source:    model.Source{IP: netip.MustParseAddr("192.0.0.1")},
		Client:    model.Client{IP: netip.MustParseAddr("192.0.0.2")}, // X-Forwarded-For
	}

	var events []model.APMEvent
	batchProcessor := model.ProcessBatchFunc(func(ctx context.Context, batch *model.Batch) error {
		events = append(events, (*batch)...)
		return nil
	})

	payload := validMetadata + "\n" + validRUMv2Span + "\n"
	p := NewProcessor(Config{
		MaxEventSize: 100 * 1024,
		Semaphore:    make(chan struct{}, 1),
	})
	err := p.HandleStream(
		context.Background(), false, baseEvent,
		strings.NewReader(payload), 10, batchProcessor,
		&Result{},
	)
	require.NoError(t, err)

	assert.Len(t, events, 1)
	assert.Equal(t, "rum-2.0", events[0].UserAgent.Original)
	assert.Equal(t, baseEvent.Source, events[0].Source)
	assert.Equal(t, baseEvent.Client, events[0].Client)
	assert.Equal(t, requestTimestamp.Add(50*time.Millisecond), events[0].Timestamp) // span's start is "50"
}

func TestLabelLeak(t *testing.T) {
	payload := `{"metadata": {"service": {"name": "testsvc", "environment": "staging", "version": null, "agent": {"name": "python", "version": "6.9.1"}, "language": {"name": "python", "version": "3.10.4"}, "runtime": {"name": "CPython", "version": "3.10.4"}, "framework": {"name": "flask", "version": "2.1.1"}}, "process": {"pid": 2112739, "ppid": 2112738, "argv": ["/home/stuart/workspace/sdh/581/venv/lib/python3.10/site-packages/flask/__main__.py", "run"], "title": null}, "system": {"hostname": "slaptop", "architecture": "x86_64", "platform": "linux"}, "labels": {"ci_commit": "unknown", "numeric": 1}}}
{"transaction": {"id": "88dee29a6571b948", "trace_id": "ba7f5d18ac4c7f39d1ff070c79b2bea5", "name": "GET /withlabels", "type": "request", "duration": 1.6199999999999999, "result": "HTTP 2xx", "timestamp": 1652185276804681, "outcome": "success", "sampled": true, "span_count": {"started": 0, "dropped": 0}, "sample_rate": 1.0, "context": {"request": {"env": {"REMOTE_ADDR": "127.0.0.1", "SERVER_NAME": "127.0.0.1", "SERVER_PORT": "5000"}, "method": "GET", "socket": {"remote_address": "127.0.0.1"}, "cookies": {}, "headers": {"host": "localhost:5000", "user-agent": "curl/7.81.0", "accept": "*/*", "app-os": "Android", "content-type": "application/json; charset=utf-8", "content-length": "29"}, "url": {"full": "http://localhost:5000/withlabels?second_with_labels", "protocol": "http:", "hostname": "localhost", "pathname": "/withlabels", "port": "5000", "search": "?second_with_labels"}}, "response": {"status_code": 200, "headers": {"Content-Type": "application/json", "Content-Length": "14"}}, "tags": {"appOs": "Android", "email_set": "hello@hello.com", "time_set": 1652185276}}}}
{"transaction": {"id": "ba5c6d6c1ab44bd1", "trace_id": "88c0a00431531a80c5ca9a41fe115f41", "name": "GET /nolabels", "type": "request", "duration": 0.652, "result": "HTTP 2xx", "timestamp": 1652185278813952, "outcome": "success", "sampled": true, "span_count": {"started": 0, "dropped": 0}, "sample_rate": 1.0, "context": {"request": {"env": {"REMOTE_ADDR": "127.0.0.1", "SERVER_NAME": "127.0.0.1", "SERVER_PORT": "5000"}, "method": "GET", "socket": {"remote_address": "127.0.0.1"}, "cookies": {}, "headers": {"host": "localhost:5000", "user-agent": "curl/7.81.0", "accept": "*/*"}, "url": {"full": "http://localhost:5000/nolabels?third_no_label", "protocol": "http:", "hostname": "localhost", "pathname": "/nolabels", "port": "5000", "search": "?third_no_label"}}, "response": {"status_code": 200, "headers": {"Content-Type": "text/html; charset=utf-8", "Content-Length": "14"}}, "tags": {}}}}`

	baseEvent := model.APMEvent{
		Host: model.Host{IP: []netip.Addr{netip.MustParseAddr("192.0.0.1")}},
	}

	var processed *model.Batch
	batchProcessor := model.ProcessBatchFunc(func(_ context.Context, b *model.Batch) error {
		processed = b
		return nil
	})

	p := NewProcessor(Config{
		MaxEventSize: 100 * 1024,
		Semaphore:    make(chan struct{}, 1),
	})
	var actualResult Result
	err := p.HandleStream(context.Background(), false, baseEvent, strings.NewReader(payload), 10, batchProcessor, &actualResult)
	require.NoError(t, err)

	txs := *processed
	assert.Len(t, txs, 2)
	// Assert first tx
	assert.Equal(t, model.NumericLabels{
		"time_set": {Value: 1652185276},
		"numeric":  {Global: true, Value: 1},
	}, txs[0].NumericLabels)
	assert.Equal(t, model.Labels{
		"appOs":     {Value: "Android"},
		"email_set": {Value: "hello@hello.com"},
		"ci_commit": {Global: true, Value: "unknown"},
	}, txs[0].Labels)

	// Assert second tx
	assert.Equal(t, model.NumericLabels{"numeric": {Global: true, Value: 1}}, txs[1].NumericLabels)
	assert.Equal(t, model.Labels{"ci_commit": {Global: true, Value: "unknown"}}, txs[1].Labels)
}

func TestConcurrentAsync(t *testing.T) {
	smallBatch := validMetadata + "\n" + validTransaction + "\n"
	bigBatch := validMetadata + "\n" + strings.Repeat(validTransaction+"\n", 2000)

	base := model.APMEvent{Host: model.Host{IP: []netip.Addr{netip.MustParseAddr("192.0.0.1")}}}
	type testCase struct {
		payload       string
		sem, requests int
		fullSem       bool
	}

	test := func(tc testCase) (pResult Result) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		p := NewProcessor(Config{
			MaxEventSize: 100 * 1024,
			Semaphore:    make(chan struct{}, tc.sem),
		})
		if tc.fullSem {
			for i := 0; i < tc.sem; i++ {
				p.sem <- struct{}{}
			}
		}
		handleStream := func(ctx context.Context, bp *accountProcessor) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result Result
				err := p.HandleStream(ctx, true, base, strings.NewReader(tc.payload), 10, bp, &result)
				if err != nil {
					result.addError(err)
				}
				if !tc.fullSem {
					select {
					case <-bp.batch:
					case <-ctx.Done():
					}
				}
				mu.Lock()
				if len(result.Errors) > 0 {
					pResult.Errors = append(pResult.Errors, result.Errors...)
				}
				mu.Unlock()
			}()
		}
		batchProcessor := &accountProcessor{batch: make(chan *model.Batch, tc.requests)}
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		for i := 0; i < tc.requests; i++ {
			handleStream(ctx, batchProcessor)
		}
		wg.Wait()
		if !tc.fullSem {
			// Try to acquire the lock to make sure all the requests have been handled
			// and the locks have been released.
			for i := 0; i < tc.sem; i++ {
				p.semAcquire(context.Background(), false)
			}
		}
		processed := atomic.LoadUint64(&batchProcessor.processed)
		pResult.Accepted += int(processed)
		return
	}

	t.Run("semaphore_full", func(t *testing.T) {
		res := test(testCase{
			sem:      2,
			requests: 3,
			fullSem:  true,
			payload:  smallBatch,
		})
		assert.Equal(t, 0, res.Accepted)
		assert.Equal(t, 3, len(res.Errors))
		for _, err := range res.Errors {
			assert.ErrorIs(t, err, ErrQueueFull)
		}
	})
	t.Run("semaphore_undersized", func(t *testing.T) {
		res := test(testCase{
			sem:      2,
			requests: 100,
			payload:  bigBatch,
		})
		// When the semaphore is full, `ErrQueueFull` is returned.
		assert.Greater(t, len(res.Errors), 0)
		for _, err := range res.Errors {
			assert.EqualError(t, err, ErrQueueFull.Error())
		}
	})
	t.Run("semaphore_empty", func(t *testing.T) {
		res := test(testCase{
			sem:      5,
			requests: 5,
			payload:  smallBatch,
		})
		assert.Equal(t, 5, res.Accepted)
		assert.Equal(t, 0, len(res.Errors))

		res = test(testCase{
			sem:      5,
			requests: 5,
			payload:  bigBatch,
		})
		assert.GreaterOrEqual(t, res.Accepted, 5)
		// all the request will return with an error since only 50 events of
		// each (5 requests * batch size) will be processed.
		assert.Equal(t, 5, len(res.Errors))
	})
	t.Run("semaphore_empty_incorrect_metadata", func(t *testing.T) {
		res := test(testCase{
			sem:      5,
			requests: 5,
			payload:  `{"metadata": {"siervice":{}}}`,
		})
		assert.Equal(t, 0, res.Accepted)
		assert.Len(t, res.Errors, 5)

		incorrectEvent := `{"metadata": {"service": {"name": "testsvc", "environment": "staging", "version": null, "agent": {"name": "python", "version": "6.9.1"}, "language": {"name": "python", "version": "3.10.4"}, "runtime": {"name": "CPython", "version": "3.10.4"}, "framework": {"name": "flask", "version": "2.1.1"}}, "process": {"pid": 2112739, "ppid": 2112738, "argv": ["/home/stuart/workspace/sdh/581/venv/lib/python3.10/site-packages/flask/__main__.py", "run"], "title": null}, "system": {"hostname": "slaptop", "architecture": "x86_64", "platform": "linux"}, "labels": {"ci_commit": "unknown", "numeric": 1}}}
{"some_incorrect_event": {}}`
		res = test(testCase{
			sem:      5,
			requests: 2,
			payload:  incorrectEvent,
		})
		assert.Equal(t, 0, res.Accepted)
		assert.Len(t, res.Errors, 2)
	})
}

type nopBatchProcessor struct{}

func (nopBatchProcessor) ProcessBatch(context.Context, *model.Batch) error {
	return nil
}

type accountProcessor struct {
	batch     chan *model.Batch
	processed uint64
}

func (p *accountProcessor) ProcessBatch(ctx context.Context, b *model.Batch) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if p.batch != nil {
		select {
		case p.batch <- b:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	atomic.AddUint64(&p.processed, 1)
	return nil
}
