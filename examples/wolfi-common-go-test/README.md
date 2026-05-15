# wolfi-common-go-test

Builds a Go test base with CGO build tooling, buf, protoc, protobuf tools, and `grpc-health-probe`. Project-specific Go test helper binaries can be installed in a downstream layer.

```sh
make -C examples/wolfi-common-go-test all
```
