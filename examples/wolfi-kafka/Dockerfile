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
ARG WOLFI_KAFKA_SLICE
ARG WOLFI_JAVA_HOME
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
        ${WOLFI_KAFKA_SLICE}_runtime; \
    mkdir -p /staging-rootfs/var/lib/kafka /staging-rootfs/tmp; \
    ln -sf "${WOLFI_JAVA_HOME}/bin/java" /staging-rootfs/usr/bin/java; \
    test -x /staging-rootfs/usr/lib/kafka/bin/kafka-server-start.sh; \
    test -x /staging-rootfs/usr/lib/kafka/bin/kafka-topics.sh; \
    test -x "/staging-rootfs${WOLFI_JAVA_HOME}/bin/java"; \
    test -f /staging-rootfs/usr/lib/kafka/config/server.properties.original; \
    test -f /staging-rootfs/etc/os-release; \
    test -f /staging-rootfs/lib/apk/db/installed

FROM scratch

COPY --from=builder /staging-rootfs /

ARG WOLFI_JAVA_HOME
ENV JAVA_HOME=${WOLFI_JAVA_HOME} \
    PATH=${WOLFI_JAVA_HOME}/bin:/usr/lib/kafka/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

EXPOSE 9092
ENTRYPOINT ["/usr/lib/kafka/bin/kafka-server-start.sh"]
CMD ["/usr/lib/kafka/config/server.properties.original"]
