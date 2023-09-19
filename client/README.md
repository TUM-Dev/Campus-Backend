# Example Go Client implementation

`client.go` is an example implementation on how to talk to the gRPC backend in Go. This kitchen sync showcases the use of the
standard system TLS CA's configured to securely connect via `api-grpc.tum.app` to the backend.

Alternatively you can also test the API using [grpcurl](https://github.com/fullstorydev/grpcurl) (also located in the `testLiveApi.sh`):
```
grpcurl -protoset <(buf build -o -) -H "x-device-id:grpc-tests" api-grpc.tum.app:443 api.Campus/GetNewsSources
```
