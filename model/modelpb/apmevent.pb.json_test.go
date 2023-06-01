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
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestAPMEventToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *APMEvent
		expected *modeljson.Document
	}{
		"empty": {
			proto:    &APMEvent{},
			expected: &modeljson.Document{},
		},
		"full": {
			proto: &APMEvent{
				NumericLabels: map[string]*NumericLabelValue{
					"foo": {
						Values: []float64{1, 2, 3},
						Value:  1,
						Global: true,
					},
				},
				Labels: map[string]*LabelValue{
					"bar": {
						Value:  "a",
						Values: []string{"a", "b", "c"},
						Global: true,
					},
				},
				Message: "message",
			},
			expected: &modeljson.Document{
				NumericLabels: map[string]modeljson.NumericLabel{
					"foo": {
						Values: []float64{1, 2, 3},
						Value:  1,
					},
				},
				Labels: map[string]modeljson.Label{
					"bar": {
						Value:  "a",
						Values: []string{"a", "b", "c"},
					},
				},
				Message: "message",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Document
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out,
				cmpopts.IgnoreFields(APMEvent{}, "Timestamp"),
				cmpopts.IgnoreFields(modeljson.Document{}, "Timestamp"),
			)
			require.Empty(t, diff)
		})
	}
}
