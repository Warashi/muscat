// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: dev/warashi/muscat/v1/muscat.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OpenRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Uri         *string                `protobuf:"bytes,1,opt,name=uri"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *OpenRequest) Reset() {
	*x = OpenRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OpenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenRequest) ProtoMessage() {}

func (x *OpenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *OpenRequest) GetUri() string {
	if x != nil {
		if x.xxx_hidden_Uri != nil {
			return *x.xxx_hidden_Uri
		}
		return ""
	}
	return ""
}

func (x *OpenRequest) SetUri(v string) {
	x.xxx_hidden_Uri = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *OpenRequest) HasUri() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *OpenRequest) ClearUri() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Uri = nil
}

type OpenRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Uri *string
}

func (b0 OpenRequest_builder) Build() *OpenRequest {
	m0 := &OpenRequest{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Uri != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Uri = b.Uri
	}
	return m0
}

type OpenResponse struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OpenResponse) Reset() {
	*x = OpenResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OpenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenResponse) ProtoMessage() {}

func (x *OpenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type OpenResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 OpenResponse_builder) Build() *OpenResponse {
	m0 := &OpenResponse{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type CopyRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Body        []byte                 `protobuf:"bytes,1,opt,name=body"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *CopyRequest) Reset() {
	*x = CopyRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CopyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CopyRequest) ProtoMessage() {}

func (x *CopyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *CopyRequest) GetBody() []byte {
	if x != nil {
		return x.xxx_hidden_Body
	}
	return nil
}

func (x *CopyRequest) SetBody(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Body = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *CopyRequest) HasBody() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *CopyRequest) ClearBody() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Body = nil
}

type CopyRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Body []byte
}

func (b0 CopyRequest_builder) Build() *CopyRequest {
	m0 := &CopyRequest{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Body != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Body = b.Body
	}
	return m0
}

type CopyResponse struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CopyResponse) Reset() {
	*x = CopyResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CopyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CopyResponse) ProtoMessage() {}

func (x *CopyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type CopyResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 CopyResponse_builder) Build() *CopyResponse {
	m0 := &CopyResponse{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type PasteRequest struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PasteRequest) Reset() {
	*x = PasteRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PasteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasteRequest) ProtoMessage() {}

func (x *PasteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type PasteRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 PasteRequest_builder) Build() *PasteRequest {
	m0 := &PasteRequest{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type PasteResponse struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Body        []byte                 `protobuf:"bytes,1,opt,name=body"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *PasteResponse) Reset() {
	*x = PasteResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PasteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasteResponse) ProtoMessage() {}

func (x *PasteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PasteResponse) GetBody() []byte {
	if x != nil {
		return x.xxx_hidden_Body
	}
	return nil
}

func (x *PasteResponse) SetBody(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Body = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *PasteResponse) HasBody() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *PasteResponse) ClearBody() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Body = nil
}

type PasteResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Body []byte
}

func (b0 PasteResponse_builder) Build() *PasteResponse {
	m0 := &PasteResponse{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Body != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Body = b.Body
	}
	return m0
}

type HealthRequest struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthRequest) Reset() {
	*x = HealthRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthRequest) ProtoMessage() {}

func (x *HealthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type HealthRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 HealthRequest_builder) Build() *HealthRequest {
	m0 := &HealthRequest{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type HealthResponse struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Pid         int64                  `protobuf:"varint,1,opt,name=pid"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HealthResponse) GetPid() int64 {
	if x != nil {
		return x.xxx_hidden_Pid
	}
	return 0
}

func (x *HealthResponse) SetPid(v int64) {
	x.xxx_hidden_Pid = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *HealthResponse) HasPid() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *HealthResponse) ClearPid() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Pid = 0
}

type HealthResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Pid *int64
}

func (b0 HealthResponse_builder) Build() *HealthResponse {
	m0 := &HealthResponse{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Pid != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Pid = *b.Pid
	}
	return m0
}

type GetInputMethodRequest struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetInputMethodRequest) Reset() {
	*x = GetInputMethodRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInputMethodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInputMethodRequest) ProtoMessage() {}

func (x *GetInputMethodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type GetInputMethodRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 GetInputMethodRequest_builder) Build() *GetInputMethodRequest {
	m0 := &GetInputMethodRequest{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type GetInputMethodResponse struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Id          *string                `protobuf:"bytes,1,opt,name=id"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *GetInputMethodResponse) Reset() {
	*x = GetInputMethodResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInputMethodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInputMethodResponse) ProtoMessage() {}

func (x *GetInputMethodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *GetInputMethodResponse) GetId() string {
	if x != nil {
		if x.xxx_hidden_Id != nil {
			return *x.xxx_hidden_Id
		}
		return ""
	}
	return ""
}

func (x *GetInputMethodResponse) SetId(v string) {
	x.xxx_hidden_Id = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *GetInputMethodResponse) HasId() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *GetInputMethodResponse) ClearId() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Id = nil
}

type GetInputMethodResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Id *string
}

func (b0 GetInputMethodResponse_builder) Build() *GetInputMethodResponse {
	m0 := &GetInputMethodResponse{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Id != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Id = b.Id
	}
	return m0
}

type SetInputMethodRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Id          *string                `protobuf:"bytes,1,opt,name=id"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *SetInputMethodRequest) Reset() {
	*x = SetInputMethodRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetInputMethodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetInputMethodRequest) ProtoMessage() {}

func (x *SetInputMethodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *SetInputMethodRequest) GetId() string {
	if x != nil {
		if x.xxx_hidden_Id != nil {
			return *x.xxx_hidden_Id
		}
		return ""
	}
	return ""
}

func (x *SetInputMethodRequest) SetId(v string) {
	x.xxx_hidden_Id = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *SetInputMethodRequest) HasId() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *SetInputMethodRequest) ClearId() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Id = nil
}

type SetInputMethodRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Id *string
}

func (b0 SetInputMethodRequest_builder) Build() *SetInputMethodRequest {
	m0 := &SetInputMethodRequest{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Id != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Id = b.Id
	}
	return m0
}

type SetInputMethodResponse struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Before      *string                `protobuf:"bytes,1,opt,name=before"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *SetInputMethodResponse) Reset() {
	*x = SetInputMethodResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetInputMethodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetInputMethodResponse) ProtoMessage() {}

func (x *SetInputMethodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *SetInputMethodResponse) GetBefore() string {
	if x != nil {
		if x.xxx_hidden_Before != nil {
			return *x.xxx_hidden_Before
		}
		return ""
	}
	return ""
}

func (x *SetInputMethodResponse) SetBefore(v string) {
	x.xxx_hidden_Before = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *SetInputMethodResponse) HasBefore() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *SetInputMethodResponse) ClearBefore() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Before = nil
}

type SetInputMethodResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Before *string
}

func (b0 SetInputMethodResponse_builder) Build() *SetInputMethodResponse {
	m0 := &SetInputMethodResponse{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Before != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Before = b.Before
	}
	return m0
}

// PortForwardRequest.body is a body of transferred data.
type PortForwardRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Body        []byte                 `protobuf:"bytes,1,opt,name=body"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *PortForwardRequest) Reset() {
	*x = PortForwardRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PortForwardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortForwardRequest) ProtoMessage() {}

func (x *PortForwardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PortForwardRequest) GetBody() []byte {
	if x != nil {
		return x.xxx_hidden_Body
	}
	return nil
}

func (x *PortForwardRequest) SetBody(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Body = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *PortForwardRequest) HasBody() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *PortForwardRequest) ClearBody() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Body = nil
}

type PortForwardRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Body []byte
}

func (b0 PortForwardRequest_builder) Build() *PortForwardRequest {
	m0 := &PortForwardRequest{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Body != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Body = b.Body
	}
	return m0
}

// PortForwardResponse.body is a body of transferred data.
type PortForwardResponse struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Body        []byte                 `protobuf:"bytes,1,opt,name=body"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *PortForwardResponse) Reset() {
	*x = PortForwardResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PortForwardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortForwardResponse) ProtoMessage() {}

func (x *PortForwardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PortForwardResponse) GetBody() []byte {
	if x != nil {
		return x.xxx_hidden_Body
	}
	return nil
}

func (x *PortForwardResponse) SetBody(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Body = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *PortForwardResponse) HasBody() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *PortForwardResponse) ClearBody() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Body = nil
}

type PortForwardResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Body []byte
}

func (b0 PortForwardResponse_builder) Build() *PortForwardResponse {
	m0 := &PortForwardResponse{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Body != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Body = b.Body
	}
	return m0
}

type ExecRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Command     *string                `protobuf:"bytes,1,opt,name=command"`
	xxx_hidden_Args        []string               `protobuf:"bytes,2,rep,name=args"`
	xxx_hidden_Stdin       []byte                 `protobuf:"bytes,3,opt,name=stdin"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *ExecRequest) Reset() {
	*x = ExecRequest{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecRequest) ProtoMessage() {}

func (x *ExecRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ExecRequest) GetCommand() string {
	if x != nil {
		if x.xxx_hidden_Command != nil {
			return *x.xxx_hidden_Command
		}
		return ""
	}
	return ""
}

func (x *ExecRequest) GetArgs() []string {
	if x != nil {
		return x.xxx_hidden_Args
	}
	return nil
}

func (x *ExecRequest) GetStdin() []byte {
	if x != nil {
		return x.xxx_hidden_Stdin
	}
	return nil
}

func (x *ExecRequest) SetCommand(v string) {
	x.xxx_hidden_Command = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 3)
}

func (x *ExecRequest) SetArgs(v []string) {
	x.xxx_hidden_Args = v
}

func (x *ExecRequest) SetStdin(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Stdin = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 2, 3)
}

func (x *ExecRequest) HasCommand() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *ExecRequest) HasStdin() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 2)
}

func (x *ExecRequest) ClearCommand() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Command = nil
}

func (x *ExecRequest) ClearStdin() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 2)
	x.xxx_hidden_Stdin = nil
}

type ExecRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Command *string
	Args    []string
	Stdin   []byte
}

func (b0 ExecRequest_builder) Build() *ExecRequest {
	m0 := &ExecRequest{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Command != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 3)
		x.xxx_hidden_Command = b.Command
	}
	x.xxx_hidden_Args = b.Args
	if b.Stdin != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 2, 3)
		x.xxx_hidden_Stdin = b.Stdin
	}
	return m0
}

type ExecResponse struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Stdout      []byte                 `protobuf:"bytes,1,opt,name=stdout"`
	xxx_hidden_Stderr      []byte                 `protobuf:"bytes,2,opt,name=stderr"`
	xxx_hidden_ExitCode    int32                  `protobuf:"varint,3,opt,name=exit_code,json=exitCode"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *ExecResponse) Reset() {
	*x = ExecResponse{}
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecResponse) ProtoMessage() {}

func (x *ExecResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dev_warashi_muscat_v1_muscat_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ExecResponse) GetStdout() []byte {
	if x != nil {
		return x.xxx_hidden_Stdout
	}
	return nil
}

func (x *ExecResponse) GetStderr() []byte {
	if x != nil {
		return x.xxx_hidden_Stderr
	}
	return nil
}

func (x *ExecResponse) GetExitCode() int32 {
	if x != nil {
		return x.xxx_hidden_ExitCode
	}
	return 0
}

func (x *ExecResponse) SetStdout(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Stdout = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 3)
}

func (x *ExecResponse) SetStderr(v []byte) {
	if v == nil {
		v = []byte{}
	}
	x.xxx_hidden_Stderr = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 1, 3)
}

func (x *ExecResponse) SetExitCode(v int32) {
	x.xxx_hidden_ExitCode = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 2, 3)
}

func (x *ExecResponse) HasStdout() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *ExecResponse) HasStderr() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 1)
}

func (x *ExecResponse) HasExitCode() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 2)
}

func (x *ExecResponse) ClearStdout() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Stdout = nil
}

func (x *ExecResponse) ClearStderr() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 1)
	x.xxx_hidden_Stderr = nil
}

func (x *ExecResponse) ClearExitCode() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 2)
	x.xxx_hidden_ExitCode = 0
}

type ExecResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Stdout   []byte
	Stderr   []byte
	ExitCode *int32
}

func (b0 ExecResponse_builder) Build() *ExecResponse {
	m0 := &ExecResponse{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Stdout != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 3)
		x.xxx_hidden_Stdout = b.Stdout
	}
	if b.Stderr != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 1, 3)
		x.xxx_hidden_Stderr = b.Stderr
	}
	if b.ExitCode != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 2, 3)
		x.xxx_hidden_ExitCode = *b.ExitCode
	}
	return m0
}

var File_dev_warashi_muscat_v1_muscat_proto protoreflect.FileDescriptor

const file_dev_warashi_muscat_v1_muscat_proto_rawDesc = "" +
	"\n" +
	"\"dev/warashi/muscat/v1/muscat.proto\x12\x15dev.warashi.muscat.v1\"\x1f\n" +
	"\vOpenRequest\x12\x10\n" +
	"\x03uri\x18\x01 \x01(\tR\x03uri\"\x0e\n" +
	"\fOpenResponse\"!\n" +
	"\vCopyRequest\x12\x12\n" +
	"\x04body\x18\x01 \x01(\fR\x04body\"\x0e\n" +
	"\fCopyResponse\"\x0e\n" +
	"\fPasteRequest\"#\n" +
	"\rPasteResponse\x12\x12\n" +
	"\x04body\x18\x01 \x01(\fR\x04body\"\x0f\n" +
	"\rHealthRequest\"\"\n" +
	"\x0eHealthResponse\x12\x10\n" +
	"\x03pid\x18\x01 \x01(\x03R\x03pid\"\x17\n" +
	"\x15GetInputMethodRequest\"(\n" +
	"\x16GetInputMethodResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"'\n" +
	"\x15SetInputMethodRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"0\n" +
	"\x16SetInputMethodResponse\x12\x16\n" +
	"\x06before\x18\x01 \x01(\tR\x06before\"(\n" +
	"\x12PortForwardRequest\x12\x12\n" +
	"\x04body\x18\x01 \x01(\fR\x04body\")\n" +
	"\x13PortForwardResponse\x12\x12\n" +
	"\x04body\x18\x01 \x01(\fR\x04body\"Q\n" +
	"\vExecRequest\x12\x18\n" +
	"\acommand\x18\x01 \x01(\tR\acommand\x12\x12\n" +
	"\x04args\x18\x02 \x03(\tR\x04args\x12\x14\n" +
	"\x05stdin\x18\x03 \x01(\fR\x05stdin\"[\n" +
	"\fExecResponse\x12\x16\n" +
	"\x06stdout\x18\x01 \x01(\fR\x06stdout\x12\x16\n" +
	"\x06stderr\x18\x02 \x01(\fR\x06stderr\x12\x1b\n" +
	"\texit_code\x18\x03 \x01(\x05R\bexitCode2\xf9\x05\n" +
	"\rMuscatService\x12U\n" +
	"\x06Health\x12$.dev.warashi.muscat.v1.HealthRequest\x1a%.dev.warashi.muscat.v1.HealthResponse\x12O\n" +
	"\x04Open\x12\".dev.warashi.muscat.v1.OpenRequest\x1a#.dev.warashi.muscat.v1.OpenResponse\x12Q\n" +
	"\x04Copy\x12\".dev.warashi.muscat.v1.CopyRequest\x1a#.dev.warashi.muscat.v1.CopyResponse(\x01\x12T\n" +
	"\x05Paste\x12#.dev.warashi.muscat.v1.PasteRequest\x1a$.dev.warashi.muscat.v1.PasteResponse0\x01\x12m\n" +
	"\x0eGetInputMethod\x12,.dev.warashi.muscat.v1.GetInputMethodRequest\x1a-.dev.warashi.muscat.v1.GetInputMethodResponse\x12m\n" +
	"\x0eSetInputMethod\x12,.dev.warashi.muscat.v1.SetInputMethodRequest\x1a-.dev.warashi.muscat.v1.SetInputMethodResponse\x12h\n" +
	"\vPortForward\x12).dev.warashi.muscat.v1.PortForwardRequest\x1a*.dev.warashi.muscat.v1.PortForwardResponse(\x010\x01\x12O\n" +
	"\x04Exec\x12\".dev.warashi.muscat.v1.ExecRequest\x1a#.dev.warashi.muscat.v1.ExecResponseB!Z\x1fgithub.com/Warashi/muscat/v2/pbb\beditionsp\xe8\a"

var file_dev_warashi_muscat_v1_muscat_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_dev_warashi_muscat_v1_muscat_proto_goTypes = []any{
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
	(*PortForwardRequest)(nil),     // 12: dev.warashi.muscat.v1.PortForwardRequest
	(*PortForwardResponse)(nil),    // 13: dev.warashi.muscat.v1.PortForwardResponse
	(*ExecRequest)(nil),            // 14: dev.warashi.muscat.v1.ExecRequest
	(*ExecResponse)(nil),           // 15: dev.warashi.muscat.v1.ExecResponse
}
var file_dev_warashi_muscat_v1_muscat_proto_depIdxs = []int32{
	6,  // 0: dev.warashi.muscat.v1.MuscatService.Health:input_type -> dev.warashi.muscat.v1.HealthRequest
	0,  // 1: dev.warashi.muscat.v1.MuscatService.Open:input_type -> dev.warashi.muscat.v1.OpenRequest
	2,  // 2: dev.warashi.muscat.v1.MuscatService.Copy:input_type -> dev.warashi.muscat.v1.CopyRequest
	4,  // 3: dev.warashi.muscat.v1.MuscatService.Paste:input_type -> dev.warashi.muscat.v1.PasteRequest
	8,  // 4: dev.warashi.muscat.v1.MuscatService.GetInputMethod:input_type -> dev.warashi.muscat.v1.GetInputMethodRequest
	10, // 5: dev.warashi.muscat.v1.MuscatService.SetInputMethod:input_type -> dev.warashi.muscat.v1.SetInputMethodRequest
	12, // 6: dev.warashi.muscat.v1.MuscatService.PortForward:input_type -> dev.warashi.muscat.v1.PortForwardRequest
	14, // 7: dev.warashi.muscat.v1.MuscatService.Exec:input_type -> dev.warashi.muscat.v1.ExecRequest
	7,  // 8: dev.warashi.muscat.v1.MuscatService.Health:output_type -> dev.warashi.muscat.v1.HealthResponse
	1,  // 9: dev.warashi.muscat.v1.MuscatService.Open:output_type -> dev.warashi.muscat.v1.OpenResponse
	3,  // 10: dev.warashi.muscat.v1.MuscatService.Copy:output_type -> dev.warashi.muscat.v1.CopyResponse
	5,  // 11: dev.warashi.muscat.v1.MuscatService.Paste:output_type -> dev.warashi.muscat.v1.PasteResponse
	9,  // 12: dev.warashi.muscat.v1.MuscatService.GetInputMethod:output_type -> dev.warashi.muscat.v1.GetInputMethodResponse
	11, // 13: dev.warashi.muscat.v1.MuscatService.SetInputMethod:output_type -> dev.warashi.muscat.v1.SetInputMethodResponse
	13, // 14: dev.warashi.muscat.v1.MuscatService.PortForward:output_type -> dev.warashi.muscat.v1.PortForwardResponse
	15, // 15: dev.warashi.muscat.v1.MuscatService.Exec:output_type -> dev.warashi.muscat.v1.ExecResponse
	8,  // [8:16] is the sub-list for method output_type
	0,  // [0:8] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_dev_warashi_muscat_v1_muscat_proto_init() }
func file_dev_warashi_muscat_v1_muscat_proto_init() {
	if File_dev_warashi_muscat_v1_muscat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_dev_warashi_muscat_v1_muscat_proto_rawDesc), len(file_dev_warashi_muscat_v1_muscat_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dev_warashi_muscat_v1_muscat_proto_goTypes,
		DependencyIndexes: file_dev_warashi_muscat_v1_muscat_proto_depIdxs,
		MessageInfos:      file_dev_warashi_muscat_v1_muscat_proto_msgTypes,
	}.Build()
	File_dev_warashi_muscat_v1_muscat_proto = out.File
	file_dev_warashi_muscat_v1_muscat_proto_goTypes = nil
	file_dev_warashi_muscat_v1_muscat_proto_depIdxs = nil
}
