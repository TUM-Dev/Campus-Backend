#!/usr/bin/env bash
set -euo pipefail

BASEDIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
cd "$BASEDIR/.."

# Match the generators to the application dependencies without changing go.mod.
PROTOBUF_VERSION="$(go list -m -f '{{.Version}}' google.golang.org/protobuf)"
GATEWAY_VERSION="$(go list -m -f '{{.Version}}' github.com/grpc-ecosystem/grpc-gateway/v2)"

# Install tools independently of the application's module graph. Keep standalone
# tool versions pinned so CI regeneration is reproducible.
go install github.com/bufbuild/buf/cmd/buf@v1.75.0
go install "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@$GATEWAY_VERSION"
go install "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@$GATEWAY_VERSION"
go install "google.golang.org/protobuf/cmd/protoc-gen-go@$PROTOBUF_VERSION"
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1
go install golang.org/x/tools/cmd/goimports@v0.52.0
