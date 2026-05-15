# Wolfi Terraform Image

Build a `scratch` image with Wolfi's Terraform package and CA certificates.

```sh
make -C examples/wolfi-terraform build
docker run --rm chisel-wolfi:terraform version
docker run --rm chisel-wolfi:terraform -help
```
