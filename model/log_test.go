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

	"github.com/stretchr/testify/assert"
)

func TestLogTransform(t *testing.T) {
	tests := []struct {
		Log    Log
		Output any
	}{
		{
			Log:    Log{},
			Output: nil,
		},
		{
			Log: Log{
				Level:  "warn",
				Logger: "bootstrap",
			},
			Output: map[string]any{
				"level":  "warn",
				"logger": "bootstrap",
			},
		},
		{
			Log: Log{
				Level:  "warn",
				Logger: "bootstrap",
				Origin: LogOrigin{
					File: LogOriginFile{
						Name: "testFile",
						Line: 12,
					},
					FunctionName: "testFunc",
				},
			},
			Output: map[string]any{
				"level":  "warn",
				"logger": "bootstrap",
				"origin": map[string]any{
					"function": "testFunc",
					"file": map[string]any{
						"name": "testFile",
						"line": 12.0,
					},
				},
			},
		},
	}

	for _, test := range tests {
		output := transformAPMEvent(APMEvent{Log: test.Log})
		assert.Equal(t, test.Output, output["log"])
	}
}
