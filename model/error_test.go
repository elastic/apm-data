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
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseException() *Exception {
	return &Exception{Message: "exception message"}
}

func (e *Exception) withCode(code string) *Exception {
	e.Code = code
	return e
}

func baseLog() *ErrorLog {
	return &ErrorLog{Message: "error log message"}
}

func TestHandleExceptionTree(t *testing.T) {
	event := APMEvent{
		Error: &Error{
			Exception: &Exception{
				Message: "message0",
				Type:    "type0",
				Stacktrace: Stacktrace{{
					Filename: "file0",
				}},
				Cause: []Exception{{
					Message: "message1",
					Type:    "type1",
				}, {
					Message: "message2",
					Type:    "type2",
					Cause: []Exception{{
						Message: "message3",
						Type:    "type3",
						Cause: []Exception{{
							Message: "message4",
							Type:    "type4",
						}, {
							Message: "message5",
							Type:    "type5",
						}},
					}},
				}, {
					Message: "message6",
					Type:    "type6",
				}},
			},
		},
	}

	m := transformAPMEvent(event)
	assert.Equal(t, map[string]any{
		"exception": []any{
			map[string]any{
				"message": "message0",
				"type":    "type0",
				"stacktrace": []any{
					map[string]any{
						"filename":              "file0",
						"exclude_from_grouping": false,
					},
				},
			},
			map[string]any{"message": "message1", "type": "type1"},
			map[string]any{"message": "message2", "type": "type2", "parent": 0.0},
			map[string]any{"message": "message3", "type": "type3"},
			map[string]any{"message": "message4", "type": "type4"},
			map[string]any{"message": "message5", "type": "type5", "parent": 3.0},
			map[string]any{"message": "message6", "type": "type6", "parent": 0.0},
		},
	}, m["error"])
}

func TestErrorFields(t *testing.T) {
	id := "45678"
	culprit := "some trigger"

	errorType := "error type"
	module := "error module"
	exMsg := "exception message"
	handled := false
	attributes := map[string]any{"k1": "val1"}
	exception := Exception{
		Type:       errorType,
		Code:       "13",
		Message:    exMsg,
		Module:     module,
		Handled:    &handled,
		Attributes: attributes,
		Stacktrace: []*StacktraceFrame{{Filename: "st file"}},
	}

	level := "level"
	loggerName := "logger"
	logMsg := "error log message"
	paramMsg := "param message"
	log := ErrorLog{
		Level:        level,
		Message:      logMsg,
		ParamMessage: paramMsg,
		LoggerName:   loggerName,
	}

	tests := map[string]struct {
		Error  Error
		Output map[string]any
	}{
		"withGroupingKey": {
			Error:  Error{GroupingKey: "foo"},
			Output: map[string]any{"grouping_key": "foo"},
		},
		"withLog": {
			Error: Error{Log: baseLog()},
			Output: map[string]any{
				"log": map[string]any{"message": "error log message"},
			},
		},
		"withLogAndException": {
			Error: Error{Exception: baseException(), Log: baseLog()},
			Output: map[string]any{
				"exception": []any{
					map[string]any{
						"message": "exception message",
					},
				},
				"log": map[string]any{"message": "error log message"},
			},
		},
		"withException": {
			Error: Error{Exception: baseException()},
			Output: map[string]any{
				"exception": []any{
					map[string]any{"message": "exception message"},
				},
			},
		},
		"stringCode": {
			Error: Error{Exception: baseException().withCode("13")},
			Output: map[string]any{
				"exception": []any{
					map[string]any{"message": "exception message", "code": "13"},
				},
			},
		},
		"withStackTrace": {
			Error: Error{StackTrace: "raw stack trace"},
			Output: map[string]any{
				"stack_trace": "raw stack trace",
			},
		},
		"withFrames": {
			Error: Error{
				ID:        id,
				Culprit:   culprit,
				Exception: &exception,
				Log:       &log,
			},
			Output: map[string]any{
				"id":      "45678",
				"culprit": "some trigger",
				"exception": []any{
					map[string]any{
						"stacktrace": []any{
							map[string]any{
								"filename":              "st file",
								"exclude_from_grouping": false,
							},
						},
						"code":       "13",
						"message":    "exception message",
						"module":     "error module",
						"attributes": map[string]any{"k1": "val1"},
						"type":       "error type",
						"handled":    false,
					},
				},
				"log": map[string]any{
					"message":       "error log message",
					"param_message": "param message",
					"logger_name":   "logger",
					"level":         "level",
				},
			},
		},
		"withLogMessageAndType": {
			Error: Error{
				Message: "error log message",
				Type:    "IllegalArgumentException",
			},
			Output: map[string]any{
				"message": "error log message",
				"type":    "IllegalArgumentException",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := transformAPMEvent(APMEvent{Error: &tc.Error})
			assert.Equal(t, tc.Output, m["error"])
		})
	}
}
