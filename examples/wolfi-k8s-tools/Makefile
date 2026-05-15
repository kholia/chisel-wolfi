EXAMPLE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT := $(abspath $(EXAMPLE_DIR)/../..)

include $(EXAMPLE_DIR)/../wolfi-common-slices.mk

IMAGE ?= chisel-wolfi:k8s-tools
COMMON_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	bash_bins \
	curl_bins \
	git_bins \
	jq_bins \
	yq_bins \
	helm_bins \
	$(WOLFI_KUBECTL_SLICE)_bins \
	kubectl-argo-rollouts_bins \
	flux_bins

include $(EXAMPLE_DIR)/../wolfi-common-image.mk

test: build
	docker run --rm "$(IMAGE)" -c 'kubectl version --client=true && helm version && kubectl-argo-rollouts version && flux --version && jq --version && yq --version && git --version && curl --version >/dev/null'
