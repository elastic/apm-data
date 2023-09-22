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

package modeldecoderutil

func Reslice[Slice ~[]model, model any](slice Slice, want int, newFn func() model) Slice {
	if diff := want - cap(slice); diff > 0 {
		extra := make([]model, diff)
		if newFn != nil {
			for i := range extra {
				extra[i] = newFn()
			}
		}
		slice = append([]model(slice)[:cap(slice)], extra...)
	}
	return slice[:want]
}
