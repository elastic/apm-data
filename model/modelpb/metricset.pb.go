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
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: metricset.proto

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

type MetricType int32

const (
	MetricType_METRIC_TYPE_UNSPECIFIED MetricType = 0
	MetricType_METRIC_TYPE_GAUGE       MetricType = 1
	MetricType_METRIC_TYPE_COUNTER     MetricType = 2
	MetricType_METRIC_TYPE_HISTOGRAM   MetricType = 3
	MetricType_METRIC_TYPE_SUMMARY     MetricType = 4
)

// Enum value maps for MetricType.
var (
	MetricType_name = map[int32]string{
		0: "METRIC_TYPE_UNSPECIFIED",
		1: "METRIC_TYPE_GAUGE",
		2: "METRIC_TYPE_COUNTER",
		3: "METRIC_TYPE_HISTOGRAM",
		4: "METRIC_TYPE_SUMMARY",
	}
	MetricType_value = map[string]int32{
		"METRIC_TYPE_UNSPECIFIED": 0,
		"METRIC_TYPE_GAUGE":       1,
		"METRIC_TYPE_COUNTER":     2,
		"METRIC_TYPE_HISTOGRAM":   3,
		"METRIC_TYPE_SUMMARY":     4,
	}
)

func (x MetricType) Enum() *MetricType {
	p := new(MetricType)
	*p = x
	return p
}

func (x MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_metricset_proto_enumTypes[0].Descriptor()
}

func (MetricType) Type() protoreflect.EnumType {
	return &file_metricset_proto_enumTypes[0]
}

func (x MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetricType.Descriptor instead.
func (MetricType) EnumDescriptor() ([]byte, []int) {
	return file_metricset_proto_rawDescGZIP(), []int{0}
}

type Metricset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Interval string             `protobuf:"bytes,2,opt,name=interval,proto3" json:"interval,omitempty"`
	Samples  []*MetricsetSample `protobuf:"bytes,3,rep,name=samples,proto3" json:"samples,omitempty"`
	DocCount uint64             `protobuf:"varint,4,opt,name=doc_count,json=docCount,proto3" json:"doc_count,omitempty"`
}

func (x *Metricset) Reset() {
	*x = Metricset{}
	mi := &file_metricset_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metricset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metricset) ProtoMessage() {}

func (x *Metricset) ProtoReflect() protoreflect.Message {
	mi := &file_metricset_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metricset.ProtoReflect.Descriptor instead.
func (*Metricset) Descriptor() ([]byte, []int) {
	return file_metricset_proto_rawDescGZIP(), []int{0}
}

func (x *Metricset) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metricset) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

func (x *Metricset) GetSamples() []*MetricsetSample {
	if x != nil {
		return x.Samples
	}
	return nil
}

func (x *Metricset) GetDocCount() uint64 {
	if x != nil {
		return x.DocCount
	}
	return 0
}

type MetricsetSample struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      MetricType     `protobuf:"varint,1,opt,name=type,proto3,enum=elastic.apm.v1.MetricType" json:"type,omitempty"`
	Name      string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Unit      string         `protobuf:"bytes,3,opt,name=unit,proto3" json:"unit,omitempty"`
	Histogram *Histogram     `protobuf:"bytes,4,opt,name=histogram,proto3" json:"histogram,omitempty"`
	Summary   *SummaryMetric `protobuf:"bytes,5,opt,name=summary,proto3" json:"summary,omitempty"`
	Value     float64        `protobuf:"fixed64,6,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MetricsetSample) Reset() {
	*x = MetricsetSample{}
	mi := &file_metricset_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetricsetSample) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsetSample) ProtoMessage() {}

func (x *MetricsetSample) ProtoReflect() protoreflect.Message {
	mi := &file_metricset_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsetSample.ProtoReflect.Descriptor instead.
func (*MetricsetSample) Descriptor() ([]byte, []int) {
	return file_metricset_proto_rawDescGZIP(), []int{1}
}

func (x *MetricsetSample) GetType() MetricType {
	if x != nil {
		return x.Type
	}
	return MetricType_METRIC_TYPE_UNSPECIFIED
}

func (x *MetricsetSample) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MetricsetSample) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

func (x *MetricsetSample) GetHistogram() *Histogram {
	if x != nil {
		return x.Histogram
	}
	return nil
}

func (x *MetricsetSample) GetSummary() *SummaryMetric {
	if x != nil {
		return x.Summary
	}
	return nil
}

func (x *MetricsetSample) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Histogram struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []float64 `protobuf:"fixed64,1,rep,packed,name=values,proto3" json:"values,omitempty"`
	Counts []uint64  `protobuf:"varint,2,rep,packed,name=counts,proto3" json:"counts,omitempty"`
}

func (x *Histogram) Reset() {
	*x = Histogram{}
	mi := &file_metricset_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Histogram) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Histogram) ProtoMessage() {}

func (x *Histogram) ProtoReflect() protoreflect.Message {
	mi := &file_metricset_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Histogram.ProtoReflect.Descriptor instead.
func (*Histogram) Descriptor() ([]byte, []int) {
	return file_metricset_proto_rawDescGZIP(), []int{2}
}

func (x *Histogram) GetValues() []float64 {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *Histogram) GetCounts() []uint64 {
	if x != nil {
		return x.Counts
	}
	return nil
}

type SummaryMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count uint64  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	Sum   float64 `protobuf:"fixed64,2,opt,name=sum,proto3" json:"sum,omitempty"`
}

func (x *SummaryMetric) Reset() {
	*x = SummaryMetric{}
	mi := &file_metricset_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SummaryMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryMetric) ProtoMessage() {}

func (x *SummaryMetric) ProtoReflect() protoreflect.Message {
	mi := &file_metricset_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryMetric.ProtoReflect.Descriptor instead.
func (*SummaryMetric) Descriptor() ([]byte, []int) {
	return file_metricset_proto_rawDescGZIP(), []int{3}
}

func (x *SummaryMetric) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *SummaryMetric) GetSum() float64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

type AggregatedDuration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count uint64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	// nanoseconds
	Sum uint64 `protobuf:"varint,2,opt,name=sum,proto3" json:"sum,omitempty"`
}

func (x *AggregatedDuration) Reset() {
	*x = AggregatedDuration{}
	mi := &file_metricset_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AggregatedDuration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatedDuration) ProtoMessage() {}

func (x *AggregatedDuration) ProtoReflect() protoreflect.Message {
	mi := &file_metricset_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatedDuration.ProtoReflect.Descriptor instead.
func (*AggregatedDuration) Descriptor() ([]byte, []int) {
	return file_metricset_proto_rawDescGZIP(), []int{4}
}

func (x *AggregatedDuration) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *AggregatedDuration) GetSum() uint64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

var File_metricset_proto protoreflect.FileDescriptor

var file_metricset_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76,
	0x31, 0x22, 0x93, 0x01, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12,
	0x39, 0x0a, 0x07, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76,
	0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x74, 0x53, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x52, 0x07, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6f,
	0x63, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x64,
	0x6f, 0x63, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xf1, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x65, 0x74, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x65, 0x6c, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75,
	0x6e, 0x69, 0x74, 0x12, 0x37, 0x0a, 0x09, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x67, 0x72, 0x61, 0x6d,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63,
	0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x67, 0x72, 0x61,
	0x6d, 0x52, 0x09, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x37, 0x0a, 0x07,
	0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x73, 0x75,
	0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3b, 0x0a, 0x09, 0x48,
	0x69, 0x73, 0x74, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x01, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x04,
	0x52, 0x06, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x22, 0x37, 0x0a, 0x0d, 0x53, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x73, 0x75,
	0x6d, 0x22, 0x3c, 0x0a, 0x12, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x73, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x73, 0x75, 0x6d, 0x2a,
	0x8d, 0x01, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b,
	0x0a, 0x17, 0x4d, 0x45, 0x54, 0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x4d,
	0x45, 0x54, 0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x41, 0x55, 0x47, 0x45,
	0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x45, 0x54, 0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x4d,
	0x45, 0x54, 0x52, 0x49, 0x43, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x49, 0x53, 0x54, 0x4f,
	0x47, 0x52, 0x41, 0x4d, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x4d, 0x45, 0x54, 0x52, 0x49, 0x43,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x55, 0x4d, 0x4d, 0x41, 0x52, 0x59, 0x10, 0x04, 0x42,
	0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6c,
	0x61, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x61, 0x70, 0x6d, 0x2d, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metricset_proto_rawDescOnce sync.Once
	file_metricset_proto_rawDescData = file_metricset_proto_rawDesc
)

func file_metricset_proto_rawDescGZIP() []byte {
	file_metricset_proto_rawDescOnce.Do(func() {
		file_metricset_proto_rawDescData = protoimpl.X.CompressGZIP(file_metricset_proto_rawDescData)
	})
	return file_metricset_proto_rawDescData
}

var file_metricset_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_metricset_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_metricset_proto_goTypes = []any{
	(MetricType)(0),            // 0: elastic.apm.v1.MetricType
	(*Metricset)(nil),          // 1: elastic.apm.v1.Metricset
	(*MetricsetSample)(nil),    // 2: elastic.apm.v1.MetricsetSample
	(*Histogram)(nil),          // 3: elastic.apm.v1.Histogram
	(*SummaryMetric)(nil),      // 4: elastic.apm.v1.SummaryMetric
	(*AggregatedDuration)(nil), // 5: elastic.apm.v1.AggregatedDuration
}
var file_metricset_proto_depIdxs = []int32{
	2, // 0: elastic.apm.v1.Metricset.samples:type_name -> elastic.apm.v1.MetricsetSample
	0, // 1: elastic.apm.v1.MetricsetSample.type:type_name -> elastic.apm.v1.MetricType
	3, // 2: elastic.apm.v1.MetricsetSample.histogram:type_name -> elastic.apm.v1.Histogram
	4, // 3: elastic.apm.v1.MetricsetSample.summary:type_name -> elastic.apm.v1.SummaryMetric
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_metricset_proto_init() }
func file_metricset_proto_init() {
	if File_metricset_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metricset_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_metricset_proto_goTypes,
		DependencyIndexes: file_metricset_proto_depIdxs,
		EnumInfos:         file_metricset_proto_enumTypes,
		MessageInfos:      file_metricset_proto_msgTypes,
	}.Build()
	File_metricset_proto = out.File
	file_metricset_proto_rawDesc = nil
	file_metricset_proto_goTypes = nil
	file_metricset_proto_depIdxs = nil
}
