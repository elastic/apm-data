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
	"time"

	"github.com/elastic/apm-data/model/modelpb"
)

var (
	// MetricsetProcessor is the Processor value that should be assigned to metricset events.
	MetricsetProcessor = Processor{Name: "metric", Event: "metric"}
)

// MetricType describes the type of a metric: gauge, counter, or histogram.
type MetricType string

// Valid MetricType values.
const (
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeCounter   MetricType = "counter"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeSummary   MetricType = "summary"
)

var metricTypeText = map[string]modelpb.MetricType{
	"gauge":     modelpb.MetricType_METRIC_TYPE_GAUGE,
	"counter":   modelpb.MetricType_METRIC_TYPE_COUNTER,
	"histogram": modelpb.MetricType_METRIC_TYPE_HISTOGRAM,
	"summary":   modelpb.MetricType_METRIC_TYPE_SUMMARY,
}

// Metricset describes a set of metrics and associated metadata.
type Metricset struct {
	// Name holds an optional name for the metricset.
	Name string
	// Interval holds the time period the metricset samples represent.
	//
	// The value is formatted in seconds, or minutes with the unit suffixed.
	// Some examples include: `10s`, `1m`, `60m`.
	Interval string
	// Samples holds the metrics in the set.
	Samples []MetricsetSample
	// DocCount holds the document count for pre-aggregated metrics.
	//
	// See https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping-doc-count-field.html
	DocCount int64
}

// MetricsetSample represents a single named metric.
type MetricsetSample struct {
	// Type holds an optional metric type.
	//
	// If Type is unspecified or invalid, it will be ignored.
	Type MetricType
	// Name holds the metric name.
	//
	// Name is required.
	Name string
	// Unit holds an optional unit:
	//
	// - "percent" (value is in the range [0,1])
	// - "byte"
	// - a time unit: "nanos", "micros", "ms", "s", "m", "h", "d"
	//
	// If Unit is unspecified or invalid, it will be ignored.
	Unit string
	// Histogram holds bucket values and counts for histogram metrics.
	Histogram
	// SummaryMetric holds a combined count and sum of aggregated
	// measurements.
	SummaryMetric
	// Value holds the metric value for single-value metrics.
	//
	// If Counts and Values are specified, then Value will be ignored.
	Value float64
}

// Histogram holds bucket values and counts for a histogram metric.
type Histogram struct {
	// Values holds the bucket values for histogram metrics.
	//
	// These values must be provided in ascending order.
	Values []float64

	// Counts holds the bucket counts for histogram metrics.
	//
	// These numbers must be positive or zero.
	//
	// If Counts is specified, then Values is expected to be
	// specified with the same number of elements, and with the
	// same order.
	Counts []int64
}

// SummaryMetric holds summary metrics (count and sum).
type SummaryMetric struct {
	// Count holds the number of aggregated measurements.
	Count int64

	// Sum holds the sum of aggregated measurements.
	Sum float64
}

// AggregatedDuration holds a count and sum of aggregated durations.
type AggregatedDuration struct {
	// Count holds the number of durations aggregated.
	Count int

	// Sum holds the sum of aggregated durations.
	Sum time.Duration
}

func (me *Metricset) toModelProtobuf(out *modelpb.Metricset) {
	var samples []*modelpb.MetricsetSample
	if n := len(me.Samples); n > 0 {
		samples = make([]*modelpb.MetricsetSample, n)
		for i, sample := range me.Samples {
			ms := &modelpb.MetricsetSample{
				Name:  sample.Name,
				Type:  metricTypeText[string(sample.Type)],
				Unit:  sample.Unit,
				Value: sample.Value,
			}

			if len(sample.Histogram.Values) != 0 || len(sample.Histogram.Counts) != 0 {
				ms.Histogram = &modelpb.Histogram{
					Values: sample.Histogram.Values,
					Counts: sample.Histogram.Counts,
				}
			}

			if !isZero(sample.SummaryMetric) {
				ms.Summary = &modelpb.SummaryMetric{
					Count: sample.SummaryMetric.Count,
					Sum:   sample.SummaryMetric.Sum,
				}
			}

			samples[i] = ms
		}
	}
	*out = modelpb.Metricset{
		Name:     me.Name,
		Interval: me.Interval,
		Samples:  samples,
		DocCount: me.DocCount,
	}
}
