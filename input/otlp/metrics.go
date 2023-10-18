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

// Portions copied from OpenTelemetry Collector (contrib), from the
// elastic exporter.
//
// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlp

import (
	"context"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/elastic/apm-data/model/modelpb"
)

// ConsumeMetricsResult contains the number of rejected data points and error message for partial success response.
type ConsumeMetricsResult struct {
	ErrorMessage       string
	RejectedDataPoints int64
}

// ConsumeMetrics calls ConsumeMetricsWithResult but ignores the result.
// It exists to satisfy the go.opentelemetry.io/collector/consumer.Metrics interface.
func (c *Consumer) ConsumeMetrics(ctx context.Context, metrics pmetric.Metrics) error {
	_, err := c.ConsumeMetricsWithResult(ctx, metrics)
	return err
}

// ConsumeMetricsWithResult consumes OpenTelemetry metrics data, converting into
// the Elastic APM metrics model and sending to the reporter.
func (c *Consumer) ConsumeMetricsWithResult(ctx context.Context, metrics pmetric.Metrics) (ConsumeMetricsResult, error) {
	totalDataPoints := int64(metrics.DataPointCount())
	totalMetrics := int64(metrics.MetricCount())
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return ConsumeMetricsResult{}, err
	}
	defer c.sem.Release(1)

	remainingDataPoints := totalDataPoints
	remainingMetrics := totalMetrics
	receiveTimestamp := time.Now()
	c.config.Logger.Debug("consuming metrics", zap.Stringer("metrics", metricsStringer(metrics)))
	batch := c.convertMetrics(metrics, receiveTimestamp, &remainingDataPoints, &remainingMetrics)
	if remainingMetrics > 0 {
		// Some metrics remained after conversion, meaning that they were dropped.
		atomic.AddInt64(&c.stats.unsupportedMetricsDropped, remainingMetrics)
	}
	if err := c.config.Processor.ProcessBatch(ctx, batch); err != nil {
		return ConsumeMetricsResult{}, err
	}
	var errMsg string
	if remainingDataPoints > 0 {
		errMsg = "unsupported data points"
	}
	return ConsumeMetricsResult{
		RejectedDataPoints: remainingDataPoints,
		ErrorMessage:       errMsg,
	}, nil
}

func (c *Consumer) convertMetrics(
	metrics pmetric.Metrics,
	receiveTimestamp time.Time,
	remainingDataPoints, remainingMetrics *int64,
) (batch *modelpb.Batch) {
	batch = &modelpb.Batch{}
	resourceMetrics := metrics.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		c.convertResourceMetrics(resourceMetrics.At(i), receiveTimestamp, batch, remainingDataPoints, remainingMetrics)
	}
	return
}

func (c *Consumer) convertResourceMetrics(
	resourceMetrics pmetric.ResourceMetrics,
	receiveTimestamp time.Time,
	out *modelpb.Batch,
	remainingDataPoints, remainingMetrics *int64,
) (
	droppedDataPoints, droppedMetrics int64) {
	baseEvent := modelpb.APMEventFromVTPool()
	baseEvent.Event = modelpb.EventFromVTPool()
	baseEvent.Event.Received = modelpb.FromTime(receiveTimestamp)

	var timeDelta time.Duration
	resource := resourceMetrics.Resource()
	translateResourceMetadata(resource, baseEvent)
	if exportTimestamp, ok := exportTimestamp(resource); ok {
		timeDelta = receiveTimestamp.Sub(exportTimestamp)
	}
	scopeMetrics := resourceMetrics.ScopeMetrics()
	for i := 0; i < scopeMetrics.Len(); i++ {
		c.convertScopeMetrics(scopeMetrics.At(i), baseEvent, timeDelta, out, remainingDataPoints, remainingMetrics)
	}
	return
}

func (c *Consumer) convertScopeMetrics(
	in pmetric.ScopeMetrics,
	baseEvent *modelpb.APMEvent,
	timeDelta time.Duration,
	out *modelpb.Batch,
	remainingDataPoints, remainingMetrics *int64,
) {
	ms := make(metricsets)
	otelMetrics := in.Metrics()
	for i := 0; i < otelMetrics.Len(); i++ {
		c.addMetric(otelMetrics.At(i), ms, remainingDataPoints, remainingMetrics)
	}
	for key, ms := range ms {
		event := baseEvent.CloneVT()
		event.Timestamp = modelpb.FromTime(key.timestamp.Add(timeDelta))
		metrs := make([]*modelpb.MetricsetSample, 0, len(ms.samples))
		for _, s := range ms.samples {
			metrs = append(metrs, s)
		}
		event.Metricset = modelpb.MetricsetFromVTPool()
		event.Metricset.Samples = metrs
		event.Metricset.Name = "app"
		if ms.attributes.Len() > 0 {
			initEventLabels(event)
			ms.attributes.Range(func(k string, v pcommon.Value) bool {
				setLabel(k, event, ifaceAttributeValue(v))
				return true
			})
			if len(event.Labels) == 0 {
				event.Labels = nil
			}
			if len(event.NumericLabels) == 0 {
				event.NumericLabels = nil
			}
		}
		*out = append(*out, event)
	}
}

func (c *Consumer) addMetric(metric pmetric.Metric, ms metricsets, remainingDataPoints, remainingMetrics *int64) {
	var anyDropped bool
	// TODO(axw) support units
	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		dps := metric.Gauge().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			if sample, ok := numberSample(dp, modelpb.MetricType_METRIC_TYPE_GAUGE); ok {
				sample.Name = metric.Name()
				ms.upsert(dp.Timestamp().AsTime(), dp.Attributes(), sample)
				*remainingDataPoints--
			} else {
				anyDropped = true
			}
		}
	case pmetric.MetricTypeSum:
		dps := metric.Sum().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			if sample, ok := numberSample(dp, modelpb.MetricType_METRIC_TYPE_COUNTER); ok {
				sample.Name = metric.Name()
				ms.upsert(dp.Timestamp().AsTime(), dp.Attributes(), sample)
				*remainingDataPoints--
			} else {
				anyDropped = true
			}
		}
	case pmetric.MetricTypeHistogram:
		dps := metric.Histogram().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			if sample, ok := histogramSample(dp.BucketCounts(), dp.ExplicitBounds()); ok {
				sample.Name = metric.Name()
				ms.upsert(dp.Timestamp().AsTime(), dp.Attributes(), sample)
				*remainingDataPoints--
			} else {
				anyDropped = true
			}
		}

	case pmetric.MetricTypeSummary:
		dps := metric.Summary().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			sample := summarySample(dp)
			sample.Name = metric.Name()
			ms.upsert(dp.Timestamp().AsTime(), dp.Attributes(), sample)
			*remainingDataPoints--
		}
	default:
		// Unsupported metric:
		// It will be recorded as dropped as remainingDataPoints and remainingMetrics are not decreased
		anyDropped = true
	}
	if !anyDropped {
		*remainingMetrics--
	}
}

func numberSample(dp pmetric.NumberDataPoint, metricType modelpb.MetricType) (*modelpb.MetricsetSample, bool) {
	var value float64
	switch dp.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		value = float64(dp.IntValue())
	case pmetric.NumberDataPointValueTypeDouble:
		value = dp.DoubleValue()
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return nil, false
		}
	default:
		return nil, false
	}
	ms := modelpb.MetricsetSampleFromVTPool()
	ms.Type = metricType
	ms.Value = value
	return ms, true
}

func summarySample(dp pmetric.SummaryDataPoint) *modelpb.MetricsetSample {
	ms := modelpb.MetricsetSampleFromVTPool()
	ms.Type = modelpb.MetricType_METRIC_TYPE_SUMMARY
	ms.Summary = modelpb.SummaryMetricFromVTPool()
	ms.Summary.Count = uint64(dp.Count())
	ms.Summary.Sum = dp.Sum()
	return ms
}

func histogramSample(bucketCounts pcommon.UInt64Slice, explicitBounds pcommon.Float64Slice) (*modelpb.MetricsetSample, bool) {
	// (From opentelemetry-proto/opentelemetry/proto/metrics/v1/metrics.proto)
	//
	// This defines size(explicit_bounds) + 1 (= N) buckets. The boundaries for
	// bucket at index i are:
	//
	// (-infinity, explicit_bounds[i]] for i == 0
	// (explicit_bounds[i-1], explicit_bounds[i]] for 0 < i < N-1
	// (explicit_bounds[i], +infinity) for i == N-1
	//
	// The values in the explicit_bounds array must be strictly increasing.
	//
	if bucketCounts.Len() != explicitBounds.Len()+1 || explicitBounds.Len() == 0 {
		return modelpb.MetricsetSampleFromVTPool(), false
	}

	// For the bucket values, we follow the approach described by Prometheus's
	// histogram_quantile function (https://prometheus.io/docs/prometheus/latest/querying/functions/#histogram_quantile)
	// to achieve consistent percentile aggregation results:
	//
	// "The histogram_quantile() function interpolates quantile values by assuming a linear
	// distribution within a bucket. (...) If a quantile is located in the highest bucket,
	// the upper bound of the second highest bucket is returned. A lower limit of the lowest
	// bucket is assumed to be 0 if the upper bound of that bucket is greater than 0. In that
	// case, the usual linear interpolation is applied within that bucket. Otherwise, the upper
	// bound of the lowest bucket is returned for quantiles located in the lowest bucket."
	values := make([]float64, 0, bucketCounts.Len())
	counts := make([]uint64, 0, bucketCounts.Len())
	for i := 0; i < bucketCounts.Len(); i++ {
		count := bucketCounts.At(i)
		if count == 0 {
			continue
		}

		var value float64
		switch i {
		// (-infinity, explicit_bounds[i]]
		case 0:
			value = explicitBounds.At(i)
			if value > 0 {
				value /= 2
			}

		// (explicit_bounds[i], +infinity)
		case bucketCounts.Len() - 1:
			value = explicitBounds.At(i - 1)

		// [explicit_bounds[i-1], explicit_bounds[i])
		default:
			// Use the midpoint between the boundaries.
			value = explicitBounds.At(i-1) + (explicitBounds.At(i)-explicitBounds.At(i-1))/2.0
		}

		counts = append(counts, uint64(count))
		values = append(values, value)
	}
	ms := modelpb.MetricsetSampleFromVTPool()
	ms.Type = modelpb.MetricType_METRIC_TYPE_HISTOGRAM
	ms.Histogram = modelpb.HistogramFromVTPool()
	ms.Histogram.Counts = counts
	ms.Histogram.Values = values
	return ms, true
}

type metricsets map[metricsetKey]metricset

type metricsetKey struct {
	timestamp time.Time
	signature string // combination of all attributes
}

type metricset struct {
	attributes pcommon.Map
	samples    map[string]*modelpb.MetricsetSample
}

// upsert searches for an existing metricset with the given timestamp and labels,
// and appends the sample to it. If there is no such existing metricset, a new one
// is created.
func (ms metricsets) upsert(timestamp time.Time, attributes pcommon.Map, sample *modelpb.MetricsetSample) {
	// We always record metrics as they are given. We also copy some
	// well-known OpenTelemetry metrics to their Elastic APM equivalents.
	ms.upsertOne(timestamp, attributes, sample)
}

func (ms metricsets) upsertOne(timestamp time.Time, attributes pcommon.Map, sample *modelpb.MetricsetSample) {
	var signatureBuilder strings.Builder
	attributes.Range(func(k string, v pcommon.Value) bool {
		signatureBuilder.WriteString(k)
		signatureBuilder.WriteString(v.AsString())
		return true
	})
	key := metricsetKey{timestamp: timestamp, signature: signatureBuilder.String()}

	m, ok := ms[key]
	if !ok {
		m = metricset{
			attributes: attributes,
			samples:    make(map[string]*modelpb.MetricsetSample),
		}
		ms[key] = m
	}
	m.samples[sample.Name] = sample
}
