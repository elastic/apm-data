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

	"github.com/elastic/apm-data/model/internal/modeljson"
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
	RowsAffected *uint32
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

func (db *DB) toModelJSON(out *modeljson.DB) {
	*out = modeljson.DB{
		Instance:     db.Instance,
		Statement:    db.Statement,
		Type:         db.Type,
		Link:         db.Link,
		RowsAffected: db.RowsAffected,
		User:         modeljson.DBUser{Name: db.UserName},
	}
}

func (c *Composite) toModelJSON(out *modeljson.SpanComposite) {
	sumDuration := time.Duration(c.Sum * float64(time.Millisecond))
	*out = modeljson.SpanComposite{
		CompressionStrategy: c.CompressionStrategy,
		Count:               c.Count,
		Sum:                 modeljson.SpanCompositeSum{US: sumDuration.Microseconds()},
	}
}

func (e *Span) toModelJSON(out *modeljson.Span) {
	*out = modeljson.Span{
		ID:                  e.ID,
		Name:                e.Name,
		Type:                e.Type,
		Subtype:             e.Subtype,
		Action:              e.Action,
		Kind:                e.Kind,
		Sync:                e.Sync,
		RepresentativeCount: e.RepresentativeCount,
		SelfTime:            modeljson.AggregatedDuration(e.SelfTime),
	}
	if e.DB != nil {
		out.DB = &modeljson.DB{}
		e.DB.toModelJSON(out.DB)
	}
	if e.Message != nil {
		out.Message = &modeljson.Message{}
		e.Message.toModelJSON(out.Message)
	}
	if e.Composite != nil {
		out.Composite = &modeljson.SpanComposite{}
		e.Composite.toModelJSON(out.Composite)
	}
	if e.DestinationService != nil {
		out.Destination = &modeljson.SpanDestination{
			Service: modeljson.SpanDestinationService{
				Type:         e.DestinationService.Type,
				Name:         e.DestinationService.Name,
				Resource:     e.DestinationService.Resource,
				ResponseTime: modeljson.AggregatedDuration(e.DestinationService.ResponseTime),
			},
		}
	}
	if n := len(e.Links); n > 0 {
		out.Links = make([]modeljson.SpanLink, n)
		for i, link := range e.Links {
			out.Links[i] = modeljson.SpanLink{
				Trace: modeljson.SpanLinkTrace{ID: link.Trace.ID},
				Span:  modeljson.SpanLinkSpan{ID: link.Span.ID},
			}
		}
	}
	if n := len(e.Stacktrace); n > 0 {
		out.Stacktrace = make([]modeljson.StacktraceFrame, n)
		for i, frame := range e.Stacktrace {
			frame.toModelJSON(&out.Stacktrace[i])
		}
	}
}
