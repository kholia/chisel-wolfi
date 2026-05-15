package apk

import (
	"fmt"
	"runtime"
)

type archPair struct {
	goArch  string
	apkArch string
}

var knownArchs = []archPair{
	{"386", "x86"},
	{"amd64", "x86_64"},
	{"arm", "armv7"},
	{"arm64", "aarch64"},
	{"ppc64le", "ppc64le"},
	{"riscv64", "riscv64"},
	{"s390x", "s390x"},
}

var platformGoArch = runtime.GOARCH

func InferArch() (string, error) {
	for _, arch := range knownArchs {
		if arch.goArch == platformGoArch {
			return arch.apkArch, nil
		}
	}
	return "", fmt.Errorf("cannot infer package architecture from current platform architecture: %s", platformGoArch)
}

func ValidateArch(apkArch string) error {
	for _, arch := range knownArchs {
		if arch.apkArch == apkArch {
			return nil
		}
	}
	return fmt.Errorf("invalid package architecture: %s", apkArch)
}
