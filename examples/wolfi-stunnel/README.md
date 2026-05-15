# Wolfi stunnel Image

Build a `scratch` image with Wolfi's `stunnel` package, sample config, state
directory, and CA certificates.

```sh
make -C examples/wolfi-stunnel build
docker run --rm chisel-wolfi:stunnel -version
```

Run with a mounted config:

```sh
docker run --rm \
  -v "$PWD/stunnel.conf:/etc/stunnel/stunnel.conf:ro" \
  chisel-wolfi:stunnel /etc/stunnel/stunnel.conf
```
