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
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientFields(t *testing.T) {
	for name, tc := range map[string]struct {
		domain string
		ip     netip.Addr
		port   int
		out    any
	}{
		"Empty":  {out: nil},
		"IPv4":   {ip: netip.MustParseAddr("192.0.0.1"), out: map[string]any{"ip": "192.0.0.1"}},
		"IPv6":   {ip: netip.MustParseAddr("2001:db8::68"), out: map[string]any{"ip": "2001:db8::68"}},
		"Port":   {port: 123, out: map[string]any{"port": 123.0}},
		"Domain": {domain: "testing.invalid", out: map[string]any{"domain": "testing.invalid"}},
	} {
		t.Run(name, func(t *testing.T) {
			out := transformAPMEvent(APMEvent{Client: Client{
				Domain: tc.domain,
				IP:     tc.ip,
				Port:   tc.port,
			}})
			assert.Equal(t, tc.out, out["client"])
		})
	}
}
