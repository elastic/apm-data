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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventFieldsEmpty(t *testing.T) {
	m := transformAPMEvent(APMEvent{})
	delete(m, "@timestamp")
	assert.Empty(t, m)
}

func TestEventFields(t *testing.T) {
	tests := map[string]struct {
		Event  Event
		Output map[string]any
	}{
		"withOutcome": {
			Event: Event{Outcome: "success"},
			Output: map[string]any{
				"outcome": "success",
			},
		},
		"withAction": {
			Event: Event{Action: "process-started"},
			Output: map[string]any{
				"action": "process-started",
			},
		},
		"withDataset": {
			Event: Event{Dataset: "access-log"},
			Output: map[string]any{
				"dataset": "access-log",
			},
		},
		"withSeverity": {
			Event: Event{Severity: 1},
			Output: map[string]any{
				"severity": 1.0,
			},
		},
		"withDuration": {
			Event: Event{Duration: time.Minute},
			Output: map[string]any{
				"duration": 60000000000.0,
			},
		},
		"withOutcomeActionDataset": {
			Event: Event{Outcome: "success", Action: "process-started", Dataset: "access-log"},
			Output: map[string]any{
				"outcome": "success",
				"action":  "process-started",
				"dataset": "access-log",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := transformAPMEvent(APMEvent{Event: tc.Event})
			assert.Equal(t, tc.Output, m["event"])
		})
	}
}
