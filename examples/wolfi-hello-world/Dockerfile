# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/go@sha256:3d8fa865f7382d566f7a13aa67a90b5f758637ab75776ecb76e3da10df13ad66 AS chisel
ARG BUILDARCH
ARG TARGETARCH
ARG GOPROXY
ARG GOSUMDB
WORKDIR /src
COPY go.mod go.sum ./
RUN GOPROXY="${GOPROXY:?set GOPROXY to your approved Go module proxy}" \
    GOSUMDB="${GOSUMDB:?set GOSUMDB to your approved checksum database or off}" \
    go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY public ./public
RUN --network=none CGO_ENABLED=0 GOOS=linux GOARCH="$BUILDARCH" \
    go build -mod=readonly -trimpath -ldflags="-s -w" -o /chisel ./cmd/chisel
COPY examples/wolfi-hello-world/hello.go ./examples/wolfi-hello-world/hello.go
RUN --network=none CGO_ENABLED=0 GOOS=linux GOARCH="$TARGETARCH" \
    go build -mod=readonly -trimpath -ldflags="-s -w" -o /hello ./examples/wolfi-hello-world/hello.go

FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/wolfi-base@sha256:441d6709305552a3411e585ad98aacb9dadda00c80f3267483c38ac6f86f49d4 AS rootfs
ARG TARGETARCH
COPY --from=chisel /chisel /usr/local/bin/chisel
COPY examples/wolfi-release/chisel.yaml /release/chisel.yaml
COPY slices /release/slices
RUN set -eux; \
    case "$TARGETARCH" in \
        amd64) arch=x86_64 ;; \
        arm64) arch=aarch64 ;; \
        riscv64) arch=riscv64 ;; \
        *) echo "unsupported TARGETARCH: $TARGETARCH" >&2; exit 1 ;; \
    esac; \
    mkdir /rootfs; \
    chisel cut --release /release --root /rootfs --arch "$arch" \
        wolfi-baselayout_base; \
    mkdir -p /rootfs/usr/bin
COPY --from=chisel /hello /rootfs/usr/bin/hello

FROM scratch
COPY --from=rootfs /rootfs /
USER 65532:65532
CMD ["/usr/bin/hello"]
