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
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: source.proto

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

type Source struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip     *IP    `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Nat    *NAT   `protobuf:"bytes,2,opt,name=nat,proto3" json:"nat,omitempty"`
	Domain string `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Port   uint32 `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Source) Reset() {
	*x = Source{}
	if protoimpl.UnsafeEnabled {
		mi := &file_source_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{0}
}

func (x *Source) GetIp() *IP {
	if x != nil {
		return x.Ip
	}
	return nil
}

func (x *Source) GetNat() *NAT {
	if x != nil {
		return x.Nat
	}
	return nil
}

func (x *Source) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Source) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type NAT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip *IP `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *NAT) Reset() {
	*x = NAT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_source_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NAT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NAT) ProtoMessage() {}

func (x *NAT) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NAT.ProtoReflect.Descriptor instead.
func (*NAT) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{1}
}

func (x *NAT) GetIp() *IP {
	if x != nil {
		return x.Ip
	}
	return nil
}

var File_source_proto protoreflect.FileDescriptor

var file_source_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x1a, 0x08,
	0x69, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x12, 0x22, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x50, 0x52, 0x02, 0x69, 0x70, 0x12, 0x25, 0x0a, 0x03, 0x6e, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70,
	0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x41, 0x54, 0x52, 0x03, 0x6e, 0x61, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x29, 0x0a, 0x03, 0x4e, 0x41, 0x54,
	0x12, 0x22, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x65,
	0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x50,
	0x52, 0x02, 0x69, 0x70, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x61, 0x70, 0x6d, 0x2d, 0x64,
	0x61, 0x74, 0x61, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_source_proto_rawDescOnce sync.Once
	file_source_proto_rawDescData = file_source_proto_rawDesc
)

func file_source_proto_rawDescGZIP() []byte {
	file_source_proto_rawDescOnce.Do(func() {
		file_source_proto_rawDescData = protoimpl.X.CompressGZIP(file_source_proto_rawDescData)
	})
	return file_source_proto_rawDescData
}

var file_source_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_source_proto_goTypes = []interface{}{
	(*Source)(nil), // 0: elastic.apm.v1.Source
	(*NAT)(nil),    // 1: elastic.apm.v1.NAT
	(*IP)(nil),     // 2: elastic.apm.v1.IP
}
var file_source_proto_depIdxs = []int32{
	2, // 0: elastic.apm.v1.Source.ip:type_name -> elastic.apm.v1.IP
	1, // 1: elastic.apm.v1.Source.nat:type_name -> elastic.apm.v1.NAT
	2, // 2: elastic.apm.v1.NAT.ip:type_name -> elastic.apm.v1.IP
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_source_proto_init() }
func file_source_proto_init() {
	if File_source_proto != nil {
		return
	}
	file_ip_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_source_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Source); i {
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
		file_source_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NAT); i {
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
			RawDescriptor: file_source_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_source_proto_goTypes,
		DependencyIndexes: file_source_proto_depIdxs,
		MessageInfos:      file_source_proto_msgTypes,
	}.Build()
	File_source_proto = out.File
	file_source_proto_rawDesc = nil
	file_source_proto_goTypes = nil
	file_source_proto_depIdxs = nil
}
