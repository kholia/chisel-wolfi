include $(dir $(abspath $(lastword $(MAKEFILE_LIST))))wolfi-versions.mk

GRYPE ?= grype

.PHONY: all build test size scan

all: build test size scan

build:
	docker build $(WOLFI_DOCKER_BUILD_ARGS) $(BUILD_ARGS) -f "$(DOCKERFILE)" -t "$(IMAGE)" "$(CONTEXT)"

size: build
	docker image ls --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}' "$(IMAGE)"
	docker image inspect "$(IMAGE)" --format 'Image size (bytes): {{.Size}}'

scan: build
	@command -v "$(GRYPE)" >/dev/null || { \
		echo "grype not found. Install grype or set GRYPE=/path/to/grype." >&2; \
		exit 127; \
	}
	$(GRYPE) "$(IMAGE)"
