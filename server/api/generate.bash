#!/bin/bash

# needs buf: https://docs.buf.build/installation#github-releases
BASEDIR=$(dirname "$0")
echo making sure that this script is run from $BASEDIR
pushd $BASEDIR > /dev/null

echo updating the generated files
export PATH="$PATH:$(go env GOPATH)/bin"
buf dep update || exit 1
buf generate || exit 1

echo making sure the openapi document points to the valid api
grep -q '"basePath": "/v1"' ./tumdev/campus_backend.swagger.json || sed -i '1 a "basePath": "/v1",' ./tumdev/campus_backend.swagger.json

echo making sure that all artifacts we don\'t need are cleaned up
rm -f google/api/*.go
rm -f google/api/*.swagger.json

echo maing sure that the generated files are formatted
go fmt tumdev/*.go || exit 1
goimports -w tumdev/*.go || exit 1
buf format -w --path tumdev || exit 1

# clean up the stack
popd > /dev/null
