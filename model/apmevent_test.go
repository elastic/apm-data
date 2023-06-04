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
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.elastic.co/fastjson"
)

func TestAPMEventFields(t *testing.T) {
	host := "host"
	hostname := "hostname"
	containerID := "container-123"
	serviceName := "myservice"
	uid := "12321"
	mail := "user@email.com"
	agentName := "elastic-node"
	destinationAddress := "1.2.3.4"
	destinationPort := 1234
	traceID := "trace_id"
	parentID := "parent_id"
	httpRequestMethod := "post"
	httpRequestBody := "<html><marquee>hello world</marquee></html>"
	coldstart := true

	for _, test := range []struct {
		input  APMEvent
		output map[string]any
	}{{
		input: APMEvent{
			Agent: Agent{
				Name:    agentName,
				Version: agentVersion,
			},
			Observer:  Observer{Type: "apm-server"},
			Container: Container{ID: containerID},
			Service: Service{
				Name: serviceName,
				Node: ServiceNode{Name: "serviceABC"},
				Origin: &ServiceOrigin{
					ID:      "abc123",
					Name:    serviceName,
					Version: "1.0",
				},
			},
			Host: Host{
				Hostname: hostname,
				Name:     host,
			},
			Client: Client{Domain: "client.domain"},
			Source: Source{
				IP:   netip.MustParseAddr("127.0.0.1"),
				Port: 1234,
				NAT: &NAT{
					IP: netip.MustParseAddr("10.10.10.10"),
				},
			},
			Destination: Destination{Address: destinationAddress, Port: destinationPort},
			Process:     Process{Pid: 1234},
			User:        User{ID: uid, Email: mail},
			Event:       Event{Outcome: "success", Duration: time.Microsecond},
			Session:     Session{ID: "session_id"},
			URL:         URL{Original: "url"},
			Labels: map[string]LabelValue{
				"a": {Value: "b"},
				"c": {Value: "true"},
				"d": {Values: []string{"true", "false"}},
			},
			NumericLabels: map[string]NumericLabelValue{
				"e": {Value: float64(1234)},
				"f": {Values: []float64{1234, 12311}},
			},
			Message:   "bottle",
			Timestamp: time.Date(2019, 1, 3, 15, 17, 4, 908.596*1e6, time.FixedZone("+0100", 3600)),
			Processor: Processor{Name: "processor_name", Event: "processor_event"},
			Trace:     Trace{ID: traceID},
			Parent:    Parent{ID: parentID},
			Child:     Child{ID: []string{"child_1", "child_2"}},
			HTTP: HTTP{
				Request: &HTTPRequest{
					Method: httpRequestMethod,
					Body:   httpRequestBody,
				},
			},
			FAAS: FAAS{
				ID:               "faasID",
				Coldstart:        &coldstart,
				Execution:        "execution",
				TriggerType:      "http",
				TriggerRequestID: "abc123",
				Name:             "faasName",
				Version:          "1.0.0",
			},
			Cloud: Cloud{
				Origin: &CloudOrigin{
					AccountID:   "accountID",
					Provider:    "aws",
					Region:      "us-west-1",
					ServiceName: "serviceName",
				},
			},
		},
		output: map[string]any{
			// common fields
			"agent":     map[string]any{"version": "1.0.0", "name": "elastic-node"},
			"observer":  map[string]any{"type": "apm-server"},
			"container": map[string]any{"id": containerID},
			"host":      map[string]any{"hostname": hostname, "name": host},
			"process":   map[string]any{"pid": 1234.0},
			"service": map[string]any{
				"name": "myservice",
				"node": map[string]any{"name": "serviceABC"},
				"origin": map[string]any{
					"id":      "abc123",
					"name":    "myservice",
					"version": "1.0",
				},
			},
			"user":   map[string]any{"id": "12321", "email": "user@email.com"},
			"client": map[string]any{"domain": "client.domain"},
			"source": map[string]any{
				"ip":   "127.0.0.1",
				"port": 1234.0,
				"nat":  map[string]any{"ip": "10.10.10.10"},
			},
			"destination": map[string]any{
				"address": "1.2.3.4",
				"ip":      "1.2.3.4",
				"port":    1234.0,
			},
			"event":   map[string]any{"outcome": "success", "duration": 1000.0},
			"session": map[string]any{"id": "session_id"},
			"url":     map[string]any{"original": "url"},
			"labels": map[string]any{
				"a": "b",
				"c": "true",
				"d": []any{"true", "false"},
			},
			"numeric_labels": map[string]any{
				"e": float64(1234),
				"f": []any{1234.0, 12311.0},
			},
			"message": "bottle",
			"trace": map[string]any{
				"id": traceID,
			},
			"processor": map[string]any{
				"name":  "processor_name",
				"event": "processor_event",
			},
			"parent": map[string]any{
				"id": parentID,
			},
			"child": map[string]any{
				"id": []any{"child_1", "child_2"},
			},
			"http": map[string]any{
				"request": map[string]any{
					"method": "post",
					"body": map[string]any{
						"original": httpRequestBody,
					},
				},
			},
			"faas": map[string]any{
				"id":        "faasID",
				"coldstart": true,
				"execution": "execution",
				"trigger": map[string]any{
					"type":       "http",
					"request_id": "abc123",
				},
				"name":    "faasName",
				"version": "1.0.0",
			},
			"cloud": map[string]any{
				"origin": map[string]any{
					"account": map[string]any{
						"id": "accountID",
					},
					"provider": "aws",
					"region":   "us-west-1",
					"service": map[string]any{
						"name": "serviceName",
					},
				},
			},
		},
	}, {
		input: APMEvent{
			Processor: TransactionProcessor,
			Timestamp: time.Date(2019, 1, 3, 15, 17, 4, 908.596*1e6, time.FixedZone("+0100", 3600)),
		},
		output: map[string]any{
			"processor": map[string]any{"name": "transaction", "event": "transaction"},
			// timestamp.us is added for transactions, spans, and errors.
			"timestamp": map[string]any{"us": 1546525024908596.0},
		},
	}} {
		m := transformAPMEvent(test.input)
		delete(m, "@timestamp")
		assert.Equal(t, test.output, m)
	}
}

func BenchmarkAPMEventToJSON(b *testing.B) {
	modelEvent := fullvent(b)

	b.Run("to-json", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			modelEvent.MarshalJSON()
		}
	})

}

func fullvent(t testing.TB) *APMEvent {
	return &APMEvent{
		Timestamp: time.Unix(1, 1),
		Span: &Span{
			Message: &Message{
				Body: "body",
				Headers: http.Header{
					"foo": []string{"bar"},
				},
				AgeMillis:  int64Ptr(2),
				QueueName:  "queuename",
				RoutingKey: "routingkey",
			},
			Composite: &Composite{
				CompressionStrategy: "exact_match",
				Count:               1,
				Sum:                 2,
			},
			DestinationService: &DestinationService{
				Type:     "destination_type",
				Name:     "destination_name",
				Resource: "destination_resource",
				ResponseTime: AggregatedDuration{
					Count: 3,
					Sum:   4 * time.Second,
				},
			},
			DB: &DB{
				RowsAffected: uintPtr(5),
				Instance:     "db_instace",
				Statement:    "db_statement",
				Type:         "db_type",
				UserName:     "db_username",
				Link:         "db_link",
			},
			Sync:    boolPtr(true),
			Kind:    "kind",
			Action:  "action",
			Subtype: "subtype",
			ID:      "id",
			Type:    "type",
			Name:    "name",
			Stacktrace: Stacktrace{
				&StacktraceFrame{
					Vars:           randomStructPb(t).AsMap(),
					Lineno:         uintPtr(1),
					Colno:          uintPtr(2),
					Filename:       "frame_filename",
					Classname:      "frame_classname",
					ContextLine:    "frame_contextline",
					Module:         "frame_module",
					Function:       "frame_function",
					AbsPath:        "frame_abspath",
					SourcemapError: "frame_sourcemaperror",
					Original: Original{
						AbsPath:      "orig_abspath",
						Filename:     "orig_filename",
						Classname:    "orig_classname",
						Lineno:       uintPtr(3),
						Colno:        uintPtr(4),
						Function:     "orig_function",
						LibraryFrame: true,
					},
					PreContext:          []string{"pre"},
					PostContext:         []string{"post"},
					LibraryFrame:        true,
					SourcemapUpdated:    true,
					ExcludeFromGrouping: true,
				},
			},
			Links: []SpanLink{
				{
					Trace: Trace{
						ID: "trace_id",
					},
					Span: Span{
						Kind:    "kind1",
						Action:  "action1",
						Subtype: "subtype1",
						ID:      "id1",
						Type:    "type1",
						Name:    "name1",
					},
				},
			},
			SelfTime: AggregatedDuration{
				Count: 6,
				Sum:   7 * time.Second,
			},
			RepresentativeCount: 8,
		},
		NumericLabels: NumericLabels{
			"foo": {
				Values: []float64{1, 2, 3},
				Value:  1,
				Global: true,
			},
		},
		Labels: Labels{
			"bar": {
				Value:  "a",
				Values: []string{"a", "b", "c"},
				Global: true,
			},
		},
		Message: "message",
		Transaction: &Transaction{
			SpanCount: SpanCount{
				Started: uintPtr(1),
				Dropped: uintPtr(2),
			},
			UserExperience: &UserExperience{
				CumulativeLayoutShift: 1,
				FirstInputDelay:       2,
				TotalBlockingTime:     3,
				Longtask: LongtaskMetrics{
					Count: 4,
					Sum:   5,
					Max:   6,
				},
			},
			// TODO investigat valid values
			Custom: nil,
			Marks: TransactionMarks{
				"foo": TransactionMark{
					"bar": 3,
				},
			},
			Message: &Message{
				Body: "body",
				Headers: http.Header{
					"foo": []string{"bar"},
				},
				AgeMillis:  int64Ptr(2),
				QueueName:  "queuename",
				RoutingKey: "routingkey",
			},
			Type:   "type",
			Name:   "name",
			Result: "result",
			ID:     "id",
			DurationHistogram: Histogram{
				Values: []float64{4},
				Counts: []int64{5},
			},
			DroppedSpansStats: []DroppedSpanStats{
				{
					DestinationServiceResource: "destinationserviceresource",
					ServiceTargetType:          "servicetargetype",
					ServiceTargetName:          "servicetargetname",
					Outcome:                    "outcome",
					Duration: AggregatedDuration{
						Count: 4,
						Sum:   5 * time.Second,
					},
				},
			},
			DurationSummary: SummaryMetric{
				Count: 6,
				Sum:   7,
			},
			RepresentativeCount: 8,
			Sampled:             true,
			Root:                true,
		},
		Metricset: &Metricset{
			Name:     "name",
			Interval: "interval",
			Samples: []MetricsetSample{
				{
					Type: MetricTypeCounter,
					Name: "name",
					Unit: "unit",
					Histogram: Histogram{
						Values: []float64{1},
						Counts: []int64{2},
					},
					SummaryMetric: SummaryMetric{
						Count: 3,
						Sum:   4,
					},
					Value: 5,
				},
			},
			DocCount: 1,
		},
		Error: &Error{
			Exception: &Exception{
				Message:    "ex_message",
				Module:     "ex_module",
				Code:       "ex_code",
				Attributes: randomStructPb(t),
				Type:       "ex_type",
				Handled:    boolPtr(true),
				Cause: []Exception{
					{
						Message: "ex1_message",
						Module:  "ex1_module",
						Code:    "ex1_code",
						Type:    "ex_type",
					},
				},
			},
			Log: &ErrorLog{
				Message:      "log_message",
				Level:        "log_level",
				ParamMessage: "log_parammessage",
				LoggerName:   "log_loggername",
			},
			ID:          "id",
			GroupingKey: "groupingkey",
			Culprit:     "culprit",
			StackTrace:  "stacktrace",
			Message:     "message",
			Type:        "type",
		},
		Cloud: Cloud{
			Origin: &CloudOrigin{
				AccountID:   "origin_accountid",
				Provider:    "origin_provider",
				Region:      "origin_region",
				ServiceName: "origin_servicename",
			},
			AccountID:        "accountid",
			AccountName:      "accountname",
			AvailabilityZone: "availabilityzone",
			InstanceID:       "instanceid",
			InstanceName:     "instancename",
			MachineType:      "machinetype",
			ProjectID:        "projectid",
			ProjectName:      "projectname",
			Provider:         "provider",
			Region:           "region",
			ServiceName:      "servicename",
		},
		Service: Service{
			Origin: &ServiceOrigin{
				ID:      "origin_id",
				Name:    "origin_name",
				Version: "origin_version",
			},
			Target: &ServiceTarget{
				Name: "target_name",
				Type: "target_type",
			},
			Language: Language{
				Name:    "language_name",
				Version: "language_version",
			},
			Runtime: Runtime{
				Name:    "runtime_name",
				Version: "runtime_version",
			},
			Framework: Framework{
				Name:    "framework_name",
				Version: "framework_version",
			},
			Name:        "name",
			Version:     "version",
			Environment: "environment",
			Node: ServiceNode{
				Name: "node_name",
			},
		},
		FAAS: FAAS{
			ID:               "id",
			Coldstart:        boolPtr(true),
			Execution:        "execution",
			TriggerType:      "triggertype",
			TriggerRequestID: "triggerrequestid",
			Name:             "name",
			Version:          "version",
		},
		Network: Network{
			Connection: NetworkConnection{
				Type:    "type",
				Subtype: "subtype",
			},
			Carrier: NetworkCarrier{
				Name: "name",
				MCC:  "mcc",
				MNC:  "mnc",
				ICC:  "icc",
			},
		},
		Container: Container{
			ID:        "id",
			Name:      "name",
			Runtime:   "runtime",
			ImageName: "imagename",
			ImageTag:  "imagetag",
		},
		User: User{
			Domain: "domain",
			ID:     "id",
			Email:  "email",
			Name:   "name",
		},
		Device: Device{
			ID: "id",
			Model: DeviceModel{
				Name:       "name",
				Identifier: "identifier",
			},
			Manufacturer: "manufacturer",
		},
		Kubernetes: Kubernetes{
			Namespace: "namespace",
			NodeName:  "nodename",
			PodName:   "podname",
			PodUID:    "poduid",
		},
		Observer: Observer{
			Hostname: "hostname",
			Name:     "name",
			Type:     "type",
			Version:  "version",
		},
		DataStream: DataStream{
			Type:      "type",
			Dataset:   "dataset",
			Namespace: "namespace",
		},
		Agent: Agent{
			Name:             "name",
			Version:          "version",
			EphemeralID:      "ephemeralid",
			ActivationMethod: "activationmethod",
		},
		Processor: Processor{
			Name:  "name",
			Event: "event",
		},
		HTTP: HTTP{
			Request: &HTTPRequest{
				Headers:  randomStructPb(t).AsMap(),
				Env:      randomStructPb(t).AsMap(),
				Cookies:  randomStructPb(t).AsMap(),
				ID:       "id",
				Method:   "method",
				Referrer: "referrer",
			},
			Response: &HTTPResponse{
				Headers:         randomStructPb(t).AsMap(),
				Finished:        boolPtr(true),
				HeadersSent:     boolPtr(true),
				TransferSize:    int64Ptr(1),
				EncodedBodySize: int64Ptr(2),
				DecodedBodySize: int64Ptr(3),
				StatusCode:      200,
			},
			Version: "version",
		},
		UserAgent: UserAgent{
			Original: "original",
			Name:     "name",
		},
		Parent: Parent{
			ID: "id",
		},
		Trace: Trace{
			ID: "id",
		},
		Host: Host{
			OS: OS{
				Name:     "name",
				Version:  "version",
				Platform: "platform",
				Full:     "full",
				Type:     "type",
			},
			Hostname:     "hostname",
			Name:         "name",
			ID:           "id",
			Architecture: "architecture",
			Type:         "type",
			IP:           []netip.Addr{netip.MustParseAddr("127.0.0.1")},
		},
		URL: URL{
			Original: "original",
			Scheme:   "scheme",
			Full:     "full",
			Domain:   "doain",
			Path:     "path",
			Query:    "query",
			Fragment: "fragment",
			Port:     443,
		},
		Log: Log{
			Level:  "level",
			Logger: "logger",
			Origin: LogOrigin{
				FunctionName: "functionname",
				File: LogOriginFile{
					Name: "name",
					Line: 1,
				},
			},
		},
		Source: Source{
			IP: netip.MustParseAddr("127.0.0.1"),
			NAT: &NAT{
				IP: netip.MustParseAddr("127.0.0.2"),
			},
			Domain: "domain",
			Port:   443,
		},
		Client: Client{
			IP:     netip.MustParseAddr("127.0.0.1"),
			Domain: "example.com",
			Port:   443,
		},
		Child: Child{
			ID: []string{"id"},
		},
		Destination: Destination{
			Address: "127.0.0.1",
			Port:    443,
		},
		Session: Session{
			ID:       "id",
			Sequence: 1,
		},
		Process: Process{
			Ppid: 1,
			Thread: ProcessThread{
				Name: "name",
				ID:   2,
			},
			Title:       "title",
			CommandLine: "commandline",
			Executable:  "executable",
			Argv:        []string{"argv"},
			Pid:         3,
		},
		Event: Event{
			Outcome:  "outcome",
			Action:   "action",
			Dataset:  "dataset",
			Kind:     "kind",
			Category: "category",
			Type:     "type",
			SuccessCount: SummaryMetric{
				Count: 1,
				Sum:   2,
			},
			Duration: 3 * time.Second,
			Severity: 4,
		},
	}
}

func BenchmarkMarshalFastJSON(b *testing.B) {
	input := APMEvent{
		Agent: Agent{
			Name:    "agent_name",
			Version: "agent_version",
		},
		Observer:  Observer{Type: "observer_type"},
		Container: Container{ID: "container_id"},
		Service: Service{
			Name: "service_name",
			Node: ServiceNode{Name: "service_node_name"},
			Origin: &ServiceOrigin{
				ID:      "abc123",
				Name:    "service_origin_name",
				Version: "1.0",
			},
		},
		Host: Host{
			Hostname: "hostname",
			Name:     "host_name",
		},
		Client: Client{Domain: "client.domain"},
		Source: Source{
			IP:   netip.MustParseAddr("127.0.0.1"),
			Port: 1234,
			NAT: &NAT{
				IP: netip.MustParseAddr("10.10.10.10"),
			},
		},
		Destination: Destination{Address: "destination_address", Port: 1234},
		Process:     Process{Pid: 1234},
		User:        User{ID: "user_id", Email: "user_email"},
		Event:       Event{Outcome: "success", Duration: time.Microsecond},
		Session:     Session{ID: "session_id"},
		URL:         URL{Original: "url"},
		Labels: map[string]LabelValue{
			"a": {Value: "b"},
			"c": {Value: "true"},
			"d": {Values: []string{"true", "false"}},
		},
		NumericLabels: map[string]NumericLabelValue{
			"e": {Value: float64(1234)},
			"f": {Values: []float64{1234, 12311}},
		},
		Message:   "bottle",
		Timestamp: time.Date(2019, 1, 3, 15, 17, 4, 908.596*1e6, time.FixedZone("+0100", 3600)),
		Processor: Processor{Name: "processor_name", Event: "processor_event"},
		Trace:     Trace{ID: "trace_id"},
		Parent:    Parent{ID: "parent_id"},
		Child:     Child{ID: []string{"child_1", "child_2"}},
		HTTP: HTTP{
			Request: &HTTPRequest{
				Method: "GET",
				Body:   "foo bar baz",
			},
		},
		FAAS: FAAS{
			ID:               "faasID",
			Execution:        "execution",
			TriggerType:      "http",
			TriggerRequestID: "abc123",
			Name:             "faasName",
			Version:          "1.0.0",
		},
		Cloud: Cloud{
			Origin: &CloudOrigin{
				AccountID:   "accountID",
				Provider:    "aws",
				Region:      "us-west-1",
				ServiceName: "serviceName",
			},
		},
	}

	var w fastjson.Writer
	for i := 0; i < b.N; i++ {
		w.Reset()
		if err := input.MarshalFastJSON(&w); err != nil {
			b.Fatal(err)
		}
	}
}

func transformAPMEvent(in APMEvent) map[string]any {
	var decoded map[string]any
	encoded := marshalJSONAPMEvent(in)
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		panic(fmt.Errorf("failed to decode %q: %w", encoded, err))
	}
	return decoded
}

func marshalJSONAPMEvent(in APMEvent) []byte {
	out, err := in.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return out
}
