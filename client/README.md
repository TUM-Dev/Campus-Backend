# Example Go Client implementation

`client.go` is an example implementation on how to talk to the gRPC backend in Go. This kitchen sync showcases the use of the
standard system TLS CA's configured to securely connect via `api-grpc.tum.app` to the backend.

Alternatively you can also test the API using [grpcurl](https://github.com/fullstorydev/grpcurl) (also located in the `testLiveApi.sh`):
```
grpcurl -protoset <(buf build -o -) -H "x-device-id:grpc-tests" api-grpc.tum.app:443 api.Campus/ListNewsSources
```

## Live API test harness

`livetest/` probes a running backend over gRPC and REST and reports anything that looks broken:
unimplemented RPCs, 5xx responses to valid or invalid requests, stale imports (news, movies, dishes, head count),
broken image URLs, REST/gRPC mismatches, student club language parity, swagger and reflection drift,
publicly enumerable upload directories, TLS expiry and latency under a small burst.

It is read-only: write RPCs are only called with requests that fail validation before anything is stored.

```bash
go run ./livetest                 # against production (api.tum.app)
go run ./livetest -local          # against a local server on 127.0.0.1:50051
go run ./livetest -run 'news/' -v # only matching checks, including notes of passing ones
go run ./livetest -json           # machine-readable output
```

It exits with 1 if any check fails (with `-strict`, also on warnings). New RPCs are skipped by the sweep
until they are added to `emptyRequestSafe` in `livetest/sweep_checks.go`.
The harness uses the generated stubs of the server version pinned in `go.mod`; bump it when the proto changes.
