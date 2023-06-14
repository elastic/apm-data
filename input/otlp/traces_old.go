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

package otlp

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"

	"github.com/elastic/apm-data/model"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
	"google.golang.org/grpc/codes"
)

// TranslateTransactionOld is used by the v2 decoder and not by the otlp input.
// It exists for compatibility reason and to migrate the otlp input to protobuf
// without touching the decoder.
func TranslateTransactionOld(
	attributes pcommon.Map,
	spanStatus ptrace.Status,
	library pcommon.InstrumentationScope,
	event *model.APMEvent,
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
		http           model.HTTP
		httpRequest    model.HTTPRequest
		httpResponse   model.HTTPResponse
	)

	var foundSpanType int
	var message model.Message

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
			setLabelOld(k, event, ifaceAttributeValue(v))
		case pcommon.ValueTypeBool:
			setLabelOld(k, event, ifaceAttributeValue(v))
		case pcommon.ValueTypeDouble:
			setLabelOld(k, event, ifaceAttributeValue(v))
		case pcommon.ValueTypeInt:
			switch kDots {
			case semconv.AttributeHTTPStatusCode:
				foundSpanType = httpSpan
				httpResponse.StatusCode = int(v.Int())
				http.Response = &httpResponse
			case semconv.AttributeNetPeerPort:
				event.Source.Port = int(v.Int())
			case semconv.AttributeNetHostPort:
				netHostPort = int(v.Int())
			case semconv.AttributeRPCGRPCStatusCode:
				foundSpanType = rpcSpan
				event.Transaction.Result = codes.Code(v.Int()).String()
			default:
				setLabelOld(k, event, ifaceAttributeValue(v))
			}
		case pcommon.ValueTypeMap:
		case pcommon.ValueTypeStr:
			stringval := truncate(v.Str())
			switch kDots {
			// http.*
			case semconv.AttributeHTTPMethod:
				foundSpanType = httpSpan
				httpRequest.Method = stringval
				http.Request = &httpRequest
			case semconv.AttributeHTTPURL, semconv.AttributeHTTPTarget, "http.path":
				foundSpanType = httpSpan
				httpURL = stringval
			case semconv.AttributeHTTPHost:
				foundSpanType = httpSpan
				httpHost = stringval
			case semconv.AttributeHTTPScheme:
				foundSpanType = httpSpan
				httpScheme = stringval
			case semconv.AttributeHTTPStatusCode:
				if intv, err := strconv.Atoi(stringval); err == nil {
					foundSpanType = httpSpan
					httpResponse.StatusCode = intv
					http.Response = &httpResponse
				}
			case "http.protocol":
				if !strings.HasPrefix(stringval, "HTTP/") {
					// Unexpected, store in labels for debugging.
					event.Labels.Set(k, stringval)
					break
				}
				stringval = strings.TrimPrefix(stringval, "HTTP/")
				fallthrough
			case semconv.AttributeHTTPFlavor:
				foundSpanType = httpSpan
				http.Version = stringval
			case semconv.AttributeHTTPServerName:
				foundSpanType = httpSpan
				httpServerName = stringval
			case semconv.AttributeHTTPClientIP:
				if ip, err := netip.ParseAddr(stringval); err == nil {
					event.Client.IP = ip
				}
			case semconv.AttributeHTTPUserAgent:
				event.UserAgent.Original = stringval

			// net.*
			case semconv.AttributeNetPeerIP:
				if ip, err := netip.ParseAddr(stringval); err == nil {
					event.Source.IP = ip
				}
			case semconv.AttributeNetPeerName:
				event.Source.Domain = stringval
			case semconv.AttributeNetHostName:
				netHostName = stringval
			case attributeNetworkConnectionType:
				event.Network.Connection.Type = stringval
			case attributeNetworkConnectionSubtype:
				event.Network.Connection.Subtype = stringval
			case attributeNetworkMCC:
				event.Network.Carrier.MCC = stringval
			case attributeNetworkMNC:
				event.Network.Carrier.MNC = stringval
			case attributeNetworkCarrierName:
				event.Network.Carrier.Name = stringval
			case attributeNetworkICC:
				event.Network.Carrier.ICC = stringval

			// messaging.*
			case "message_bus.destination", semconv.AttributeMessagingDestination:
				message.QueueName = stringval
				foundSpanType = messagingSpan

			// rpc.*
			//
			// TODO(axw) add RPC fieldset to ECS? Currently we drop these
			// attributes, and rely on the operation name like we do with
			// Elastic APM agents.
			case semconv.AttributeRPCSystem:
				foundSpanType = rpcSpan
			case semconv.AttributeRPCGRPCStatusCode:
				foundSpanType = rpcSpan
			case semconv.AttributeRPCService:
			case semconv.AttributeRPCMethod:

			// miscellaneous
			case "type":
				event.Transaction.Type = stringval
			case "session.id":
				event.Session.ID = stringval
			case semconv.AttributeServiceVersion:
				// NOTE support for sending service.version as a span tag
				// is deprecated, and will be removed in 8.0. Instrumentation
				// should set this as a resource attribute (OTel) or tracer
				// tag (Jaeger).
				event.Service.Version = stringval
			default:
				event.Labels.Set(k, stringval)
			}
		}
		return true
	})

	if event.Transaction.Type == "" {
		switch foundSpanType {
		case httpSpan, rpcSpan:
			event.Transaction.Type = "request"
		case messagingSpan:
			event.Transaction.Type = "messaging"
		default:
			event.Transaction.Type = "unknown"
		}
	}

	switch foundSpanType {
	case httpSpan:
		event.HTTP = http

		// Set outcome nad result from status code.
		if statusCode := httpResponse.StatusCode; statusCode > 0 {
			if event.Event.Outcome == outcomeUnknown {
				event.Event.Outcome = serverHTTPStatusCodeOutcome(statusCode)
			}
			if event.Transaction.Result == "" {
				event.Transaction.Result = httpStatusCodeResult(statusCode)
			}
		}

		// Build the model.URL from http{URL,Host,Scheme}.
		httpHost := httpHost
		if httpHost == "" {
			httpHost = httpServerName
			if httpHost == "" {
				httpHost = netHostName
				if httpHost == "" {
					httpHost = event.Host.Hostname
				}
			}
			if httpHost != "" && netHostPort > 0 {
				httpHost = net.JoinHostPort(httpHost, strconv.Itoa(netHostPort))
			}
		}
		event.URL = model.ParseURL(httpURL, httpHost, httpScheme)
	case messagingSpan:
		event.Transaction.Message = &message
	}

	if !event.Client.IP.IsValid() {
		event.Client = model.Client{IP: event.Source.IP, Port: event.Source.Port, Domain: event.Source.Domain}
	}

	if samplerType != (pcommon.Value{}) {
		// The client has reported its sampling rate, so we can use it to extrapolate span metrics.
		parseSamplerAttributesOld(samplerType, samplerParam, event)
	}

	if event.Transaction.Result == "" {
		event.Transaction.Result = spanStatusResult(spanStatus)
	}
	if name := library.Name(); name != "" {
		event.Service.Framework.Name = name
		event.Service.Framework.Version = library.Version()
	}
}

func parseSamplerAttributesOld(samplerType, samplerParam pcommon.Value, event *model.APMEvent) {
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
		event.Transaction.RepresentativeCount = 0
		event.Labels.Set("sampler_type", samplerType)
		switch samplerParam.Type() {
		case pcommon.ValueTypeBool:
			event.Labels.Set("sampler_param", strconv.FormatBool(samplerParam.Bool()))
		case pcommon.ValueTypeDouble:
			event.NumericLabels.Set("sampler_param", samplerParam.Double())
		}
	}
}

func setLabelOld(key string, event *model.APMEvent, v interface{}) {
	switch v := v.(type) {
	case string:
		event.Labels.Set(key, v)
	case bool:
		event.Labels.Set(key, strconv.FormatBool(v))
	case float64:
		event.NumericLabels.Set(key, v)
	case int64:
		event.NumericLabels.Set(key, float64(v))
	case []interface{}:
		if len(v) == 0 {
			return
		}
		switch v[0].(type) {
		case string:
			value := make([]string, len(v))
			for i := range v {
				value[i] = v[i].(string)
			}
			event.Labels.SetSlice(key, value)
		case float64:
			value := make([]float64, len(v))
			for i := range v {
				value[i] = v[i].(float64)
			}
			event.NumericLabels.SetSlice(key, value)
		}
	}
}

func TranslateSpanOld(spanKind ptrace.SpanKind, attributes pcommon.Map, event *model.APMEvent) {
	isJaeger := strings.HasPrefix(event.Agent.Name, "Jaeger")

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
		httpURL    string
		httpHost   string
		httpTarget string
		httpScheme = "http"
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

	var http model.HTTP
	var httpRequest model.HTTPRequest
	var httpResponse model.HTTPResponse
	var message model.Message
	var db model.DB
	var destinationService model.DestinationService
	var serviceTarget model.ServiceTarget
	var foundSpanType int
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
			setLabelOld(k, event, ifaceAttributeValueSlice(v.Slice()))
		case pcommon.ValueTypeBool:
			switch kDots {
			case semconv.AttributeMessagingTempDestination:
				messageTempDestination = v.Bool()
				fallthrough
			default:
				setLabelOld(k, event, strconv.FormatBool(v.Bool()))
			}
		case pcommon.ValueTypeDouble:
			setLabelOld(k, event, v.Double())
		case pcommon.ValueTypeInt:
			switch kDots {
			case "http.status_code":
				httpResponse.StatusCode = int(v.Int())
				http.Response = &httpResponse
				foundSpanType = httpSpan
			case semconv.AttributeNetPeerPort, "peer.port":
				netPeerPort = int(v.Int())
			case semconv.AttributeRPCGRPCStatusCode:
				rpcSystem = "grpc"
				foundSpanType = rpcSpan
			default:
				setLabelOld(k, event, v.Int())
			}
		case pcommon.ValueTypeStr:
			stringval := truncate(v.Str())

			switch kDots {
			// http.*
			case semconv.AttributeHTTPHost:
				httpHost = stringval
				foundSpanType = httpSpan
			case semconv.AttributeHTTPScheme:
				httpScheme = stringval
				foundSpanType = httpSpan
			case semconv.AttributeHTTPTarget:
				httpTarget = stringval
				foundSpanType = httpSpan
			case semconv.AttributeHTTPURL:
				httpURL = stringval
				foundSpanType = httpSpan
			case semconv.AttributeHTTPMethod:
				httpRequest.Method = stringval
				http.Request = &httpRequest
				foundSpanType = httpSpan

			// db.*
			case "sql.query":
				if db.Type == "" {
					db.Type = "sql"
				}
				fallthrough
			case semconv.AttributeDBStatement:
				// Statement should not be truncated, use original string value.
				db.Statement = v.Str()
				foundSpanType = dbSpan
			case semconv.AttributeDBName, "db.instance":
				db.Instance = stringval
				foundSpanType = dbSpan
			case semconv.AttributeDBSystem, "db.type":
				db.Type = stringval
				foundSpanType = dbSpan
			case semconv.AttributeDBUser:
				db.UserName = stringval
				foundSpanType = dbSpan

			// net.*
			case semconv.AttributeNetPeerName, "peer.hostname":
				netPeerName = stringval
			case semconv.AttributeNetPeerIP, "peer.ipv4", "peer.ipv6":
				netPeerIP = stringval
			case "peer.address":
				peerAddress = stringval
			case attributeNetworkConnectionType:
				event.Network.Connection.Type = stringval
			case attributeNetworkConnectionSubtype:
				event.Network.Connection.Subtype = stringval
			case attributeNetworkMCC:
				event.Network.Carrier.MCC = stringval
			case attributeNetworkMNC:
				event.Network.Carrier.MNC = stringval
			case attributeNetworkCarrierName:
				event.Network.Carrier.Name = stringval
			case attributeNetworkICC:
				event.Network.Carrier.ICC = stringval

			// session.*
			case "session.id":
				event.Session.ID = stringval

			// messaging.*
			case "message_bus.destination", semconv.AttributeMessagingDestination:
				message.QueueName = stringval
				foundSpanType = messagingSpan
			case semconv.AttributeMessagingOperation:
				messageOperation = stringval
				foundSpanType = messagingSpan
			case semconv.AttributeMessagingSystem:
				messageSystem = stringval
				foundSpanType = messagingSpan

			// rpc.*
			//
			// TODO(axw) add RPC fieldset to ECS? Currently we drop these
			// attributes, and rely on the operation name and span type/subtype
			// like we do with Elastic APM agents.
			case semconv.AttributeRPCSystem:
				rpcSystem = stringval
				foundSpanType = rpcSpan
			case semconv.AttributeRPCService:
				rpcService = stringval
				foundSpanType = rpcSpan
			case semconv.AttributeRPCGRPCStatusCode:
				rpcSystem = "grpc"
				foundSpanType = rpcSpan
			case semconv.AttributeRPCMethod:

			// miscellaneous
			case "span.kind": // filter out
			case semconv.AttributePeerService:
				peerService = stringval
			default:
				event.Labels.Set(k, stringval)
			}
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
	} else if httpTarget != "" {
		// Build http.url from http.scheme, http.target, etc.
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

	switch foundSpanType {
	case httpSpan:
		if httpResponse.StatusCode > 0 {
			if event.Event.Outcome == outcomeUnknown {
				event.Event.Outcome = clientHTTPStatusCodeOutcome(httpResponse.StatusCode)
			}
		}
		event.Span.Type = "external"
		subtype := "http"
		event.Span.Subtype = subtype
		event.HTTP = http
		event.URL.Original = httpURL
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
	case dbSpan:
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
		event.Span.DB = &db
	case messagingSpan:
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
		event.Span.Message = &message
	case rpcSpan:
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
		event.Destination = model.Destination{Address: destAddr, Port: destPort}
	}
	if destinationService != (model.DestinationService{}) {
		if destinationService.Type == "" {
			// Copy span type to destination.service.type.
			destinationService.Type = event.Span.Type
		}
		event.Span.DestinationService = &destinationService
	}

	if serviceTarget != (model.ServiceTarget{}) {
		event.Service.Target = &serviceTarget
	}

	if samplerType != (pcommon.Value{}) {
		// The client has reported its sampling rate, so we can use it to extrapolate transaction metrics.
		parseSamplerAttributesOld(samplerType, samplerParam, event)
	}
}

