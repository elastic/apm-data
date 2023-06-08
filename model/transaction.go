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
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	// TransactionProcessor is the Processor value that should be assigned to transaction events.
	TransactionProcessor = Processor{Name: "transaction", Event: "transaction"}
)

// Transaction holds values for transaction.* fields. This may be used in
// transaction, span, and error events (i.e. transaction.id), as well as
// internal metrics such as breakdowns (i.e. including transaction.name).
type Transaction struct {
	SpanCount      SpanCount
	UserExperience *UserExperience
	Custom         map[string]any
	Marks          TransactionMarks
	Message        *Message
	// Type holds the transaction type: "request", "message", etc.
	Type string
	// Name holds the transaction name: "GET /foo", etc.
	Name string
	// Result holds the transaction result: "HTTP 2xx", "OK", "Error", etc.
	Result string
	ID     string
	// DurationHistogram holds a transaction duration histogram,
	// with bucket values measured in microseconds, for transaction
	// duration metrics.
	DurationHistogram Histogram
	// DroppedSpanStats holds a list of the spans that were dropped by an
	// agent; not indexed.
	DroppedSpansStats []DroppedSpanStats
	// DurationSummary holds an aggregated transaction duration summary,
	// for service metrics. The DurationSummary.Sum field has microsecond
	// resolution.
	//
	// NOTE(axw) this is used only for service metrics, which are in technical
	// preview. Do not use this field without discussion, as the field mapping
	// is subject to removal.
	DurationSummary SummaryMetric
	// RepresentativeCount holds the approximate number of
	// transactions that this transaction represents for aggregation.
	// This is used for scaling metrics.
	RepresentativeCount float64
	// Sampled holds the transaction's sampling decision.
	//
	// If Sampled is false, then it will be omitted from the output event.
	Sampled bool
	// Root indicates whether or not the transaction is the trace root.
	//
	// If Root is false, it will be omitted from the output event.
	Root bool
}

type SpanCount struct {
	Dropped *uint32
	Started *uint32
}

func (e *Transaction) toModelProtobuf(out *modelpb.Transaction, metricset bool) {
	var marks map[string]*modelpb.TransactionMark
	if n := len(e.Marks); n > 0 {
		marks = make(map[string]*modelpb.TransactionMark, n)
		for k, mark := range e.Marks {
			marks[k] = &modelpb.TransactionMark{
				Measurements: mark,
			}
		}
	}
	var message *modelpb.Message
	if e.Message != nil {
		message = &modelpb.Message{}
		e.Message.toModelProtobuf(message)
	}
	var userExperience *modelpb.UserExperience
	if e.UserExperience != nil {
		userExperience = &modelpb.UserExperience{
			CumulativeLayoutShift: e.UserExperience.CumulativeLayoutShift,
			FirstInputDelay:       e.UserExperience.FirstInputDelay,
			TotalBlockingTime:     e.UserExperience.TotalBlockingTime,
		}
		if !isZero(e.UserExperience.Longtask) {
			userExperience.LongTask = &modelpb.LongtaskMetrics{
				Count: int64(e.UserExperience.Longtask.Count),
				Sum:   e.UserExperience.Longtask.Sum,
				Max:   e.UserExperience.Longtask.Max,
			}
		}
	}
	var droppedSpansStats []*modelpb.DroppedSpanStats
	if metricset {
		// DroppedSpansStats is only indexed for metric documents, never for events.
		if n := len(e.DroppedSpansStats); n > 0 {
			droppedSpansStats = make([]*modelpb.DroppedSpanStats, n)
			for i, dss := range e.DroppedSpansStats {
				dssProto := modelpb.DroppedSpanStats{
					DestinationServiceResource: dss.DestinationServiceResource,
					ServiceTargetType:          dss.ServiceTargetType,
					ServiceTargetName:          dss.ServiceTargetName,
					Outcome:                    dss.Outcome,
				}
				if !isZero(dss.Duration) {
					dssProto.Duration = &modelpb.AggregatedDuration{
						Count: int64(dss.Duration.Count),
						Sum:   durationpb.New(dss.Duration.Sum),
					}
				}
				droppedSpansStats[i] = &dssProto
			}
		}
	}
	var durationHistogram *modelpb.Histogram
	if len(e.DurationHistogram.Values) != 0 || len(e.DurationHistogram.Counts) != 0 {
		durationHistogram = &modelpb.Histogram{
			Values: e.DurationHistogram.Values,
			Counts: e.DurationHistogram.Counts,
		}
	}
	var durationSummary *modelpb.SummaryMetric
	if !isZero(e.DurationSummary) {
		durationSummary = &modelpb.SummaryMetric{
			Count: e.DurationSummary.Count,
			Sum:   e.DurationSummary.Sum,
		}
	}
	var spanCount *modelpb.SpanCount
	if !isZero(e.SpanCount) {
		spanCount = &modelpb.SpanCount{
			Dropped: e.SpanCount.Dropped,
			Started: e.SpanCount.Started,
		}
	}
	var custom *structpb.Struct
	if len(e.Custom) != 0 {
		if s, err := structpb.NewStruct(e.Custom); err == nil {
			custom = s
		}
	}
	*out = modelpb.Transaction{
		Id:                  e.ID,
		Type:                e.Type,
		Name:                e.Name,
		Result:              e.Result,
		Sampled:             e.Sampled,
		Root:                e.Root,
		RepresentativeCount: e.RepresentativeCount,

		DurationHistogram: durationHistogram,
		DurationSummary:   durationSummary,
		SpanCount:         spanCount,
		DroppedSpansStats: droppedSpansStats,

		Marks:          marks,
		Custom:         custom,
		Message:        message,
		UserExperience: userExperience,
	}
}

type TransactionMarks map[string]TransactionMark

type TransactionMark map[string]float64

type DroppedSpanStats struct {
	DestinationServiceResource string
	ServiceTargetType          string
	ServiceTargetName          string
	Outcome                    string
	Duration                   AggregatedDuration
}
