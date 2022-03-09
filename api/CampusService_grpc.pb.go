// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: CampusService.proto

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
	GetNewsSources(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NewsSourceArray, error)
	SearchRooms(ctx context.Context, in *SearchRoomsRequest, opts ...grpc.CallOption) (*SearchRoomsReply, error)
	// a location is a campus location/building, e.g. "Garching Forschungszentrum"
	GetLocations(ctx context.Context, in *GetLocationsRequest, opts ...grpc.CallOption) (*GetLocationsReply, error)
	GetRoomMaps(ctx context.Context, in *GetRoomMapsRequest, opts ...grpc.CallOption) (*GetRoomMapsReply, error)
	GetRoomCoordinates(ctx context.Context, in *GetRoomCoordinatesRequest, opts ...grpc.CallOption) (*GetRoomCoordinatesReply, error)
	GetRoomSchedule(ctx context.Context, in *GetRoomScheduleRequest, opts ...grpc.CallOption) (*GetRoomScheduleReply, error)
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

func (c *campusClient) GetNewsSources(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*NewsSourceArray, error) {
	out := new(NewsSourceArray)
	err := c.cc.Invoke(ctx, "/api.Campus/GetNewsSources", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *campusClient) SearchRooms(ctx context.Context, in *SearchRoomsRequest, opts ...grpc.CallOption) (*SearchRoomsReply, error) {
	out := new(SearchRoomsReply)
	err := c.cc.Invoke(ctx, "/api.Campus/SearchRooms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *campusClient) GetLocations(ctx context.Context, in *GetLocationsRequest, opts ...grpc.CallOption) (*GetLocationsReply, error) {
	out := new(GetLocationsReply)
	err := c.cc.Invoke(ctx, "/api.Campus/GetLocations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *campusClient) GetRoomMaps(ctx context.Context, in *GetRoomMapsRequest, opts ...grpc.CallOption) (*GetRoomMapsReply, error) {
	out := new(GetRoomMapsReply)
	err := c.cc.Invoke(ctx, "/api.Campus/GetRoomMaps", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *campusClient) GetRoomCoordinates(ctx context.Context, in *GetRoomCoordinatesRequest, opts ...grpc.CallOption) (*GetRoomCoordinatesReply, error) {
	out := new(GetRoomCoordinatesReply)
	err := c.cc.Invoke(ctx, "/api.Campus/GetRoomCoordinates", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *campusClient) GetRoomSchedule(ctx context.Context, in *GetRoomScheduleRequest, opts ...grpc.CallOption) (*GetRoomScheduleReply, error) {
	out := new(GetRoomScheduleReply)
	err := c.cc.Invoke(ctx, "/api.Campus/GetRoomSchedule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CampusServer is the server API for Campus service.
// All implementations must embed UnimplementedCampusServer
// for forward compatibility
type CampusServer interface {
	GetTopNews(context.Context, *emptypb.Empty) (*GetTopNewsReply, error)
	GetNewsSources(context.Context, *emptypb.Empty) (*NewsSourceArray, error)
	SearchRooms(context.Context, *SearchRoomsRequest) (*SearchRoomsReply, error)
	// a location is a campus location/building, e.g. "Garching Forschungszentrum"
	GetLocations(context.Context, *GetLocationsRequest) (*GetLocationsReply, error)
	GetRoomMaps(context.Context, *GetRoomMapsRequest) (*GetRoomMapsReply, error)
	GetRoomCoordinates(context.Context, *GetRoomCoordinatesRequest) (*GetRoomCoordinatesReply, error)
	GetRoomSchedule(context.Context, *GetRoomScheduleRequest) (*GetRoomScheduleReply, error)
	mustEmbedUnimplementedCampusServer()
}

// UnimplementedCampusServer must be embedded to have forward compatible implementations.
type UnimplementedCampusServer struct {
}

func (UnimplementedCampusServer) GetTopNews(context.Context, *emptypb.Empty) (*GetTopNewsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopNews not implemented")
}
func (UnimplementedCampusServer) GetNewsSources(context.Context, *emptypb.Empty) (*NewsSourceArray, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNewsSources not implemented")
}
func (UnimplementedCampusServer) SearchRooms(context.Context, *SearchRoomsRequest) (*SearchRoomsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchRooms not implemented")
}
func (UnimplementedCampusServer) GetLocations(context.Context, *GetLocationsRequest) (*GetLocationsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocations not implemented")
}
func (UnimplementedCampusServer) GetRoomMaps(context.Context, *GetRoomMapsRequest) (*GetRoomMapsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomMaps not implemented")
}
func (UnimplementedCampusServer) GetRoomCoordinates(context.Context, *GetRoomCoordinatesRequest) (*GetRoomCoordinatesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomCoordinates not implemented")
}
func (UnimplementedCampusServer) GetRoomSchedule(context.Context, *GetRoomScheduleRequest) (*GetRoomScheduleReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomSchedule not implemented")
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

func _Campus_GetNewsSources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CampusServer).GetNewsSources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Campus/GetNewsSources",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CampusServer).GetNewsSources(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Campus_SearchRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CampusServer).SearchRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Campus/SearchRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CampusServer).SearchRooms(ctx, req.(*SearchRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Campus_GetLocations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLocationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CampusServer).GetLocations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Campus/GetLocations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CampusServer).GetLocations(ctx, req.(*GetLocationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Campus_GetRoomMaps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomMapsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CampusServer).GetRoomMaps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Campus/GetRoomMaps",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CampusServer).GetRoomMaps(ctx, req.(*GetRoomMapsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Campus_GetRoomCoordinates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomCoordinatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CampusServer).GetRoomCoordinates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Campus/GetRoomCoordinates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CampusServer).GetRoomCoordinates(ctx, req.(*GetRoomCoordinatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Campus_GetRoomSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CampusServer).GetRoomSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Campus/GetRoomSchedule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CampusServer).GetRoomSchedule(ctx, req.(*GetRoomScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
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
		{
			MethodName: "GetNewsSources",
			Handler:    _Campus_GetNewsSources_Handler,
		},
		{
			MethodName: "SearchRooms",
			Handler:    _Campus_SearchRooms_Handler,
		},
		{
			MethodName: "GetLocations",
			Handler:    _Campus_GetLocations_Handler,
		},
		{
			MethodName: "GetRoomMaps",
			Handler:    _Campus_GetRoomMaps_Handler,
		},
		{
			MethodName: "GetRoomCoordinates",
			Handler:    _Campus_GetRoomCoordinates_Handler,
		},
		{
			MethodName: "GetRoomSchedule",
			Handler:    _Campus_GetRoomSchedule_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "CampusService.proto",
}
