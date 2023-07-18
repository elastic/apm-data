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

package common

import structpb "google.golang.org/protobuf/types/known/structpb"

func KvToMap(h []*KeyValue) map[string]any {
	if len(h) == 0 {
		return nil
	}
	m := make(map[string]any, len(h))
	for _, v := range h {
		switch v.Value.GetKind().(type) {
		case *structpb.Value_NullValue:
			m[v.Key] = nil
		case *structpb.Value_NumberValue:
			m[v.Key] = v.Value.GetNumberValue()
		case *structpb.Value_StringValue:
			m[v.Key] = v.Value.GetStringValue()
		case *structpb.Value_BoolValue:
			m[v.Key] = v.Value.GetBoolValue()
		case *structpb.Value_StructValue:
			m[v.Key] = v.Value.GetStructValue()
		case *structpb.Value_ListValue:
			m[v.Key] = v.Value.GetListValue()
		}
	}
	return m
}
