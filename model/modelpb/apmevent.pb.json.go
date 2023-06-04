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
	out := modeljson.Document{
		Span: &modeljson.Span{
			Message:     &modeljson.Message{},
			Composite:   &modeljson.SpanComposite{},
			Destination: &modeljson.SpanDestination{},
			DB:          &modeljson.DB{},
		},
		Transaction: &modeljson.Transaction{
			UserExperience: &modeljson.UserExperience{},
			Message:        &modeljson.Message{},
		},
		Metricset: &modeljson.Metricset{},
		Error: &modeljson.Error{
			Exception: &modeljson.Exception{},
			Log:       &modeljson.ErrorLog{},
		},
		TimestampStruct: &modeljson.Timestamp{},
		Cloud:           &modeljson.Cloud{},
		Service: &modeljson.Service{
			Node:      &modeljson.ServiceNode{},
			Language:  &modeljson.Language{},
			Runtime:   &modeljson.Runtime{},
			Framework: &modeljson.Framework{},
			Origin:    &modeljson.ServiceOrigin{},
			Target:    &modeljson.ServiceTarget{},
		},
		FAAS:       &modeljson.FAAS{},
		Network:    &modeljson.Network{},
		Container:  &modeljson.Container{},
		User:       &modeljson.User{},
		Device:     &modeljson.Device{},
		Kubernetes: &modeljson.Kubernetes{},
		Observer:   &modeljson.Observer{},
		Agent:      &modeljson.Agent{},
		HTTP: &modeljson.HTTP{
			Request: &modeljson.HTTPRequest{
				Body: &modeljson.HTTPRequestBody{},
			},
			Response: &modeljson.HTTPResponse{},
		},
		UserAgent: &modeljson.UserAgent{},
		Parent:    &modeljson.Parent{},
		Trace:     &modeljson.Trace{},
		Host: &modeljson.Host{
			OS: &modeljson.OS{},
		},
		URL:         &modeljson.URL{},
		Log:         &modeljson.Log{},
		Source:      &modeljson.Source{},
		Client:      &modeljson.Client{},
		Child:       &modeljson.Child{},
		Destination: &modeljson.Destination{},
		Session:     &modeljson.Session{},
		Process:     &modeljson.Process{},
		Event:       &modeljson.Event{},
	}
	e.updateModelJSON(&out)
	return ErrInvalidLength
	//return out.MarshalFastJSON(w)
}

func (e *APMEvent) updateModelJSON(doc *modeljson.Document) {
	doc.Timestamp = modeljson.Time(e.Timestamp.AsTime())
	doc.Message = e.Message

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
		e.Transaction.toModelJSON(doc.Transaction, e.Processor.Name == "metric" && e.Processor.Event == "metric")
	} else {
		doc.Transaction = nil
	}

	if e.Span != nil {
		e.Span.toModelJSON(doc.Span)
	} else {
		doc.Span = nil
	}

	if e.Metricset != nil {
		e.Metricset.toModelJSON(doc.Metricset)
		doc.DocCount = e.Metricset.DocCount
	} else {
		doc.Metricset = nil
	}

	if e.Error != nil {
		e.Error.toModelJSON(doc.Error)
	} else {
		doc.Error = nil
	}

	if e.Event != nil {
		e.Event.toModelJSON(doc.Event)
	} else {
		doc.Event = nil
	}

	// Set high resolution timestamp.
	//
	// TODO(axw) change @timestamp to use date_nanos, and remove this field.
	var timestampStruct modeljson.Timestamp
	var timestampSet bool = false
	if !e.Timestamp.AsTime().IsZero() {
		if e.Processor != nil {
			processorName := e.Processor.Name
			processorEvent := e.Processor.Event
			if (processorName == "error" && processorEvent == "error") || (processorName == "transaction" && (processorEvent == "transaction" || processorEvent == "span")) {
				timestampStruct.US = int(e.Timestamp.AsTime().UnixNano() / 1000)
				doc.TimestampStruct = &timestampStruct
				timestampSet = true
			}
		}
	}
	if !timestampSet {
		doc.TimestampStruct = nil
	}

	if e.Cloud != nil {
		e.Cloud.toModelJSON(doc.Cloud)
	} else {
		doc.Cloud = nil
	}

	if e.Faas != nil {
		e.Faas.toModelJSON(doc.FAAS)
	} else {
		doc.FAAS = nil
	}

	if e.Device != nil {
		e.Device.toModelJSON(doc.Device)
	} else {
		doc.Device = nil
	}

	if e.Network != nil {
		e.Network.toModelJSON(doc.Network)
	} else {
		doc.Network = nil
	}

	if e.Observer != nil {
		e.Observer.toModelJSON(doc.Observer)
	} else {
		doc.Observer = nil
	}

	if e.Container != nil {
		e.Container.toModelJSON(doc.Container)
	} else {
		doc.Container = nil
	}

	if e.Kubernetes != nil {
		e.Kubernetes.toModelJSON(doc.Kubernetes)
	} else {
		doc.Kubernetes = nil
	}

	if e.Agent != nil {
		e.Agent.toModelJSON(doc.Agent)
	} else {
		doc.Agent = nil
	}

	if e.Trace != nil {
		doc.Trace.ID = e.Trace.Id
	} else {
		doc.Trace = nil
	}

	if e.User != nil {
		e.User.toModelJSON(doc.User)
	} else {
		doc.User = nil
	}

	if e.Source != nil {
		e.Source.toModelJSON(doc.Source)
	} else {
		doc.Source = nil
	}

	if e.Parent != nil {
		doc.Parent.ID = e.Parent.Id
	} else {
		doc.Parent = nil
	}

	if e.Child != nil {
		doc.Child.ID = e.Child.Id
	} else {
		doc.Child = nil
	}

	if e.Client != nil {
		e.Client.toModelJSON(doc.Client)
	} else {
		doc.Client = nil
	}

	if e.UserAgent != nil {
		doc.UserAgent.Original = e.UserAgent.Original
		doc.UserAgent.Name = e.UserAgent.Name
	} else {
		doc.UserAgent = nil
	}

	if e.Service != nil {
		e.Service.toModelJSON(doc.Service)
	} else {
		doc.Service = nil
	}

	if e.Http != nil {
		e.Http.toModelJSON(doc.HTTP)
	} else {
		doc.HTTP = nil
	}

	if e.Host != nil {
		e.Host.toModelJSON(doc.Host)
	} else {
		doc.Host = nil
	}

	if e.Url != nil {
		e.Url.toModelJSON(doc.URL)
	} else {
		doc.URL = nil
	}

	if e.Log != nil {
		e.Log.toModelJSON(doc.Log)
	} else {
		doc.Log = nil
	}

	if e.Process != nil {
		e.Process.toModelJSON(doc.Process)
	} else {
		doc.Process = nil
	}

	if e.Destination != nil {
		e.Destination.toModelJSON(doc.Destination)
	} else {
		doc.Destination = nil
	}

	if e.Session != nil {
		doc.Session.ID = e.Session.Id
		doc.Session.Sequence = int(e.Session.Sequence)
	} else {
		doc.Session = nil
	}
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
