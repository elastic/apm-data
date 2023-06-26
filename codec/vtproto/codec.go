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

package vtproto

// VTProtoMarshaler expects a struct to have MarshalVT() method that
// returns a byte slice representation the object
type VTProtoMarshaler interface {
	MarshalVT() ([]byte, error)
}

// VTProtoUnmarshaler expects a struct to have UnmarshalVT() method that
// parses a given byte slice and places the decoded results in the object.
type VTProtoUnmarshaler interface {
	UnmarshalVT([]byte) error
}

// Codec is a composite of Encoder and Decoder.
type Codec struct {
}

// Encode encodes vtprotoMessage type into byte slice
func (Codec) Encode(in VTProtoMarshaler) ([]byte, error) {
	return in.MarshalVT()
}

// Decode decodes a byte slice into vtprotoMessage type.
func (Codec) Decode(in []byte, out VTProtoUnmarshaler) error {
	return out.UnmarshalVT(in)
}
