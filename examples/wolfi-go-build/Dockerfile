# syntax=docker/dockerfile:1

ARG CHISEL_BUILDER_IMAGE=golang:alpine
ARG WOLFI_IMAGE=cgr.dev/chainguard/wolfi-base

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
        busybox_shell \
        bash_bins \
        coreutils_bins \
        findutils_bins \
        grep_bins \
        sed_bins \
        gnutar_bins \
        gzip_bins \
        unzip_bins \
        git_bins \
        curl_bins \
        ${WOLFI_GO_SLICE}_bins \
        build-base_dev \
        grpc-health-probe_bins; \
    mkdir -p /staging-rootfs/go /staging-rootfs/work; \
    test -L /staging-rootfs/usr/bin/go; \
    test -x /staging-rootfs/usr/lib/go/bin/go; \
    test -x /staging-rootfs/usr/bin/gcc; \
    test -x /staging-rootfs/usr/bin/g++; \
    test -x /staging-rootfs/usr/bin/git; \
    test -x /staging-rootfs/usr/bin/unzip; \
    test -x /staging-rootfs/usr/bin/grpc-health-probe; \
    test -f /staging-rootfs/etc/ssl/certs/ca-certificates.crt; \
    test -f /staging-rootfs/etc/os-release; \
    test -f /staging-rootfs/lib/apk/db/installed

FROM scratch

COPY --from=builder /staging-rootfs /

ENV GOROOT=/usr/lib/go \
    GOPATH=/go \
    GOTOOLCHAIN=local \
    PATH=/go/bin:/usr/lib/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /work
ENTRYPOINT ["/bin/sh"]
