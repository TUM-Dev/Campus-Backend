package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

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
)

const httpPort = ":50051"

const shutdownTimeout = 10 * time.Second

// Version is injected at build time by the compiler with the correct git-commit-sha or "dev" in development
var Version = "dev"

//go:embed swagger
var swagfs embed.FS

var (
	grpcRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "campus_grpc_requests_total",
		Help: "Total number of gRPC requests handled by method and status code.",
	}, []string{"method", "code"})
	grpcRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "campus_grpc_request_duration_seconds",
		Help:    "Latency distribution of gRPC requests by method.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})
)

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
	httpServer := &http.Server{
		Handler:           httpMux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	httpMux.HandleFunc("/imprint", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})

	httpMux.HandleFunc("/health", healthHandler(db))
	httpMux.Handle("/metrics", promhttp.Handler())

	httpMux.Handle("/", http.RedirectHandler("/swagger/", http.StatusTemporaryRedirect))
	httpMux.Handle("/files/", http.StripPrefix("/files/", utils.PublicFileServer("/Storage")))
	httpMux.Handle("/rating-images/", http.StripPrefix("/rating-images/", campusService.RatingImageHandler()))
	httpMux.Handle("/swagger/", http.FileServer(http.FS(swagfs)))

	// Main GRPC Server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(UnaryRequestLogger, UnaryPanicRecovery),
		grpc.ChainStreamInterceptor(StreamRequestLogger, StreamPanicRecovery),
	)
	pb.RegisterCampusServer(grpcServer, campusService)
	reflection.Register(grpcServer)

	// GRPC Gateway for HTTP REST -> GRPC
	grpcGatewayMux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(
		func(key string) (string, bool) {
			key = textproto.CanonicalMIMEHeaderKey(key)
			switch key {
			case "X-Device-Id", "X-Os-Version", "X-App-Version":
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
	}
	if err := pb.RegisterCampusHandlerFromEndpoint(context.Background(), grpcGatewayMux, httpPort, opts); err != nil {
		log.WithError(err).Fatal("could not RegisterCampusHandlerFromEndpoint")
	}
	httpMux.Handle("/v1/", http.StripPrefix("/v1", grpcGatewayMux))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run periodic device flusher in the background so it does not block graceful shutdown.
	go func() {
		if err := campusService.RunDeviceFlusher(); err != nil {
			log.WithError(err).Error("device flusher exited with error")
		}
	}()

	// Start each server in its own goroutine and shut down gracefully on SIGINT/SIGTERM.
	g := errgroup.Group{}
	g.Go(func() error {
		if err := grpcServer.Serve(grpcListener); err != nil && !errors.Is(err, cmux.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := httpServer.Serve(httpListener); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, cmux.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := mux.Serve(); err != nil && !errors.Is(err, net.ErrClosed) && !errors.Is(err, cmux.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error { return cronService.RunContext(ctx) })
	g.Go(func() error {
		<-ctx.Done()
		log.Info("shutdown signal received, draining servers")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		done := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(done)
		}()
		select {
		case <-done:
		case <-shutdownCtx.Done():
			log.Warn("gRPC graceful stop timed out, forcing stop")
			grpcServer.Stop()
		}

		if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("HTTP server shutdown error")
		}
		mux.Close()
		_ = listener.Close()
		return nil
	})

	log.Info("running server")
	if err := g.Wait(); err != nil {
		log.WithError(err).Error("encountered issue while running the server")
	}
}

func healthHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}
		sqlDB, err := db.DB()
		if err != nil {
			log.WithError(err).Warn("health check failed to obtain sql.DB")
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}
		pingCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(pingCtx); err != nil {
			log.WithError(err).Warn("health check database ping failed")
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("healthy"))
	}
}

func UnaryPanicRecovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			method := strings.TrimPrefix(info.FullMethod, "/api.Campus/")
			log.WithFields(log.Fields{"method": method, "panic": r, "stack": string(debug.Stack())}).Error("recovered from panic in gRPC unary handler")
			sentry.CurrentHub().Recover(r)
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

func StreamPanicRecovery(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			method := strings.TrimPrefix(info.FullMethod, "/api.Campus/")
			log.WithFields(log.Fields{"method": method, "panic": r, "stack": string(debug.Stack())}).Error("recovered from panic in gRPC stream handler")
			sentry.CurrentHub().Recover(r)
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return handler(srv, stream)
}

func UnaryRequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	elapsed := time.Since(start)
	method := strings.TrimPrefix(info.FullMethod, "/api.Campus/")
	grpcRequestsTotal.WithLabelValues(method, status.Code(err).String()).Inc()
	grpcRequestDuration.WithLabelValues(method).Observe(elapsed.Seconds())
	fields := log.Fields{"elapsed": elapsed, "method": method}
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
	elapsed := time.Since(start)
	method := strings.TrimPrefix(info.FullMethod, "/api.Campus/")
	grpcRequestsTotal.WithLabelValues(method, status.Code(err).String()).Inc()
	grpcRequestDuration.WithLabelValues(method).Observe(elapsed.Seconds())
	fields := log.Fields{"elapsed": elapsed, "method": method}
	if err == nil {
		log.WithFields(fields).Info("request")
	} else {
		log.WithFields(fields).WithError(err).Warn("request")
	}
	return err
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
