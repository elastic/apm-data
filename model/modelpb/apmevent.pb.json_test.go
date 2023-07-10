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
	"testing"
	"time"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BenchmarkAPMEventToJSON(b *testing.B) {
	event := fullEvent(b)

	b.Run("to-json", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			event.MarshalJSON()
		}
	})
}

func fullEvent(t testing.TB) *APMEvent {
	return &APMEvent{
		Timestamp: timestamppb.New(time.Unix(1, 1)),
		Span: &Span{
			Message: &Message{
				Body: "body",
				Headers: []*HTTPHeader{
					{
						Key:   "foo",
						Value: []string{"bar"},
					},
				},
				AgeMillis:  int64Ptr(2),
				QueueName:  "queuename",
				RoutingKey: "routingkey",
			},
			Composite: &Composite{
				CompressionStrategy: CompressionStrategy_COMPRESSION_STRATEGY_EXACT_MATCH,
				Count:               1,
				Sum:                 2,
			},
			DestinationService: &DestinationService{
				Type:     "destination_type",
				Name:     "destination_name",
				Resource: "destination_resource",
				ResponseTime: &AggregatedDuration{
					Count: 3,
					Sum:   durationpb.New(4 * time.Second),
				},
			},
			Db: &DB{
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
			Id:      "id",
			Type:    "type",
			Name:    "name",
			Stacktrace: []*StacktraceFrame{
				{
					Vars:           randomStructPb(t),
					Lineno:         uintPtr(1),
					Colno:          uintPtr(2),
					Filename:       "frame_filename",
					Classname:      "frame_classname",
					ContextLine:    "frame_contextline",
					Module:         "frame_module",
					Function:       "frame_function",
					AbsPath:        "frame_abspath",
					SourcemapError: "frame_sourcemaperror",
					Original: &Original{
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
			Links: []*SpanLink{
				{
					TraceId: "trace_id",
					SpanId:  "id1",
				},
			},
			SelfTime: &AggregatedDuration{
				Count: 6,
				Sum:   durationpb.New(7 * time.Second),
			},
			RepresentativeCount: 8,
		},
		NumericLabels: map[string]*NumericLabelValue{
			"foo": {
				Values: []float64{1, 2, 3},
				Value:  1,
				Global: true,
			},
		},
		Labels: map[string]*LabelValue{
			"bar": {
				Value:  "a",
				Values: []string{"a", "b", "c"},
				Global: true,
			},
		},
		Message: "message",
		Transaction: &Transaction{
			SpanCount: &SpanCount{
				Started: uintPtr(1),
				Dropped: uintPtr(2),
			},
			UserExperience: &UserExperience{
				CumulativeLayoutShift: 1,
				FirstInputDelay:       2,
				TotalBlockingTime:     3,
				LongTask: &LongtaskMetrics{
					Count: 4,
					Sum:   5,
					Max:   6,
				},
			},
			// TODO investigat valid values
			Custom: nil,
			Marks: map[string]*TransactionMark{
				"foo": {
					Measurements: map[string]float64{
						"bar": 3,
					},
				},
			},
			Message: &Message{
				Body: "body",
				Headers: []*HTTPHeader{
					{
						Key:   "foo",
						Value: []string{"bar"},
					},
				},
				AgeMillis:  int64Ptr(2),
				QueueName:  "queuename",
				RoutingKey: "routingkey",
			},
			Type:   "type",
			Name:   "name",
			Result: "result",
			Id:     "id",
			DurationHistogram: &Histogram{
				Values: []float64{4},
				Counts: []int64{5},
			},
			DroppedSpansStats: []*DroppedSpanStats{
				{
					DestinationServiceResource: "destinationserviceresource",
					ServiceTargetType:          "servicetargetype",
					ServiceTargetName:          "servicetargetname",
					Outcome:                    "outcome",
					Duration: &AggregatedDuration{
						Count: 4,
						Sum:   durationpb.New(5 * time.Second),
					},
				},
			},
			DurationSummary: &SummaryMetric{
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
			Samples: []*MetricsetSample{
				{
					Type: MetricType_METRIC_TYPE_COUNTER,
					Name: "name",
					Unit: "unit",
					Histogram: &Histogram{
						Values: []float64{1},
						Counts: []int64{2},
					},
					Summary: &SummaryMetric{
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
				Cause: []*Exception{
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
			Id:          "id",
			GroupingKey: "groupingkey",
			Culprit:     "culprit",
			StackTrace:  "stacktrace",
			Message:     "message",
			Type:        "type",
		},
		Cloud: &Cloud{
			Origin: &CloudOrigin{
				AccountId:   "origin_accountid",
				Provider:    "origin_provider",
				Region:      "origin_region",
				ServiceName: "origin_servicename",
			},
			AccountId:        "accountid",
			AccountName:      "accountname",
			AvailabilityZone: "availabilityzone",
			InstanceId:       "instanceid",
			InstanceName:     "instancename",
			MachineType:      "machinetype",
			ProjectId:        "projectid",
			ProjectName:      "projectname",
			Provider:         "provider",
			Region:           "region",
			ServiceName:      "servicename",
		},
		Service: &Service{
			Origin: &ServiceOrigin{
				Id:      "origin_id",
				Name:    "origin_name",
				Version: "origin_version",
			},
			Target: &ServiceTarget{
				Name: "target_name",
				Type: "target_type",
			},
			Language: &Language{
				Name:    "language_name",
				Version: "language_version",
			},
			Runtime: &Runtime{
				Name:    "runtime_name",
				Version: "runtime_version",
			},
			Framework: &Framework{
				Name:    "framework_name",
				Version: "framework_version",
			},
			Name:        "name",
			Version:     "version",
			Environment: "environment",
			Node: &ServiceNode{
				Name: "node_name",
			},
		},
		Faas: &Faas{
			Id:               "id",
			ColdStart:        boolPtr(true),
			Execution:        "execution",
			TriggerType:      "triggertype",
			TriggerRequestId: "triggerrequestid",
			Name:             "name",
			Version:          "version",
		},
		Network: &Network{
			Connection: &NetworkConnection{
				Type:    "type",
				Subtype: "subtype",
			},
			Carrier: &NetworkCarrier{
				Name: "name",
				Mcc:  "mcc",
				Mnc:  "mnc",
				Icc:  "icc",
			},
		},
		Container: &Container{
			Id:        "id",
			Name:      "name",
			Runtime:   "runtime",
			ImageName: "imagename",
			ImageTag:  "imagetag",
		},
		User: &User{
			Domain: "domain",
			Id:     "id",
			Email:  "email",
			Name:   "name",
		},
		Device: &Device{
			Id: "id",
			Model: &DeviceModel{
				Name:       "name",
				Identifier: "identifier",
			},
			Manufacturer: "manufacturer",
		},
		Kubernetes: &Kubernetes{
			Namespace: "namespace",
			NodeName:  "nodename",
			PodName:   "podname",
			PodUid:    "poduid",
		},
		Observer: &Observer{
			Hostname: "hostname",
			Name:     "name",
			Type:     "type",
			Version:  "version",
		},
		DataStream: &DataStream{
			Type:      "type",
			Dataset:   "dataset",
			Namespace: "namespace",
		},
		Agent: &Agent{
			Name:             "name",
			Version:          "version",
			EphemeralId:      "ephemeralid",
			ActivationMethod: "activationmethod",
		},
		Processor: &Processor{
			Name:  "name",
			Event: "event",
		},
		Http: &HTTP{
			Request: &HTTPRequest{
				Headers:  randomHTTPHeaders(t),
				Env:      randomStructPb(t),
				Cookies:  randomStructPb(t),
				Id:       "id",
				Method:   "method",
				Referrer: "referrer",
			},
			Response: &HTTPResponse{
				Headers:         randomHTTPHeaders(t),
				Finished:        boolPtr(true),
				HeadersSent:     boolPtr(true),
				TransferSize:    int64Ptr(1),
				EncodedBodySize: int64Ptr(2),
				DecodedBodySize: int64Ptr(3),
				StatusCode:      200,
			},
			Version: "version",
		},
		UserAgent: &UserAgent{
			Original: "original",
			Name:     "name",
		},
		ParentId: "id",
		Trace: &Trace{
			Id: "id",
		},
		Host: &Host{
			Os: &OS{
				Name:     "name",
				Version:  "version",
				Platform: "platform",
				Full:     "full",
				Type:     "type",
			},
			Hostname:     "hostname",
			Name:         "name",
			Id:           "id",
			Architecture: "architecture",
			Type:         "type",
			Ip: []*IP{
				{
					IpAddr: &IP_V4{2130706433},
				},
			},
		},
		Url: &URL{
			Original: "original",
			Scheme:   "scheme",
			Full:     "full",
			Domain:   "doain",
			Path:     "path",
			Query:    "query",
			Fragment: "fragment",
			Port:     443,
		},
		Log: &Log{
			Level:  "level",
			Logger: "logger",
			Origin: &LogOrigin{
				FunctionName: "functionname",
				File: &LogOriginFile{
					Name: "name",
					Line: 1,
				},
			},
		},
		Source: &Source{
			Ip: &IP{
				IpAddr: &IP_V4{2130706433},
			},
			Nat: &NAT{
				Ip: &IP{
					IpAddr: &IP_V4{2130706434},
				},
			},
			Domain: "domain",
			Port:   443,
		},
		Client: &Client{
			Ip: &IP{
				IpAddr: &IP_V4{2130706433},
			},
			Domain: "example.com",
			Port:   443,
		},
		ChildIds: []string{"id"},
		Destination: &Destination{
			Address: "127.0.0.1",
			Port:    443,
		},
		Session: &Session{
			Id:       "id",
			Sequence: 1,
		},
		Process: &Process{
			Ppid: 1,
			Thread: &ProcessThread{
				Name: "name",
				Id:   2,
			},
			Title:       "title",
			CommandLine: "commandline",
			Executable:  "executable",
			Argv:        []string{"argv"},
			Pid:         3,
		},
		Event: &Event{
			Outcome:  "outcome",
			Action:   "action",
			Dataset:  "dataset",
			Kind:     "kind",
			Category: "category",
			Type:     "type",
			SuccessCount: &SummaryMetric{
				Count: 1,
				Sum:   2,
			},
			Duration: durationpb.New(3 * time.Second),
			Severity: 4,
		},
	}
}
