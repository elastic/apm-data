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
	"github.com/stretchr/testify/require"
)

func TestMessageToModelJSON(t *testing.T) {
	millis := uint64(1)

	testCases := map[string]struct {
		proto    *Message
		expected *modeljson.Message
	}{
		"empty": {
			proto:    &Message{},
			expected: &modeljson.Message{},
		},
		"full": {
			proto: &Message{
				Body: "body",
				Headers: []*HTTPHeader{
					{
						Key:   "foo",
						Value: []string{"bar"},
					},
				},
				AgeMillis:  &millis,
				QueueName:  "queuename",
				RoutingKey: "routingkey",
			},
			expected: &modeljson.Message{
				Body: "body",
				Headers: map[string][]string{
					"foo": {"bar"},
				},
				Age: modeljson.MessageAge{
					Millis: &millis,
				},
				Queue: modeljson.MessageQueue{
					Name: "queuename",
				},
				RoutingKey: "routingkey",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Message
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
