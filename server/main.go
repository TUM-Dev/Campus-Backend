package main

import (
	"context"
	"embed"
	"encoding/json"
	"net"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"google.golang.org/grpc/reflection"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/getsentry/sentry-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const httpPort = ":50051"

// Version is injected at build time by the compiler with the correct git-commit-sha or "dev" in development
var Version = "dev"

//go:embed swagger
var swagfs embed.FS

func main() {
	utils.SetupTelemetry(Version)
	defer sentry.Flush(10 * time.Second) // make sure that sentry handles shutdowns gracefully

	db := utils.SetupDB()

	// Create any other background services (these shouldn't do any long-running work here)
	cronService := cron.New(db)
	campusService := backend.New(db)

	// Listen to our configured ports
	listener, err := net.Listen("tcp", httpPort)
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}
	mux := cmux.New(listener)
	grpcListener := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := mux.Match(cmux.HTTP1Fast())

	// HTTP Stuff
	httpMux := http.NewServeMux()
	httpServer := &http.Server{Handler: httpMux}
	httpMux.HandleFunc("/imprint", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})

	httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("healthy"))
	})
	httpMux.Handle("/metrics", promhttp.Handler())

	httpMux.Handle("/", http.RedirectHandler("/swagger/", http.StatusTemporaryRedirect))
	httpMux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("/Storage"))))
	httpMux.Handle("/swagger/", http.FileServer(http.FS(swagfs)))

	// Main GRPC Server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(UnaryRequestLogger), grpc.StreamInterceptor(StreamRequestLogger))
	pb.RegisterCampusServer(grpcServer, campusService)
	reflection.Register(grpcServer)

	// GRPC Gateway for HTTP REST -> GRPC
	grpcGatewayMux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(
		func(key string) (string, bool) {
			key = textproto.CanonicalMIMEHeaderKey(key)
			if key == "X-Device-Id" {
				return key, true
			}
			// don't filter headers (pass all to gRPC handlers)
			return runtime.DefaultHeaderMatcher(key)
		}),
		runtime.WithErrorHandler(errorHandler),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent("internal"),
		grpc.WithUnaryInterceptor(addMethodNameInterceptor),
	}
	if err := pb.RegisterCampusHandlerFromEndpoint(context.Background(), grpcGatewayMux, httpPort, opts); err != nil {
		log.WithError(err).Fatal("could not RegisterCampusHandlerFromEndpoint")
	}
	httpMux.Handle("/v1/", http.StripPrefix("/v1", grpcGatewayMux))

	// Start each server in its own go routine and logs any errors
	g := errgroup.Group{}
	g.Go(func() error { return grpcServer.Serve(grpcListener) })
	g.Go(func() error { return httpServer.Serve(httpListener) })
	g.Go(func() error { return mux.Serve() })
	g.Go(func() error { return cronService.Run() })                // Setup cron jobs
	g.Go(func() error { return campusService.RunDeviceFlusher() }) // Setup campus service

	log.Info("running server")
	if err := g.Wait(); err != nil {
		log.WithError(err).Error("encountered issue while running the server")
	}
}

func UnaryRequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	fields := log.Fields{"elapsed": time.Since(start), "method": strings.TrimPrefix(info.FullMethod, "/api.Campus/")}
	if err == nil {
		log.WithFields(fields).Info("request")
	} else {
		log.WithFields(fields).WithError(err).Warn("request")
	}
	return resp, err
}
func StreamRequestLogger(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, stream)
	fields := log.Fields{"elapsed": time.Since(start), "method": strings.TrimPrefix(info.FullMethod, "/api.Campus/")}
	if err != nil {
		log.WithFields(fields).Info("request")
	} else {
		log.WithFields(fields).WithError(err).Warn("request")
	}
	return err
}

// addMethodNameInterceptor adds the method name (e.g. "ListNewsSources") to the metadata as x-campus-method for later use (currently logging the devices api usage)
func addMethodNameInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-campus-method", method)
	return invoker(ctx, method, req, reply, cc, opts...)
}

// errorHandler translates gRPC raised by the backend into HTTP errors.
func errorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	errorResp := utils.GrpcErrorToWebError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorResp.StatusCode)

	if err = json.NewEncoder(w).Encode(errorResp); err != nil {
		log.WithError(err).Error("Marshal error response failed")
	}
}
