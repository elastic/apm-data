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

package v2

import (
	"encoding/json"

	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/nullable"
)

var (
	patternAlphaNumericExt = `^[a-zA-Z0-9 _-]+$`
	patternNoAsteriskQuote = `^[^*"]*$` //do not allow '*' '"'

	enumOutcome = []string{"success", "failure", "unknown"}
)

// entry points

// errorRoot requires an error event to be present
type errorRoot struct {
	Error errorEvent `json:"error" validate:"required"`
}

// metadatatRoot requires a metadata event to be present
type metadataRoot struct {
	Metadata metadata `json:"metadata" validate:"required"`
}

// metricsetRoot requires a metricset event to be present
type metricsetRoot struct {
	Metricset metricset `json:"metricset" validate:"required"`
}

// spanRoot requires a span event to be present
type spanRoot struct {
	Span span `json:"span" validate:"required"`
}

// transactionRoot requires a transaction event to be present
type transactionRoot struct {
	Transaction transaction `json:"transaction" validate:"required"`
}

// logRoot requires a log event to be present
type logRoot struct {
	Log log `json:"log" validate:"required"`
}

// other structs

type context struct {
	// Custom can contain additional metadata to be stored with the event.
	// The format is unspecified and can be deeply nested objects.
	// The information will not be indexed or searchable in Elasticsearch.
	Custom map[string]any `json:"custom"`
	// Tags are a flat mapping of user-defined tags. On the agent side, tags
	// are called labels. Allowed value types are string, boolean and number
	// values. Tags are indexed and searchable.
	Tags map[string]any `json:"tags" validate:"inputTypesVals=string;bool;number,maxLengthVals=1024"`
	// Service related information can be sent per event. Information provided
	// here will override the more generic information retrieved from metadata,
	// missing service fields will be retrieved from the metadata information.
	Service contextService `json:"service"`
	// Cloud holds fields related to the cloud or infrastructure the events
	// are coming from.
	Cloud contextCloud `json:"cloud"`
	// User holds information about the correlated user for this event. If
	// user data are provided here, all user related information from metadata
	// is ignored, otherwise the metadata's user information will be stored
	// with the event.
	User user `json:"user"`
	// Page holds information related to the current page and page referers.
	// It is only sent from RUM agents.
	Page contextPage `json:"page"`
	// Request describes the HTTP request information in case the event was
	// created as a result of an HTTP request.
	Request contextRequest `json:"request"`
	// Message holds details related to message receiving and publishing
	// if the captured event integrates with a messaging system
	Message contextMessage `json:"message"`
	// Response describes the HTTP response information in case the event was
	// created as a result of an HTTP request.
	Response contextResponse `json:"response"`
}

type faas struct {
	// A unique identifier of the invoked serverless function.
	ID nullable.String `json:"id"`
	// The request id of the function invocation.
	Execution nullable.String `json:"execution"`
	// Trigger attributes.
	Trigger trigger `json:"trigger"`
	// The lambda function name.
	Name nullable.String `json:"name"`
	// The lambda function version.
	Version nullable.String `json:"version"`
	// Indicates whether a function invocation was a cold start or not.
	Coldstart nullable.Bool `json:"coldstart"`
}

type trigger struct {
	// The trigger type.
	Type nullable.String `json:"type"`
	// The id of the origin trigger request.
	RequestID nullable.String `json:"request_id"`
}

type contextCloud struct {
	// Origin contains the self-nested field groups for cloud.
	Origin contextCloudOrigin `json:"origin"`
}

type contextCloudOrigin struct {
	// The cloud account or organization id used to identify
	// different entities in a multi-tenant environment.
	Account contextCloudOriginAccount `json:"account"`
	// Name of the cloud provider.
	Provider nullable.String `json:"provider"`
	// Region in which this host, resource, or service is located.
	Region nullable.String `json:"region"`
	// The cloud service name is intended to distinguish services running
	// on different platforms within a provider.
	Service contextCloudOriginService `json:"service"`
}

type contextCloudOriginAccount struct {
	// The cloud account or organization id used to identify
	// different entities in a multi-tenant environment.
	ID nullable.String `json:"id"`
}

type contextCloudOriginService struct {
	// The cloud service name is intended to distinguish services running
	// on different platforms within a provider.
	Name nullable.String `json:"name"`
}

type contextMessage struct {
	// Headers received with the message, similar to HTTP request headers.
	Headers nullable.HTTPHeader `json:"headers"`
	// Body of the received message, similar to an HTTP request body
	Body nullable.String `json:"body"`
	// Queue holds information about the message queue where the message is received.
	Queue contextMessageQueue `json:"queue"`
	// RoutingKey holds the optional routing key of the received message as set
	// on the queuing system, such as in RabbitMQ.
	RoutingKey nullable.String `json:"routing_key"`
	// Age of the message. If the monitored messaging framework provides a
	// timestamp for the message, agents may use it. Otherwise, the sending
	// agent can add a timestamp in milliseconds since the Unix epoch to the
	// message's metadata to be retrieved by the receiving agent. If a
	// timestamp is not available, agents should omit this field.
	Age contextMessageAge `json:"age"`
}

type contextMessageAge struct {
	// Age of the message in milliseconds.
	Milliseconds nullable.Int `json:"ms"`
}

type contextMessageQueue struct {
	// Name holds the name of the message queue where the message is received.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
}

type contextPage struct {
	// Referer holds the URL of the page that 'linked' to the current page.
	Referer nullable.String `json:"referer"`
	// URL of the current page
	URL nullable.String `json:"url"`
}

type contextRequest struct {
	// Cookies used by the request, parsed as key-value objects.
	Cookies map[string]any `json:"cookies"`
	// Env holds environment variable information passed to the monitored service.
	Env map[string]any `json:"env"`
	// Body only contais the request bod, not the query string information.
	// It can either be a dictionary (for standard HTTP requests) or a raw
	// request body.
	Body nullable.Interface `json:"body" validate:"inputTypes=string;object"`
	// Headers includes any HTTP headers sent by the requester. Cookies will
	// be taken by headers if supplied.
	Headers nullable.HTTPHeader `json:"headers"`
	// URL holds information sucha as the raw URL, scheme, host and path.
	URL contextRequestURL `json:"url"`
	// HTTPVersion holds information about the used HTTP version.
	HTTPVersion nullable.String `json:"http_version" validate:"maxLength=1024"`
	// Method holds information about the method of the HTTP request.
	Method nullable.String `json:"method" validate:"required,maxLength=1024"`
	// Socket holds information related to the recorded request,
	// such as whether or not data were encrypted and the remote address.
	Socket contextRequestSocket `json:"socket"`
}

type contextRequestURL struct {
	// Port of the request, e.g. '443'. Can be sent as string or int.
	Port nullable.Interface `json:"port" validate:"inputTypes=string;int,targetType=int,maxLength=1024"`
	// Full, possibly agent-assembled URL of the request,
	// e.g. https://example.com:443/search?q=elasticsearch#top.
	Full nullable.String `json:"full" validate:"maxLength=1024"`
	// Hash of the request URL, e.g. 'top'
	Hash nullable.String `json:"hash" validate:"maxLength=1024"`
	// Hostname information of the request, e.g. 'example.com'."
	Hostname nullable.String `json:"hostname" validate:"maxLength=1024"`
	// Path of the request, e.g. '/search'
	Path nullable.String `json:"pathname" validate:"maxLength=1024"`
	// Protocol information for the recorded request, e.g. 'https:'.
	Protocol nullable.String `json:"protocol" validate:"maxLength=1024"`
	// Raw unparsed URL of the HTTP request line,
	// e.g https://example.com:443/search?q=elasticsearch. This URL may be
	// absolute or relative. For more details, see
	// https://www.w3.org/Protocols/rfc2616/rfc2616-sec5.html#sec5.1.2.
	Raw nullable.String `json:"raw" validate:"maxLength=1024"`
	// Search contains the query string information of the request. It is
	// expected to have values delimited by ampersands.
	Search nullable.String `json:"search" validate:"maxLength=1024"`
}

type contextRequestSocket struct {
	// RemoteAddress holds the network address sending the request. It should
	// be obtained through standard APIs and not be parsed from any headers
	// like 'Forwarded'.
	RemoteAddress nullable.String `json:"remote_address"`
	// Encrypted indicates whether a request was sent as TLS/HTTPS request.
	// DEPRECATED: this field will be removed in a future release.
	Encrypted nullable.Bool `json:"encrypted"`
}

type contextResponse struct {
	// Headers holds the http headers sent in the http response.
	Headers nullable.HTTPHeader `json:"headers"`
	// StatusCode sent in the http response.
	StatusCode nullable.Int `json:"status_code"`
	// TransferSize holds the total size of the payload.
	TransferSize nullable.Int `json:"transfer_size"`
	// DecodedBodySize holds the size of the decoded payload.
	DecodedBodySize nullable.Int `json:"decoded_body_size"`
	// EncodedBodySize holds the size of the encoded payload.
	EncodedBodySize nullable.Int `json:"encoded_body_size"`
	// Finished indicates whether the response was finished or not.
	Finished nullable.Bool `json:"finished"`
	// HeadersSent indicates whether http headers were sent.
	HeadersSent nullable.Bool `json:"headers_sent"`
}

type contextService struct {
	// Agent holds information about the APM agent capturing the event.
	Agent contextServiceAgent `json:"agent"`
	// Environment in which the monitored service is running,
	// e.g. `production` or `staging`.
	Environment nullable.String `json:"environment" validate:"maxLength=1024"`
	// Framework holds information about the framework used in the
	// monitored service.
	Framework contextServiceFramework `json:"framework"`
	// ID holds a unique identifier for the service.
	ID nullable.String `json:"id"`
	// Language holds information about the programming language of the
	// monitored service.
	Language contextServiceLanguage `json:"language"`
	// Name of the monitored service.
	Name nullable.String `json:"name" validate:"maxLength=1024,pattern=patternAlphaNumericExt"`
	// Node must be a unique meaningful name of the service node.
	Node contextServiceNode `json:"node"`
	// Origin contains the self-nested field groups for service.
	Origin contextServiceOrigin `json:"origin"`
	// Runtime holds information about the language runtime running the
	// monitored service
	Runtime contextServiceRuntime `json:"runtime"`
	// Target holds information about the outgoing service in case of
	// an outgoing event
	Target contextServiceTarget `json:"target"`
	// Version of the monitored service.
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type contextServiceTarget struct {
	_ struct{} `validate:"requiredAnyOf=type;name"`
	// Immutable name of the target service for the event
	Name nullable.String `json:"name"`
	// Immutable type of the target service for the event
	Type nullable.String `json:"type"`
}

type contextServiceOrigin struct {
	// Immutable id of the service emitting this event.
	ID nullable.String `json:"id"`
	// Immutable name of the service emitting this event.
	Name nullable.String `json:"name"`
	// The version of the service the data was collected from.
	Version nullable.String `json:"version"`
}

type contextServiceAgent struct {
	// EphemeralID is a free format ID used for metrics correlation by agents
	EphemeralID nullable.String `json:"ephemeral_id" validate:"maxLength=1024"`
	// Name of the APM agent capturing information.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Version of the APM agent capturing information.
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type contextServiceFramework struct {
	// Name of the used framework
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Version of the used framework
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type contextServiceLanguage struct {
	// Name of the used programming language
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Version of the used programming language
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type contextServiceNode struct {
	// Name of the service node
	Name nullable.String `json:"configured_name" validate:"maxLength=1024"`
}

type contextServiceRuntime struct {
	// Name of the language runtime
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Version of the language runtime
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

// errorEvent represents an error or a logged error message,
// captured by an APM agent in a monitored service.
type errorEvent struct {
	_ struct{} `validate:"requiredAnyOf=exception;log"`
	// Timestamp holds the recorded time of the event, UTC based and formatted
	// as microseconds since Unix epoch.
	Timestamp nullable.TimeMicrosUnix `json:"timestamp"`
	// Log holds additional information added when the error is logged.
	Log errorLog `json:"log"`
	// Culprit identifies the function call which was the primary perpetrator
	// of this event.
	Culprit nullable.String `json:"culprit" validate:"maxLength=1024"`
	// ID holds the hex encoded 128 random bits ID of the event.
	ID nullable.String `json:"id" validate:"required,maxLength=1024"`
	// ParentID holds the hex encoded 64 random bits ID of the parent
	// transaction or span.
	ParentID nullable.String `json:"parent_id" validate:"requiredIfAny=transaction_id;trace_id,maxLength=1024"`
	// TraceID holds the hex encoded 128 random bits ID of the correlated trace.
	TraceID nullable.String `json:"trace_id" validate:"requiredIfAny=transaction_id;parent_id,maxLength=1024"`
	// TransactionID holds the hex encoded 64 random bits ID of the correlated
	// transaction.
	TransactionID nullable.String `json:"transaction_id" validate:"maxLength=1024"`
	// Exception holds information about the original error.
	// The information is language specific.
	Exception errorException `json:"exception"`
	// Transaction holds information about the correlated transaction.
	Transaction errorTransactionRef `json:"transaction"`
	// Context holds arbitrary contextual information for the event.
	Context context `json:"context"`
}

type errorException struct {
	_ struct{} `validate:"requiredAnyOf=message;type"`
	// Attributes of the exception.
	Attributes map[string]any `json:"attributes"`
	// Code that is set when the error happened, e.g. database error code.
	Code nullable.Interface `json:"code" validate:"inputTypes=string;int,maxLength=1024"`
	// Cause can hold a collection of error exceptions representing chained
	// exceptions. The chain starts with the outermost exception, followed
	// by its cause, and so on.
	Cause []errorException `json:"cause"`
	// Stacktrace information of the captured exception.
	Stacktrace []stacktraceFrame `json:"stacktrace"`
	// Message contains the originally captured error message.
	Message nullable.String `json:"message"`
	// Module describes the exception type's module namespace.
	Module nullable.String `json:"module" validate:"maxLength=1024"`
	// Type of the exception.
	Type nullable.String `json:"type" validate:"maxLength=1024"`
	// Handled indicates whether the error was caught in the code or not.
	Handled nullable.Bool `json:"handled"`
}

type errorLog struct {
	// Level represents the severity of the recorded log.
	Level nullable.String `json:"level" validate:"maxLength=1024"`
	// LoggerName holds the name of the used logger instance.
	LoggerName nullable.String `json:"logger_name" validate:"maxLength=1024"`
	// Message of the logged error. In case a parameterized message is captured,
	// Message should contain the same information, but with any placeholders
	// being replaced.
	Message nullable.String `json:"message" validate:"required"`
	// ParamMessage should contain the same information as Message, but with
	// placeholders where parameters were logged, e.g. 'error connecting to %s'.
	// The string is not interpreted, allowing differnt placeholders per client
	// languange. The information might be used to group errors together.
	ParamMessage nullable.String `json:"param_message" validate:"maxLength=1024"`
	// Stacktrace information of the captured error.
	Stacktrace []stacktraceFrame `json:"stacktrace"`
}

type errorTransactionRef struct {
	// Name is the generic designation of a transaction in the scope of a
	// single service, eg: 'GET /users/:id'.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Type expresses the correlated transaction's type as keyword that has
	// specific relevance within the service's domain,
	// eg: 'request', 'backgroundjob'.
	Type nullable.String `json:"type" validate:"maxLength=1024"`
	// Sampled indicates whether or not the full information for a transaction
	// is captured. If a transaction is unsampled no spans and less context
	// information will be reported.
	Sampled nullable.Bool `json:"sampled"`
}

type metadata struct {
	// Labels are a flat mapping of user-defined tags. Allowed value types are
	// string, boolean and number values. Labels are indexed and searchable.
	Labels map[string]any `json:"labels" validate:"inputTypesVals=string;bool;number,maxLengthVals=1024"`
	// Service metadata about the monitored service.
	Service metadataService `json:"service" validate:"required"`
	// Cloud metadata about where the monitored service is running.
	Cloud metadataCloud `json:"cloud"`
	// System metadata
	System metadataSystem `json:"system"`
	// User metadata, which can be overwritten on a per event basis.
	User user `json:"user"`
	// Network holds information about the network over which the
	// monitored service is communicating.
	Network network `json:"network"`
	// Process metadata about the monitored service.
	Process metadataProcess `json:"process"`
}

type metadataCloud struct {
	// Account where the monitored service is running.
	Account metadataCloudAccount `json:"account"`
	// AvailabilityZone where the monitored service is running, e.g. us-east-1a
	AvailabilityZone nullable.String `json:"availability_zone" validate:"maxLength=1024"`
	// Instance on which the monitored service is running.
	Instance metadataCloudInstance `json:"instance"`
	// Machine on which the monitored service is running.
	Machine metadataCloudMachine `json:"machine"`
	// Project in which the monitored service is running.
	Project metadataCloudProject `json:"project"`
	// Provider that is used, e.g. aws, azure, gcp, digitalocean.
	Provider nullable.String `json:"provider" validate:"required,maxLength=1024"`
	// Region where the monitored service is running, e.g. us-east-1
	Region nullable.String `json:"region" validate:"maxLength=1024"`
	// Service that is monitored on cloud
	Service metadataCloudService `json:"service"`
}

type metadataCloudAccount struct {
	// ID of the cloud account.
	ID nullable.String `json:"id" validate:"maxLength=1024"`
	// Name of the cloud account.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
}

type metadataCloudInstance struct {
	// ID of the cloud instance.
	ID nullable.String `json:"id" validate:"maxLength=1024"`
	// Name of the cloud instance.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
}

type metadataCloudMachine struct {
	// ID of the cloud machine.
	Type nullable.String `json:"type" validate:"maxLength=1024"`
}

type metadataCloudProject struct {
	// ID of the cloud project.
	ID nullable.String `json:"id" validate:"maxLength=1024"`
	// Name of the cloud project.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
}

type metadataCloudService struct {
	// Name of the cloud service, intended to distinguish services running on
	// different platforms within a provider, eg AWS EC2 vs Lambda,
	// GCP GCE vs App Engine, Azure VM vs App Server.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
}

type metadataProcess struct {
	// Argv holds the command line arguments used to start this process.
	Argv []string `json:"argv"`
	// Title is the process title. It can be the same as process name.
	Title nullable.String `json:"title" validate:"maxLength=1024"`
	// PID holds the process ID of the service.
	Pid nullable.Int `json:"pid" validate:"required"`
	// Ppid holds the parent process ID of the service.
	Ppid nullable.Int `json:"ppid"`
}

type metadataService struct {
	// Agent holds information about the APM agent capturing the event.
	Agent metadataServiceAgent `json:"agent" validate:"required"`
	// Environment in which the monitored service is running,
	// e.g. `production` or `staging`.
	Environment nullable.String `json:"environment" validate:"maxLength=1024"`
	// Framework holds information about the framework used in the
	// monitored service.
	Framework metadataServiceFramework `json:"framework"`
	// ID holds a unique identifier for the running service.
	ID nullable.String `json:"id"`
	// Language holds information about the programming language of the
	// monitored service.
	Language metadataServiceLanguage `json:"language"`
	// Name of the monitored service.
	Name nullable.String `json:"name" validate:"required,minLength=1,maxLength=1024,pattern=patternAlphaNumericExt"`
	// Node must be a unique meaningful name of the service node.
	Node metadataServiceNode `json:"node"`
	// Runtime holds information about the language runtime running the
	// monitored service
	Runtime metadataServiceRuntime `json:"runtime"`
	// Version of the monitored service.
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type metadataServiceAgent struct {
	// ActivationMethod of the APM agent capturing information.
	ActivationMethod nullable.String `json:"activation_method" validate:"maxLength=1024"`
	// EphemeralID is a free format ID used for metrics correlation by agents
	EphemeralID nullable.String `json:"ephemeral_id" validate:"maxLength=1024"`
	// Name of the APM agent capturing information.
	Name nullable.String `json:"name" validate:"required,minLength=1,maxLength=1024"`
	// Version of the APM agent capturing information.
	Version nullable.String `json:"version" validate:"required,maxLength=1024"`
}

type metadataServiceFramework struct {
	// Name of the used framework
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Version of the used framework
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type metadataServiceLanguage struct {
	// Name of the used programming language
	Name nullable.String `json:"name" validate:"required,maxLength=1024"`
	// Version of the used programming language
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type metadataServiceNode struct {
	// Name of the service node
	Name nullable.String `json:"configured_name" validate:"maxLength=1024"`
}

type metadataServiceRuntime struct {
	// Name of the language runtime
	Name nullable.String `json:"name" validate:"required,maxLength=1024"`
	// Name of the language runtime
	Version nullable.String `json:"version" validate:"required,maxLength=1024"`
}

type metadataSystem struct {
	// Architecture of the system the monitored service is running on.
	Architecture nullable.String `json:"architecture" validate:"maxLength=1024"`
	// ConfiguredHostname is the configured name of the host the monitored
	// service is running on. It should only be sent when configured by the
	// user. If given, it is used as the event's hostname.
	ConfiguredHostname nullable.String `json:"configured_hostname" validate:"maxLength=1024"`
	// Container holds the system's container ID if available.
	Container metadataSystemContainer `json:"container"`
	// DetectedHostname is the hostname detected by the APM agent. It usually
	// contains what the hostname command returns on the host machine.
	// It will be used as the event's hostname if ConfiguredHostname is not present.
	DetectedHostname nullable.String `json:"detected_hostname" validate:"maxLength=1024"`
	// Deprecated: Use ConfiguredHostname and DetectedHostname instead.
	// DeprecatedHostname is the host name of the system the service is
	// running on. It does not distinguish between configured and detected
	// hostname and therefore is deprecated and only used if no other hostname
	// information is available.
	DeprecatedHostname nullable.String `json:"hostname" validate:"maxLength=1024"`
	// Kubernetes system information if the monitored service runs on Kubernetes.
	Kubernetes metadataSystemKubernetes `json:"kubernetes"`
	// Platform name of the system platform the monitored service is running on.
	Platform nullable.String `json:"platform" validate:"maxLength=1024"`
	// The OpenTelemetry semantic conventions compliant "host.id" attribute, if available.
	HostID nullable.String `json:"host_id" validate:"maxLength=1024"`
}

type metadataSystemContainer struct {
	// ID of the container the monitored service is running in.
	ID nullable.String `json:"id" validate:"maxLength=1024"`
}

type metadataSystemKubernetes struct {
	// Namespace of the Kubernetes resource the monitored service is run on.
	Namespace nullable.String `json:"namespace" validate:"maxLength=1024"`
	// Node related information
	Node metadataSystemKubernetesNode `json:"node"`
	// Pod related information
	Pod metadataSystemKubernetesPod `json:"pod"`
}

type metadataSystemKubernetesNode struct {
	// Name of the Kubernetes Node
	Name nullable.String `json:"name" validate:"maxLength=1024"`
}

type metadataSystemKubernetesPod struct {
	// Name of the Kubernetes Pod
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// UID is the system-generated string uniquely identifying the Pod.
	UID nullable.String `json:"uid" validate:"maxLength=1024"`
}

type network struct {
	Connection networkConnection `json:"connection"`
}

type networkConnection struct {
	Type nullable.String `json:"type" validate:"maxLength=1024"`
}

type metricset struct {
	// Timestamp holds the recorded time of the event, UTC based and formatted
	// as microseconds since Unix epoch
	Timestamp nullable.TimeMicrosUnix `json:"timestamp"`
	// Samples hold application metrics collected from the agent.
	Samples map[string]metricsetSampleValue `json:"samples" validate:"required,patternKeys=patternNoAsteriskQuote"`
	// Span holds selected information about the correlated transaction.
	Span metricsetSpanRef `json:"span"`
	// Tags are a flat mapping of user-defined tags. On the agent side, tags
	// are called labels. Allowed value types are string, boolean and number
	// values. Tags are indexed and searchable.
	Tags map[string]any `json:"tags" validate:"inputTypesVals=string;bool;number,maxLengthVals=1024"`
	// Transaction holds selected information about the correlated transaction.
	Transaction metricsetTransactionRef `json:"transaction"`
	// Service holds selected information about the correlated service.
	Service metricsetServiceRef `json:"service"`
	// FAAS holds fields related to Function as a Service events.
	FAAS faas `json:"faas"`
}

type metricsetSampleValue struct {
	// At least one of value or values must be specified.
	_ struct{} `validate:"requiredAnyOf=value;values"`
	// Type holds an optional metric type: gauge, counter, or histogram.
	//
	// If Type is unknown, it will be ignored.
	Type nullable.String `json:"type"`

	// Unit holds an optional unit for the metric.
	//
	// - "percent" (value is in the range [0,1])
	// - "byte"
	// - a time unit: "nanos", "micros", "ms", "s", "m", "h", "d"
	//
	// If Unit is unknown, it will be ignored.
	Unit nullable.String `json:"unit"`

	// Values holds the bucket values for histogram metrics.
	//
	// Values must be provided in ascending order; failure to do
	// so will result in the metric being discarded.
	Values []float64 `json:"values" validate:"requiredIfAny=counts"`

	// Counts holds the bucket counts for histogram metrics.
	//
	// These numbers must be positive or zero.
	//
	// If Counts is specified, then Values is expected to be
	// specified with the same number of elements, and with the
	// same order.
	Counts []uint64 `json:"counts" validate:"requiredIfAny=values"`
	// Value holds the value of a single metric sample.
	Value nullable.Float64 `json:"value"`
}

type metricsetSpanRef struct {
	// Subtype is a further sub-division of the type (e.g. postgresql, elasticsearch)
	Subtype nullable.String `json:"subtype" validate:"maxLength=1024"`
	// Type expresses the correlated span's type as keyword that has specific
	// relevance within the service's domain, eg: 'request', 'backgroundjob'.
	Type nullable.String `json:"type" validate:"maxLength=1024"`
}

type metricsetTransactionRef struct {
	// Name of the correlated transaction.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Type expresses the correlated transaction's type as keyword that has specific
	// relevance within the service's domain, eg: 'request', 'backgroundjob'.
	Type nullable.String `json:"type" validate:"maxLength=1024"`
}

type metricsetServiceRef struct {
	// Name of the correlated service.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Version of the correlated service.
	Version nullable.String `json:"version" validate:"maxLength=1024"`
}

type span struct {
	_ struct{} `validate:"requiredAnyOf=start;timestamp"`
	// Timestamp holds the recorded time of the event, UTC based and formatted
	// as microseconds since Unix epoch
	Timestamp nullable.TimeMicrosUnix `json:"timestamp"`
	// OTel contains unmapped OpenTelemetry attributes.
	OTel otel `json:"otel"`
	// ID holds the hex encoded 64 random bits ID of the event.
	ID nullable.String `json:"id" validate:"required,maxLength=1024"`
	// TraceID holds the hex encoded 128 random bits ID of the correlated trace.
	TraceID nullable.String `json:"trace_id" validate:"required,maxLength=1024"`
	// Action holds the specific kind of event within the sub-type represented
	// by the span (e.g. query, connect)
	Action nullable.String `json:"action" validate:"maxLength=1024"`
	// Name is the generic designation of a span in the scope of a transaction.
	Name nullable.String `json:"name" validate:"required,maxLength=1024"`
	// Outcome of the span: success, failure, or unknown. Outcome may be one of
	// a limited set of permitted values describing the success or failure of
	// the span. It can be used for calculating error rates for outgoing requests.
	Outcome nullable.String `json:"outcome" validate:"enum=enumOutcome"`
	// ChildIDs holds a list of successor transactions and/or spans.
	ChildIDs []string `json:"child_ids" validate:"maxLength=1024"`
	// ParentID holds the hex encoded 64 random bits ID of the parent
	// transaction or span.
	ParentID nullable.String `json:"parent_id" validate:"required,maxLength=1024"`
	// Links holds links to other spans, potentially in other traces.
	Links []spanLink `json:"links"`
	// Stacktrace connected to this span event.
	Stacktrace []stacktraceFrame `json:"stacktrace"`
	// Type holds the span's type, and can have specific keywords
	// within the service's domain (eg: 'request', 'backgroundjob', etc)
	Type nullable.String `json:"type" validate:"required,maxLength=1024"`
	// Subtype is a further sub-division of the type (e.g. postgresql, elasticsearch)
	Subtype nullable.String `json:"subtype" validate:"maxLength=1024"`
	// TransactionID holds the hex encoded 64 random bits ID of the correlated
	// transaction.
	TransactionID nullable.String `json:"transaction_id" validate:"maxLength=1024"`
	// Composite holds details on a group of spans represented by a single one.
	Composite spanComposite `json:"composite"`
	// Context holds arbitrary contextual information for the event.
	Context spanContext `json:"context"`
	// Start is the offset relative to the transaction's timestamp identifying
	// the start of the span, in milliseconds.
	Start nullable.Float64 `json:"start"`
	// SampleRate applied to the monitored service at the time where this span
	// was recorded.
	SampleRate nullable.Float64 `json:"sample_rate"`
	// Duration of the span in milliseconds. When the span is a composite one,
	// duration is the gross duration, including "whitespace" in between spans.
	Duration nullable.Float64 `json:"duration" validate:"required,min=0"`
	// Sync indicates whether the span was executed synchronously or asynchronously.
	Sync nullable.Bool `json:"sync"`
}

type spanContext struct {
	// Tags are a flat mapping of user-defined tags. On the agent side, tags
	// are called labels. Allowed value types are string, boolean and number
	// values. Tags are indexed and searchable.
	Tags map[string]any `json:"tags" validate:"inputTypesVals=string;bool;number,maxLengthVals=1024"`
	// Service related information can be sent per span. Information provided
	// here will override the more generic information retrieved from metadata,
	// missing service fields will be retrieved from the metadata information.
	Service contextService `json:"service"`
	// Message holds details related to message receiving and publishing
	// if the captured event integrates with a messaging system
	Message contextMessage `json:"message"`
	// Database contains contextual data for database spans
	Database spanContextDatabase `json:"db"`
	// Destination contains contextual data about the destination of spans
	Destination spanContextDestination `json:"destination"`
	// HTTP contains contextual information when the span concerns an HTTP request.
	HTTP spanContextHTTP `json:"http"`
}

type spanContextDatabase struct {
	// Instance name of the database.
	Instance nullable.String `json:"instance"`
	// Link to the database server.
	Link nullable.String `json:"link" validate:"maxLength=1024"`
	// Statement of the recorded database event, e.g. query.
	Statement nullable.String `json:"statement"`
	// Type of the recorded database event., e.g. sql, cassandra, hbase, redis.
	Type nullable.String `json:"type"`
	// User is the username with which the database is accessed.
	User nullable.String `json:"user"`
	// RowsAffected shows the number of rows affected by the statement.
	RowsAffected nullable.Int `json:"rows_affected"`
}

type spanContextDestination struct {
	// Service describes the destination service
	Service spanContextDestinationService `json:"service"`
	// Address is the destination network address:
	// hostname (e.g. 'localhost'),
	// FQDN (e.g. 'elastic.co'),
	// IPv4 (e.g. '127.0.0.1')
	// IPv6 (e.g. '::1')
	Address nullable.String `json:"address" validate:"maxLength=1024"`
	// Port is the destination network port (e.g. 443)
	Port nullable.Int `json:"port"`
}

type spanContextDestinationService struct {
	// Name is the identifier for the destination service,
	// e.g. 'http://elastic.co', 'elasticsearch', 'rabbitmq' (
	// DEPRECATED: this field will be removed in a future release
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Resource identifies the destination service resource being operated on
	// e.g. 'http://elastic.co:80', 'elasticsearch', 'rabbitmq/queue_name'
	// DEPRECATED: this field will be removed in a future release
	Resource nullable.String `json:"resource" validate:"required,maxLength=1024"`
	// Type of the destination service, e.g. db, elasticsearch. Should
	// typically be the same as span.type.
	// DEPRECATED: this field will be removed in a future release
	Type nullable.String `json:"type" validate:"maxLength=1024"`
}

type spanContextHTTP struct {
	// Request describes the HTTP request information.
	Request spanContextHTTPRequest `json:"request"`
	// Method holds information about the method of the HTTP request.
	Method nullable.String `json:"method" validate:"maxLength=1024"`
	// URL is the raw url of the correlating HTTP request.
	URL nullable.String `json:"url"`
	// Response describes the HTTP response information in case the event was
	// created as a result of an HTTP request.
	Response spanContextHTTPResponse `json:"response"`
	// Deprecated: Use Response.StatusCode instead.
	// StatusCode sent in the http response.
	StatusCode nullable.Int `json:"status_code"`
}

type spanContextHTTPRequest struct {
	// The http request body usually as a string, but may be a dictionary for multipart/form-data content
	Body nullable.Interface `json:"body"`
	// ID holds the unique identifier for the http request.
	ID nullable.String `json:"id"`
}

type spanContextHTTPResponse struct {
	// Headers holds the http headers sent in the http response.
	Headers nullable.HTTPHeader `json:"headers"`
	// DecodedBodySize holds the size of the decoded payload.
	DecodedBodySize nullable.Int `json:"decoded_body_size"`
	// EncodedBodySize holds the size of the encoded payload.
	EncodedBodySize nullable.Int `json:"encoded_body_size"`
	// StatusCode sent in the http response.
	StatusCode nullable.Int `json:"status_code"`
	// TransferSize holds the total size of the payload.
	TransferSize nullable.Int `json:"transfer_size"`
}

type stacktraceFrame struct {
	_ struct{} `validate:"requiredAnyOf=classname;filename"`
	// Vars is a flat mapping of local variables of the frame.
	Vars map[string]any `json:"vars"`
	// Filename is the relative name of the frame's file.
	Filename nullable.String `json:"filename"`
	// AbsPath is the absolute path of the frame's file.
	AbsPath nullable.String `json:"abs_path"`
	// Classname of the frame.
	Classname nullable.String `json:"classname"`
	// ContextLine is the line from the frame's file.
	ContextLine nullable.String `json:"context_line"`
	// Function represented by the frame.
	Function nullable.String `json:"function"`
	// Module to which the frame belongs to.
	Module nullable.String `json:"module"`
	// PostContext is a slice of code lines immediately before the line
	// from the frame's file.
	PostContext []string `json:"post_context"`
	// PreContext is a slice of code lines immediately after the line
	// from the frame's file.
	PreContext []string `json:"pre_context"`
	// LineNumber of the frame.
	LineNumber nullable.Int `json:"lineno"`
	// ColumnNumber of the frame.
	ColumnNumber nullable.Int `json:"colno"`
	// LibraryFrame indicates whether the frame is from a third party library.
	LibraryFrame nullable.Bool `json:"library_frame"`
}

type spanComposite struct {
	// A string value indicating which compression strategy was used. The valid
	// values are `exact_match` and `same_kind`.
	CompressionStrategy nullable.String `json:"compression_strategy" validate:"required"`
	// Count is the number of compressed spans the composite span represents.
	// The minimum count is 2, as a composite span represents at least two spans.
	Count nullable.Int `json:"count" validate:"required,min=2"`
	// Sum is the durations of all compressed spans this composite span
	// represents in milliseconds.
	Sum nullable.Float64 `json:"sum" validate:"required,min=0"`
}

type transaction struct {
	// Marks capture the timing of a significant event during the lifetime of
	// a transaction. Marks are organized into groups and can be set by the
	// user or the agent. Marks are only reported by RUM agents.
	Marks transactionMarks `json:"marks"`
	// Timestamp holds the recorded time of the event, UTC based and formatted
	// as microseconds since Unix epoch
	Timestamp nullable.TimeMicrosUnix `json:"timestamp"`
	// OTel contains unmapped OpenTelemetry attributes.
	OTel otel `json:"otel"`
	// Links holds links to other spans, potentially in other traces.
	Links []spanLink `json:"links"`
	// TraceID holds the hex encoded 128 random bits ID of the correlated trace.
	TraceID nullable.String `json:"trace_id" validate:"required,maxLength=1024"`
	// ID holds the hex encoded 64 random bits ID of the event.
	ID nullable.String `json:"id" validate:"required,maxLength=1024"`
	// ParentID holds the hex encoded 64 random bits ID of the parent
	// transaction or span.
	ParentID nullable.String `json:"parent_id" validate:"maxLength=1024"`
	// Name is the generic designation of a transaction in the scope of a
	// single service, eg: 'GET /users/:id'.
	Name nullable.String `json:"name" validate:"maxLength=1024"`
	// Type expresses the transaction's type as keyword that has specific
	// relevance within the service's domain, eg: 'request', 'backgroundjob'.
	Type nullable.String `json:"type" validate:"required,maxLength=1024"`
	// Result of the transaction. For HTTP-related transactions, this should
	// be the status code formatted like 'HTTP 2xx'.
	Result nullable.String `json:"result" validate:"maxLength=1024"`
	// DroppedSpanStats holds information about spans that were dropped
	// (for example due to transaction_max_spans or exit_span_min_duration).
	DroppedSpanStats []transactionDroppedSpanStats `json:"dropped_spans_stats"`
	// Outcome of the transaction with a limited set of permitted values,
	// describing the success or failure of the transaction from the service's
	// perspective. It is used for calculating error rates for incoming requests.
	// Permitted values: success, failure, unknown.
	Outcome nullable.String `json:"outcome" validate:"enum=enumOutcome"`
	// FAAS holds fields related to Function as a Service events.
	FAAS faas `json:"faas"`
	// Session holds optional transaction session information for RUM.
	Session transactionSession `json:"session"`
	// Context holds arbitrary contextual information for the event.
	Context context `json:"context"`
	// UserExperience holds metrics for measuring real user experience.
	// This information is only sent by RUM agents.
	UserExperience transactionUserExperience `json:"experience"`
	// SpanCount counts correlated spans.
	SpanCount transactionSpanCount `json:"span_count" validate:"required"`
	// SampleRate applied to the monitored service at the time where this transaction
	// was recorded. Allowed values are [0..1]. A SampleRate <1 indicates that
	// not all spans are recorded.
	SampleRate nullable.Float64 `json:"sample_rate"`
	// Duration how long the transaction took to complete, in milliseconds
	// with 3 decimal points.
	Duration nullable.Float64 `json:"duration" validate:"required,min=0"`
	// Sampled indicates whether or not the full information for a transaction
	// is captured. If a transaction is unsampled no spans and less context
	// information will be reported.
	Sampled nullable.Bool `json:"sampled"`
}

type log struct {
	// Labels are a flat mapping of user-defined key-value pairs.
	Labels map[string]any `json:"labels" validate:"inputTypesVals=string;bool;number,maxLengthVals=1024"`
	// Timestamp holds the recorded time of the event, UTC based and formatted
	// as microseconds since Unix epoch
	Timestamp nullable.TimeMicrosUnix `json:"@timestamp"`
	// Below embedded fields are added to enable supporting both nested and flat JSON.
	// This is achieved by generating code using static analysis of these structs.
	// The logic parses JSON tag of each struct field to produce a code which, at runtime,
	// checks the nested map to retrieve the required value for each field.
	EcsLogServiceFields
	EcsLogErrorFields
	EcsLogEventFields
	EcsLogProcessFields

	// TraceID holds the ID of the correlated trace.
	TraceID nullable.String `json:"trace.id" validate:"maxLength=1024"`
	// TransactionID holds the ID of the correlated transaction.
	TransactionID nullable.String `json:"transaction.id" validate:"maxLength=1024"`
	// SpanID holds the ID of the correlated span.
	SpanID nullable.String `json:"span.id" validate:"maxLength=1024"`
	// Message logged as part of the log. In case a parameterized message is
	// captured, Message should contain the same information, but with any placeholders
	// being replaced.
	Message nullable.String `json:"message"`
	// FAAS holds fields related to Function as a Service events.
	FAAS faas `json:"faas"`
	// Below embedded fields are added to enable supporting both nested and flat JSON.
	// This is achieved by generating code using static analysis of these structs.
	// The logic parses JSON tag of each struct field to produce a code which, at runtime,
	// checks the nested map to retrieve the required value for each field.
	EcsLogLogFields
}

// EcsLogEventFields holds event.* fields for supporting ECS logging format and enables
// parsing them in flat as well as nested notation.
type EcsLogEventFields struct {
	NestedStruct map[string]interface{} `json:"event" nested:"true"`
	// ProcessThreadName represents the name of the thread.
	EventDataset nullable.String `json:"event.dataset" validate:"maxLength=1024"`
}

// EcsLogProcessFields holds process.* fields for supporting ECS logging format and enables
// parsing them in flat as well as nested notation.
type EcsLogProcessFields struct {
	NestedStruct map[string]interface{} `json:"process" nested:"true"`
	// ProcessThreadName represents the name of the thread.
	ProcessThreadName nullable.String `json:"process.thread.name" validate:"maxLength=1024"`
}

// EcsLogErrorFields holds error.* fields for supporting ECS logging format and enables
// parsing them in flat as well as nested notation.
type EcsLogErrorFields struct {
	NestedStruct map[string]interface{} `json:"error" nested:"true"`
	// ErrorType represents the type of the error if the log line represents an error.
	ErrorType nullable.String `json:"error.type"`
	// ErrorMessage represents the message contained in the error if the log line
	// represents an error.
	ErrorMessage nullable.String `json:"error.message"`
	// ErrorStacktrace represents the plain text stacktrace of the error the log line
	// represents.
	ErrorStacktrace nullable.String `json:"error.stack_trace"`
}

// EcsLogLogFields holds log.* fields for supporting ECS logging format and enables
// parsing them in flat as well as nested notation.
type EcsLogLogFields struct {
	NestedStruct map[string]interface{} `json:"log" nested:"true"`
	// Level represents the severity of the recorded log.
	Level nullable.String `json:"log.level" validate:"maxLength=1024"`
	// Logger represents the name of the used logger instance.
	Logger nullable.String `json:"log.logger" validate:"maxLength=1024"`
	// OriginFileName represents the filename containing the sourcecode where the log
	// originated.
	OriginFileName nullable.String `json:"log.origin.file.name" validate:"maxLength=1024"`
	// OriginFunction represents the function name where the log originated.
	OriginFunction nullable.String `json:"log.origin.function"`
	// OriginFileLine represents the line number in the file containing the sourcecode
	// where the log originated.
	OriginFileLine nullable.Int `json:"log.origin.file.line"`
}

// EcsLogServiceFields holds service.* fields for supporting ECS logging format and
// enables parsing them in flat as well as nested notation.
type EcsLogServiceFields struct {
	NestedStruct map[string]interface{} `json:"service" nested:"true"`
	// ServiceName represents name of the service which originated the log line.
	ServiceName nullable.String `json:"service.name" validate:"maxLength=1024"`
	// ServiceVersion represents the version of the service which originated the log
	// line.
	ServiceVersion nullable.String `json:"service.version" validate:"maxLength=1024"`
	// ServiceEnvironment represents the environment the service which originated the
	// log line is running in.
	ServiceEnvironment nullable.String `json:"service.environment" validate:"maxLength=1024"`
	// ServiceNodeName represents a unique node name per host for the service which
	// originated the log line.
	ServiceNodeName nullable.String `json:"service.node.name" validate:"maxLength=1024"`
}

type otel struct {
	// Attributes hold the unmapped OpenTelemetry attributes.
	Attributes map[string]interface{} `json:"attributes"`
	// SpanKind holds the incoming OpenTelemetry span kind.
	SpanKind nullable.String `json:"span_kind"`
}

type transactionSession struct {
	// ID holds a session ID for grouping a set of related transactions.
	ID nullable.String `json:"id" validate:"required,maxLength=1024"`

	// Sequence holds an optional sequence number for a transaction within
	// a session. It is not meaningful to compare sequences across two
	// different sessions.
	Sequence nullable.Int `json:"sequence" validate:"min=1"`
}

type transactionMarks struct {
	Events map[string]transactionMarkEvents `json:"-"`
}

func (m *transactionMarks) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Events)
}

type transactionMarkEvents struct {
	Measurements map[string]float64 `json:"-"`
}

func (m *transactionMarkEvents) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Measurements)
}

type transactionSpanCount struct {
	// Dropped is the number of correlated spans that have been dropped by
	// the APM agent recording the transaction.
	Dropped nullable.Int `json:"dropped"`
	// Started is the number of correlated spans that are recorded.
	Started nullable.Int `json:"started" validate:"required"`
}

// transactionUserExperience holds real user (browser) experience metrics.
type transactionUserExperience struct {
	// CumulativeLayoutShift holds the Cumulative Layout Shift (CLS) metric value,
	// or a negative value if CLS is unknown. See https://web.dev/cls/
	CumulativeLayoutShift nullable.Float64 `json:"cls" validate:"min=0"`
	// FirstInputDelay holds the First Input Delay (FID) metric value,
	// or a negative value if FID is unknown. See https://web.dev/fid/
	FirstInputDelay nullable.Float64 `json:"fid" validate:"min=0"`
	// Longtask holds longtask duration/count metrics.
	Longtask longtaskMetrics `json:"longtask"`
	// TotalBlockingTime holds the Total Blocking Time (TBT) metric value,
	// or a negative value if TBT is unknown. See https://web.dev/tbt/
	TotalBlockingTime nullable.Float64 `json:"tbt" validate:"min=0"`
}

type longtaskMetrics struct {
	// Count is the total number of of longtasks.
	Count nullable.Int64 `json:"count" validate:"required,min=0"`
	// Max longtask duration
	Max nullable.Float64 `json:"max" validate:"required,min=0"`
	// Sum of longtask durations
	Sum nullable.Float64 `json:"sum" validate:"required,min=0"`
}

type user struct {
	// Domain of the logged in user
	Domain nullable.String `json:"domain" validate:"maxLength=1024"`
	// ID identifies the logged in user, e.g. can be the primary key of the user
	ID nullable.Interface `json:"id" validate:"maxLength=1024,inputTypes=string;int"`
	// Email of the user.
	Email nullable.String `json:"email" validate:"maxLength=1024"`
	// Name of the user.
	Name nullable.String `json:"username" validate:"maxLength=1024"`
}

type transactionDroppedSpanStats struct {
	// DestinationServiceResource identifies the destination service resource
	// being operated on. e.g. 'http://elastic.co:80', 'elasticsearch', 'rabbitmq/queue_name'.
	DestinationServiceResource nullable.String `json:"destination_service_resource" validate:"maxLength=1024"`
	// ServiceTargetType identifies the type of the target service being operated on
	// e.g. 'oracle', 'rabbitmq'
	ServiceTargetType nullable.String `json:"service_target_type" validate:"maxLength=512"`
	// ServiceTargetName identifies the instance name of the target service being operated on
	ServiceTargetName nullable.String `json:"service_target_name" validate:"maxLength=512"`
	// Outcome of the span: success, failure, or unknown. Outcome may be one of
	// a limited set of permitted values describing the success or failure of
	// the span. It can be used for calculating error rates for outgoing requests.
	Outcome nullable.String `json:"outcome" validate:"enum=enumOutcome"`
	// Duration holds duration aggregations about the dropped span.
	Duration transactionDroppedSpansDuration `json:"duration"`
}

type transactionDroppedSpansDuration struct {
	// Count holds the number of times the dropped span happened.
	Count nullable.Int `json:"count" validate:"min=1"`
	// Sum holds dimensions about the dropped span's duration.
	Sum transactionDroppedSpansDurationSum `json:"sum"`
}

type transactionDroppedSpansDurationSum struct {
	// Us represents the summation of the span duration.
	Us nullable.Int `json:"us" validate:"min=0"`
}

type spanLink struct {
	// SpanID holds the ID of the linked span.
	SpanID nullable.String `json:"span_id" validate:"required,maxLength=1024"`

	// TraceID holds the ID of the linked span's trace.
	TraceID nullable.String `json:"trace_id" validate:"required,maxLength=1024"`
}
