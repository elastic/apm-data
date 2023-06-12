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

package modelprocessor_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/apm-data/model/modelpb"
	"github.com/elastic/apm-data/model/modelprocessor"
)

func TestSetDataStream(t *testing.T) {
	tests := []struct {
		input  *modelpb.APMEvent
		output *modelpb.DataStream
	}{{
		input:  &modelpb.APMEvent{},
		output: &modelpb.DataStream{Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.TransactionProcessor()},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.SpanProcessor()},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.TransactionProcessor(), Agent: &modelpb.Agent{Name: "js-base"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm.rum", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.SpanProcessor(), Agent: &modelpb.Agent{Name: "js-base"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm.rum", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.TransactionProcessor(), Agent: &modelpb.Agent{Name: "rum-js"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm.rum", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.SpanProcessor(), Agent: &modelpb.Agent{Name: "rum-js"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm.rum", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.TransactionProcessor(), Agent: &modelpb.Agent{Name: "iOS/swift"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm.rum", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.SpanProcessor(), Agent: &modelpb.Agent{Name: "iOS/swift"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm.rum", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.TransactionProcessor(), Agent: &modelpb.Agent{Name: "go"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.SpanProcessor(), Agent: &modelpb.Agent{Name: "go"}},
		output: &modelpb.DataStream{Type: "traces", Dataset: "apm", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.ErrorProcessor()},
		output: &modelpb.DataStream{Type: "logs", Dataset: "apm.error", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Processor: modelpb.LogProcessor(),
		},
		output: &modelpb.DataStream{Type: "logs", Dataset: "apm.app.unknown", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Processor: modelpb.LogProcessor(),
			Service:   &modelpb.Service{Name: "service-name"},
		},
		output: &modelpb.DataStream{Type: "logs", Dataset: "apm.app.service_name", Namespace: "custom"},
	}, {
		input:  &modelpb.APMEvent{Processor: modelpb.ErrorProcessor(), Agent: &modelpb.Agent{Name: "iOS/swift"}},
		output: &modelpb.DataStream{Type: "logs", Dataset: "apm.error", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Processor: modelpb.LogProcessor(),
			Agent:     &modelpb.Agent{Name: "iOS/swift"},
			Service:   &modelpb.Service{Name: "service-name"},
		},
		output: &modelpb.DataStream{Type: "logs", Dataset: "apm.app.service_name", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Agent:       &modelpb.Agent{Name: "rum-js"},
			Processor:   modelpb.MetricsetProcessor(),
			Service:     &modelpb.Service{Name: "service-name"},
			Metricset:   &modelpb.Metricset{},
			Transaction: &modelpb.Transaction{Name: "foo"},
		},
		output: &modelpb.DataStream{Type: "metrics", Dataset: "apm.internal", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Agent:     &modelpb.Agent{Name: "rum-js"},
			Processor: modelpb.MetricsetProcessor(),
			Service:   &modelpb.Service{Name: "service-name"},
			Metricset: &modelpb.Metricset{
				Samples: []*modelpb.MetricsetSample{
					{Name: "system.memory.total"}, // known agent metric
				},
			},
		},
		output: &modelpb.DataStream{Type: "metrics", Dataset: "apm.internal", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Agent:     &modelpb.Agent{Name: "rum-js"},
			Processor: modelpb.MetricsetProcessor(),
			Service:   &modelpb.Service{Name: "service-name"},
			Metricset: &modelpb.Metricset{
				Samples: []*modelpb.MetricsetSample{
					{Name: "system.memory.total"}, // known agent metric
					{Name: "custom_metric"},       // custom metric
				},
			},
		},
		output: &modelpb.DataStream{Type: "metrics", Dataset: "apm.app.service_name", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Processor:   modelpb.MetricsetProcessor(),
			Service:     &modelpb.Service{Name: "service-name"},
			Metricset:   &modelpb.Metricset{},
			Transaction: &modelpb.Transaction{Name: "foo"},
		},
		output: &modelpb.DataStream{Type: "metrics", Dataset: "apm.internal", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Processor:   modelpb.MetricsetProcessor(),
			Service:     &modelpb.Service{Name: "service-name"},
			Metricset:   &modelpb.Metricset{Name: "transaction", Interval: "1m"},
			Transaction: &modelpb.Transaction{Name: "foo"},
		},
		output: &modelpb.DataStream{Type: "metrics", Dataset: "apm.transaction.1m", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Processor:   modelpb.MetricsetProcessor(),
			Service:     &modelpb.Service{Name: "service-name"},
			Metricset:   &modelpb.Metricset{Name: "transaction", Interval: "60m"},
			Transaction: &modelpb.Transaction{Name: "foo"},
		},
		output: &modelpb.DataStream{Type: "metrics", Dataset: "apm.transaction.60m", Namespace: "custom"},
	}, {
		input: &modelpb.APMEvent{
			Processor: modelpb.MetricsetProcessor(),
			Metricset: &modelpb.Metricset{
				Name: "agent_config",
				Samples: []*modelpb.MetricsetSample{
					{Name: "agent_config_applied", Value: 1},
				},
			},
		},
		output: &modelpb.DataStream{Type: "metrics", Dataset: "apm.internal", Namespace: "custom"},
	}}

	for _, test := range tests {
		batch := modelpb.Batch{test.input}
		processor := modelprocessor.SetDataStream{Namespace: "custom"}
		err := processor.ProcessBatch(context.Background(), &batch)
		assert.NoError(t, err)
		assert.Equal(t, test.output, batch[0].DataStream)
	}

}

func TestSetDataStreamInternalMetricsetTypeUnit(t *testing.T) {
	batch := modelpb.Batch{{
		Agent:     &modelpb.Agent{Name: "rum-js"},
		Processor: modelpb.MetricsetProcessor(),
		Service:   &modelpb.Service{Name: "service-name"},
		Metricset: &modelpb.Metricset{
			Samples: []*modelpb.MetricsetSample{
				{Name: "system.memory.total", Type: modelpb.MetricType_METRIC_TYPE_GAUGE, Unit: "byte"},
				{Name: "system.process.memory.size", Type: modelpb.MetricType_METRIC_TYPE_GAUGE, Unit: "byte"},
			},
		},
	}}

	processor := modelprocessor.SetDataStream{Namespace: "custom"}
	err := processor.ProcessBatch(context.Background(), &batch)
	assert.NoError(t, err)
	for _, sample := range batch[0].Metricset.Samples {
		assert.Zero(t, sample.Type)
		assert.Zero(t, sample.Unit)
	}
}

func TestSetDataStreamServiceName(t *testing.T) {
	processor := modelprocessor.SetDataStream{Namespace: "custom"}
	batch := modelpb.Batch{{
		Processor: modelpb.MetricsetProcessor(),
		Service:   &modelpb.Service{Name: "UPPER-CASE"},
	}, {
		Processor: modelpb.MetricsetProcessor(),
		Service:   &modelpb.Service{Name: "\\/*?\"<>| ,#:"},
	}}

	err := processor.ProcessBatch(context.Background(), &batch)
	assert.NoError(t, err)
	assert.Equal(t, "apm.app.upper_case", batch[0].DataStream.Dataset)
	assert.Equal(t, "apm.app.____________", batch[1].DataStream.Dataset)
}
