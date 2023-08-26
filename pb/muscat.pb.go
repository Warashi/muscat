// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: dev/warashi/muscat/v1/muscat.proto

package pb

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

type OpenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uri string `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
}

func (x *OpenRequest) Reset() {
	*x = OpenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenRequest) ProtoMessage() {}

func (x *OpenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenRequest.ProtoReflect.Descriptor instead.
func (*OpenRequest) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{0}
}

func (x *OpenRequest) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

type OpenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OpenResponse) Reset() {
	*x = OpenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenResponse) ProtoMessage() {}

func (x *OpenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenResponse.ProtoReflect.Descriptor instead.
func (*OpenResponse) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{1}
}

type CopyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body []byte `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *CopyRequest) Reset() {
	*x = CopyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CopyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CopyRequest) ProtoMessage() {}

func (x *CopyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CopyRequest.ProtoReflect.Descriptor instead.
func (*CopyRequest) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{2}
}

func (x *CopyRequest) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type CopyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CopyResponse) Reset() {
	*x = CopyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CopyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CopyResponse) ProtoMessage() {}

func (x *CopyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CopyResponse.ProtoReflect.Descriptor instead.
func (*CopyResponse) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{3}
}

type PasteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PasteRequest) Reset() {
	*x = PasteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PasteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasteRequest) ProtoMessage() {}

func (x *PasteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasteRequest.ProtoReflect.Descriptor instead.
func (*PasteRequest) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{4}
}

type PasteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body []byte `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *PasteResponse) Reset() {
	*x = PasteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PasteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasteResponse) ProtoMessage() {}

func (x *PasteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasteResponse.ProtoReflect.Descriptor instead.
func (*PasteResponse) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{5}
}

func (x *PasteResponse) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type HealthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthRequest) Reset() {
	*x = HealthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthRequest) ProtoMessage() {}

func (x *HealthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthRequest.ProtoReflect.Descriptor instead.
func (*HealthRequest) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{6}
}

type HealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pid int64 `protobuf:"varint,1,opt,name=pid,proto3" json:"pid,omitempty"`
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthResponse.ProtoReflect.Descriptor instead.
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{7}
}

func (x *HealthResponse) GetPid() int64 {
	if x != nil {
		return x.Pid
	}
	return 0
}

type GetInputMethodRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetInputMethodRequest) Reset() {
	*x = GetInputMethodRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInputMethodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInputMethodRequest) ProtoMessage() {}

func (x *GetInputMethodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInputMethodRequest.ProtoReflect.Descriptor instead.
func (*GetInputMethodRequest) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{8}
}

type GetInputMethodResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetInputMethodResponse) Reset() {
	*x = GetInputMethodResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInputMethodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInputMethodResponse) ProtoMessage() {}

func (x *GetInputMethodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInputMethodResponse.ProtoReflect.Descriptor instead.
func (*GetInputMethodResponse) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{9}
}

func (x *GetInputMethodResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type SetInputMethodRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SetInputMethodRequest) Reset() {
	*x = SetInputMethodRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetInputMethodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetInputMethodRequest) ProtoMessage() {}

func (x *SetInputMethodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetInputMethodRequest.ProtoReflect.Descriptor instead.
func (*SetInputMethodRequest) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{10}
}

func (x *SetInputMethodRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type SetInputMethodResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Before string `protobuf:"bytes,1,opt,name=before,proto3" json:"before,omitempty"`
}

func (x *SetInputMethodResponse) Reset() {
	*x = SetInputMethodResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetInputMethodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetInputMethodResponse) ProtoMessage() {}

func (x *SetInputMethodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetInputMethodResponse.ProtoReflect.Descriptor instead.
func (*SetInputMethodResponse) Descriptor() ([]byte, []int) {
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP(), []int{11}
}

func (x *SetInputMethodResponse) GetBefore() string {
	if x != nil {
		return x.Before
	}
	return ""
}

var File_dev_warashi_muscat_v1_muscat_proto protoreflect.FileDescriptor

var file_dev_warashi_muscat_v1_muscat_proto_rawDesc = []byte{
	0x0a, 0x22, 0x64, 0x65, 0x76, 0x2f, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2f, 0x6d, 0x75,
	0x73, 0x63, 0x61, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x64, 0x65, 0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68,
	0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x76, 0x31, 0x22, 0x1f, 0x0a, 0x0b, 0x4f,
	0x70, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x22, 0x0e, 0x0a, 0x0c,
	0x4f, 0x70, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x21, 0x0a, 0x0b,
	0x43, 0x6f, 0x70, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22,
	0x0e, 0x0a, 0x0c, 0x43, 0x6f, 0x70, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x0e, 0x0a, 0x0c, 0x50, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x23, 0x0a, 0x0d, 0x50, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x22, 0x0f, 0x0a, 0x0d, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x22, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x70, 0x69, 0x64, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65, 0x74,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x28, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x27, 0x0a, 0x15,
	0x53, 0x65, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x30, 0x0a, 0x16, 0x53, 0x65, 0x74, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x32, 0xbe, 0x04, 0x0a, 0x0d, 0x4d, 0x75, 0x73, 0x63,
	0x61, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x06, 0x48, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x12, 0x24, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68,
	0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x64, 0x65, 0x76, 0x2e,
	0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x4f, 0x0a, 0x04, 0x4f, 0x70, 0x65, 0x6e, 0x12, 0x22, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x77,
	0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x4f, 0x70, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x64,
	0x65, 0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x70, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x51, 0x0a, 0x04, 0x43, 0x6f, 0x70, 0x79, 0x12, 0x22, 0x2e, 0x64, 0x65, 0x76, 0x2e,
	0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6f, 0x70, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63,
	0x61, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x70, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x28, 0x01, 0x12, 0x54, 0x0a, 0x05, 0x50, 0x61, 0x73, 0x74, 0x65, 0x12, 0x23, 0x2e,
	0x64, 0x65, 0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63,
	0x61, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x73, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x24, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69,
	0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x73, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x6d, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x2c, 0x2e, 0x64,
	0x65, 0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x64, 0x65, 0x76,
	0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6d, 0x0a, 0x0e, 0x53, 0x65, 0x74,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x2c, 0x2e, 0x64, 0x65,
	0x76, 0x2e, 0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x64, 0x65, 0x76, 0x2e,
	0x77, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2e, 0x6d, 0x75, 0x73, 0x63, 0x61, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x65, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1e, 0x5a, 0x1c, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x57, 0x61, 0x72, 0x61, 0x73, 0x68, 0x69, 0x2f, 0x6d,
	0x75, 0x73, 0x63, 0x61, 0x74, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dev_warashi_muscat_v1_muscat_proto_rawDescOnce sync.Once
	file_dev_warashi_muscat_v1_muscat_proto_rawDescData = file_dev_warashi_muscat_v1_muscat_proto_rawDesc
)

func file_dev_warashi_muscat_v1_muscat_proto_rawDescGZIP() []byte {
	file_dev_warashi_muscat_v1_muscat_proto_rawDescOnce.Do(func() {
		file_dev_warashi_muscat_v1_muscat_proto_rawDescData = protoimpl.X.CompressGZIP(file_dev_warashi_muscat_v1_muscat_proto_rawDescData)
	})
	return file_dev_warashi_muscat_v1_muscat_proto_rawDescData
}

var file_dev_warashi_muscat_v1_muscat_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_dev_warashi_muscat_v1_muscat_proto_goTypes = []interface{}{
	(*OpenRequest)(nil),            // 0: dev.warashi.muscat.v1.OpenRequest
	(*OpenResponse)(nil),           // 1: dev.warashi.muscat.v1.OpenResponse
	(*CopyRequest)(nil),            // 2: dev.warashi.muscat.v1.CopyRequest
	(*CopyResponse)(nil),           // 3: dev.warashi.muscat.v1.CopyResponse
	(*PasteRequest)(nil),           // 4: dev.warashi.muscat.v1.PasteRequest
	(*PasteResponse)(nil),          // 5: dev.warashi.muscat.v1.PasteResponse
	(*HealthRequest)(nil),          // 6: dev.warashi.muscat.v1.HealthRequest
	(*HealthResponse)(nil),         // 7: dev.warashi.muscat.v1.HealthResponse
	(*GetInputMethodRequest)(nil),  // 8: dev.warashi.muscat.v1.GetInputMethodRequest
	(*GetInputMethodResponse)(nil), // 9: dev.warashi.muscat.v1.GetInputMethodResponse
	(*SetInputMethodRequest)(nil),  // 10: dev.warashi.muscat.v1.SetInputMethodRequest
	(*SetInputMethodResponse)(nil), // 11: dev.warashi.muscat.v1.SetInputMethodResponse
}
var file_dev_warashi_muscat_v1_muscat_proto_depIdxs = []int32{
	6,  // 0: dev.warashi.muscat.v1.MuscatService.Health:input_type -> dev.warashi.muscat.v1.HealthRequest
	0,  // 1: dev.warashi.muscat.v1.MuscatService.Open:input_type -> dev.warashi.muscat.v1.OpenRequest
	2,  // 2: dev.warashi.muscat.v1.MuscatService.Copy:input_type -> dev.warashi.muscat.v1.CopyRequest
	4,  // 3: dev.warashi.muscat.v1.MuscatService.Paste:input_type -> dev.warashi.muscat.v1.PasteRequest
	8,  // 4: dev.warashi.muscat.v1.MuscatService.GetInputMethod:input_type -> dev.warashi.muscat.v1.GetInputMethodRequest
	10, // 5: dev.warashi.muscat.v1.MuscatService.SetInputMethod:input_type -> dev.warashi.muscat.v1.SetInputMethodRequest
	7,  // 6: dev.warashi.muscat.v1.MuscatService.Health:output_type -> dev.warashi.muscat.v1.HealthResponse
	1,  // 7: dev.warashi.muscat.v1.MuscatService.Open:output_type -> dev.warashi.muscat.v1.OpenResponse
	3,  // 8: dev.warashi.muscat.v1.MuscatService.Copy:output_type -> dev.warashi.muscat.v1.CopyResponse
	5,  // 9: dev.warashi.muscat.v1.MuscatService.Paste:output_type -> dev.warashi.muscat.v1.PasteResponse
	9,  // 10: dev.warashi.muscat.v1.MuscatService.GetInputMethod:output_type -> dev.warashi.muscat.v1.GetInputMethodResponse
	11, // 11: dev.warashi.muscat.v1.MuscatService.SetInputMethod:output_type -> dev.warashi.muscat.v1.SetInputMethodResponse
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_dev_warashi_muscat_v1_muscat_proto_init() }
func file_dev_warashi_muscat_v1_muscat_proto_init() {
	if File_dev_warashi_muscat_v1_muscat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpenRequest); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpenResponse); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CopyRequest); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CopyResponse); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PasteRequest); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PasteResponse); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthRequest); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthResponse); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInputMethodRequest); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInputMethodResponse); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetInputMethodRequest); i {
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
		file_dev_warashi_muscat_v1_muscat_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetInputMethodResponse); i {
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
			RawDescriptor: file_dev_warashi_muscat_v1_muscat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dev_warashi_muscat_v1_muscat_proto_goTypes,
		DependencyIndexes: file_dev_warashi_muscat_v1_muscat_proto_depIdxs,
		MessageInfos:      file_dev_warashi_muscat_v1_muscat_proto_msgTypes,
	}.Build()
	File_dev_warashi_muscat_v1_muscat_proto = out.File
	file_dev_warashi_muscat_v1_muscat_proto_rawDesc = nil
	file_dev_warashi_muscat_v1_muscat_proto_goTypes = nil
	file_dev_warashi_muscat_v1_muscat_proto_depIdxs = nil
}
