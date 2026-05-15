# Wolfi Envoy gperftools Image

This opt-in example builds Envoy from git sources using the official Envoy
Ubuntu build image, adds `--define tcmalloc=gperftools` to the Bazel build, and
copies the resulting Envoy binary into a chiseled Wolfi runtime rootfs.

Build from the repository root:

```sh
make -C examples/wolfi-envoy-gperftools build
```

Check the Envoy binary:

```sh
docker run --rm chisel-wolfi:envoy-gperftools --version
```

Run Envoy with the included default config:

```sh
docker run --rm -p 10000:10000 -p 9901:9901 chisel-wolfi:envoy-gperftools
```

The default build fetches `v1.38.0` from `https://github.com/envoyproxy/envoy.git`
and uses `envoyproxy/envoy-build-ubuntu:86873047235e9b8232df989a5999b9bebf9db69c`.
Override the git ref, source repository, build image, or Bazel parallelism with
Docker build args:

```sh
make -C examples/wolfi-envoy-gperftools build BUILD_ARGS='--build-arg ENVOY_GIT_REF=main --build-arg BAZEL_JOBS=8'
```

The Dockerfile defaults to `linux/arm64` when using Docker's legacy builder,
which does not provide BuildKit's automatic platform args. For an amd64 legacy
build, pass matching platform and arch args:

```sh
make -C examples/wolfi-envoy-gperftools build BUILD_ARGS='--build-arg BUILDPLATFORM=linux/amd64 --build-arg TARGETPLATFORM=linux/amd64 --build-arg BUILDARCH=amd64 --build-arg TARGETARCH=amd64'
```

Envoy source builds are large and slow, so this example is intentionally not
part of the default `examples/Makefile` suite or publish workflow. The final
image keeps Wolfi's `os-release`, APK database, Envoy config, and CA
certificates, but the Envoy binary is built from git instead of cut from the
Wolfi Envoy APK package.
