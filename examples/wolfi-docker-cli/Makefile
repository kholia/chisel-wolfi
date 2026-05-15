EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:docker-cli
COMMON_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	bash_bins \
	docker-cli_bins \
	docker-cli-buildx_plugin \
	docker-compose_bins

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'docker --version && docker buildx version && docker-compose version'
