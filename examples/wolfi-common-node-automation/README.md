# wolfi-common-node-automation

Builds a Node.js automation base with npm, Python, native build tools, AWS CLI, Docker CLI, kubectl, Helm, Terraform, OpenBao, yq, git, curl, and archive utilities. The upstream `@anthropic-ai/claude-code` npm install belongs in a downstream layer so this example stays package-reproducible from Wolfi slices.

```sh
make -C examples/wolfi-common-node-automation all
```
