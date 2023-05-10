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

func TestDeviceFields(t *testing.T) {
	tests := []struct {
		Device Device
		Output any
	}{
		{
			Device: Device{},
			Output: nil,
		},
		{
			Device: Device{
				ID: "E733F41E-DF47-4BB4-AAF0-FD784FD95653",
				Model: DeviceModel{
					Name:       "Samsung Galaxy S6",
					Identifier: "SM-G920F",
				},
				Manufacturer: "Samsung",
			},
			Output: map[string]any{
				"id": "E733F41E-DF47-4BB4-AAF0-FD784FD95653",
				"model": map[string]any{
					"name":       "Samsung Galaxy S6",
					"identifier": "SM-G920F",
				},
				"manufacturer": "Samsung",
			},
		},
	}

	for _, test := range tests {
		output := transformAPMEvent(APMEvent{Device: test.Device})
		assert.Equal(t, test.Output, output["device"])
	}
}
