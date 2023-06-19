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
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"

	"github.com/elastic/apm-data/model/modelpb"
)

const (
	AgentNameJaeger = "Jaeger"
)

var (
	serviceNameInvalidRegexp = regexp.MustCompile("[^a-zA-Z0-9 _-]")
)

func translateResourceMetadata(resource pcommon.Resource, out *modelpb.APMEvent) {
	var exporterVersion string
	resource.Attributes().Range(func(k string, v pcommon.Value) bool {
		switch k {
		// service.*
		case semconv.AttributeServiceName:
			out.Service = populateNil(out.Service)
			out.Service.Name = cleanServiceName(v.Str())
		case semconv.AttributeServiceVersion:
			out.Service = populateNil(out.Service)
			out.Service.Version = truncate(v.Str())
		case semconv.AttributeServiceInstanceID:
			out.Service = populateNil(out.Service)
			out.Service.Node = populateNil(out.Service.Node)
			out.Service.Node.Name = truncate(v.Str())

		// deployment.*
		case semconv.AttributeDeploymentEnvironment:
			out.Service = populateNil(out.Service)
			out.Service.Environment = truncate(v.Str())

		// telemetry.sdk.*
		case semconv.AttributeTelemetrySDKName:
			out.Agent = populateNil(out.Agent)
			out.Agent.Name = truncate(v.Str())
		case semconv.AttributeTelemetrySDKVersion:
			out.Agent = populateNil(out.Agent)
			out.Agent.Version = truncate(v.Str())
		case semconv.AttributeTelemetrySDKLanguage:
			out.Service = populateNil(out.Service)
			out.Service.Language = populateNil(out.Service.Language)
			out.Service.Language.Name = truncate(v.Str())

		// cloud.*
		case semconv.AttributeCloudProvider:
			out.Cloud = populateNil(out.Cloud)
			out.Cloud.Provider = truncate(v.Str())
		case semconv.AttributeCloudAccountID:
			out.Cloud = populateNil(out.Cloud)
			out.Cloud.AccountId = truncate(v.Str())
		case semconv.AttributeCloudRegion:
			out.Cloud = populateNil(out.Cloud)
			out.Cloud.Region = truncate(v.Str())
		case semconv.AttributeCloudAvailabilityZone:
			out.Cloud = populateNil(out.Cloud)
			out.Cloud.AvailabilityZone = truncate(v.Str())
		case semconv.AttributeCloudPlatform:
			out.Cloud = populateNil(out.Cloud)
			out.Cloud.ServiceName = truncate(v.Str())

		// container.*
		case semconv.AttributeContainerName:
			out.Container = populateNil(out.Container)
			out.Container.Name = truncate(v.Str())
		case semconv.AttributeContainerID:
			out.Container = populateNil(out.Container)
			out.Container.Id = truncate(v.Str())
		case semconv.AttributeContainerImageName:
			out.Container = populateNil(out.Container)
			out.Container.ImageName = truncate(v.Str())
		case semconv.AttributeContainerImageTag:
			out.Container = populateNil(out.Container)
			out.Container.ImageTag = truncate(v.Str())
		case "container.runtime":
			out.Container = populateNil(out.Container)
			out.Container.Runtime = truncate(v.Str())

		// k8s.*
		case semconv.AttributeK8SNamespaceName:
			out.Kubernetes = populateNil(out.Kubernetes)
			out.Kubernetes.Namespace = truncate(v.Str())
		case semconv.AttributeK8SNodeName:
			out.Kubernetes = populateNil(out.Kubernetes)
			out.Kubernetes.NodeName = truncate(v.Str())
		case semconv.AttributeK8SPodName:
			out.Kubernetes = populateNil(out.Kubernetes)
			out.Kubernetes.PodName = truncate(v.Str())
		case semconv.AttributeK8SPodUID:
			out.Kubernetes = populateNil(out.Kubernetes)
			out.Kubernetes.PodUid = truncate(v.Str())

		// host.*
		case semconv.AttributeHostName:
			out.Host = populateNil(out.Host)
			out.Host.Hostname = truncate(v.Str())
		case semconv.AttributeHostID:
			out.Host = populateNil(out.Host)
			out.Host.Id = truncate(v.Str())
		case semconv.AttributeHostType:
			out.Host = populateNil(out.Host)
			out.Host.Type = truncate(v.Str())
		case "host.arch":
			out.Host = populateNil(out.Host)
			out.Host.Architecture = truncate(v.Str())

		// process.*
		case semconv.AttributeProcessPID:
			out.Process = populateNil(out.Process)
			out.Process.Pid = uint32(v.Int())
		case semconv.AttributeProcessCommandLine:
			out.Process = populateNil(out.Process)
			out.Process.CommandLine = truncate(v.Str())
		case semconv.AttributeProcessExecutablePath:
			out.Process = populateNil(out.Process)
			out.Process.Executable = truncate(v.Str())
		case "process.runtime.name":
			out.Service = populateNil(out.Service)
			out.Service.Runtime = populateNil(out.Service.Runtime)
			out.Service.Runtime.Name = truncate(v.Str())
		case "process.runtime.version":
			out.Service = populateNil(out.Service)
			out.Service.Runtime = populateNil(out.Service.Runtime)
			out.Service.Runtime.Version = truncate(v.Str())

		// os.*
		case semconv.AttributeOSType:
			out.Host = populateNil(out.Host)
			out.Host.Os = populateNil(out.Host.Os)
			out.Host.Os.Platform = strings.ToLower(truncate(v.Str()))
		case semconv.AttributeOSDescription:
			out.Host = populateNil(out.Host)
			out.Host.Os = populateNil(out.Host.Os)
			out.Host.Os.Full = truncate(v.Str())
		case semconv.AttributeOSName:
			out.Host = populateNil(out.Host)
			out.Host.Os = populateNil(out.Host.Os)
			out.Host.Os.Name = truncate(v.Str())
		case semconv.AttributeOSVersion:
			out.Host = populateNil(out.Host)
			out.Host.Os = populateNil(out.Host.Os)
			out.Host.Os.Version = truncate(v.Str())

		// device.*
		case semconv.AttributeDeviceID:
			out.Device = populateNil(out.Device)
			out.Device.Id = truncate(v.Str())
		case semconv.AttributeDeviceModelIdentifier:
			out.Device = populateNil(out.Device)
			out.Device.Model = populateNil(out.Device.Model)
			out.Device.Model.Identifier = truncate(v.Str())
		case semconv.AttributeDeviceModelName:
			out.Device = populateNil(out.Device)
			out.Device.Model = populateNil(out.Device.Model)
			out.Device.Model.Name = truncate(v.Str())
		case "device.manufacturer":
			out.Device = populateNil(out.Device)
			out.Device.Manufacturer = truncate(v.Str())

		// Legacy OpenCensus attributes.
		case "opencensus.exporterversion":
			exporterVersion = v.Str()

		// timestamp attribute to deal with time skew on mobile
		// devices. APM server should drop this field.
		case "telemetry.sdk.elastic_export_timestamp":
			// Do nothing.

		default:
			if out.Labels == nil {
				out.Labels = make(modelpb.Labels)
			}
			if out.NumericLabels == nil {
				out.NumericLabels = make(modelpb.NumericLabels)
			}
			setLabel(replaceDots(k), out, ifaceAttributeValue(v))
		}
		return true
	})

	// https://www.elastic.co/guide/en/ecs/current/ecs-os.html#field-os-type:
	//
	// "One of these following values should be used (lowercase): linux, macos, unix, windows.
	// If the OS you’re dealing with is not in the list, the field should not be populated."
	switch out.GetHost().GetOs().GetPlatform() {
	case "windows", "linux":
		out.Host.Os.Type = out.Host.Os.Platform
	case "darwin":
		out.Host.Os.Type = "macos"
	case "aix", "hpux", "solaris":
		out.Host.Os.Type = "unix"
	}

	switch out.GetHost().GetOs().GetName() {
	case "Android":
		out.Host.Os.Type = "android"
	case "iOS":
		out.Host.Os.Type = "ios"
	}

	if strings.HasPrefix(exporterVersion, "Jaeger") {
		// version is of format `Jaeger-<agentlanguage>-<version>`, e.g. `Jaeger-Go-2.20.0`
		const nVersionParts = 3
		versionParts := strings.SplitN(exporterVersion, "-", nVersionParts)
		if out.GetService().GetLanguage().GetName() == "" && len(versionParts) == nVersionParts {
			out.Service = populateNil(out.Service)
			out.Service.Language = populateNil(out.Service.Language)
			out.Service.Language.Name = versionParts[1]
		}
		out.Agent = populateNil(out.Agent)
		if v := versionParts[len(versionParts)-1]; v != "" {
			out.Agent.Version = v
		}
		out.Agent.Name = AgentNameJaeger

		// Translate known Jaeger labels.
		if clientUUID, ok := out.Labels["client-uuid"]; ok {
			out.Agent.EphemeralId = clientUUID.Value
			delete(out.Labels, "client-uuid")
		}
		if systemIP, ok := out.Labels["ip"]; ok {
			if ip, err := netip.ParseAddr(systemIP.Value); err == nil {
				out.Host.Ip = []string{ip.String()}
			}
			delete(out.Labels, "ip")
		}
	}

	if out.GetService().GetName() == "" {
		out.Service = populateNil(out.Service)
		// service.name is a required field.
		out.Service.Name = "unknown"
	}
	if out.GetAgent().GetName() == "" {
		out.Agent = populateNil(out.Agent)
		// agent.name is a required field.
		out.Agent.Name = "otlp"
	}
	if out.Agent.Version == "" {
		// agent.version is a required field.
		out.Agent.Version = "unknown"
	}
	if out.GetService().GetLanguage().GetName() != "" {
		out.Agent = populateNil(out.Agent)
		out.Agent.Name = fmt.Sprintf("%s/%s", out.Agent.Name, out.Service.Language.Name)
	} else {
		out.Service = populateNil(out.Service)
		out.Service.Language = populateNil(out.Service.Language)
		out.Service.Language.Name = "unknown"
	}

	// Set the decoded labels as "global" -- defined at the service level.
	for k, v := range out.Labels {
		v.Global = true
		out.Labels[k] = v
	}
	for k, v := range out.NumericLabels {
		v.Global = true
		out.NumericLabels[k] = v
	}
}

func cleanServiceName(name string) string {
	return serviceNameInvalidRegexp.ReplaceAllString(truncate(name), "_")
}

func ifaceAttributeValue(v pcommon.Value) interface{} {
	switch v.Type() {
	case pcommon.ValueTypeStr:
		return truncate(v.Str())
	case pcommon.ValueTypeBool:
		return strconv.FormatBool(v.Bool())
	case pcommon.ValueTypeInt:
		return float64(v.Int())
	case pcommon.ValueTypeDouble:
		return v.Double()
	case pcommon.ValueTypeSlice:
		return ifaceAttributeValueSlice(v.Slice())
	}
	return nil
}

func ifaceAttributeValueSlice(slice pcommon.Slice) []interface{} {
	values := make([]interface{}, slice.Len())
	for i := range values {
		values[i] = ifaceAttributeValue(slice.At(i))
	}
	return values
}

// initEventLabels initializes an event-specific labels from an event.
func initEventLabels(e *modelpb.APMEvent) {
	e.Labels = modelpb.Labels(e.Labels).Clone()
	e.NumericLabels = modelpb.NumericLabels(e.NumericLabels).Clone()
}

func setLabel(key string, event *modelpb.APMEvent, v interface{}) {
	switch v := v.(type) {
	case string:
		modelpb.Labels(event.Labels).Set(key, v)
	case bool:
		modelpb.Labels(event.Labels).Set(key, strconv.FormatBool(v))
	case float64:
		modelpb.NumericLabels(event.NumericLabels).Set(key, v)
	case int64:
		modelpb.NumericLabels(event.NumericLabels).Set(key, float64(v))
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
			modelpb.Labels(event.Labels).SetSlice(key, value)
		case float64:
			value := make([]float64, len(v))
			for i := range v {
				value[i] = v[i].(float64)
			}
			modelpb.NumericLabels(event.NumericLabels).SetSlice(key, value)
		}
	}
}

func populateNil[T any](a *T) *T {
	if a == nil {
		return new(T)
	}
	return a
}
