# Wolfi Fluent Bit Image

Build a `scratch` Fluent Bit image from Wolfi APK slices.

```sh
make -C examples/wolfi-fluent-bit build
docker run --rm chisel-wolfi:fluent-bit --version
docker run --rm chisel-wolfi:fluent-bit --dry-run -c /fluent-bit/etc/fluent-bit.conf
```

The image includes Wolfi's Fluent Bit binary, bundled Fluent Bit config and
parser files under `/fluent-bit/etc`, CA certificates, `os-release`, and the
APK installed database for scanners.

Run with the default Wolfi config:

```sh
docker run --rm -p 2020:2020 chisel-wolfi:fluent-bit
```

Mount your own config by overriding the command:

```sh
docker run --rm \
  -v "$PWD/fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf:ro" \
  chisel-wolfi:fluent-bit -c /fluent-bit/etc/fluent-bit.conf
```
