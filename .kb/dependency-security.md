# Preface

This document describes dependency upgrades and vulnerability verification for Chisel.

Read the top-level `.kb/agents.md` file before continuing below.

# Overview

`.github/workflows/security.yaml` checks source and test dependencies with pinned `govulncheck`, builds a Linux binary, checks its vulnerable symbols, and scans repository dependencies, CI actions, and the binary with Grype. Findings rated low or higher fail the Grype checks; reachable vulnerabilities fail the Go checks.

# Important

- Use the maintained `github.com/ProtonMail/go-crypto/openpgp` packages for archive signatures. Do not reintroduce `golang.org/x/crypto/openpgp`; GO-2026-5932 affects every version of those unmaintained packages.
- YAML v3 is maintained at `go.yaml.in/yaml/v3`. Preserve slice parsing and rendering compatibility when upgrading it.
- Keep `go.mod`, the Snap and Workshop Go channels, and `CHISEL_BUILDER_IMAGE` in `examples/wolfi-versions.mk` aligned with a supported, patched toolchain. Wolfi's `WOLFI_GO_SLICE` selects a separate packaged toolchain and requires matching slice definitions before changing tracks.
- Pin GitHub Actions to verified release commits, including actions in the performance workflows.
- `anchore/scan-action` treats its `path` input as a directory (`dir:`). Build release binaries into a dedicated directory and pass that directory to the action; pass the individual binary to `govulncheck -mode=binary`.
- Wolfi package aliases in APKINDEX are not resolved by Chisel. Slice filenames, package fields, and essential references must use actual package names. Keep the glibc runtime, loader, and development slices on the same track; current examples use 2.44. OpenSSL 3 compatibility slices also pull OpenSSL 4 libraries because current programs need both ABIs.
- Upgrade protobuf, protoc, their versioned shared-library packages, and Abseil together. The protobuf example smoke test generates code with all three UPB plugins; plugin `--version` is unsupported and must not be used with `|| true` to hide failures.
- Use `scripts/check-wolfi-slices.py` to rerun the CI extraction matrix and filesystem assertions locally, then run the relevant Docker example smoke tests to catch missing shared libraries and symbol versions. The script keeps roots and logs and accepts architecture, case, concurrency, cache, and output directory options.
- Before pushing Wolfi slice or example build changes, run `scripts/check-wolfi-examples.py` for fresh Docker builds and smoke tests. It pulls base images and disables the builder-stage cache to detect changes in rolling APK indexes. Use `BUILD_FLAGS` for these flags when invoking Make directly; overriding `BUILD_ARGS` can erase the common images' required slice selection.
- Several graphics packages (`libpng`, `libx11`, `lcms2`, `libxext`, `libxi`, `libxrender`, and `libxtst`) are now APK metapackages. Keep their existing slice names as dependency wrappers and select files from the split runtime and data packages that own them. Refresh indexes and check the complete Java dependency chain when these packages change.
- Select ICU libraries and their matching data package per consumer ABI. Keep parallel ICU slice versions when rolling packages require different ABIs; a global replacement can fix PostgreSQL while breaking Node.js. Extraction checks alone do not detect a mismatched shared-library ABI, so run executable smoke tests too.
- Grype's module-level GO-2026-5932 match can persist because ProtonMail's implementation depends on `golang.org/x/crypto`. Check `go list -deps -test ./...` for the affected import paths and run `govulncheck -test ./...` before classifying it. No ignore rule suppresses this advisory.

For local validation, run race tests, the configured `golangci-lint`, and `go run .github/scripts/external-packages-license-check.go`. Scan both `go.mod` and the repository, then build and scan the Linux amd64, arm64, and riscv64 release binaries. A `go.mod` scan alone omits embedded standard-library and CI-action vulnerabilities. Run `go test ./internal/archive/ --real-archive` with network access when changing signature verification.
