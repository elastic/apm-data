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
	"go.elastic.co/fastjson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	var outProto modelpb.APMEvent
	e.ToModelProtobuf(&outProto)
	return outProto.MarshalFastJSON(w)
}

func (e *APMEvent) ToModelProtobuf(out *modelpb.APMEvent) {
	var labels map[string]*modelpb.LabelValue
	if n := len(e.Labels); n > 0 {
		labels = make(map[string]*modelpb.LabelValue, n)
		for k, v := range e.Labels {
			pbv := &modelpb.LabelValue{
				Value:  v.Value,
				Values: v.Values,
				Global: v.Global,
			}
			labels[k] = pbv
		}
	}

	var numericLabels map[string]*modelpb.NumericLabelValue
	if n := len(e.NumericLabels); n > 0 {
		numericLabels = make(map[string]*modelpb.NumericLabelValue, n)
		for k, v := range e.NumericLabels {
			numericLabels[k] = &modelpb.NumericLabelValue{
				Value:  v.Value,
				Values: v.Values,
				Global: v.Global,
			}
		}
	}

	*out = modelpb.APMEvent{
		Timestamp:     timestamppb.New(e.Timestamp),
		NumericLabels: numericLabels,
		Labels:        labels,
		Message:       e.Message,
	}

	if !isZero(e.DataStream) {
		out.DataStream = &modelpb.DataStream{
			Type:      e.DataStream.Type,
			Dataset:   e.DataStream.Dataset,
			Namespace: e.DataStream.Namespace,
		}
	}

	if !isZero(e.Processor) {
		out.Processor = &modelpb.Processor{
			Name:  e.Processor.Name,
			Event: e.Processor.Event,
		}
	}

	var transaction modelpb.Transaction
	if e.Transaction != nil {
		e.Transaction.toModelProtobuf(&transaction, e.Processor == MetricsetProcessor)
		out.Transaction = &transaction
	}

	var span modelpb.Span
	if e.Span != nil {
		e.Span.toModelProtobuf(&span)
		out.Span = &span
	}

	var metricset modelpb.Metricset
	if e.Metricset != nil {
		e.Metricset.toModelProtobuf(&metricset)
		out.Metricset = &metricset
	}

	var errorStruct modelpb.Error
	if e.Error != nil {
		e.Error.toModelProtobuf(&errorStruct)
		out.Error = &errorStruct
	}

	var event modelpb.Event
	if !isZero(e.Event) {
		e.Event.toModelProtobuf(&event)
		out.Event = &event
	}

	if !isZero(e.Cloud) {
		cloud := modelpb.Cloud{
			AvailabilityZone: e.Cloud.AvailabilityZone,
			Provider:         e.Cloud.Provider,
			Region:           e.Cloud.Region,
			AccountId:        e.Cloud.AccountID,
			AccountName:      e.Cloud.AccountName,
			ServiceName:      e.Cloud.ServiceName,
			ProjectId:        e.Cloud.ProjectID,
			ProjectName:      e.Cloud.ProjectName,
			InstanceId:       e.Cloud.InstanceID,
			InstanceName:     e.Cloud.InstanceName,
			MachineType:      e.Cloud.MachineType,
		}
		if e.Cloud.Origin != nil {
			cloud.Origin = &modelpb.CloudOrigin{
				Provider:    e.Cloud.Origin.Provider,
				Region:      e.Cloud.Origin.Region,
				AccountId:   e.Cloud.Origin.AccountID,
				ServiceName: e.Cloud.Origin.ServiceName,
			}
		}
		out.Cloud = &cloud
	}

	if !isZero(e.FAAS) {
		out.Faas = &modelpb.Faas{
			Id:               e.FAAS.ID,
			Name:             e.FAAS.Name,
			Version:          e.FAAS.Version,
			Execution:        e.FAAS.Execution,
			ColdStart:        e.FAAS.Coldstart,
			TriggerType:      e.FAAS.TriggerType,
			TriggerRequestId: e.FAAS.TriggerRequestID,
		}
	}

	if !isZero(e.Device) {
		device := modelpb.Device{
			Id:           e.Device.ID,
			Manufacturer: e.Device.Manufacturer,
		}
		if !isZero(e.Device.Model) {
			device.Model = &modelpb.DeviceModel{
				Name:       e.Device.Model.Name,
				Identifier: e.Device.Model.Identifier,
			}
		}
		out.Device = &device
	}

	if !isZero(e.Network) {
		network := modelpb.Network{}
		if !isZero(e.Network.Connection) {
			network.Connection = &modelpb.NetworkConnection{
				Type:    e.Network.Connection.Type,
				Subtype: e.Network.Connection.Subtype,
			}
		}
		if !isZero(e.Network.Carrier) {
			network.Carrier = &modelpb.NetworkCarrier{
				Name: e.Network.Carrier.Name,
				Mcc:  e.Network.Carrier.MCC,
				Mnc:  e.Network.Carrier.MNC,
				Icc:  e.Network.Carrier.ICC,
			}
		}
		out.Network = &network
	}

	if !isZero(e.Observer) {
		out.Observer = &modelpb.Observer{
			Hostname: e.Observer.Hostname,
			Name:     e.Observer.Name,
			Type:     e.Observer.Type,
			Version:  e.Observer.Version,
		}
	}

	if !isZero(e.Container) {
		out.Container = &modelpb.Container{
			Id:        e.Container.ID,
			Name:      e.Container.Name,
			Runtime:   e.Container.Runtime,
			ImageName: e.Container.ImageName,
			ImageTag:  e.Container.ImageTag,
		}
	}

	if !isZero(e.Kubernetes) {
		out.Kubernetes = &modelpb.Kubernetes{
			Namespace: e.Kubernetes.Namespace,
			NodeName:  e.Kubernetes.NodeName,
			PodName:   e.Kubernetes.PodName,
			PodUid:    e.Kubernetes.PodUID,
		}
	}

	if !isZero(e.Agent) {
		out.Agent = &modelpb.Agent{
			Name:             e.Agent.Name,
			Version:          e.Agent.Version,
			EphemeralId:      e.Agent.EphemeralID,
			ActivationMethod: e.Agent.ActivationMethod,
		}
	}

	if !isZero(e.Trace) {
		out.Trace = &modelpb.Trace{
			Id: e.Trace.ID,
		}
	}

	if !isZero(e.User) {
		out.User = &modelpb.User{
			Domain: e.User.Domain,
			Id:     e.User.ID,
			Email:  e.User.Email,
			Name:   e.User.Name,
		}
	}

	if !isZero(e.Source) {
		source := modelpb.Source{
			Ip:     e.Source.IP.String(),
			Domain: e.Source.Domain,
			Port:   uint32(e.Source.Port),
		}
		if e.Source.NAT != nil {
			source.Nat = &modelpb.NAT{
				Ip: e.Source.NAT.IP.String(),
			}
		}
		out.Source = &source
	}

	if !isZero(e.Parent) {
		out.Parent = &modelpb.Parent{
			Id: e.Parent.ID,
		}
	}

	if len(e.Child.ID) > 0 {
		out.Child = &modelpb.Child{
			Id: e.Child.ID,
		}
	}

	if !isZero(e.Client) {
		client := modelpb.Client{
			Domain: e.Client.Domain,
			Port:   uint32(e.Client.Port),
		}
		if e.Client.IP.IsValid() {
			client.Ip = e.Client.IP.String()
		}
		out.Client = &client
	}

	if !isZero(e.UserAgent) {
		out.UserAgent = &modelpb.UserAgent{
			Original: e.UserAgent.Original,
			Name:     e.UserAgent.Name,
		}
	}

	if !isZero(e.Service) {
		service := modelpb.Service{
			Name:        e.Service.Name,
			Version:     e.Service.Version,
			Environment: e.Service.Environment,
		}
		if e.Service.Origin != nil {
			service.Origin = &modelpb.ServiceOrigin{
				Id:      e.Service.Origin.ID,
				Name:    e.Service.Origin.Name,
				Version: e.Service.Origin.Version,
			}
		}
		if e.Service.Target != nil {
			service.Target = &modelpb.ServiceTarget{
				Name: e.Service.Target.Name,
				Type: e.Service.Target.Type,
			}
		}
		if !isZero(e.Service.Language) {
			service.Language = &modelpb.Language{
				Name:    e.Service.Language.Name,
				Version: e.Service.Language.Version,
			}
		}
		if !isZero(e.Service.Runtime) {
			service.Runtime = &modelpb.Runtime{
				Name:    e.Service.Runtime.Name,
				Version: e.Service.Runtime.Version,
			}
		}
		if !isZero(e.Service.Framework) {
			service.Framework = &modelpb.Framework{
				Name:    e.Service.Framework.Name,
				Version: e.Service.Framework.Version,
			}
		}
		if !isZero(e.Service.Node) {
			service.Node = &modelpb.ServiceNode{
				Name: e.Service.Node.Name,
			}
		}
		out.Service = &service
	}

	if !isZero(e.HTTP) {
		http := modelpb.HTTP{
			Version: e.HTTP.Version,
		}
		if e.HTTP.Request != nil {
			httpRequest := modelpb.HTTPRequest{
				Id:       e.HTTP.Request.ID,
				Method:   e.HTTP.Request.Method,
				Referrer: e.HTTP.Request.Referrer,
			}
			if e.HTTP.Request.Body != nil {
				if v, err := structpb.NewValue(e.HTTP.Request.Body); err == nil {
					httpRequest.Body = v
				}
			}
			if len(e.HTTP.Request.Headers) != 0 {
				if m, err := structpb.NewStruct(e.HTTP.Request.Headers); err == nil {
					httpRequest.Headers = m
				}
			}
			if len(e.HTTP.Request.Env) != 0 {
				if m, err := structpb.NewStruct(e.HTTP.Request.Env); err == nil {
					httpRequest.Env = m
				}
			}
			if len(e.HTTP.Request.Cookies) != 0 {
				if m, err := structpb.NewStruct(e.HTTP.Request.Cookies); err == nil {
					httpRequest.Cookies = m
				}
			}
			http.Request = &httpRequest
		}
		if e.HTTP.Response != nil {
			httpResponse := modelpb.HTTPResponse{
				StatusCode:      int32(e.HTTP.Response.StatusCode),
				Finished:        e.HTTP.Response.Finished,
				HeadersSent:     e.HTTP.Response.HeadersSent,
				TransferSize:    e.HTTP.Response.TransferSize,
				EncodedBodySize: e.HTTP.Response.EncodedBodySize,
				DecodedBodySize: e.HTTP.Response.DecodedBodySize,
			}
			if len(e.HTTP.Response.Headers) != 0 {
				if m, err := structpb.NewStruct(e.HTTP.Response.Headers); err == nil {
					httpResponse.Headers = m
				}
			}
			http.Response = &httpResponse
		}
		out.Http = &http
	}

	if !isZero(e.Host.OS) || !isZero(e.Host.Hostname) || !isZero(e.Host.Name) || !isZero(e.Host.Name) || !isZero(e.Host.ID) ||
		!isZero(e.Host.Architecture) || !isZero(e.Host.Type) || len(e.Host.IP) != 0 {
		host := modelpb.Host{
			Hostname:     e.Host.Hostname,
			Name:         e.Host.Name,
			Id:           e.Host.ID,
			Architecture: e.Host.Architecture,
			Type:         e.Host.Type,
			Ip:           make([]string, 0, len(e.Host.IP)),
		}
		for _, ip := range e.Host.IP {
			if ip.IsValid() {
				host.Ip = append(host.Ip, ip.String())
			}
		}
		if len(host.Ip) == 0 {
			host.Ip = nil
		}
		if !isZero(e.Host.OS) {
			host.Os = &modelpb.OS{
				Name:     e.Host.OS.Name,
				Version:  e.Host.OS.Version,
				Platform: e.Host.OS.Platform,
				Full:     e.Host.OS.Full,
				Type:     e.Host.OS.Type,
			}
		}
		out.Host = &host
	}

	if !isZero(e.URL) {
		out.Url = &modelpb.URL{
			Original: e.URL.Original,
			Scheme:   e.URL.Scheme,
			Full:     e.URL.Full,
			Domain:   e.URL.Domain,
			Path:     e.URL.Path,
			Query:    e.URL.Query,
			Fragment: e.URL.Fragment,
			Port:     uint32(e.URL.Port),
		}
	}

	if !isZero(e.Log) {
		log := modelpb.Log{
			Level:  e.Log.Level,
			Logger: e.Log.Logger,
		}
		if !isZero(e.Log.Origin) {
			log.Origin = &modelpb.LogOrigin{
				FunctionName: e.Log.Origin.FunctionName,
			}
			if !isZero(e.Log.Origin.File) {
				log.Origin.File = &modelpb.LogOriginFile{
					Name: e.Log.Origin.File.Name,
					Line: int32(e.Log.Origin.File.Line),
				}
			}
		}
		out.Log = &log
	}

	if !isZero(e.Process.Pid) || !isZero(e.Process.Title) || !isZero(e.Process.CommandLine) || !isZero(e.Process.Executable) ||
		len(e.Process.Argv) != 0 || !isZero(e.Process.Thread) || !isZero(e.Process.Ppid) {
		process := modelpb.Process{
			Ppid:        e.Process.Ppid,
			Title:       e.Process.Title,
			CommandLine: e.Process.CommandLine,
			Executable:  e.Process.Executable,
			Argv:        e.Process.Argv,
			Pid:         uint32(e.Process.Pid),
		}
		if !isZero(e.Process.Thread) {
			process.Thread = &modelpb.ProcessThread{
				Name: e.Process.Thread.Name,
				Id:   int32(e.Process.Thread.ID),
			}
		}
		out.Process = &process
	}

	if !isZero(e.Destination) {
		out.Destination = &modelpb.Destination{
			Address: e.Destination.Address,
			Port:    uint32(e.Destination.Port),
		}
	}

	if !isZero(e.Session) {
		out.Session = &modelpb.Session{
			Id:       e.Session.ID,
			Sequence: int64(e.Session.Sequence),
		}
	}
}

func isZero[T comparable](t T) bool {
	var zero T
	return t == zero
}
