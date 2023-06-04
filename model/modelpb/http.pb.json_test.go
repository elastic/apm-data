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
	"testing"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestHTTPToModelJSON(t *testing.T) {
	headers, headersMap := randomStruct(t)
	headers2, headersMap2 := randomStruct(t)
	cookies, cookiesMap := randomStruct(t)
	envs, envsMap := randomStruct(t)
	tru := true

	testCases := map[string]struct {
		proto    *HTTP
		expected *modeljson.HTTP
	}{
		"empty": {
			proto:    &HTTP{},
			expected: &modeljson.HTTP{},
		},
		"request": {
			proto: &HTTP{
				Request: &HTTPRequest{
					Headers:  headers,
					Env:      envs,
					Cookies:  cookies,
					Id:       "id",
					Method:   "method",
					Referrer: "referrer",
				},
			},
			expected: &modeljson.HTTP{
				Request: &modeljson.HTTPRequest{
					Headers:  headersMap,
					Env:      envsMap,
					Cookies:  cookiesMap,
					ID:       "id",
					Method:   "method",
					Referrer: "referrer",
				},
			},
		},
		"response": {
			proto: &HTTP{
				Response: &HTTPResponse{
					Headers:         headers2,
					Finished:        &tru,
					HeadersSent:     &tru,
					TransferSize:    int64Ptr(1),
					EncodedBodySize: int64Ptr(2),
					DecodedBodySize: int64Ptr(3),
					StatusCode:      200,
				},
			},
			expected: &modeljson.HTTP{
				Response: &modeljson.HTTPResponse{
					Finished:        &tru,
					HeadersSent:     &tru,
					TransferSize:    int64Ptr(1),
					EncodedBodySize: int64Ptr(2),
					DecodedBodySize: int64Ptr(3),
					Headers:         headersMap2,
					StatusCode:      200,
				},
			},
		},
		"no pointers": {
			proto: &HTTP{
				Request: &HTTPRequest{
					Id:       "id",
					Method:   "method",
					Referrer: "referrer",
				},
				Response: &HTTPResponse{
					StatusCode: 200,
				},
				Version: "version",
			},
			expected: &modeljson.HTTP{
				Request: &modeljson.HTTPRequest{
					ID:       "id",
					Method:   "method",
					Referrer: "referrer",
				},
				Response: &modeljson.HTTPResponse{
					StatusCode: 200,
				},
				Version: "version",
			},
		},
		"full": {
			proto: &HTTP{
				Request: &HTTPRequest{
					Headers:  headers,
					Env:      envs,
					Cookies:  cookies,
					Id:       "id",
					Method:   "method",
					Referrer: "referrer",
				},
				Response: &HTTPResponse{
					Headers:         headers2,
					Finished:        &tru,
					HeadersSent:     &tru,
					TransferSize:    int64Ptr(1),
					EncodedBodySize: int64Ptr(2),
					DecodedBodySize: int64Ptr(3),
					StatusCode:      200,
				},
				Version: "version",
			},
			expected: &modeljson.HTTP{
				Request: &modeljson.HTTPRequest{
					Headers:  headersMap,
					Env:      envsMap,
					Cookies:  cookiesMap,
					ID:       "id",
					Method:   "method",
					Referrer: "referrer",
				},
				Response: &modeljson.HTTPResponse{
					Finished:        &tru,
					HeadersSent:     &tru,
					TransferSize:    int64Ptr(1),
					EncodedBodySize: int64Ptr(2),
					DecodedBodySize: int64Ptr(3),
					Headers:         headersMap2,
					StatusCode:      200,
				},
				Version: "version",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			out := modeljson.HTTP{
				Request:  &modeljson.HTTPRequest{},
				Response: &modeljson.HTTPResponse{},
			}
			tc.proto.toModelJSON(&out)
			diff := cmp.Diff(*tc.expected, out)
			require.Empty(t, diff)
		})
	}
}
