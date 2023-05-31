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

import "github.com/elastic/apm-data/model/internal/modeljson"

type Stacktrace []*StacktraceFrame

type StacktraceFrame struct {
	Vars                map[string]any
	Lineno              *uint32
	Colno               *uint32
	Filename            string
	Classname           string
	ContextLine         string
	Module              string
	Function            string
	AbsPath             string
	SourcemapError      string
	Original            Original
	PreContext          []string
	PostContext         []string
	LibraryFrame        bool
	SourcemapUpdated    bool
	ExcludeFromGrouping bool
}

type Original struct {
	AbsPath      string
	Filename     string
	Classname    string
	Lineno       *uint32
	Colno        *uint32
	Function     string
	LibraryFrame bool
}

func (s *StacktraceFrame) toModelJSON(out *modeljson.StacktraceFrame) {
	*out = modeljson.StacktraceFrame{
		Filename:            s.Filename,
		Classname:           s.Classname,
		AbsPath:             s.AbsPath,
		Module:              s.Module,
		Function:            s.Function,
		LibraryFrame:        s.LibraryFrame,
		ExcludeFromGrouping: s.ExcludeFromGrouping,
	}

	if len(s.Vars) != 0 {
		out.Vars = s.Vars
	}

	if len(s.PreContext) != 0 || len(s.PostContext) != 0 {
		out.Context = &modeljson.StacktraceFrameContext{
			Pre:  s.PreContext,
			Post: s.PostContext,
		}
	}

	if s.Lineno != nil || s.Colno != nil || s.ContextLine != "" {
		out.Line = &modeljson.StacktraceFrameLine{
			Number:  s.Lineno,
			Column:  s.Colno,
			Context: s.ContextLine,
		}
	}

	sourcemap := modeljson.StacktraceFrameSourcemap{
		Updated: s.SourcemapUpdated,
		Error:   s.SourcemapError,
	}
	if sourcemap != (modeljson.StacktraceFrameSourcemap{}) {
		out.Sourcemap = &sourcemap
	}

	orig := modeljson.StacktraceFrameOriginal{LibraryFrame: s.Original.LibraryFrame}
	if s.SourcemapUpdated {
		orig.Filename = s.Original.Filename
		orig.Classname = s.Original.Classname
		orig.AbsPath = s.Original.AbsPath
		orig.Function = s.Original.Function
		orig.Colno = s.Original.Colno
		orig.Lineno = s.Original.Lineno
	}
	if orig != (modeljson.StacktraceFrameOriginal{}) {
		out.Original = &orig
	}
}
