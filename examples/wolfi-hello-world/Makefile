EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

IMAGE ?= chisel-wolfi:hello-world
DOCKERFILE ?= $(EXAMPLE_DIR)/Dockerfile
CONTEXT ?= $(REPO_ROOT)
BUILDX ?= docker-buildx
GRYPE ?= grype

.PHONY: all build test size scan

all: build test size scan

build:
	$(BUILDX) build --load $(BUILD_ARGS) -f "$(DOCKERFILE)" -t "$(IMAGE)" "$(CONTEXT)"

test: build
	test "$$(docker run --rm "$(IMAGE)")" = "Hello, World!"

size: build
	docker image ls --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}' "$(IMAGE)"
	docker image inspect "$(IMAGE)" --format 'Image size (bytes): {{.Size}}'

scan: build
	@command -v "$(GRYPE)" >/dev/null || { \
		echo "grype not found. Install grype or set GRYPE=/path/to/grype." >&2; \
		exit 127; \
	}
	$(GRYPE) "$(IMAGE)"
