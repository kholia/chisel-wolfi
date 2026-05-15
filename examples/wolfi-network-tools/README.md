# Wolfi Network Tools Image

Build a `scratch` network/debug toolbox from Wolfi APK slices.

Included tools:

- `curl`
- `wget`
- `openssl`
- `jq`
- `dig`, `host`, `nslookup`
- `ip`, `ss`, `ping`, `tracepath`
- `file`
- `strace`
- `stunnel`
- `unbound` and Unbound tools

```sh
make -C examples/wolfi-network-tools build
docker run --rm chisel-wolfi:network-tools -c 'curl --version; wget --version | busybox head -n 1; dig -v; jq --version'
```
