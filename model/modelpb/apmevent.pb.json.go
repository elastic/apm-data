package modelpb

import (
	"net"
	"net/netip"

	"github.com/elastic/apm-data/model/internal/modeljson"
)

func (e *APMEvent) toModelJSON(out *modeljson.Document) {
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
		Timestamp:           modeljson.Time(e.Timestamp.AsTime()),
		DataStreamType:      e.DataStream.Type,
		DataStreamDataset:   e.DataStream.Dataset,
		DataStreamNamespace: e.DataStream.Namespace,
		Processor: modeljson.Processor{
			Name:  e.Processor.Name,
			Event: e.Processor.Event,
		},
		Labels:        labels,
		NumericLabels: numericLabels,
		Message:       e.Message,
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
	if !isZero(e.Event) {
		e.Event.toModelJSON(&event)
		doc.Event = &event
	}

	// Set high resolution timestamp.
	//
	// TODO(axw) change @timestamp to use date_nanos, and remove this field.
	var timestampStruct modeljson.Timestamp
	if !e.Timestamp.AsTime().IsZero() {
		processorName := e.Processor.Name
		processorEvent := e.Processor.Event
		if (processorName == "error" && processorEvent == "error") || (processorName == "transaction" && (processorEvent == "transaction" || processorEvent == "span")) {
			timestampStruct.US = int(e.Timestamp.AsTime().UnixNano() / 1000)
			doc.TimestampStruct = &timestampStruct
		}
	}

	cloud := modeljson.Cloud{
		AvailabilityZone: e.Cloud.AvailabilityZone,
		Provider:         e.Cloud.Provider,
		Region:           e.Cloud.Region,
		Account: modeljson.CloudAccount{
			ID:   e.Cloud.AccountId,
			Name: e.Cloud.AccountName,
		},
		Service: modeljson.CloudService{
			Name: e.Cloud.ServiceName,
		},
		Project: modeljson.CloudProject{
			ID:   e.Cloud.ProjectId,
			Name: e.Cloud.ProjectName,
		},
		Instance: modeljson.CloudInstance{
			ID:   e.Cloud.InstanceId,
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
				ID: e.Cloud.Origin.AccountId,
			},
			Service: modeljson.CloudService{
				Name: e.Cloud.Origin.ServiceName,
			},
		}
	}
	setNonZero(&doc.Cloud, &cloud)

	faas := modeljson.FAAS{
		ID:        e.Faas.Id,
		Name:      e.Faas.Name,
		Version:   e.Faas.Version,
		Execution: e.Faas.Execution,
		Coldstart: e.Faas.ColdStart,
		Trigger: modeljson.FAASTrigger{
			Type:      e.Faas.TriggerType,
			RequestID: e.Faas.TriggerRequestId,
		},
	}
	setNonZero(&doc.FAAS, &faas)

	device := modeljson.Device{
		ID:           e.Device.Id,
		Manufacturer: e.Device.Manufacturer,
		Model: modeljson.DeviceModel{
			Name:       e.Device.Model.Name,
			Identifier: e.Device.Model.Identifier,
		},
	}
	setNonZero(&doc.Device, &device)

	network := modeljson.Network{
		Connection: modeljson.NetworkConnection{
			Type:    e.Network.Connection.Type,
			Subtype: e.Network.Connection.Subtype,
		},
		Carrier: modeljson.NetworkCarrier{
			Name: e.Network.Carrier.Name,
			MCC:  e.Network.Carrier.Mcc,
			MNC:  e.Network.Carrier.Mnc,
			ICC:  e.Network.Carrier.Icc,
		},
	}
	setNonZero(&doc.Network, &network)

	observer := modeljson.Observer{
		Hostname: e.Observer.Hostname,
		Name:     e.Observer.Name,
		Type:     e.Observer.Type,
		Version:  e.Observer.Version,
	}
	setNonZero(&doc.Observer, &observer)

	container := modeljson.Container{
		ID:      e.Container.Id,
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
			UID:  e.Kubernetes.PodUid,
		},
	}
	setNonZero(&doc.Kubernetes, &kubernetes)

	agent := modeljson.Agent{
		Name:             e.Agent.Name,
		Version:          e.Agent.Version,
		EphemeralID:      e.Agent.EphemeralId,
		ActivationMethod: e.Agent.ActivationMethod,
	}
	setNonZero(&doc.Agent, &agent)

	trace := modeljson.Trace{
		ID: e.Trace.Id,
	}
	setNonZero(&doc.Trace, &trace)

	user := modeljson.User{
		Domain: e.User.Domain,
		ID:     e.User.Id,
		Email:  e.User.Email,
		Name:   e.User.Name,
	}
	setNonZero(&doc.User, &user)

	source := modeljson.Source{
		Domain: e.Source.Domain,
		Port:   int(e.Source.Port),
	}
	if ip, err := netip.ParseAddr(e.Source.Ip); err == nil {
		source.IP = modeljson.IP(ip)
	}
	if e.Source.Nat != nil {
		if ip, err := netip.ParseAddr(e.Source.Nat.Ip); err == nil {
			source.NAT.IP = modeljson.IP(ip)
		}
	}
	setNonZero(&doc.Source, &source)

	parent := modeljson.Parent{
		ID: e.Parent.Id,
	}
	setNonZero(&doc.Parent, &parent)

	child := modeljson.Child{
		ID: e.Child.Id,
	}
	if len(child.ID) > 0 {
		doc.Child = &child
	}

	client := modeljson.Client{Domain: e.Client.Domain, Port: int(e.Client.Port)}
	if _, err := netip.ParseAddr(e.Client.Ip); err == nil {
		client.IP = e.Client.Ip
	}
	setNonZero(&doc.Client, &client)

	userAgent := modeljson.UserAgent{
		Original: e.UserAgent.Original,
		Name:     e.UserAgent.Name,
	}
	setNonZero(&doc.UserAgent, &userAgent)

	service := modeljson.Service{
		Name:        e.Service.Name,
		Version:     e.Service.Version,
		Environment: e.Service.Environment,
	}
	serviceNode := modeljson.ServiceNode{
		Name: e.Service.Node.Name,
	}
	setNonZero(&service.Node, &serviceNode)
	serviceLanguage := modeljson.Language{
		Name:    e.Service.Language.Name,
		Version: e.Service.Language.Version,
	}
	setNonZero(&service.Language, &serviceLanguage)
	serviceRuntime := modeljson.Runtime{
		Name:    e.Service.Runtime.Name,
		Version: e.Service.Runtime.Version,
	}
	setNonZero(&service.Runtime, &serviceRuntime)
	serviceFramework := modeljson.Framework{
		Name:    e.Service.Framework.Name,
		Version: e.Service.Framework.Version,
	}
	setNonZero(&service.Framework, &serviceFramework)
	var serviceOrigin modeljson.ServiceOrigin
	var serviceTarget modeljson.ServiceTarget
	if e.Service.Origin != nil {
		serviceOrigin = modeljson.ServiceOrigin{
			ID:      e.Service.Origin.Id,
			Name:    e.Service.Origin.Name,
			Version: e.Service.Origin.Version,
		}
		service.Origin = &serviceOrigin
	}
	if e.Service.Target != nil {
		serviceTarget = modeljson.ServiceTarget{
			Name: e.Service.Target.Name,
			Type: e.Service.Target.Type,
		}
		service.Target = &serviceTarget
	}
	setNonZero(&doc.Service, &service)

	http := modeljson.HTTP{
		Version: e.Http.Version,
	}
	var httpRequest modeljson.HTTPRequest
	var httpRequestBody modeljson.HTTPRequestBody
	var httpResponse modeljson.HTTPResponse
	if e.Http.Request != nil {
		httpRequest = modeljson.HTTPRequest{
			ID:       e.Http.Request.Id,
			Method:   e.Http.Request.Method,
			Referrer: e.Http.Request.Referrer,
			Headers:  e.Http.Request.Headers.AsMap(),
			Env:      e.Http.Request.Env.AsMap(),
			Cookies:  e.Http.Request.Cookies.AsMap(),
		}
		if e.Http.Request.Body != nil {
			httpRequestBody.Original = e.Http.Request.Body
			httpRequest.Body = &httpRequestBody
		}
		http.Request = &httpRequest
	}
	if e.Http.Response != nil {
		httpResponse = modeljson.HTTPResponse{
			StatusCode:      int(e.Http.Response.StatusCode),
			Headers:         e.Http.Response.Headers.AsMap(),
			Finished:        e.Http.Response.Finished,
			HeadersSent:     e.Http.Response.HeadersSent,
			TransferSize:    e.Http.Response.TransferSize,
			EncodedBodySize: e.Http.Response.EncodedBodySize,
			DecodedBodySize: e.Http.Response.DecodedBodySize,
		}
		http.Response = &httpResponse
	}
	setNonZero(&doc.HTTP, &http)

	host := modeljson.Host{
		Hostname:     e.Host.Hostname,
		Name:         e.Host.Name,
		ID:           e.Host.Id,
		Architecture: e.Host.Architecture,
		Type:         e.Host.Type,
		IP:           make([]string, 0, len(e.Host.Ip)),
	}
	for _, ip := range e.Host.Ip {
		if _, err := netip.ParseAddr(ip); err == nil {
			host.IP = append(host.IP, ip)
		}
	}
	if len(host.IP) == 0 {
		host.IP = nil
	}
	hostOS := modeljson.OS{
		Name:     e.Host.Os.Name,
		Version:  e.Host.Os.Version,
		Platform: e.Host.Os.Platform,
		Full:     e.Host.Os.Full,
		Type:     e.Host.Os.Type,
	}
	setNonZero(&host.OS, &hostOS)
	if !isZero(host.OS) || !isZero(host.Hostname) || !isZero(host.Name) || !isZero(host.Name) || !isZero(host.ID) ||
		!isZero(host.Architecture) || !isZero(host.Type) || len(host.IP) != 0 {
		doc.Host = &host
	}

	url := modeljson.URL{
		Original: e.Url.Original,
		Scheme:   e.Url.Scheme,
		Full:     e.Url.Full,
		Domain:   e.Url.Domain,
		Path:     e.Url.Path,
		Query:    e.Url.Query,
		Fragment: e.Url.Fragment,
		Port:     int(e.Url.Port),
	}
	setNonZero(&doc.URL, &url)

	log := modeljson.Log{
		Level:  e.Log.Level,
		Logger: e.Log.Logger,
		Origin: modeljson.LogOrigin{
			Function: e.Log.Origin.FunctionName,
			File: modeljson.LogOriginFile{
				Name: e.Log.Origin.File.Name,
				Line: int(e.Log.Origin.File.Line),
			},
		},
	}
	setNonZero(&doc.Log, &log)

	process := modeljson.Process{
		Pid:         int(e.Process.Pid),
		Title:       e.Process.Title,
		CommandLine: e.Process.CommandLine,
		Executable:  e.Process.Executable,
		Args:        e.Process.Argv,
		Thread: modeljson.ProcessThread{
			Name: e.Process.Thread.Name,
			ID:   int(e.Process.Thread.Id),
		},
		Parent: modeljson.ProcessParent{Pid: e.Process.Ppid},
	}
	if !isZero(process.Pid) || !isZero(process.Title) || !isZero(process.CommandLine) || !isZero(process.Executable) ||
		len(process.Args) != 0 || !isZero(process.Thread) || !isZero(process.Parent) {
		doc.Process = &process
	}

	destination := modeljson.Destination{
		Address: e.Destination.Address,
		Port:    int(e.Destination.Port),
	}
	if e.Destination.Address != "" {
		if ip := net.ParseIP(e.Destination.Address); ip != nil {
			destination.IP = e.Destination.Address
		}
	}
	setNonZero(&doc.Destination, &destination)

	session := modeljson.Session{
		ID:       e.Session.Id,
		Sequence: int(e.Session.Sequence),
	}
	if session.ID != "" {
		doc.Session = &session
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
