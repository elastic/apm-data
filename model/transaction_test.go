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

func TestTransactionTransform(t *testing.T) {
	id := "123"
	dropped, startedSpans := 5, 14
	name := "mytransaction"
	duration := 65980 * time.Microsecond

	tests := []struct {
		Processor   Processor
		Transaction Transaction
		Output      map[string]any
		Msg         string
	}{
		{
			Transaction: Transaction{
				ID:   id,
				Type: "tx",
			},
			Output: map[string]any{
				"id":   id,
				"type": "tx",
			},
			Msg: "SpanCount empty",
		},
		{
			Transaction: Transaction{
				ID:        id,
				Type:      "tx",
				SpanCount: SpanCount{Started: &startedSpans},
			},
			Output: map[string]any{
				"id":         id,
				"type":       "tx",
				"span_count": map[string]any{"started": 14.0},
			},
			Msg: "SpanCount only contains `started`",
		},
		{
			Transaction: Transaction{
				ID:        id,
				Type:      "tx",
				SpanCount: SpanCount{Dropped: &dropped},
			},
			Output: map[string]any{
				"id":         id,
				"type":       "tx",
				"span_count": map[string]any{"dropped": 5.0},
			},
			Msg: "SpanCount only contains `dropped`",
		},
		{
			Transaction: Transaction{
				ID:                  id,
				Name:                name,
				Type:                "tx",
				Result:              "tx result",
				Sampled:             true,
				SpanCount:           SpanCount{Started: &startedSpans, Dropped: &dropped},
				Root:                true,
				RepresentativeCount: 1.23,
			},
			Output: map[string]any{
				"id":                   id,
				"name":                 "mytransaction",
				"type":                 "tx",
				"result":               "tx result",
				"span_count":           map[string]any{"started": 14.0, "dropped": 5.0},
				"sampled":              true,
				"root":                 true,
				"representative_count": 1.23,
			},
			Msg: "Full Event",
		},
		{
			Transaction: Transaction{
				ID:        id,
				Name:      name,
				Type:      "tx",
				Result:    "tx result",
				Sampled:   true,
				SpanCount: SpanCount{Started: &startedSpans, Dropped: &dropped},
				DroppedSpansStats: []DroppedSpanStats{
					{
						DestinationServiceResource: "mysql://server:3306",
						Outcome:                    "success",
						Duration: AggregatedDuration{
							Count: 5,
							Sum:   duration,
						},
					},
					{
						DestinationServiceResource: "http://elasticsearch:9200",
						Outcome:                    "unknown",
						Duration: AggregatedDuration{
							Count: 15,
							Sum:   duration,
						},
					},
				},
				Root: true,
			},
			Output: map[string]any{
				"id":         id,
				"name":       "mytransaction",
				"type":       "tx",
				"result":     "tx result",
				"span_count": map[string]any{"started": 14.0, "dropped": 5.0},
				"sampled":    true,
				"root":       true,
			},
			Msg: "Full Event With Dropped Spans Statistics",
		},
		{
			Processor: MetricsetProcessor,
			Transaction: Transaction{
				ID:        id,
				Name:      name,
				Type:      "tx",
				Result:    "tx result",
				Sampled:   true,
				SpanCount: SpanCount{Started: &startedSpans, Dropped: &dropped},
				DroppedSpansStats: []DroppedSpanStats{
					{
						DestinationServiceResource: "mysql://server:3306",
						Outcome:                    "success",
						Duration: AggregatedDuration{
							Count: 5,
							Sum:   duration,
						},
					},
					{
						DestinationServiceResource: "http://elasticsearch:9200",
						Outcome:                    "unknown",
						Duration: AggregatedDuration{
							Count: 15,
							Sum:   duration,
						},
					},
				},
				Root: true,
			},
			Output: map[string]any{
				"id":         id,
				"name":       "mytransaction",
				"type":       "tx",
				"result":     "tx result",
				"span_count": map[string]any{"started": 14.0, "dropped": 5.0},
				"dropped_spans_stats": []any{
					map[string]any{
						"destination_service_resource": "mysql://server:3306",
						"duration":                     map[string]any{"count": 5.0, "sum.us": 65980.0},
						"outcome":                      "success",
					},
					map[string]any{
						"destination_service_resource": "http://elasticsearch:9200",
						"duration":                     map[string]any{"count": 15.0, "sum.us": 65980.0},
						"outcome":                      "unknown",
					},
				},
				"sampled": true,
				"root":    true,
			},
			Msg: "Full Event With Dropped Spans Statistics using the MetricsetProcessor processor",
		},
	}

	for idx, test := range tests {
		event := APMEvent{
			Processor:   TransactionProcessor,
			Transaction: &test.Transaction,
		}
		emptyProc := Processor{}
		if test.Processor != emptyProc {
			event.Processor = test.Processor
		}
		m := transformAPMEvent(event)
		assert.Equal(t, test.Output, m["transaction"], fmt.Sprintf("Failed at idx %v; %s", idx, test.Msg))
	}
}

func TestEventsTransformWithMetadata(t *testing.T) {
	hostname := "a.b.c"
	architecture := "darwin"
	platform := "x64"
	id, name, userAgent := "123", "jane", "node-js-2.3"
	url := "https://localhost"
	serviceName, serviceNodeName, serviceVersion := "myservice", "service-123", "2.1.3"

	m := transformAPMEvent(APMEvent{
		Processor: TransactionProcessor,
		Service: Service{
			Name:    serviceName,
			Version: serviceVersion,
			Node:    ServiceNode{Name: serviceNodeName},
		},
		Host: Host{
			Name:         name,
			Hostname:     hostname,
			Architecture: architecture,
			OS:           OS{Platform: platform},
		},
		User:      User{ID: id, Name: name},
		UserAgent: UserAgent{Original: userAgent},
		URL:       URL{Original: url},
		Transaction: &Transaction{
			Custom:  map[string]any{"foo.bar": "baz"},
			Message: &Message{QueueName: "routeUser"},
			Sampled: true,
		},
	})
	delete(m, "@timestamp")

	assert.Equal(t, map[string]any{
		"processor":  map[string]any{"name": "transaction", "event": "transaction"},
		"user":       map[string]any{"id": "123", "name": "jane"},
		"user_agent": map[string]any{"original": userAgent},
		"host": map[string]any{
			"architecture": "darwin",
			"hostname":     "a.b.c",
			"name":         "jane",
			"os": map[string]any{
				"platform": "x64",
			},
		},
		"service": map[string]any{
			"name":    serviceName,
			"version": serviceVersion,
			"node":    map[string]any{"name": serviceNodeName},
		},
		"transaction": map[string]any{
			"sampled": true,
			"custom": map[string]any{
				"foo_bar": "baz",
			},
			"message": map[string]any{"queue": map[string]any{"name": "routeUser"}},
		},
		"url": map[string]any{
			"original": url,
		},
	}, m)
}

func TestTransactionTransformMarks(t *testing.T) {
	out := marshalJSONAPMEvent(APMEvent{
		Timestamp: time.Time{},
		Transaction: &Transaction{
			Marks: TransactionMarks{
				"a.b": TransactionMark{
					"c.d": 123,
				},
			},
		},
	})

	assert.JSONEq(t, `{"@timestamp":"0001-01-01T00:00:00.000Z","transaction":{
		"marks": {
		  "a_b": {
		    "c_d": 123
		  }
		}
	}}`, string(out))
}
