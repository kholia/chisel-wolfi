# Wolfi Unbound Image

Build a `scratch` image with Wolfi's Unbound daemon, tools, and default config.

```sh
make -C examples/wolfi-unbound build
docker run --rm chisel-wolfi:unbound -V
docker run --rm --entrypoint /usr/bin/unbound-checkconf chisel-wolfi:unbound /etc/unbound/unbound.conf
```

Run Unbound:

```sh
docker run --rm -p 5353:53/udp -p 5353:53/tcp chisel-wolfi:unbound
```
