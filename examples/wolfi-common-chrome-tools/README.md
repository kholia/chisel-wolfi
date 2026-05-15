# wolfi-common-chrome-tools

This slices the public tools used around the Chrome automation image: git, wget, unzip, jq, and OpenBao. The upstream image is based on `chromedp/headless-shell`; a fully chiseled Chromium runtime still needs the large GUI/browser dependency closure to be sliced.

```sh
make -C examples/wolfi-common-chrome-tools all
```
