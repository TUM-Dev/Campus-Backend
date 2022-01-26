#!/bin/sh

# needs buf: https://docs.buf.build/installation#github-releases

buf mod update
buf generate

cp ./gen/openapiv2/CampusService.swagger.json ../server/swagger/swagger.json