package web

import (
	"context"
	"encoding/json"
	"errors"
	gw "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/backend"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func HTTPServe(l net.Listener, grpcPort string) error {
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(
		func(s string) (string, bool) {
			return s, true
		}),
		runtime.WithErrorHandler(errorHandler),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterCampusHandlerFromEndpoint(context.TODO(), mux, grpcPort, opts)
	if err != nil {
		return err
	}

	s := &http.Server{Handler: mux}
	return s.Serve(l)
}

// errorHandler translates gRPC raised by the backend into HTTP errors.
func errorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	httpStatus := http.StatusInternalServerError
	httpResponse := "Internal Server Error"
	if errors.Is(err, backend.ErrNoDeviceID) {
		httpStatus = http.StatusForbidden
		httpResponse = "No device id"
	}
	w.WriteHeader(httpStatus)
	// Marshal won't fail, we know all inputs.
	resp, _ := json.Marshal(errorResponse{Error: httpResponse})
	_, err = w.Write(resp)
	if err != nil {
		log.WithError(err).Error("Error writing response")
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
