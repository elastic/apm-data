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

import (
	"github.com/elastic/apm-data/model/modelpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func ToKv(m map[string]any) []*modelpb.KeyValue {
	m = normalizeMap(m)
	if len(m) == 0 {
		return nil
	}

	out := make([]*modelpb.KeyValue, len(m))

	i := 0
	for k, v := range m {
		value, _ := structpb.NewValue(v)
		out[i] = &modelpb.KeyValue{
			Key:   k,
			Value: value,
		}
		i++
	}

	return out
}
