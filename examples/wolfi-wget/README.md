# Wolfi wget Image

Build a `scratch` image with Wolfi's `wget` package and CA certificates.

```sh
make -C examples/wolfi-wget build
docker run --rm chisel-wolfi:wget --version
docker run --rm chisel-wolfi:wget -S --spider https://www.chainguard.dev/
```
