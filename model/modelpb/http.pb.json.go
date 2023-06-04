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

import "github.com/elastic/apm-data/model/internal/modeljson"

func (h *HTTP) toModelJSON(out *modeljson.HTTP) {
	out.Version = h.Version
	if h.Request != nil {
		out.Request.ID = h.Request.Id
		out.Request.Method = h.Request.Method
		out.Request.Referrer = h.Request.Referrer
		if len(h.Request.Headers.AsMap()) != 0 {
			out.Request.Headers = h.Request.Headers.AsMap()
		}
		if len(h.Request.Env.AsMap()) != 0 {
			out.Request.Env = h.Request.Env.AsMap()
		}
		if len(h.Request.Cookies.AsMap()) != 0 {
			out.Request.Cookies = h.Request.Cookies.AsMap()
		}
		if h.Request.Body != nil {
			out.Request.Body.Original = h.Request.Body
		}
	} else {
		out.Request = nil
	}
	if h.Response != nil {
		out.Response.StatusCode = int(h.Response.StatusCode)
		out.Response.Finished = h.Response.Finished
		out.Response.HeadersSent = h.Response.HeadersSent
		out.Response.TransferSize = h.Response.TransferSize
		out.Response.EncodedBodySize = h.Response.EncodedBodySize
		out.Response.DecodedBodySize = h.Response.DecodedBodySize
		if len(h.Response.Headers.AsMap()) != 0 {
			out.Response.Headers = h.Response.Headers.AsMap()
		}
	} else {
		out.Response = nil
	}
}
