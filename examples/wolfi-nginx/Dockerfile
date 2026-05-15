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
ARG WOLFI_NGINX_SLICE
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
        wolfi-baselayout_base \
        ca-certificates-bundle_certs \
        ${WOLFI_NGINX_SLICE}_bins \
        ${WOLFI_NGINX_SLICE}_html \
        ${WOLFI_NGINX_SLICE}-config_config; \
    printf '%s\n' 'nginx:x:65532:' >> /staging-rootfs/etc/group; \
    printf '%s\n' 'nginx:x:65532:65532:nginx:/var/lib/nginx:/sbin/nologin' >> /staging-rootfs/etc/passwd; \
    mkdir -p /staging-rootfs/var/lib/nginx/logs /staging-rootfs/var/lib/nginx/tmp /staging-rootfs/run; \
    test -x /staging-rootfs/usr/bin/nginx; \
    test -f /staging-rootfs/etc/nginx/nginx.conf; \
    test -f /staging-rootfs/etc/nginx/mime.types; \
    test -f /staging-rootfs/var/lib/nginx/html/index.html; \
    test -f /staging-rootfs/etc/os-release; \
    test -f /staging-rootfs/lib/apk/db/installed

FROM scratch

COPY --from=builder /staging-rootfs /

EXPOSE 80
ENTRYPOINT ["/usr/bin/nginx"]
CMD ["-g", "daemon off;"]
