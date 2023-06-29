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
	"time"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestTransactionToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto               *Transaction
		expectedNoMetricset *modeljson.Transaction
		expectedMetricset   *modeljson.Transaction
	}{
		"empty": {
			proto:               &Transaction{},
			expectedNoMetricset: &modeljson.Transaction{},
			expectedMetricset:   &modeljson.Transaction{},
		},
		"no pointers": {
			proto: &Transaction{
				Type:                "type",
				Name:                "name",
				Result:              "result",
				Id:                  "id",
				RepresentativeCount: 8,
				Sampled:             true,
				Root:                true,
			},
			expectedNoMetricset: &modeljson.Transaction{
				Type:                "type",
				Name:                "name",
				Result:              "result",
				ID:                  "id",
				RepresentativeCount: 8,
				Sampled:             true,
				Root:                true,
			},
			expectedMetricset: &modeljson.Transaction{
				Type:                "type",
				Name:                "name",
				Result:              "result",
				ID:                  "id",
				RepresentativeCount: 8,
				Sampled:             true,
				Root:                true,
			},
		},
		"full": {
			proto: &Transaction{
				SpanCount: &SpanCount{
					Started: uintPtr(1),
					Dropped: uintPtr(2),
				},
				// TODO investigat valid values
				Custom: nil,
				Marks: map[string]*TransactionMark{
					"foo": {
						Measurements: map[string]float64{
							"bar": 3,
						},
					},
				},
				Type:   "type",
				Name:   "name",
				Result: "result",
				Id:     "id",
				DurationHistogram: &Histogram{
					Values: []float64{4},
					Counts: []int64{5},
				},
				DroppedSpansStats: []*DroppedSpanStats{
					{
						DestinationServiceResource: "destinationserviceresource",
						ServiceTargetType:          "servicetargetype",
						ServiceTargetName:          "servicetargetname",
						Outcome:                    "outcome",
						Duration: &AggregatedDuration{
							Count: 4,
							Sum:   durationpb.New(5 * time.Second),
						},
					},
				},
				DurationSummary: &SummaryMetric{
					Count: 6,
					Sum:   7,
				},
				RepresentativeCount: 8,
				Sampled:             true,
				Root:                true,
			},
			expectedNoMetricset: &modeljson.Transaction{
				SpanCount: modeljson.SpanCount{
					Started: uintPtr(1),
					Dropped: uintPtr(2),
				},
				// TODO investigat valid values
				Custom: nil,
				Marks: map[string]map[string]float64{
					"foo": {
						"bar": 3,
					},
				},
				Type:   "type",
				Name:   "name",
				Result: "result",
				ID:     "id",
				DurationHistogram: modeljson.Histogram{
					Values: []float64{4},
					Counts: []int64{5},
				},
				DurationSummary: modeljson.SummaryMetric{
					Count: 6,
					Sum:   7,
				},
				RepresentativeCount: 8,
				Sampled:             true,
				Root:                true,
			},
			expectedMetricset: &modeljson.Transaction{
				SpanCount: modeljson.SpanCount{
					Started: uintPtr(1),
					Dropped: uintPtr(2),
				},
				// TODO investigat valid values
				Custom: nil,
				Marks: map[string]map[string]float64{
					"foo": {
						"bar": 3,
					},
				},
				Type:   "type",
				Name:   "name",
				Result: "result",
				ID:     "id",
				DurationHistogram: modeljson.Histogram{
					Values: []float64{4},
					Counts: []int64{5},
				},
				DroppedSpansStats: []modeljson.DroppedSpanStats{
					{
						DestinationServiceResource: "destinationserviceresource",
						ServiceTargetType:          "servicetargetype",
						ServiceTargetName:          "servicetargetname",
						Outcome:                    "outcome",
						Duration: modeljson.AggregatedDuration{
							Count: 4,
							Sum:   5 * time.Second,
						},
					},
				},
				DurationSummary: modeljson.SummaryMetric{
					Count: 6,
					Sum:   7,
				},
				RepresentativeCount: 8,
				Sampled:             true,
				Root:                true,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Transaction
			tc.proto.toModelJSON(&out, false)
			require.Empty(t, cmp.Diff(*tc.expectedNoMetricset, out))

			var out2 modeljson.Transaction
			tc.proto.toModelJSON(&out2, true)
			require.Empty(t, cmp.Diff(*tc.expectedMetricset, out2))
		})
	}
}
