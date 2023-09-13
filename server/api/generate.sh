#!/bin/sh

# needs buf: https://docs.buf.build/installation#github-releases

echo updating the generated files
buf mod update
buf generate

echo making sure the openapi document points to the valid api
sed -i '1 a "basePath": "/v1",' ./tumdev/campus_backend.swagger.json

echo making sure that all artifacts we don\'t need are cleaned up
rm -f google/api/*.go
rm -f google/api/*.swagger.json

