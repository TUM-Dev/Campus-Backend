#!/bin/bash

grpcurl -protoset <(buf build -o -) -plaintext -H "x-device-id:grpc-tests" api.tum.app:50052 api.Campus/GetNewsSources
