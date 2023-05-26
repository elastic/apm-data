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

package modelpb

import (
	"time"

	"github.com/elastic/apm-data/model/internal/modeljson"
)

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
		CompressionStrategy: c.CompressionStrategy.String(),
		Count:               int(c.Count),
		Sum:                 modeljson.SpanCompositeSum{US: sumDuration.Microseconds()},
	}
}

func (e *Span) toModelJSON(out *modeljson.Span) {
	*out = modeljson.Span{
		ID:                  e.Id,
		Name:                e.Name,
		Type:                e.Type,
		Subtype:             e.Subtype,
		Action:              e.Action,
		Kind:                e.Kind,
		Sync:                e.Sync,
		RepresentativeCount: e.RepresentativeCount,
	}
	if e.SelfTime != nil {
		out.SelfTime = modeljson.AggregatedDuration{
			Count: int(e.SelfTime.Count),
			Sum:   e.SelfTime.Sum.AsDuration(),
		}
	}
	if e.Db != nil {
		out.DB = &modeljson.DB{}
		e.Db.toModelJSON(out.DB)
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
				Type:     e.DestinationService.Type,
				Name:     e.DestinationService.Name,
				Resource: e.DestinationService.Resource,
			},
		}
		if e.DestinationService.ResponseTime != nil {
			out.Destination.Service.ResponseTime = modeljson.AggregatedDuration{
				Count: int(e.DestinationService.ResponseTime.Count),
				Sum:   e.DestinationService.ResponseTime.Sum.AsDuration(),
			}
		}
	}
	if n := len(e.Links); n > 0 {
		out.Links = make([]modeljson.SpanLink, n)
		for i, link := range e.Links {
			out.Links[i] = modeljson.SpanLink{
				Trace: modeljson.SpanLinkTrace{ID: link.Trace.Id},
				Span:  modeljson.SpanLinkSpan{ID: link.Span.Id},
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
