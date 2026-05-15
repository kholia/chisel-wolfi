include $(dir $(abspath $(lastword $(MAKEFILE_LIST))))wolfi-versions.mk

COMMON_MINIMAL_BASE_SLICES := \
	wolfi-baselayout_base-library \
	ca-certificates-bundle_certs \
	busybox_shell

COMMON_BASE_UTILITY_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	bash_bins \
	coreutils_bins \
	findutils_bins \
	grep_bins \
	sed_bins \
	gnutar_bins \
	gzip_bins \
	unzip_bins \
	zip_bins \
	wget_bins \
	curl_bins \
	git_bins \
	jq_bins \
	file_bins

COMMON_SECRET_TOOL_SLICES := \
	openbao_bins \
	openbao-compat_vault

COMMON_DOCKER_BASE_SLICES := \
	$(COMMON_BASE_UTILITY_SLICES) \
	make_bins \
	rsync_bins \
	aws-cli-2_bins \
	terraform_bins \
	$(COMMON_SECRET_TOOL_SLICES) \
	docker-cli_bins \
	docker-cli-buildx_plugin \
	docker-compose_bins \
	helm_bins \
	yq_bins \
	yamllint_bins \
	$(WOLFI_KUBECTL_SLICE)_bins \
	kubectl-argo-rollouts_bins \
	flux_bins \
	buf_bins \
	uv_bins

COMMON_DOCKER_CI_BASE_SLICES := \
	$(COMMON_BASE_UTILITY_SLICES) \
	make_bins \
	rsync_bins \
	$(WOLFI_GETTEXT_SLICE)_bins \
	aws-cli-2_bins \
	$(COMMON_SECRET_TOOL_SLICES) \
	docker-cli_bins \
	docker-cli-buildx_plugin \
	docker-compose_bins \
	$(WOLFI_PIP_SLICE)_bins \
	uv_bins \
	$(WOLFI_SEMVER_SLICE)_bins \
	helm_bins \
	yq_bins \
	$(WOLFI_KUBECTL_SLICE)_bins \
	kubectl-argo-rollouts_bins \
	google-cloud-sdk-core_bins \
	gke-gcloud-auth-plugin_bins \
	yamllint_bins

COMMON_INFRA_BASE_SLICES := \
	$(COMMON_DOCKER_BASE_SLICES) \
	grpc-health-probe_bins \
	google-cloud-sdk-core_bins \
	gke-gcloud-auth-plugin_bins

COMMON_AWS_LAMBDA_PYTHON_SLICES := \
	$(COMMON_BASE_UTILITY_SLICES) \
	$(WOLFI_PYTHON_SLICE)_runtime \
	uv_bins \
	aws-cli-2_bins \
	terraform_bins \
	$(COMMON_SECRET_TOOL_SLICES) \
	docker-cli_bins \
	$(WOLFI_JAVA_SLICE)_runtime

COMMON_GO_BASE_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	bash_bins \
	coreutils_bins \
	findutils_bins \
	gnutar_bins \
	gzip_bins \
	unzip_bins \
	git_bins \
	curl_bins \
	$(WOLFI_GO_SLICE)_bins

COMMON_GO_BUILD_SLICES := \
	$(COMMON_GO_BASE_SLICES) \
	make_bins \
	build-base_dev

COMMON_GO_BUF_SLICES := \
	$(COMMON_GO_BUILD_SLICES) \
	buf_bins

COMMON_GO_PROTOC_SLICES := \
	$(COMMON_GO_BUF_SLICES) \
	protoc_bins \
	protobuf_bins

COMMON_GO_TEST_SLICES := \
	$(COMMON_GO_PROTOC_SLICES) \
	grpc-health-probe_bins

COMMON_GO_RUNTIME_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	grpc-health-probe_bins

COMMON_JAVA_BUILD_SLICES := \
	$(COMMON_BASE_UTILITY_SLICES) \
	make_bins \
	build-base_dev \
	$(WOLFI_JAVA_SLICE)_jdk \
	$(WOLFI_MAVEN_SLICE)_bins \
	docker-cli_bins \
	aws-cli-2_bins \
	$(COMMON_SECRET_TOOL_SLICES)

COMMON_JAVA_RUNTIME_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	$(WOLFI_JAVA_SLICE)_runtime

COMMON_NODE_BUILD_SLICES := \
	$(COMMON_BASE_UTILITY_SLICES) \
	$(WOLFI_NODE_SLICE)_bins \
	npm_bins \
	yarn_bins \
	pnpm_bins \
	$(WOLFI_PYTHON_SLICE)_runtime \
	build-base_dev

COMMON_NODE_UI_SLICES := \
	$(COMMON_NODE_BUILD_SLICES)

COMMON_NODE_PM2_SLICES := \
	$(COMMON_NODE_UI_SLICES) \
	su-exec_bins

COMMON_NODE_CLAUDE_SLICES := \
	$(COMMON_NODE_BUILD_SLICES) \
	aws-cli-2_bins \
	docker-cli_bins \
	helm_bins \
	$(WOLFI_KUBECTL_SLICE)_bins \
	terraform_bins \
	$(COMMON_SECRET_TOOL_SLICES) \
	yq_bins

COMMON_NODE_RUNTIME_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	$(WOLFI_NODE_MINIMAL_SLICE)_runtime

COMMON_UNBOUND_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	bind-tools_bins \
	unbound_bins \
	unbound-config_config

COMMON_UBUNTU_BASE_SLICES := \
	$(COMMON_MINIMAL_BASE_SLICES) \
	wget_bins \
	grpc-health-probe_bins

COMMON_CHROME_SLICES := \
	$(COMMON_BASE_UTILITY_SLICES) \
	$(COMMON_SECRET_TOOL_SLICES)

CI_BASE_SHELL_SLICES := \
	wolfi-baselayout_base-library \
	busybox_shell

CI_NETWORK_TOOLS_SLICES := \
	wolfi-baselayout_base \
	curl_bins \
	wget_bins \
	stunnel_bins \
	stunnel_config \
	stunnel_state \
	unbound_bins

CI_CLOUD_TOOLS_SLICES := \
	wolfi-baselayout_base \
	ca-certificates-bundle_certs \
	aws-cli-2_bins \
	terraform_bins \
	$(COMMON_SECRET_TOOL_SLICES)

CI_ENVOY_SLICES := \
	wolfi-baselayout_base-library \
	ca-certificates-bundle_certs \
	$(WOLFI_ENVOY_SLICE)_bins \
	$(WOLFI_ENVOY_CONFIG_SLICE)_config

CI_LANGUAGE_RUNTIME_SLICES := \
	wolfi-baselayout_base-library \
	ca-certificates-bundle_certs \
	grpc-health-probe_bins \
	$(WOLFI_NODE_MINIMAL_SLICE)_runtime \
	$(WOLFI_JAVA_SLICE)_runtime

CI_LANGUAGE_BUILD_SLICES := \
	wolfi-baselayout_base-library \
	ca-certificates-bundle_certs \
	busybox_shell \
	bash_bins \
	coreutils_bins \
	findutils_bins \
	grep_bins \
	sed_bins \
	gnutar_bins \
	gzip_bins \
	unzip_bins \
	git_bins \
	curl_bins \
	$(WOLFI_GO_SLICE)_bins \
	build-base_dev \
	$(WOLFI_NODE_SLICE)_bins \
	npm_bins \
	yarn_bins \
	pnpm_bins \
	$(WOLFI_JAVA_SLICE)_jdk \
	$(WOLFI_MAVEN_SLICE)_bins

CI_COMMON_CI_TOOLS_SLICES := \
	$(COMMON_BASE_UTILITY_SLICES) \
	make_bins \
	rsync_bins \
	aws-cli-2_bins \
	terraform_bins \
	$(COMMON_SECRET_TOOL_SLICES) \
	docker-cli_bins \
	docker-cli-buildx_plugin \
	docker-compose_bins \
	helm_bins \
	yq_bins \
	$(WOLFI_KUBECTL_SLICE)_bins \
	kubectl-argo-rollouts_bins \
	flux_bins \
	buf_bins \
	uv_bins

CI_DOCKER_CI_BASE_SLICES := \
	$(COMMON_DOCKER_CI_BASE_SLICES)

.PHONY: print-slices

print-slices:
	@test -n "$(SLICE_SET)" || { echo "SLICE_SET is required" >&2; exit 1; }
	@test -n "$(strip $($(SLICE_SET)))" || { echo "unknown or empty slice set: $(SLICE_SET)" >&2; exit 1; }
	@printf '%s\n' "$(strip $($(SLICE_SET)))"
