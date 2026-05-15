# Central version knobs for Wolfi examples.
#
# Update this file when bumping the helper images used to build examples or
# when moving the example set to newer versioned Wolfi packages.

ifndef WOLFI_VERSIONS_MK_INCLUDED
WOLFI_VERSIONS_MK_INCLUDED := 1

CHISEL_BUILDER_IMAGE ?= golang:1.26-alpine
WOLFI_IMAGE ?= cgr.dev/chainguard/wolfi-base:latest

WOLFI_GO_SLICE ?= go-1.26
WOLFI_NODE_SLICE ?= nodejs-24
WOLFI_NODE_MINIMAL_SLICE ?= nodejs-24-minimal
WOLFI_PYTHON_SLICE ?= python-3.14-base
WOLFI_PYTHON_BIN ?= /usr/bin/python3.14
WOLFI_FALLBACK_PYTHON_BIN ?= /usr/bin/python3.13
WOLFI_PIP_SLICE ?= py3.13-pip
WOLFI_SEMVER_SLICE ?= py3.13-semver-bin
WOLFI_GETTEXT_SLICE ?= gettext-envsubst

WOLFI_JAVA_SLICE ?= openjdk-21
WOLFI_JAVA_HOME ?= /usr/lib/jvm/java-21-openjdk
WOLFI_MAVEN_SLICE ?= maven-3.9
WOLFI_MAVEN_HOME ?= /usr/share/java/maven

WOLFI_ENVOY_SLICE ?= envoy-1.38
WOLFI_ENVOY_CONFIG_SLICE ?= envoy-1.38-config
WOLFI_FLUENT_BIT_SLICE ?= fluent-bit-5.0
WOLFI_PROMETHEUS_SLICE ?= prometheus-3.11
WOLFI_GRAFANA_SLICE ?= grafana-13.0
WOLFI_ALERTMANAGER_SLICE ?= prometheus-alertmanager
WOLFI_THANOS_SLICE ?= thanos
WOLFI_OTEL_COLLECTOR_SLICE ?= opentelemetry-collector-contrib
WOLFI_KUBE_STATE_METRICS_SLICE ?= kube-state-metrics
WOLFI_KUBE_RBAC_PROXY_SLICE ?= kube-rbac-proxy
WOLFI_GRAFANA_ALLOY_SLICE ?= grafana-alloy
WOLFI_OPA_SLICE ?= opa
WOLFI_OPA_ENVOY_SLICE ?= opa-envoy
WOLFI_GATEKEEPER_SLICE ?= gatekeeper-3.21
WOLFI_NGINX_SLICE ?= nginx-mainline
WOLFI_VALKEY_SLICE ?= valkey-9.0
WOLFI_MEMCACHED_SLICE ?= memcached
WOLFI_POSTGRES_SLICE ?= postgresql-18
WOLFI_KAFKA_SLICE ?= kafka-4.2
WOLFI_NATS_SLICE ?= nats-server
WOLFI_COREDNS_SLICE ?= coredns-1.14
WOLFI_KUBECTL_SLICE ?= kubectl-1.35

WOLFI_CLOUDSDK_PYTHON_BIN ?= /usr/bin/python3.13

WOLFI_BUILD_ARG_NAMES := \
	CHISEL_BUILDER_IMAGE \
	WOLFI_IMAGE \
	WOLFI_GO_SLICE \
	WOLFI_NODE_SLICE \
	WOLFI_NODE_MINIMAL_SLICE \
	WOLFI_PYTHON_SLICE \
	WOLFI_PYTHON_BIN \
	WOLFI_FALLBACK_PYTHON_BIN \
	WOLFI_PIP_SLICE \
	WOLFI_SEMVER_SLICE \
	WOLFI_GETTEXT_SLICE \
	WOLFI_JAVA_SLICE \
	WOLFI_JAVA_HOME \
	WOLFI_MAVEN_SLICE \
	WOLFI_MAVEN_HOME \
	WOLFI_ENVOY_SLICE \
	WOLFI_ENVOY_CONFIG_SLICE \
	WOLFI_FLUENT_BIT_SLICE \
	WOLFI_PROMETHEUS_SLICE \
	WOLFI_GRAFANA_SLICE \
	WOLFI_ALERTMANAGER_SLICE \
	WOLFI_THANOS_SLICE \
	WOLFI_OTEL_COLLECTOR_SLICE \
	WOLFI_KUBE_STATE_METRICS_SLICE \
	WOLFI_KUBE_RBAC_PROXY_SLICE \
	WOLFI_GRAFANA_ALLOY_SLICE \
	WOLFI_OPA_SLICE \
	WOLFI_OPA_ENVOY_SLICE \
	WOLFI_GATEKEEPER_SLICE \
	WOLFI_NGINX_SLICE \
	WOLFI_VALKEY_SLICE \
	WOLFI_MEMCACHED_SLICE \
	WOLFI_POSTGRES_SLICE \
	WOLFI_KAFKA_SLICE \
	WOLFI_NATS_SLICE \
	WOLFI_COREDNS_SLICE \
	WOLFI_KUBECTL_SLICE \
	WOLFI_CLOUDSDK_PYTHON_BIN

WOLFI_DOCKER_BUILD_ARGS = $(foreach arg,$(WOLFI_BUILD_ARG_NAMES),--build-arg $(arg)=$($(arg)))

.PHONY: print-build-args print-var

print-build-args:
	@$(foreach arg,$(WOLFI_BUILD_ARG_NAMES),printf '%s\n' '--build-arg' '$(arg)=$($(arg))';)

print-var:
	@test -n "$(VAR)" || { echo "VAR is required" >&2; exit 1; }
	@test -n "$($(VAR))" || { echo "unknown or empty version variable: $(VAR)" >&2; exit 1; }
	@printf '%s\n' "$($(VAR))"

endif
