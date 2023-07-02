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
	if !ok {
		return ErrUnauthorized
	}

	if len(md["authorization"]) == 0 {
		return ErrUnauthorized
	}

	envApiKey := env.ApiKey()

	for _, authorizationKey := range md["authorization"] {
		if authorizationKey == envApiKey {
			return nil
		}
	}

	return ErrUnauthorized
}
