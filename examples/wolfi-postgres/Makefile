EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

IMAGE ?= chisel-wolfi:postgres
DOCKERFILE ?= $(EXAMPLE_DIR)/Dockerfile
CONTEXT ?= $(REPO_ROOT)

include $(EXAMPLE_DIR)/../common.mk

test: build
	docker run --rm "$(IMAGE)" --version
	docker run --rm --entrypoint /usr/bin/initdb "$(IMAGE)" --version
	docker run --rm --entrypoint /usr/bin/pg_ctl "$(IMAGE)" --version
