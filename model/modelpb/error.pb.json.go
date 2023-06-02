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
	out.ID = e.Id
	out.GroupingKey = e.GroupingKey
	out.Culprit = e.Culprit
	out.Message = e.Message
	out.Type = e.Type
	out.StackTrace = e.StackTrace
	if e.Custom != nil {
		out.Custom = customFields(e.Custom.AsMap())
	}
	if e.Exception != nil {
		e.Exception.toModelJSON(out.Exception)
	}
	if e.Log != nil {
		out.Log.Message = e.Log.Message
		out.Log.ParamMessage = e.Log.ParamMessage
		out.Log.LoggerName = e.Log.LoggerName
		out.Log.Level = e.Log.Level
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
	out.Message = e.Message
	out.Module = e.Module
	out.Code = e.Code
	out.Type = e.Type
	out.Handled = e.Handled
	if e.Attributes != nil && len(e.Attributes.AsMap()) != 0 {
		out.Attributes = e.Attributes.AsMap()
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