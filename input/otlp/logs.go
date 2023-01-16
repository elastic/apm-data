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

package otlp

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
	"go.uber.org/zap"

	"github.com/elastic/apm-data/model"
)

func (c *Consumer) ConsumeLogs(ctx context.Context, logs plog.Logs) error {
	receiveTimestamp := time.Now()
	c.config.Logger.Debug("consuming logs", zap.Stringer("logs", logsStringer(logs)))
	resourceLogs := logs.ResourceLogs()
	batch := make(model.Batch, 0, resourceLogs.Len())
	for i := 0; i < resourceLogs.Len(); i++ {
		c.convertResourceLogs(resourceLogs.At(i), receiveTimestamp, &batch)
	}
	return c.config.Processor.ProcessBatch(ctx, &batch)
}

func (c *Consumer) convertResourceLogs(resourceLogs plog.ResourceLogs, receiveTimestamp time.Time, out *model.Batch) {
	var timeDelta time.Duration
	resource := resourceLogs.Resource()
	baseEvent := model.APMEvent{Processor: model.LogProcessor}
	translateResourceMetadata(resource, &baseEvent)

	if exportTimestamp, ok := exportTimestamp(resource); ok {
		timeDelta = receiveTimestamp.Sub(exportTimestamp)
	}
	scopeLogs := resourceLogs.ScopeLogs()
	for i := 0; i < scopeLogs.Len(); i++ {
		c.convertInstrumentationLibraryLogs(scopeLogs.At(i), baseEvent, timeDelta, out)
	}
}

func (c *Consumer) convertInstrumentationLibraryLogs(
	in plog.ScopeLogs,
	baseEvent model.APMEvent,
	timeDelta time.Duration,
	out *model.Batch,
) {
	otelLogs := in.LogRecords()
	for i := 0; i < otelLogs.Len(); i++ {
		event := c.convertLogRecord(otelLogs.At(i), baseEvent, timeDelta)
		*out = append(*out, event)
	}
}

func (c *Consumer) convertLogRecord(
	record plog.LogRecord,
	baseEvent model.APMEvent,
	timeDelta time.Duration,
) model.APMEvent {
	event := baseEvent
	initEventLabels(&event)
	event.Timestamp = record.Timestamp().AsTime().Add(timeDelta)
	event.Event.Severity = int64(record.SeverityNumber())
	event.Log.Level = record.SeverityText()
	if body := record.Body(); body.Type() != pcommon.ValueTypeEmpty {
		event.Message = body.AsString()
		if body.Type() == pcommon.ValueTypeMap {
			setLabels(body.Map(), &event)
		}
	}
	if traceID := record.TraceID(); !traceID.IsEmpty() {
		event.Trace.ID = traceID.HexString()
	}
	if spanID := record.SpanID(); !spanID.IsEmpty() {
		if event.Span == nil {
			event.Span = &model.Span{}
		}
		event.Span.ID = spanID.HexString()
	}
	attrs := record.Attributes()
	setLabels(attrs, &event)
	convertErrorLogRecord(attrs, &event)
	return event
}

func convertErrorLogRecord(m pcommon.Map, event *model.APMEvent) {
	exceptionMessage, exceptionMessageOk := m.Get(semconv.AttributeExceptionMessage)
	exceptionType, exceptionTypeOk := m.Get(semconv.AttributeExceptionType)

	if exceptionMessageOk && exceptionTypeOk {
		exceptionStacktrace, _ := m.Get(semconv.AttributeExceptionStacktrace)
		exceptionEscaped, _ := m.Get(semconv.AttributeExceptionEscaped)
		// Per OpenTelemetry semantic conventions:
		//   `At least one of the following sets of attributes is required:
		//   - exception.type
		//   - exception.message`
		event.Error = convertOpenTelemetryExceptionSpanEvent(
			exceptionType.Str(), exceptionMessage.Str(), exceptionStacktrace.Str(),
			exceptionEscaped.Bool(), event.Service.Language.Name,
		)
	}

	eventName, eventNameOk := m.Get("event.name")
	if eventNameOk {
		if event.Error == nil {
			event.Error = &model.Error{}
		}
		if eventName.Str() == "crash" {
			event.Event.Category = "device"
		}
		event.Error.Type = eventName.Str()
	}

	if event.Error != nil {
		event.Processor = model.ErrorProcessor
		event.Event.Kind = "event"
		event.Event.Type = "error"
	}
}

func setLabels(m pcommon.Map, event *model.APMEvent) {
	m.Range(func(k string, v pcommon.Value) bool {
		switch k {
		case semconv.AttributeExceptionMessage:
		case semconv.AttributeExceptionStacktrace:
		case semconv.AttributeExceptionType:
		case semconv.AttributeExceptionEscaped:
		default:
			setLabel(replaceDots(k), event, ifaceAttributeValue(v))
		}
		return true
	})
}
