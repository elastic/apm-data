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

func TestObserverOmitempty(t *testing.T) {
	output := transformAPMEvent(APMEvent{})
	assert.NotContains(t, output, "observer")
}

func TestObserverFields(t *testing.T) {
	output := transformAPMEvent(APMEvent{
		Observer: Observer{
			Hostname: "observer_hostname",
			Name:     "observer_name",
			Type:     "observer_type",
			Version:  "observer_version",
		},
	})

	assert.Equal(t, map[string]any{
		"hostname": "observer_hostname",
		"name":     "observer_name",
		"type":     "observer_type",
		"version":  "observer_version",
	}, output["observer"])
}
