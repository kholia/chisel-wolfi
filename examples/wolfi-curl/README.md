# Wolfi curl Image

Build a `scratch` image with Wolfi's `curl` package and CA certificates.

```sh
make -C examples/wolfi-curl build
docker run --rm chisel-wolfi:curl --version
docker run --rm chisel-wolfi:curl -I https://www.chainguard.dev/
```
