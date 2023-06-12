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
	"encoding/hex"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/elastic/apm-data/model/modelpb"
)

func (c *Consumer) ConsumeLogs(ctx context.Context, logs plog.Logs) error {
	receiveTimestamp := time.Now()
	c.config.Logger.Debug("consuming logs", zap.Stringer("logs", logsStringer(logs)))
	resourceLogs := logs.ResourceLogs()
	batch := make(modelpb.Batch, 0, resourceLogs.Len())
	for i := 0; i < resourceLogs.Len(); i++ {
		c.convertResourceLogs(resourceLogs.At(i), receiveTimestamp, &batch)
	}
	return c.config.Processor.ProcessBatch(ctx, &batch)
}

func (c *Consumer) convertResourceLogs(resourceLogs plog.ResourceLogs, receiveTimestamp time.Time, out *modelpb.Batch) {
	var timeDelta time.Duration
	resource := resourceLogs.Resource()
	baseEvent := modelpb.APMEvent{Processor: modelpb.LogProcessor()}
	translateResourceMetadata(resource, &baseEvent)

	if exportTimestamp, ok := exportTimestamp(resource); ok {
		timeDelta = receiveTimestamp.Sub(exportTimestamp)
	}
	scopeLogs := resourceLogs.ScopeLogs()
	for i := 0; i < scopeLogs.Len(); i++ {
		c.convertInstrumentationLibraryLogs(scopeLogs.At(i), &baseEvent, timeDelta, out)
	}
}

func (c *Consumer) convertInstrumentationLibraryLogs(
	in plog.ScopeLogs,
	baseEvent *modelpb.APMEvent,
	timeDelta time.Duration,
	out *modelpb.Batch,
) {
	otelLogs := in.LogRecords()
	for i := 0; i < otelLogs.Len(); i++ {
		event := c.convertLogRecord(otelLogs.At(i), baseEvent, timeDelta)
		*out = append(*out, event)
	}
}

func (c *Consumer) convertLogRecord(
	record plog.LogRecord,
	baseEvent *modelpb.APMEvent,
	timeDelta time.Duration,
) *modelpb.APMEvent {
	event := proto.Clone(baseEvent).(*modelpb.APMEvent)
	initEventLabels(event)
	event.Timestamp = timestamppb.New(record.Timestamp().AsTime().Add(timeDelta))
	if event.Event == nil {
		event.Event = &modelpb.Event{}
	}
	event.Event.Severity = int64(record.SeverityNumber())
	if event.Log == nil {
		event.Log = &modelpb.Log{}
	}
	event.Log.Level = record.SeverityText()
	if body := record.Body(); body.Type() != pcommon.ValueTypeEmpty {
		event.Message = body.AsString()
		if body.Type() == pcommon.ValueTypeMap {
			setLabels(body.Map(), event)
		}
	}
	if traceID := record.TraceID(); !traceID.IsEmpty() {
		event.Trace = &modelpb.Trace{
			Id: hex.EncodeToString(traceID[:]),
		}
	}
	if spanID := record.SpanID(); !spanID.IsEmpty() {
		if event.Span == nil {
			event.Span = &modelpb.Span{}
		}
		event.Span.Id = hex.EncodeToString(spanID[:])
	}
	attrs := record.Attributes()

	var exceptionMessage string
	var exceptionStacktrace string
	var exceptionType string
	var exceptionEscaped bool
	var eventName string
	var eventDomain string
	attrs.Range(func(k string, v pcommon.Value) bool {
		switch k {
		case semconv.AttributeExceptionMessage:
			exceptionMessage = v.Str()
		case semconv.AttributeExceptionStacktrace:
			exceptionStacktrace = v.Str()
		case semconv.AttributeExceptionType:
			exceptionType = v.Str()
		case semconv.AttributeExceptionEscaped:
			exceptionEscaped = v.Bool()
		case "event.name":
			eventName = v.Str()
		case "event.domain":
			eventDomain = v.Str()
		case "session.id":
			event.Session.Id = v.Str()
		default:
			setLabel(replaceDots(k), event, ifaceAttributeValue(v))
		}
		return true
	})

	if exceptionMessage != "" && exceptionType != "" {
		// Per OpenTelemetry semantic conventions:
		//   `At least one of the following sets of attributes is required:
		//   - exception.type
		//   - exception.message`
		event.Error = convertOpenTelemetryExceptionSpanEvent(
			exceptionType, exceptionMessage, exceptionStacktrace,
			exceptionEscaped, event.Service.Language.Name,
		)
	}

	if eventDomain == "device" && eventName != "" {
		event.Event.Category = "device"
		if eventName == "crash" {
			if event.Error == nil {
				event.Error = &modelpb.Error{}
			}
			event.Error.Type = "crash"
		} else {
			event.Event.Kind = "event"
			event.Event.Action = eventName
		}
	}

	if event.Error != nil {
		event.Processor = modelpb.ErrorProcessor()
		event.Event.Kind = "event"
		event.Event.Type = "error"
	}

	return event
}

func setLabels(m pcommon.Map, event *modelpb.APMEvent) {
	m.Range(func(k string, v pcommon.Value) bool {
		setLabel(replaceDots(k), event, ifaceAttributeValue(v))
		return true
	})
}
