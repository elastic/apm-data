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

// Reslice increases the slice's capacity, if necessary, to guarantee space for n elements.
// If specified, the newFn function is used to populate the extra (optionally) allocated space.
// The method returns a slice with len(slice)==n.
func Reslice[Slice ~[]model, model any](slice Slice, n int, newFn func() model) Slice {
	// check if there is enough capacity to hold n elements
	if diff := n - cap(slice); diff > 0 {
		// start of the extra space
		idx := cap(slice)

		// Grow the slice
		// Note: append gives no guarantee on the capacity of the resulting slice
		// and might overallocate as long as there's enough space for n elements.
		slice = append([]model(slice)[:cap(slice)], make([]model, diff)...)
		if newFn != nil {
			// extend the slice to its capacity
			// Note: reset the slice length to n before returning
			slice = slice[:cap(slice)]

			// populate the extra space
			for ; idx < len(slice); idx++ {
				slice[idx] = newFn()
			}
		}
	}

	// slice has enough capacity to hold n elements and there is no need to grow the slice.
	// Do not make assumption on the length of the slice and always return a slice with len(slice)==n.
	// Most of the time len(slice)==0 because the argument is a newly reset slice from a
	// protobuf model.
	return slice[:n]
}
