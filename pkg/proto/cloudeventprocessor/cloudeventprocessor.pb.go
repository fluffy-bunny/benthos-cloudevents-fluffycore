// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: pkg/proto/cloudeventprocessor/cloudeventprocessor.proto

package cloudeventprocessor

import (
	cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
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

type ByteBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages [][]byte `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *ByteBatch) Reset() {
	*x = ByteBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ByteBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ByteBatch) ProtoMessage() {}

func (x *ByteBatch) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ByteBatch.ProtoReflect.Descriptor instead.
func (*ByteBatch) Descriptor() ([]byte, []int) {
	return file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescGZIP(), []int{0}
}

func (x *ByteBatch) GetMessages() [][]byte {
	if x != nil {
		return x.Messages
	}
	return nil
}

type ProcessCloudEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Channel  string                       `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	Batch    *cloudevents.CloudEventBatch `protobuf:"bytes,2,opt,name=batch,proto3" json:"batch,omitempty"`
	BadBatch *ByteBatch                   `protobuf:"bytes,3,opt,name=bad_batch,json=badBatch,proto3" json:"bad_batch,omitempty"`
}

func (x *ProcessCloudEventsRequest) Reset() {
	*x = ProcessCloudEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessCloudEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessCloudEventsRequest) ProtoMessage() {}

func (x *ProcessCloudEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessCloudEventsRequest.ProtoReflect.Descriptor instead.
func (*ProcessCloudEventsRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessCloudEventsRequest) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *ProcessCloudEventsRequest) GetBatch() *cloudevents.CloudEventBatch {
	if x != nil {
		return x.Batch
	}
	return nil
}

func (x *ProcessCloudEventsRequest) GetBadBatch() *ByteBatch {
	if x != nil {
		return x.BadBatch
	}
	return nil
}

type ProcessCloudEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProcessCloudEventsResponse) Reset() {
	*x = ProcessCloudEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessCloudEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessCloudEventsResponse) ProtoMessage() {}

func (x *ProcessCloudEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessCloudEventsResponse.ProtoReflect.Descriptor instead.
func (*ProcessCloudEventsResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescGZIP(), []int{2}
}

var File_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto protoreflect.FileDescriptor

var file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDesc = []byte{
	0x0a, 0x37, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x2f,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x70, 0x6b, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x1a, 0x27, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x27, 0x0a, 0x09, 0x42, 0x79, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1a,
	0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c,
	0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0xb6, 0x01, 0x0a, 0x19, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x12, 0x38, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x12, 0x45, 0x0a, 0x09,
	0x62, 0x61, 0x64, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x28, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x2e,
	0x42, 0x79, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x08, 0x62, 0x61, 0x64, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x22, 0x1c, 0x0a, 0x1a, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x43, 0x6c,
	0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0xa3, 0x01, 0x0a, 0x13, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x8b, 0x01, 0x0a, 0x12, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x12, 0x38, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e, 0x70, 0x6b, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xd8, 0x01, 0x0a, 0x2b, 0x63, 0x6f, 0x6d, 0x2e,
	0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x62, 0x65, 0x6e, 0x74,
	0x68, 0x6f, 0x73, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x42, 0x13, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x50, 0x01, 0x5a, 0x68,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66,
	0x79, 0x2d, 0x62, 0x75, 0x6e, 0x6e, 0x79, 0x2f, 0x62, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x73, 0x2d,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2d, 0x66, 0x6c, 0x75, 0x66,
	0x66, 0x79, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x6f, 0x72, 0x3b, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0xaa, 0x02, 0x27, 0x46, 0x6c, 0x75, 0x66, 0x66,
	0x79, 0x42, 0x75, 0x6e, 0x6e, 0x79, 0x2e, 0x42, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x73, 0x2e, 0x43,
	0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescOnce sync.Once
	file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescData = file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDesc
)

func file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescGZIP() []byte {
	file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescOnce.Do(func() {
		file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescData)
	})
	return file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDescData
}

var file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_goTypes = []interface{}{
	(*ByteBatch)(nil),                   // 0: pkg.proto.cloudeventprocessor.ByteBatch
	(*ProcessCloudEventsRequest)(nil),   // 1: pkg.proto.cloudeventprocessor.ProcessCloudEventsRequest
	(*ProcessCloudEventsResponse)(nil),  // 2: pkg.proto.cloudeventprocessor.ProcessCloudEventsResponse
	(*cloudevents.CloudEventBatch)(nil), // 3: io.cloudevents.v1.CloudEventBatch
}
var file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_depIdxs = []int32{
	3, // 0: pkg.proto.cloudeventprocessor.ProcessCloudEventsRequest.batch:type_name -> io.cloudevents.v1.CloudEventBatch
	0, // 1: pkg.proto.cloudeventprocessor.ProcessCloudEventsRequest.bad_batch:type_name -> pkg.proto.cloudeventprocessor.ByteBatch
	1, // 2: pkg.proto.cloudeventprocessor.CloudEventProcessor.ProcessCloudEvents:input_type -> pkg.proto.cloudeventprocessor.ProcessCloudEventsRequest
	2, // 3: pkg.proto.cloudeventprocessor.CloudEventProcessor.ProcessCloudEvents:output_type -> pkg.proto.cloudeventprocessor.ProcessCloudEventsResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_init() }
func file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_init() {
	if File_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ByteBatch); i {
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
		file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessCloudEventsRequest); i {
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
		file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessCloudEventsResponse); i {
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
			RawDescriptor: file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_goTypes,
		DependencyIndexes: file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_depIdxs,
		MessageInfos:      file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_msgTypes,
	}.Build()
	File_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto = out.File
	file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_rawDesc = nil
	file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_goTypes = nil
	file_pkg_proto_cloudeventprocessor_cloudeventprocessor_proto_depIdxs = nil
}
