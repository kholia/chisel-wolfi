# syntax=docker/dockerfile:1

# Standard boilerplate: build the local chisel binary.
ARG CHISEL_BUILDER_IMAGE=golang:alpine
ARG WOLFI_IMAGE=cgr.dev/chainguard/wolfi-base

FROM --platform=$BUILDPLATFORM ${CHISEL_BUILDER_IMAGE} AS chisel
ARG BUILDARCH
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY public ./public
RUN CGO_ENABLED=0 GOOS=linux GOARCH="${BUILDARCH}" \
    go build -trimpath -ldflags="-s -w" -o /chisel ./cmd/chisel

# Standard boilerplate: prepare a scratch rootfs from Wolfi slices.
FROM --platform=$BUILDPLATFORM ${WOLFI_IMAGE} AS rootfs
ARG TARGETARCH

RUN apk add --no-cache ca-certificates-bundle

COPY --from=chisel /chisel /usr/local/bin/chisel
COPY examples/wolfi-release/chisel.yaml /release/chisel.yaml
COPY slices /release/slices

# Unbound-specific bits below: the `unbound_bins` slice and its state dir.
RUN set -eux; \
    case "${TARGETARCH}" in \
        amd64) apk_arch=x86_64 ;; \
        arm64) apk_arch=aarch64 ;; \
        *) echo "unsupported TARGETARCH=${TARGETARCH}" >&2; exit 1 ;; \
    esac; \
    mkdir /rootfs; \
    chisel cut --release /release --root /rootfs --arch "${apk_arch}" \
        wolfi-baselayout_base \
        unbound_bins; \
    mkdir -p /rootfs/var/lib/unbound; \
    chown 65532:65532 /rootfs/var/lib/unbound; \
    test -x /rootfs/usr/bin/unbound; \
    test -x /rootfs/usr/bin/unbound-checkconf; \
    test -f /rootfs/etc/unbound/unbound.conf

# Standard boilerplate: copy the cut rootfs into an empty final image.
FROM scratch

COPY --from=rootfs /rootfs /

# Unbound's runtime surface: DNS on port 53, foreground mode, default config.
EXPOSE 53/tcp 53/udp

ENTRYPOINT ["/usr/bin/unbound"]
CMD ["-d", "-c", "/etc/unbound/unbound.conf"]
