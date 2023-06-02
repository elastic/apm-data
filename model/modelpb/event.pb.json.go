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

func (e *Event) toModelJSON(out *modeljson.Event) {
	out.Outcome = e.Outcome
	out.Action = e.Action
	out.Dataset = e.Dataset
	out.Kind = e.Kind
	out.Category = e.Category
	out.Type = e.Type
	out.Duration = int64(e.Duration.AsDuration().Nanoseconds())
	out.Severity = e.Severity
	if e.SuccessCount != nil {
		out.SuccessCount = modeljson.SummaryMetric{
			Count: e.SuccessCount.Count,
			Sum:   e.SuccessCount.Sum,
		}
	}
}