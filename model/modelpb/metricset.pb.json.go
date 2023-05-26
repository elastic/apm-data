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

func (me *Metricset) toModelJSON(out *modeljson.Metricset) {
	var samples []modeljson.MetricsetSample
	if n := len(me.Samples); n > 0 {
		samples = make([]modeljson.MetricsetSample, n)
		for i, sample := range me.Samples {
			samples[i] = modeljson.MetricsetSample{
				Name:  sample.Name,
				Type:  string(sample.Type),
				Unit:  sample.Unit,
				Value: sample.Value,
				Histogram: modeljson.Histogram{
					Values: sample.Histogram.Values,
					Counts: sample.Histogram.Counts,
				},
				Summary: modeljson.SummaryMetric{
					Count: sample.Summary.Count,
					Sum:   sample.Summary.Sum,
				},
			}
		}
	}
	*out = modeljson.Metricset{
		Name:     me.Name,
		Interval: me.Interval,
		Samples:  samples,
	}
}
