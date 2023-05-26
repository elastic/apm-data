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

package modeljson

type HTTP struct {
	Request  *HTTPRequest  `json:"request,omitempty"`
	Response *HTTPResponse `json:"response,omitempty"`
	Version  string        `json:"version,omitempty"`
}

type HTTPRequest struct {
	Body     *HTTPRequestBody `json:"body,omitempty"`
	Headers  map[string]any   `json:"headers,omitempty"` // Non-ECS field.
	Env      map[string]any   `json:"env,omitempty"`     // Non-ECS field.
	Cookies  map[string]any   `json:"cookies,omitempty"` // Non-ECS field.
	ID       string           `json:"id,omitempty"`
	Method   string           `json:"method,omitempty"`
	Referrer string           `json:"referrer,omitempty"`
}

type HTTPRequestBody struct {
	Original any `json:"original,omitempty"` // Non-ECS field.
}

type HTTPResponse struct {
	Finished        *bool          `json:"finished,omitempty"`          // Non-ECS field.
	HeadersSent     *bool          `json:"headers_sent,omitempty"`      // Non-ECS field.
	TransferSize    *int64         `json:"transfer_size,omitempty"`     // Non-ECS field.
	EncodedBodySize *int64         `json:"encoded_body_size,omitempty"` // Non-ECS field.
	DecodedBodySize *int64         `json:"decoded_body_size,omitempty"` // Non-ECS field.
	Headers         map[string]any `json:"headers,omitempty"`           // Non-ECS field.
	StatusCode      int            `json:"status_code,omitempty"`
}
