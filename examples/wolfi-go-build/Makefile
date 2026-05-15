EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

IMAGE ?= chisel-wolfi:go-build
DOCKERFILE ?= $(EXAMPLE_DIR)/Dockerfile
CONTEXT ?= $(REPO_ROOT)

include $(EXAMPLE_DIR)/../common.mk

test: build
	docker run --rm "$(IMAGE)" -c 'go version && go env GOROOT && gcc --version >/dev/null && g++ --version >/dev/null && git --version && unzip -v >/dev/null && grpc_health_probe -version'
