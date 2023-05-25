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
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: url.proto

package modelpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type URL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Original string `protobuf:"bytes,1,opt,name=original,proto3" json:"original,omitempty"`
	Scheme   string `protobuf:"bytes,2,opt,name=scheme,proto3" json:"scheme,omitempty"`
	Full     string `protobuf:"bytes,3,opt,name=full,proto3" json:"full,omitempty"`
	Domain   string `protobuf:"bytes,4,opt,name=domain,proto3" json:"domain,omitempty"`
	Path     string `protobuf:"bytes,5,opt,name=path,proto3" json:"path,omitempty"`
	Query    string `protobuf:"bytes,6,opt,name=query,proto3" json:"query,omitempty"`
	Fragment string `protobuf:"bytes,7,opt,name=fragment,proto3" json:"fragment,omitempty"`
	Port     uint32 `protobuf:"varint,8,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *URL) Reset() {
	*x = URL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *URL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*URL) ProtoMessage() {}

func (x *URL) ProtoReflect() protoreflect.Message {
	mi := &file_url_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use URL.ProtoReflect.Descriptor instead.
func (*URL) Descriptor() ([]byte, []int) {
	return file_url_proto_rawDescGZIP(), []int{0}
}

func (x *URL) GetOriginal() string {
	if x != nil {
		return x.Original
	}
	return ""
}

func (x *URL) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

func (x *URL) GetFull() string {
	if x != nil {
		return x.Full
	}
	return ""
}

func (x *URL) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *URL) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *URL) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *URL) GetFragment() string {
	if x != nil {
		return x.Fragment
	}
	return ""
}

func (x *URL) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

var File_url_proto protoreflect.FileDescriptor

var file_url_proto_rawDesc = []byte{
	0x0a, 0x09, 0x75, 0x72, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x6c, 0x61,
	0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x22, 0xbf, 0x01, 0x0a, 0x03,
	0x55, 0x52, 0x4c, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x75, 0x6c, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x75, 0x6c, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x1a, 0x0a,
	0x08, 0x66, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x2b, 0x5a,
	0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2f, 0x61, 0x70, 0x6d, 0x2d, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_url_proto_rawDescOnce sync.Once
	file_url_proto_rawDescData = file_url_proto_rawDesc
)

func file_url_proto_rawDescGZIP() []byte {
	file_url_proto_rawDescOnce.Do(func() {
		file_url_proto_rawDescData = protoimpl.X.CompressGZIP(file_url_proto_rawDescData)
	})
	return file_url_proto_rawDescData
}

var file_url_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_url_proto_goTypes = []interface{}{
	(*URL)(nil), // 0: elastic.apm.v1.URL
}
var file_url_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_url_proto_init() }
func file_url_proto_init() {
	if File_url_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_url_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*URL); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_url_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_url_proto_goTypes,
		DependencyIndexes: file_url_proto_depIdxs,
		MessageInfos:      file_url_proto_msgTypes,
	}.Build()
	File_url_proto = out.File
	file_url_proto_rawDesc = nil
	file_url_proto_goTypes = nil
	file_url_proto_depIdxs = nil
}