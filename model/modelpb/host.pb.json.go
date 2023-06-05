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
	"net/netip"

	"github.com/elastic/apm-data/model/internal/modeljson"
)

func (h *Host) toModelJSON(out *modeljson.Host) {
	*out = modeljson.Host{
		Hostname:     h.Hostname,
		Name:         h.Name,
		ID:           h.Id,
		Architecture: h.Architecture,
		Type:         h.Type,
	}

	if n := len(h.Ip); n != 0 {
		ips := make([]string, 0, n)
		for _, ip := range h.Ip {
			if _, err := netip.ParseAddr(ip); err == nil {
				ips = append(ips, ip)
			}
		}
		out.IP = ips
	}

	if h.Os != nil {
		var os modeljson.OS
		h.Os.toModelJSON(&os)
		out.OS = &os
	}
}
