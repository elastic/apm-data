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
	"net/netip"
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestSourceToModelJSON(t *testing.T) {
	testCases := map[string]struct {
		proto    *Source
		expected *modeljson.Source
	}{
		"empty": {
			proto:    &Source{},
			expected: &modeljson.Source{},
		},
		"full": {
			proto: &Source{
				Ip: MustParseIP("127.0.0.1"),
				Nat: &NAT{
					Ip: MustParseIP("127.0.0.2"),
				},
				Domain: "domain",
				Port:   443,
			},
			expected: &modeljson.Source{
				IP: modeljson.IP(netip.MustParseAddr("127.0.0.1")),
				NAT: modeljson.NAT{
					IP: modeljson.IP(netip.MustParseAddr("127.0.0.2")),
				},
				Domain: "domain",
				Port:   443,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var out modeljson.Source
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out,
				cmp.Comparer(func(a modeljson.IP, b modeljson.IP) bool {
					return netip.Addr(a).Compare(netip.Addr(b)) == 0
				}),
			)
			require.Empty(t, diff)
		})
	}
}
