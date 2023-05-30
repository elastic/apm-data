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

func (s *StacktraceFrame) toModelProtobuf(out *modelpb.StacktraceFrame) {
	*out = modelpb.StacktraceFrame{
		Lineno:              s.Lineno,
		Colno:               s.Colno,
		Filename:            s.Filename,
		Classname:           s.Classname,
		ContextLine:         s.ContextLine,
		Module:              s.Module,
		Function:            s.Function,
		AbsPath:             s.AbsPath,
		SourcemapError:      s.SourcemapError,
		PreContext:          s.PreContext,
		PostContext:         s.PostContext,
		LibraryFrame:        s.LibraryFrame,
		SourcemapUpdated:    s.SourcemapUpdated,
		ExcludeFromGrouping: s.ExcludeFromGrouping,
	}

	if !isZero(s.Original) {
		out.Original = &modelpb.Original{
			AbsPath:      s.Original.AbsPath,
			Filename:     s.Original.Filename,
			Classname:    s.Original.Classname,
			Lineno:       s.Original.Lineno,
			Colno:        s.Original.Colno,
			Function:     s.Original.Function,
			LibraryFrame: s.Original.LibraryFrame,
		}
	}

	if len(s.Vars) != 0 {
		if v, err := structpb.NewStruct(s.Vars); err == nil {
			out.Vars = v
		}
	}
}
