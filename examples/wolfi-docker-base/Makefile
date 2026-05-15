EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:docker-base
COMMON_SLICES := $(COMMON_DOCKER_CI_BASE_SLICES)

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'docker --version && docker buildx version && docker-compose version && make --version >/dev/null && curl --version >/dev/null && wget --version >/dev/null && unzip -v >/dev/null && git --version && aws --version && bash --version >/dev/null && rsync --version >/dev/null && envsubst --version && yamllint --version && python3 --version && pip3 --version && uv --version && uvx --version && bao --version && vault version && jq --version && yq --version && semver --version && helm version && kubectl version --client=true && gcloud --version && gke-gcloud-auth-plugin --version && kubectl-argo-rollouts version'
