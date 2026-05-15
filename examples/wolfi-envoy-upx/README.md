# Wolfi Envoy UPX Image

This is an opt-in variant of `wolfi-envoy` that compresses `/usr/bin/envoy`
with UPX during the build stage.

```sh
make -C examples/wolfi-envoy-upx build
docker run --rm chisel-wolfi:envoy-upx --version
docker run --rm chisel-wolfi:envoy-upx --mode validate -c /etc/envoy/envoy.yaml
```

The final image does not include UPX. It keeps Wolfi's `os-release`, APK
database, Envoy config, and CA certificates, but the Envoy binary bytes are
post-processed after Chisel extracts them from the signed Wolfi APK. Use the
regular `wolfi-envoy` example when exact package-file provenance matters more
than binary-size experimentation.

UPX can reduce the unpacked filesystem size while increasing the compressed
container artifact size because packed executables are less compressible by the
image layer compressor. This example is intentionally not part of the default
example set or publish workflow.

The default uses `--best` to maximize the experiment's unpacked binary-size
reduction. Override `UPX_FLAGS` when testing different compression settings:

```sh
make -C examples/wolfi-envoy-upx build BUILD_ARGS='--build-arg UPX_FLAGS=--fast'
```
