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
	"encoding/binary"
	"net/netip"
)

// ParseIP turns a string IP address into a valid proto IP object
func ParseIP(s string) (*IP, error) {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return nil, err
	}

	return Addr2IP(addr), nil
}

func MustParseIP(s string) *IP {
	ip, err := ParseIP(s)
	if err != nil {
		panic(err)
	}

	return ip
}

func Addr2IP(addr netip.Addr) *IP {
	if addr.Is4() {
		return &IP{
			V4: binary.BigEndian.Uint32(addr.AsSlice()),
		}
	}

	return &IP{
		V6: addr.AsSlice(),
	}
}

func IP2Addr(i *IP) netip.Addr {
	var ip []byte

	if i.GetV6() != nil {
		ip = i.V6
	}

	if i.GetV4() != 0 {
		ip = make([]byte, 4)
		binary.BigEndian.PutUint32(ip, i.V4)
	}

	if addr, ok := netip.AddrFromSlice(ip); ok {
		return addr
	}

	return netip.Addr{}
}

func IP2String(i *IP) string {
	var ip []byte

	if i.GetV6() != nil {
		ip = i.V6
	}

	if i.GetV4() != 0 {
		ip = make([]byte, 4)
		binary.BigEndian.PutUint32(ip, i.V4)
	}

	if addr, ok := netip.AddrFromSlice(ip); ok {
		return addr.String()
	}

	return ""
}
