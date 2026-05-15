EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:common-aws-lambda-python312
COMMON_SLICES := $(COMMON_AWS_LAMBDA_PYTHON_SLICES)

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'python3 --version && uv --version && aws --version && terraform version && bao --version && vault version && docker --version && java -version'
