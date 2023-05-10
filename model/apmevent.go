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

package model

import (
	"time"

	"go.elastic.co/fastjson"
)

const (
	// timestampFormat formats timestamps according to Elasticsearch's
	// strict_date_optional_time date format, which includes a fractional
	// seconds component.
	timestampFormat = "2006-01-02T15:04:05.000Z07:00"
)

// APMEvent holds the details of an APM event.
//
// Exactly one of the event fields should be non-nil.
type APMEvent struct {
	// Timestamp holds the event timestamp.
	//
	// See https://www.elastic.co/guide/en/ecs/current/ecs-base.html#field-timestamp
	Timestamp time.Time
	Span      *Span
	// NumericLabels holds the numeric (scaled_float) labels to apply to the event.
	// Supports slice values.
	NumericLabels NumericLabels
	// Labels holds the string (keyword) labels to apply to the event, stored as
	// keywords. Supports slice values.
	//
	// See https://www.elastic.co/guide/en/ecs/current/ecs-base.html#field-labels
	Labels      Labels
	Transaction *Transaction
	Metricset   *Metricset
	Error       *Error
	Cloud       Cloud
	Service     Service
	FAAS        FAAS
	Network     Network
	Container   Container
	User        User
	Device      Device
	Kubernetes  Kubernetes
	Observer    Observer
	// DataStream optionally holds data stream identifiers.
	DataStream DataStream
	Agent      Agent
	Processor  Processor
	HTTP       HTTP
	UserAgent  UserAgent
	Parent     Parent
	// Message holds the message for log events.
	//
	// See https://www.elastic.co/guide/en/ecs/current/ecs-base.html#field-message
	Message     string
	Trace       Trace
	Host        Host
	URL         URL
	Log         Log
	Source      Source
	Client      Client
	Child       Child
	Destination Destination
	Session     Session
	Process     Process
	Event       Event
}

// TODO(axw) remove this and generate the JSON encoding directly.
func (e *APMEvent) fields() map[string]any {
	var fields mapStr
	if e.Transaction != nil {
		fields.maybeSetMapStr("transaction", e.Transaction.fields(e.Processor))
	}
	if e.Span != nil {
		fields.maybeSetMapStr("span", e.Span.fields())
	}
	if e.Metricset != nil {
		e.Metricset.setFields(&fields)
	}
	if e.Error != nil {
		fields.maybeSetMapStr("error", e.Error.fields())
	}

	// Set high resolution timestamp.
	//
	// TODO(axw) change @timestamp to use date_nanos, and remove this field.
	if !e.Timestamp.IsZero() {
		switch e.Processor {
		case TransactionProcessor, SpanProcessor, ErrorProcessor:
			fields.set("timestamp", map[string]any{
				"us": int(e.Timestamp.UnixNano() / 1000),
			})
		}
	}

	// Set top-level field sets.
	e.DataStream.setFields(&fields)
	fields.maybeSetMapStr("service", e.Service.Fields())
	fields.maybeSetMapStr("agent", e.Agent.fields())
	fields.maybeSetMapStr("observer", e.Observer.Fields())
	fields.maybeSetMapStr("host", e.Host.fields())
	fields.maybeSetMapStr("device", e.Device.fields())
	fields.maybeSetMapStr("process", e.Process.fields())
	fields.maybeSetMapStr("user", e.User.fields())
	fields.maybeSetMapStr("client", e.Client.fields())
	fields.maybeSetMapStr("source", e.Source.fields())
	fields.maybeSetMapStr("destination", e.Destination.fields())
	fields.maybeSetMapStr("user_agent", e.UserAgent.fields())
	fields.maybeSetMapStr("container", e.Container.fields())
	fields.maybeSetMapStr("kubernetes", e.Kubernetes.fields())
	fields.maybeSetMapStr("cloud", e.Cloud.fields())
	fields.maybeSetMapStr("network", e.Network.fields())
	fields.maybeSetMapStr("labels", e.Labels.fields())
	fields.maybeSetMapStr("numeric_labels", e.NumericLabels.fields())
	fields.maybeSetMapStr("event", e.Event.fields())
	fields.maybeSetMapStr("url", e.URL.fields())
	fields.maybeSetMapStr("session", e.Session.fields())
	fields.maybeSetMapStr("parent", e.Parent.fields())
	fields.maybeSetMapStr("child", e.Child.fields())
	fields.maybeSetMapStr("processor", e.Processor.fields())
	fields.maybeSetMapStr("trace", e.Trace.fields())
	fields.maybeSetString("message", e.Message)
	fields.maybeSetMapStr("http", e.HTTP.fields())
	fields.maybeSetMapStr("faas", e.FAAS.fields())
	fields.maybeSetMapStr("log", e.Log.fields())
	return map[string]any(fields)
}

// MarshalJSON marshals e as JSON.
func (e *APMEvent) MarshalJSON() ([]byte, error) {
	var w fastjson.Writer
	if err := e.MarshalFastJSON(&w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// MarshalFastJSON marshals e as JSON, writing the result to w.
func (e *APMEvent) MarshalFastJSON(w *fastjson.Writer) error {
	w.RawByte('{')
	w.RawString(`"@timestamp":"`)
	w.Time(e.Timestamp, timestampFormat)
	w.RawByte('"')
	for k, v := range e.fields() {
		w.RawByte(',')
		w.String(k)
		w.RawByte(':')
		if err := encodeAny(v, w); err != nil {
			return err
		}
	}
	w.RawByte('}')
	return nil
}

func encodeAny(v any, out *fastjson.Writer) error {
	switch v := v.(type) {
	case map[string]any:
		return encodeMap(v, out)
	default:
		return fastjson.Marshal(out, v)
	}
}

func encodeMap(v map[string]any, out *fastjson.Writer) error {
	out.RawByte('{')
	first := true
	for k, v := range v {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.String(k)
		out.RawByte(':')
		if err := encodeAny(v, out); err != nil {
			return err
		}
	}
	out.RawByte('}')
	return nil
}
