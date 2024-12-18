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
// source: transaction.proto

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

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SpanCount             *SpanCount                  `protobuf:"bytes,1,opt,name=span_count,json=spanCount,proto3" json:"span_count,omitempty"`
	UserExperience        *UserExperience             `protobuf:"bytes,2,opt,name=user_experience,json=userExperience,proto3" json:"user_experience,omitempty"`
	Custom                []*KeyValue                 `protobuf:"bytes,3,rep,name=custom,proto3" json:"custom,omitempty"`
	Marks                 map[string]*TransactionMark `protobuf:"bytes,4,rep,name=marks,proto3" json:"marks,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Message               *Message                    `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
	Type                  string                      `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	Name                  string                      `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	Result                string                      `protobuf:"bytes,8,opt,name=result,proto3" json:"result,omitempty"`
	Id                    string                      `protobuf:"bytes,9,opt,name=id,proto3" json:"id,omitempty"`
	DurationHistogram     *Histogram                  `protobuf:"bytes,10,opt,name=duration_histogram,json=durationHistogram,proto3" json:"duration_histogram,omitempty"`
	DroppedSpansStats     []*DroppedSpanStats         `protobuf:"bytes,11,rep,name=dropped_spans_stats,json=droppedSpansStats,proto3" json:"dropped_spans_stats,omitempty"`
	DurationSummary       *SummaryMetric              `protobuf:"bytes,12,opt,name=duration_summary,json=durationSummary,proto3" json:"duration_summary,omitempty"`
	RepresentativeCount   float64                     `protobuf:"fixed64,13,opt,name=representative_count,json=representativeCount,proto3" json:"representative_count,omitempty"`
	Sampled               bool                        `protobuf:"varint,14,opt,name=sampled,proto3" json:"sampled,omitempty"`
	Root                  bool                        `protobuf:"varint,15,opt,name=root,proto3" json:"root,omitempty"`
	ProfilerStackTraceIds []string                    `protobuf:"bytes,16,rep,name=profiler_stack_trace_ids,json=profilerStackTraceIds,proto3" json:"profiler_stack_trace_ids,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	mi := &file_transaction_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetSpanCount() *SpanCount {
	if x != nil {
		return x.SpanCount
	}
	return nil
}

func (x *Transaction) GetUserExperience() *UserExperience {
	if x != nil {
		return x.UserExperience
	}
	return nil
}

func (x *Transaction) GetCustom() []*KeyValue {
	if x != nil {
		return x.Custom
	}
	return nil
}

func (x *Transaction) GetMarks() map[string]*TransactionMark {
	if x != nil {
		return x.Marks
	}
	return nil
}

func (x *Transaction) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *Transaction) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Transaction) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Transaction) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *Transaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Transaction) GetDurationHistogram() *Histogram {
	if x != nil {
		return x.DurationHistogram
	}
	return nil
}

func (x *Transaction) GetDroppedSpansStats() []*DroppedSpanStats {
	if x != nil {
		return x.DroppedSpansStats
	}
	return nil
}

func (x *Transaction) GetDurationSummary() *SummaryMetric {
	if x != nil {
		return x.DurationSummary
	}
	return nil
}

func (x *Transaction) GetRepresentativeCount() float64 {
	if x != nil {
		return x.RepresentativeCount
	}
	return 0
}

func (x *Transaction) GetSampled() bool {
	if x != nil {
		return x.Sampled
	}
	return false
}

func (x *Transaction) GetRoot() bool {
	if x != nil {
		return x.Root
	}
	return false
}

func (x *Transaction) GetProfilerStackTraceIds() []string {
	if x != nil {
		return x.ProfilerStackTraceIds
	}
	return nil
}

type SpanCount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dropped *uint32 `protobuf:"varint,1,opt,name=dropped,proto3,oneof" json:"dropped,omitempty"`
	Started *uint32 `protobuf:"varint,2,opt,name=started,proto3,oneof" json:"started,omitempty"`
}

func (x *SpanCount) Reset() {
	*x = SpanCount{}
	mi := &file_transaction_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SpanCount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SpanCount) ProtoMessage() {}

func (x *SpanCount) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SpanCount.ProtoReflect.Descriptor instead.
func (*SpanCount) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{1}
}

func (x *SpanCount) GetDropped() uint32 {
	if x != nil && x.Dropped != nil {
		return *x.Dropped
	}
	return 0
}

func (x *SpanCount) GetStarted() uint32 {
	if x != nil && x.Started != nil {
		return *x.Started
	}
	return 0
}

type TransactionMark struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Measurements map[string]float64 `protobuf:"bytes,1,rep,name=measurements,proto3" json:"measurements,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
}

func (x *TransactionMark) Reset() {
	*x = TransactionMark{}
	mi := &file_transaction_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionMark) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionMark) ProtoMessage() {}

func (x *TransactionMark) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionMark.ProtoReflect.Descriptor instead.
func (*TransactionMark) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{2}
}

func (x *TransactionMark) GetMeasurements() map[string]float64 {
	if x != nil {
		return x.Measurements
	}
	return nil
}

type DroppedSpanStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DestinationServiceResource string              `protobuf:"bytes,1,opt,name=destination_service_resource,json=destinationServiceResource,proto3" json:"destination_service_resource,omitempty"`
	ServiceTargetType          string              `protobuf:"bytes,2,opt,name=service_target_type,json=serviceTargetType,proto3" json:"service_target_type,omitempty"`
	ServiceTargetName          string              `protobuf:"bytes,3,opt,name=service_target_name,json=serviceTargetName,proto3" json:"service_target_name,omitempty"`
	Outcome                    string              `protobuf:"bytes,4,opt,name=outcome,proto3" json:"outcome,omitempty"`
	Duration                   *AggregatedDuration `protobuf:"bytes,5,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *DroppedSpanStats) Reset() {
	*x = DroppedSpanStats{}
	mi := &file_transaction_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DroppedSpanStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DroppedSpanStats) ProtoMessage() {}

func (x *DroppedSpanStats) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DroppedSpanStats.ProtoReflect.Descriptor instead.
func (*DroppedSpanStats) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{3}
}

func (x *DroppedSpanStats) GetDestinationServiceResource() string {
	if x != nil {
		return x.DestinationServiceResource
	}
	return ""
}

func (x *DroppedSpanStats) GetServiceTargetType() string {
	if x != nil {
		return x.ServiceTargetType
	}
	return ""
}

func (x *DroppedSpanStats) GetServiceTargetName() string {
	if x != nil {
		return x.ServiceTargetName
	}
	return ""
}

func (x *DroppedSpanStats) GetOutcome() string {
	if x != nil {
		return x.Outcome
	}
	return ""
}

func (x *DroppedSpanStats) GetDuration() *AggregatedDuration {
	if x != nil {
		return x.Duration
	}
	return nil
}

var File_transaction_proto protoreflect.FileDescriptor

var file_transaction_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d,
	0x2e, 0x76, 0x31, 0x1a, 0x10, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xde, 0x06, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x0a, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6c, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x52, 0x09, 0x73, 0x70, 0x61, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x47, 0x0a, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65, 0x6e,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x45, 0x78,
	0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x45, 0x78,
	0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x06, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x06, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x12, 0x3c, 0x0a, 0x05, 0x6d, 0x61,
	0x72, 0x6b, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x65, 0x6c, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x05, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x12, 0x31, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x65, 0x6c, 0x61, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x48, 0x0a, 0x12, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x67, 0x72, 0x61,
	0x6d, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69,
	0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x67, 0x72,
	0x61, 0x6d, 0x52, 0x11, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x69, 0x73, 0x74,
	0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x50, 0x0a, 0x13, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64,
	0x5f, 0x73, 0x70, 0x61, 0x6e, 0x73, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x0b, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x20, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x53, 0x70, 0x61, 0x6e, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x11, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x53, 0x70, 0x61,
	0x6e, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x48, 0x0a, 0x10, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x52, 0x0f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x12, 0x31, 0x0a, 0x14, 0x72, 0x65, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x61, 0x74,
	0x69, 0x76, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x13, 0x72, 0x65, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x64, 0x18,
	0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x72, 0x6f,
	0x6f, 0x74, 0x12, 0x37, 0x0a, 0x18, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x72, 0x5f, 0x73,
	0x74, 0x61, 0x63, 0x6b, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x10,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x15, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x63, 0x6b, 0x54, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x73, 0x1a, 0x59, 0x0a, 0x0a, 0x4d,
	0x61, 0x72, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x35, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x65, 0x6c, 0x61,
	0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x72, 0x6b, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x61, 0x0a, 0x09, 0x53, 0x70, 0x61, 0x6e, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x07, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x07, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x48, 0x01, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x88, 0x01,
	0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x22, 0xa9, 0x01, 0x0a, 0x0f, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x72, 0x6b, 0x12, 0x55, 0x0a,
	0x0c, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x61, 0x70,
	0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x4d, 0x61, 0x72, 0x6b, 0x2e, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x1a, 0x3f, 0x0a, 0x11, 0x4d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8e, 0x02, 0x0a, 0x10, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65,
	0x64, 0x53, 0x70, 0x61, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x40, 0x0a, 0x1c, 0x64, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x1a, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x13,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x13,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6f, 0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f,
	0x75, 0x74, 0x63, 0x6f, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x65, 0x6c, 0x61, 0x73, 0x74,
	0x69, 0x63, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x65, 0x64, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6c, 0x61, 0x73, 0x74, 0x69, 0x63, 0x2f, 0x61, 0x70, 0x6d,
	0x2d, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transaction_proto_rawDescOnce sync.Once
	file_transaction_proto_rawDescData = file_transaction_proto_rawDesc
)

func file_transaction_proto_rawDescGZIP() []byte {
	file_transaction_proto_rawDescOnce.Do(func() {
		file_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_transaction_proto_rawDescData)
	})
	return file_transaction_proto_rawDescData
}

var file_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_transaction_proto_goTypes = []any{
	(*Transaction)(nil),        // 0: elastic.apm.v1.Transaction
	(*SpanCount)(nil),          // 1: elastic.apm.v1.SpanCount
	(*TransactionMark)(nil),    // 2: elastic.apm.v1.TransactionMark
	(*DroppedSpanStats)(nil),   // 3: elastic.apm.v1.DroppedSpanStats
	nil,                        // 4: elastic.apm.v1.Transaction.MarksEntry
	nil,                        // 5: elastic.apm.v1.TransactionMark.MeasurementsEntry
	(*UserExperience)(nil),     // 6: elastic.apm.v1.UserExperience
	(*KeyValue)(nil),           // 7: elastic.apm.v1.KeyValue
	(*Message)(nil),            // 8: elastic.apm.v1.Message
	(*Histogram)(nil),          // 9: elastic.apm.v1.Histogram
	(*SummaryMetric)(nil),      // 10: elastic.apm.v1.SummaryMetric
	(*AggregatedDuration)(nil), // 11: elastic.apm.v1.AggregatedDuration
}
var file_transaction_proto_depIdxs = []int32{
	1,  // 0: elastic.apm.v1.Transaction.span_count:type_name -> elastic.apm.v1.SpanCount
	6,  // 1: elastic.apm.v1.Transaction.user_experience:type_name -> elastic.apm.v1.UserExperience
	7,  // 2: elastic.apm.v1.Transaction.custom:type_name -> elastic.apm.v1.KeyValue
	4,  // 3: elastic.apm.v1.Transaction.marks:type_name -> elastic.apm.v1.Transaction.MarksEntry
	8,  // 4: elastic.apm.v1.Transaction.message:type_name -> elastic.apm.v1.Message
	9,  // 5: elastic.apm.v1.Transaction.duration_histogram:type_name -> elastic.apm.v1.Histogram
	3,  // 6: elastic.apm.v1.Transaction.dropped_spans_stats:type_name -> elastic.apm.v1.DroppedSpanStats
	10, // 7: elastic.apm.v1.Transaction.duration_summary:type_name -> elastic.apm.v1.SummaryMetric
	5,  // 8: elastic.apm.v1.TransactionMark.measurements:type_name -> elastic.apm.v1.TransactionMark.MeasurementsEntry
	11, // 9: elastic.apm.v1.DroppedSpanStats.duration:type_name -> elastic.apm.v1.AggregatedDuration
	2,  // 10: elastic.apm.v1.Transaction.MarksEntry.value:type_name -> elastic.apm.v1.TransactionMark
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_transaction_proto_init() }
func file_transaction_proto_init() {
	if File_transaction_proto != nil {
		return
	}
	file_experience_proto_init()
	file_message_proto_init()
	file_metricset_proto_init()
	file_keyvalue_proto_init()
	file_transaction_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transaction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_transaction_proto_goTypes,
		DependencyIndexes: file_transaction_proto_depIdxs,
		MessageInfos:      file_transaction_proto_msgTypes,
	}.Build()
	File_transaction_proto = out.File
	file_transaction_proto_rawDesc = nil
	file_transaction_proto_goTypes = nil
	file_transaction_proto_depIdxs = nil
}
