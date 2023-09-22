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
		ip := IPFromVTPool()
		ip.V4 = binary.BigEndian.Uint32(addr.AsSlice())
		return ip
	}

	ip := IPFromVTPool()
	ip.V6 = addr.AsSlice()
	return ip
}

func IP2Addr(i *IP) netip.Addr {
	if addr := i.GetV6(); len(addr) == 16 {
		return netip.AddrFrom16([16]byte(addr))
	}
	if i.GetV4() != 0 {
		var addr [4]byte
		binary.BigEndian.PutUint32(addr[:], i.V4)
		return netip.AddrFrom4(addr)
	}
	return netip.Addr{}
}

func IP2String(i *IP) string {
	addr := IP2Addr(i)
	if addr.IsValid() {
		return addr.String()
	}
	return ""
}
