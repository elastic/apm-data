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

// Portions copied from OpenTelemetry Collector (contrib), from the
// elastic exporter.
//
// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlp

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv12 "go.opentelemetry.io/otel/semconv/v1.12.0"
	semconv16 "go.opentelemetry.io/otel/semconv/v1.16.0"
	semconv18 "go.opentelemetry.io/otel/semconv/v1.18.0"
	semconv25 "go.opentelemetry.io/otel/semconv/v1.25.0"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"google.golang.org/grpc/codes"

	"github.com/elastic/apm-data/model/modelpb"
)

const (
	keywordLength = 1024
	dot           = "."
	underscore    = "_"

	outcomeSuccess = "success"
	outcomeFailure = "failure"
	outcomeUnknown = "unknown"

	attributeHTTPPath                   = "http.path"
	attributeDbElasticsearchClusterName = "db.elasticsearch.cluster.name"
	attributeDataStreamDataset          = "data_stream.dataset"
	attributeDataStreamNamespace        = "data_stream.namespace"
)

// ConsumeTracesResult contains the number of rejected spans and error message for partial success response.
type ConsumeTracesResult struct {
	ErrorMessage  string
	RejectedSpans int64
}

// ConsumeTraces calls ConsumeTracesWithResult but ignores the result.
// It exists to satisfy the go.opentelemetry.io/collector/consumer.Traces interface.
func (c *Consumer) ConsumeTraces(ctx context.Context, traces ptrace.Traces) error {
	_, err := c.ConsumeTracesWithResult(ctx, traces)
	return err
}

// ConsumeTracesWithResult consumes OpenTelemetry trace data,
// converting into Elastic APM events and reporting to the Elastic APM schema.
func (c *Consumer) ConsumeTracesWithResult(ctx context.Context, traces ptrace.Traces) (ConsumeTracesResult, error) {
	if err := semAcquire(ctx, c.tracer, c.sem, 1); err != nil {
		return ConsumeTracesResult{}, err
	}
	defer c.sem.Release(1)

	receiveTimestamp := time.Now()

	resourceSpans := traces.ResourceSpans()
	batch := make(modelpb.Batch, 0, resourceSpans.Len())
	for i := 0; i < resourceSpans.Len(); i++ {
		c.convertResourceSpans(resourceSpans.At(i), receiveTimestamp, &batch)
	}
	if err := c.config.Processor.ProcessBatch(ctx, &batch); err != nil {
		return ConsumeTracesResult{}, err
	}
	return ConsumeTracesResult{RejectedSpans: 0}, nil
}

func (c *Consumer) convertResourceSpans(
	resourceSpans ptrace.ResourceSpans,
	receiveTimestamp time.Time,
	out *modelpb.Batch,
) {
	baseEvent := modelpb.APMEvent{}
	baseEvent.Event = &modelpb.Event{}
	baseEvent.Event.Received = modelpb.FromTime(receiveTimestamp)
	var timeDelta time.Duration
	resource := resourceSpans.Resource()
	translateResourceMetadata(resource, &baseEvent)
	if exportTimestamp, ok := exportTimestamp(resource); ok {
		timeDelta = receiveTimestamp.Sub(exportTimestamp)
	}
	scopeSpans := resourceSpans.ScopeSpans()
	for i := 0; i < scopeSpans.Len(); i++ {
		c.convertScopeSpans(scopeSpans.At(i), &baseEvent, timeDelta, out)
	}
}

func (c *Consumer) convertScopeSpans(
	in ptrace.ScopeSpans,
	baseEvent *modelpb.APMEvent,
	timeDelta time.Duration,
	out *modelpb.Batch,
) {
	otelSpans := in.Spans()
	for i := 0; i < otelSpans.Len(); i++ {
		c.convertSpan(otelSpans.At(i), in.Scope(), baseEvent, timeDelta, out)
	}
}

func (c *Consumer) convertSpan(
	otelSpan ptrace.Span,
	otelLibrary pcommon.InstrumentationScope,
	baseEvent *modelpb.APMEvent,
	timeDelta time.Duration,
	out *modelpb.Batch,
) {
	root := otelSpan.ParentSpanID().IsEmpty()
	var parentID string
	if !root {
		parentID = hexSpanID(otelSpan.ParentSpanID())
	}

	startTime := otelSpan.StartTimestamp().AsTime()
	endTime := otelSpan.EndTimestamp().AsTime()
	duration := endTime.Sub(startTime)

	// Message consumption results in either a transaction or a span based
	// on whether the consumption is active or passive. Otel spans
	// currently do not have the metadata to make this distinction. For
	// now, we assume that the majority of consumption is passive, and
	// therefore start a transaction whenever span kind == consumer.
	name := otelSpan.Name()
	spanID := hexSpanID(otelSpan.SpanID())
	representativeCount := getRepresentativeCountFromTracestateHeader(otelSpan.TraceState().AsRaw())
	event := baseEvent.CloneVT()

	translateScopeMetadata(otelLibrary, event)

	initEventLabels(event)
	event.Timestamp = modelpb.FromTime(startTime.Add(timeDelta))
	if id := hexTraceID(otelSpan.TraceID()); id != "" {
		event.Trace = &modelpb.Trace{}
		event.Trace.Id = id
	}
	if event.Event == nil {
		event.Event = &modelpb.Event{}
	}
	event.Event.Duration = uint64(duration)
	event.Event.Outcome = spanStatusOutcome(otelSpan.Status())
	if parentID != "" {
		event.ParentId = parentID
	}
	if root || otelSpan.Kind() == ptrace.SpanKindServer || otelSpan.Kind() == ptrace.SpanKindConsumer {
		event.Transaction = &modelpb.Transaction{}
		event.Transaction.Id = spanID
		event.Transaction.Name = name
		event.Transaction.Sampled = true
		event.Transaction.RepresentativeCount = representativeCount
		if spanID != "" {
			event.Span = &modelpb.Span{}
			event.Span.Id = spanID
		}

		TranslateTransaction(otelSpan.Attributes(), otelSpan.Status(), otelLibrary, event)
	} else {
		event.Span = &modelpb.Span{}
		event.Span.Id = spanID
		event.Span.Name = name
		event.Span.RepresentativeCount = representativeCount
		TranslateSpan(otelSpan.Kind(), otelSpan.Attributes(), event)
	}
	translateSpanLinks(event, otelSpan.Links())
	if len(event.Labels) == 0 {
		event.Labels = nil
	}
	if len(event.NumericLabels) == 0 {
		event.NumericLabels = nil
	}
	*out = append(*out, event)

	events := otelSpan.Events()
	event = event.CloneVT()
	event.Labels = baseEvent.Labels               // only copy common labels to span events
	event.NumericLabels = baseEvent.NumericLabels // only copy common labels to span events
	event.Event = &modelpb.Event{}
	event.Event.Received = baseEvent.Event.Received // only copy event.received to span events
	event.Destination = nil                         // don't set destination for span events
	for i := 0; i < events.Len(); i++ {
		*out = append(*out, c.convertSpanEvent(events.At(i), event, timeDelta))
	}
}

// TranslateTransaction converts incoming otlp/otel trace data into the
// expected elasticsearch format.
func TranslateTransaction(
	attributes pcommon.Map,
	spanStatus ptrace.Status,
	library pcommon.InstrumentationScope,
	event *modelpb.APMEvent,
) {
	isJaeger := strings.HasPrefix(event.Agent.Name, "Jaeger")

	var (
		netHostName string
		netHostPort int
	)

	var (
		httpScheme     string
		httpURL        string
		httpServerName string
		httpHost       string
		http           modelpb.HTTP
		httpRequest    modelpb.HTTPRequest
		httpResponse   modelpb.HTTPResponse
		urlPath        string
		urlQuery       string
	)

	var isHTTP, isRPC, isMessaging bool
	var messagingQueueName string

	var samplerType, samplerParam pcommon.Value
	attributes.Range(func(kDots string, v pcommon.Value) bool {
		if isJaeger {
			switch kDots {
			case "sampler.type":
				samplerType = v
				return true
			case "sampler.param":
				samplerParam = v
				return true
			}
		}

		k := replaceDots(kDots)
		switch v.Type() {
		case pcommon.ValueTypeSlice:
			switch kDots {
			case "elastic.profiler_stack_trace_ids":
				var vSlice = v.Slice()
				event.Transaction.ProfilerStackTraceIds = slices.Grow(event.Transaction.ProfilerStackTraceIds, vSlice.Len())
				for i := 0; i < vSlice.Len(); i++ {
					var idVal = vSlice.At(i)
					if idVal.Type() == pcommon.ValueTypeStr {
						event.Transaction.ProfilerStackTraceIds = append(event.Transaction.ProfilerStackTraceIds, idVal.Str())
					}
				}
			default:
				setLabel(k, event, v)
			}
		case pcommon.ValueTypeInt:
			switch kDots {
			case string(semconv25.HTTPStatusCodeKey), string(semconv.HTTPResponseStatusCodeKey):
				isHTTP = true
				httpResponse.StatusCode = uint32(v.Int())
				http.Response = &httpResponse
			case string(semconv25.NetPeerPortKey), string(semconv.ClientPortKey):
				if event.Source == nil {
					event.Source = &modelpb.Source{}
				}
				event.Source.Port = uint32(v.Int())
			case string(semconv25.NetHostPortKey), string(semconv.ServerPortKey):
				netHostPort = int(v.Int())
			case string(semconv.RPCGRPCStatusCodeKey):
				isRPC = true
				event.Transaction.Result = codes.Code(v.Int()).String()
			default:
				setLabel(k, event, v)
			}
		case pcommon.ValueTypeMap:
		case pcommon.ValueTypeStr:
			stringval := truncate(v.Str())
			switch kDots {
			// http.*
			case string(semconv25.HTTPMethodKey), string(semconv.HTTPRequestMethodKey):
				isHTTP = true
				httpRequest.Method = stringval
				http.Request = &httpRequest
			case string(semconv25.HTTPURLKey), string(semconv.URLFullKey), string(semconv25.HTTPTargetKey), attributeHTTPPath:
				isHTTP = true
				httpURL = stringval
			case string(semconv.URLPathKey):
				isHTTP = true
				urlPath = stringval
			case string(semconv.URLQueryKey):
				isHTTP = true
				urlQuery = stringval
			case string(semconv12.HTTPHostKey): //removed after 1.12 (stable)
				isHTTP = true
				httpHost = stringval
			case string(semconv25.HTTPSchemeKey), string(semconv.URLSchemeKey):
				isHTTP = true
				httpScheme = stringval
			case string(semconv25.HTTPStatusCodeKey), string(semconv.HTTPResponseStatusCodeKey):
				if intv, err := strconv.Atoi(stringval); err == nil {
					isHTTP = true
					httpResponse.StatusCode = uint32(intv)
					http.Response = &httpResponse
				}
			case "http.protocol":
				if !strings.HasPrefix(stringval, "HTTP/") {
					// Unexpected, store in labels for debugging.
					modelpb.Labels(event.Labels).Set(k, stringval)
					break
				}
				stringval = strings.TrimPrefix(stringval, "HTTP/")
				fallthrough
			case string(semconv25.HTTPFlavorKey): //removed after 1.25 (experimental)
				isHTTP = true
				http.Version = stringval
			case string(semconv12.HTTPServerNameKey): //removed after 1.12 (stable)
				isHTTP = true
				httpServerName = stringval
			case string(semconv18.HTTPClientIPKey): //removed after 1.18 (stable)
				if ip, err := modelpb.ParseIP(stringval); err == nil {
					if event.Client == nil {
						event.Client = &modelpb.Client{}
					}
					event.Client.Ip = ip
				}
			case string(semconv25.HTTPUserAgentKey), string(semconv.UserAgentOriginalKey):
				if event.UserAgent == nil {
					event.UserAgent = &modelpb.UserAgent{}
				}
				event.UserAgent.Original = stringval

			// net.*
			case string(semconv12.NetPeerIPKey), string(semconv25.NetSockPeerAddrKey):
				if event.Source == nil {
					event.Source = &modelpb.Source{}
				}
				if ip, err := modelpb.ParseIP(stringval); err == nil {
					event.Source.Ip = ip
				}
			case string(semconv25.NetPeerNameKey), string(semconv.ClientAddressKey):
				if event.Source == nil {
					event.Source = &modelpb.Source{}
				}
				event.Source.Domain = stringval
			case string(semconv25.NetHostNameKey), string(semconv.ServerAddressKey):
				netHostName = stringval
			case string(semconv.NetworkConnectionTypeKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Connection == nil {
					event.Network.Connection = &modelpb.NetworkConnection{}
				}
				event.Network.Connection.Type = stringval
			case string(semconv.NetworkConnectionSubtypeKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Connection == nil {
					event.Network.Connection = &modelpb.NetworkConnection{}
				}
				event.Network.Connection.Subtype = stringval
			case string(semconv.NetworkCarrierMccKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Mcc = stringval
			case string(semconv.NetworkCarrierMncKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Mnc = stringval
			case string(semconv.NetworkCarrierNameKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Name = stringval
			case string(semconv.NetworkCarrierIccKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Icc = stringval

			// messaging.*
			//
			// messaging.destination is now called messaging.destination.name in the latest semconv
			// https://opentelemetry.io/docs/specs/semconv/attributes-registry/messaging
			// keep both of them for the backward compatibility
			case string(semconv16.MessagingDestinationKey), string(semconv.MessagingDestinationNameKey), "message_bus.destination":
				isMessaging = true
				messagingQueueName = stringval
			case string(semconv.MessagingSystemKey):
				isMessaging = true
				modelpb.Labels(event.Labels).Set(k, stringval)
			case string(semconv25.MessagingOperationKey), string(semconv.MessagingOperationTypeKey), string(semconv.MessagingOperationNameKey):
				isMessaging = true
				modelpb.Labels(event.Labels).Set(k, stringval)

			// rpc.*
			//
			// TODO(axw) add RPC fieldset to ECS? Currently we drop these
			// attributes, and rely on the operation name like we do with
			// Elastic APM agents.
			case string(semconv.RPCSystemKey):
				isRPC = true
			case string(semconv.RPCGRPCStatusCodeKey):
				isRPC = true
			case string(semconv.RPCServiceKey):
			case string(semconv.RPCMethodKey):

			// miscellaneous
			case "type":
				event.Transaction.Type = stringval
			case "session.id":
				if event.Session == nil {
					event.Session = &modelpb.Session{}
				}
				event.Session.Id = stringval
			case string(semconv.ServiceVersionKey):
				// NOTE support for sending service.version as a span tag
				// is deprecated, and will be removed in 8.0. Instrumentation
				// should set this as a resource attribute (OTel) or tracer
				// tag (Jaeger).
				event.Service.Version = stringval

			// data_stream.*
			case attributeDataStreamDataset:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Dataset = sanitizeDataStreamDataset(stringval)
			case attributeDataStreamNamespace:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Namespace = sanitizeDataStreamNamespace(stringval)
			default:
				modelpb.Labels(event.Labels).Set(k, stringval)
			}
		default:
			setLabel(k, event, v)
		}
		return true
	})

	if event.Transaction.Type == "" {
		switch {
		case isMessaging:
			event.Transaction.Type = "messaging"
		case isHTTP, isRPC:
			event.Transaction.Type = "request"
		default:
			event.Transaction.Type = "unknown"
		}
	}

	if isHTTP {
		if http.SizeVT() != 0 {
			event.Http = &http
		}

		// Set outcome nad result from status code.
		if statusCode := httpResponse.StatusCode; statusCode > 0 {
			if event.Event.Outcome == outcomeUnknown {
				event.Event.Outcome = serverHTTPStatusCodeOutcome(int(statusCode))
			}
			if event.Transaction.Result == "" {
				event.Transaction.Result = httpStatusCodeResult(int(statusCode))
			}
		}

		httpHost := httpHost
		if httpHost == "" {
			httpHost = httpServerName
			if httpHost == "" {
				httpHost = netHostName
				if httpHost == "" {
					httpHost = event.GetHost().GetHostname()
				}
			}
			if httpHost != "" && netHostPort > 0 {
				httpHost = net.JoinHostPort(httpHost, strconv.Itoa(netHostPort))
			}
		}

		// Build a relative url from the UrlPath and UrlQuery.
		httpURL := httpURL
		if httpURL == "" && urlPath != "" {
			httpURL = urlPath
			if urlQuery != "" {
				httpURL += "?" + urlQuery
			}
		}

		// Build the modelpb.URL from http{URL,Host,Scheme}.
		event.Url = modelpb.ParseURL(httpURL, httpHost, httpScheme)
	}

	if isMessaging {
		// Overwrite existing event.Transaction.Message
		event.Transaction.Message = nil
		if messagingQueueName != "" {
			event.Transaction.Message = &modelpb.Message{}
			event.Transaction.Message.QueueName = messagingQueueName
		}
	}

	if event.Client == nil && event.Source != nil {
		event.Client = &modelpb.Client{}
		event.Client.Ip = event.Source.Ip
		event.Client.Port = event.Source.Port
		event.Client.Domain = event.Source.Domain
	}

	if samplerType != (pcommon.Value{}) {
		// The client has reported its sampling rate, so we can use it to extrapolate span metrics.
		parseSamplerAttributes(samplerType, samplerParam, event)
	}

	if event.Transaction.Result == "" {
		event.Transaction.Result = spanStatusResult(spanStatus)
	}

	// if outcome and result are still not assigned, assign success
	if event.Event.Outcome == outcomeUnknown {
		event.Event.Outcome = outcomeSuccess
		if event.Transaction.Result == "" {
			event.Transaction.Result = "Success"
		}
	}
}

// TranslateSpan converts incoming otlp/otel trace data into the
// expected elasticsearch format.
func TranslateSpan(spanKind ptrace.SpanKind, attributes pcommon.Map, event *modelpb.APMEvent) {
	isJaeger := strings.HasPrefix(event.GetAgent().GetName(), "Jaeger")

	var (
		netPeerName string
		netPeerIP   string
		netPeerPort int
	)

	var (
		peerService string
		peerAddress string
	)

	var (
		httpURL      string
		httpHost     string
		httpTarget   string
		httpURLPath  string
		httpURLQuery string
		httpScheme   = "http"
	)

	var (
		messageSystem          string
		messageOperation       string
		messageTempDestination bool
	)

	var (
		rpcSystem  string
		rpcService string
	)

	var (
		genAiSystem string
	)

	var http modelpb.HTTP
	var httpRequest modelpb.HTTPRequest
	var httpResponse modelpb.HTTPResponse
	var message modelpb.Message
	var db modelpb.DB
	var destinationService modelpb.DestinationService
	var serviceTarget modelpb.ServiceTarget
	var isHTTP, isDatabase, isRPC, isMessaging, isGenAi bool
	var samplerType, samplerParam pcommon.Value
	attributes.Range(func(kDots string, v pcommon.Value) bool {
		if isJaeger {
			switch kDots {
			case "sampler.type":
				samplerType = v
				return true
			case "sampler.param":
				samplerParam = v
				return true
			}
		}

		k := replaceDots(kDots)
		switch v.Type() {
		case pcommon.ValueTypeBool:
			switch kDots {
			case string(semconv16.MessagingTempDestinationKey), string(semconv.MessagingDestinationTemporaryKey):
				messageTempDestination = v.Bool()
				fallthrough
			default:
				setLabel(k, event, v)
			}
		case pcommon.ValueTypeInt:
			switch kDots {
			case string(semconv25.HTTPStatusCodeKey), string(semconv.HTTPResponseStatusCodeKey):
				httpResponse.StatusCode = uint32(v.Int())
				http.Response = &httpResponse
				isHTTP = true
			case string(semconv25.NetPeerPortKey), string(semconv.ServerPortKey), "peer.port":
				netPeerPort = int(v.Int())
			case string(semconv.RPCGRPCStatusCodeKey):
				rpcSystem = "grpc"
				isRPC = true
			default:
				setLabel(k, event, v)
			}
		case pcommon.ValueTypeStr:
			stringval := truncate(v.Str())

			switch kDots {
			// http.*
			case string(semconv12.HTTPHostKey): //removed after 1.12 (stable)
				httpHost = stringval
				isHTTP = true
			case string(semconv25.HTTPSchemeKey), string(semconv.URLSchemeKey):
				httpScheme = stringval
				isHTTP = true
			case string(semconv25.HTTPTargetKey):
				httpTarget = stringval
				isHTTP = true
			case string(semconv.URLPathKey):
				httpURLPath = stringval
				isHTTP = true
			case string(semconv.URLQueryKey):
				httpURLQuery = stringval
				isHTTP = true
			case string(semconv25.HTTPURLKey), string(semconv.URLFullKey):
				httpURL = stringval
				isHTTP = true
			case string(semconv25.HTTPMethodKey), string(semconv.HTTPRequestMethodKey):
				httpRequest.Method = stringval
				http.Request = &httpRequest
				isHTTP = true

			// db.*
			case "sql.query":
				if db.Type == "" {
					db.Type = "sql"
				}
				fallthrough
			case string(semconv25.DBStatementKey), string(semconv.DBQueryTextKey):
				// Statement should not be truncated, use original string value.
				db.Statement = v.Str()
				isDatabase = true
			case string(semconv25.DBNameKey), string(semconv.DBNamespaceKey), "db.instance", attributeDbElasticsearchClusterName:
				db.Instance = stringval
				isDatabase = true
			case string(semconv.DBSystemKey), "db.type":
				db.Type = stringval
				isDatabase = true
			case string(semconv25.DBUserKey): //removed after 1.25 (experimental)
				db.UserName = stringval
				isDatabase = true

			// net.*
			case string(semconv25.NetPeerNameKey), "peer.hostname", string(semconv.ServerAddressKey):
				netPeerName = stringval
			case string(semconv12.NetPeerIPKey), string(semconv25.NetSockPeerAddrKey), string(semconv.NetworkPeerAddressKey), "peer.ipv4", "peer.ipv6":
				netPeerIP = stringval
			case "peer.address":
				peerAddress = stringval
			case string(semconv.NetworkConnectionTypeKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Connection == nil {
					event.Network.Connection = &modelpb.NetworkConnection{}
				}
				event.Network.Connection.Type = stringval
			case string(semconv.NetworkConnectionSubtypeKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Connection == nil {
					event.Network.Connection = &modelpb.NetworkConnection{}
				}
				event.Network.Connection.Subtype = stringval
			case string(semconv.NetworkCarrierMccKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Mcc = stringval
			case string(semconv.NetworkCarrierMncKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Mnc = stringval
			case string(semconv.NetworkCarrierNameKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Name = stringval
			case string(semconv.NetworkCarrierIccKey):
				if event.Network == nil {
					event.Network = &modelpb.Network{}
				}
				if event.Network.Carrier == nil {
					event.Network.Carrier = &modelpb.NetworkCarrier{}
				}
				event.Network.Carrier.Icc = stringval

			// session.*
			case "session.id":
				if event.Session == nil {
					event.Session = &modelpb.Session{}
				}
				event.Session.Id = stringval

			// messaging.*
			//
			// messaging.destination is now called messaging.destination.name in the latest semconv
			// https://opentelemetry.io/docs/specs/semconv/attributes-registry/messaging
			// keep both of them for the backward compatibility
			case string(semconv16.MessagingDestinationKey), string(semconv.MessagingDestinationNameKey), "message_bus.destination":
				message.QueueName = stringval
				isMessaging = true
			case string(semconv25.MessagingOperationKey), string(semconv.MessagingOperationTypeKey):
				messageOperation = stringval
				isMessaging = true
			case string(semconv.MessagingOperationNameKey):
				modelpb.Labels(event.Labels).Set(k, stringval)
				isMessaging = true
			case string(semconv.MessagingSystemKey):
				messageSystem = stringval
				isMessaging = true

			// rpc.*
			//
			// TODO(axw) add RPC fieldset to ECS? Currently we drop these
			// attributes, and rely on the operation name and span type/subtype
			// like we do with Elastic APM agents.
			case string(semconv.RPCSystemKey):
				rpcSystem = stringval
				isRPC = true
			case string(semconv.RPCServiceKey):
				rpcService = stringval
				isRPC = true
			case string(semconv.RPCGRPCStatusCodeKey):
				rpcSystem = "grpc"
				isRPC = true
			case string(semconv.RPCMethodKey):

			case string(semconv.CodeStacktraceKey):
				if event.Code == nil {
					event.Code = &modelpb.Code{}
				}
				// stacktrace is expected to be large thus un-truncated value is needed
				event.Code.Stacktrace = v.Str()

			// gen_ai.*
			case string(semconv.GenAISystemKey):
				genAiSystem = stringval
				isGenAi = true

			// miscellaneous
			case "span.kind": // filter out
			case string(semconv.PeerServiceKey):
				peerService = stringval

			// data_stream.*
			case attributeDataStreamDataset:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Dataset = sanitizeDataStreamDataset(stringval)
			case attributeDataStreamNamespace:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Namespace = sanitizeDataStreamNamespace(stringval)
			default:
				setLabel(k, event, v)
			}
		default:
			setLabel(k, event, v)
		}
		return true
	})

	if netPeerName == "" && (!strings.ContainsRune(peerAddress, ':') || net.ParseIP(peerAddress) != nil) {
		// peer.address is not necessarily a hostname
		// or IP address; it could be something like
		// a JDBC connection string or ip:port. Ignore
		// values containing colons, except for IPv6.
		netPeerName = peerAddress
	}

	destPort := netPeerPort
	destAddr := netPeerName
	if destAddr == "" {
		destAddr = netPeerIP
	}

	var fullURL *url.URL
	if httpURL != "" {
		fullURL, _ = url.Parse(httpURL)
	} else if httpTarget != "" || httpURLPath != "" {
		// Build http.url from http.scheme, http.target, etc.
		// http.target was supported in older semconv versions, but was replaced by url.path and url.query
		if httpTarget == "" {
			httpTarget = httpURLPath
			if httpURLQuery != "" {
				httpTarget += "?" + httpURLQuery
			}
		}

		if u, err := url.Parse(httpTarget); err == nil {
			fullURL = u
			fullURL.Scheme = httpScheme
			if httpHost == "" {
				// Set host from net.peer.*
				httpHost = destAddr
				if destPort > 0 {
					httpHost = net.JoinHostPort(httpHost, strconv.Itoa(destPort))
				}
			}
			fullURL.Host = httpHost
			httpURL = fullURL.String()
		}
	}
	if fullURL != nil {
		var port int
		portString := fullURL.Port()
		if portString != "" {
			port, _ = strconv.Atoi(portString)
		} else {
			port = schemeDefaultPort(fullURL.Scheme)
		}

		// Set destination.{address,port} from the HTTP URL,
		// replacing peer.* based values to ensure consistency.
		destAddr = truncate(fullURL.Hostname())
		if port > 0 {
			destPort = port
		}
	}

	serviceTarget.Name = peerService
	destinationService.Name = peerService
	destinationService.Resource = peerService
	if peerAddress != "" {
		destinationService.Resource = peerAddress
	}

	if isHTTP {
		if httpResponse.StatusCode > 0 && event.Event.Outcome == outcomeUnknown {
			event.Event.Outcome = clientHTTPStatusCodeOutcome(int(httpResponse.StatusCode))
		}
		if http.SizeVT() != 0 {
			event.Http = &http
		}
		if event.Url == nil {
			event.Url = &modelpb.URL{}
		}
		event.Url.Original = httpURL
	}
	if isDatabase {
		event.Span.Db = &db
	}
	if isMessaging {
		event.Span.Message = &message
	}

	switch {
	case isDatabase:
		event.Span.Type = "db"
		event.Span.Subtype = db.Type
		serviceTarget.Type = event.Span.Type
		if event.Span.Subtype != "" {
			serviceTarget.Type = event.Span.Subtype
			if destinationService.Name == "" {
				// For database requests, we currently just identify
				// the destination service by db.system.
				destinationService.Name = event.Span.Subtype
				destinationService.Resource = event.Span.Subtype
			}
		}
		if db.Instance != "" {
			serviceTarget.Name = db.Instance
		}
	case isMessaging:
		event.Span.Type = "messaging"
		event.Span.Subtype = messageSystem
		if messageOperation == "" && spanKind == ptrace.SpanKindProducer {
			messageOperation = "send"
		}
		event.Span.Action = messageOperation
		serviceTarget.Type = event.Span.Type
		if event.Span.Subtype != "" {
			serviceTarget.Type = event.Span.Subtype
			if destinationService.Name == "" {
				destinationService.Name = event.Span.Subtype
				destinationService.Resource = event.Span.Subtype
			}
		}
		if destinationService.Resource != "" && message.QueueName != "" {
			destinationService.Resource += "/" + message.QueueName
		}
		if message.QueueName != "" && !messageTempDestination {
			serviceTarget.Name = message.QueueName
		}
	case isRPC:
		event.Span.Type = "external"
		event.Span.Subtype = rpcSystem
		serviceTarget.Type = event.Span.Type
		if event.Span.Subtype != "" {
			serviceTarget.Type = event.Span.Subtype
		}
		// Set destination.service.* from the peer address, unless peer.service was specified.
		if destinationService.Name == "" {
			destHostPort := net.JoinHostPort(destAddr, strconv.Itoa(destPort))
			destinationService.Name = destHostPort
			destinationService.Resource = destHostPort
		}
		if rpcService != "" {
			serviceTarget.Name = rpcService
		}
	case isHTTP:
		event.Span.Type = "external"
		subtype := "http"
		event.Span.Subtype = subtype
		serviceTarget.Type = event.Span.Subtype
		if fullURL != nil {
			url := url.URL{Scheme: fullURL.Scheme, Host: fullURL.Host}
			resource := url.Host
			if destPort == schemeDefaultPort(url.Scheme) {
				if fullURL.Port() != "" {
					// Remove the default port from destination.service.name
					url.Host = destAddr
				} else {
					// Add the default port to destination.service.resource
					resource = fmt.Sprintf("%s:%d", resource, destPort)
				}
			}

			serviceTarget.Name = resource
			if destinationService.Name == "" {
				destinationService.Name = url.String()
				destinationService.Resource = resource
			}
		}
	case isGenAi:
		event.Span.Type = "genai"
		event.Span.Subtype = genAiSystem
		serviceTarget.Type = event.Span.Type
	default:
		// Only set event.Span.Type if not already set
		if event.Span.Type == "" {
			switch spanKind {
			case ptrace.SpanKindInternal:
				event.Span.Type = "app"
				event.Span.Subtype = "internal"
			default:
				event.Span.Type = "unknown"
			}
		}
	}

	if destAddr != "" {
		event.Destination = &modelpb.Destination{
			Address: destAddr,
			Port:    uint32(destPort),
		}
	}
	if destinationService.SizeVT() != 0 {
		if destinationService.Type == "" {
			// Copy span type to destination.service.type.
			destinationService.Type = event.Span.Type
		}
		event.Span.DestinationService = &destinationService
	}

	if serviceTarget.SizeVT() != 0 {
		event.Service.Target = &serviceTarget
	}

	if samplerType != (pcommon.Value{}) {
		// The client has reported its sampling rate, so we can use it to extrapolate transaction metrics.
		parseSamplerAttributes(samplerType, samplerParam, event)
	}

	// if outcome is still not assigned, assign success
	if event.Event.Outcome == outcomeUnknown {
		event.Event.Outcome = outcomeSuccess
	}
}

func parseSamplerAttributes(samplerType, samplerParam pcommon.Value, event *modelpb.APMEvent) {
	switch samplerType := samplerType.Str(); samplerType {
	case "probabilistic":
		probability := samplerParam.Double()
		if probability > 0 && probability <= 1 {
			if event.Span != nil {
				event.Span.RepresentativeCount = 1 / probability
			}
			if event.Transaction != nil {
				event.Transaction.RepresentativeCount = 1 / probability
			}
		}
	default:
		if event.Span != nil {
			event.Span.RepresentativeCount = 0
		}
		if event.Transaction != nil {
			event.Transaction.RepresentativeCount = 0
		}
		modelpb.Labels(event.Labels).Set("sampler_type", samplerType)
		switch samplerParam.Type() {
		case pcommon.ValueTypeBool:
			modelpb.Labels(event.Labels).Set("sampler_param", strconv.FormatBool(samplerParam.Bool()))
		case pcommon.ValueTypeDouble:
			modelpb.NumericLabels(event.NumericLabels).Set("sampler_param", samplerParam.Double())
		}
	}
}

func (c *Consumer) convertSpanEvent(
	spanEvent ptrace.SpanEvent,
	parent *modelpb.APMEvent, // either span or transaction
	timeDelta time.Duration,
) *modelpb.APMEvent {
	event := parent.CloneVT()
	initEventLabels(event)
	event.Transaction = nil // populate fields as required from parent
	event.Span = nil        // populate fields as required from parent
	event.ParentId = ""     // populate fields as required from parent

	// Remove unnecessary fields from span event
	if event.Service != nil {
		event.Service.Target = nil
		event.Service.Origin = nil
	}

	event.Timestamp = modelpb.FromTime(spanEvent.Timestamp().AsTime().Add(timeDelta))

	isJaeger := strings.HasPrefix(parent.Agent.Name, "Jaeger")
	if isJaeger {
		event.Error = c.convertJaegerErrorSpanEvent(spanEvent, event)
	} else if spanEvent.Name() == "exception" {
		// Translate exception span events to errors.
		//
		// If it's not Jaeger, we assume OpenTelemetry semantic semconv.
		// Per OpenTelemetry semantic conventions:
		//   `The name of the event MUST be "exception"`
		var exceptionEscaped bool
		var exceptionMessage, exceptionStacktrace, exceptionType string
		spanEvent.Attributes().Range(func(k string, v pcommon.Value) bool {
			switch k {
			case string(semconv.ExceptionMessageKey):
				exceptionMessage = v.Str()
			case string(semconv.ExceptionStacktraceKey):
				exceptionStacktrace = v.Str()
			case string(semconv.ExceptionTypeKey):
				exceptionType = v.Str()
			case "exception.escaped":
				exceptionEscaped = v.Bool()

			// data_stream.*
			// Note: fields are parsed but dataset will be overridden by SetDataStream because it is an error
			case attributeDataStreamDataset:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Dataset = sanitizeDataStreamDataset(v.Str())
			case attributeDataStreamNamespace:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Namespace = sanitizeDataStreamNamespace(v.Str())

			default:
				setLabel(replaceDots(k), event, v)
			}
			return true
		})
		if exceptionMessage != "" || exceptionType != "" {
			// Per OpenTelemetry semantic conventions:
			//   `At least one of the following sets of attributes is required:
			//   - exception.type
			//   - exception.message`
			event.Error = convertOpenTelemetryExceptionSpanEvent(
				exceptionType, exceptionMessage, exceptionStacktrace,
				exceptionEscaped, parent.Service.Language.Name,
			)
		}
	}

	if event.Error != nil {
		setErrorContext(event, parent)
	} else {
		// Set "event.kind" to indicate this is a log event.
		if event.Event == nil {
			event.Event = &modelpb.Event{}
		}
		event.Event.Kind = "event"
		event.Message = spanEvent.Name()
		setLogContext(event, parent)
		spanEvent.Attributes().Range(func(k string, v pcommon.Value) bool {
			switch k {
			// data_stream.*
			case attributeDataStreamDataset:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Dataset = sanitizeDataStreamDataset(v.Str())
			case attributeDataStreamNamespace:
				if event.DataStream == nil {
					event.DataStream = &modelpb.DataStream{}
				}
				event.DataStream.Namespace = sanitizeDataStreamNamespace(v.Str())
			default:
				k = replaceDots(k)
				if isJaeger && k == "message" {
					event.Message = truncate(v.Str())
					return true
				}
				setLabel(k, event, v)
			}
			return true
		})
	}
	return event
}

func (c *Consumer) convertJaegerErrorSpanEvent(event ptrace.SpanEvent, apmEvent *modelpb.APMEvent) *modelpb.Error {
	var isError bool
	var exMessage, exType string
	var logMessage string

	if name := truncate(event.Name()); name == "error" {
		isError = true // according to opentracing spec
	} else {
		// Jaeger seems to send the message in the 'event' field.
		//
		// In case 'message' is sent we will use that, otherwise
		// we will use 'event'.
		logMessage = name
	}

	event.Attributes().Range(func(k string, v pcommon.Value) bool {
		if v.Type() != pcommon.ValueTypeStr {
			return true
		}
		stringval := truncate(v.Str())
		switch k {
		case "error", "error.object":
			exMessage = stringval
			isError = true
		case "error.kind":
			exType = stringval
			isError = true
		case "level":
			isError = stringval == "error"
		case "message":
			logMessage = stringval

		// data_stream.*
		// Note: fields are parsed but dataset will be overridden by SetDataStream because it is an error
		case attributeDataStreamDataset:
			if apmEvent.DataStream == nil {
				apmEvent.DataStream = &modelpb.DataStream{}
			}
			apmEvent.DataStream.Dataset = sanitizeDataStreamDataset(v.Str())
		case attributeDataStreamNamespace:
			if apmEvent.DataStream == nil {
				apmEvent.DataStream = &modelpb.DataStream{}
			}
			apmEvent.DataStream.Namespace = sanitizeDataStreamNamespace(v.Str())

		default:
			setLabel(replaceDots(k), apmEvent, v)
		}
		return true
	})
	if !isError {
		return nil
	}
	if logMessage == "" && exMessage == "" && exType == "" {
		return nil
	}
	e := modelpb.Error{}
	if logMessage != "" {
		e.Log = &modelpb.ErrorLog{}
		e.Log.Message = logMessage
	}
	if exMessage != "" || exType != "" {
		e.Exception = &modelpb.Exception{}
		e.Exception.Message = exMessage
		e.Exception.Type = exType
		if id, err := newUniqueID(); err == nil {
			e.Id = id
		}
	}
	return &e
}

func setErrorContext(out *modelpb.APMEvent, parent *modelpb.APMEvent) {
	out.Trace.Id = parent.Trace.Id
	out.Http = parent.Http
	out.Url = parent.Url
	if parent.Transaction != nil {
		out.Transaction = &modelpb.Transaction{}
		out.Transaction.Id = parent.Transaction.Id
		out.Transaction.Sampled = parent.Transaction.Sampled
		out.Transaction.Type = parent.Transaction.Type
		out.Span = &modelpb.Span{}
		out.Span.Id = parent.Transaction.Id
		out.Error.Custom = parent.Transaction.Custom
		out.ParentId = parent.Transaction.Id
	}
	if parent.Span != nil {
		out.ParentId = parent.Span.Id
	}
}

func setLogContext(out *modelpb.APMEvent, parent *modelpb.APMEvent) {
	if parent.Transaction != nil {
		out.Transaction = &modelpb.Transaction{}
		out.Transaction.Id = parent.Transaction.Id
		out.Span = &modelpb.Span{}
		out.Span.Id = parent.Transaction.Id
	}
	if parent.Span != nil {
		out.Span = &modelpb.Span{}
		out.Span.Id = parent.Span.Id
	}
}

func translateSpanLinks(out *modelpb.APMEvent, in ptrace.SpanLinkSlice) {
	n := in.Len()
	if n == 0 {
		return
	}
	if out.Span == nil {
		out.Span = &modelpb.Span{}
	}
	out.Span.Links = make([]*modelpb.SpanLink, 0, n)
	for i := 0; i < n; i++ {
		link := in.At(i)
		// When a link has the elastic.is_child attribute set, it is stored in the child_ids instead
		elChildAttribVal, elChildAttribPresent := link.Attributes().Get("elastic.is_child")
		// alternatively, we also look for just "is_child" without the elastic. prefix
		childAttribVal, childAttribPresent := link.Attributes().Get("is_child")
		if (elChildAttribPresent && elChildAttribVal.Bool()) || (childAttribPresent && childAttribVal.Bool()) {
			out.ChildIds = append(out.ChildIds, hexSpanID(link.SpanID()))
		} else {
			sl := modelpb.SpanLink{}
			sl.SpanId = hexSpanID(link.SpanID())
			sl.TraceId = hexTraceID(link.TraceID())
			out.Span.Links = append(out.Span.Links, &sl)
		}
	}
}

func hexSpanID(id pcommon.SpanID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

func hexTraceID(id pcommon.TraceID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

func replaceDots(s string) string {
	return strings.ReplaceAll(s, dot, underscore)
}

// spanStatusOutcome returns the outcome for transactions and spans based on
// the given OTLP span status.
func spanStatusOutcome(status ptrace.Status) string {
	switch status.Code() {
	case ptrace.StatusCodeOk:
		return outcomeSuccess
	case ptrace.StatusCodeError:
		return outcomeFailure
	}
	return outcomeUnknown
}

// spanStatusResult returns the result for transactions based on the given
// OTLP span status. If the span status is unknown, an empty result string
// is returned.
func spanStatusResult(status ptrace.Status) string {
	switch status.Code() {
	case ptrace.StatusCodeOk:
		return "Success"
	case ptrace.StatusCodeError:
		return "Error"
	}
	return ""
}

var standardStatusCodeResults = [...]string{
	"HTTP 1xx",
	"HTTP 2xx",
	"HTTP 3xx",
	"HTTP 4xx",
	"HTTP 5xx",
}

// httpStatusCodeResult returns the transaction result value to use for the
// given HTTP status code.
func httpStatusCodeResult(statusCode int) string {
	switch i := statusCode / 100; i {
	case 1, 2, 3, 4, 5:
		return standardStatusCodeResults[i-1]
	}
	return fmt.Sprintf("HTTP %d", statusCode)
}

// serverHTTPStatusCodeOutcome returns the transaction outcome value to use for
// the given HTTP status code.
func serverHTTPStatusCodeOutcome(statusCode int) string {
	if statusCode >= 500 {
		return outcomeFailure
	}
	return outcomeSuccess
}

// clientHTTPStatusCodeOutcome returns the span outcome value to use for the
// given HTTP status code.
func clientHTTPStatusCodeOutcome(statusCode int) string {
	if statusCode >= 400 {
		return outcomeFailure
	}
	return outcomeSuccess
}

// truncate returns s truncated at n runes, and the number of runes in the resulting string (<= n).
func truncate(s string) string {
	var j int
	for i := range s {
		if j == keywordLength {
			return s[:i]
		}
		j++
	}
	return s
}

func schemeDefaultPort(scheme string) int {
	switch scheme {
	case "http":
		return 80
	case "https":
		return 443
	}
	return 0
}

// parses traceparent header, which is expected to be in the W3C Trace-Context
// and searches for the p-value as specified in
//
//	https://opentelemetry.io/docs/reference/specification/trace/tracestate-probability-sampling/#p-value
//
// to calculate the adjusted count (i.e. representative count)
//
// If the p-value is missing or invalid in the tracestate we assume
// a sampling rate of 100% and a representative count of 1.
func getRepresentativeCountFromTracestateHeader(tracestace string) float64 {
	// Default p-value is 0, leading to a default representative count of 1.
	var p uint64 = 0

	otValue := getValueForKeyInString(tracestace, "ot", ',', '=')

	if otValue != "" {
		pValue := getValueForKeyInString(otValue, "p", ';', ':')

		if pValue != "" {
			p, _ = strconv.ParseUint(pValue, 10, 6)
		}
	}

	if p > 62 {
		return 0.0
	}

	return math.Pow(2, float64(p))
}

func getValueForKeyInString(str string, key string, separator rune, assignChar rune) string {
	for {
		str = strings.TrimSpace(str)
		if str == "" {
			break
		}
		kv := str
		if sepIdx := strings.IndexRune(str, separator); sepIdx != -1 {
			kv = strings.TrimSpace(str[:sepIdx])
			str = str[sepIdx+1:]
		} else {
			str = ""
		}
		equal := strings.IndexRune(kv, assignChar)
		if equal != -1 && kv[:equal] == key {
			return kv[equal+1:]
		}
	}

	return ""
}
