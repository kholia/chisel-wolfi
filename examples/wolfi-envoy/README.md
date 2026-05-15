# Wolfi Envoy Image

This example builds a `scratch` Envoy image from Wolfi APK packages using
Chisel. It builds Chisel from this checkout, cuts the Envoy binary and
default Wolfi Envoy config, then copies the result into the final image.

Build from the repository root:

```sh
make -C examples/wolfi-envoy build
```

Check the Envoy binary:

```sh
docker run --rm chisel-wolfi:envoy --version
```

Run Envoy with the included default config:

```sh
docker run --rm -p 10000:10000 -p 9901:9901 chisel-wolfi:envoy
```

The image exposes Envoy's listener on port `10000` and admin interface on
port `9901`, matching the configured Wolfi Envoy package.
