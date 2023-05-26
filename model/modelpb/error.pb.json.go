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

func (e *Error) toModelJSON(out *modeljson.Error) {
	*out = modeljson.Error{
		ID:          e.Id,
		GroupingKey: e.GroupingKey,
		Culprit:     e.Culprit,
		Message:     e.Message,
		Type:        e.Type,
		StackTrace:  e.StackTrace,
		Custom:      customFields(e.Custom.AsMap()),
	}
	if e.Exception != nil {
		out.Exception = &modeljson.Exception{}
		e.Exception.toModelJSON(out.Exception)
	}
	if e.Log != nil {
		out.Log = &modeljson.ErrorLog{
			Message:      e.Log.Message,
			ParamMessage: e.Log.ParamMessage,
			LoggerName:   e.Log.LoggerName,
			Level:        e.Log.Level,
		}
		if n := len(e.Log.Stacktrace); n > 0 {
			out.Log.Stacktrace = make([]modeljson.StacktraceFrame, n)
			for i, frame := range e.Log.Stacktrace {
				if frame != nil {
					frame.toModelJSON(&out.Log.Stacktrace[i])
				}
			}
		}
	}
}

func (e *Exception) toModelJSON(out *modeljson.Exception) {
	*out = modeljson.Exception{
		Message:    e.Message,
		Module:     e.Module,
		Code:       e.Code,
		Attributes: e.Attributes,
		Type:       e.Type,
		Handled:    e.Handled,
	}
	if n := len(e.Cause); n > 0 {
		out.Cause = make([]modeljson.Exception, n)
		for i, cause := range e.Cause {
			if cause != nil {
				cause.toModelJSON(&out.Cause[i])
			}
		}
	}
	if n := len(e.Stacktrace); n > 0 {
		out.Stacktrace = make([]modeljson.StacktraceFrame, n)
		for i, frame := range e.Stacktrace {
			if frame != nil {
				frame.toModelJSON(&out.Stacktrace[i])
			}
		}
	}
}
