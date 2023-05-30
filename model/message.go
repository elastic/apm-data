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

import (
	"net/http"

	"github.com/elastic/apm-data/model/modelpb"
)

// Message holds information about a recorded message, such as the message body and meta information
type Message struct {
	Body       string
	Headers    http.Header
	AgeMillis  *int64
	QueueName  string
	RoutingKey string
}

func (m *Message) toModelProtobuf(out *modelpb.Message) {
	var headers map[string]*modelpb.HTTPHeaderValue
	if n := len(m.Headers); n > 0 {
		headers = make(map[string]*modelpb.HTTPHeaderValue, n)
		for k, v := range m.Headers {
			headers[k] = &modelpb.HTTPHeaderValue{
				Values: v,
			}
		}
	}
	*out = modelpb.Message{
		Body:       m.Body,
		Headers:    headers,
		AgeMillis:  m.AgeMillis,
		QueueName:  m.QueueName,
		RoutingKey: m.RoutingKey,
	}
}
