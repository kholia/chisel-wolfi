# wolfi-docker-base

Builds a Docker CI base image from public Wolfi packages. It includes Docker CLI, Buildx, Compose, AWS CLI, gcloud, GKE auth plugin, kubectl, Helm, argo rollouts, Python, pip, uv, OpenBao with `vault` compatibility, envsubst, yamllint, jq, yq, semver, git, curl, wget, rsync, make, zip, and unzip.

```sh
make -C examples/wolfi-docker-base all
```

This intentionally leaves private downloads and non-Wolfi-packaged tools to downstream layers.
