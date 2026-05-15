make -C examples/wolfi-hello-world test BUILD_ARGS="--build-arg GOPROXY=https://proxy.golang.org,direct --build-arg GOSUMDB=sum.golang.org"
