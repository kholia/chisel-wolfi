# Examples

These examples build `scratch` images from Wolfi APK packages using Chisel.
Run each build command from the repository root.

All Wolfi examples use the shared release at `examples/wolfi-release/`. Its
`chisel.yaml` pins the Wolfi RSA public key used to verify the signed APK
repository index.

Each example has a Makefile with `build`, `test`, `size`, and `scan` targets:

```sh
make -C examples/wolfi-curl build
make -C examples/wolfi-curl test
make -C examples/wolfi-curl size
make -C examples/wolfi-curl scan
```

The `scan` target uses `grype`; set `GRYPE=/path/to/grype` to override the
scanner command. The top-level `examples/Makefile` can run the same target
across every example, for example `make -C examples test`.

Helper image tags and versioned Wolfi package slice names are controlled in
one place: `examples/wolfi-versions.mk`. Update that file for Go, Node.js,
OpenJDK, Maven, Envoy, kubectl, Python, or Wolfi base image bumps.

The Wolfi examples target small, scanner-friendly images by combining Chisel
slices with a rolling APK package stream. See
[`docs/wolfi-security-model.md`](../docs/wolfi-security-model.md) for the
zero-CVE philosophy and the FIPS/STIG boundary.

For GHCR test image publishing and the `YYYY-MM-DD_N` tag convention, see
[`docs/publishing-wolfi-images.md`](../docs/publishing-wolfi-images.md).

| Example | Purpose |
| --- | --- |
| [wolfi-dockerfile](wolfi-dockerfile) | Generic Wolfi Dockerfile with configurable slices. |
| [wolfi-envoy](wolfi-envoy) | Envoy image with Wolfi's default Envoy config. |
| [wolfi-envoy-gperftools](wolfi-envoy-gperftools) | Experimental Envoy image built from git sources with `--define tcmalloc=gperftools`. Build explicitly; it is not part of the default example set. |
| [wolfi-envoy-upx](wolfi-envoy-upx) | Experimental Envoy image with the Envoy binary compressed by UPX after slicing. Build explicitly; it is not part of the default example set. |
| [wolfi-fluent-bit](wolfi-fluent-bit) | Fluent Bit log processor image with Wolfi's default config and parser files. |
| [wolfi-prometheus](wolfi-prometheus) | Prometheus server and promtool with Wolfi's default config. |
| [wolfi-grafana](wolfi-grafana) | Grafana server with bundled assets and default config. |
| [wolfi-alertmanager](wolfi-alertmanager) | Prometheus Alertmanager and amtool with default config. |
| [wolfi-thanos](wolfi-thanos) | Thanos CLI/server binary. |
| [wolfi-otel-collector](wolfi-otel-collector) | OpenTelemetry Collector Contrib with a minimal OTLP config. |
| [wolfi-kube-state-metrics](wolfi-kube-state-metrics) | kube-state-metrics binary. |
| [wolfi-kube-rbac-proxy](wolfi-kube-rbac-proxy) | kube-rbac-proxy binary. |
| [wolfi-grafana-alloy](wolfi-grafana-alloy) | Grafana Alloy with Wolfi's default config. |
| [wolfi-hello-world](wolfi-hello-world) | Chiseled Wolfi hello-world image. |
| [wolfi-opa](wolfi-opa) | Open Policy Agent binary. |
| [wolfi-opa-envoy](wolfi-opa-envoy) | OPA Envoy-capable binary. |
| [wolfi-gatekeeper](wolfi-gatekeeper) | Gatekeeper manager binary. |
| [wolfi-nginx](wolfi-nginx) | NGINX mainline with default config and static HTML. |
| [wolfi-k8s-tools](wolfi-k8s-tools) | Kubernetes operations toolbox with kubectl, Helm, Argo Rollouts, Flux, jq, yq, Git, and curl. |
| [wolfi-docker-cli](wolfi-docker-cli) | Docker CLI toolbox with Buildx and Compose. |
| [wolfi-valkey](wolfi-valkey) | Valkey server and CLI. |
| [wolfi-memcached](wolfi-memcached) | Memcached server. |
| [wolfi-postgres](wolfi-postgres) | PostgreSQL 18 server tools. |
| [wolfi-kafka](wolfi-kafka) | Apache Kafka with OpenJDK runtime slices. |
| [wolfi-nats-server](wolfi-nats-server) | NATS Server with Wolfi config. |
| [wolfi-coredns](wolfi-coredns) | CoreDNS with a minimal forwarding Corefile. |
| [wolfi-go-build](wolfi-go-build) | Go build image with Git, curl, unzip, CGO toolchain, CA certificates, and grpc health probe. |
| [wolfi-go-runtime](wolfi-go-runtime) | Small Go service runtime base with CA certificates and grpc health probe. |
| [wolfi-node-build](wolfi-node-build) | Node.js build image with npm, yarn, pnpm, Python, Git, curl, and native build tools. |
| [wolfi-node-runtime](wolfi-node-runtime) | Minimal Node.js runtime image with CA certificates. |
| [wolfi-java-build](wolfi-java-build) | OpenJDK build image with Maven, Git, curl, jq, unzip, and native build tools. |
| [wolfi-java-runtime](wolfi-java-runtime) | OpenJDK runtime image with Java trust-store data and CA certificates. |
| [wolfi-aws-cli](wolfi-aws-cli) | AWS CLI v2 image with Python runtime and CA certificates. |
| [wolfi-terraform](wolfi-terraform) | Terraform image with CA certificates. |
| [wolfi-vault](wolfi-vault) | Vault-compatible OpenBao image with CA certificates. |
| [wolfi-curl](wolfi-curl) | Minimal curl image with CA certificates. |
| [wolfi-wget](wolfi-wget) | Minimal wget image with CA certificates. |
| [wolfi-stunnel](wolfi-stunnel) | stunnel binary, sample config, state directory, and CA certificates. |
| [wolfi-unbound](wolfi-unbound) | Unbound daemon and tools with default Wolfi config. |
| [wolfi-network-tools](wolfi-network-tools) | Combined network/debug toolbox with curl, wget, DNS tools, iproute2, ping, jq, file, strace, stunnel, and unbound. |

## Common Wolfi Images

These examples cover common CI, runtime, language, and operations image shapes
using public Wolfi packages. Organization-specific downloads or credentials
belong in downstream layers.

| Example | Purpose |
| --- | --- |
| [wolfi-common-aws-lambda-python312](wolfi-common-aws-lambda-python312) | Python cloud tooling base with AWS CLI, Terraform, OpenBao, Docker CLI, Java, and `uv`. |
| [wolfi-common-chrome-tools](wolfi-common-chrome-tools) | Browser automation companion tools; full Chromium slicing still needs the browser GUI dependency closure. |
| [wolfi-common-unbound-server](wolfi-common-unbound-server) | Unbound DNS server base with Wolfi config and DNS tools. |
| [wolfi-docker-base](wolfi-docker-base) | Docker CI base with Docker CLI/buildx/compose, AWS CLI, gcloud, kubectl, Helm, Python, uv, OpenBao, envsubst, yamllint, jq, yq, and semver. |
| [wolfi-common-docker-toolbox](wolfi-common-docker-toolbox) | Docker CI toolbox with Docker CLI/buildx/compose, cloud tools, Kubernetes tools, and archive utilities. |
| [wolfi-common-go-service-runtime](wolfi-common-go-service-runtime) | Small Go service runtime base with CA certificates and `grpc-health-probe`. |
| [wolfi-common-go-minimal-runtime](wolfi-common-go-minimal-runtime) | Alternate minimal Go runtime base for service images. |
| [wolfi-common-go-base](wolfi-common-go-base) | Go build base with git, curl, tar, gzip, and unzip. |
| [wolfi-common-go-buf](wolfi-common-go-buf) | Go build base plus the `buf` CLI. |
| [wolfi-common-go-protoc](wolfi-common-go-protoc) | Go build base plus `buf`, `protoc`, and protobuf tools. |
| [wolfi-common-go-test](wolfi-common-go-test) | Go test base with CGO build tooling, protobuf tools, and `grpc-health-probe`. |
| [wolfi-common-infra-toolbox](wolfi-common-infra-toolbox) | Infra toolbox with Docker, AWS CLI, gcloud, kubectl, Helm, Terraform, OpenBao, Flux, and argo rollouts. |
| [wolfi-common-infra-toolbox-extended](wolfi-common-infra-toolbox-extended) | Extended infra toolbox base for layering extra organization-specific CLIs. |
| [wolfi-common-java-ci](wolfi-common-java-ci) | Java CI base with OpenJDK, Maven, native build tools, Docker CLI, AWS CLI, and OpenBao. |
| [wolfi-common-java-runtime](wolfi-common-java-runtime) | Small Java runtime base. |
| [wolfi-common-java-build](wolfi-common-java-build) | Java build base with OpenJDK, Maven, git, curl, jq, tar, and unzip. |
| [wolfi-common-node-ui](wolfi-common-node-ui) | Node.js UI build base with npm, yarn, pnpm, Python, and native build tools. |
| [wolfi-common-node-ui-extended](wolfi-common-node-ui-extended) | Alternate Node.js UI build base for frontend applications. |
| [wolfi-common-node-spa](wolfi-common-node-spa) | Single-page app build base with the shared Node.js toolchain. |
| [wolfi-common-node-web](wolfi-common-node-web) | Web app build base with the shared Node.js toolchain. |
| [wolfi-common-node-service](wolfi-common-node-service) | Node.js service build base with `su-exec`. |
| [wolfi-common-node-service-extended](wolfi-common-node-service-extended) | Alternate Node.js service build base with `su-exec`. |
| [wolfi-common-node-runtime](wolfi-common-node-runtime) | Minimal Node.js runtime base. |
| [wolfi-common-node-automation](wolfi-common-node-automation) | Node.js automation base with cloud and Kubernetes CLIs. |
| [wolfi-common-unbound-cache](wolfi-common-unbound-cache) | Unbound cache image with Wolfi config. |
| [wolfi-common-utility-base](wolfi-common-utility-base) | Utility base with wget and `grpc-health-probe`. |
