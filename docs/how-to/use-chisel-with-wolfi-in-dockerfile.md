# Use Chisel With Wolfi In A Dockerfile

Chisel can cut packages from APK repositories such as Wolfi. In a Dockerfile,
the usual pattern is a multi-stage build:

1. Build or install `chisel`.
2. Create a local Chisel release with `kind: apk`.
3. Cut Wolfi slices into a staging root file system.
4. Copy the staging root into a `scratch` final image.

The working example lives in `examples/wolfi-dockerfile/`.

For the security model behind these examples, including the rolling-release
zero-CVE philosophy and the boundary for FIPS and STIG support, see
[Wolfi Security Model](../wolfi-security-model.md).

## Build The Image

Run this from the repository root:

```sh
make -C examples/wolfi-dockerfile build
```

The Makefile passes the shared build arguments from
`examples/wolfi-versions.mk`. The Dockerfile builds Chisel from this source tree
so that the example uses the same APK support and Wolfi slices as the checkout.

## Release Configuration

The example uses this local Chisel release:

```yaml
format: v3

maintenance:
  standard: 2025-01-01
  end-of-life: 2100-01-01

archives:
  wolfi:
    kind: apk
    url: https://packages.wolfi.dev/os
```

The `kind: apk` archive tells Chisel to read `APKINDEX.tar.gz` and `.apk`
packages from the configured Wolfi repository.

## Select Slices

By default, the Dockerfile cuts a small Wolfi shell root:

```text
wolfi-baselayout_base-library busybox_shell
```

The `wolfi-baselayout_base-library` slice provides Wolfi's merged directory
layout and common runtime libraries. The `busybox_shell` slice adds `/bin/sh`
through Wolfi's `/bin -> /usr/bin` layout.

Use `CHISEL_SLICES` to build a different image:

```sh
docker build \
  $(make -s -f examples/wolfi-versions.mk print-build-args | xargs) \
  -f examples/wolfi-dockerfile/Dockerfile \
  --build-arg CHISEL_SLICES="wolfi-baselayout_base-library busybox_shell stunnel_bins unbound_bins" \
  -t chisel-wolfi:network-tools .
```

## Run The Image

```sh
docker run --rm chisel-wolfi:busybox -c 'echo hello from Wolfi; busybox uname -m'
```

For multi-platform builds, Docker provides `TARGETARCH` values such as `amd64`
and `arm64`. Wolfi APK indexes use `x86_64` and `aarch64`, so the example maps
the architecture before calling `chisel cut --arch`.

The default image includes `/bin/sh`, but it does not run APK triggers to create
every BusyBox applet symlink. Use `busybox <applet>` for applets that are not
present as standalone commands.
