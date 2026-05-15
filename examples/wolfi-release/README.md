# Shared Wolfi Release

This directory contains the single `chisel.yaml` used by every Wolfi Docker
example. Keep the Wolfi archive URL and pinned RSA public key here so all
examples verify the signed APK repository index consistently.

The `slices` symlink points at the repository-level `slices/` directory, so this
directory can be used directly as a Chisel release root.
