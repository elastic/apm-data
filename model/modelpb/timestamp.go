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

import "time"

// TimeToPBTimestamp encodes a time.Time to Unix epoch nanos in uint64 for protobuf.
func TimeToPBTimestamp(t time.Time) uint64 {
	return uint64(t.UnixNano())
}

// PBTimestampToTime decodes a uint64 of Unix epoch nanos to a time.Time for protobuf.
func PBTimestampToTime(timestamp uint64) time.Time {
	return time.Unix(0, int64(timestamp)).UTC()
}

func PBTimestampNow() uint64 {
	return TimeToPBTimestamp(time.Now())
}

func PBTimestampTruncate(t uint64, d time.Duration) uint64 {
	ns := uint64(d.Nanoseconds())
	return t - t%ns
}

func PBTimestampAdd(t uint64, d time.Duration) uint64 {
	ns := uint64(d.Nanoseconds())
	return t + ns
}

func PBTimestampSub(t uint64, tt time.Time) time.Duration {
	b := TimeToPBTimestamp(tt)
	return time.Duration(t - b)
}
