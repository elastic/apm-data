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
	"go.elastic.co/fastjson"
)

func (e *APMEvent) MarshalJSON() ([]byte, error) {
	var w fastjson.Writer
	if err := e.MarshalFastJSON(&w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (e *APMEvent) MarshalFastJSON(w *fastjson.Writer) error {
	var labels map[string]modeljson.Label
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
	}

	var numericLabels map[string]modeljson.NumericLabel
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
	}

	doc := modeljson.Document{
		Timestamp:     modeljson.Time(e.Timestamp.AsTime()),
		Labels:        labels,
		NumericLabels: numericLabels,
		Message:       e.Message,
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

	var transaction modeljson.Transaction
	if e.Transaction != nil {
		e.Transaction.toModelJSON(&transaction, e.Processor.Name == "metric" && e.Processor.Event == "metric")
		doc.Transaction = &transaction
	}

	var span modeljson.Span
	if e.Span != nil {
		e.Span.toModelJSON(&span)
		doc.Span = &span
	}

	var metricset modeljson.Metricset
	if e.Metricset != nil {
		e.Metricset.toModelJSON(&metricset)
		doc.Metricset = &metricset
		doc.DocCount = e.Metricset.DocCount
	}

	var errorStruct modeljson.Error
	if e.Error != nil {
		e.Error.toModelJSON(&errorStruct)
		doc.Error = &errorStruct
	}

	var event modeljson.Event
	if e.Event != nil {
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

	var cloud modeljson.Cloud
	if e.Cloud != nil {
		e.Cloud.toModelJSON(&cloud)
		doc.Cloud = &cloud
	}

	var fass modeljson.FAAS
	if e.Faas != nil {
		e.Faas.toModelJSON(&fass)
		doc.FAAS = &fass
	}

	var device modeljson.Device
	if e.Device != nil {
		e.Device.toModelJSON(&device)
		doc.Device = &device
	}

	var network modeljson.Network
	if e.Network != nil {
		e.Network.toModelJSON(&network)
		doc.Network = &network
	}

	var observer modeljson.Observer
	if e.Observer != nil {
		e.Observer.toModelJSON(&observer)
		doc.Observer = &observer
	}

	var container modeljson.Container
	if e.Container != nil {
		e.Container.toModelJSON(&container)
		doc.Container = &container
	}

	var kubernetes modeljson.Kubernetes
	if e.Kubernetes != nil {
		e.Kubernetes.toModelJSON(&kubernetes)
		doc.Kubernetes = &kubernetes
	}

	var agent modeljson.Agent
	if e.Agent != nil {
		e.Agent.toModelJSON(&agent)
		doc.Agent = &agent
	}

	if e.Trace != nil {
		doc.Trace = &modeljson.Trace{
			ID: e.Trace.Id,
		}
	}

	var user modeljson.User
	if e.User != nil {
		e.User.toModelJSON(&user)
		doc.User = &user
	}

	var source modeljson.Source
	if e.Source != nil {
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

	var client modeljson.Client
	if e.Client != nil {
		e.Client.toModelJSON(&client)
		doc.Client = &client
	}

	if e.UserAgent != nil {
		doc.UserAgent = &modeljson.UserAgent{
			Original: e.UserAgent.Original,
			Name:     e.UserAgent.Name,
		}
	}

	var service modeljson.Service
	if e.Service != nil {
		e.Service.toModelJSON(&service)
		doc.Service = &service
	}

	var httpStruct modeljson.HTTP
	if e.Http != nil {
		e.Http.toModelJSON(&httpStruct)
		doc.HTTP = &httpStruct
	}

	var host modeljson.Host
	if e.Host != nil {
		e.Host.toModelJSON(&host)
		doc.Host = &host
	}

	var url modeljson.URL
	if e.Url != nil {
		e.Url.toModelJSON(&url)
		doc.URL = &url
	}

	var logStruct modeljson.Log
	if e.Log != nil {
		e.Log.toModelJSON(&logStruct)
		doc.Log = &logStruct
	}

	var process modeljson.Process
	if e.Process != nil {
		e.Process.toModelJSON(&process)
		doc.Process = &process
	}

	var destination modeljson.Destination
	if e.Destination != nil {
		e.Destination.toModelJSON(&destination)
		doc.Destination = &destination
	}

	if e.Session != nil {
		doc.Session = &modeljson.Session{
			ID:       e.Session.Id,
			Sequence: int(e.Session.Sequence),
		}
	}

	return doc.MarshalFastJSON(w)
}
