# Updating Protocol Buffers

To update or extend the API `Protocol Buffers`, edit the `CampusService.proto` file.

## Adding a new Endpoint
Add a `rpc` call inside the `service Campus{...}` namespace to add a new endpoint.

Example:
```proto
rpc GetTopNews (google.protobuf.Empty) returns (GetTopNewsReply) {
  option (google.api.http) = {
    get: "/news/top"
  };
}
```

Once done, specify the two messages or the one in case you specified `google.protobuf.Empty` for the request.

Example:
```proto
message GetTopNewsReply {
  string image_url = 1;
  string link = 2;
  google.protobuf.Timestamp created = 3;
  google.protobuf.Timestamp from = 4;
  google.protobuf.Timestamp to = 5;
}
```

# Generating Protocol Buffers

Once the `CampusService.proto` is up-to-date, you must generate the up-to-date proto files.

## System Requirements

### Arch

```bash
yay -S go protobuf
```

### Fedora

```bash
sudo dnf install go protobuf
```

## Buf
Buf is required to generate protobuf files automatically.
Follow the installation instructions detailed here: https://docs.buf.build/installation#github-releases

## Go Dependencies

```bash
go get google.golang.org/protobuf/cmd/protoc-gen-go \
       google.golang.org/grpc/cmd/protoc-gen-go-grpc

go install \
   github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
   github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
   google.golang.org/protobuf/cmd/protoc-gen-go \
   google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

## Generating

Once you have installed all dependencies, run `./generate.sh` to update the client and server `Protocol Buffers` definitions.

# Common Issues

**Issue**
```
Failure: plugin openapiv2: could not find protoc plugin for name openapiv2 - please make sure protoc-gen-openapiv2 is installed and present on your $PATH
```

**Solution**
Make sure the go binary directory is a part of your path. Execute the following piece of code, or add it to your `.bashrc`. Then try again.
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```