EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

IMAGE ?= chisel-wolfi:envoy
DOCKERFILE ?= $(EXAMPLE_DIR)/Dockerfile
CONTEXT ?= $(REPO_ROOT)

include $(EXAMPLE_DIR)/../common.mk

test: build
	docker run --rm "$(IMAGE)" --version
	docker run --rm "$(IMAGE)" --mode validate -c /etc/envoy/envoy.yaml
