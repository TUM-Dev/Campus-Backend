#!/bin/sh

# needs buf: https://docs.buf.build/installation#github-releases

buf mod update
buf generate

sed -i '1 a "basePath": "/v1",' ./gen/openapiv2/CampusService.swagger.json
cp ./gen/openapiv2/CampusService.swagger.json ../server/swagger/swagger.json

