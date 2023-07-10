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

	"github.com/stretchr/testify/assert"
)

func TestParseIP(t *testing.T) {
	for _, tt := range []struct {
		name    string
		address string

		expectedIP       *IP
		expectedErr      error
		expectedStringIP string
	}{
		{
			name:    "with a valid IPv4 address",
			address: "127.0.0.1",

			expectedIP: &IP{
				IpAddr: &IP_V4{
					V4: 2130706433,
				},
			},
			expectedStringIP: "127.0.0.1",
		},
		{
			name:    "with a valid IPv6 address",
			address: "::1",

			expectedIP: &IP{
				IpAddr: &IP_V6{
					V6: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				},
			},
			expectedStringIP: "::1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := ParseIP(tt.address)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedIP, ip)
			assert.Equal(t, tt.expectedStringIP, IP2String(ip))
		})
	}
}
