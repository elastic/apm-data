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

import (
	"github.com/elastic/apm-data/model/internal/modeljson"
)

func (e *APMEvent) toModelJSON(out *modeljson.Document) {
	doc := modeljson.Document{
		Timestamp: modeljson.Time(e.Timestamp.AsTime()),
		Message:   e.Message,
	}

	if n := len(e.Labels); n > 0 {
		labels := make(map[string]modeljson.Label)
		for k, label := range e.Labels {
			if label != nil {
				labels[sanitizeLabelKey(k)] = modeljson.Label{
					Value:  label.Value,
					Values: label.Values,
				}
			}
		}
		doc.Labels = labels
	}

	if n := len(e.NumericLabels); n > 0 {
		numericLabels := make(map[string]modeljson.NumericLabel)
		for k, label := range e.NumericLabels {
			if label != nil {
				numericLabels[sanitizeLabelKey(k)] = modeljson.NumericLabel{
					Value:  label.Value,
					Values: label.Values,
				}
			}
		}
		doc.NumericLabels = numericLabels
	}

	if e.DataStream != nil {
		doc.DataStreamType = e.DataStream.Type
		doc.DataStreamDataset = e.DataStream.Dataset
		doc.DataStreamNamespace = e.DataStream.Namespace
	}

	if e.Processor != nil {
		doc.Processor = modeljson.Processor{
			Name:  e.Processor.Name,
			Event: e.Processor.Event,
		}
	}

	if e.Transaction != nil {
		var transaction modeljson.Transaction
		e.Transaction.toModelJSON(&transaction, e.Processor.Name == "metric" && e.Processor.Event == "metric")
		doc.Transaction = &transaction
	}

	if e.Span != nil {
		var span modeljson.Span
		e.Span.toModelJSON(&span)
		doc.Span = &span
	}

	if e.Metricset != nil {
		var metricset modeljson.Metricset
		e.Metricset.toModelJSON(&metricset)
		doc.Metricset = &metricset
		doc.DocCount = e.Metricset.DocCount
	}

	if e.Error != nil {
		var errorStruct modeljson.Error
		e.Error.toModelJSON(&errorStruct)
		doc.Error = &errorStruct
	}

	if e.Event != nil {
		var event modeljson.Event
		e.Event.toModelJSON(&event)
		doc.Event = &event
	}

	// Set high resolution timestamp.
	//
	// TODO(axw) change @timestamp to use date_nanos, and remove this field.
	var timestampStruct modeljson.Timestamp
	if !e.Timestamp.AsTime().IsZero() {
		if e.Processor != nil {
			processorName := e.Processor.Name
			processorEvent := e.Processor.Event
			if (processorName == "error" && processorEvent == "error") || (processorName == "transaction" && (processorEvent == "transaction" || processorEvent == "span")) {
				timestampStruct.US = int(e.Timestamp.AsTime().UnixNano() / 1000)
				doc.TimestampStruct = &timestampStruct
			}
		}
	}

	if e.Cloud != nil {
		var cloud modeljson.Cloud
		e.Cloud.toModelJSON(&cloud)
		doc.Cloud = &cloud
	}

	if e.Faas != nil {
		var faas modeljson.FAAS
		e.Faas.toModelJSON(&faas)
		doc.FAAS = &faas
	}

	if e.Device != nil {
		var device modeljson.Device
		e.Device.toModelJSON(&device)
		doc.Device = &device
	}

	if e.Network != nil {
		var network modeljson.Network
		e.Network.toModelJSON(&network)
		doc.Network = &network
	}

	if e.Observer != nil {
		var observer modeljson.Observer
		e.Observer.toModelJSON(&observer)
		doc.Observer = &observer
	}

	if e.Container != nil {
		var container modeljson.Container
		e.Container.toModelJSON(&container)
		doc.Container = &container
	}

	if e.Kubernetes != nil {
		var kube modeljson.Kubernetes
		e.Kubernetes.toModelJSON(&kube)
		doc.Kubernetes = &kube
	}

	if e.Agent != nil {
		var agent modeljson.Agent
		e.Agent.toModelJSON(&agent)
		doc.Agent = &agent
	}

	if e.Trace != nil {
		doc.Trace = &modeljson.Trace{
			ID: e.Trace.Id,
		}
	}

	if e.User != nil {
		var user modeljson.User
		e.User.toModelJSON(&user)
		doc.User = &user
	}

	if e.Source != nil {
		var source modeljson.Source
		e.Source.toModelJSON(&source)
		doc.Source = &source
	}

	if e.Parent != nil {
		doc.Parent = &modeljson.Parent{
			ID: e.Parent.Id,
		}
	}

	if e.Child != nil {
		doc.Child = &modeljson.Child{
			ID: e.Child.Id,
		}
	}

	if e.Client != nil {
		var client modeljson.Client
		e.Client.toModelJSON(&client)
		doc.Client = &client
	}

	if e.UserAgent != nil {
		doc.UserAgent = &modeljson.UserAgent{
			Original: e.UserAgent.Original,
			Name:     e.UserAgent.Name,
		}
	}

	if e.Service != nil {
		var service modeljson.Service
		e.Service.toModelJSON(&service)
		doc.Service = &service
	}

	if e.Http != nil {
		var httpStruct modeljson.HTTP
		e.Http.toModelJSON(&httpStruct)
		doc.HTTP = &httpStruct
	}

	if e.Host != nil {
		var host modeljson.Host
		e.Host.toModelJSON(&host)
		doc.Host = &host
	}

	if e.Url != nil {
		var urlStruct modeljson.URL
		e.Url.toModelJSON(&urlStruct)
		doc.URL = &urlStruct
	}

	if e.Log != nil {
		var logStruct modeljson.Log
		e.Log.toModelJSON(&logStruct)
		doc.Log = &logStruct
	}

	if e.Process != nil {
		var process modeljson.Process
		e.Process.toModelJSON(&process)
		doc.Process = &process
	}

	if e.Destination != nil {
		var destination modeljson.Destination
		e.Destination.toModelJSON(&destination)
		doc.Destination = &destination
	}

	if e.Session != nil {
		doc.Session = &modeljson.Session{
			ID:       e.Session.Id,
			Sequence: int(e.Session.Sequence),
		}
	}

	*out = doc
}

func setNonZero[T comparable](to **T, from *T) {
	if !isZero(*from) {
		*to = from
	}
}

func isZero[T comparable](t T) bool {
	var zero T
	return t == zero
}
