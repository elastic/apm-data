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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStacktraceTransform(t *testing.T) {
	originalLineno := 111
	originalColno := 222
	originalFunction := "original function"
	originalFilename := "original filename"
	originalModule := "original module"
	originalClassname := "original classname"
	originalAbsPath := "original path"

	mappedLineno := 333
	mappedColno := 444
	mappedFunction := "mapped function"
	mappedFilename := "mapped filename"
	mappedClassname := "mapped classname"
	mappedAbsPath := "mapped path"

	contextLine := "context line"

	tests := []struct {
		Stacktrace Stacktrace
		Output     any
		Msg        string
	}{
		{
			Stacktrace: Stacktrace{},
			Output:     nil,
			Msg:        "Empty Stacktrace",
		},
		{
			Stacktrace: Stacktrace{&StacktraceFrame{}},
			Output:     []any{map[string]any{"exclude_from_grouping": false}},
			Msg:        "Stacktrace with empty Frame",
		},
		{
			Stacktrace: Stacktrace{{
				Colno:        &originalColno,
				Lineno:       &originalLineno,
				Filename:     originalFilename,
				Function:     originalFunction,
				Classname:    originalClassname,
				Module:       originalModule,
				AbsPath:      originalAbsPath,
				LibraryFrame: true,
				Vars:         map[string]any{"a": "abc", "b": 123},
			}},
			Output: []any{
				map[string]any{
					"abs_path":  "original path",
					"filename":  "original filename",
					"function":  "original function",
					"classname": "original classname",
					"module":    "original module",
					"line": map[string]any{
						"number": 111.0,
						"column": 222.0,
					},
					"exclude_from_grouping": false,
					"library_frame":         true,
					"vars":                  map[string]any{"a": "abc", "b": 123.0},
				},
			},
			Msg: "unmapped stacktrace",
		},
		{
			Stacktrace: Stacktrace{{
				Colno:     &mappedColno,
				Lineno:    &mappedLineno,
				Filename:  mappedFilename,
				Function:  mappedFunction,
				Classname: mappedClassname,
				AbsPath:   mappedAbsPath,
				Original: Original{
					Colno:     &originalColno,
					Lineno:    &originalLineno,
					Filename:  originalFilename,
					Function:  originalFunction,
					Classname: originalClassname,
					AbsPath:   originalAbsPath,
				},
				ExcludeFromGrouping: true,
				SourcemapUpdated:    true,
				SourcemapError:      "boom",
				ContextLine:         contextLine,
				PreContext:          []string{"before1", "before2"},
				PostContext:         []string{"after1", "after2"},
			}},
			Output: []any{
				map[string]any{
					"abs_path":  "mapped path",
					"filename":  "mapped filename",
					"function":  "mapped function",
					"classname": "mapped classname",
					"line": map[string]any{
						"number":  333.0,
						"column":  444.0,
						"context": "context line",
					},
					"context": map[string]any{
						"pre":  []any{"before1", "before2"},
						"post": []any{"after1", "after2"},
					},
					"original": map[string]any{
						"abs_path":  "original path",
						"filename":  "original filename",
						"function":  "original function",
						"classname": "original classname",
						"lineno":    111.0,
						"colno":     222.0,
					},
					"exclude_from_grouping": true,
					"sourcemap": map[string]any{
						"updated": true,
						"error":   "boom",
					},
				},
			},
			Msg: "mapped stacktrace",
		},
	}

	for idx, test := range tests {
		output := transformAPMEvent(APMEvent{Span: &Span{Stacktrace: test.Stacktrace}})
		span := output["span"].(map[string]any)
		assert.Equal(t, test.Output, span["stacktrace"], fmt.Sprintf("Failed at idx %v; %s", idx, test.Msg))
	}
}
