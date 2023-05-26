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
	var marks map[string]map[string]float64
	if n := len(e.Marks); n > 0 {
		marks = make(map[string]map[string]float64, n)
		for k, mark := range e.Marks {
			sanitizedMark := make(map[string]float64, len(mark.Measurements))
			for k, v := range mark.Measurements {
				sanitizedMark[sanitizeLabelKey(k)] = v
			}
			marks[sanitizeLabelKey(k)] = sanitizedMark
		}
	}
	var message *modeljson.Message
	if e.Message != nil {
		message = &modeljson.Message{}
		e.Message.toModelJSON(message)
	}
	var userExperience *modeljson.UserExperience
	if e.UserExperience != nil {
		userExperience = &modeljson.UserExperience{
			CumulativeLayoutShift: e.UserExperience.CumulativeLayoutShift,
			FirstInputDelay:       e.UserExperience.FirstInputDelay,
			TotalBlockingTime:     e.UserExperience.TotalBlockingTime,
			Longtask: modeljson.LongtaskMetrics{
				Count: int(e.UserExperience.LongTask.Count),
				Sum:   e.UserExperience.LongTask.Sum,
				Max:   e.UserExperience.LongTask.Max,
			},
		}
	}
	var droppedSpansStats []modeljson.DroppedSpanStats
	if metricset {
		// DroppedSpansStats is only indexed for metric documents, never for events.
		if n := len(e.DroppedSpansStats); n > 0 {
			droppedSpansStats = make([]modeljson.DroppedSpanStats, n)
			for i, dss := range e.DroppedSpansStats {
				droppedSpansStats[i] = modeljson.DroppedSpanStats{
					DestinationServiceResource: dss.DestinationServiceResource,
					ServiceTargetType:          dss.ServiceTargetType,
					ServiceTargetName:          dss.ServiceTargetName,
					Outcome:                    dss.Outcome,
					Duration: modeljson.AggregatedDuration{
						Count: int(dss.Duration.Count),
						Sum:   dss.Duration.Sum.AsDuration(),
					},
				}
			}
		}
	}
	*out = modeljson.Transaction{
		ID:                  e.Id,
		Type:                e.Type,
		Name:                e.Name,
		Result:              e.Result,
		Sampled:             e.Sampled,
		Root:                e.Root,
		RepresentativeCount: e.RepresentativeCount,

		DurationHistogram: modeljson.Histogram{
			Values: e.DurationHistogram.Values,
			Counts: e.DurationHistogram.Counts,
		},
		DurationSummary: modeljson.SummaryMetric{
			Count: e.DurationSummary.Count,
			Sum:   e.DurationSummary.Sum,
		},
		SpanCount: modeljson.SpanCount{
			Dropped: e.SpanCount.Dropped,
			Started: e.SpanCount.Started,
		},
		DroppedSpansStats: droppedSpansStats,

		Marks:          marks,
		Custom:         customFields(e.Custom.AsMap()),
		Message:        message,
		UserExperience: userExperience,
	}
}
