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

func (s *StacktraceFrame) toModelJSON(out *modeljson.StacktraceFrame) {
	out.Filename = s.Filename
	out.Classname = s.Classname
	out.AbsPath = s.AbsPath
	out.Module = s.Module
	out.Function = s.Function
	out.LibraryFrame = s.LibraryFrame
	out.ExcludeFromGrouping = s.ExcludeFromGrouping

	if len(s.Vars.AsMap()) != 0 {
		out.Vars = s.Vars.AsMap()
	}

	if len(s.PreContext) != 0 || len(s.PostContext) != 0 {
		if out.Context == nil {
			out.Context = &modeljson.StacktraceFrameContext{}
		}
		out.Context.Pre = s.PreContext
		out.Context.Post = s.PostContext
	}

	if s.Lineno != nil || s.Colno != nil || s.ContextLine != "" {
		if out.Line == nil {
			out.Line = &modeljson.StacktraceFrameLine{}
		}
		out.Line.Number = s.Lineno
		out.Line.Column = s.Colno
		out.Line.Context = s.ContextLine
	}

	sourcemap := modeljson.StacktraceFrameSourcemap{
		Updated: s.SourcemapUpdated,
		Error:   s.SourcemapError,
	}
	if sourcemap != (modeljson.StacktraceFrameSourcemap{}) {
		if out.Sourcemap == nil {
			out.Sourcemap = &modeljson.StacktraceFrameSourcemap{}
		}
		out.Sourcemap.Updated = s.SourcemapUpdated
		out.Sourcemap.Error = s.SourcemapError
	}

	if s.Original != nil {
		if out.Original == nil {
			out.Original = &modeljson.StacktraceFrameOriginal{}
		}
		out.Original.LibraryFrame = s.Original.LibraryFrame
		if s.SourcemapUpdated {
			out.Original.Filename = s.Original.Filename
			out.Original.Classname = s.Original.Classname
			out.Original.AbsPath = s.Original.AbsPath
			out.Original.Function = s.Original.Function
			out.Original.Colno = s.Original.Colno
			out.Original.Lineno = s.Original.Lineno
		}
	}
}
