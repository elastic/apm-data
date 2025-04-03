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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.28.3
// source: cloud.proto

package modelpb

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Cloud struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	Origin           *CloudOrigin           `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
	AccountId        string                 `protobuf:"bytes,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	AccountName      string                 `protobuf:"bytes,3,opt,name=account_name,json=accountName,proto3" json:"account_name,omitempty"`
	AvailabilityZone string                 `protobuf:"bytes,4,opt,name=availability_zone,json=availabilityZone,proto3" json:"availability_zone,omitempty"`
	InstanceId       string                 `protobuf:"bytes,5,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	InstanceName     string                 `protobuf:"bytes,6,opt,name=instance_name,json=instanceName,proto3" json:"instance_name,omitempty"`
	MachineType      string                 `protobuf:"bytes,7,opt,name=machine_type,json=machineType,proto3" json:"machine_type,omitempty"`
	ProjectId        string                 `protobuf:"bytes,8,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	ProjectName      string                 `protobuf:"bytes,9,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
	Provider         string                 `protobuf:"bytes,10,opt,name=provider,proto3" json:"provider,omitempty"`
	Region           string                 `protobuf:"bytes,11,opt,name=region,proto3" json:"region,omitempty"`
	ServiceName      string                 `protobuf:"bytes,12,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *Cloud) Reset() {
	*x = Cloud{}
	mi := &file_cloud_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cloud) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cloud) ProtoMessage() {}

func (x *Cloud) ProtoReflect() protoreflect.Message {
	mi := &file_cloud_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cloud.ProtoReflect.Descriptor instead.
func (*Cloud) Descriptor() ([]byte, []int) {
	return file_cloud_proto_rawDescGZIP(), []int{0}
}

func (x *Cloud) GetOrigin() *CloudOrigin {
	if x != nil {
		return x.Origin
	}
	return nil
}

func (x *Cloud) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Cloud) GetAccountName() string {
	if x != nil {
		return x.AccountName
	}
	return ""
}

func (x *Cloud) GetAvailabilityZone() string {
	if x != nil {
		return x.AvailabilityZone
	}
	return ""
}

func (x *Cloud) GetInstanceId() string {
	if x != nil {
		return x.InstanceId
	}
	return ""
}

func (x *Cloud) GetInstanceName() string {
	if x != nil {
		return x.InstanceName
	}
	return ""
}

func (x *Cloud) GetMachineType() string {
	if x != nil {
		return x.MachineType
	}
	return ""
}

func (x *Cloud) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *Cloud) GetProjectName() string {
	if x != nil {
		return x.ProjectName
	}
	return ""
}

func (x *Cloud) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *Cloud) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Cloud) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

type CloudOrigin struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccountId     string                 `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Provider      string                 `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	Region        string                 `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	ServiceName   string                 `protobuf:"bytes,4,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CloudOrigin) Reset() {
	*x = CloudOrigin{}
	mi := &file_cloud_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CloudOrigin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloudOrigin) ProtoMessage() {}

func (x *CloudOrigin) ProtoReflect() protoreflect.Message {
	mi := &file_cloud_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloudOrigin.ProtoReflect.Descriptor instead.
func (*CloudOrigin) Descriptor() ([]byte, []int) {
	return file_cloud_proto_rawDescGZIP(), []int{1}
}

func (x *CloudOrigin) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *CloudOrigin) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *CloudOrigin) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *CloudOrigin) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

var File_cloud_proto protoreflect.FileDescriptor

const file_cloud_proto_rawDesc = "" +
	"\n" +
	"\vcloud.proto\x12\x0eelastic.apm.v1\"\xad\x03\n" +
	"\x05Cloud\x123\n" +
	"\x06origin\x18\x01 \x01(\v2\x1b.elastic.apm.v1.CloudOriginR\x06origin\x12\x1d\n" +
	"\n" +
	"account_id\x18\x02 \x01(\tR\taccountId\x12!\n" +
	"\faccount_name\x18\x03 \x01(\tR\vaccountName\x12+\n" +
	"\x11availability_zone\x18\x04 \x01(\tR\x10availabilityZone\x12\x1f\n" +
	"\vinstance_id\x18\x05 \x01(\tR\n" +
	"instanceId\x12#\n" +
	"\rinstance_name\x18\x06 \x01(\tR\finstanceName\x12!\n" +
	"\fmachine_type\x18\a \x01(\tR\vmachineType\x12\x1d\n" +
	"\n" +
	"project_id\x18\b \x01(\tR\tprojectId\x12!\n" +
	"\fproject_name\x18\t \x01(\tR\vprojectName\x12\x1a\n" +
	"\bprovider\x18\n" +
	" \x01(\tR\bprovider\x12\x16\n" +
	"\x06region\x18\v \x01(\tR\x06region\x12!\n" +
	"\fservice_name\x18\f \x01(\tR\vserviceName\"\x83\x01\n" +
	"\vCloudOrigin\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\x12\x1a\n" +
	"\bprovider\x18\x02 \x01(\tR\bprovider\x12\x16\n" +
	"\x06region\x18\x03 \x01(\tR\x06region\x12!\n" +
	"\fservice_name\x18\x04 \x01(\tR\vserviceNameB+Z)github.com/elastic/apm-data/model/modelpbb\x06proto3"

var (
	file_cloud_proto_rawDescOnce sync.Once
	file_cloud_proto_rawDescData []byte
)

func file_cloud_proto_rawDescGZIP() []byte {
	file_cloud_proto_rawDescOnce.Do(func() {
		file_cloud_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_cloud_proto_rawDesc), len(file_cloud_proto_rawDesc)))
	})
	return file_cloud_proto_rawDescData
}

var file_cloud_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_cloud_proto_goTypes = []any{
	(*Cloud)(nil),       // 0: elastic.apm.v1.Cloud
	(*CloudOrigin)(nil), // 1: elastic.apm.v1.CloudOrigin
}
var file_cloud_proto_depIdxs = []int32{
	1, // 0: elastic.apm.v1.Cloud.origin:type_name -> elastic.apm.v1.CloudOrigin
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_cloud_proto_init() }
func file_cloud_proto_init() {
	if File_cloud_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_cloud_proto_rawDesc), len(file_cloud_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cloud_proto_goTypes,
		DependencyIndexes: file_cloud_proto_depIdxs,
		MessageInfos:      file_cloud_proto_msgTypes,
	}.Build()
	File_cloud_proto = out.File
	file_cloud_proto_goTypes = nil
	file_cloud_proto_depIdxs = nil
}
