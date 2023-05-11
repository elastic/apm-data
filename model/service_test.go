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

var (
	version, environment   = "5.1.3", "staging"
	langName, langVersion  = "ecmascript", "8"
	rtName, rtVersion      = "node", "8.0.0"
	fwName, fwVersion      = "Express", "1.2.3"
	targetName, targetType = "testdb", "oracle"
)

func TestServiceTransform(t *testing.T) {
	serviceName, serviceNodeName := "myService", "abc"

	tests := []struct {
		Service Service
		Fields  any
	}{
		{
			Service: Service{},
			Fields:  nil,
		},
		{
			Service: Service{
				Name:        serviceName,
				Version:     version,
				Environment: environment,
				Language: Language{
					Name:    langName,
					Version: langVersion,
				},
				Runtime: Runtime{
					Name:    rtName,
					Version: rtVersion,
				},
				Framework: Framework{
					Name:    fwName,
					Version: fwVersion,
				},
				Node: ServiceNode{Name: serviceNodeName},
				Target: &ServiceTarget{
					Name: targetName,
					Type: targetType,
				},
			},
			Fields: map[string]any{
				"name":        "myService",
				"version":     "5.1.3",
				"environment": "staging",
				"language": map[string]any{
					"name":    "ecmascript",
					"version": "8",
				},
				"runtime": map[string]any{
					"name":    "node",
					"version": "8.0.0",
				},
				"framework": map[string]any{
					"name":    "Express",
					"version": "1.2.3",
				},
				"node": map[string]any{"name": serviceNodeName},
				"target": map[string]any{
					"name": "testdb",
					"type": "oracle",
				},
			},
		},
		{
			Service: Service{
				Name: serviceName,
				Node: ServiceNode{Name: serviceNodeName},
				Target: &ServiceTarget{
					Name: "",
					Type: "test",
				},
			},
			Fields: map[string]any{
				"name": "myService",
				"node": map[string]any{"name": serviceNodeName},
				"target": map[string]any{
					"type": "test",
				},
			},
		},
		{
			Service: Service{
				Name: serviceName,
				Node: ServiceNode{Name: serviceNodeName},
				Target: &ServiceTarget{
					Name: "test",
					Type: "",
				},
			},
			Fields: map[string]any{
				"name": "myService",
				"node": map[string]any{"name": serviceNodeName},
				"target": map[string]any{
					"name": "test",
					"type": "",
				},
			},
		},
	}

	for _, test := range tests {
		output := transformAPMEvent(APMEvent{Service: test.Service})
		assert.Equal(t, test.Fields, output["service"])
	}
}
