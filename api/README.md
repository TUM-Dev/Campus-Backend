
# Install

1. Setup Go and `protoc` on your system: `yay -S go protobuf`
1. Install go plugins for protoc:
```
$ go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
1. Run ./generate.sh to update server/client code