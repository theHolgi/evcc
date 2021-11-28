// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StreamServiceClient is the client API for StreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamServiceClient interface {
	FetchResponse(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_FetchResponseClient, error)
	SendUpdate(ctx context.Context, in *Request, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type streamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamServiceClient(cc grpc.ClientConnInterface) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) FetchResponse(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_FetchResponseClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[0], "/StreamService/FetchResponse", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceFetchResponseClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamService_FetchResponseClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type streamServiceFetchResponseClient struct {
	grpc.ClientStream
}

func (x *streamServiceFetchResponseClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamServiceClient) SendUpdate(ctx context.Context, in *Request, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/StreamService/SendUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StreamServiceServer is the server API for StreamService service.
// All implementations must embed UnimplementedStreamServiceServer
// for forward compatibility
type StreamServiceServer interface {
	FetchResponse(*Request, StreamService_FetchResponseServer) error
	SendUpdate(context.Context, *Request) (*emptypb.Empty, error)
	mustEmbedUnimplementedStreamServiceServer()
}

// UnimplementedStreamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamServiceServer struct {
}

func (UnimplementedStreamServiceServer) FetchResponse(*Request, StreamService_FetchResponseServer) error {
	return status.Errorf(codes.Unimplemented, "method FetchResponse not implemented")
}
func (UnimplementedStreamServiceServer) SendUpdate(context.Context, *Request) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUpdate not implemented")
}
func (UnimplementedStreamServiceServer) mustEmbedUnimplementedStreamServiceServer() {}

// UnsafeStreamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServiceServer will
// result in compilation errors.
type UnsafeStreamServiceServer interface {
	mustEmbedUnimplementedStreamServiceServer()
}

func RegisterStreamServiceServer(s grpc.ServiceRegistrar, srv StreamServiceServer) {
	s.RegisterService(&StreamService_ServiceDesc, srv)
}

func _StreamService_FetchResponse_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).FetchResponse(m, &streamServiceFetchResponseServer{stream})
}

type StreamService_FetchResponseServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type streamServiceFetchResponseServer struct {
	grpc.ServerStream
}

func (x *streamServiceFetchResponseServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func _StreamService_SendUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamServiceServer).SendUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/StreamService/SendUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamServiceServer).SendUpdate(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// StreamService_ServiceDesc is the grpc.ServiceDesc for StreamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendUpdate",
			Handler:    _StreamService_SendUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FetchResponse",
			Handler:       _StreamService_FetchResponse_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/connect.proto",
}