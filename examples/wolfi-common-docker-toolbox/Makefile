EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:common-docker-toolbox
COMMON_SLICES := $(COMMON_DOCKER_BASE_SLICES)

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'docker --version && docker buildx version && docker-compose version && make --version >/dev/null && curl --version >/dev/null && wget --version >/dev/null && unzip -v >/dev/null && git --version && aws --version && terraform version && bao --version && vault version && jq --version && yq --version && yamllint --version && helm version && kubectl version --client=true && kubectl-argo-rollouts version && flux --version && buf --version && uv --version'
