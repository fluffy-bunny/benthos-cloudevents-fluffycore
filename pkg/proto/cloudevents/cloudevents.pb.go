//*
// CloudEvent Protobuf Format
//
// - Required context attributes are explicitly represented.
// - Optional and Extension context attributes are carried in a map structure.
// - Data may be represented as binary, text, or protobuf messages.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.0
// source: pkg/proto/cloudevents/cloudevents.proto

package cloudevents

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CloudEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required Attributes
	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Source      string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"` // URI-reference
	SpecVersion string `protobuf:"bytes,3,opt,name=spec_version,json=specVersion,proto3" json:"spec_version,omitempty"`
	Type        string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	// Optional & Extension Attributes
	Attributes map[string]*CloudEvent_CloudEventAttributeValue `protobuf:"bytes,5,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// -- CloudEvent Data (Bytes, Text, or Proto)
	//
	// Types that are assignable to Data:
	//
	//	*CloudEvent_BinaryData
	//	*CloudEvent_TextData
	//	*CloudEvent_ProtoData
	Data isCloudEvent_Data `protobuf_oneof:"data"`
}

func (x *CloudEvent) Reset() {
	*x = CloudEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloudEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloudEvent) ProtoMessage() {}

func (x *CloudEvent) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloudEvent.ProtoReflect.Descriptor instead.
func (*CloudEvent) Descriptor() ([]byte, []int) {
	return file_pkg_proto_cloudevents_cloudevents_proto_rawDescGZIP(), []int{0}
}

func (x *CloudEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CloudEvent) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *CloudEvent) GetSpecVersion() string {
	if x != nil {
		return x.SpecVersion
	}
	return ""
}

func (x *CloudEvent) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *CloudEvent) GetAttributes() map[string]*CloudEvent_CloudEventAttributeValue {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (m *CloudEvent) GetData() isCloudEvent_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *CloudEvent) GetBinaryData() []byte {
	if x, ok := x.GetData().(*CloudEvent_BinaryData); ok {
		return x.BinaryData
	}
	return nil
}

func (x *CloudEvent) GetTextData() string {
	if x, ok := x.GetData().(*CloudEvent_TextData); ok {
		return x.TextData
	}
	return ""
}

func (x *CloudEvent) GetProtoData() *anypb.Any {
	if x, ok := x.GetData().(*CloudEvent_ProtoData); ok {
		return x.ProtoData
	}
	return nil
}

type isCloudEvent_Data interface {
	isCloudEvent_Data()
}

type CloudEvent_BinaryData struct {
	BinaryData []byte `protobuf:"bytes,6,opt,name=binary_data,json=binaryData,proto3,oneof"`
}

type CloudEvent_TextData struct {
	TextData string `protobuf:"bytes,7,opt,name=text_data,json=textData,proto3,oneof"`
}

type CloudEvent_ProtoData struct {
	ProtoData *anypb.Any `protobuf:"bytes,8,opt,name=proto_data,json=protoData,proto3,oneof"`
}

func (*CloudEvent_BinaryData) isCloudEvent_Data() {}

func (*CloudEvent_TextData) isCloudEvent_Data() {}

func (*CloudEvent_ProtoData) isCloudEvent_Data() {}

type CloudEventBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*CloudEvent `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *CloudEventBatch) Reset() {
	*x = CloudEventBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloudEventBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloudEventBatch) ProtoMessage() {}

func (x *CloudEventBatch) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloudEventBatch.ProtoReflect.Descriptor instead.
func (*CloudEventBatch) Descriptor() ([]byte, []int) {
	return file_pkg_proto_cloudevents_cloudevents_proto_rawDescGZIP(), []int{1}
}

func (x *CloudEventBatch) GetEvents() []*CloudEvent {
	if x != nil {
		return x.Events
	}
	return nil
}

type CloudEvent_CloudEventAttributeValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Attr:
	//
	//	*CloudEvent_CloudEventAttributeValue_CeBoolean
	//	*CloudEvent_CloudEventAttributeValue_CeInteger
	//	*CloudEvent_CloudEventAttributeValue_CeString
	//	*CloudEvent_CloudEventAttributeValue_CeBytes
	//	*CloudEvent_CloudEventAttributeValue_CeUri
	//	*CloudEvent_CloudEventAttributeValue_CeUriRef
	//	*CloudEvent_CloudEventAttributeValue_CeTimestamp
	Attr isCloudEvent_CloudEventAttributeValue_Attr `protobuf_oneof:"attr"`
}

func (x *CloudEvent_CloudEventAttributeValue) Reset() {
	*x = CloudEvent_CloudEventAttributeValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloudEvent_CloudEventAttributeValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloudEvent_CloudEventAttributeValue) ProtoMessage() {}

func (x *CloudEvent_CloudEventAttributeValue) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloudEvent_CloudEventAttributeValue.ProtoReflect.Descriptor instead.
func (*CloudEvent_CloudEventAttributeValue) Descriptor() ([]byte, []int) {
	return file_pkg_proto_cloudevents_cloudevents_proto_rawDescGZIP(), []int{0, 1}
}

func (m *CloudEvent_CloudEventAttributeValue) GetAttr() isCloudEvent_CloudEventAttributeValue_Attr {
	if m != nil {
		return m.Attr
	}
	return nil
}

func (x *CloudEvent_CloudEventAttributeValue) GetCeBoolean() bool {
	if x, ok := x.GetAttr().(*CloudEvent_CloudEventAttributeValue_CeBoolean); ok {
		return x.CeBoolean
	}
	return false
}

func (x *CloudEvent_CloudEventAttributeValue) GetCeInteger() int32 {
	if x, ok := x.GetAttr().(*CloudEvent_CloudEventAttributeValue_CeInteger); ok {
		return x.CeInteger
	}
	return 0
}

func (x *CloudEvent_CloudEventAttributeValue) GetCeString() string {
	if x, ok := x.GetAttr().(*CloudEvent_CloudEventAttributeValue_CeString); ok {
		return x.CeString
	}
	return ""
}

func (x *CloudEvent_CloudEventAttributeValue) GetCeBytes() []byte {
	if x, ok := x.GetAttr().(*CloudEvent_CloudEventAttributeValue_CeBytes); ok {
		return x.CeBytes
	}
	return nil
}

func (x *CloudEvent_CloudEventAttributeValue) GetCeUri() string {
	if x, ok := x.GetAttr().(*CloudEvent_CloudEventAttributeValue_CeUri); ok {
		return x.CeUri
	}
	return ""
}

func (x *CloudEvent_CloudEventAttributeValue) GetCeUriRef() string {
	if x, ok := x.GetAttr().(*CloudEvent_CloudEventAttributeValue_CeUriRef); ok {
		return x.CeUriRef
	}
	return ""
}

func (x *CloudEvent_CloudEventAttributeValue) GetCeTimestamp() *timestamppb.Timestamp {
	if x, ok := x.GetAttr().(*CloudEvent_CloudEventAttributeValue_CeTimestamp); ok {
		return x.CeTimestamp
	}
	return nil
}

type isCloudEvent_CloudEventAttributeValue_Attr interface {
	isCloudEvent_CloudEventAttributeValue_Attr()
}

type CloudEvent_CloudEventAttributeValue_CeBoolean struct {
	CeBoolean bool `protobuf:"varint,1,opt,name=ce_boolean,json=ceBoolean,proto3,oneof"`
}

type CloudEvent_CloudEventAttributeValue_CeInteger struct {
	CeInteger int32 `protobuf:"varint,2,opt,name=ce_integer,json=ceInteger,proto3,oneof"`
}

type CloudEvent_CloudEventAttributeValue_CeString struct {
	CeString string `protobuf:"bytes,3,opt,name=ce_string,json=ceString,proto3,oneof"`
}

type CloudEvent_CloudEventAttributeValue_CeBytes struct {
	CeBytes []byte `protobuf:"bytes,4,opt,name=ce_bytes,json=ceBytes,proto3,oneof"`
}

type CloudEvent_CloudEventAttributeValue_CeUri struct {
	CeUri string `protobuf:"bytes,5,opt,name=ce_uri,json=ceUri,proto3,oneof"`
}

type CloudEvent_CloudEventAttributeValue_CeUriRef struct {
	CeUriRef string `protobuf:"bytes,6,opt,name=ce_uri_ref,json=ceUriRef,proto3,oneof"`
}

type CloudEvent_CloudEventAttributeValue_CeTimestamp struct {
	CeTimestamp *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=ce_timestamp,json=ceTimestamp,proto3,oneof"`
}

func (*CloudEvent_CloudEventAttributeValue_CeBoolean) isCloudEvent_CloudEventAttributeValue_Attr() {}

func (*CloudEvent_CloudEventAttributeValue_CeInteger) isCloudEvent_CloudEventAttributeValue_Attr() {}

func (*CloudEvent_CloudEventAttributeValue_CeString) isCloudEvent_CloudEventAttributeValue_Attr() {}

func (*CloudEvent_CloudEventAttributeValue_CeBytes) isCloudEvent_CloudEventAttributeValue_Attr() {}

func (*CloudEvent_CloudEventAttributeValue_CeUri) isCloudEvent_CloudEventAttributeValue_Attr() {}

func (*CloudEvent_CloudEventAttributeValue_CeUriRef) isCloudEvent_CloudEventAttributeValue_Attr() {}

func (*CloudEvent_CloudEventAttributeValue_CeTimestamp) isCloudEvent_CloudEventAttributeValue_Attr() {
}

var File_pkg_proto_cloudevents_cloudevents_proto protoreflect.FileDescriptor

var file_pkg_proto_cloudevents_cloudevents_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x69, 0x6f, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x05, 0x0a, 0x0a, 0x43, 0x6c, 0x6f,
	0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x73, 0x70, 0x65, 0x63, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x70, 0x65, 0x63, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x69, 0x6f, 0x2e,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x0b, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x5f,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0a, 0x62, 0x69,
	0x6e, 0x61, 0x72, 0x79, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x09, 0x74, 0x65, 0x78, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x74,
	0x65, 0x78, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x35, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x48, 0x00, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x75,
	0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x4c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x36, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x9a, 0x02, 0x0a, 0x18, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x63, 0x65, 0x5f, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x63, 0x65, 0x42, 0x6f, 0x6f, 0x6c,
	0x65, 0x61, 0x6e, 0x12, 0x1f, 0x0a, 0x0a, 0x63, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x09, 0x63, 0x65, 0x49, 0x6e, 0x74,
	0x65, 0x67, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x09, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x63, 0x65, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x08, 0x63, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x07, 0x63, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x12, 0x17, 0x0a, 0x06, 0x63, 0x65, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x05, 0x63, 0x65, 0x55, 0x72, 0x69, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x65, 0x5f,
	0x75, 0x72, 0x69, 0x5f, 0x72, 0x65, 0x66, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x08, 0x63, 0x65, 0x55, 0x72, 0x69, 0x52, 0x65, 0x66, 0x12, 0x3f, 0x0a, 0x0c, 0x63, 0x65, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x0b, 0x63,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0x0a, 0x04, 0x61, 0x74,
	0x74, 0x72, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x48, 0x0a, 0x0f, 0x43, 0x6c,
	0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x35, 0x0a,
	0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x42, 0xc4, 0x01, 0x0a, 0x1e, 0x70, 0x6b, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x2d, 0x62, 0x75, 0x6e,
	0x6e, 0x79, 0x2f, 0x62, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x73, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2d, 0x66, 0x6c, 0x75, 0x66, 0x66, 0x79, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0xaa, 0x02, 0x1a, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x4e,
	0x61, 0x74, 0x69, 0x76, 0x65, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x17, 0x49, 0x6f, 0x5c, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0xea, 0x02,
	0x1a, 0x49, 0x6f, 0x3a, 0x3a, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x3a, 0x3a, 0x56, 0x31, 0x3a, 0x3a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_cloudevents_cloudevents_proto_rawDescOnce sync.Once
	file_pkg_proto_cloudevents_cloudevents_proto_rawDescData = file_pkg_proto_cloudevents_cloudevents_proto_rawDesc
)

func file_pkg_proto_cloudevents_cloudevents_proto_rawDescGZIP() []byte {
	file_pkg_proto_cloudevents_cloudevents_proto_rawDescOnce.Do(func() {
		file_pkg_proto_cloudevents_cloudevents_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_cloudevents_cloudevents_proto_rawDescData)
	})
	return file_pkg_proto_cloudevents_cloudevents_proto_rawDescData
}

var file_pkg_proto_cloudevents_cloudevents_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_proto_cloudevents_cloudevents_proto_goTypes = []interface{}{
	(*CloudEvent)(nil),      // 0: io.cloudevents.v1.CloudEvent
	(*CloudEventBatch)(nil), // 1: io.cloudevents.v1.CloudEventBatch
	nil,                     // 2: io.cloudevents.v1.CloudEvent.AttributesEntry
	(*CloudEvent_CloudEventAttributeValue)(nil), // 3: io.cloudevents.v1.CloudEvent.CloudEventAttributeValue
	(*anypb.Any)(nil),             // 4: google.protobuf.Any
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_pkg_proto_cloudevents_cloudevents_proto_depIdxs = []int32{
	2, // 0: io.cloudevents.v1.CloudEvent.attributes:type_name -> io.cloudevents.v1.CloudEvent.AttributesEntry
	4, // 1: io.cloudevents.v1.CloudEvent.proto_data:type_name -> google.protobuf.Any
	0, // 2: io.cloudevents.v1.CloudEventBatch.events:type_name -> io.cloudevents.v1.CloudEvent
	3, // 3: io.cloudevents.v1.CloudEvent.AttributesEntry.value:type_name -> io.cloudevents.v1.CloudEvent.CloudEventAttributeValue
	5, // 4: io.cloudevents.v1.CloudEvent.CloudEventAttributeValue.ce_timestamp:type_name -> google.protobuf.Timestamp
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pkg_proto_cloudevents_cloudevents_proto_init() }
func file_pkg_proto_cloudevents_cloudevents_proto_init() {
	if File_pkg_proto_cloudevents_cloudevents_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CloudEvent); i {
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
		file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CloudEventBatch); i {
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
		file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CloudEvent_CloudEventAttributeValue); i {
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
	file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*CloudEvent_BinaryData)(nil),
		(*CloudEvent_TextData)(nil),
		(*CloudEvent_ProtoData)(nil),
	}
	file_pkg_proto_cloudevents_cloudevents_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*CloudEvent_CloudEventAttributeValue_CeBoolean)(nil),
		(*CloudEvent_CloudEventAttributeValue_CeInteger)(nil),
		(*CloudEvent_CloudEventAttributeValue_CeString)(nil),
		(*CloudEvent_CloudEventAttributeValue_CeBytes)(nil),
		(*CloudEvent_CloudEventAttributeValue_CeUri)(nil),
		(*CloudEvent_CloudEventAttributeValue_CeUriRef)(nil),
		(*CloudEvent_CloudEventAttributeValue_CeTimestamp)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_proto_cloudevents_cloudevents_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_cloudevents_cloudevents_proto_goTypes,
		DependencyIndexes: file_pkg_proto_cloudevents_cloudevents_proto_depIdxs,
		MessageInfos:      file_pkg_proto_cloudevents_cloudevents_proto_msgTypes,
	}.Build()
	File_pkg_proto_cloudevents_cloudevents_proto = out.File
	file_pkg_proto_cloudevents_cloudevents_proto_rawDesc = nil
	file_pkg_proto_cloudevents_cloudevents_proto_goTypes = nil
	file_pkg_proto_cloudevents_cloudevents_proto_depIdxs = nil
}
