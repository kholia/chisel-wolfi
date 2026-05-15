# Wolfi AWS CLI Image

Build a `scratch` image with Wolfi's AWS CLI v2 package.

```sh
make -C examples/wolfi-aws-cli build
docker run --rm chisel-wolfi:aws-cli --version
docker run --rm chisel-wolfi:aws-cli sts get-caller-identity
```
