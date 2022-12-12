package main

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
	"github.com/getsentry/sentry-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	httpPort = ":50051"
)

var Version = "dev"

//go:embed swagger
var swagfs embed.FS

func main() {
	// Connect to DB
	var conn gorm.Dialector
	shouldAutoMigrate := false
	if dbHost := os.Getenv("DB_DSN"); dbHost != "" {
		log.Info("Connecting to dsn")
		conn = mysql.Open(dbHost)
	} else {
		log.Info("Switching to test.db")
		conn = sqlite.Open("test.db")
		shouldAutoMigrate = true
	}

	environment := "development"
	if Version != "dev" {
		environment = "production"
	}
	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:         os.Getenv("SENTRY_DSN"),
			Release:     Version,
			Environment: environment,
		}); err != nil {
			log.WithError(err).Error("Sentry initialization failed")
		}
	} else {
		log.Println("continuing without sentry")
	}
	db, err := gorm.Open(conn, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	tumMigrator := migration.New(db, shouldAutoMigrate)
	err = tumMigrator.Migrate()
	if err != nil {
		log.WithError(err).Fatal("Failed to migrate database")
		return
	}

	var mensaCronActivated = true
	if len(os.Args) > 2 && os.Args[1] == "-MensaCron" && os.Args[2] == "0" {
		mensaCronActivated = false
		log.Println("Cronjobs for the cafeteria rating are deactivated. Remove commandline argument <-MensaCron 0> or set it to 1.", len(os.Args))
	}
	// Create any other background services (these shouldn't do any long running work here)
	cronService := cron.New(db, mensaCronActivated)
	campusService := backend.New(db)

	// Listen to our configured ports
	listener, err := net.Listen("tcp", httpPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	m := cmux.New(listener)
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	// HTTP Stuff
	mux := http.NewServeMux()
	httpServer := &http.Server{Handler: mux}
	mux.HandleFunc("/imprint", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})
	static, _ := fs.Sub(swagfs, "swagger")
	mux.Handle("/", http.FileServer(http.FS(static)))

	// Main GRPC Server
	grpcS := grpc.NewServer()
	pb.RegisterCampusServer(grpcS, campusService)

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
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent("internal"),
		grpc.WithUnaryInterceptor(addMethodNameInterceptor),
	}
	if err := pb.RegisterCampusHandlerFromEndpoint(ctx, grpcGatewayMux, httpPort, opts); err != nil {
		panic(err)
	}
	restPrefix := "/v1"
	mux.Handle("/v1/", http.StripPrefix(restPrefix, grpcGatewayMux))

	// Start each server in its own go routine and logs any errors
	g := errgroup.Group{}
	g.Go(func() error { return grpcS.Serve(grpcListener) })
	g.Go(func() error { return httpServer.Serve(httpListener) })
	g.Go(func() error { return m.Serve() })
	g.Go(func() error { return cronService.Run() })                // Setup cron jobs
	g.Go(func() error { return campusService.RunDeviceFlusher() }) // Setup campus service

	log.Println("running server")
	err = g.Wait()
	if err != nil {
		log.Error(err)
	}
}

// addMethodNameInterceptor adds the method name (e.g. "GetNewsSources") to the metadata as x-campus-method for later use (currently logging the devices api usage)
func addMethodNameInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-campus-method", method)
	return invoker(ctx, method, req, reply, cc, opts...)
}

// errorHandler translates gRPC raised by the backend into HTTP errors.
func errorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	httpStatus := http.StatusInternalServerError
	httpResponse := "Internal Server Error"
	if strings.HasPrefix(err.Error(), "no device id") {
		httpStatus = http.StatusForbidden
		httpResponse = "Not Authorized"
	}
	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.NotFound:
			httpStatus = http.StatusNotFound
			httpResponse = "Not Found"
		case codes.Unimplemented:
			httpStatus = http.StatusNotImplemented
			httpResponse = "Not Implemented"
		}
	}

	w.WriteHeader(httpStatus)
	resp, err := json.Marshal(errorResponse{Error: httpResponse})
	if err != nil {
		log.WithError(err).Error("Marshal error response failed")
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		log.WithError(err).Error("Error writing response")
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
