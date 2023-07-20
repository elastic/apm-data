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

package modeljson

import (
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestServiceToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Service
		expected *modeljson.Service
	}{
		"empty": {
			proto:    &Service{},
			expected: &modeljson.Service{},
		},
		"no pointers": {
			proto: &Service{
				Name:        "name",
				Version:     "version",
				Environment: "environment",
			},
			expected: &modeljson.Service{
				Name:        "name",
				Version:     "version",
				Environment: "environment",
			},
		},
		"full": {
			proto: &Service{
				Origin: &ServiceOrigin{
					Id:      "origin_id",
					Name:    "origin_name",
					Version: "origin_version",
				},
				Target: &ServiceTarget{
					Name: "target_name",
					Type: "target_type",
				},
				Language: &Language{
					Name:    "language_name",
					Version: "language_version",
				},
				Runtime: &Runtime{
					Name:    "runtime_name",
					Version: "runtime_version",
				},
				Framework: &Framework{
					Name:    "framework_name",
					Version: "framework_version",
				},
				Name:        "name",
				Version:     "version",
				Environment: "environment",
				Node: &ServiceNode{
					Name: "node_name",
				},
			},
			expected: &modeljson.Service{
				Origin: &modeljson.ServiceOrigin{
					ID:      "origin_id",
					Name:    "origin_name",
					Version: "origin_version",
				},
				Target: &modeljson.ServiceTarget{
					Name: "target_name",
					Type: "target_type",
				},
				Language: &modeljson.Language{
					Name:    "language_name",
					Version: "language_version",
				},
				Runtime: &modeljson.Runtime{
					Name:    "runtime_name",
					Version: "runtime_version",
				},
				Framework: &modeljson.Framework{
					Name:    "framework_name",
					Version: "framework_version",
				},
				Name:        "name",
				Version:     "version",
				Environment: "environment",
				Node: &modeljson.ServiceNode{
					Name: "node_name",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			out := modeljson.Service{
				Node:      &modeljson.ServiceNode{},
				Language:  &modeljson.Language{},
				Runtime:   &modeljson.Runtime{},
				Framework: &modeljson.Framework{},
				Origin:    &modeljson.ServiceOrigin{},
				Target:    &modeljson.ServiceTarget{},
			}
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
