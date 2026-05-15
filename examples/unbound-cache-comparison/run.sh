#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
ARTIFACT_DIR="${SCRIPT_DIR}/artifacts"
IMAGE_DIR="${ARTIFACT_DIR}/images"

PLATFORM="${PLATFORM:-linux/arm64}"
BUILDX_BUILDER="${BUILDX_BUILDER:-colima}"
UBUNTU_2404_IMAGE="${UBUNTU_2404_IMAGE:-chisel-demo:ubuntu-2404-unbound-cache}"
UBUNTU_2604_IMAGE="${UBUNTU_2604_IMAGE:-${UBUNTU_IMAGE:-chisel-demo:ubuntu-2604-unbound-cache}}"
CHISEL_IMAGE="${CHISEL_IMAGE:-chisel-demo:chisel-wolfi-unbound-cache}"
CHISEL_SLICES="${CHISEL_SLICES:-wolfi-baselayout_base-library ca-certificates-bundle_certs bind-tools_bins unbound_bins unbound-config_config}"

mkdir -p "${ARTIFACT_DIR}" "${IMAGE_DIR}"

run_log() {
    local log_file="$1"
    shift
    "$@" 2>&1 | tee "${log_file}"
}

write_versions() {
    {
        printf 'Generated: '
        date -u '+%Y-%m-%dT%H:%M:%SZ'
        printf 'Repository: '
        git -C "${REPO_ROOT}" rev-parse --show-toplevel
        printf 'Commit: '
        git -C "${REPO_ROOT}" rev-parse HEAD
        printf 'Platform: %s\n' "${PLATFORM}"
        printf 'Buildx builder: %s\n\n' "${BUILDX_BUILDER}"
        docker version
        printf '\n'
        docker-buildx version
        printf '\n'
        grype version
        printf '\nWolfi slices:\n%s\n' "${CHISEL_SLICES}"
    } > "${ARTIFACT_DIR}/versions.txt"
}

build_ubuntu() {
    local version="$1"
    local image="$2"

    run_log "${ARTIFACT_DIR}/ubuntu-${version}-unbound-cache.build.log" \
        docker-buildx build \
            --builder "${BUILDX_BUILDER}" \
            --platform "${PLATFORM}" \
            --progress=plain \
            --pull \
            --provenance=false \
            --sbom=false \
            --load \
            --metadata-file "${ARTIFACT_DIR}/ubuntu-${version}-unbound-cache.build-metadata.json" \
            -f "${SCRIPT_DIR}/Dockerfile.ubuntu-${version}" \
            -t "${image}" \
            "${REPO_ROOT}"
}

build_chisel() {
    local wolfi_build_args=()

    while IFS= read -r arg; do
        wolfi_build_args+=("${arg}")
    done < <(make -s -C "${REPO_ROOT}" -f examples/wolfi-versions.mk print-build-args)

    run_log "${ARTIFACT_DIR}/chisel-wolfi-unbound-cache.build.log" \
        docker-buildx build \
            --builder "${BUILDX_BUILDER}" \
            --platform "${PLATFORM}" \
            --progress=plain \
            --pull \
            --provenance=false \
            --sbom=false \
            --load \
            --metadata-file "${ARTIFACT_DIR}/chisel-wolfi-unbound-cache.build-metadata.json" \
            "${wolfi_build_args[@]}" \
            --build-arg "CHISEL_SLICES=${CHISEL_SLICES}" \
            -f "${SCRIPT_DIR}/Dockerfile.chisel-wolfi" \
            -t "${CHISEL_IMAGE}" \
            "${REPO_ROOT}"
}

smoke_test() {
    {
        smoke_test_ubuntu "${UBUNTU_2404_IMAGE}"
        printf '\n'
        smoke_test_ubuntu "${UBUNTU_2604_IMAGE}"

        printf '\n$ docker run --rm --entrypoint /usr/bin/unbound %s -V\n' "${CHISEL_IMAGE}"
        docker run --rm --entrypoint /usr/bin/unbound "${CHISEL_IMAGE}" -V
        printf '\n$ docker run --rm --entrypoint /usr/bin/unbound-checkconf %s /etc/unbound/unbound.conf\n' "${CHISEL_IMAGE}"
        docker run --rm --entrypoint /usr/bin/unbound-checkconf "${CHISEL_IMAGE}" /etc/unbound/unbound.conf
        printf '\n$ docker run --rm --entrypoint /usr/bin/dig %s -v\n' "${CHISEL_IMAGE}"
        docker run --rm --entrypoint /usr/bin/dig "${CHISEL_IMAGE}" -v
    } > "${ARTIFACT_DIR}/smoke-tests.txt" 2>&1
}

smoke_test_ubuntu() {
    local image="$1"

    printf '$ docker run --rm --entrypoint /usr/sbin/unbound %s -V\n' "${image}"
    docker run --rm --entrypoint /usr/sbin/unbound "${image}" -V
    printf '\n$ docker run --rm --entrypoint /usr/sbin/unbound-checkconf %s /etc/unbound/unbound.conf\n' "${image}"
    docker run --rm --entrypoint /usr/sbin/unbound-checkconf "${image}" /etc/unbound/unbound.conf
    printf '\n$ docker run --rm --entrypoint /usr/bin/dig %s -v\n' "${image}"
    docker run --rm --entrypoint /usr/bin/dig "${image}" -v
}

save_metadata() {
    docker image inspect "${UBUNTU_2404_IMAGE}" > "${ARTIFACT_DIR}/ubuntu-24.04-unbound-cache.inspect.json"
    docker image inspect "${UBUNTU_2604_IMAGE}" > "${ARTIFACT_DIR}/ubuntu-26.04-unbound-cache.inspect.json"
    docker image inspect "${CHISEL_IMAGE}" > "${ARTIFACT_DIR}/chisel-wolfi-unbound-cache.inspect.json"

    {
        printf 'image\tid\tcreated\tbytes\n'
        docker image inspect "${UBUNTU_2404_IMAGE}" \
            --format '{{index .RepoTags 0}}	{{.Id}}	{{.Created}}	{{.Size}}'
        docker image inspect "${UBUNTU_2604_IMAGE}" \
            --format '{{index .RepoTags 0}}	{{.Id}}	{{.Created}}	{{.Size}}'
        docker image inspect "${CHISEL_IMAGE}" \
            --format '{{index .RepoTags 0}}	{{.Id}}	{{.Created}}	{{.Size}}'
    } > "${ARTIFACT_DIR}/image-sizes.tsv"

    docker save -o "${IMAGE_DIR}/ubuntu-24.04-unbound-cache.tar" "${UBUNTU_2404_IMAGE}"
    gzip -c "${IMAGE_DIR}/ubuntu-24.04-unbound-cache.tar" > "${IMAGE_DIR}/ubuntu-24.04-unbound-cache.tar.gz"
    docker save -o "${IMAGE_DIR}/ubuntu-26.04-unbound-cache.tar" "${UBUNTU_2604_IMAGE}"
    gzip -c "${IMAGE_DIR}/ubuntu-26.04-unbound-cache.tar" > "${IMAGE_DIR}/ubuntu-26.04-unbound-cache.tar.gz"
    docker save -o "${IMAGE_DIR}/chisel-wolfi-unbound-cache.tar" "${CHISEL_IMAGE}"
    gzip -c "${IMAGE_DIR}/chisel-wolfi-unbound-cache.tar" > "${IMAGE_DIR}/chisel-wolfi-unbound-cache.tar.gz"
}

scan_image() {
    local archive="$1"
    local name="$2"

    grype "docker-archive:${archive}" -o table > "${ARTIFACT_DIR}/${name}.grype.txt"
    grype "docker-archive:${archive}" -o json > "${ARTIFACT_DIR}/${name}.grype.json"
}

write_scan_summary() {
    {
        printf 'image\tmatches\tcritical\thigh\tmedium\tlow\tnegligible\tunknown\tfixed\tnot-fixed\twont-fix\tunknown-fix\n'
        scan_summary_row "${UBUNTU_2404_IMAGE}" "${ARTIFACT_DIR}/ubuntu-24.04-unbound-cache.grype.json"
        scan_summary_row "${UBUNTU_2604_IMAGE}" "${ARTIFACT_DIR}/ubuntu-26.04-unbound-cache.grype.json"
        scan_summary_row "${CHISEL_IMAGE}" "${ARTIFACT_DIR}/chisel-wolfi-unbound-cache.grype.json"
    } > "${ARTIFACT_DIR}/grype-summary.tsv"
}

scan_summary_row() {
    local image="$1"
    local report="$2"

    jq -r --arg image "${image}" '
        def sevcount($s): ([.matches[]? | (.vulnerability.severity // "unknown" | ascii_downcase) | select(. == $s)] | length);
        def statecount($s): ([.matches[]? | (.vulnerability.fix.state // "unknown" | ascii_downcase) | select(. == $s)] | length);
        [$image, (.matches | length), sevcount("critical"), sevcount("high"), sevcount("medium"), sevcount("low"), sevcount("negligible"), sevcount("unknown"), statecount("fixed"), statecount("not-fixed"), statecount("wont-fix"), statecount("unknown")] | @tsv
    ' "${report}"
}

main() {
    write_versions
    build_ubuntu "24.04" "${UBUNTU_2404_IMAGE}"
    build_ubuntu "26.04" "${UBUNTU_2604_IMAGE}"
    build_chisel
    smoke_test
    save_metadata
    scan_image "${IMAGE_DIR}/ubuntu-24.04-unbound-cache.tar" "ubuntu-24.04-unbound-cache"
    scan_image "${IMAGE_DIR}/ubuntu-26.04-unbound-cache.tar" "ubuntu-26.04-unbound-cache"
    scan_image "${IMAGE_DIR}/chisel-wolfi-unbound-cache.tar" "chisel-wolfi-unbound-cache"
    write_scan_summary
}

main "$@"
