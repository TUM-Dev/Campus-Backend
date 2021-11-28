package web

import (
	"context"
	gw "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func HTTPServe(l net.Listener, grpcPort string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterCampusHandlerFromEndpoint(context.TODO(), mux, grpcPort, opts)
	if err != nil {
		return err
	}

	s := &http.Server{Handler: mux}
	return s.Serve(l)
}
