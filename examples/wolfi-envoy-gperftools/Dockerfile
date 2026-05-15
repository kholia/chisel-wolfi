# syntax=docker/dockerfile:1

ARG CHISEL_BUILDER_IMAGE=golang:alpine
ARG ENVOY_BUILD_IMAGE=envoyproxy/envoy-build-ubuntu:86873047235e9b8232df989a5999b9bebf9db69c
ARG WOLFI_IMAGE=cgr.dev/chainguard/wolfi-base
ARG BUILDPLATFORM=linux/arm64
ARG TARGETPLATFORM=linux/arm64
ARG BUILDARCH=arm64
ARG TARGETARCH=arm64

FROM --platform=$BUILDPLATFORM ${CHISEL_BUILDER_IMAGE} AS chisel-build
ARG BUILDARCH
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY public ./public

RUN CGO_ENABLED=0 GOOS=linux GOARCH="${BUILDARCH}" \
    go build -trimpath -ldflags="-s -w" -o /out/chisel ./cmd/chisel

FROM --platform=$TARGETPLATFORM ${ENVOY_BUILD_IMAGE} AS envoy-build
ARG ENVOY_GIT_REF=v1.38.0
ARG ENVOY_REPOSITORY=https://github.com/envoyproxy/envoy.git
ARG BAZEL_JOBS=
SHELL ["/bin/bash", "-euxo", "pipefail", "-c"]

RUN if ! getent passwd 1000 >/dev/null; then \
        useradd -m -d /build -u 1000 -s /bin/bash envoy-build; \
    fi; \
    mkdir -p /build /out /source; \
    chown -R 1000:1000 /build /out /source

USER 1000
ENV HOME=/build
WORKDIR /source

RUN git init .; \
    git remote add origin "${ENVOY_REPOSITORY}"; \
    git fetch --depth=1 origin "${ENVOY_GIT_REF}"; \
    git checkout --detach FETCH_HEAD

RUN { \
        printf '%s\n' 'build --define tcmalloc=gperftools'; \
        if [[ -n "${BAZEL_JOBS}" ]]; then \
            printf 'build --jobs=%s\n' "${BAZEL_JOBS}"; \
        fi; \
    } > user.bazelrc; \
    ./ci/do_ci.sh release.server_only; \
    envoy_binary="$(find "${HOME}/.cache/envoy-bazel/envoy" -path '*/source/exe/envoy-static' -type f -perm -0100 | head -n 1)"; \
    test -n "${envoy_binary}"; \
    install -m 0555 "${envoy_binary}" /out/envoy; \
    /out/envoy --version

FROM --platform=$BUILDPLATFORM ${WOLFI_IMAGE} AS builder
ARG TARGETARCH
ARG WOLFI_ENVOY_CONFIG_SLICE
ARG WOLFI_ENVOY_SLICE
ARG WOLFI_GO_SLICE
ARG WOLFI_JAVA_HOME
ARG WOLFI_JAVA_SLICE
ARG WOLFI_MAVEN_HOME
ARG WOLFI_MAVEN_SLICE
ARG WOLFI_NODE_MINIMAL_SLICE
ARG WOLFI_NODE_SLICE
ARG WOLFI_PYTHON_BIN
ARG WOLFI_FALLBACK_PYTHON_BIN
ARG WOLFI_PYTHON_SLICE
ARG WOLFI_CLOUDSDK_PYTHON_BIN
SHELL ["/bin/sh", "-eux", "-c"]

RUN apk add --no-cache ca-certificates-bundle

COPY --from=chisel-build /out/chisel /usr/local/bin/chisel
COPY examples/wolfi-release/chisel.yaml /chisel-release/chisel.yaml
COPY slices /chisel-release/slices

RUN case "${TARGETARCH}" in \
        amd64) APK_ARCH=x86_64 ;; \
        arm64) APK_ARCH=aarch64 ;; \
        *) echo "unsupported TARGETARCH: ${TARGETARCH}" >&2; exit 1 ;; \
    esac; \
    mkdir /staging-rootfs; \
    chisel cut --release /chisel-release --root /staging-rootfs --arch "${APK_ARCH}" \
        wolfi-baselayout_base-library \
        ca-certificates-bundle_certs \
        ${WOLFI_ENVOY_CONFIG_SLICE}_config; \
    mkdir -p /staging-rootfs/usr/bin; \
    test -f /staging-rootfs/etc/envoy/envoy.yaml; \
    test -f /staging-rootfs/etc/os-release; \
    test -f /staging-rootfs/lib/apk/db/installed

COPY --from=envoy-build --chown=0:0 /out/envoy /staging-rootfs/usr/bin/envoy

RUN test -x /staging-rootfs/usr/bin/envoy

FROM scratch

COPY --from=builder /staging-rootfs /

EXPOSE 10000 9901

ENTRYPOINT ["/usr/bin/envoy"]
CMD ["-c", "/etc/envoy/envoy.yaml"]
