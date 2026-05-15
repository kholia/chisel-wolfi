#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Publish one Wolfi example image to GitHub Container Registry.

Usage:
  scripts/publish-wolfi-image.sh --owner <github-user-or-org> --example <wolfi-example> [options]

Required:
  --owner OWNER          GitHub user or organization that owns the GHCR package.
  --example EXAMPLE      Example directory under examples/, such as wolfi-vault.

Options:
  --tag TAG              Immutable tag. Defaults to today's UTC date with _1.
  --revision N           Revision for the default date tag. Defaults to 1.
  --token TOKEN          GHCR token. Defaults to GHCR_TOKEN from the environment.
  --platforms LIST       Publish with buildx for these platforms.
                         Example: linux/amd64,linux/arm64
  --multi-platform       Shortcut for --platforms linux/amd64,linux/arm64.
  --skip-scan            Skip the native make scan target.
  --skip-native-test     Skip native build/test/size before publishing.
  --dry-run              Print commands without running them.
  -h, --help             Show this help.

Examples:
  GHCR_TOKEN=... scripts/publish-wolfi-image.sh --owner my-org --example wolfi-vault

  GHCR_TOKEN=... scripts/publish-wolfi-image.sh \
    --owner my-org \
    --example wolfi-envoy \
    --tag 2026-05-15_2 \
    --multi-platform
USAGE
}

die() {
  echo "error: $*" >&2
  exit 1
}

run() {
  printf '+'
  printf ' %q' "$@"
  printf '\n'
  if [ "$DRY_RUN" = "0" ]; then
    "$@"
  fi
}

require_cmd() {
  command -v "$1" >/dev/null || die "missing required command: $1"
}

load_wolfi_build_args() {
  WOLFI_BUILD_ARGS=()
  local arg
  while IFS= read -r arg; do
    WOLFI_BUILD_ARGS+=("$arg")
  done < <(make -s -f examples/wolfi-versions.mk print-build-args)
}

verify_multi_platform_manifest() {
  local image_ref="$1"
  local platforms="$2"

  if [ "$DRY_RUN" = "1" ]; then
    echo "+ docker buildx imagetools inspect $image_ref"
    echo "+ verify manifest contains only requested platforms: $platforms"
    return
  fi

  local inspect_output
  inspect_output="$(docker buildx imagetools inspect "$image_ref")"
  printf '%s\n' "$inspect_output"

  if printf '%s\n' "$inspect_output" | grep -q 'unknown/unknown'; then
    die "$image_ref contains an unknown/unknown manifest, usually from Buildx attestations"
  fi

  local platform
  for platform in ${platforms//,/ }; do
    if ! printf '%s\n' "$inspect_output" | grep -Eq "Platform:[[:space:]]+$platform([[:space:]]|$)"; then
      die "$image_ref is missing platform $platform"
    fi
  done
}

OWNER=""
EXAMPLE=""
TAG=""
REVISION="1"
TOKEN="${GHCR_TOKEN:-}"
PLATFORMS=""
SKIP_SCAN="0"
SKIP_NATIVE_TEST="0"
DRY_RUN="0"

while [ "$#" -gt 0 ]; do
  case "$1" in
    --owner)
      OWNER="${2:-}"
      shift 2
      ;;
    --example)
      EXAMPLE="${2:-}"
      shift 2
      ;;
    --tag)
      TAG="${2:-}"
      shift 2
      ;;
    --revision)
      REVISION="${2:-}"
      shift 2
      ;;
    --token)
      TOKEN="${2:-}"
      shift 2
      ;;
    --platforms)
      PLATFORMS="${2:-}"
      shift 2
      ;;
    --multi-platform)
      PLATFORMS="linux/amd64,linux/arm64"
      shift
      ;;
    --skip-scan)
      SKIP_SCAN="1"
      shift
      ;;
    --skip-native-test)
      SKIP_NATIVE_TEST="1"
      shift
      ;;
    --dry-run)
      DRY_RUN="1"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      die "unknown argument: $1"
      ;;
  esac
done

[ -n "$OWNER" ] || die "--owner is required"
[ -n "$EXAMPLE" ] || die "--example is required"
[ -d "examples/$EXAMPLE" ] || die "examples/$EXAMPLE does not exist"
[ -f "examples/$EXAMPLE/Dockerfile" ] || die "examples/$EXAMPLE/Dockerfile does not exist"

case "$EXAMPLE" in
  wolfi-*) ;;
  *) die "--example must use the wolfi-* example directory name" ;;
esac

case "$REVISION" in
  ''|*[!0-9]*) die "--revision must be a positive integer" ;;
  0) die "--revision must be greater than zero" ;;
esac

if [ -z "$TAG" ]; then
  TAG="$(date -u +%F)_$REVISION"
fi

if ! [[ "$TAG" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}_[0-9]+$ ]]; then
  die "--tag must use YYYY-MM-DD_N format"
fi

require_cmd docker
require_cmd git
require_cmd make
load_wolfi_build_args

IMAGE_NAME="${EXAMPLE#wolfi-}"
IMAGE="ghcr.io/$OWNER/chisel-wolfi-$IMAGE_NAME"
LOCAL_IMAGE="chisel-wolfi:${EXAMPLE}-publish-test"
FULL_SHA="$(git rev-parse HEAD)"

echo "Example:     $EXAMPLE"
echo "Image:       $IMAGE"
echo "Tag:         $TAG"
echo "Git SHA:     $FULL_SHA"
echo "Platforms:   ${PLATFORMS:-native}"

if [ "$DRY_RUN" = "0" ]; then
  [ -n "$TOKEN" ] || die "set GHCR_TOKEN or pass --token"
  printf '%s' "$TOKEN" | docker login ghcr.io -u "$OWNER" --password-stdin
else
  echo "+ docker login ghcr.io -u $OWNER --password-stdin"
fi

if [ "$SKIP_NATIVE_TEST" = "0" ]; then
  if [ "$SKIP_SCAN" = "0" ]; then
    run make -C "examples/$EXAMPLE" IMAGE="$LOCAL_IMAGE" build test size scan
  else
    run make -C "examples/$EXAMPLE" IMAGE="$LOCAL_IMAGE" build test size
  fi
fi

if [ -n "$PLATFORMS" ]; then
  if docker buildx inspect chisel-wolfi-builder >/dev/null 2>&1; then
    run docker buildx use chisel-wolfi-builder
  else
    run docker buildx create --use --name chisel-wolfi-builder
  fi
  run docker buildx build \
    --platform "$PLATFORMS" \
    "${WOLFI_BUILD_ARGS[@]}" \
    -f "examples/$EXAMPLE/Dockerfile" \
    -t "$IMAGE:$TAG" \
    --provenance=false \
    --sbom=false \
    --label "org.opencontainers.image.source=https://github.com/$OWNER/chisel" \
    --label "org.opencontainers.image.revision=$FULL_SHA" \
    --label "org.opencontainers.image.version=$TAG" \
    --push \
    .

  verify_multi_platform_manifest "$IMAGE:$TAG" "$PLATFORMS"
else
  run docker tag "$LOCAL_IMAGE" "$IMAGE:$TAG"
  run docker push "$IMAGE:$TAG"
fi

echo "Published:"
echo "  $IMAGE:$TAG"
