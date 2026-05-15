EXAMPLES := \
	wolfi-alertmanager \
	wolfi-aws-cli \
	wolfi-coredns \
	wolfi-curl \
	wolfi-docker-cli \
	wolfi-docker-base \
	wolfi-dockerfile \
	wolfi-envoy \
	wolfi-fluent-bit \
	wolfi-gatekeeper \
	wolfi-grafana \
	wolfi-grafana-alloy \
	wolfi-go-build \
	wolfi-go-runtime \
	wolfi-hello-world \
	wolfi-java-build \
	wolfi-java-runtime \
	wolfi-k8s-tools \
	wolfi-kafka \
	wolfi-kube-rbac-proxy \
	wolfi-kube-state-metrics \
	wolfi-memcached \
	wolfi-nats-server \
	wolfi-nginx \
	wolfi-network-tools \
	wolfi-node-build \
	wolfi-node-runtime \
	wolfi-opa \
	wolfi-opa-envoy \
	wolfi-otel-collector \
	wolfi-postgres \
	wolfi-prometheus \
	wolfi-stunnel \
	wolfi-terraform \
	wolfi-thanos \
	wolfi-unbound \
	wolfi-valkey \
	wolfi-vault \
	wolfi-wget \
	wolfi-common-aws-lambda-python312 \
	wolfi-common-chrome-tools \
	wolfi-common-unbound-server \
	wolfi-common-docker-toolbox \
	wolfi-common-go-service-runtime \
	wolfi-common-go-minimal-runtime \
	wolfi-common-go-base \
	wolfi-common-go-buf \
	wolfi-common-go-protoc \
	wolfi-common-go-test \
	wolfi-common-infra-toolbox \
	wolfi-common-infra-toolbox-extended \
	wolfi-common-java-ci \
	wolfi-common-java-runtime \
	wolfi-common-java-build \
	wolfi-common-node-ui \
	wolfi-common-node-ui-extended \
	wolfi-common-node-spa \
	wolfi-common-node-web \
	wolfi-common-node-service \
	wolfi-common-node-service-extended \
	wolfi-common-node-runtime \
	wolfi-common-node-automation \
	wolfi-common-unbound-cache \
	wolfi-common-utility-base

.PHONY: all build test size scan $(EXAMPLES)

all:
	@set -e; \
	for example in $(EXAMPLES); do \
		$(MAKE) -C "$$example" all; \
	done

build test size scan:
	@set -e; \
	for example in $(EXAMPLES); do \
		$(MAKE) -C "$$example" $@; \
	done

$(EXAMPLES):
	$(MAKE) -C "$@"
