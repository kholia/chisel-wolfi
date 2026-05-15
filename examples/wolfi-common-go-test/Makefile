EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:common-go-test
COMMON_SLICES := $(COMMON_GO_TEST_SLICES)

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'go version && gcc --version >/dev/null && protoc --version && grpc_health_probe -version'
