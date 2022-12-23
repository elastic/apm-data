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
)

var (
	// SpanProcessor is the Processor value that should be assigned to span events.
	SpanProcessor = Processor{Name: "transaction", Event: "span"}
)

type Span struct {
	Message            *Message
	Composite          *Composite
	DestinationService *DestinationService
	DB                 *DB
	Sync               *bool
	// Kind holds the span kind: "CLIENT", "SERVER", "PRODUCER", "CONSUMER" and "INTERNAL".
	Kind string
	// Action holds the span action: "query", "execute", etc.
	Action string
	// Subtype holds the span subtype: "http", "sql", etc.
	Subtype string
	ID      string
	// Type holds the span type: "external", "db", etc.
	Type string
	// Name holds the span name: "SELECT FROM table_name", etc.
	Name       string
	Stacktrace Stacktrace
	Links      []SpanLink
	// SelfTime holds the aggregated span durations, for breakdown metrics.
	SelfTime AggregatedDuration
	// RepresentativeCount holds the approximate number of spans that
	// this span represents for aggregation. This will only be set when
	// the sampling rate is known, and is used for scaling metrics.
	RepresentativeCount float64
}

// DB contains information related to a database query of a span event
type DB struct {
	RowsAffected *int
	Instance     string
	Statement    string
	Type         string
	UserName     string
	Link         string
}

// DestinationService contains information about the destination service of a span event
type DestinationService struct {
	Type     string // Deprecated
	Name     string // Deprecated
	Resource string // Deprecated

	// ResponseTime holds aggregated span durations for the destination service resource.
	ResponseTime AggregatedDuration
}

// Composite holds details on a group of spans compressed into one.
type Composite struct {
	CompressionStrategy string
	Count               int
	Sum                 float64 // milliseconds
}

func (db *DB) fields() map[string]any {
	if db == nil {
		return nil
	}
	var fields, user mapStr
	fields.maybeSetString("instance", db.Instance)
	fields.maybeSetString("statement", db.Statement)
	fields.maybeSetString("type", db.Type)
	fields.maybeSetString("link", db.Link)
	fields.maybeSetIntptr("rows_affected", db.RowsAffected)
	if user.maybeSetString("name", db.UserName) {
		fields.set("user", map[string]any(user))
	}
	return map[string]any(fields)
}

func (d *DestinationService) fields() map[string]any {
	if d == nil {
		return nil
	}
	var fields mapStr
	fields.maybeSetString("type", d.Type)
	fields.maybeSetString("name", d.Name)
	fields.maybeSetString("resource", d.Resource)
	fields.maybeSetMapStr("response_time", d.ResponseTime.fields())
	return map[string]any(fields)
}

func (c *Composite) fields() map[string]any {
	if c == nil {
		return nil
	}
	var fields mapStr
	sumDuration := time.Duration(c.Sum * float64(time.Millisecond))
	fields.set("sum", map[string]any{"us": int(sumDuration.Microseconds())})
	fields.set("count", c.Count)
	fields.set("compression_strategy", c.CompressionStrategy)

	return map[string]any(fields)
}

func (e *Span) fields() map[string]any {
	var span mapStr
	span.maybeSetString("name", e.Name)
	span.maybeSetString("type", e.Type)
	span.maybeSetString("id", e.ID)
	span.maybeSetString("kind", e.Kind)
	span.maybeSetString("subtype", e.Subtype)
	span.maybeSetString("action", e.Action)
	span.maybeSetBool("sync", e.Sync)
	span.maybeSetMapStr("db", e.DB.fields())
	span.maybeSetMapStr("message", e.Message.Fields())
	span.maybeSetMapStr("composite", e.Composite.fields())
	if destinationServiceFields := e.DestinationService.fields(); len(destinationServiceFields) > 0 {
		destinationMap, ok := span["destination"].(map[string]any)
		if !ok {
			destinationMap = make(map[string]any)
			span.set("destination", destinationMap)
		}
		destinationMap["service"] = destinationServiceFields
	}
	if st := e.Stacktrace.transform(); len(st) > 0 {
		span.set("stacktrace", st)
	}
	span.maybeSetMapStr("self_time", e.SelfTime.fields())
	if n := len(e.Links); n > 0 {
		links := make([]map[string]any, n)
		for i, link := range e.Links {
			links[i] = link.fields()
		}
		span.set("links", links)
	}
	if e.RepresentativeCount > 0 {
		span.set("representative_count", e.RepresentativeCount)
	}
	return map[string]any(span)
}
