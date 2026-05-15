EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:common-infra-toolbox-extended
COMMON_SLICES := $(COMMON_INFRA_BASE_SLICES)

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'aws --version && gcloud --version && kubectl version --client=true && terraform version && bao --version && vault version'
