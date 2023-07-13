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
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/apm-data/input/elasticapm/internal/decoder"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/modeldecoderutil"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/netutil"
	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/nullable"
	"github.com/elastic/apm-data/input/otlp"
	"github.com/elastic/apm-data/model/modelpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var (
	errorRootPool = sync.Pool{
		New: func() interface{} {
			return &errorRoot{}
		},
	}
	metadataRootPool = sync.Pool{
		New: func() interface{} {
			return &metadataRoot{}
		},
	}
	metricsetRootPool = sync.Pool{
		New: func() interface{} {
			return &metricsetRoot{}
		},
	}
	spanRootPool = sync.Pool{
		New: func() interface{} {
			return &spanRoot{}
		},
	}
	transactionRootPool = sync.Pool{
		New: func() interface{} {
			return &transactionRoot{}
		},
	}
	logRootPool = sync.Pool{
		New: func() interface{} {
			return &logRoot{}
		},
	}
)

var compressionStrategyText = map[string]modelpb.CompressionStrategy{
	"exact_match": modelpb.CompressionStrategy_COMPRESSION_STRATEGY_EXACT_MATCH,
	"same_kind":   modelpb.CompressionStrategy_COMPRESSION_STRATEGY_SAME_KIND,
}

var metricTypeText = map[string]modelpb.MetricType{
	"gauge":     modelpb.MetricType_METRIC_TYPE_GAUGE,
	"counter":   modelpb.MetricType_METRIC_TYPE_COUNTER,
	"histogram": modelpb.MetricType_METRIC_TYPE_HISTOGRAM,
	"summary":   modelpb.MetricType_METRIC_TYPE_SUMMARY,
}

var (
	// reForServiceTargetExpr regex will capture service target type and name
	// Service target type comprises of only lowercase alphabets
	// Service target name comprises of all word characters
	reForServiceTargetExpr = regexp.MustCompile(`^([a-z0-9]+)(?:/(\w+))?$`)
)

func fetchErrorRoot() *errorRoot {
	return errorRootPool.Get().(*errorRoot)
}

func releaseErrorRoot(root *errorRoot) {
	root.Reset()
	errorRootPool.Put(root)
}

func fetchMetadataRoot() *metadataRoot {
	return metadataRootPool.Get().(*metadataRoot)
}

func releaseMetadataRoot(root *metadataRoot) {
	root.Reset()
	metadataRootPool.Put(root)
}

func fetchMetricsetRoot() *metricsetRoot {
	return metricsetRootPool.Get().(*metricsetRoot)
}

func releaseMetricsetRoot(root *metricsetRoot) {
	root.Reset()
	metricsetRootPool.Put(root)
}

func fetchSpanRoot() *spanRoot {
	return spanRootPool.Get().(*spanRoot)
}

func releaseSpanRoot(root *spanRoot) {
	root.Reset()
	spanRootPool.Put(root)
}

func fetchTransactionRoot() *transactionRoot {
	return transactionRootPool.Get().(*transactionRoot)
}

func releaseTransactionRoot(root *transactionRoot) {
	root.Reset()
	transactionRootPool.Put(root)
}

func fetchLogRoot() *logRoot {
	return logRootPool.Get().(*logRoot)
}

func releaseLogRoot(root *logRoot) {
	root.Reset()
	logRootPool.Put(root)
}

// DecodeMetadata decodes metadata from d, updating out.
//
// DecodeMetadata should be used when the the stream in the decoder does not contain the
// `metadata` key, but only the metadata data.
func DecodeMetadata(d decoder.Decoder, out *modelpb.APMEvent) error {
	return decodeMetadata(decodeIntoMetadata, d, out)
}

// DecodeNestedMetadata decodes metadata from d, updating out.
//
// DecodeNestedMetadata should be used when the stream in the decoder contains the `metadata` key
func DecodeNestedMetadata(d decoder.Decoder, out *modelpb.APMEvent) error {
	return decodeMetadata(decodeIntoMetadataRoot, d, out)
}

// DecodeNestedError decodes an error from d, appending it to batch.
//
// DecodeNestedError should be used when the stream in the decoder contains the `error` key
func DecodeNestedError(d decoder.Decoder, input *modeldecoder.Input, batch *modelpb.Batch) error {
	root := fetchErrorRoot()
	defer releaseErrorRoot(root)
	err := d.Decode(root)
	if err != nil && err != io.EOF {
		return modeldecoder.NewDecoderErrFromJSONIter(err)
	}
	if err := root.validate(); err != nil {
		return modeldecoder.NewValidationErr(err)
	}
	event := input.Base.CloneVT()
	mapToErrorModel(&root.Error, event)
	*batch = append(*batch, event)
	return err
}

// DecodeNestedMetricset decodes a metricset from d, appending it to batch.
//
// DecodeNestedMetricset should be used when the stream in the decoder contains the `metricset` key
func DecodeNestedMetricset(d decoder.Decoder, input *modeldecoder.Input, batch *modelpb.Batch) error {
	root := fetchMetricsetRoot()
	defer releaseMetricsetRoot(root)
	var err error
	if err = d.Decode(root); err != nil && err != io.EOF {
		return modeldecoder.NewDecoderErrFromJSONIter(err)
	}
	if err := root.validate(); err != nil {
		return modeldecoder.NewValidationErr(err)
	}
	event := input.Base.CloneVT()
	if mapToMetricsetModel(&root.Metricset, event) {
		*batch = append(*batch, event)
	}
	return err
}

// DecodeNestedSpan decodes a span from d, appending it to batch.
//
// DecodeNestedSpan should be used when the stream in the decoder contains the `span` key
func DecodeNestedSpan(d decoder.Decoder, input *modeldecoder.Input, batch *modelpb.Batch) error {
	root := fetchSpanRoot()
	defer releaseSpanRoot(root)
	var err error
	if err = d.Decode(root); err != nil && err != io.EOF {
		return modeldecoder.NewDecoderErrFromJSONIter(err)
	}
	if err := root.validate(); err != nil {
		return modeldecoder.NewValidationErr(err)
	}
	event := input.Base.CloneVT()
	mapToSpanModel(&root.Span, event)
	*batch = append(*batch, event)
	return err
}

// DecodeNestedTransaction decodes a transaction from d, appending it to batch.
//
// DecodeNestedTransaction should be used when the stream in the decoder contains the `transaction` key
func DecodeNestedTransaction(d decoder.Decoder, input *modeldecoder.Input, batch *modelpb.Batch) error {
	root := fetchTransactionRoot()
	defer releaseTransactionRoot(root)
	var err error
	if err = d.Decode(root); err != nil && err != io.EOF {
		return modeldecoder.NewDecoderErrFromJSONIter(err)
	}
	if err := root.validate(); err != nil {
		return modeldecoder.NewValidationErr(err)
	}
	event := input.Base.CloneVT()
	mapToTransactionModel(&root.Transaction, event)
	*batch = append(*batch, event)
	return err
}

// DecodeNestedLog decodes a log event from d, appending it to batch.
//
// DecodeNestedLog should be used when the stream in the decoder contains the `log` key
func DecodeNestedLog(d decoder.Decoder, input *modeldecoder.Input, batch *modelpb.Batch) error {
	root := fetchLogRoot()
	defer releaseLogRoot(root)
	var err error
	if err = d.Decode(root); err != nil && err != io.EOF {
		return modeldecoder.NewDecoderErrFromJSONIter(err)
	}
	// Flatten any nested source to set the values for the flat fields
	if err := root.processNestedSource(); err != nil {
		return err
	}
	if err := root.validate(); err != nil {
		return modeldecoder.NewValidationErr(err)
	}
	event := input.Base.CloneVT()
	mapToLogModel(&root.Log, event)
	*batch = append(*batch, event)
	return err
}

func decodeMetadata(decFn func(d decoder.Decoder, m *metadataRoot) error, d decoder.Decoder, out *modelpb.APMEvent) error {
	m := fetchMetadataRoot()
	defer releaseMetadataRoot(m)
	var err error
	if err = decFn(d, m); err != nil && err != io.EOF {
		return modeldecoder.NewDecoderErrFromJSONIter(err)
	}
	if err := m.validate(); err != nil {
		return modeldecoder.NewValidationErr(err)
	}
	mapToMetadataModel(&m.Metadata, out)
	return err
}

func decodeIntoMetadata(d decoder.Decoder, m *metadataRoot) error {
	return d.Decode(&m.Metadata)
}

func decodeIntoMetadataRoot(d decoder.Decoder, m *metadataRoot) error {
	return d.Decode(m)
}

func mapToFAASModel(from faas, faas *modelpb.Faas) {
	if from.ID.IsSet() {
		faas.Id = from.ID.Val
	}
	if from.Coldstart.IsSet() {
		valCopy := from.Coldstart
		faas.ColdStart = &valCopy.Val
	}
	if from.Execution.IsSet() {
		faas.Execution = from.Execution.Val
	}
	if from.Trigger.Type.IsSet() {
		faas.TriggerType = from.Trigger.Type.Val
	}
	if from.Trigger.RequestID.IsSet() {
		faas.TriggerRequestId = from.Trigger.RequestID.Val
	}
	if from.Name.IsSet() {
		faas.Name = from.Name.Val
	}
	if from.Version.IsSet() {
		faas.Version = from.Version.Val
	}
}

func mapToDroppedSpansModel(from []transactionDroppedSpanStats, tx *modelpb.Transaction) {
	for _, f := range from {
		if f.IsSet() {
			var to modelpb.DroppedSpanStats
			if f.DestinationServiceResource.IsSet() {
				to.DestinationServiceResource = f.DestinationServiceResource.Val
			}
			if f.Outcome.IsSet() {
				to.Outcome = f.Outcome.Val
			}
			if f.Duration.IsSet() {
				to.Duration = &modelpb.AggregatedDuration{}
				to.Duration.Count = int64(f.Duration.Count.Val)
				sum := f.Duration.Sum
				if sum.IsSet() {
					to.Duration.Sum = durationpb.New(time.Duration(sum.Us.Val) * time.Microsecond)
				}
			}
			if f.ServiceTargetType.IsSet() {
				to.ServiceTargetType = f.ServiceTargetType.Val
			}
			if f.ServiceTargetName.IsSet() {
				to.ServiceTargetName = f.ServiceTargetName.Val
			}

			tx.DroppedSpansStats = append(tx.DroppedSpansStats, &to)
		}
	}
}

func mapToCloudModel(from contextCloud, cloud *modelpb.Cloud) {
	cloudOrigin := &modelpb.CloudOrigin{}
	if from.Origin.Account.ID.IsSet() {
		cloudOrigin.AccountId = from.Origin.Account.ID.Val
	}
	if from.Origin.Provider.IsSet() {
		cloudOrigin.Provider = from.Origin.Provider.Val
	}
	if from.Origin.Region.IsSet() {
		cloudOrigin.Region = from.Origin.Region.Val
	}
	if from.Origin.Service.Name.IsSet() {
		cloudOrigin.ServiceName = from.Origin.Service.Name.Val
	}
	cloud.Origin = cloudOrigin
}

func mapToClientModel(from contextRequest, source **modelpb.Source, client **modelpb.Client) {
	// http.Request.Headers and http.Request.Socket are only set for backend events.
	if _, err := netip.ParseAddr((*source).GetIp()); err != nil {
		ip, port := netutil.SplitAddrPort(from.Socket.RemoteAddress.Val)
		if ip.IsValid() {
			*source = populateNil(*source)
			(*source).Ip, (*source).Port = ip.String(), uint32(port)
		}
	}
	if _, err := netip.ParseAddr((*client).GetIp()); err != nil {
		if (*source).GetIp() != "" {
			*client = populateNil(*client)
			(*client).Ip = (*source).Ip
		}
		if ip, port := netutil.ClientAddrFromHeaders(from.Headers.Val); ip.IsValid() {
			if (*source).GetIp() != "" {
				(*source).Nat = &modelpb.NAT{Ip: (*source).Ip}
			}
			*client = populateNil(*client)
			(*client).Ip, (*client).Port = ip.String(), uint32(port)
			*source = populateNil(*source)
			(*source).Ip, (*source).Port = ip.String(), uint32(port)
		}
	}
}

func mapToErrorModel(from *errorEvent, event *modelpb.APMEvent) {
	out := &modelpb.Error{}
	event.Error = out
	event.Processor = modelpb.ErrorProcessor()

	// overwrite metadata with event specific information
	mapToServiceModel(from.Context.Service, &event.Service)
	mapToAgentModel(from.Context.Service.Agent, &event.Agent)
	overwriteUserInMetadataModel(from.Context.User, event)
	mapToUserAgentModel(from.Context.Request.Headers, &event.UserAgent)
	mapToClientModel(from.Context.Request, &event.Source, &event.Client)

	// map errorEvent specific data

	if from.Context.IsSet() {
		if len(from.Context.Tags) > 0 {
			modeldecoderutil.MergeLabels(from.Context.Tags, event)
		}
		if from.Context.Request.IsSet() {
			event.Http = populateNil(event.Http)
			event.Http.Request = &modelpb.HTTPRequest{}
			mapToRequestModel(from.Context.Request, event.Http.Request)
			if from.Context.Request.HTTPVersion.IsSet() {
				event.Http.Version = from.Context.Request.HTTPVersion.Val
			}
		}
		if from.Context.Response.IsSet() {
			event.Http = populateNil(event.Http)
			event.Http.Response = &modelpb.HTTPResponse{}
			mapToResponseModel(from.Context.Response, event.Http.Response)
		}
		if from.Context.Request.URL.IsSet() {
			event.Url = populateNil(event.Url)
			mapToRequestURLModel(from.Context.Request.URL, event.Url)
		}
		if from.Context.Page.IsSet() {
			if from.Context.Page.URL.IsSet() && !from.Context.Request.URL.IsSet() {
				event.Url = modelpb.ParseURL(from.Context.Page.URL.Val, "", "")
			}
			if from.Context.Page.Referer.IsSet() {
				event.Http = populateNil(event.Http)
				event.Http.Request = populateNil(event.Http.Request)
				if event.Http.Request.Referrer == "" {
					event.Http.Request.Referrer = from.Context.Page.Referer.Val
				}
			}
		}
		if len(from.Context.Custom) > 0 {
			out.Custom = modeldecoderutil.ToStruct(from.Context.Custom)
		}
	}
	if from.Culprit.IsSet() {
		out.Culprit = from.Culprit.Val
	}
	if from.Exception.IsSet() {
		out.Exception = &modelpb.Exception{}
		mapToExceptionModel(from.Exception, out.Exception)
	}
	if from.ID.IsSet() {
		out.Id = from.ID.Val
	}
	if from.Log.IsSet() {
		log := modelpb.ErrorLog{}
		if from.Log.Level.IsSet() {
			log.Level = from.Log.Level.Val
		}
		if from.Log.LoggerName.IsSet() {
			log.LoggerName = from.Log.LoggerName.Val
		}
		if from.Log.Message.IsSet() {
			log.Message = from.Log.Message.Val
		}
		if from.Log.ParamMessage.IsSet() {
			log.ParamMessage = from.Log.ParamMessage.Val
		}
		if len(from.Log.Stacktrace) > 0 {
			log.Stacktrace = make([]*modelpb.StacktraceFrame, len(from.Log.Stacktrace))
			mapToStracktraceModel(from.Log.Stacktrace, log.Stacktrace)
		}
		out.Log = &log
	}
	if from.ParentID.IsSet() {
		event.ParentId = from.ParentID.Val
	}
	if !from.Timestamp.Val.IsZero() {
		event.Timestamp = timestamppb.New(from.Timestamp.Val)
	}
	if from.TraceID.IsSet() {
		event.Trace = &modelpb.Trace{
			Id: from.TraceID.Val,
		}
	}
	if from.TransactionID.IsSet() || from.Transaction.IsSet() {
		event.Transaction = &modelpb.Transaction{}
	}
	if from.TransactionID.IsSet() {
		event.Transaction.Id = from.TransactionID.Val
		event.Span = &modelpb.Span{
			Id: from.TransactionID.Val,
		}
	}
	if from.Transaction.IsSet() {
		if from.Transaction.Sampled.IsSet() {
			event.Transaction.Sampled = from.Transaction.Sampled.Val
		}
		if from.Transaction.Name.IsSet() {
			event.Transaction.Name = from.Transaction.Name.Val
		}
		if from.Transaction.Type.IsSet() {
			event.Transaction.Type = from.Transaction.Type.Val
		}
	}
}

func mapToExceptionModel(from errorException, out *modelpb.Exception) {
	if len(from.Attributes) > 0 {
		out.Attributes = modeldecoderutil.ToStruct(from.Attributes)
	}
	if from.Code.IsSet() {
		out.Code = modeldecoderutil.ExceptionCodeString(from.Code.Val)
	}
	if len(from.Cause) > 0 {
		out.Cause = make([]*modelpb.Exception, len(from.Cause))
		for i := 0; i < len(from.Cause); i++ {
			if from.Cause[i].IsSet() {
				var ex modelpb.Exception
				mapToExceptionModel(from.Cause[i], &ex)
				out.Cause[i] = &ex
			}
		}
	}
	if from.Handled.IsSet() {
		out.Handled = &from.Handled.Val
	}
	if from.Message.IsSet() {
		out.Message = from.Message.Val
	}
	if from.Module.IsSet() {
		out.Module = from.Module.Val
	}
	if len(from.Stacktrace) > 0 {
		out.Stacktrace = make([]*modelpb.StacktraceFrame, len(from.Stacktrace))
		mapToStracktraceModel(from.Stacktrace, out.Stacktrace)
	}
	if from.Type.IsSet() {
		out.Type = from.Type.Val
	}
}

func mapToMetadataModel(from *metadata, out *modelpb.APMEvent) {
	// Cloud
	if from.Cloud.IsSet() {
		out.Cloud = populateNil(out.Cloud)
		if from.Cloud.Account.ID.IsSet() {
			out.Cloud.AccountId = from.Cloud.Account.ID.Val
		}
		if from.Cloud.Account.Name.IsSet() {
			out.Cloud.AccountName = from.Cloud.Account.Name.Val
		}
		if from.Cloud.AvailabilityZone.IsSet() {
			out.Cloud.AvailabilityZone = from.Cloud.AvailabilityZone.Val
		}
		if from.Cloud.Instance.ID.IsSet() {
			out.Cloud.InstanceId = from.Cloud.Instance.ID.Val
		}
		if from.Cloud.Instance.Name.IsSet() {
			out.Cloud.InstanceName = from.Cloud.Instance.Name.Val
		}
		if from.Cloud.Machine.Type.IsSet() {
			out.Cloud.MachineType = from.Cloud.Machine.Type.Val
		}
		if from.Cloud.Project.ID.IsSet() {
			out.Cloud.ProjectId = from.Cloud.Project.ID.Val
		}
		if from.Cloud.Project.Name.IsSet() {
			out.Cloud.ProjectName = from.Cloud.Project.Name.Val
		}
		if from.Cloud.Provider.IsSet() {
			out.Cloud.Provider = from.Cloud.Provider.Val
		}
		if from.Cloud.Region.IsSet() {
			out.Cloud.Region = from.Cloud.Region.Val
		}
		if from.Cloud.Service.Name.IsSet() {
			out.Cloud.ServiceName = from.Cloud.Service.Name.Val
		}
	}

	// Labels
	if len(from.Labels) > 0 {
		modeldecoderutil.GlobalLabelsFrom(from.Labels, out)
	}

	// Process
	if len(from.Process.Argv) > 0 {
		out.Process = populateNil(out.Process)
		out.Process.Argv = append(out.Process.Argv[:0], from.Process.Argv...)
	}
	if from.Process.Pid.IsSet() {
		out.Process = populateNil(out.Process)
		out.Process.Pid = uint32(from.Process.Pid.Val)
	}
	if from.Process.Ppid.IsSet() {
		var pid = uint32(from.Process.Ppid.Val)
		out.Process = populateNil(out.Process)
		out.Process.Ppid = pid
	}
	if from.Process.Title.IsSet() {
		out.Process = populateNil(out.Process)
		out.Process.Title = from.Process.Title.Val
	}

	// Service
	if from.Service.Agent.IsSet() {
		out.Agent = populateNil(out.Agent)
		if from.Service.Agent.ActivationMethod.IsSet() {
			out.Agent.ActivationMethod = from.Service.Agent.ActivationMethod.Val
		}
		if from.Service.Agent.EphemeralID.IsSet() {
			out.Agent.EphemeralId = from.Service.Agent.EphemeralID.Val
		}
		if from.Service.Agent.Name.IsSet() {
			out.Agent.Name = from.Service.Agent.Name.Val
		}
		if from.Service.Agent.Version.IsSet() {
			out.Agent.Version = from.Service.Agent.Version.Val
		}
	}
	if from.Service.IsSet() {
		out.Service = populateNil(out.Service)
		if from.Service.Environment.IsSet() {
			out.Service.Environment = from.Service.Environment.Val
		}
		if from.Service.Framework.IsSet() {
			out.Service.Framework = populateNil(out.Service.Framework)
			if from.Service.Framework.Name.IsSet() {
				out.Service.Framework.Name = from.Service.Framework.Name.Val
			}
			if from.Service.Framework.Version.IsSet() {
				out.Service.Framework.Version = from.Service.Framework.Version.Val
			}
		}
		if from.Service.Language.IsSet() {
			out.Service.Language = populateNil(out.Service.Language)
			if from.Service.Language.Name.IsSet() {
				out.Service.Language.Name = from.Service.Language.Name.Val
			}
			if from.Service.Language.Version.IsSet() {
				out.Service.Language.Version = from.Service.Language.Version.Val
			}
		}
		if from.Service.Name.IsSet() {
			out.Service.Name = from.Service.Name.Val
		}
		if from.Service.Node.Name.IsSet() {
			out.Service.Node = populateNil(out.Service.Node)
			out.Service.Node.Name = from.Service.Node.Name.Val
		}
		if from.Service.Runtime.IsSet() {
			out.Service.Runtime = populateNil(out.Service.Runtime)
			if from.Service.Runtime.Name.IsSet() {
				out.Service.Runtime.Name = from.Service.Runtime.Name.Val
			}
			if from.Service.Runtime.Version.IsSet() {
				out.Service.Runtime.Version = from.Service.Runtime.Version.Val
			}
		}
		if from.Service.Version.IsSet() {
			out.Service.Version = from.Service.Version.Val
		}
	}

	// System
	if from.System.Architecture.IsSet() {
		out.Host = populateNil(out.Host)
		out.Host.Architecture = from.System.Architecture.Val
	}
	if from.System.ConfiguredHostname.IsSet() {
		out.Host = populateNil(out.Host)
		out.Host.Name = from.System.ConfiguredHostname.Val
	}
	if from.System.Container.ID.IsSet() {
		out.Container = populateNil(out.Container)
		out.Container.Id = from.System.Container.ID.Val
	}
	if from.System.DetectedHostname.IsSet() {
		out.Host = populateNil(out.Host)
		out.Host.Hostname = from.System.DetectedHostname.Val
	}
	if !from.System.ConfiguredHostname.IsSet() && !from.System.DetectedHostname.IsSet() &&
		from.System.DeprecatedHostname.IsSet() {
		out.Host = populateNil(out.Host)
		out.Host.Hostname = from.System.DeprecatedHostname.Val
	}
	if from.System.Kubernetes.Namespace.IsSet() {
		out.Kubernetes = populateNil(out.Kubernetes)
		out.Kubernetes.Namespace = from.System.Kubernetes.Namespace.Val
	}
	if from.System.Kubernetes.Node.Name.IsSet() {
		out.Kubernetes = populateNil(out.Kubernetes)
		out.Kubernetes.NodeName = from.System.Kubernetes.Node.Name.Val
	}
	if from.System.Kubernetes.Pod.Name.IsSet() {
		out.Kubernetes = populateNil(out.Kubernetes)
		out.Kubernetes.PodName = from.System.Kubernetes.Pod.Name.Val
	}
	if from.System.Kubernetes.Pod.UID.IsSet() {
		out.Kubernetes = populateNil(out.Kubernetes)
		out.Kubernetes.PodUid = from.System.Kubernetes.Pod.UID.Val
	}
	if from.System.Platform.IsSet() {
		out.Host = populateNil(out.Host)
		out.Host.Os = populateNil(out.Host.Os)
		out.Host.Os.Platform = from.System.Platform.Val
	}

	// User
	if from.User.Domain.IsSet() {
		out.User = populateNil(out.User)
		out.User.Domain = fmt.Sprint(from.User.Domain.Val)
	}
	if from.User.ID.IsSet() {
		out.User = populateNil(out.User)
		out.User.Id = fmt.Sprint(from.User.ID.Val)
	}
	if from.User.Email.IsSet() {
		out.User = populateNil(out.User)
		out.User.Email = from.User.Email.Val
	}
	if from.User.Name.IsSet() {
		out.User = populateNil(out.User)
		out.User.Name = from.User.Name.Val
	}

	// Network
	if from.Network.Connection.Type.IsSet() {
		out.Network = populateNil(out.Network)
		out.Network.Connection = populateNil(out.Network.Connection)
		out.Network.Connection.Type = from.Network.Connection.Type.Val
	}
}

func mapToMetricsetModel(from *metricset, event *modelpb.APMEvent) bool {
	event.Metricset = &modelpb.Metricset{Name: "app"}
	event.Processor = modelpb.MetricsetProcessor()

	if !from.Timestamp.Val.IsZero() {
		event.Timestamp = timestamppb.New(from.Timestamp.Val)
	}

	if len(from.Samples) > 0 {
		samples := make([]*modelpb.MetricsetSample, 0, len(from.Samples))
		for name, sample := range from.Samples {
			var counts []int64
			var values []float64
			var histogram *modelpb.Histogram
			if n := len(sample.Values); n > 0 {
				values = make([]float64, n)
				copy(values, sample.Values)
			}
			if n := len(sample.Counts); n > 0 {
				counts = make([]int64, n)
				copy(counts, sample.Counts)
			}
			if len(counts) != 0 || len(values) != 0 {
				histogram = &modelpb.Histogram{
					Values: values,
					Counts: counts,
				}
			}

			samples = append(samples, &modelpb.MetricsetSample{
				Type:      metricTypeText[sample.Type.Val],
				Name:      name,
				Unit:      sample.Unit.Val,
				Value:     sample.Value.Val,
				Histogram: histogram,
			})
		}
		event.Metricset.Samples = samples
	}

	if len(from.Tags) > 0 {
		modeldecoderutil.MergeLabels(from.Tags, event)
	}

	if from.Span.IsSet() {
		event.Span = &modelpb.Span{}
		if from.Span.Subtype.IsSet() {
			event.Span.Subtype = from.Span.Subtype.Val
		}
		if from.Span.Type.IsSet() {
			event.Span.Type = from.Span.Type.Val
		}
	}

	ok := true
	if from.Transaction.IsSet() {
		event.Transaction = &modelpb.Transaction{}
		if from.Transaction.Name.IsSet() {
			event.Transaction.Name = from.Transaction.Name.Val
		}
		if from.Transaction.Type.IsSet() {
			event.Transaction.Type = from.Transaction.Type.Val
		}
		// Transaction fields specified: this is an APM-internal metricset.
		// If there are no known metric samples, we return false so the
		// metricset is not added to the batch.
		ok = modeldecoderutil.SetInternalMetrics(event)
	}

	if from.Service.Name.IsSet() {
		event.Service = populateNil(event.Service)
		event.Service.Name = from.Service.Name.Val
		event.Service.Version = from.Service.Version.Val
	}

	if from.FAAS.IsSet() {
		event.Faas = populateNil(event.Faas)
		mapToFAASModel(from.FAAS, event.Faas)
	}

	return ok
}

func mapToRequestModel(from contextRequest, out *modelpb.HTTPRequest) {
	if from.Method.IsSet() {
		out.Method = from.Method.Val
	}
	if len(from.Env) > 0 {
		out.Env = modeldecoderutil.ToStruct(from.Env)
	}
	if from.Body.IsSet() {
		out.Body = modeldecoderutil.ToValue(modeldecoderutil.NormalizeHTTPRequestBody(from.Body.Val))
	}
	if len(from.Cookies) > 0 {
		out.Cookies = modeldecoderutil.ToStruct(from.Cookies)
	}
	if from.Headers.IsSet() {
		out.Headers = modeldecoderutil.HTTPHeadersToModelpb(from.Headers.Val)
	}
}

func mapToRequestURLModel(from contextRequestURL, out *modelpb.URL) {
	if from.Raw.IsSet() {
		out.Original = from.Raw.Val
	}
	if from.Full.IsSet() {
		out.Full = from.Full.Val
	}
	if from.Hostname.IsSet() {
		out.Domain = from.Hostname.Val
	}
	if from.Path.IsSet() {
		out.Path = from.Path.Val
	}
	if from.Search.IsSet() {
		out.Query = from.Search.Val
	}
	if from.Hash.IsSet() {
		out.Fragment = from.Hash.Val
	}
	if from.Protocol.IsSet() {
		out.Scheme = strings.TrimSuffix(from.Protocol.Val, ":")
	}
	if from.Port.IsSet() {
		// should never result in an error, type is checked when decoding
		port, err := strconv.Atoi(fmt.Sprint(from.Port.Val))
		if err == nil {
			out.Port = uint32(port)
		}
	}
}

func mapToResponseModel(from contextResponse, out *modelpb.HTTPResponse) {
	if from.Finished.IsSet() {
		val := from.Finished.Val
		out.Finished = &val
	}
	if from.Headers.IsSet() {
		out.Headers = modeldecoderutil.HTTPHeadersToModelpb(from.Headers.Val)
	}
	if from.HeadersSent.IsSet() {
		val := from.HeadersSent.Val
		out.HeadersSent = &val
	}
	if from.StatusCode.IsSet() {
		out.StatusCode = int32(from.StatusCode.Val)
	}
	if from.TransferSize.IsSet() {
		val := int64(from.TransferSize.Val)
		out.TransferSize = &val
	}
	if from.EncodedBodySize.IsSet() {
		val := int64(from.EncodedBodySize.Val)
		out.EncodedBodySize = &val
	}
	if from.DecodedBodySize.IsSet() {
		val := int64(from.DecodedBodySize.Val)
		out.DecodedBodySize = &val
	}
}

func mapToServiceModel(from contextService, outPtr **modelpb.Service) {
	var out *modelpb.Service
	*outPtr = populateNil(*outPtr)
	out = *outPtr
	if from.Environment.IsSet() {
		out.Environment = from.Environment.Val
	}
	if from.Framework.IsSet() {
		out.Framework = populateNil(out.Framework)
		if from.Framework.Name.IsSet() {
			out.Framework.Name = from.Framework.Name.Val
		}
		if from.Framework.Version.IsSet() {
			out.Framework.Version = from.Framework.Version.Val
		}
	}
	if from.Language.IsSet() {
		out.Language = populateNil(out.Language)
		if from.Language.Name.IsSet() {
			out.Language.Name = from.Language.Name.Val
		}
		if from.Language.Version.IsSet() {
			out.Language.Version = from.Language.Version.Val
		}
	}
	if from.Name.IsSet() {
		out.Name = from.Name.Val
		out.Version = from.Version.Val
	}
	if from.Node.Name.IsSet() {
		out.Node = populateNil(out.Node)
		out.Node.Name = from.Node.Name.Val
	}
	if from.Runtime.IsSet() {
		out.Runtime = populateNil(out.Runtime)
		if from.Runtime.Name.IsSet() {
			out.Runtime.Name = from.Runtime.Name.Val
		}
		if from.Runtime.Version.IsSet() {
			out.Runtime.Version = from.Runtime.Version.Val
		}
	}
	if from.Origin.IsSet() {
		outOrigin := &modelpb.ServiceOrigin{}
		if from.Origin.ID.IsSet() {
			outOrigin.Id = from.Origin.ID.Val
		}
		if from.Origin.Name.IsSet() {
			outOrigin.Name = from.Origin.Name.Val
		}
		if from.Origin.Version.IsSet() {
			outOrigin.Version = from.Origin.Version.Val
		}
		out.Origin = outOrigin
	}
	if from.Target.IsSet() {
		outTarget := &modelpb.ServiceTarget{}
		if from.Target.Name.IsSet() {
			outTarget.Name = from.Target.Name.Val
		}
		if from.Target.Type.IsSet() {
			outTarget.Type = from.Target.Type.Val
		}
		out.Target = outTarget
	}
}

func mapToAgentModel(from contextServiceAgent, out **modelpb.Agent) {
	*out = populateNil(*out)
	if from.Name.IsSet() {
		(*out).Name = from.Name.Val
	}
	if from.Version.IsSet() {
		(*out).Version = from.Version.Val
	}
	if from.EphemeralID.IsSet() {
		(*out).EphemeralId = from.EphemeralID.Val
	}
}

func mapToSpanModel(from *span, event *modelpb.APMEvent) {
	out := &modelpb.Span{}
	event.Span = out
	event.Processor = modelpb.SpanProcessor()

	// map span specific data
	if !from.Action.IsSet() && !from.Subtype.IsSet() {
		sep := "."
		before, after, ok := strings.Cut(from.Type.Val, sep)
		out.Type = before
		if ok {
			out.Subtype, out.Action, _ = strings.Cut(after, sep)
		}
	} else {
		if from.Action.IsSet() {
			out.Action = from.Action.Val
		}
		if from.Subtype.IsSet() {
			out.Subtype = from.Subtype.Val
		}
		if from.Type.IsSet() {
			out.Type = from.Type.Val
		}
	}
	if from.Composite.IsSet() {
		composite := modelpb.Composite{}
		if from.Composite.Count.IsSet() {
			composite.Count = uint32(from.Composite.Count.Val)
		}
		if from.Composite.Sum.IsSet() {
			composite.Sum = from.Composite.Sum.Val
		}
		if strategy, ok := compressionStrategyText[from.Composite.CompressionStrategy.Val]; ok {
			composite.CompressionStrategy = strategy
		}
		out.Composite = &composite
	}
	if len(from.ChildIDs) > 0 {
		event.ChildIds = make([]string, len(from.ChildIDs))
		copy(event.ChildIds, from.ChildIDs)
	}
	if from.Context.Database.IsSet() {
		db := modelpb.DB{}
		if from.Context.Database.Instance.IsSet() {
			db.Instance = from.Context.Database.Instance.Val
		}
		if from.Context.Database.Link.IsSet() {
			db.Link = from.Context.Database.Link.Val
		}
		if from.Context.Database.RowsAffected.IsSet() {
			val := uint32(from.Context.Database.RowsAffected.Val)
			db.RowsAffected = &val
		}
		if from.Context.Database.Statement.IsSet() {
			db.Statement = from.Context.Database.Statement.Val
		}
		if from.Context.Database.Type.IsSet() {
			db.Type = from.Context.Database.Type.Val
		}
		if from.Context.Database.User.IsSet() {
			db.UserName = from.Context.Database.User.Val
		}
		out.Db = &db
	}
	if from.Context.Destination.Address.IsSet() || from.Context.Destination.Port.IsSet() {
		if from.Context.Destination.Address.IsSet() {
			event.Destination = populateNil(event.Destination)
			event.Destination.Address = from.Context.Destination.Address.Val
		}
		if from.Context.Destination.Port.IsSet() {
			event.Destination = populateNil(event.Destination)
			event.Destination.Port = uint32(from.Context.Destination.Port.Val)
		}
	}
	if from.Context.Destination.Service.IsSet() {
		service := modelpb.DestinationService{}
		if from.Context.Destination.Service.Name.IsSet() {
			service.Name = from.Context.Destination.Service.Name.Val
		}
		if from.Context.Destination.Service.Resource.IsSet() {
			service.Resource = from.Context.Destination.Service.Resource.Val
		}
		if from.Context.Destination.Service.Type.IsSet() {
			service.Type = from.Context.Destination.Service.Type.Val
		}
		out.DestinationService = &service
	}
	if from.Context.HTTP.IsSet() {
		if from.Context.HTTP.Method.IsSet() {
			event.Http = populateNil(event.Http)
			event.Http.Request = populateNil(event.Http.Request)
			event.Http.Request.Method = from.Context.HTTP.Method.Val
		}
		if from.Context.HTTP.Request.ID.IsSet() {
			event.Http = populateNil(event.Http)
			event.Http.Request = populateNil(event.Http.Request)
			event.Http.Request.Id = from.Context.HTTP.Request.ID.Val
		}
		if from.Context.HTTP.Response.IsSet() {
			event.Http = populateNil(event.Http)
			response := modelpb.HTTPResponse{}
			if from.Context.HTTP.Response.DecodedBodySize.IsSet() {
				val := int64(from.Context.HTTP.Response.DecodedBodySize.Val)
				response.DecodedBodySize = &val
			}
			if from.Context.HTTP.Response.EncodedBodySize.IsSet() {
				val := int64(from.Context.HTTP.Response.EncodedBodySize.Val)
				response.EncodedBodySize = &val
			}
			if from.Context.HTTP.Response.Headers.IsSet() {
				response.Headers = modeldecoderutil.HTTPHeadersToModelpb(from.Context.HTTP.Response.Headers.Val)
			}
			if from.Context.HTTP.Response.StatusCode.IsSet() {
				response.StatusCode = int32(from.Context.HTTP.Response.StatusCode.Val)
			}
			if from.Context.HTTP.Response.TransferSize.IsSet() {
				val := int64(from.Context.HTTP.Response.TransferSize.Val)
				response.TransferSize = &val
			}
			event.Http.Response = &response
		}
		if from.Context.HTTP.StatusCode.IsSet() {
			event.Http = populateNil(event.Http)
			event.Http.Response = populateNil(event.Http.Response)
			event.Http.Response.StatusCode = int32(from.Context.HTTP.StatusCode.Val)
		}
		if from.Context.HTTP.URL.IsSet() {
			event.Url = populateNil(event.Url)
			event.Url.Original = from.Context.HTTP.URL.Val
		}
	}
	if from.Context.Message.IsSet() {
		message := modelpb.Message{}
		if from.Context.Message.Body.IsSet() {
			message.Body = from.Context.Message.Body.Val
		}
		if from.Context.Message.Headers.IsSet() {
			message.Headers = modeldecoderutil.HTTPHeadersToModelpb(from.Context.Message.Headers.Val)
		}
		if from.Context.Message.Age.Milliseconds.IsSet() {
			val := int64(from.Context.Message.Age.Milliseconds.Val)
			message.AgeMillis = &val
		}
		if from.Context.Message.Queue.Name.IsSet() {
			message.QueueName = from.Context.Message.Queue.Name.Val
		}
		if from.Context.Message.RoutingKey.IsSet() {
			message.RoutingKey = from.Context.Message.RoutingKey.Val
		}
		out.Message = &message
	}
	if from.Context.Service.IsSet() {
		mapToServiceModel(from.Context.Service, &event.Service)
		mapToAgentModel(from.Context.Service.Agent, &event.Agent)
	}
	if !from.Context.Service.Target.Type.IsSet() && from.Context.Destination.Service.Resource.IsSet() {
		event.Service = populateNil(event.Service)
		outTarget := targetFromDestinationResource(from.Context.Destination.Service.Resource.Val)
		event.Service.Target = &outTarget
	}
	if len(from.Context.Tags) > 0 {
		modeldecoderutil.MergeLabels(from.Context.Tags, event)
	}
	if from.Duration.IsSet() {
		event.Event = populateNil(event.Event)
		duration := time.Duration(from.Duration.Val * float64(time.Millisecond))
		event.Event.Duration = durationpb.New(duration)
	}
	if from.ID.IsSet() {
		out.Id = from.ID.Val
	}
	if from.Name.IsSet() {
		out.Name = from.Name.Val
	}
	event.Event = populateNil(event.Event)
	if from.Outcome.IsSet() {
		event.Event.Outcome = from.Outcome.Val
	} else {
		if from.Context.HTTP.StatusCode.IsSet() {
			statusCode := from.Context.HTTP.StatusCode.Val
			if statusCode >= http.StatusBadRequest {
				event.Event.Outcome = "failure"
			} else {
				event.Event.Outcome = "success"
			}
		} else {
			event.Event.Outcome = "unknown"
		}
	}
	if from.ParentID.IsSet() {
		event.ParentId = from.ParentID.Val
	}
	if from.SampleRate.IsSet() {
		if from.SampleRate.Val > 0 {
			out.RepresentativeCount = 1 / from.SampleRate.Val
		}
	} else {
		// NOTE: we don't know the sample rate, so we need to make an assumption.
		//
		// Agents never send spans for non-sampled transactions, so assuming a
		// representative count of 1 (i.e. sampling rate of 100%) may be invalid; but
		// this is more useful than producing no metrics at all (i.e. assuming 0).
		out.RepresentativeCount = 1
	}
	if len(from.Stacktrace) > 0 {
		out.Stacktrace = make([]*modelpb.StacktraceFrame, len(from.Stacktrace))
		mapToStracktraceModel(from.Stacktrace, out.Stacktrace)
	}
	if from.Sync.IsSet() {
		val := from.Sync.Val
		out.Sync = &val
	}
	if !from.Timestamp.Val.IsZero() {
		event.Timestamp = timestamppb.New(from.Timestamp.Val)
	} else if from.Start.IsSet() {
		// event.Timestamp should have been initialized to the time the
		// payload was received; offset that by "start" milliseconds for
		// RUM.
		base := time.Time{}
		if event.Timestamp != nil {
			base = event.Timestamp.AsTime()
		}
		event.Timestamp = timestamppb.New(base.Add(
			time.Duration(float64(time.Millisecond) * from.Start.Val),
		))
	}
	if from.TraceID.IsSet() {
		event.Trace = &modelpb.Trace{
			Id: from.TraceID.Val,
		}
	}
	if from.TransactionID.IsSet() {
		event.Transaction = &modelpb.Transaction{Id: from.TransactionID.Val}
	}
	if from.OTel.IsSet() {
		mapOTelAttributesSpan(from.OTel, event)
	}
	if len(from.Links) > 0 {
		mapSpanLinks(from.Links, &out.Links)
	}
	if out.Type == "" {
		out.Type = "unknown"
	}
}

func mapToStracktraceModel(from []stacktraceFrame, out []*modelpb.StacktraceFrame) {
	for idx, eventFrame := range from {
		fr := modelpb.StacktraceFrame{}
		if eventFrame.AbsPath.IsSet() {
			fr.AbsPath = eventFrame.AbsPath.Val
		}
		if eventFrame.Classname.IsSet() {
			fr.Classname = eventFrame.Classname.Val
		}
		if eventFrame.ColumnNumber.IsSet() {
			val := uint32(eventFrame.ColumnNumber.Val)
			fr.Colno = &val
		}
		if eventFrame.ContextLine.IsSet() {
			fr.ContextLine = eventFrame.ContextLine.Val
		}
		if eventFrame.Filename.IsSet() {
			fr.Filename = eventFrame.Filename.Val
		}
		if eventFrame.Function.IsSet() {
			fr.Function = eventFrame.Function.Val
		}
		if eventFrame.LibraryFrame.IsSet() {
			val := eventFrame.LibraryFrame.Val
			fr.LibraryFrame = val
		}
		if eventFrame.LineNumber.IsSet() {
			val := uint32(eventFrame.LineNumber.Val)
			fr.Lineno = &val
		}
		if eventFrame.Module.IsSet() {
			fr.Module = eventFrame.Module.Val
		}
		if len(eventFrame.PostContext) > 0 {
			fr.PostContext = make([]string, len(eventFrame.PostContext))
			copy(fr.PostContext, eventFrame.PostContext)
		}
		if len(eventFrame.PreContext) > 0 {
			fr.PreContext = make([]string, len(eventFrame.PreContext))
			copy(fr.PreContext, eventFrame.PreContext)
		}
		if len(eventFrame.Vars) > 0 {
			fr.Vars = modeldecoderutil.ToStruct(eventFrame.Vars)
		}
		out[idx] = &fr
	}
}

func mapToTransactionModel(from *transaction, event *modelpb.APMEvent) {
	out := &modelpb.Transaction{}
	event.Processor = modelpb.TransactionProcessor()
	event.Transaction = out

	// overwrite metadata with event specific information
	mapToServiceModel(from.Context.Service, &event.Service)
	mapToAgentModel(from.Context.Service.Agent, &event.Agent)
	overwriteUserInMetadataModel(from.Context.User, event)
	mapToUserAgentModel(from.Context.Request.Headers, &event.UserAgent)
	mapToClientModel(from.Context.Request, &event.Source, &event.Client)
	if from.FAAS.IsSet() {
		event.Faas = populateNil(event.Faas)
		mapToFAASModel(from.FAAS, event.Faas)
	}
	if from.Context.Cloud.IsSet() {
		event.Cloud = populateNil(event.Cloud)
		mapToCloudModel(from.Context.Cloud, event.Cloud)
	}
	mapToDroppedSpansModel(from.DroppedSpanStats, event.Transaction)

	// map transaction specific data

	if from.Context.IsSet() {
		if len(from.Context.Custom) > 0 {
			out.Custom = modeldecoderutil.ToStruct(from.Context.Custom)
		}
		if len(from.Context.Tags) > 0 {
			modeldecoderutil.MergeLabels(from.Context.Tags, event)
		}
		if from.Context.Message.IsSet() {
			out.Message = &modelpb.Message{}
			if from.Context.Message.Age.IsSet() {
				val := int64(from.Context.Message.Age.Milliseconds.Val)
				out.Message.AgeMillis = &val
			}
			if from.Context.Message.Body.IsSet() {
				out.Message.Body = from.Context.Message.Body.Val
			}
			if from.Context.Message.Headers.IsSet() {
				out.Message.Headers = modeldecoderutil.HTTPHeadersToModelpb(from.Context.Message.Headers.Val)
			}
			if from.Context.Message.Queue.IsSet() && from.Context.Message.Queue.Name.IsSet() {
				out.Message.QueueName = from.Context.Message.Queue.Name.Val
			}
			if from.Context.Message.RoutingKey.IsSet() {
				out.Message.RoutingKey = from.Context.Message.RoutingKey.Val
			}
		}
		if from.Context.Request.IsSet() {
			event.Http = populateNil(event.Http)
			event.Http.Request = &modelpb.HTTPRequest{}
			mapToRequestModel(from.Context.Request, event.Http.Request)
			if from.Context.Request.HTTPVersion.IsSet() {
				event.Http.Version = from.Context.Request.HTTPVersion.Val
			}
		}
		if from.Context.Request.URL.IsSet() {
			event.Url = populateNil(event.Url)
			mapToRequestURLModel(from.Context.Request.URL, event.Url)
		}
		if from.Context.Response.IsSet() {
			event.Http = populateNil(event.Http)
			event.Http.Response = &modelpb.HTTPResponse{}
			mapToResponseModel(from.Context.Response, event.Http.Response)
		}
		if from.Context.Page.IsSet() {
			if from.Context.Page.URL.IsSet() && !from.Context.Request.URL.IsSet() {
				event.Url = modelpb.ParseURL(from.Context.Page.URL.Val, "", "")
			}
			if from.Context.Page.Referer.IsSet() {
				event.Http = populateNil(event.Http)
				event.Http.Request = populateNil(event.Http.Request)
				if event.Http.Request.Referrer == "" {
					event.Http.Request.Referrer = from.Context.Page.Referer.Val
				}
			}
		}
	}
	if from.Duration.IsSet() {
		duration := time.Duration(from.Duration.Val * float64(time.Millisecond))
		event.Event = populateNil(event.Event)
		event.Event.Duration = durationpb.New(duration)
	}
	if from.ID.IsSet() {
		out.Id = from.ID.Val
		event.Span = &modelpb.Span{
			Id: from.ID.Val,
		}
	}
	if from.Marks.IsSet() {
		out.Marks = make(map[string]*modelpb.TransactionMark, len(from.Marks.Events))
		for event, val := range from.Marks.Events {
			if len(val.Measurements) > 0 {
				out.Marks[event] = &modelpb.TransactionMark{
					Measurements: val.Measurements,
				}
			}
		}
	}
	if from.Name.IsSet() {
		out.Name = from.Name.Val
	}
	event.Event = populateNil(event.Event)
	if from.Outcome.IsSet() {
		event.Event.Outcome = from.Outcome.Val
	} else {
		if from.Context.Response.StatusCode.IsSet() {
			statusCode := from.Context.Response.StatusCode.Val
			if statusCode >= http.StatusInternalServerError {
				event.Event.Outcome = "failure"
			} else {
				event.Event.Outcome = "success"
			}
		} else {
			event.Event.Outcome = "unknown"
		}
	}
	if from.ParentID.IsSet() {
		event.ParentId = from.ParentID.Val
	}
	if from.Result.IsSet() {
		out.Result = from.Result.Val
	}
	sampled := true
	if from.Sampled.IsSet() {
		sampled = from.Sampled.Val
	}
	out.Sampled = sampled
	if from.SampleRate.IsSet() {
		if from.SampleRate.Val > 0 {
			out.RepresentativeCount = 1 / from.SampleRate.Val
		}
	} else {
		out.RepresentativeCount = 1
	}
	if from.Session.ID.IsSet() {
		event.Session = &modelpb.Session{
			Id:       from.Session.ID.Val,
			Sequence: int64(from.Session.Sequence.Val),
		}
	}
	if from.SpanCount.Dropped.IsSet() {
		out.SpanCount = populateNil(out.SpanCount)
		dropped := uint32(from.SpanCount.Dropped.Val)
		out.SpanCount.Dropped = &dropped
	}
	if from.SpanCount.Started.IsSet() {
		out.SpanCount = populateNil(out.SpanCount)
		started := uint32(from.SpanCount.Started.Val)
		out.SpanCount.Started = &started
	}
	if !from.Timestamp.Val.IsZero() {
		event.Timestamp = timestamppb.New(from.Timestamp.Val)
	}
	if from.TraceID.IsSet() {
		event.Trace = &modelpb.Trace{
			Id: from.TraceID.Val,
		}
	}
	if from.Type.IsSet() {
		out.Type = from.Type.Val
	}
	if from.UserExperience.IsSet() {
		out.UserExperience = &modelpb.UserExperience{
			CumulativeLayoutShift: -1,
			FirstInputDelay:       -1,
			TotalBlockingTime:     -1,
			LongTask:              &modelpb.LongtaskMetrics{Count: -1},
		}
		if from.UserExperience.CumulativeLayoutShift.IsSet() {
			out.UserExperience.CumulativeLayoutShift = from.UserExperience.CumulativeLayoutShift.Val
		}
		if from.UserExperience.FirstInputDelay.IsSet() {
			out.UserExperience.FirstInputDelay = from.UserExperience.FirstInputDelay.Val

		}
		if from.UserExperience.TotalBlockingTime.IsSet() {
			out.UserExperience.TotalBlockingTime = from.UserExperience.TotalBlockingTime.Val
		}
		if from.UserExperience.Longtask.IsSet() {
			out.UserExperience.LongTask = &modelpb.LongtaskMetrics{
				Count: int64(from.UserExperience.Longtask.Count.Val),
				Sum:   from.UserExperience.Longtask.Sum.Val,
				Max:   from.UserExperience.Longtask.Max.Val,
			}
		}
	}

	if from.OTel.IsSet() {
		event.Span = populateNil(event.Span)
		mapOTelAttributesTransaction(from.OTel, event)
	}

	if len(from.Links) > 0 {
		event.Span = populateNil(event.Span)
		mapSpanLinks(from.Links, &event.Span.Links)
	}
	if out.Type == "" {
		out.Type = "unknown"
	}
}

func mapToLogModel(from *log, event *modelpb.APMEvent) {
	event.Processor = modelpb.LogProcessor()

	if from.FAAS.IsSet() {
		event.Faas = populateNil(event.Faas)
		mapToFAASModel(from.FAAS, event.Faas)
	}
	if !from.Timestamp.Val.IsZero() {
		event.Timestamp = timestamppb.New(from.Timestamp.Val)
	}
	if from.TraceID.IsSet() {
		event.Trace = &modelpb.Trace{
			Id: from.TraceID.Val,
		}
	}
	if from.TransactionID.IsSet() {
		event.Transaction = &modelpb.Transaction{
			Id: from.TransactionID.Val,
		}
		event.Span = &modelpb.Span{
			Id: from.TransactionID.Val,
		}
	}
	if from.SpanID.IsSet() {
		event.Span = &modelpb.Span{
			Id: from.SpanID.Val,
		}
	}
	if from.Message.IsSet() {
		event.Message = from.Message.Val
	}
	if from.Level.IsSet() {
		event.Log = populateNil(event.Log)
		event.Log.Level = from.Level.Val
	}
	if from.Logger.IsSet() {
		event.Log = populateNil(event.Log)
		event.Log.Logger = from.Logger.Val
	}
	if from.OriginFunction.IsSet() {
		event.Log = populateNil(event.Log)
		event.Log.Origin = populateNil(event.Log.Origin)
		event.Log.Origin.FunctionName = from.OriginFunction.Val
	}
	if from.OriginFileLine.IsSet() {
		event.Log = populateNil(event.Log)
		event.Log.Origin = populateNil(event.Log.Origin)
		event.Log.Origin.File = populateNil(event.Log.Origin.File)
		event.Log.Origin.File.Line = int32(from.OriginFileLine.Val)
	}
	if from.OriginFileName.IsSet() {
		event.Log = populateNil(event.Log)
		event.Log.Origin = populateNil(event.Log.Origin)
		event.Log.Origin.File = populateNil(event.Log.Origin.File)
		event.Log.Origin.File.Name = from.OriginFileName.Val
	}
	if from.ErrorType.IsSet() ||
		from.ErrorMessage.IsSet() ||
		from.ErrorStacktrace.IsSet() {
		event.Error = &modelpb.Error{
			Message:    from.ErrorMessage.Val,
			Type:       from.ErrorType.Val,
			StackTrace: from.ErrorStacktrace.Val,
		}
	}
	if from.ServiceName.IsSet() {
		event.Service = populateNil(event.Service)
		event.Service.Name = from.ServiceName.Val
	}
	if from.ServiceVersion.IsSet() {
		event.Service = populateNil(event.Service)
		event.Service.Version = from.ServiceVersion.Val
	}
	if from.ServiceEnvironment.IsSet() {
		event.Service = populateNil(event.Service)
		event.Service.Environment = from.ServiceEnvironment.Val
	}
	if from.ServiceNodeName.IsSet() {
		event.Service = populateNil(event.Service)
		event.Service.Node = populateNil(event.Service.Node)
		event.Service.Node.Name = from.ServiceNodeName.Val
	}
	if from.ProcessThreadName.IsSet() {
		event.Process = populateNil(event.Process)
		event.Process.Thread = populateNil(event.Process.Thread)
		event.Process.Thread.Name = from.ProcessThreadName.Val
	}
	if from.EventDataset.IsSet() {
		event.Event = populateNil(event.Event)
		event.Event.Dataset = from.EventDataset.Val
	}
	if len(from.Labels) > 0 {
		modeldecoderutil.MergeLabels(from.Labels, event)
	}
}

func mapOTelAttributesTransaction(from otel, out *modelpb.APMEvent) {
	scope := pcommon.NewInstrumentationScope()
	m := otelAttributeMap(&from)
	if from.SpanKind.IsSet() {
		out.Span.Kind = from.SpanKind.Val
	}
	if out.Labels == nil {
		out.Labels = make(modelpb.Labels)
	}
	if out.NumericLabels == nil {
		out.NumericLabels = make(modelpb.NumericLabels)
	}
	// TODO: Does this work? Is there a way we can infer the status code,
	// potentially in the actual attributes map?
	spanStatus := ptrace.NewStatus()
	spanStatus.SetCode(ptrace.StatusCodeUnset)
	otlp.TranslateTransaction(m, spanStatus, scope, out)

	if out.Span.Kind == "" {
		switch out.Transaction.Type {
		case "messaging":
			out.Span.Kind = "CONSUMER"
		case "request":
			out.Span.Kind = "SERVER"
		default:
			out.Span.Kind = "INTERNAL"
		}
	}
}

const (
	spanKindInternal = "INTERNAL"
	spanKindServer   = "SERVER"
	spanKindClient   = "CLIENT"
	spanKindProducer = "PRODUCER"
	spanKindConsumer = "CONSUMER"
)

func mapOTelAttributesSpan(from otel, out *modelpb.APMEvent) {
	m := otelAttributeMap(&from)
	if out.Labels == nil {
		out.Labels = make(modelpb.Labels)
	}
	if out.NumericLabels == nil {
		out.NumericLabels = make(modelpb.NumericLabels)
	}
	var spanKind ptrace.SpanKind
	if from.SpanKind.IsSet() {
		switch from.SpanKind.Val {
		case spanKindInternal:
			spanKind = ptrace.SpanKindInternal
		case spanKindServer:
			spanKind = ptrace.SpanKindServer
		case spanKindClient:
			spanKind = ptrace.SpanKindClient
		case spanKindProducer:
			spanKind = ptrace.SpanKindProducer
		case spanKindConsumer:
			spanKind = ptrace.SpanKindConsumer
		default:
			spanKind = ptrace.SpanKindUnspecified
		}
		out.Span.Kind = from.SpanKind.Val
	}
	otlp.TranslateSpan(spanKind, m, out)

	if spanKind == ptrace.SpanKindUnspecified {
		switch out.Span.Type {
		case "db", "external", "storage":
			out.Span.Kind = "CLIENT"
		default:
			out.Span.Kind = "INTERNAL"
		}
	}
}

func mapToUserAgentModel(from nullable.HTTPHeader, out **modelpb.UserAgent) {
	// overwrite userAgent information if available
	if from.IsSet() {
		if h := from.Val.Values(textproto.CanonicalMIMEHeaderKey("User-Agent")); len(h) > 0 {
			*out = populateNil(*out)
			(*out).Original = strings.Join(h, ", ")
		}
	}
}

func overwriteUserInMetadataModel(from user, out *modelpb.APMEvent) {
	// overwrite User specific values if set
	// either populate all User fields or none to avoid mixing
	// different user data
	if !from.Domain.IsSet() && !from.ID.IsSet() && !from.Email.IsSet() && !from.Name.IsSet() {
		return
	}
	out.User = &modelpb.User{}
	if from.Domain.IsSet() {
		out.User.Domain = fmt.Sprint(from.Domain.Val)
	}
	if from.ID.IsSet() {
		out.User.Id = fmt.Sprint(from.ID.Val)
	}
	if from.Email.IsSet() {
		out.User.Email = from.Email.Val
	}
	if from.Name.IsSet() {
		out.User.Name = from.Name.Val
	}
}

func otelAttributeMap(o *otel) pcommon.Map {
	m := pcommon.NewMap()
	for k, v := range o.Attributes {
		if attr, ok := otelAttributeValue(k, v); ok {
			attr.CopyTo(m.PutEmpty(k))
		}
	}
	return m
}

func otelAttributeValue(k string, v interface{}) (pcommon.Value, bool) {
	// According to the spec, these are the allowed primitive types
	// Additionally, homogeneous arrays (single type) of primitive types are allowed
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/common/common.md#attributes
	switch v := v.(type) {
	case string:
		return pcommon.NewValueStr(v), true
	case bool:
		return pcommon.NewValueBool(v), true
	case json.Number:
		// Semantic conventions have specified types, and we rely on this
		// in processor/otel when mapping to our data model. For example,
		// `http.status_code` is expected to be an int.
		if !isOTelDoubleAttribute(k) {
			if v, err := v.Int64(); err == nil {
				return pcommon.NewValueInt(v), true
			}
		}
		if v, err := v.Float64(); err == nil {
			return pcommon.NewValueDouble(v), true
		}
	case []interface{}:
		array := pcommon.NewValueSlice()
		array.Slice().EnsureCapacity(len(v))
		for i := range v {
			if elem, ok := otelAttributeValue(k, v[i]); ok {
				elem.CopyTo(array.Slice().AppendEmpty())
			}
		}
		return array, true
	}
	return pcommon.Value{}, false
}

// isOTelDoubleAttribute indicates whether k is an OpenTelemetry semantic convention attribute
// known to have type "double". As this list grows over time, we should consider generating
// the mapping with OpenTelemetry's semconvgen build tool.
//
// For the canonical semantic convention definitions, see
// https://github.com/open-telemetry/opentelemetry-specification/tree/main/semantic_conventions/trace
func isOTelDoubleAttribute(k string) bool {
	switch k {
	case "aws.dynamodb.provisioned_read_capacity":
		return true
	case "aws.dynamodb.provisioned_write_capacity":
		return true
	}
	return false
}

func mapSpanLinks(from []spanLink, out *[]*modelpb.SpanLink) {
	*out = make([]*modelpb.SpanLink, len(from))
	for i, link := range from {
		(*out)[i] = &modelpb.SpanLink{
			SpanId:  link.SpanID.Val,
			TraceId: link.TraceID.Val,
		}
	}
}

func targetFromDestinationResource(res string) (target modelpb.ServiceTarget) {
	submatch := reForServiceTargetExpr.FindStringSubmatch(res)
	switch len(submatch) {
	case 3:
		target.Type = submatch[1]
		target.Name = submatch[2]
	default:
		target.Name = res
	}
	return
}

func populateNil[T any](a *T) *T {
	if a == nil {
		return new(T)
	}
	return a
}
