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
ARG UPX_FLAGS=--best
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

RUN apk add --no-cache ca-certificates-bundle upx

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
        ${WOLFI_ENVOY_SLICE}_bins \
        ${WOLFI_ENVOY_CONFIG_SLICE}_config; \
    test -x /staging-rootfs/usr/bin/envoy; \
    chmod u+w /staging-rootfs/usr/bin/envoy; \
    upx ${UPX_FLAGS} /staging-rootfs/usr/bin/envoy; \
    upx -t /staging-rootfs/usr/bin/envoy; \
    chmod 0555 /staging-rootfs/usr/bin/envoy; \
    test -x /staging-rootfs/usr/bin/envoy; \
    test -f /staging-rootfs/etc/envoy/envoy.yaml; \
    test -f /staging-rootfs/etc/os-release; \
    test -f /staging-rootfs/lib/apk/db/installed

FROM scratch

COPY --from=builder /staging-rootfs /

EXPOSE 10000 9901

ENTRYPOINT ["/usr/bin/envoy"]
CMD ["-c", "/etc/envoy/envoy.yaml"]
