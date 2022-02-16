// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AgentServiceClient is the client API for AgentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentServiceClient interface {
	Do(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	CreateService(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*CreateRequestResponse, error)
	Stream(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (AgentService_StreamClient, error)
}

type agentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentServiceClient(cc grpc.ClientConnInterface) AgentServiceClient {
	return &agentServiceClient{cc}
}

func (c *agentServiceClient) Do(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.AgentService/Do", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) CreateService(ctx context.Context, in *CreateServiceRequest, opts ...grpc.CallOption) (*CreateRequestResponse, error) {
	out := new(CreateRequestResponse)
	err := c.cc.Invoke(ctx, "/proto.AgentService/CreateService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentServiceClient) Stream(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (AgentService_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &AgentService_ServiceDesc.Streams[0], "/proto.AgentService/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentServiceStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AgentService_StreamClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type agentServiceStreamClient struct {
	grpc.ClientStream
}

func (x *agentServiceStreamClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AgentServiceServer is the server API for AgentService service.
// All implementations must embed UnimplementedAgentServiceServer
// for forward compatibility
type AgentServiceServer interface {
	Do(context.Context, *empty.Empty) (*empty.Empty, error)
	CreateService(context.Context, *CreateServiceRequest) (*CreateRequestResponse, error)
	Stream(*empty.Empty, AgentService_StreamServer) error
	mustEmbedUnimplementedAgentServiceServer()
}

// UnimplementedAgentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServiceServer struct {
}

func (UnimplementedAgentServiceServer) Do(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Do not implemented")
}
func (UnimplementedAgentServiceServer) CreateService(context.Context, *CreateServiceRequest) (*CreateRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateService not implemented")
}
func (UnimplementedAgentServiceServer) Stream(*empty.Empty, AgentService_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedAgentServiceServer) mustEmbedUnimplementedAgentServiceServer() {}

// UnsafeAgentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServiceServer will
// result in compilation errors.
type UnsafeAgentServiceServer interface {
	mustEmbedUnimplementedAgentServiceServer()
}

func RegisterAgentServiceServer(s grpc.ServiceRegistrar, srv AgentServiceServer) {
	s.RegisterService(&AgentService_ServiceDesc, srv)
}

func _AgentService_Do_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).Do(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AgentService/Do",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).Do(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_CreateService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServiceServer).CreateService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.AgentService/CreateService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServiceServer).CreateService(ctx, req.(*CreateServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AgentServiceServer).Stream(m, &agentServiceStreamServer{stream})
}

type AgentService_StreamServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

type agentServiceStreamServer struct {
	grpc.ServerStream
}

func (x *agentServiceStreamServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// AgentService_ServiceDesc is the grpc.ServiceDesc for AgentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Do",
			Handler:    _AgentService_Do_Handler,
		},
		{
			MethodName: "CreateService",
			Handler:    _AgentService_CreateService_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _AgentService_Stream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/proto/agent.proto",
}