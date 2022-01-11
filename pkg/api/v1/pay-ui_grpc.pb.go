// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// PayUIClient is the client API for PayUI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PayUIClient interface {
	// ListEngineSpecs returns a list of Message Engine(s) that can be started through the UI.
	ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (PayUI_ListEngineSpecsClient, error)
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error)
}

type payUIClient struct {
	cc grpc.ClientConnInterface
}

func NewPayUIClient(cc grpc.ClientConnInterface) PayUIClient {
	return &payUIClient{cc}
}

func (c *payUIClient) ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (PayUI_ListEngineSpecsClient, error) {
	stream, err := c.cc.NewStream(ctx, &PayUI_ServiceDesc.Streams[0], "/v1.PayUI/ListEngineSpecs", opts...)
	if err != nil {
		return nil, err
	}
	x := &payUIListEngineSpecsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PayUI_ListEngineSpecsClient interface {
	Recv() (*ListEngineSpecsResponse, error)
	grpc.ClientStream
}

type payUIListEngineSpecsClient struct {
	grpc.ClientStream
}

func (x *payUIListEngineSpecsClient) Recv() (*ListEngineSpecsResponse, error) {
	m := new(ListEngineSpecsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *payUIClient) IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error) {
	out := new(IsReadOnlyResponse)
	err := c.cc.Invoke(ctx, "/v1.PayUI/IsReadOnly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PayUIServer is the server API for PayUI service.
// All implementations must embed UnimplementedPayUIServer
// for forward compatibility
type PayUIServer interface {
	// ListEngineSpecs returns a list of Message Engine(s) that can be started through the UI.
	ListEngineSpecs(*ListEngineSpecsRequest, PayUI_ListEngineSpecsServer) error
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error)
	mustEmbedUnimplementedPayUIServer()
}

// UnimplementedPayUIServer must be embedded to have forward compatible implementations.
type UnimplementedPayUIServer struct {
}

func (UnimplementedPayUIServer) ListEngineSpecs(*ListEngineSpecsRequest, PayUI_ListEngineSpecsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListEngineSpecs not implemented")
}
func (UnimplementedPayUIServer) IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReadOnly not implemented")
}
func (UnimplementedPayUIServer) mustEmbedUnimplementedPayUIServer() {}

// UnsafePayUIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PayUIServer will
// result in compilation errors.
type UnsafePayUIServer interface {
	mustEmbedUnimplementedPayUIServer()
}

func RegisterPayUIServer(s grpc.ServiceRegistrar, srv PayUIServer) {
	s.RegisterService(&PayUI_ServiceDesc, srv)
}

func _PayUI_ListEngineSpecs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEngineSpecsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PayUIServer).ListEngineSpecs(m, &payUIListEngineSpecsServer{stream})
}

type PayUI_ListEngineSpecsServer interface {
	Send(*ListEngineSpecsResponse) error
	grpc.ServerStream
}

type payUIListEngineSpecsServer struct {
	grpc.ServerStream
}

func (x *payUIListEngineSpecsServer) Send(m *ListEngineSpecsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PayUI_IsReadOnly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadOnlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PayUIServer).IsReadOnly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PayUI/IsReadOnly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PayUIServer).IsReadOnly(ctx, req.(*IsReadOnlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PayUI_ServiceDesc is the grpc.ServiceDesc for PayUI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PayUI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.PayUI",
	HandlerType: (*PayUIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReadOnly",
			Handler:    _PayUI_IsReadOnly_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListEngineSpecs",
			Handler:       _PayUI_ListEngineSpecs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pay-ui.proto",
}
