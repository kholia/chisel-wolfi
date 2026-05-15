# Wolfi Go Build

Builds a chiseled Wolfi Go build image with Git, curl, unzip, a C/C++
toolchain for CGO, CA certificates, and `grpc_health_probe`.

```sh
make -C examples/wolfi-go-build build
make -C examples/wolfi-go-build test
make -C examples/wolfi-go-build size
make -C examples/wolfi-go-build scan
```
