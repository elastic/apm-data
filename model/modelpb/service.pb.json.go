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

func (s *Service) toModelJSON(out *modeljson.Service) {
	out.Name = s.Name
	out.Version = s.Version
	out.Environment = s.Environment
	if s.Node != nil {
		out.Node.Name = s.Node.Name
	} else {
		out.Node = nil
	}
	if s.Language != nil {
		out.Language.Name = s.Language.Name
		out.Language.Version = s.Language.Version
	} else {
		out.Language = nil
	}
	if s.Runtime != nil {
		out.Runtime.Name = s.Runtime.Name
		out.Runtime.Version = s.Runtime.Version
	} else {
		out.Runtime = nil
	}
	if s.Framework != nil {
		out.Framework.Name = s.Framework.Name
		out.Framework.Version = s.Framework.Version
	} else {
		out.Framework = nil
	}
	if s.Origin != nil {
		out.Origin.ID = s.Origin.Id
		out.Origin.Name = s.Origin.Name
		out.Origin.Version = s.Origin.Version
	} else {
		out.Origin = nil
	}
	if s.Target != nil {
		out.Target.Name = s.Target.Name
		out.Target.Type = s.Target.Type
	} else {
		out.Target = nil
	}
}
