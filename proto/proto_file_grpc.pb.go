// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: proto_file.proto

package proto_file_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	OrchAgent_AgentOrchGet_FullMethodName  = "/proto.OrchAgent/AgentOrchGet"
	OrchAgent_AgentOrchPost_FullMethodName = "/proto.OrchAgent/AgentOrchPost"
)

// OrchAgentClient is the client API for OrchAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrchAgentClient interface {
	AgentOrchGet(ctx context.Context, in *AgentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Task], error)
	AgentOrchPost(ctx context.Context, in *ResponseOfSecondServer, opts ...grpc.CallOption) (*OrchResponse, error)
}

type orchAgentClient struct {
	cc grpc.ClientConnInterface
}

func NewOrchAgentClient(cc grpc.ClientConnInterface) OrchAgentClient {
	return &orchAgentClient{cc}
}

func (c *orchAgentClient) AgentOrchGet(ctx context.Context, in *AgentRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Task], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &OrchAgent_ServiceDesc.Streams[0], OrchAgent_AgentOrchGet_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[AgentRequest, Task]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type OrchAgent_AgentOrchGetClient = grpc.ServerStreamingClient[Task]

func (c *orchAgentClient) AgentOrchPost(ctx context.Context, in *ResponseOfSecondServer, opts ...grpc.CallOption) (*OrchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrchResponse)
	err := c.cc.Invoke(ctx, OrchAgent_AgentOrchPost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrchAgentServer is the server API for OrchAgent service.
// All implementations must embed UnimplementedOrchAgentServer
// for forward compatibility.
type OrchAgentServer interface {
	AgentOrchGet(*AgentRequest, grpc.ServerStreamingServer[Task]) error
	AgentOrchPost(context.Context, *ResponseOfSecondServer) (*OrchResponse, error)
	mustEmbedUnimplementedOrchAgentServer()
}

// UnimplementedOrchAgentServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOrchAgentServer struct{}

func (UnimplementedOrchAgentServer) AgentOrchGet(*AgentRequest, grpc.ServerStreamingServer[Task]) error {
	return status.Errorf(codes.Unimplemented, "method AgentOrchGet not implemented")
}
func (UnimplementedOrchAgentServer) AgentOrchPost(context.Context, *ResponseOfSecondServer) (*OrchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AgentOrchPost not implemented")
}
func (UnimplementedOrchAgentServer) mustEmbedUnimplementedOrchAgentServer() {}
func (UnimplementedOrchAgentServer) testEmbeddedByValue()                   {}

// UnsafeOrchAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrchAgentServer will
// result in compilation errors.
type UnsafeOrchAgentServer interface {
	mustEmbedUnimplementedOrchAgentServer()
}

func RegisterOrchAgentServer(s grpc.ServiceRegistrar, srv OrchAgentServer) {
	// If the following call pancis, it indicates UnimplementedOrchAgentServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OrchAgent_ServiceDesc, srv)
}

func _OrchAgent_AgentOrchGet_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AgentRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OrchAgentServer).AgentOrchGet(m, &grpc.GenericServerStream[AgentRequest, Task]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type OrchAgent_AgentOrchGetServer = grpc.ServerStreamingServer[Task]

func _OrchAgent_AgentOrchPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseOfSecondServer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchAgentServer).AgentOrchPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrchAgent_AgentOrchPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchAgentServer).AgentOrchPost(ctx, req.(*ResponseOfSecondServer))
	}
	return interceptor(ctx, in, info, handler)
}

// OrchAgent_ServiceDesc is the grpc.ServiceDesc for OrchAgent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrchAgent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.OrchAgent",
	HandlerType: (*OrchAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AgentOrchPost",
			Handler:    _OrchAgent_AgentOrchPost_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AgentOrchGet",
			Handler:       _OrchAgent_AgentOrchGet_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto_file.proto",
}
