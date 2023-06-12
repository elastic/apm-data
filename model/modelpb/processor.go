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

const (
	spanProcessorName         = "transaction"
	spanProcessorEvent        = "span"
	transactionProcessorName  = "transaction"
	transactionProcessorEvent = "transaction"
	errorProcessorName        = "error"
	errorProcessorEvent       = "error"
	logProcessorName          = "log"
	logProcessorEvent         = "log"
	metricsetProcessorName    = "metric"
	metricsetProcessorEvent   = "metric"
)

func SpanProcessor() *Processor {
	return &Processor{
		Name:  spanProcessorName,
		Event: spanProcessorEvent,
	}
}

func (p *Processor) IsSpan() bool {
	return p.Name == spanProcessorName && p.Event == spanProcessorEvent
}

func TransactionProcessor() *Processor {
	return &Processor{
		Name:  transactionProcessorName,
		Event: transactionProcessorEvent,
	}
}

func (p *Processor) IsTransaction() bool {
	return p.Name == transactionProcessorName && p.Event == transactionProcessorEvent
}

func ErrorProcessor() *Processor {
	return &Processor{
		Name:  errorProcessorName,
		Event: errorProcessorEvent,
	}
}

func (p *Processor) IsError() bool {
	return p.Name == errorProcessorName && p.Event == errorProcessorEvent
}

func LogProcessor() *Processor {
	return &Processor{
		Name:  logProcessorName,
		Event: logProcessorEvent,
	}
}

func (p *Processor) IsLog() bool {
	return p.Name == logProcessorName && p.Event == logProcessorEvent
}

func MetricsetProcessor() *Processor {
	return &Processor{
		Name:  metricsetProcessorName,
		Event: metricsetProcessorEvent,
	}
}

func (p *Processor) IsMetricset() bool {
	return p.Name == metricsetProcessorName && p.Event == metricsetProcessorEvent
}
