EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:common-go-protoc
COMMON_SLICES := $(COMMON_GO_PROTOC_SLICES)

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -ec 'go version; \
		protoc --version; \
		printf "%s\n" "syntax = \"proto3\"; message Smoke { string value = 1; }" > /tmp/smoke.proto; \
		protoc -I/tmp --upb_out=/tmp --upb_minitable_out=/tmp --upbdefs_out=/tmp /tmp/smoke.proto; \
		test -s /tmp/smoke.upb.c; \
		test -s /tmp/smoke.upb_minitable.c; \
		test -s /tmp/smoke.upbdefs.c'
