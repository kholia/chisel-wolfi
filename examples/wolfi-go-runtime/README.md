# Wolfi Go Runtime

Builds a small chiseled Wolfi runtime base for Go services with CA
certificates and `grpc_health_probe`.

```sh
make -C examples/wolfi-go-runtime build
make -C examples/wolfi-go-runtime test
make -C examples/wolfi-go-runtime size
make -C examples/wolfi-go-runtime scan
```
