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

func (e *Transaction) toModelJSON(out *modeljson.Transaction, metricset bool) {
	*out = modeljson.Transaction{
		ID:                  e.Id,
		Type:                e.Type,
		Name:                e.Name,
		Result:              e.Result,
		Sampled:             e.Sampled,
		Root:                e.Root,
		RepresentativeCount: e.RepresentativeCount,
	}

	if e.Custom != nil {
		m := kvToMap(e.Custom)
		updateFields(m)
		out.Custom = m
	}

	if n := len(e.Marks); n > 0 {
		marks := make(map[string]map[string]float64, n)
		for k, mark := range e.Marks {
			sanitizedMark := make(map[string]float64, len(mark.Measurements))
			for k, v := range mark.Measurements {
				sanitizedMark[sanitizeLabelKey(k)] = v
			}
			marks[sanitizeLabelKey(k)] = sanitizedMark
		}
		out.Marks = marks
	}
	if e.Message != nil {
		var message modeljson.Message
		e.Message.toModelJSON(&message)
		out.Message = &message
	}
	if e.UserExperience != nil {
		var userExperience modeljson.UserExperience
		e.UserExperience.toModelJSON(&userExperience)
		out.UserExperience = &userExperience
	}
	if metricset {
		// DroppedSpansStats is only indexed for metric documents, never for events.
		if n := len(e.DroppedSpansStats); n > 0 {
			droppedSpansStats := make([]modeljson.DroppedSpanStats, n)
			for i, dss := range e.DroppedSpansStats {
				dssJson := modeljson.DroppedSpanStats{
					DestinationServiceResource: dss.DestinationServiceResource,
					ServiceTargetType:          dss.ServiceTargetType,
					ServiceTargetName:          dss.ServiceTargetName,
					Outcome:                    dss.Outcome,
				}

				if dss.Duration != nil {
					dssJson.Duration = modeljson.AggregatedDuration{
						Count: int(dss.Duration.Count),
						Sum:   dss.Duration.Sum.AsDuration(),
					}
				}

				droppedSpansStats[i] = dssJson
			}
			out.DroppedSpansStats = droppedSpansStats
		}
	}

	if e.DurationHistogram != nil {
		out.DurationHistogram = modeljson.Histogram{
			Values: e.DurationHistogram.Values,
			Counts: e.DurationHistogram.Counts,
		}
	}
	if e.DurationSummary != nil {
		out.DurationSummary = modeljson.SummaryMetric{
			Count: e.DurationSummary.Count,
			Sum:   e.DurationSummary.Sum,
		}
	}
	if e.SpanCount != nil {
		out.SpanCount = modeljson.SpanCount{
			Dropped: e.SpanCount.Dropped,
			Started: e.SpanCount.Started,
		}
	}
}
