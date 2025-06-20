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
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv26 "go.opentelemetry.io/otel/semconv/v1.26.0"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"

	"github.com/elastic/apm-data/model/modelpb"
)

const (
	AgentNameJaeger          = "Jaeger"
	MaxDataStreamBytes       = 100
	DisallowedNamespaceRunes = "\\/*?\"<>| ,#:"
	DisallowedDatasetRunes   = "-\\/*?\"<>| ,#:"
)

var (
	serviceNameInvalidRegexp = regexp.MustCompile("[^a-zA-Z0-9 _-]")
)

func translateResourceMetadata(resource pcommon.Resource, out *modelpb.APMEvent) {
	var exporterVersion string
	resource.Attributes().Range(func(k string, v pcommon.Value) bool {
		switch k {
		// service.*
		case string(semconv.ServiceNameKey):
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			out.Service.Name = cleanServiceName(v.Str())
		case string(semconv.ServiceVersionKey):
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			out.Service.Version = truncate(v.Str())
		case string(semconv.ServiceInstanceIDKey):
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			if out.Service.Node == nil {
				out.Service.Node = &modelpb.ServiceNode{}
			}
			out.Service.Node.Name = truncate(v.Str())

		// deployment.*
		case string(semconv26.DeploymentEnvironmentKey), string(semconv.DeploymentEnvironmentNameKey):
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			out.Service.Environment = truncate(v.Str())

		// telemetry.sdk.*
		case string(semconv.TelemetrySDKNameKey):
			if out.Agent == nil {
				out.Agent = &modelpb.Agent{}
			}
			out.Agent.Name = truncate(v.Str())
		case string(semconv.TelemetrySDKVersionKey):
			if out.Agent == nil {
				out.Agent = &modelpb.Agent{}
			}
			out.Agent.Version = truncate(v.Str())
		case string(semconv.TelemetrySDKLanguageKey):
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			if out.Service.Language == nil {
				out.Service.Language = &modelpb.Language{}
			}
			out.Service.Language.Name = truncate(v.Str())

		// cloud.*
		case string(semconv.CloudProviderKey):
			if out.Cloud == nil {
				out.Cloud = &modelpb.Cloud{}
			}
			out.Cloud.Provider = truncate(v.Str())
		case string(semconv.CloudAccountIDKey):
			if out.Cloud == nil {
				out.Cloud = &modelpb.Cloud{}
			}
			out.Cloud.AccountId = truncate(v.Str())
		case string(semconv.CloudRegionKey):
			if out.Cloud == nil {
				out.Cloud = &modelpb.Cloud{}
			}
			out.Cloud.Region = truncate(v.Str())
		case string(semconv.CloudAvailabilityZoneKey):
			if out.Cloud == nil {
				out.Cloud = &modelpb.Cloud{}
			}
			out.Cloud.AvailabilityZone = truncate(v.Str())
		case string(semconv.CloudPlatformKey):
			if out.Cloud == nil {
				out.Cloud = &modelpb.Cloud{}
			}
			out.Cloud.ServiceName = truncate(v.Str())

		// container.*
		case string(semconv.ContainerNameKey):
			if out.Container == nil {
				out.Container = &modelpb.Container{}
			}
			out.Container.Name = truncate(v.Str())
		case string(semconv.ContainerIDKey):
			if out.Container == nil {
				out.Container = &modelpb.Container{}
			}
			out.Container.Id = truncate(v.Str())
		case string(semconv.ContainerImageNameKey):
			if out.Container == nil {
				out.Container = &modelpb.Container{}
			}
			out.Container.ImageName = truncate(v.Str())
		case "container.image.tag":
			if out.Container == nil {
				out.Container = &modelpb.Container{}
			}
			out.Container.ImageTag = truncate(v.Str())
		case "container.runtime":
			if out.Container == nil {
				out.Container = &modelpb.Container{}
			}
			out.Container.Runtime = truncate(v.Str())

		// k8s.*
		case string(semconv.K8SNamespaceNameKey):
			if out.Kubernetes == nil {
				out.Kubernetes = &modelpb.Kubernetes{}
			}
			out.Kubernetes.Namespace = truncate(v.Str())
		case string(semconv.K8SNodeNameKey):
			if out.Kubernetes == nil {
				out.Kubernetes = &modelpb.Kubernetes{}
			}
			out.Kubernetes.NodeName = truncate(v.Str())
		case string(semconv.K8SPodNameKey):
			if out.Kubernetes == nil {
				out.Kubernetes = &modelpb.Kubernetes{}
			}
			out.Kubernetes.PodName = truncate(v.Str())
		case string(semconv.K8SPodUIDKey):
			if out.Kubernetes == nil {
				out.Kubernetes = &modelpb.Kubernetes{}
			}
			out.Kubernetes.PodUid = truncate(v.Str())

		// host.*
		case string(semconv.HostNameKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			out.Host.Hostname = truncate(v.Str())
		case string(semconv.HostIDKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			out.Host.Id = truncate(v.Str())
		case string(semconv.HostTypeKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			out.Host.Type = truncate(v.Str())
		case "host.arch":
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			out.Host.Architecture = truncate(v.Str())
		case string(semconv.HostIPKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			if v.Type() != pcommon.ValueTypeSlice {
				break // switch
			}
			slice := v.Slice()
			result := make([]*modelpb.IP, 0, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				ip, err := modelpb.ParseIP(slice.At(i).Str())
				if err == nil {
					result = append(result, ip)
				}
			}
			out.Host.Ip = result

		// process.*
		case string(semconv.ProcessPIDKey):
			if out.Process == nil {
				out.Process = &modelpb.Process{}
			}
			out.Process.Pid = uint32(v.Int())
		case string(semconv.ProcessCommandLineKey):
			if out.Process == nil {
				out.Process = &modelpb.Process{}
			}
			out.Process.CommandLine = truncate(v.Str())
		case string(semconv.ProcessExecutablePathKey):
			if out.Process == nil {
				out.Process = &modelpb.Process{}
			}
			out.Process.Executable = truncate(v.Str())
		case "process.runtime.name":
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			if out.Service.Runtime == nil {
				out.Service.Runtime = &modelpb.Runtime{}
			}
			out.Service.Runtime.Name = truncate(v.Str())
		case "process.runtime.version":
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			if out.Service.Runtime == nil {
				out.Service.Runtime = &modelpb.Runtime{}
			}
			out.Service.Runtime.Version = truncate(v.Str())
		case string(semconv.ProcessOwnerKey):
			if out.User == nil {
				out.User = &modelpb.User{}
			}
			out.User.Name = truncate(v.Str())

		// os.*
		case string(semconv.OSTypeKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			if out.Host.Os == nil {
				out.Host.Os = &modelpb.OS{}
			}
			out.Host.Os.Platform = strings.ToLower(truncate(v.Str()))
		case string(semconv.OSDescriptionKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			if out.Host.Os == nil {
				out.Host.Os = &modelpb.OS{}
			}
			out.Host.Os.Full = truncate(v.Str())
		case string(semconv.OSNameKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			if out.Host.Os == nil {
				out.Host.Os = &modelpb.OS{}
			}
			out.Host.Os.Name = truncate(v.Str())
		case string(semconv.OSVersionKey):
			if out.Host == nil {
				out.Host = &modelpb.Host{}
			}
			if out.Host.Os == nil {
				out.Host.Os = &modelpb.OS{}
			}
			out.Host.Os.Version = truncate(v.Str())

		// device.*
		case string(semconv.DeviceIDKey):
			if out.Device == nil {
				out.Device = &modelpb.Device{}
			}
			out.Device.Id = truncate(v.Str())
		case string(semconv.DeviceModelIdentifierKey):
			if out.Device == nil {
				out.Device = &modelpb.Device{}
			}
			if out.Device.Model == nil {
				out.Device.Model = &modelpb.DeviceModel{}
			}
			out.Device.Model.Identifier = truncate(v.Str())
		case string(semconv.DeviceModelNameKey):
			if out.Device == nil {
				out.Device = &modelpb.Device{}
			}
			if out.Device.Model == nil {
				out.Device.Model = &modelpb.DeviceModel{}
			}
			out.Device.Model.Name = truncate(v.Str())
		case "device.manufacturer":
			if out.Device == nil {
				out.Device = &modelpb.Device{}
			}
			out.Device.Manufacturer = truncate(v.Str())

		// Legacy OpenCensus attributes.
		case "opencensus.exporterversion":
			exporterVersion = v.Str()

		// timestamp attribute to deal with time skew on mobile
		// devices. APM server should drop this field.
		case "telemetry.sdk.elastic_export_timestamp":
			// Do nothing.
		case "telemetry.distro.name":
		case "telemetry.distro.version":
			//distro version & name are handled below and should not end up as labels

		// data_stream.*
		case attributeDataStreamDataset:
			if out.DataStream == nil {
				out.DataStream = &modelpb.DataStream{}
			}
			out.DataStream.Dataset = sanitizeDataStreamDataset(v.Str())
		case attributeDataStreamNamespace:
			if out.DataStream == nil {
				out.DataStream = &modelpb.DataStream{}
			}
			out.DataStream.Namespace = sanitizeDataStreamNamespace(v.Str())
		default:
			if out.Labels == nil {
				out.Labels = make(modelpb.Labels)
			}
			if out.NumericLabels == nil {
				out.NumericLabels = make(modelpb.NumericLabels)
			}
			setLabel(replaceDots(k), out, v)
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
			if out.Service == nil {
				out.Service = &modelpb.Service{}
			}
			if out.Service.Language == nil {
				out.Service.Language = &modelpb.Language{}
			}
			out.Service.Language.Name = versionParts[1]
		}
		if out.Agent == nil {
			out.Agent = &modelpb.Agent{}
		}
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
			if ip, err := modelpb.ParseIP(systemIP.Value); err == nil {
				out.Host.Ip = []*modelpb.IP{ip}
			}
			delete(out.Labels, "ip")
		}
	}

	if out.GetService().GetName() == "" {
		if out.Service == nil {
			out.Service = &modelpb.Service{}
		}
		// service.name is a required field.
		out.Service.Name = "unknown"
	}
	if out.Agent == nil {
		out.Agent = &modelpb.Agent{}
	}
	if out.GetAgent().GetName() == "" {
		// agent.name is a required field.
		out.Agent.Name = "otlp"
	}
	if out.Agent.Version == "" {
		// agent.version is a required field.
		out.Agent.Version = "unknown"
	}

	distroName, distroNameSet := resource.Attributes().Get("telemetry.distro.name")
	distroVersion, distroVersionSet := resource.Attributes().Get("telemetry.distro.version")

	if distroNameSet && distroName.Str() != "" {
		agentLang := "unknown"
		if out.GetService().GetLanguage().GetName() != "" {
			agentLang = out.GetService().GetLanguage().GetName()
		}

		out.Agent.Name = fmt.Sprintf("%s/%s/%s", out.Agent.Name, agentLang, distroName.Str())

		//we intentionally do not want to fallback to the Otel SDK version if we have a distro name, this would only cause confusion
		out.Agent.Version = "unknown"
		if distroVersionSet && distroVersion.Str() != "" {
			out.Agent.Version = distroVersion.Str()
		}
	} else {
		//distro is not set, use just the language as suffix if present
		if out.GetService().GetLanguage().GetName() != "" {
			out.Agent.Name = fmt.Sprintf("%s/%s", out.Agent.Name, out.GetService().GetLanguage().GetName())
		}
	}

	if out.GetService().GetLanguage().GetName() == "" {
		if out.Service == nil {
			out.Service = &modelpb.Service{}
		}
		if out.Service.Language == nil {
			out.Service.Language = &modelpb.Language{}
		}
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

func translateScopeMetadata(scope pcommon.InstrumentationScope, out *modelpb.APMEvent) {
	scope.Attributes().Range(func(k string, v pcommon.Value) bool {
		switch k {
		// data_stream.*
		case attributeDataStreamDataset:
			if out.DataStream == nil {
				out.DataStream = &modelpb.DataStream{}
			}
			out.DataStream.Dataset = sanitizeDataStreamDataset(v.Str())
		case attributeDataStreamNamespace:
			if out.DataStream == nil {
				out.DataStream = &modelpb.DataStream{}
			}
			out.DataStream.Namespace = sanitizeDataStreamNamespace(v.Str())
		}
		return true
	})

	if name := scope.Name(); name != "" {
		if out.Service == nil {
			out.Service = &modelpb.Service{}
		}
		out.Service.Framework = &modelpb.Framework{
			Name:    name,
			Version: scope.Version(),
		}
	}
}

func cleanServiceName(name string) string {
	return serviceNameInvalidRegexp.ReplaceAllString(truncate(name), "_")
}

// initEventLabels initializes an event-specific labels from an event.
func initEventLabels(e *modelpb.APMEvent) {
	e.Labels = modelpb.Labels(e.Labels).Clone()
	e.NumericLabels = modelpb.NumericLabels(e.NumericLabels).Clone()
}

func setLabel(key string, event *modelpb.APMEvent, v pcommon.Value) {
	switch v.Type() {
	case pcommon.ValueTypeStr:
		modelpb.Labels(event.Labels).Set(key, truncate(v.Str()))
	case pcommon.ValueTypeBool:
		modelpb.Labels(event.Labels).Set(key, strconv.FormatBool(v.Bool()))
	case pcommon.ValueTypeInt:
		modelpb.NumericLabels(event.NumericLabels).Set(key, float64(v.Int()))
	case pcommon.ValueTypeDouble:
		modelpb.NumericLabels(event.NumericLabels).Set(key, v.Double())
	case pcommon.ValueTypeSlice:
		s := v.Slice()
		if s.Len() == 0 {
			return
		}
		switch s.At(0).Type() {
		case pcommon.ValueTypeStr:
			result := make([]string, 0, s.Len())
			for i := 0; i < s.Len(); i++ {
				r := s.At(i)
				if r.Type() == pcommon.ValueTypeStr {
					result = append(result, truncate(r.Str()))
				}
			}
			modelpb.Labels(event.Labels).SetSlice(key, result)
		case pcommon.ValueTypeBool:
			result := make([]string, 0, s.Len())
			for i := 0; i < s.Len(); i++ {
				r := s.At(i)
				if r.Type() == pcommon.ValueTypeBool {
					result = append(result, strconv.FormatBool(r.Bool()))
				}
			}
			modelpb.Labels(event.Labels).SetSlice(key, result)
		case pcommon.ValueTypeDouble:
			result := make([]float64, 0, s.Len())
			for i := 0; i < s.Len(); i++ {
				r := s.At(i)
				if r.Type() == pcommon.ValueTypeDouble {
					result = append(result, r.Double())
				}
			}
			modelpb.NumericLabels(event.NumericLabels).SetSlice(key, result)
		case pcommon.ValueTypeInt:
			result := make([]float64, 0, s.Len())
			for i := 0; i < s.Len(); i++ {
				r := s.At(i)
				if r.Type() == pcommon.ValueTypeInt {
					result = append(result, float64(r.Int()))
				}
			}
			modelpb.NumericLabels(event.NumericLabels).SetSlice(key, result)
		}
	}
}

// Sanitize the datastream fields (dataset, namespace) to apply restrictions
// as outlined in https://www.elastic.co/guide/en/ecs/current/ecs-data_stream.html
func sanitizeDataStreamDataset(field string) string {
	field = strings.Map(replaceReservedRune(DisallowedDatasetRunes), field)
	if len(field) > MaxDataStreamBytes {
		return field[:MaxDataStreamBytes]
	}

	return field
}

// Sanitize the datastream fields (dataset, namespace) to apply restrictions
// as outlined in https://www.elastic.co/guide/en/ecs/current/ecs-data_stream.html
func sanitizeDataStreamNamespace(field string) string {
	field = strings.Map(replaceReservedRune(DisallowedNamespaceRunes), field)
	if len(field) > MaxDataStreamBytes {
		return field[:MaxDataStreamBytes]
	}
	return field
}

func replaceReservedRune(disallowedRunes string) func(r rune) rune {
	return func(r rune) rune {
		if strings.ContainsRune(disallowedRunes, r) {
			return '_'
		}
		return unicode.ToLower(r)
	}
}
