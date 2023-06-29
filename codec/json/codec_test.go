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

package json

import (
	"testing"

	"github.com/elastic/apm-data/codec/internal/testutils"
	"github.com/elastic/apm-data/model/modelpb"
)

func BenchmarkCodec(b *testing.B) {
	codec := Codec{}
	b.Run("Encode", func(b *testing.B) {
		var output int64
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				by, _ := codec.Encode(testutils.FullEvent(b))
				output += int64(len(by))
			}
		})

		bytePerOp := float64(output) / float64(b.N)
		b.ReportMetric(bytePerOp, "bytes/op")
	})

	b.Run("Decode", func(b *testing.B) {
		encoded, _ := codec.Encode(testutils.FullEvent(b))
		var ev *modelpb.APMEvent
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				codec.Decode(encoded, ev)
			}
		})
	})
}
