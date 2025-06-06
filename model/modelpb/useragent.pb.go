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
// source: useragent.proto

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

type UserAgent struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Original      string                 `protobuf:"bytes,1,opt,name=original,proto3" json:"original,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserAgent) Reset() {
	*x = UserAgent{}
	mi := &file_useragent_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserAgent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserAgent) ProtoMessage() {}

func (x *UserAgent) ProtoReflect() protoreflect.Message {
	mi := &file_useragent_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserAgent.ProtoReflect.Descriptor instead.
func (*UserAgent) Descriptor() ([]byte, []int) {
	return file_useragent_proto_rawDescGZIP(), []int{0}
}

func (x *UserAgent) GetOriginal() string {
	if x != nil {
		return x.Original
	}
	return ""
}

func (x *UserAgent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_useragent_proto protoreflect.FileDescriptor

const file_useragent_proto_rawDesc = "" +
	"\n" +
	"\x0fuseragent.proto\x12\x0eelastic.apm.v1\";\n" +
	"\tUserAgent\x12\x1a\n" +
	"\boriginal\x18\x01 \x01(\tR\boriginal\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04nameB+Z)github.com/elastic/apm-data/model/modelpbb\x06proto3"

var (
	file_useragent_proto_rawDescOnce sync.Once
	file_useragent_proto_rawDescData []byte
)

func file_useragent_proto_rawDescGZIP() []byte {
	file_useragent_proto_rawDescOnce.Do(func() {
		file_useragent_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_useragent_proto_rawDesc), len(file_useragent_proto_rawDesc)))
	})
	return file_useragent_proto_rawDescData
}

var file_useragent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_useragent_proto_goTypes = []any{
	(*UserAgent)(nil), // 0: elastic.apm.v1.UserAgent
}
var file_useragent_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_useragent_proto_init() }
func file_useragent_proto_init() {
	if File_useragent_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_useragent_proto_rawDesc), len(file_useragent_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_useragent_proto_goTypes,
		DependencyIndexes: file_useragent_proto_depIdxs,
		MessageInfos:      file_useragent_proto_msgTypes,
	}.Build()
	File_useragent_proto = out.File
	file_useragent_proto_goTypes = nil
	file_useragent_proto_depIdxs = nil
}
