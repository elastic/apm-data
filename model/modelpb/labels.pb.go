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
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: labels.proto

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

type LabelValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value  string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	Global bool     `protobuf:"varint,3,opt,name=global,proto3" json:"global,omitempty"`
}

func (x *LabelValue) Reset() {
	*x = LabelValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_labels_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelValue) ProtoMessage() {}

func (x *LabelValue) ProtoReflect() protoreflect.Message {
	mi := &file_labels_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelValue.ProtoReflect.Descriptor instead.
func (*LabelValue) Descriptor() ([]byte, []int) {
	return file_labels_proto_rawDescGZIP(), []int{0}
}

func (x *LabelValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *LabelValue) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *LabelValue) GetGlobal() bool {
	if x != nil {
		return x.Global
	}
	return false
}

type NumericLabelValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []float64 `protobuf:"fixed64,1,rep,packed,name=values,proto3" json:"values,omitempty"`
	Value  float64   `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
	Global bool      `protobuf:"varint,3,opt,name=global,proto3" json:"global,omitempty"`
}

func (x *NumericLabelValue) Reset() {
	*x = NumericLabelValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_labels_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NumericLabelValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NumericLabelValue) ProtoMessage() {}

func (x *NumericLabelValue) ProtoReflect() protoreflect.Message {
	mi := &file_labels_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NumericLabelValue.ProtoReflect.Descriptor instead.
func (*NumericLabelValue) Descriptor() ([]byte, []int) {
	return file_labels_proto_rawDescGZIP(), []int{1}
}

func (x *NumericLabelValue) GetValues() []float64 {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *NumericLabelValue) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *NumericLabelValue) GetGlobal() bool {
	if x != nil {
		return x.Global
	}
	return false
}

var File_labels_proto protoreflect.FileDescriptor

var file_labels_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x22, 0x52,
	0x0a, 0x0a, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x6c,
	0x6f, 0x62, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x67, 0x6c, 0x6f, 0x62,
	0x61, 0x6c, 0x22, 0x59, 0x0a, 0x11, 0x4e, 0x75, 0x6d, 0x65, 0x72, 0x69, 0x63, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x42, 0x2b, 0x5a,
	0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2f, 0x61, 0x70, 0x6d, 0x2d, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_labels_proto_rawDescOnce sync.Once
	file_labels_proto_rawDescData = file_labels_proto_rawDesc
)

func file_labels_proto_rawDescGZIP() []byte {
	file_labels_proto_rawDescOnce.Do(func() {
		file_labels_proto_rawDescData = protoimpl.X.CompressGZIP(file_labels_proto_rawDescData)
	})
	return file_labels_proto_rawDescData
}

var file_labels_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_labels_proto_goTypes = []interface{}{
	(*LabelValue)(nil),        // 0: elastic.apm.v1.LabelValue
	(*NumericLabelValue)(nil), // 1: elastic.apm.v1.NumericLabelValue
}
var file_labels_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_labels_proto_init() }
func file_labels_proto_init() {
	if File_labels_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_labels_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelValue); i {
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
		file_labels_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NumericLabelValue); i {
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
			RawDescriptor: file_labels_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_labels_proto_goTypes,
		DependencyIndexes: file_labels_proto_depIdxs,
		MessageInfos:      file_labels_proto_msgTypes,
	}.Build()
	File_labels_proto = out.File
	file_labels_proto_rawDesc = nil
	file_labels_proto_goTypes = nil
	file_labels_proto_depIdxs = nil
}
