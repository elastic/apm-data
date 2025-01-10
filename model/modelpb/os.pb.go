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
// 	protoc-gen-go v1.36.2
// 	protoc        v5.28.3
// source: os.proto

package modelpb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OS struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version       string                 `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Platform      string                 `protobuf:"bytes,3,opt,name=platform,proto3" json:"platform,omitempty"`
	Full          string                 `protobuf:"bytes,4,opt,name=full,proto3" json:"full,omitempty"`
	Type          string                 `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OS) Reset() {
	*x = OS{}
	mi := &file_os_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OS) ProtoMessage() {}

func (x *OS) ProtoReflect() protoreflect.Message {
	mi := &file_os_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OS.ProtoReflect.Descriptor instead.
func (*OS) Descriptor() ([]byte, []int) {
	return file_os_proto_rawDescGZIP(), []int{0}
}

func (x *OS) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OS) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *OS) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *OS) GetFull() string {
	if x != nil {
		return x.Full
	}
	return ""
}

func (x *OS) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_os_proto protoreflect.FileDescriptor

var file_os_proto_rawDesc = []byte{
	0x0a, 0x08, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x6c, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x22, 0x76, 0x0a, 0x02, 0x4f, 0x53,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x75,
	0x6c, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x75, 0x6c, 0x6c, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x61, 0x70, 0x6d, 0x2d, 0x64, 0x61, 0x74,
	0x61, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_os_proto_rawDescOnce sync.Once
	file_os_proto_rawDescData = file_os_proto_rawDesc
)

func file_os_proto_rawDescGZIP() []byte {
	file_os_proto_rawDescOnce.Do(func() {
		file_os_proto_rawDescData = protoimpl.X.CompressGZIP(file_os_proto_rawDescData)
	})
	return file_os_proto_rawDescData
}

var file_os_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_os_proto_goTypes = []any{
	(*OS)(nil), // 0: elastic.apm.v1.OS
}
var file_os_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_os_proto_init() }
func file_os_proto_init() {
	if File_os_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_os_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_os_proto_goTypes,
		DependencyIndexes: file_os_proto_depIdxs,
		MessageInfos:      file_os_proto_msgTypes,
	}.Build()
	File_os_proto = out.File
	file_os_proto_rawDesc = nil
	file_os_proto_goTypes = nil
	file_os_proto_depIdxs = nil
}
