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
	"net"
	"time"

	"github.com/elastic/apm-data/model/internal/modeljson"
	"go.elastic.co/fastjson"
)

// APMEvent holds the details of an APM event.
//
// Exactly one of the event fields should be non-nil.
type APMEvent struct {
	// Timestamp holds the event timestamp.
	//
	// See https://www.elastic.co/guide/en/ecs/current/ecs-base.html#field-timestamp
	Timestamp time.Time
	Span      *Span
	// NumericLabels holds the numeric (scaled_float) labels to apply to the event.
	// Supports slice values.
	NumericLabels NumericLabels
	// Labels holds the string (keyword) labels to apply to the event, stored as
	// keywords. Supports slice values.
	//
	// See https://www.elastic.co/guide/en/ecs/current/ecs-base.html#field-labels
	Labels      Labels
	Transaction *Transaction
	Metricset   *Metricset
	Error       *Error
	Cloud       Cloud
	Service     Service
	FAAS        FAAS
	Network     Network
	Container   Container
	User        User
	Device      Device
	Kubernetes  Kubernetes
	Observer    Observer
	// DataStream optionally holds data stream identifiers.
	DataStream DataStream
	Agent      Agent
	Processor  Processor
	HTTP       HTTP
	UserAgent  UserAgent
	Parent     Parent
	// Message holds the message for log events.
	//
	// See https://www.elastic.co/guide/en/ecs/current/ecs-base.html#field-message
	Message     string
	Trace       Trace
	Host        Host
	URL         URL
	Log         Log
	Source      Source
	Client      Client
	Child       Child
	Destination Destination
	Session     Session
	Process     Process
	Event       Event
}

// MarshalJSON marshals e as JSON.
func (e *APMEvent) MarshalJSON() ([]byte, error) {
	var w fastjson.Writer
	if err := e.MarshalFastJSON(&w); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// MarshalFastJSON marshals e as JSON, writing the result to w.
func (e *APMEvent) MarshalFastJSON(w *fastjson.Writer) error {
	var labels map[string]modeljson.Label
	if n := len(e.Labels); n > 0 {
		labels = make(map[string]modeljson.Label)
		for k, label := range e.Labels {
			labels[sanitizeLabelKey(k)] = modeljson.Label{
				Value:  label.Value,
				Values: label.Values,
			}
		}
	}

	var numericLabels map[string]modeljson.NumericLabel
	if n := len(e.NumericLabels); n > 0 {
		numericLabels = make(map[string]modeljson.NumericLabel)
		for k, label := range e.NumericLabels {
			numericLabels[sanitizeLabelKey(k)] = modeljson.NumericLabel{
				Value:  label.Value,
				Values: label.Values,
			}
		}
	}

	doc := modeljson.Document{
		Timestamp:           modeljson.Time(e.Timestamp),
		DataStreamType:      e.DataStream.Type,
		DataStreamDataset:   e.DataStream.Dataset,
		DataStreamNamespace: e.DataStream.Namespace,
		Processor:           modeljson.Processor(e.Processor),
		Labels:              labels,
		NumericLabels:       numericLabels,
		Message:             e.Message,
	}

	var transaction modeljson.Transaction
	if e.Transaction != nil {
		e.Transaction.toModelJSON(&transaction, e.Processor == MetricsetProcessor)
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
	if !isZero(e.Event) {
		e.Event.toModelJSON(&event)
		doc.Event = &event
	}

	// Set high resolution timestamp.
	//
	// TODO(axw) change @timestamp to use date_nanos, and remove this field.
	var timestampStruct modeljson.Timestamp
	if !e.Timestamp.IsZero() {
		switch e.Processor {
		case TransactionProcessor, SpanProcessor, ErrorProcessor:
			timestampStruct.US = int(e.Timestamp.UnixNano() / 1000)
			doc.TimestampStruct = &timestampStruct
		}
	}

	cloud := modeljson.Cloud{
		AvailabilityZone: e.Cloud.AvailabilityZone,
		Provider:         e.Cloud.Provider,
		Region:           e.Cloud.Region,
		Account: modeljson.CloudAccount{
			ID:   e.Cloud.AccountID,
			Name: e.Cloud.AccountName,
		},
		Service: modeljson.CloudService{
			Name: e.Cloud.ServiceName,
		},
		Project: modeljson.CloudProject{
			ID:   e.Cloud.ProjectID,
			Name: e.Cloud.ProjectName,
		},
		Instance: modeljson.CloudInstance{
			ID:   e.Cloud.InstanceID,
			Name: e.Cloud.InstanceName,
		},
		Machine: modeljson.CloudMachine{
			Type: e.Cloud.MachineType,
		},
	}
	if e.Cloud.Origin != nil {
		cloud.Origin = modeljson.CloudOrigin{
			Provider: e.Cloud.Origin.Provider,
			Region:   e.Cloud.Origin.Region,
			Account: modeljson.CloudAccount{
				ID: e.Cloud.Origin.AccountID,
			},
			Service: modeljson.CloudService{
				Name: e.Cloud.Origin.ServiceName,
			},
		}
	}
	setNonZero(&doc.Cloud, &cloud)

	faas := modeljson.FAAS{
		ID:        e.FAAS.ID,
		Name:      e.FAAS.Name,
		Version:   e.FAAS.Version,
		Execution: e.FAAS.Execution,
		Coldstart: e.FAAS.Coldstart,
		Trigger: modeljson.FAASTrigger{
			Type:      e.FAAS.TriggerType,
			RequestID: e.FAAS.TriggerRequestID,
		},
	}
	setNonZero(&doc.FAAS, &faas)

	device := modeljson.Device{
		ID:           e.Device.ID,
		Manufacturer: e.Device.Manufacturer,
		Model:        modeljson.DeviceModel(e.Device.Model),
	}
	setNonZero(&doc.Device, &device)

	network := modeljson.Network{
		Connection: modeljson.NetworkConnection(e.Network.Connection),
		Carrier:    modeljson.NetworkCarrier(e.Network.Carrier),
	}
	setNonZero(&doc.Network, &network)

	observer := modeljson.Observer(e.Observer)
	setNonZero(&doc.Observer, &observer)

	container := modeljson.Container{
		ID:      e.Container.ID,
		Name:    e.Container.Name,
		Runtime: e.Container.Runtime,
		Image: modeljson.ContainerImage{
			Name: e.Container.ImageName,
			Tag:  e.Container.ImageTag,
		},
	}
	setNonZero(&doc.Container, &container)

	kubernetes := modeljson.Kubernetes{
		Namespace: e.Kubernetes.Namespace,
		Node: modeljson.KubernetesNode{
			Name: e.Kubernetes.NodeName,
		},
		Pod: modeljson.KubernetesPod{
			Name: e.Kubernetes.PodName,
			UID:  e.Kubernetes.PodUID,
		},
	}
	setNonZero(&doc.Kubernetes, &kubernetes)

	agent := modeljson.Agent(e.Agent)
	setNonZero(&doc.Agent, &agent)

	trace := modeljson.Trace(e.Trace)
	setNonZero(&doc.Trace, &trace)

	user := modeljson.User(e.User)
	setNonZero(&doc.User, &user)

	source := modeljson.Source{
		IP:     modeljson.IP(e.Source.IP),
		Domain: e.Source.Domain,
		Port:   e.Source.Port,
	}
	if e.Source.NAT != nil {
		source.NAT.IP = modeljson.IP(e.Source.NAT.IP)
	}
	setNonZero(&doc.Source, &source)

	parent := modeljson.Parent(e.Parent)
	setNonZero(&doc.Parent, &parent)

	child := modeljson.Child(e.Child)
	if len(child.ID) > 0 {
		doc.Child = &child
	}

	client := modeljson.Client{Domain: e.Client.Domain, Port: e.Client.Port}
	if e.Client.IP.IsValid() {
		client.IP = e.Client.IP.String()
	}
	setNonZero(&doc.Client, &client)

	userAgent := modeljson.UserAgent(e.UserAgent)
	setNonZero(&doc.UserAgent, &userAgent)

	service := modeljson.Service{
		Name:        e.Service.Name,
		Version:     e.Service.Version,
		Environment: e.Service.Environment,
	}
	serviceNode := modeljson.ServiceNode(e.Service.Node)
	setNonZero(&service.Node, &serviceNode)
	serviceLanguage := modeljson.Language(e.Service.Language)
	setNonZero(&service.Language, &serviceLanguage)
	serviceRuntime := modeljson.Runtime(e.Service.Runtime)
	setNonZero(&service.Runtime, &serviceRuntime)
	serviceFramework := modeljson.Framework(e.Service.Framework)
	setNonZero(&service.Framework, &serviceFramework)
	var serviceOrigin modeljson.ServiceOrigin
	var serviceTarget modeljson.ServiceTarget
	if e.Service.Origin != nil {
		serviceOrigin = modeljson.ServiceOrigin(*e.Service.Origin)
		service.Origin = &serviceOrigin
	}
	if e.Service.Target != nil {
		serviceTarget = modeljson.ServiceTarget(*e.Service.Target)
		service.Target = &serviceTarget
	}
	setNonZero(&doc.Service, &service)

	http := modeljson.HTTP{
		Version: e.HTTP.Version,
	}
	var httpRequest modeljson.HTTPRequest
	var httpRequestBody modeljson.HTTPRequestBody
	var httpResponse modeljson.HTTPResponse
	if e.HTTP.Request != nil {
		httpRequest = modeljson.HTTPRequest{
			ID:       e.HTTP.Request.ID,
			Method:   e.HTTP.Request.Method,
			Referrer: e.HTTP.Request.Referrer,
			Headers:  e.HTTP.Request.Headers,
			Env:      e.HTTP.Request.Env,
			Cookies:  e.HTTP.Request.Cookies,
		}
		if e.HTTP.Request.Body != nil {
			httpRequestBody.Original = e.HTTP.Request.Body
			httpRequest.Body = &httpRequestBody
		}
		http.Request = &httpRequest
	}
	if e.HTTP.Response != nil {
		httpResponse = modeljson.HTTPResponse{
			StatusCode:      e.HTTP.Response.StatusCode,
			Headers:         e.HTTP.Response.Headers,
			Finished:        e.HTTP.Response.Finished,
			HeadersSent:     e.HTTP.Response.HeadersSent,
			TransferSize:    e.HTTP.Response.TransferSize,
			EncodedBodySize: e.HTTP.Response.EncodedBodySize,
			DecodedBodySize: e.HTTP.Response.DecodedBodySize,
		}
		http.Response = &httpResponse
	}
	setNonZero(&doc.HTTP, &http)

	host := modeljson.Host{
		Hostname:     e.Host.Hostname,
		Name:         e.Host.Name,
		ID:           e.Host.ID,
		Architecture: e.Host.Architecture,
		Type:         e.Host.Type,
		IP:           make([]string, 0, len(e.Host.IP)),
	}
	for _, ip := range e.Host.IP {
		if ip.IsValid() {
			host.IP = append(host.IP, ip.String())
		}
	}
	if len(host.IP) == 0 {
		host.IP = nil
	}
	hostOS := modeljson.OS(e.Host.OS)
	setNonZero(&host.OS, &hostOS)
	if !isZero(host.OS) || !isZero(host.Hostname) || !isZero(host.Name) || !isZero(host.Name) || !isZero(host.ID) ||
		!isZero(host.Architecture) || !isZero(host.Type) || len(host.IP) != 0 {
		doc.Host = &host
	}

	url := modeljson.URL(e.URL)
	setNonZero(&doc.URL, &url)

	log := modeljson.Log{
		Level:  e.Log.Level,
		Logger: e.Log.Logger,
		Origin: modeljson.LogOrigin{
			Function: e.Log.Origin.FunctionName,
			File:     modeljson.LogOriginFile(e.Log.Origin.File),
		},
	}
	setNonZero(&doc.Log, &log)

	process := modeljson.Process{
		Pid:         e.Process.Pid,
		Title:       e.Process.Title,
		CommandLine: e.Process.CommandLine,
		Executable:  e.Process.Executable,
		Args:        e.Process.Argv,
		Thread:      modeljson.ProcessThread(e.Process.Thread),
		Parent:      modeljson.ProcessParent{Pid: e.Process.Ppid},
	}
	if !isZero(process.Pid) || !isZero(process.Title) || !isZero(process.CommandLine) || !isZero(process.Executable) ||
		len(process.Args) != 0 || !isZero(process.Thread) || !isZero(process.Parent) {
		doc.Process = &process
	}

	destination := modeljson.Destination{
		Address: e.Destination.Address,
		Port:    e.Destination.Port,
	}
	if e.Destination.Address != "" {
		if ip := net.ParseIP(e.Destination.Address); ip != nil {
			destination.IP = e.Destination.Address
		}
	}
	setNonZero(&doc.Destination, &destination)

	session := modeljson.Session(e.Session)
	if session.ID != "" {
		doc.Session = &session
	}

	return net.ErrClosed
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
