package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/env"
	"github.com/makasim/sentryhook"
	"io/fs"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"

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
	"gorm.io/gorm"
)

const httpPort = ":50051"

// Version is injected at build time by the compiler with the correct git-commit-sha or "dev" in development
var Version = "dev"

//go:embed swagger
var swagfs embed.FS

func main() {
	setupTelemetry()
	defer sentry.Flush(2 * time.Second) // make sure that sentry handles shutdowns gracefully

	// Connect to DB
	var conn gorm.Dialector
	shouldAutoMigrate := false
	dbHost := os.Getenv("DB_DSN")
	if dbHost != "" {
		log.Info("Connecting to dsn")
		conn = mysql.Open(dbHost)
	} else {
		log.Error("Failed to start! The 'DB_DSN' environment variable is not defined. Take a look at the README.md for more details.")
		os.Exit(-1)
	}

	// initializing connection to InfluxDB
	err := backend.ConnectToInfluxDB()
	if errors.Is(err, backend.ErrInfluxTokenNotConfigured) {
		log.Warn("InfluxDB token not configured - continuing without InfluxDB")
	} else if errors.Is(err, backend.ErrInfluxURLNotConfigured) {
		log.Warn("InfluxDB url not configured - continuing without InfluxDB")
	} else if err != nil {
		log.WithError(err).Error("InfluxDB connection failed - health check failed")
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

	// Create any other background services (these shouldn't do any long-running work here)
	cronService := cron.New(db)
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

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("healthy"))
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

	log.Info("running server")
	err = g.Wait()
	if err != nil {
		log.Error(err)
	}
}

// setupTelemetry initializes our telemetry stack
// - sentry to be connected with log
// - logrus to
func setupTelemetry() {
	environment := "development"
	log.SetLevel(log.TraceLevel)
	if env.IsProd() {
		log.SetLevel(log.InfoLevel)
		environment = "production"
		log.SetFormatter(&log.JSONFormatter{}) // simpler to query but harder to parse in the console
	}

	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			AttachStacktrace: true,
			Release:          Version,
			Dist:             Version, // see https://github.com/getsentry/sentry-react-native/issues/516 why this is equal
			Environment:      environment,
		}); err != nil {
			log.WithError(err).Error("Sentry initialization failed")
		}
		log.AddHook(sentryhook.New([]log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel}))
	} else {
		log.Info("continuing without sentry")
	}
}

// addMethodNameInterceptor adds the method name (e.g. "GetNewsSources") to the metadata as x-campus-method for later use (currently logging the devices api usage)
func addMethodNameInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-campus-method", method)
	return invoker(ctx, method, req, reply, cc, opts...)
}

// errorHandler translates gRPC raised by the backend into HTTP errors.
func errorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	errorResp := errorResponse{Error: "Internal Server Error", StatusCode: http.StatusInternalServerError}

	if strings.HasPrefix(err.Error(), "no device id") {
		errorResp.Error = "Not Authorized"
		errorResp.StatusCode = http.StatusForbidden
	}

	if s, ok := status.FromError(err); ok {
		switch s.Code() {
		case codes.NotFound:
			errorResp.Error = "Not Found"
			errorResp.StatusCode = http.StatusNotFound
		case codes.Unimplemented:
			errorResp.Error = "Not Implemented"
			errorResp.StatusCode = http.StatusNotImplemented
		case codes.InvalidArgument:
			errorResp.Error = "Invalid Argument"
			errorResp.StatusCode = http.StatusBadRequest
		case codes.Internal:
			errorResp.Error = "Internal Server Error"
			errorResp.StatusCode = http.StatusInternalServerError
		}

		if s.Message() != "" {
			errorResp.Error = s.Message()
		}

		if s.Details() != nil && len(s.Details()) > 0 {
			errorResp.Details = s.Details()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorResp.StatusCode)

	err = json.NewEncoder(w).Encode(errorResp)

	if err != nil {
		log.WithError(err).Error("Marshal error response failed")
		return
	}
}

type errorResponse struct {
	Error      string        `json:"error"`
	Details    []interface{} `json:"details,omitempty"`
	StatusCode int           `json:"-"`
}
