// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: CampusService.proto

/*
Package api is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package api

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_Campus_GetTopNews_0(ctx context.Context, marshaler runtime.Marshaler, client CampusClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq emptypb.Empty
	var metadata runtime.ServerMetadata

	msg, err := client.GetTopNews(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Campus_GetTopNews_0(ctx context.Context, marshaler runtime.Marshaler, server CampusServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq emptypb.Empty
	var metadata runtime.ServerMetadata

	msg, err := server.GetTopNews(ctx, &protoReq)
	return msg, metadata, err

}

func request_Campus_GetNewsSources_0(ctx context.Context, marshaler runtime.Marshaler, client CampusClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq emptypb.Empty
	var metadata runtime.ServerMetadata

	msg, err := client.GetNewsSources(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Campus_GetNewsSources_0(ctx context.Context, marshaler runtime.Marshaler, server CampusServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq emptypb.Empty
	var metadata runtime.ServerMetadata

	msg, err := server.GetNewsSources(ctx, &protoReq)
	return msg, metadata, err

}

func request_Campus_SearchRooms_0(ctx context.Context, marshaler runtime.Marshaler, client CampusClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SearchRoomsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["query"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "query")
	}

	protoReq.Query, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "query", err)
	}

	msg, err := client.SearchRooms(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Campus_SearchRooms_0(ctx context.Context, marshaler runtime.Marshaler, server CampusServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq SearchRoomsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["query"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "query")
	}

	protoReq.Query, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "query", err)
	}

	msg, err := server.SearchRooms(ctx, &protoReq)
	return msg, metadata, err

}

func request_Campus_GetLocations_0(ctx context.Context, marshaler runtime.Marshaler, client CampusClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq LocationsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["location"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "location")
	}

	protoReq.Location, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "location", err)
	}

	msg, err := client.GetLocations(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Campus_GetLocations_0(ctx context.Context, marshaler runtime.Marshaler, server CampusServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq LocationsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["location"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "location")
	}

	protoReq.Location, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "location", err)
	}

	msg, err := server.GetLocations(ctx, &protoReq)
	return msg, metadata, err

}

func request_Campus_GetAvailableMaps_0(ctx context.Context, marshaler runtime.Marshaler, client CampusClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetAvailableMapsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["arch_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "arch_id")
	}

	protoReq.ArchId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "arch_id", err)
	}

	msg, err := client.GetAvailableMaps(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Campus_GetAvailableMaps_0(ctx context.Context, marshaler runtime.Marshaler, server CampusServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetAvailableMapsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["arch_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "arch_id")
	}

	protoReq.ArchId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "arch_id", err)
	}

	msg, err := server.GetAvailableMaps(ctx, &protoReq)
	return msg, metadata, err

}

func request_Campus_GetCoordinates_0(ctx context.Context, marshaler runtime.Marshaler, client CampusClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetCoordinatesRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["arch_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "arch_id")
	}

	protoReq.ArchId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "arch_id", err)
	}

	msg, err := client.GetCoordinates(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Campus_GetCoordinates_0(ctx context.Context, marshaler runtime.Marshaler, server CampusServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetCoordinatesRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["arch_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "arch_id")
	}

	protoReq.ArchId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "arch_id", err)
	}

	msg, err := server.GetCoordinates(ctx, &protoReq)
	return msg, metadata, err

}

func request_Campus_GetRoomSchedule_0(ctx context.Context, marshaler runtime.Marshaler, client CampusClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetRoomScheduleRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["room_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "room_id")
	}

	protoReq.RoomId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "room_id", err)
	}

	val, ok = pathParams["start"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "start")
	}

	protoReq.Start, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "start", err)
	}

	val, ok = pathParams["end"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "end")
	}

	protoReq.End, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "end", err)
	}

	msg, err := client.GetRoomSchedule(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Campus_GetRoomSchedule_0(ctx context.Context, marshaler runtime.Marshaler, server CampusServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetRoomScheduleRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["room_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "room_id")
	}

	protoReq.RoomId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "room_id", err)
	}

	val, ok = pathParams["start"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "start")
	}

	protoReq.Start, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "start", err)
	}

	val, ok = pathParams["end"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "end")
	}

	protoReq.End, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "end", err)
	}

	msg, err := server.GetRoomSchedule(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterCampusHandlerServer registers the http handlers for service Campus to "mux".
// UnaryRPC     :call CampusServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterCampusHandlerFromEndpoint instead.
func RegisterCampusHandlerServer(ctx context.Context, mux *runtime.ServeMux, server CampusServer) error {

	mux.Handle("GET", pattern_Campus_GetTopNews_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/api.Campus/GetTopNews", runtime.WithHTTPPathPattern("/topnews"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Campus_GetTopNews_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetTopNews_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetNewsSources_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/api.Campus/GetNewsSources", runtime.WithHTTPPathPattern("/news/sources"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Campus_GetNewsSources_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetNewsSources_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetNewsSources_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_SearchRooms_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/api.Campus/SearchRooms", runtime.WithHTTPPathPattern("/roomfinder/room/search/{query}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Campus_SearchRooms_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_SearchRooms_0(ctx, mux, outboundMarshaler, w, req, response_Campus_SearchRooms_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetLocations_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/api.Campus/GetLocations", runtime.WithHTTPPathPattern("/locations/{location}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Campus_GetLocations_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetLocations_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetLocations_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetAvailableMaps_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/api.Campus/GetAvailableMaps", runtime.WithHTTPPathPattern("/roomfinder/room/availableMaps/{arch_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Campus_GetAvailableMaps_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetAvailableMaps_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetAvailableMaps_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetCoordinates_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/api.Campus/GetCoordinates", runtime.WithHTTPPathPattern("/roomfinder/room/coordinates/{arch_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Campus_GetCoordinates_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetCoordinates_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetRoomSchedule_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/api.Campus/GetRoomSchedule", runtime.WithHTTPPathPattern("/roomfinder/room/scheduleById/{room_id}/{start}/{end}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Campus_GetRoomSchedule_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetRoomSchedule_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetRoomSchedule_0{resp}, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterCampusHandlerFromEndpoint is same as RegisterCampusHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterCampusHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterCampusHandler(ctx, mux, conn)
}

// RegisterCampusHandler registers the http handlers for service Campus to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterCampusHandler(ctx context.Context, mux *runtime.ServeMux, conn grpc.ClientConnInterface) error {
	return RegisterCampusHandlerClient(ctx, mux, NewCampusClient(conn))
}

// RegisterCampusHandlerClient registers the http handlers for service Campus
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "CampusClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "CampusClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "CampusClient" to call the correct interceptors.
func RegisterCampusHandlerClient(ctx context.Context, mux *runtime.ServeMux, client CampusClient) error {

	mux.Handle("GET", pattern_Campus_GetTopNews_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/api.Campus/GetTopNews", runtime.WithHTTPPathPattern("/topnews"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Campus_GetTopNews_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetTopNews_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetNewsSources_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/api.Campus/GetNewsSources", runtime.WithHTTPPathPattern("/news/sources"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Campus_GetNewsSources_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetNewsSources_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetNewsSources_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_SearchRooms_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/api.Campus/SearchRooms", runtime.WithHTTPPathPattern("/roomfinder/room/search/{query}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Campus_SearchRooms_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_SearchRooms_0(ctx, mux, outboundMarshaler, w, req, response_Campus_SearchRooms_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetLocations_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/api.Campus/GetLocations", runtime.WithHTTPPathPattern("/locations/{location}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Campus_GetLocations_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetLocations_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetLocations_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetAvailableMaps_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/api.Campus/GetAvailableMaps", runtime.WithHTTPPathPattern("/roomfinder/room/availableMaps/{arch_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Campus_GetAvailableMaps_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetAvailableMaps_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetAvailableMaps_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetCoordinates_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/api.Campus/GetCoordinates", runtime.WithHTTPPathPattern("/roomfinder/room/coordinates/{arch_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Campus_GetCoordinates_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetCoordinates_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Campus_GetRoomSchedule_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/api.Campus/GetRoomSchedule", runtime.WithHTTPPathPattern("/roomfinder/room/scheduleById/{room_id}/{start}/{end}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Campus_GetRoomSchedule_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Campus_GetRoomSchedule_0(ctx, mux, outboundMarshaler, w, req, response_Campus_GetRoomSchedule_0{resp}, mux.GetForwardResponseOptions()...)

	})

	return nil
}

type response_Campus_GetNewsSources_0 struct {
	proto.Message
}

func (m response_Campus_GetNewsSources_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*NewsSourceArray)
	return response.Sources
}

type response_Campus_SearchRooms_0 struct {
	proto.Message
}

func (m response_Campus_SearchRooms_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*SearchRoomsReply)
	return response.Rooms
}

type response_Campus_GetLocations_0 struct {
	proto.Message
}

func (m response_Campus_GetLocations_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*LocationsReply)
	return response.Locations
}

type response_Campus_GetAvailableMaps_0 struct {
	proto.Message
}

func (m response_Campus_GetAvailableMaps_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*GetAvailableMapsReply)
	return response.Maps
}

type response_Campus_GetRoomSchedule_0 struct {
	proto.Message
}

func (m response_Campus_GetRoomSchedule_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*GetRoomScheduleReply)
	return response.Events
}

var (
	pattern_Campus_GetTopNews_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"topnews"}, ""))

	pattern_Campus_GetNewsSources_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"news", "sources"}, ""))

	pattern_Campus_SearchRooms_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"roomfinder", "room", "search", "query"}, ""))

	pattern_Campus_GetLocations_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"locations", "location"}, ""))

	pattern_Campus_GetAvailableMaps_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"roomfinder", "room", "availableMaps", "arch_id"}, ""))

	pattern_Campus_GetCoordinates_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"roomfinder", "room", "coordinates", "arch_id"}, ""))

	pattern_Campus_GetRoomSchedule_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3, 1, 0, 4, 1, 5, 4, 1, 0, 4, 1, 5, 5}, []string{"roomfinder", "room", "scheduleById", "room_id", "start", "end"}, ""))
)

var (
	forward_Campus_GetTopNews_0 = runtime.ForwardResponseMessage

	forward_Campus_GetNewsSources_0 = runtime.ForwardResponseMessage

	forward_Campus_SearchRooms_0 = runtime.ForwardResponseMessage

	forward_Campus_GetLocations_0 = runtime.ForwardResponseMessage

	forward_Campus_GetAvailableMaps_0 = runtime.ForwardResponseMessage

	forward_Campus_GetCoordinates_0 = runtime.ForwardResponseMessage

	forward_Campus_GetRoomSchedule_0 = runtime.ForwardResponseMessage
)
