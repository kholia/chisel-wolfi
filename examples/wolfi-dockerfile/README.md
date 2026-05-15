# Use Chisel With Wolfi In A Dockerfile

This example builds a `scratch` image from Wolfi APK packages using Chisel.
It builds the Chisel binary from this source tree, creates a local Wolfi
release from `examples/wolfi-release/`, cuts the
selected slices into `/staging-rootfs`, then copies that root file system into
the final image.

Build from the repository root:

```sh
make -C examples/wolfi-dockerfile build
```

Run the chiseled Wolfi image:

```sh
docker run --rm chisel-wolfi:busybox -c 'echo hello from Wolfi; busybox uname -m'
```

The default slice set is:

```text
wolfi-baselayout_base-library busybox_shell
```

To include more Wolfi slices, override `CHISEL_SLICES`:

```sh
docker build \
  $(make -s -f examples/wolfi-versions.mk print-build-args | xargs) \
  -f examples/wolfi-dockerfile/Dockerfile \
  --build-arg CHISEL_SLICES="wolfi-baselayout_base-library busybox_shell stunnel_bins unbound_bins" \
  -t chisel-wolfi:network-tools .
```

For multi-platform builds, Docker exposes `TARGETARCH` as `amd64` or `arm64`.
The Dockerfile maps those values to Wolfi APK architecture names, `x86_64` and
`aarch64`, before running `chisel cut`.

The default image includes `/bin/sh`, but it does not run APK triggers to create
every BusyBox applet symlink. Use `busybox <applet>` for applets that are not
present as standalone commands.
