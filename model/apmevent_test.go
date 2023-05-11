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
