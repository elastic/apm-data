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
	"github.com/elastic/apm-data/model/modelpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	// SpanProcessor is the Processor value that should be assigned to span events.
	SpanProcessor = Processor{Name: "transaction", Event: "span"}

	compressionStrategyText = map[string]modelpb.CompressionStrategy{
		"exact_match": modelpb.CompressionStrategy_COMPRESSION_STRATEGY_EXACT_MATCH,
		"same_kind":   modelpb.CompressionStrategy_COMPRESSION_STRATEGY_SAME_KIND,
	}
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

func (db *DB) toModelProtobuf(out *modelpb.DB) {
	*out = modelpb.DB{
		RowsAffected: db.RowsAffected,
		Instance:     db.Instance,
		Statement:    db.Statement,
		Type:         db.Type,
		UserName:     db.UserName,
		Link:         db.Link,
	}
}

func (c *Composite) toModelProtobuf(out *modelpb.Composite) {
	*out = modelpb.Composite{
		CompressionStrategy: compressionStrategyText[c.CompressionStrategy],
		Count:               uint32(c.Count),
		Sum:                 c.Sum,
	}
}

func (e *Span) toModelProtobuf(out *modelpb.Span) {
	*out = modelpb.Span{
		Sync:                e.Sync,
		Kind:                e.Kind,
		Action:              e.Action,
		Subtype:             e.Subtype,
		Id:                  e.ID,
		Type:                e.Type,
		Name:                e.Name,
		RepresentativeCount: e.RepresentativeCount,
	}
	if !isZero(e.SelfTime) {
		out.SelfTime = &modelpb.AggregatedDuration{
			Count: int64(e.SelfTime.Count),
			Sum:   durationpb.New(e.SelfTime.Sum),
		}
	}
	if e.DB != nil {
		out.Db = &modelpb.DB{}
		e.DB.toModelProtobuf(out.Db)
	}
	if e.Message != nil {
		out.Message = &modelpb.Message{}
		e.Message.toModelProtobuf(out.Message)
	}
	if e.Composite != nil {
		out.Composite = &modelpb.Composite{}
		e.Composite.toModelProtobuf(out.Composite)
	}
	if e.DestinationService != nil {
		out.DestinationService = &modelpb.DestinationService{
			Type:     e.DestinationService.Type,
			Name:     e.DestinationService.Name,
			Resource: e.DestinationService.Resource,
		}
		if !isZero(e.DestinationService.ResponseTime) {
			out.DestinationService.ResponseTime = &modelpb.AggregatedDuration{
				Count: int64(e.DestinationService.ResponseTime.Count),
				Sum:   durationpb.New(e.DestinationService.ResponseTime.Sum),
			}
		}
	}
	if n := len(e.Links); n > 0 {
		out.Links = make([]*modelpb.SpanLink, n)
		for i, link := range e.Links {
			out.Links[i] = &modelpb.SpanLink{
				Trace: &modelpb.Trace{Id: link.Trace.ID},
				Span:  &modelpb.Span{Id: link.Span.ID},
			}
		}
	}
	if n := len(e.Stacktrace); n > 0 {
		out.Stacktrace = make([]*modelpb.StacktraceFrame, n)
		for i, frame := range e.Stacktrace {
			outFrame := modelpb.StacktraceFrame{}
			frame.toModelProtobuf(&outFrame)
			out.Stacktrace[i] = &outFrame
		}
	}
}
