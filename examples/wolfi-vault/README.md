# Wolfi Vault-Compatible OpenBao Image

Build a `scratch` image named for Vault discoverability, backed by Wolfi's
OpenBao package and `openbao-compat` Vault-compatible CLI symlink.

```sh
make -C examples/wolfi-vault build
docker run --rm chisel-wolfi:vault --version
docker run --rm --entrypoint /usr/bin/bao chisel-wolfi:vault --version
docker run --rm chisel-wolfi:vault server -dev -dev-listen-address=0.0.0.0:8200
```
