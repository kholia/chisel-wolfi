# Publishing Wolfi Images

The recommended tag format for test and launch images is:

```text
YYYY-MM-DD_N
```

For example:

```text
2026-05-15_1
2026-05-15_2
```

This gives us an immutable, human-readable build tag that lines up with the
rolling-release model. The date is the UTC publish date, and the numeric suffix
is incremented when we republish images on the same date. Do not publish extra
floating tags or SHA tags for these images.

## GHCR Test Images

Use the `Publish Wolfi Images` workflow from GitHub Actions. It publishes images
to GitHub Container Registry using this image naming pattern:

```text
ghcr.io/<owner>/chisel-wolfi-<example-name>:<tag>
```

For example, publishing `wolfi-vault` with `tag_date=2026-05-15` and
`tag_revision=1` creates:

```text
ghcr.io/<owner>/chisel-wolfi-vault:2026-05-15_1
```

The workflow supports these image sets:

- `smoke`: curl, wget, vault-compatible OpenBao, and Envoy
- `popular`: the smoke set plus stunnel, Unbound, network tools, AWS CLI, and Terraform
- `all`: every standalone Dockerfile example

The default platform set is `linux/amd64,linux/arm64`. Public Wolfi does not
currently publish a `riscv64` APK index at `https://packages.wolfi.dev/os`, so
the GHCR workflow does not publish `linux/riscv64` images yet.

The workflow disables Buildx provenance and SBOM attestations for these test
publishes. Otherwise GHCR may show an extra `unknown/unknown` manifest entry
next to the real `linux/amd64` and `linux/arm64` images.

The example Dockerfiles build and run the temporary `chisel` tool on the Buildx
builder platform, while using `TARGETARCH` only to select the Wolfi APK
architecture to cut. This avoids running `chisel` under QEMU during local
multi-platform builds.

After publishing, the workflow and helper script inspect the pushed manifest and
fail if `unknown/unknown` is present or if any requested platform is missing.
If a tag was published before this check existed, republish the same
`YYYY-MM-DD_N` tag with the current script or workflow so the tag points at a
clean manifest list.

## Manual GHCR Push

Use the helper script from the repository root when you want to publish one image
manually:

```sh
GHCR_TOKEN="<pat-with-write:packages>" \
  scripts/publish-wolfi-image.sh \
  --owner "<github-user-or-org>" \
  --example wolfi-vault
```

For a multi-platform push:

```sh
GHCR_TOKEN="<pat-with-write:packages>" \
  scripts/publish-wolfi-image.sh \
  --owner "<github-user-or-org>" \
  --example wolfi-vault \
  --multi-platform
```

The script supports `--tag`, `--revision`, `--platforms`, `--skip-scan`,
`--skip-native-test`, and `--dry-run`.

The helper script and workflow both read Docker build arguments from
`examples/wolfi-versions.mk`, which is the central place for helper image tags
and versioned Wolfi package slice names.

The equivalent manual commands are below.

Set the target repository, example, and immutable tag:

```sh
export GHCR_OWNER="<github-user-or-org>"
export GHCR_TOKEN="<pat-with-write:packages>"
export EXAMPLE="wolfi-vault"
export TAG="$(date -u +%F)_1"
export IMAGE_NAME="${EXAMPLE#wolfi-}"
export IMAGE="ghcr.io/$GHCR_OWNER/chisel-wolfi-$IMAGE_NAME"
```

Log in to GitHub Container Registry:

```sh
echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_OWNER" --password-stdin
```

Build, test, size, scan, and push the native local image:

```sh
make -C "examples/$EXAMPLE" IMAGE="$IMAGE:$TAG" build test size scan
docker push "$IMAGE:$TAG"
```

For a multi-platform image, use Buildx:

```sh
docker buildx create --use --name chisel-wolfi-builder 2>/dev/null || \
  docker buildx use chisel-wolfi-builder
mapfile -t wolfi_build_args < <(make -s -f examples/wolfi-versions.mk print-build-args)

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  "${wolfi_build_args[@]}" \
  -f "examples/$EXAMPLE/Dockerfile" \
  -t "$IMAGE:$TAG" \
  --provenance=false \
  --sbom=false \
  --label "org.opencontainers.image.source=https://github.com/$GHCR_OWNER/chisel" \
  --label "org.opencontainers.image.revision=$(git rev-parse HEAD)" \
  --label "org.opencontainers.image.version=$TAG" \
  --push \
  .
```

For example, with `EXAMPLE=wolfi-vault` and `TAG=2026-05-15_1`, the pushed
references are:

```text
ghcr.io/<owner>/chisel-wolfi-vault:2026-05-15_1
```

Use another standalone Dockerfile example, such as `wolfi-curl`, `wolfi-envoy`,
or `wolfi-aws-cli`, by changing `EXAMPLE`.
