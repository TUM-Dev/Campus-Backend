// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// CampusClient is the client API for Campus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CampusClient interface {
	GetTopNews(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTopNewsReply, error)
	GetNewsSources(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Campus_GetNewsSourcesClient, error)
}

type campusClient struct {
	cc grpc.ClientConnInterface
}

func NewCampusClient(cc grpc.ClientConnInterface) CampusClient {
	return &campusClient{cc}
}

func (c *campusClient) GetTopNews(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTopNewsReply, error) {
	out := new(GetTopNewsReply)
	err := c.cc.Invoke(ctx, "/api.Campus/GetTopNews", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *campusClient) GetNewsSources(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Campus_GetNewsSourcesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Campus_ServiceDesc.Streams[0], "/api.Campus/GetNewsSources", opts...)
	if err != nil {
		return nil, err
	}
	x := &campusGetNewsSourcesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Campus_GetNewsSourcesClient interface {
	Recv() (*NewsSource, error)
	grpc.ClientStream
}

type campusGetNewsSourcesClient struct {
	grpc.ClientStream
}

func (x *campusGetNewsSourcesClient) Recv() (*NewsSource, error) {
	m := new(NewsSource)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CampusServer is the server API for Campus service.
// All implementations must embed UnimplementedCampusServer
// for forward compatibility
type CampusServer interface {
	GetTopNews(context.Context, *emptypb.Empty) (*GetTopNewsReply, error)
	GetNewsSources(*emptypb.Empty, Campus_GetNewsSourcesServer) error
	mustEmbedUnimplementedCampusServer()
}

// UnimplementedCampusServer must be embedded to have forward compatible implementations.
type UnimplementedCampusServer struct {
}

func (UnimplementedCampusServer) GetTopNews(context.Context, *emptypb.Empty) (*GetTopNewsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopNews not implemented")
}
func (UnimplementedCampusServer) GetNewsSources(*emptypb.Empty, Campus_GetNewsSourcesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetNewsSources not implemented")
}
func (UnimplementedCampusServer) mustEmbedUnimplementedCampusServer() {}

// UnsafeCampusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CampusServer will
// result in compilation errors.
type UnsafeCampusServer interface {
	mustEmbedUnimplementedCampusServer()
}

func RegisterCampusServer(s grpc.ServiceRegistrar, srv CampusServer) {
	s.RegisterService(&Campus_ServiceDesc, srv)
}

func _Campus_GetTopNews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CampusServer).GetTopNews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Campus/GetTopNews",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CampusServer).GetTopNews(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Campus_GetNewsSources_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CampusServer).GetNewsSources(m, &campusGetNewsSourcesServer{stream})
}

type Campus_GetNewsSourcesServer interface {
	Send(*NewsSource) error
	grpc.ServerStream
}

type campusGetNewsSourcesServer struct {
	grpc.ServerStream
}

func (x *campusGetNewsSourcesServer) Send(m *NewsSource) error {
	return x.ServerStream.SendMsg(m)
}

// Campus_ServiceDesc is the grpc.ServiceDesc for Campus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Campus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Campus",
	HandlerType: (*CampusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTopNews",
			Handler:    _Campus_GetTopNews_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetNewsSources",
			Handler:       _Campus_GetNewsSources_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "CampusService.proto",
}
