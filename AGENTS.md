# Preface

This repository hosts Chisel, a software tool for carving and cutting Debian packages. This `AGENTS.md` is the root index into the dispersed knowledge base; read it to orient yourself, then follow the links below for details.

Read the top-level `.kb/agents.md` file before continuing below.


# Overview

Chisel cuts fine-grained slices out of Debian packages. It takes a selection of slices defined in chisel-releases, resolves their dependencies, downloads the corresponding packages from Ubuntu archives, extracts the selected files into a target root filesystem, runs Starlark mutation scripts, and generates manifests at paths requested by the selected slice definitions. The codebase is split across command entry points, an internals tree holding the carving pipeline, and public packages defining the manifest format.


# Important

- Prefer deterministic, predictable behavior that fails on the side of correctness and preserves existing user content by default.
- Keep Chisel as independent as practical from incidental package archive contents.
- Treat changes to the chisel-releases format and the public manifest packages as compatibility-sensitive.
- Validate Wolfi slice and image changes locally with fresh extraction before pushing upstream; follow `.kb/dependency-security.md`.


# Directory

- `cmd/` - Binary entry points for the `chisel` application.
- `slices/` - Wolfi APK slice definitions, including explicit versioned runtime packages.
- `examples/` - Wolfi release configuration, shared version settings, and container smoke tests.
- `scripts/` - Local Wolfi slice and Docker matrix checks, plus image publishing helpers.
- `.github/workflows/` - Build, test, and vulnerability scan workflows.
- `internal/` - Core internal packages: slice orchestration, setup, extraction, archive, cache, and supporting utilities.
- `public/` - Public API packages for external consumers.
- `snap/` - Snap packaging configuration for Chisel.
- `tests/` - End-to-end integration test suite.
- `docs/` - Documentation assets.
- `go.mod` - Go module definition.
- `go.sum` - Go dependency checksums.
- `spread.yaml` - Spread test configuration.
- `workshop.yaml` - Workshop configuration.
- `.golangci.yaml` - golangci-lint configuration.
- `README.md` - Project readme.
- `SECURITY.md` - Security policy.
- `LICENSE` - Project license.


# Documents

- `.kb/agents.md` - General rules for the knowledge base reading and writing.
- `.kb/dependency-security.md` - Dependency upgrades, toolchain alignment, and vulnerability verification.
- `cmd/AGENTS.md` - CLI entry points.
- `internal/AGENTS.md` - Core internal packages.
- `public/AGENTS.md` - Public API packages.
- `snap/AGENTS.md` - Snap packaging.
- `tests/AGENTS.md` - Integration testing.
