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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessaging_Fields(t *testing.T) {
	ageMillis := 1577958057123
	out := transformAPMEvent(APMEvent{
		Transaction: &Transaction{
			Message: &Message{
				QueueName:  "orders",
				RoutingKey: "a_routing_key",
				Body:       "order confirmed",
				Headers:    http.Header{"Internal": []string{"false"}, "Services": []string{"user", "order"}},
				AgeMillis:  &ageMillis,
			},
		},
	})

	transaction := out["transaction"].(map[string]any)
	assert.Equal(t, map[string]any{
		"queue":       map[string]any{"name": "orders"},
		"routing_key": "a_routing_key",
		"body":        "order confirmed",
		"headers": map[string]any{
			"Internal": []any{"false"},
			"Services": []any{"user", "order"},
		},
		"age": map[string]any{
			"ms": 1577958057123.0,
		},
	}, transaction["message"])
}
