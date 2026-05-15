# Wolfi Security Model

Chisel's Wolfi support is designed around two ideas that work well together:
cut only the files a workload needs, and rebuild from a rolling APK repository
that moves quickly to current fixed package versions.

## Zero-CVE Philosophy

Zero-CVE images are enabled by a rolling release strategy. When a fix lands
upstream, a rolling distribution can publish the fixed package version directly.
A rebuilt image then carries the new version, and vulnerability scanners can
usually see that the affected package is fixed.

That is different from a backporting strategy. Stable distributions often keep
the same upstream version and carry selected security patches on top. This is a
valid maintenance model, but scanner output can remain noisy because the visible
package version still appears affected, or because issues are triaged as not
fixed, not affected, or won't fix for that release stream. Those findings tend
to accumulate over time in broad base images.

A point-in-time scan of `ubuntu:noble` shows the pattern:

```text
$ grype ubuntu:noble
Parsed image sha256:3c220491de2e7cd3946456f3e8b2a2d2e0f52383235f5cab93ead70cca9ee9ed
Cataloged contents
  Packages:      92 packages
  Executables:   722 executables
  File metadata: 2,263 locations
  File digests:  2,263 files
Scanned for vulnerabilities: 80 vulnerability matches
  by severity: 0 critical, 0 high, 60 medium, 18 low, 2 negligible
  by status:   3 fixed, 77 not-fixed, 0 ignored
```

The goal of the Wolfi examples is to make that output smaller in two ways:
start from packages that are updated in a rolling model, then use Chisel slices
to remove files and packages that the workload does not need.

## What Chisel Can And Cannot Promise

Chisel reduces the filesystem and records the package metadata for what it cut.
It does not make an unfixed package fixed, and it does not make a non-validated
cryptographic module FIPS validated.

The security posture comes from the full chain:

- signed APK repository indexes
- package signature verification when packages are signed
- current packages from the selected repository
- small, workload-specific slice selections
- retained `/etc/os-release` and APK database metadata for scanner visibility
- CI rebuilds and vulnerability scans

## FIPS Support

FIPS support is a property of validated cryptographic modules, runtime
configuration, and the package stream used to build the image. Chisel can cut
FIPS-capable packages from a configured APK repository, but it cannot create FIPS
validation by slicing public non-FIPS packages.

For FIPS images, use a release configuration that points at a repository or base
image stream containing the validated FIPS provider and hardened OpenSSL
configuration, then keep the FIPS package metadata in the final image.

## STIG Compliance

STIG compliance is broader than package selection. It includes image contents,
configuration, users, permissions, runtime defaults, and operational controls.
Chisel helps by reducing the image contents and preserving package metadata, but
STIG evidence still needs to come from a hardened image profile and a compliance
scan of the final artifact.

The practical launch posture is:

- use rolling Wolfi packages for the default zero-CVE path
- keep FIPS and STIG images as explicit, separately validated variants
- publish scanner output for both default and compliance-focused images
