package http_header_authorization

import (
	"context"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/env"
	"google.golang.org/grpc/metadata"
)

var (
	ErrUnauthorized = errors.New("unauthorized! API Key is missing or invalid")
)

func CheckApiKeyAuthorization(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 {
		return ErrUnauthorized
	}

	envApiKey := env.ApiKey()

	// If no API key is set, allow all requests
	if !envApiKey.Valid {
		return nil
	}

	for _, authorizationKey := range md["authorization"] {
		if authorizationKey == envApiKey.String {
			return nil
		}
	}

	return ErrUnauthorized
}
