EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:common-utility-base
COMMON_SLICES := $(COMMON_UBUNTU_BASE_SLICES)

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'wget --version >/dev/null && grpc_health_probe -version'
