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

	"github.com/elastic/apm-data/model/modelpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Event holds information about an event, in ECS terms.
//
// https://www.elastic.co/guide/en/ecs/current/ecs-event.html
type Event struct {
	// Outcome holds the event outcome: "success", "failure", or "unknown".
	Outcome string
	// Action holds the action captured by the event for log events.
	Action string
	// Dataset holds the the dataset which produces the events. If an event
	// source publishes more than one type of log or events (e.g. access log,
	// error log), the dataset is used to specify which one the event comes from.
	Dataset string

	// Kind holds the kind of the event. The highest categorization field in the hierarchy.
	Kind string

	// Category holds the event category. The second categorization field in the hierarchy.
	Category string

	// Type holds the event type. The third categorization field in the hierarchy.
	Type string

	// SuccessCount holds an aggregated count of events with different outcomes.
	// A "failure" adds to the Count. A "success" adds to both the Count and the
	// Sum. An "unknown" has no effect. If Count is zero, it will be omitted
	// from the output event.
	SuccessCount SummaryMetric
	// Duration holds the event duration.
	//
	// Duration is only added as a field (`duration`) if greater than zero.
	Duration time.Duration
	// Severity holds the numeric severity of the event for log events.
	Severity int64
}

func (e *Event) toModelProtobuf(out *modelpb.Event) {
	*out = modelpb.Event{
		Outcome:  e.Outcome,
		Action:   e.Action,
		Dataset:  e.Dataset,
		Kind:     e.Kind,
		Category: e.Category,
		Type:     e.Type,
		Severity: e.Severity,
	}
	if !isZero(e.Duration) {
		out.Duration = durationpb.New(e.Duration)
	}
	if !isZero(e.SuccessCount) {
		out.SuccessCount = &modelpb.SummaryMetric{
			Count: e.SuccessCount.Count,
			Sum:   e.SuccessCount.Sum,
		}
	}
}
