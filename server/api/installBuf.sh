#!/bin/bash

go get github.com/bufbuild/buf/cmd/buf \
       google.golang.org/protobuf/cmd/protoc-gen-go \
       google.golang.org/grpc/cmd/protoc-gen-go-grpc

go install \
       github.com/bufbuild/buf/cmd/buf \
       github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
       github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
       google.golang.org/protobuf/cmd/protoc-gen-go \
       google.golang.org/grpc/cmd/protoc-gen-go-grpc

go mod tidy
