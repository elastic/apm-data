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

func TestSpanTransform(t *testing.T) {
	path := "test/path"
	hexID := "0147258369012345"
	subtype := "amqp"
	action := "publish"
	timestamp := time.Date(2019, 1, 3, 15, 17, 4, 908.596*1e6, time.FixedZone("+0100", 3600))
	instance, statement, dbType, user, rowsAffected := "db01", "select *", "sql", "jane", uint32(5)
	destServiceType, destServiceName, destServiceResource := "db", "elasticsearch", "elasticsearch"
	links := []SpanLink{{Span: Span{ID: "linked_span"}, Trace: Trace{ID: "linked_trace"}}}

	m := transformAPMEvent(APMEvent{
		Processor: SpanProcessor,
		Timestamp: timestamp,
		Span: &Span{
			ID:                  hexID,
			Name:                "myspan",
			Type:                "myspantype",
			Kind:                "CLIENT",
			Subtype:             subtype,
			Action:              action,
			RepresentativeCount: 5.6,
			Stacktrace:          Stacktrace{{AbsPath: path}},
			DB: &DB{
				Instance:     instance,
				Statement:    statement,
				Type:         dbType,
				UserName:     user,
				RowsAffected: &rowsAffected,
			},
			DestinationService: &DestinationService{
				Type:     destServiceType,
				Name:     destServiceName,
				Resource: destServiceResource,
			},
			Message:   &Message{QueueName: "users"},
			Composite: &Composite{Count: 10, Sum: 1.1, CompressionStrategy: "exact_match"},
			Links:     links,
		},
	})
	delete(m, "@timestamp")

	assert.Equal(t, map[string]any{
		"processor": map[string]any{"name": "transaction", "event": "span"},
		"span": map[string]any{
			"id":      hexID,
			"name":    "myspan",
			"kind":    "CLIENT",
			"type":    "myspantype",
			"subtype": subtype,
			"action":  action,
			"stacktrace": []any{
				map[string]any{
					"exclude_from_grouping": false,
					"abs_path":              path,
				},
			},
			"db": map[string]any{
				"instance":      instance,
				"statement":     statement,
				"type":          dbType,
				"user":          map[string]any{"name": user},
				"rows_affected": 5.0,
			},
			"destination": map[string]any{
				"service": map[string]any{
					"type":     destServiceType,
					"name":     destServiceName,
					"resource": destServiceResource,
				},
			},
			"message": map[string]any{"queue": map[string]any{"name": "users"}},
			"composite": map[string]any{
				"count":                10.0,
				"sum":                  map[string]any{"us": 1100.0},
				"compression_strategy": "exact_match",
			},
			"links": []any{
				map[string]any{
					"span":  map[string]any{"id": "linked_span"},
					"trace": map[string]any{"id": "linked_trace"},
				},
			},
			"representative_count": 5.6,
		},
		"timestamp": map[string]any{"us": float64(timestamp.UnixNano() / 1000)},
	}, m)
}

func TestSpanHTTPFields(t *testing.T) {
	m := transformAPMEvent(APMEvent{
		Processor: SpanProcessor,
		HTTP: HTTP{
			Version: "2.0",
			Request: &HTTPRequest{
				Method: "get",
			},
			Response: &HTTPResponse{
				StatusCode: 200,
			},
		},
		URL: URL{Original: "http://localhost"},
	})
	delete(m, "@timestamp")
	assert.Equal(t, map[string]any{
		"processor": map[string]any{
			"name":  "transaction",
			"event": "span",
		},
		"http": map[string]any{
			"version": "2.0",
			"request": map[string]any{
				"method": "get",
			},
			"response": map[string]any{
				"status_code": 200.0,
			},
		},
		"url": map[string]any{
			"original": "http://localhost",
		},
	}, m)
}
