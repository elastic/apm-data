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
	"github.com/tidwall/sjson"
)

func TestHost(t *testing.T) {
	out := marshalJSONAPMEvent(APMEvent{
		Host: Host{
			Hostname:     "host_hostname",
			Name:         "host_name",
			Architecture: "amd",
			Type:         "t2.medium",
			IP: []netip.Addr{
				netip.MustParseAddr("127.0.0.1"),
				netip.MustParseAddr("::1"),
			},
			OS: OS{
				Name:     "macOS",
				Version:  "10.14.6",
				Platform: "osx",
				Full:     "Mac OS Mojave",
				Type:     "macos",
			},
		},
	})
	out, _ = sjson.DeleteBytes(out, "\\@timestamp")

	assert.JSONEq(t, `{"host":{
		"hostname": "host_hostname",
		"name": "host_name",
		"architecture": "amd",
		"type": "t2.medium",
		"ip": ["127.0.0.1", "::1"],
		"os": {
		  "name": "macOS",
		  "version": "10.14.6",
		  "platform": "osx",
		  "full": "Mac OS Mojave",
		  "type": "macos"
		}
	}}`, string(out))
}
