#!/bin/bash

BASEDIR=$(dirname "$0")
echo "making sure that this script is run from $BASEDIR"
pushd "$BASEDIR" > /dev/null || exit

go work use ./

echo downloading...
go get github.com/bufbuild/buf/cmd/buf
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

echo installing...
GO111MODULE=on go install github.com/bufbuild/buf/cmd/buf
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc


echo tidiing up
go mod tidy

popd || exit
