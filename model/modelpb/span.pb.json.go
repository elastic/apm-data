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
	out.Instance = db.Instance
	out.Statement = db.Statement
	out.Type = db.Type
	out.Link = db.Link
	out.RowsAffected = db.RowsAffected
	out.User = modeljson.DBUser{Name: db.UserName}
}

func (c *Composite) toModelJSON(out *modeljson.SpanComposite) {
	sumDuration := time.Duration(c.Sum * float64(time.Millisecond))
	out.CompressionStrategy = c.CompressionStrategy.String()
	out.Count = int(c.Count)
	out.Sum = modeljson.SpanCompositeSum{US: sumDuration.Microseconds()}
}

func (e *Span) toModelJSON(out *modeljson.Span) {
	out.ID = e.Id
	out.Name = e.Name
	out.Type = e.Type
	out.Subtype = e.Subtype
	out.Action = e.Action
	out.Kind = e.Kind
	out.Sync = e.Sync
	out.RepresentativeCount = e.RepresentativeCount
	if e.SelfTime != nil {
		out.SelfTime = modeljson.AggregatedDuration{
			Count: int(e.SelfTime.Count),
			Sum:   e.SelfTime.Sum.AsDuration(),
		}
	}
	if e.Db != nil {
		e.Db.toModelJSON(out.DB)
	} else {
		out.DB = nil
	}
	if e.Message != nil {
		e.Message.toModelJSON(out.Message)
	} else {
		out.Message = nil
	}
	if e.Composite != nil {
		e.Composite.toModelJSON(out.Composite)
	} else {
		out.Composite = nil
	}
	if e.DestinationService != nil {
		out.Destination.Service.Type = e.DestinationService.Type
		out.Destination.Service.Name = e.DestinationService.Name
		out.Destination.Service.Resource = e.DestinationService.Resource
		if e.DestinationService.ResponseTime != nil {
			out.Destination.Service.ResponseTime = modeljson.AggregatedDuration{
				Count: int(e.DestinationService.ResponseTime.Count),
				Sum:   e.DestinationService.ResponseTime.Sum.AsDuration(),
			}
		}
	} else {
		out.Destination = nil
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
