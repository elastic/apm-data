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
// source: source.proto

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

type Source struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ip            *IP                    `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Nat           *NAT                   `protobuf:"bytes,2,opt,name=nat,proto3" json:"nat,omitempty"`
	Domain        string                 `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Port          uint32                 `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Source) Reset() {
	*x = Source{}
	mi := &file_source_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[0]
	if x != nil {
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ip            *IP                    `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NAT) Reset() {
	*x = NAT{}
	mi := &file_source_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NAT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NAT) ProtoMessage() {}

func (x *NAT) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[1]
	if x != nil {
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

const file_source_proto_rawDesc = "" +
	"\n" +
	"\fsource.proto\x12\x0eelastic.apm.v1\x1a\bip.proto\"\x7f\n" +
	"\x06Source\x12\"\n" +
	"\x02ip\x18\x01 \x01(\v2\x12.elastic.apm.v1.IPR\x02ip\x12%\n" +
	"\x03nat\x18\x02 \x01(\v2\x13.elastic.apm.v1.NATR\x03nat\x12\x16\n" +
	"\x06domain\x18\x03 \x01(\tR\x06domain\x12\x12\n" +
	"\x04port\x18\x04 \x01(\rR\x04port\")\n" +
	"\x03NAT\x12\"\n" +
	"\x02ip\x18\x01 \x01(\v2\x12.elastic.apm.v1.IPR\x02ipB+Z)github.com/elastic/apm-data/model/modelpbb\x06proto3"

var (
	file_source_proto_rawDescOnce sync.Once
	file_source_proto_rawDescData []byte
)

func file_source_proto_rawDescGZIP() []byte {
	file_source_proto_rawDescOnce.Do(func() {
		file_source_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_source_proto_rawDesc), len(file_source_proto_rawDesc)))
	})
	return file_source_proto_rawDescData
}

var file_source_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_source_proto_goTypes = []any{
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_source_proto_rawDesc), len(file_source_proto_rawDesc)),
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
	file_source_proto_goTypes = nil
	file_source_proto_depIdxs = nil
}
