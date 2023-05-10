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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricset(t *testing.T) {
	tests := []struct {
		Metricset *Metricset
		Output    map[string]any
		Msg       string
	}{
		{
			Metricset: &Metricset{Name: "raj"},
			Output: map[string]any{
				"metricset": map[string]any{"name": "raj"},
			},
			Msg: "Payload with metricset name.",
		},
		{
			Metricset: &Metricset{
				Samples: []MetricsetSample{
					{Name: "a.counter", Value: 612},
					{Name: "some.gauge", Value: 9.16},
				},
			},
			Output: map[string]any{
				"metricset": map[string]any{
					"samples": []any{
						map[string]any{
							"name":  "a.counter",
							"value": 612.0,
						},
						map[string]any{
							"name":  "some.gauge",
							"value": 9.16,
						},
					},
				},
			},
			Msg: "Payload with valid metric.",
		},
		{
			Metricset: &Metricset{
				DocCount: 6,
			},
			Output: map[string]any{
				"_doc_count": 6.0,
				"metricset":  map[string]any{},
			},
			Msg: "_doc_count",
		},
		{
			Metricset: &Metricset{
				Samples: []MetricsetSample{
					{
						Name: "latency_histogram",
						Type: "histogram",
						Unit: "s",
						Histogram: Histogram{
							Counts: []int64{1, 2, 3},
							Values: []float64{1.1, 2.2, 3.3},
						},
					},
					{
						Name: "request_summary",
						Type: "summary",
						SummaryMetric: SummaryMetric{
							Count: 10,
							Sum:   123.456,
						},
					},
					{
						Name:  "just_type",
						Type:  "counter",
						Value: 123,
					},
					{
						Name:  "just_unit",
						Unit:  "percent",
						Value: 0.99,
					},
				},
			},
			Output: map[string]any{
				"metricset": map[string]any{
					"samples": []any{
						map[string]any{
							"name":   "latency_histogram",
							"type":   "histogram",
							"unit":   "s",
							"counts": []any{1.0, 2.0, 3.0},
							"values": []any{1.1, 2.2, 3.3},
						},
						map[string]any{
							"name":        "request_summary",
							"type":        "summary",
							"sum":         123.456,
							"value_count": 10.0,
						},
						map[string]any{
							"name":  "just_type",
							"type":  "counter",
							"value": 123.0,
						},
						map[string]any{
							"name":  "just_unit",
							"unit":  "percent",
							"value": 0.99,
						},
					},
				},
			},
			Msg: "Payload with metric type and unit.",
		},
	}

	for idx, test := range tests {
		m := transformAPMEvent(APMEvent{Metricset: test.Metricset})
		delete(m, "@timestamp")
		assert.Equal(t, test.Output, m, fmt.Sprintf("Failed at idx %v; %s", idx, test.Msg))
	}
}

func TestTransformMetricsetTransaction(t *testing.T) {
	m := transformAPMEvent(APMEvent{
		Processor: MetricsetProcessor,
		Transaction: &Transaction{
			Name:   "transaction_name",
			Type:   "transaction_type",
			Result: "transaction_result",
			DurationHistogram: Histogram{
				Counts: []int64{1, 2, 3},
				Values: []float64{4.5, 6.0, 9.0},
			},
		},
		Metricset: &Metricset{Name: "transaction"},
	})
	assert.Equal(t, map[string]any{
		"@timestamp": "0001-01-01T00:00:00.000Z",
		"processor":  map[string]any{"name": "metric", "event": "metric"},
		"metricset":  map[string]any{"name": "transaction"},
		"transaction": map[string]any{
			"name":   "transaction_name",
			"type":   "transaction_type",
			"result": "transaction_result",
			"duration.histogram": map[string]any{
				"counts": []any{1.0, 2.0, 3.0},
				"values": []any{4.5, 6.0, 9.0},
			},
		},
	}, m)
}

func TestTransformMetricsetSpan(t *testing.T) {
	m := transformAPMEvent(APMEvent{
		Processor: MetricsetProcessor,
		Span: &Span{
			Type:    "span_type",
			Subtype: "span_subtype",
			SelfTime: AggregatedDuration{
				Count: 123,
				Sum:   time.Millisecond,
			},
			DestinationService: &DestinationService{
				Resource: "destination_service_resource",
				ResponseTime: AggregatedDuration{
					Count: 456,
					Sum:   time.Second,
				},
			},
		},
		Metricset: &Metricset{Name: "span"},
	})
	assert.Equal(t, map[string]any{
		"@timestamp": "0001-01-01T00:00:00.000Z",
		"processor":  map[string]any{"name": "metric", "event": "metric"},
		"metricset":  map[string]any{"name": "span"},
		"span": map[string]any{
			"type":    "span_type",
			"subtype": "span_subtype",
			"self_time": map[string]any{
				"count":  123.0,
				"sum.us": 1000.0,
			},
			"destination": map[string]any{
				"service": map[string]any{
					"resource": "destination_service_resource",
					"response_time": map[string]any{
						"count":  456.0,
						"sum.us": 1000000.0,
					},
				},
			},
		},
	}, m)
}
