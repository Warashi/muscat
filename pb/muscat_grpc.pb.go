// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: muscat.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MuscatClient is the client API for Muscat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MuscatClient interface {
	Open(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenResponse, error)
	Copy(ctx context.Context, opts ...grpc.CallOption) (Muscat_CopyClient, error)
	Paste(ctx context.Context, in *PasteRequest, opts ...grpc.CallOption) (Muscat_PasteClient, error)
	Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error)
}

type muscatClient struct {
	cc grpc.ClientConnInterface
}

func NewMuscatClient(cc grpc.ClientConnInterface) MuscatClient {
	return &muscatClient{cc}
}

func (c *muscatClient) Open(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenResponse, error) {
	out := new(OpenResponse)
	err := c.cc.Invoke(ctx, "/dev.warashi.muscat.Muscat/Open", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *muscatClient) Copy(ctx context.Context, opts ...grpc.CallOption) (Muscat_CopyClient, error) {
	stream, err := c.cc.NewStream(ctx, &Muscat_ServiceDesc.Streams[0], "/dev.warashi.muscat.Muscat/Copy", opts...)
	if err != nil {
		return nil, err
	}
	x := &muscatCopyClient{stream}
	return x, nil
}

type Muscat_CopyClient interface {
	Send(*CopyRequest) error
	CloseAndRecv() (*CopyResponse, error)
	grpc.ClientStream
}

type muscatCopyClient struct {
	grpc.ClientStream
}

func (x *muscatCopyClient) Send(m *CopyRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *muscatCopyClient) CloseAndRecv() (*CopyResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(CopyResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *muscatClient) Paste(ctx context.Context, in *PasteRequest, opts ...grpc.CallOption) (Muscat_PasteClient, error) {
	stream, err := c.cc.NewStream(ctx, &Muscat_ServiceDesc.Streams[1], "/dev.warashi.muscat.Muscat/Paste", opts...)
	if err != nil {
		return nil, err
	}
	x := &muscatPasteClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Muscat_PasteClient interface {
	Recv() (*PasteResponse, error)
	grpc.ClientStream
}

type muscatPasteClient struct {
	grpc.ClientStream
}

func (x *muscatPasteClient) Recv() (*PasteResponse, error) {
	m := new(PasteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *muscatClient) Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, "/dev.warashi.muscat.Muscat/Health", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MuscatServer is the server API for Muscat service.
// All implementations must embed UnimplementedMuscatServer
// for forward compatibility
type MuscatServer interface {
	Open(context.Context, *OpenRequest) (*OpenResponse, error)
	Copy(Muscat_CopyServer) error
	Paste(*PasteRequest, Muscat_PasteServer) error
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	mustEmbedUnimplementedMuscatServer()
}

// UnimplementedMuscatServer must be embedded to have forward compatible implementations.
type UnimplementedMuscatServer struct {
}

func (UnimplementedMuscatServer) Open(context.Context, *OpenRequest) (*OpenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Open not implemented")
}
func (UnimplementedMuscatServer) Copy(Muscat_CopyServer) error {
	return status.Errorf(codes.Unimplemented, "method Copy not implemented")
}
func (UnimplementedMuscatServer) Paste(*PasteRequest, Muscat_PasteServer) error {
	return status.Errorf(codes.Unimplemented, "method Paste not implemented")
}
func (UnimplementedMuscatServer) Health(context.Context, *HealthRequest) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedMuscatServer) mustEmbedUnimplementedMuscatServer() {}

// UnsafeMuscatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MuscatServer will
// result in compilation errors.
type UnsafeMuscatServer interface {
	mustEmbedUnimplementedMuscatServer()
}

func RegisterMuscatServer(s grpc.ServiceRegistrar, srv MuscatServer) {
	s.RegisterService(&Muscat_ServiceDesc, srv)
}

func _Muscat_Open_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MuscatServer).Open(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.warashi.muscat.Muscat/Open",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MuscatServer).Open(ctx, req.(*OpenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Muscat_Copy_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MuscatServer).Copy(&muscatCopyServer{stream})
}

type Muscat_CopyServer interface {
	SendAndClose(*CopyResponse) error
	Recv() (*CopyRequest, error)
	grpc.ServerStream
}

type muscatCopyServer struct {
	grpc.ServerStream
}

func (x *muscatCopyServer) SendAndClose(m *CopyResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *muscatCopyServer) Recv() (*CopyRequest, error) {
	m := new(CopyRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Muscat_Paste_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PasteRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MuscatServer).Paste(m, &muscatPasteServer{stream})
}

type Muscat_PasteServer interface {
	Send(*PasteResponse) error
	grpc.ServerStream
}

type muscatPasteServer struct {
	grpc.ServerStream
}

func (x *muscatPasteServer) Send(m *PasteResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Muscat_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MuscatServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.warashi.muscat.Muscat/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MuscatServer).Health(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Muscat_ServiceDesc is the grpc.ServiceDesc for Muscat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Muscat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dev.warashi.muscat.Muscat",
	HandlerType: (*MuscatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Open",
			Handler:    _Muscat_Open_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _Muscat_Health_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Copy",
			Handler:       _Muscat_Copy_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Paste",
			Handler:       _Muscat_Paste_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "muscat.proto",
}
