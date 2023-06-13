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
	"github.com/elastic/apm-data/model/modelpb"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	// ErrorProcessor is the Processor value that should be assigned to error events.
	ErrorProcessor = Processor{Name: "error", Event: "error"}
)

type Error struct {
	Custom      map[string]any
	Exception   *Exception
	Log         *ErrorLog
	ID          string
	GroupingKey string
	Culprit     string
	// StackTrace holds an unparsed stack trace.
	//
	// This may be set when a stack trace cannot be parsed.
	StackTrace string
	// Message holds an error message.
	//
	// Message is the ECS field equivalent of the APM field `error.log.message`.
	Message string
	// Type holds the type of the error.
	Type string
}

type Exception struct {
	Message    string
	Module     string
	Code       string
	Attributes interface{}
	Stacktrace Stacktrace
	Type       string
	Handled    *bool
	Cause      []Exception
}

type ErrorLog struct {
	Message      string
	Level        string
	ParamMessage string
	LoggerName   string
	Stacktrace   Stacktrace
}

func (e *Error) toModelProtobuf(out *modelpb.Error) {
	*out = modelpb.Error{
		Id:          e.ID,
		GroupingKey: e.GroupingKey,
		Culprit:     e.Culprit,
		Message:     e.Message,
		Type:        e.Type,
		StackTrace:  e.StackTrace,
	}
	if len(e.Custom) != 0 {
		if cf, err := structpb.NewStruct(e.Custom); err == nil {
			out.Custom = cf
		}
	}
	if e.Exception != nil {
		out.Exception = &modelpb.Exception{}
		e.Exception.toModelProtobuf(out.Exception)
	}
	if e.Log != nil {
		out.Log = &modelpb.ErrorLog{
			Message:      e.Log.Message,
			ParamMessage: e.Log.ParamMessage,
			LoggerName:   e.Log.LoggerName,
			Level:        e.Log.Level,
		}
		if n := len(e.Log.Stacktrace); n > 0 {
			out.Log.Stacktrace = make([]*modelpb.StacktraceFrame, n)
			for i, frame := range e.Log.Stacktrace {
				s := modelpb.StacktraceFrame{}
				frame.toModelProtobuf(&s)
				out.Log.Stacktrace[i] = &s
			}
		}
	}
}

func (e *Exception) toModelProtobuf(out *modelpb.Exception) {
	*out = modelpb.Exception{
		Message: e.Message,
		Module:  e.Module,
		Code:    e.Code,
		Type:    e.Type,
		Handled: e.Handled,
	}
	if m, ok := e.Attributes.(map[string]any); ok && len(m) != 0 {
		if a, err := structpb.NewStruct(m); err == nil {
			out.Attributes = a
		}
	}
	if n := len(e.Cause); n > 0 {
		out.Cause = make([]*modelpb.Exception, n)
		for i, cause := range e.Cause {
			outCause := modelpb.Exception{}
			cause.toModelProtobuf(&outCause)
			out.Cause[i] = &outCause
		}
	}
	if n := len(e.Stacktrace); n > 0 {
		out.Stacktrace = make([]*modelpb.StacktraceFrame, n)
		for i, frame := range e.Stacktrace {
			outFrame := modelpb.StacktraceFrame{}
			frame.toModelProtobuf(&outFrame)
			out.Stacktrace[i] = &outFrame
		}
	}
}
