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

// Portions copied from OpenTelemetry Collector (contrib), from the
// elastic exporter.
//
// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlp_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/elastic/apm-data/input/otlp"
	"github.com/elastic/apm-data/model/modelpb"
)

func TestConsumerConsumeLogs(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var processor modelpb.ProcessBatchFunc = func(_ context.Context, batch *modelpb.Batch) error {
			assert.Empty(t, batch)
			return nil
		}

		consumer := otlp.NewConsumer(otlp.ConsumerConfig{
			Processor: processor,
			Semaphore: semaphore.NewWeighted(100),
		})
		logs := plog.NewLogs()
		assert.NoError(t, consumer.ConsumeLogs(context.Background(), logs))
	})

	commonEvent := modelpb.APMEvent{
		Processor: modelpb.LogProcessor(),
		Agent: &modelpb.Agent{
			Name:    "otlp/go",
			Version: "unknown",
		},
		Service: &modelpb.Service{
			Name:     "unknown",
			Language: &modelpb.Language{Name: "go"},
		},
		Message: "a random log message",
		Event: &modelpb.Event{
			Severity: int64(plog.SeverityNumberInfo),
		},
		Log:           &modelpb.Log{Level: "Info"},
		Span:          &modelpb.Span{Id: "0200000000000000"},
		Trace:         &modelpb.Trace{Id: "01000000000000000000000000000000"},
		Labels:        modelpb.Labels{},
		NumericLabels: modelpb.NumericLabels{},
	}
	test := func(name string, body interface{}, expectedMessage string) {
		t.Run(name, func(t *testing.T) {
			logs := plog.NewLogs()
			resourceLogs := logs.ResourceLogs().AppendEmpty()
			logs.ResourceLogs().At(0).Resource().Attributes().PutStr(semconv.AttributeTelemetrySDKLanguage, "go")
			scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
			newLogRecord(body).CopyTo(scopeLogs.LogRecords().AppendEmpty())

			var processed modelpb.Batch
			var processor modelpb.ProcessBatchFunc = func(_ context.Context, batch *modelpb.Batch) error {
				if processed != nil {
					panic("already processes batch")
				}
				processed = *batch
				assert.NotNil(t, processed[0].Timestamp)
				processed[0].Timestamp = nil
				return nil
			}
			consumer := otlp.NewConsumer(otlp.ConsumerConfig{
				Processor: processor,
				Semaphore: semaphore.NewWeighted(100),
			})
			assert.NoError(t, consumer.ConsumeLogs(context.Background(), logs))

			now := time.Now().Unix()
			for _, e := range processed {
				assert.InDelta(t, now, e.Event.Received.AsTime().Unix(), 2)
				e.Event.Received = nil
			}

			expected := proto.Clone(&commonEvent).(*modelpb.APMEvent)
			expected.Message = expectedMessage
			assert.Empty(t, cmp.Diff(modelpb.Batch{expected}, processed, protocmp.Transform()))
		})
	}
	test("string_body", "a random log message", "a random log message")
	test("int_body", 1234, "1234")
	test("float_body", 1234.1234, "1234.1234")
	test("bool_body", true, "true")
	// TODO(marclop): How to test map body
}

func TestConsumeLogsSemaphore(t *testing.T) {
	logs := plog.NewLogs()
	var batches []*modelpb.Batch

	doneCh := make(chan struct{})
	recorder := modelpb.ProcessBatchFunc(func(ctx context.Context, batch *modelpb.Batch) error {
		<-doneCh
		batches = append(batches, batch)
		return nil
	})
	consumer := otlp.NewConsumer(otlp.ConsumerConfig{
		Processor: recorder,
		Semaphore: semaphore.NewWeighted(1),
	})

	startCh := make(chan struct{})
	go func() {
		close(startCh)
		assert.NoError(t, consumer.ConsumeLogs(context.Background(), logs))
	}()

	<-startCh
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	assert.Equal(t, consumer.ConsumeLogs(ctx, logs).Error(), "context deadline exceeded")
	close(doneCh)

	assert.NoError(t, consumer.ConsumeLogs(context.Background(), logs))
}

func TestConsumerConsumeLogsException(t *testing.T) {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	resourceAttrs := logs.ResourceLogs().At(0).Resource().Attributes()
	resourceAttrs.PutStr(semconv.AttributeTelemetrySDKLanguage, "java")
	resourceAttrs.PutStr("key0", "zero")
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	record1 := newLogRecord("foo")
	record1.Attributes().PutStr("key1", "one")
	record1.CopyTo(scopeLogs.LogRecords().AppendEmpty())

	record2 := newLogRecord("bar")
	record2.Attributes().PutStr("event.name", "crash")
	record2.Attributes().PutStr("event.domain", "device")
	record2.Attributes().PutStr("exception.type", "HighLevelException")
	record2.Attributes().PutStr("exception.message", "MidLevelException: LowLevelException")
	record2.Attributes().PutStr("exception.stacktrace", `
HighLevelException: MidLevelException: LowLevelException
	at Junk.a(Junk.java:13)
	at Junk.main(Junk.java:4)
Caused by: MidLevelException: LowLevelException
	at Junk.c(Junk.java:23)
	at Junk.b(Junk.java:17)
	at Junk.a(Junk.java:11)
	... 1 more
	Suppressed: java.lang.ArithmeticException: / by zero
		at Junk.c(Junk.java:25)
		... 3 more
Caused by: LowLevelException
	at Junk.e(Junk.java:37)
	at Junk.d(Junk.java:34)
	at Junk.c(Junk.java:21)
	... 3 more`[1:])
	record2.CopyTo(scopeLogs.LogRecords().AppendEmpty())

	var processed modelpb.Batch
	var processor modelpb.ProcessBatchFunc = func(_ context.Context, batch *modelpb.Batch) error {
		if processed != nil {
			panic("already processes batch")
		}
		processed = *batch
		assert.NotNil(t, processed[0].Timestamp)
		processed[0].Timestamp = nil
		return nil
	}
	consumer := otlp.NewConsumer(otlp.ConsumerConfig{
		Processor: processor,
		Semaphore: semaphore.NewWeighted(100),
	})
	assert.NoError(t, consumer.ConsumeLogs(context.Background(), logs))

	now := time.Now().Unix()
	for _, e := range processed {
		assert.InDelta(t, now, e.Event.Received.AsTime().Unix(), 2)
		e.Event.Received = nil
	}

	assert.Len(t, processed, 2)
	assert.Equal(t, modelpb.Labels{"key0": {Global: true, Value: "zero"}, "key1": {Value: "one"}}, modelpb.Labels(processed[0].Labels))
	assert.Empty(t, processed[0].NumericLabels)
	processed[1].Timestamp = nil
	out := cmp.Diff(&modelpb.APMEvent{
		Service: &modelpb.Service{
			Name: "unknown",
			Language: &modelpb.Language{
				Name: "java",
			},
		},
		Agent: &modelpb.Agent{
			Name:    "otlp/java",
			Version: "unknown",
		},
		Event: &modelpb.Event{
			Severity: int64(plog.SeverityNumberInfo),
			Kind:     "event",
			Category: "device",
			Type:     "error",
		},
		Labels:        modelpb.Labels{"key0": {Global: true, Value: "zero"}},
		NumericLabels: modelpb.NumericLabels{},
		Processor:     modelpb.ErrorProcessor(),
		Message:       "bar",
		Trace:         &modelpb.Trace{Id: "01000000000000000000000000000000"},
		Span:          &modelpb.Span{Id: "0200000000000000"},
		Log: &modelpb.Log{
			Level: "Info",
		},
		Error: &modelpb.Error{
			Type: "crash",
			Exception: &modelpb.Exception{
				Type:    "HighLevelException",
				Message: "MidLevelException: LowLevelException",
				Handled: newBool(true),
				Stacktrace: []*modelpb.StacktraceFrame{{
					Classname: "Junk",
					Function:  "a",
					Filename:  "Junk.java",
					Lineno:    newUint32(13),
				}, {
					Classname: "Junk",
					Function:  "main",
					Filename:  "Junk.java",
					Lineno:    newUint32(4),
				}},
				Cause: []*modelpb.Exception{{
					Message: "MidLevelException: LowLevelException",
					Handled: newBool(true),
					Stacktrace: []*modelpb.StacktraceFrame{{
						Classname: "Junk",
						Function:  "c",
						Filename:  "Junk.java",
						Lineno:    newUint32(23),
					}, {
						Classname: "Junk",
						Function:  "b",
						Filename:  "Junk.java",
						Lineno:    newUint32(17),
					}, {
						Classname: "Junk",
						Function:  "a",
						Filename:  "Junk.java",
						Lineno:    newUint32(11),
					}, {
						Classname: "Junk",
						Function:  "main",
						Filename:  "Junk.java",
						Lineno:    newUint32(4),
					}},
					Cause: []*modelpb.Exception{{
						Message: "LowLevelException",
						Handled: newBool(true),
						Stacktrace: []*modelpb.StacktraceFrame{{
							Classname: "Junk",
							Function:  "e",
							Filename:  "Junk.java",
							Lineno:    newUint32(37),
						}, {
							Classname: "Junk",
							Function:  "d",
							Filename:  "Junk.java",
							Lineno:    newUint32(34),
						}, {
							Classname: "Junk",
							Function:  "c",
							Filename:  "Junk.java",
							Lineno:    newUint32(21),
						}, {
							Classname: "Junk",
							Function:  "b",
							Filename:  "Junk.java",
							Lineno:    newUint32(17),
						}, {
							Classname: "Junk",
							Function:  "a",
							Filename:  "Junk.java",
							Lineno:    newUint32(11),
						}, {
							Classname: "Junk",
							Function:  "main",
							Filename:  "Junk.java",
							Lineno:    newUint32(4),
						}},
					}},
				}},
			},
		},
	}, processed[1],
		protocmp.Transform(),
		protocmp.IgnoreFields(&modelpb.Error{}, "id"),
	)
	assert.Empty(t, out)
}

func TestConsumerConsumeOTelEventLogs(t *testing.T) {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	resourceAttrs := logs.ResourceLogs().At(0).Resource().Attributes()
	resourceAttrs.PutStr(semconv.AttributeTelemetrySDKLanguage, "swift")
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	record1 := newLogRecord("") // no log body
	record1.Attributes().PutStr("event.domain", "device")
	record1.Attributes().PutStr("event.name", "MyEvent")

	record1.CopyTo(scopeLogs.LogRecords().AppendEmpty())

	var processed modelpb.Batch
	var processor modelpb.ProcessBatchFunc = func(_ context.Context, batch *modelpb.Batch) error {
		if processed != nil {
			panic("already processes batch")
		}
		processed = *batch
		assert.NotNil(t, processed[0].Timestamp)
		processed[0].Timestamp = timestamppb.New(time.Time{})
		return nil
	}
	consumer := otlp.NewConsumer(otlp.ConsumerConfig{
		Processor: processor,
		Semaphore: semaphore.NewWeighted(100),
	})
	assert.NoError(t, consumer.ConsumeLogs(context.Background(), logs))

	assert.Len(t, processed, 1)
	assert.Equal(t, "event", processed[0].Event.Kind)
	assert.Equal(t, "device", processed[0].Event.Category)
	assert.Equal(t, "MyEvent", processed[0].Event.Action)
}

func TestConsumerConsumeLogsLabels(t *testing.T) {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	resourceAttrs := logs.ResourceLogs().At(0).Resource().Attributes()
	resourceAttrs.PutStr(semconv.AttributeTelemetrySDKLanguage, "go")
	resourceAttrs.PutStr("key0", "zero")
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	record1 := newLogRecord("whatever")
	record1.Attributes().PutStr("key1", "one")
	record1.CopyTo(scopeLogs.LogRecords().AppendEmpty())

	record2 := newLogRecord("andever")
	record2.Attributes().PutDouble("key2", 2)
	record2.CopyTo(scopeLogs.LogRecords().AppendEmpty())

	record3 := newLogRecord("amen")
	record3.Attributes().PutStr("key3", "three")
	record3.Attributes().PutInt("key4", 4)
	record3.CopyTo(scopeLogs.LogRecords().AppendEmpty())

	var processed modelpb.Batch
	var processor modelpb.ProcessBatchFunc = func(_ context.Context, batch *modelpb.Batch) error {
		if processed != nil {
			panic("already processes batch")
		}
		processed = *batch
		assert.NotNil(t, processed[0].Timestamp)
		processed[0].Timestamp = timestamppb.New(time.Time{})
		return nil
	}
	consumer := otlp.NewConsumer(otlp.ConsumerConfig{
		Processor: processor,
		Semaphore: semaphore.NewWeighted(100),
	})
	assert.NoError(t, consumer.ConsumeLogs(context.Background(), logs))

	assert.Len(t, processed, 3)
	assert.Equal(t, modelpb.Labels{"key0": {Global: true, Value: "zero"}, "key1": {Value: "one"}}, modelpb.Labels(processed[0].Labels))
	assert.Empty(t, processed[0].NumericLabels)
	assert.Equal(t, modelpb.Labels{"key0": {Global: true, Value: "zero"}}, modelpb.Labels(processed[1].Labels))
	assert.Equal(t, modelpb.NumericLabels{"key2": {Value: 2}}, modelpb.NumericLabels(processed[1].NumericLabels))
	assert.Equal(t, modelpb.Labels{"key0": {Global: true, Value: "zero"}, "key3": {Value: "three"}}, modelpb.Labels(processed[2].Labels))
	assert.Equal(t, modelpb.NumericLabels{"key4": {Value: 4}}, modelpb.NumericLabels(processed[2].NumericLabels))
}

func newLogRecord(body interface{}) plog.LogRecord {
	otelLogRecord := plog.NewLogRecord()
	otelLogRecord.SetTraceID(pcommon.TraceID{1})
	otelLogRecord.SetSpanID(pcommon.SpanID{2})
	otelLogRecord.SetSeverityNumber(plog.SeverityNumberInfo)
	otelLogRecord.SetSeverityText("Info")
	otelLogRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	switch b := body.(type) {
	case string:
		otelLogRecord.Body().SetStr(b)
	case int:
		otelLogRecord.Body().SetInt(int64(b))
	case float64:
		otelLogRecord.Body().SetDouble(b)
	case bool:
		otelLogRecord.Body().SetBool(b)
		// case map[string]string:
		// TODO(marclop) figure out how to set the body since it cannot be set
		// as a map.
		// otelLogRecord.Body()
	}
	return otelLogRecord
}
