# Unbound Cache Comparison

This comparison builds three DNS cache images:

- `chisel-demo:ubuntu-2404-unbound-cache`, a conventional Ubuntu 24.04 image
  with `unbound` and DNS tooling installed by `apt`.
- `chisel-demo:ubuntu-2604-unbound-cache`, a conventional Ubuntu 26.04 image
  with `unbound` and DNS tooling installed by `apt`.
- `chisel-demo:chisel-wolfi-unbound-cache`, a `scratch` image cut from Wolfi APK
  packages by this repository. It keeps Unbound, its config, CA certificates,
  and `dig`, but no shell.

Regenerate every artifact from the repository root:

```sh
examples/unbound-cache-comparison/run.sh
```

Build only the Chisel/Wolfi Unbound cache image with the same defaults:

```sh
docker-buildx build \
    --builder colima \
    --platform linux/arm64 \
    --progress=plain \
    --pull \
    --provenance=false \
    --sbom=false \
    --load \
    --metadata-file examples/unbound-cache-comparison/artifacts/chisel-wolfi-unbound-cache.build-metadata.json \
    $(make -s -f examples/wolfi-versions.mk print-build-args) \
    --build-arg "CHISEL_SLICES=wolfi-baselayout_base-library ca-certificates-bundle_certs bind-tools_bins unbound_bins unbound-config_config" \
    -f examples/unbound-cache-comparison/Dockerfile.chisel-wolfi \
    -t chisel-demo:chisel-wolfi-unbound-cache \
    .
```

```
export DOCKER_HOST="unix://$HOME/.colima/default/docker.sock"

grype chisel-wolfi:wolfi-unbound-publish-test
```

The script saves build logs, image metadata, smoke-test output, image sizes, and
Grype reports under `artifacts/`. `docker save` archives and compressed copies
are written to `artifacts/images/` and ignored by Git because they are large
generated files.
