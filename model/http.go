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

package model

// HTTP holds information about an HTTP request and/or response.
type HTTP struct {
	Request  *HTTPRequest
	Response *HTTPResponse
	Version  string
}

// HTTPRequest holds information about an HTTP request.
type HTTPRequest struct {
	Body     interface{}
	Headers  map[string]any // Non-ECS field.
	Env      map[string]any // Non-ECS field.
	Cookies  map[string]any // Non-ECS field.
	ID       string
	Method   string
	Referrer string
}

// HTTPResponse holds information about an HTTP response.
type HTTPResponse struct {
	Headers         map[string]any // Non-ECS field.
	Finished        *bool          // Non-ECS field.
	HeadersSent     *bool          // Non-ECS field.
	TransferSize    *int           // Non-ECS field.
	EncodedBodySize *int           // Non-ECS field.
	DecodedBodySize *int           // Non-ECS field.
	StatusCode      int
}

func (h *HTTP) fields() map[string]any {
	var fields mapStr
	fields.maybeSetString("version", h.Version)
	if h.Request != nil {
		fields.maybeSetMapStr("request", h.Request.fields())
	}
	if h.Response != nil {
		fields.maybeSetMapStr("response", h.Response.fields())
	}
	return map[string]any(fields)
}

func (h *HTTPRequest) fields() map[string]any {
	var fields mapStr
	fields.maybeSetString("id", h.ID)
	fields.maybeSetString("method", h.Method)
	fields.maybeSetString("referrer", h.Referrer)
	fields.maybeSetMapStr("headers", h.Headers)
	fields.maybeSetMapStr("env", h.Env)
	fields.maybeSetMapStr("cookies", h.Cookies)
	if h.Body != nil {
		var body mapStr
		body.set("original", h.Body)
		fields.set("body", body)
	}
	return map[string]any(fields)
}

func (h *HTTPResponse) fields() map[string]any {
	var fields mapStr
	if h.StatusCode > 0 {
		fields.set("status_code", h.StatusCode)
	}
	fields.maybeSetMapStr("headers", h.Headers)
	fields.maybeSetBool("finished", h.Finished)
	fields.maybeSetBool("headers_sent", h.HeadersSent)
	fields.maybeSetIntptr("transfer_size", h.TransferSize)
	fields.maybeSetIntptr("encoded_body_size", h.EncodedBodySize)
	fields.maybeSetIntptr("decoded_body_size", h.DecodedBodySize)
	return map[string]any(fields)
}
