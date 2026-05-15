EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

IMAGE ?= chisel-wolfi:java-build
DOCKERFILE ?= $(EXAMPLE_DIR)/Dockerfile
CONTEXT ?= $(REPO_ROOT)

include $(EXAMPLE_DIR)/../common.mk

test: build
	docker run --rm "$(IMAGE)" -c 'java --version && javac --version && jar --version && mvn --version && git --version && gcc --version >/dev/null'
