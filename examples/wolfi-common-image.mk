DOCKERFILE ?= $(EXAMPLE_DIR)/../wolfi-common-image.Dockerfile
CONTEXT ?= $(REPO_ROOT)
BUILD_ARGS ?= --build-arg WOLFI_SLICES='$(COMMON_SLICES)'

include $(EXAMPLE_DIR)/../common.mk
