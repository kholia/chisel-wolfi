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
ARG WOLFI_SLICES
SHELL ["/bin/sh", "-eux", "-c"]

RUN apk add --no-cache ca-certificates-bundle

COPY --from=chisel-build /out/chisel /usr/local/bin/chisel
COPY examples/wolfi-release/chisel.yaml /chisel-release/chisel.yaml
COPY slices /chisel-release/slices

RUN test -n "${WOLFI_SLICES}"; \
    case "${TARGETARCH}" in \
        amd64) APK_ARCH=x86_64 ;; \
        arm64) APK_ARCH=aarch64 ;; \
        riscv64) APK_ARCH=riscv64 ;; \
        *) echo "unsupported TARGETARCH: ${TARGETARCH}" >&2; exit 1 ;; \
    esac; \
    mkdir /staging-rootfs; \
    chisel cut --release /chisel-release --root /staging-rootfs --arch "${APK_ARCH}" ${WOLFI_SLICES}; \
    mkdir -p /staging-rootfs/work /staging-rootfs/go /staging-rootfs/app; \
    if [ -x "/staging-rootfs${WOLFI_PYTHON_BIN}" ]; then ln -sf "${WOLFI_PYTHON_BIN}" /staging-rootfs/usr/bin/python3; fi; \
    if [ ! -e /staging-rootfs/usr/bin/python3 ] && [ -x "/staging-rootfs${WOLFI_FALLBACK_PYTHON_BIN}" ]; then ln -sf "${WOLFI_FALLBACK_PYTHON_BIN}" /staging-rootfs/usr/bin/python3; fi; \
    if [ -x /staging-rootfs/usr/bin/pysemver ]; then ln -sf /usr/bin/pysemver /staging-rootfs/usr/bin/semver; fi; \
    if [ -x "/staging-rootfs${WOLFI_JAVA_HOME}/bin/java" ]; then \
        ln -sf "${WOLFI_JAVA_HOME}/bin/java" /staging-rootfs/usr/bin/java; \
        ln -sf "${WOLFI_JAVA_HOME}/bin/jar" /staging-rootfs/usr/bin/jar; \
        ln -sf "${WOLFI_JAVA_HOME}/bin/keytool" /staging-rootfs/usr/bin/keytool; \
    fi; \
    if [ -x "/staging-rootfs${WOLFI_JAVA_HOME}/bin/javac" ]; then \
        ln -sf "${WOLFI_JAVA_HOME}/bin/javac" /staging-rootfs/usr/bin/javac; \
        ln -sf "${WOLFI_JAVA_HOME}/bin/jlink" /staging-rootfs/usr/bin/jlink; \
    fi; \
    test -f /staging-rootfs/etc/os-release; \
    test -f /staging-rootfs/lib/apk/db/installed

FROM scratch

COPY --from=builder /staging-rootfs /

ARG WOLFI_CLOUDSDK_PYTHON_BIN
ARG WOLFI_JAVA_HOME
ARG WOLFI_MAVEN_HOME

ENV GOROOT=/usr/lib/go \
    GOPATH=/go \
    GOTOOLCHAIN=local \
    JAVA_HOME=${WOLFI_JAVA_HOME} \
    MAVEN_HOME=${WOLFI_MAVEN_HOME} \
    CLOUDSDK_PYTHON=${WOLFI_CLOUDSDK_PYTHON_BIN} \
    PATH=/go/bin:/usr/lib/go/bin:${WOLFI_MAVEN_HOME}/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /work
ENTRYPOINT ["/bin/sh"]
