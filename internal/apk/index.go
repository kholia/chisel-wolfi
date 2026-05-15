package apk

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Index struct {
	Packages []*Package
}

type Package struct {
	Name      string
	Version   string
	Arch      string
	Checksum  string
	Size      int64
	Installed int64
}

func ParseIndex(reader io.Reader) (*Index, error) {
	scanner := bufio.NewScanner(reader)
	index := &Index{}
	pkg := &Package{}
	addPackage := func() error {
		if *pkg == (Package{}) {
			return nil
		}
		if pkg.Name == "" {
			return fmt.Errorf("package record missing P field")
		}
		if pkg.Version == "" {
			return fmt.Errorf("package %q record missing V field", pkg.Name)
		}
		index.Packages = append(index.Packages, pkg)
		pkg = &Package{}
		return nil
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if err := addPackage(); err != nil {
				return nil, err
			}
			continue
		}
		if len(line) < 2 || line[1] != ':' {
			return nil, fmt.Errorf("invalid APKINDEX line %q", line)
		}
		key, value := line[0], line[2:]
		switch key {
		case 'C':
			pkg.Checksum = value
		case 'P':
			pkg.Name = value
		case 'V':
			pkg.Version = value
		case 'A':
			pkg.Arch = value
		case 'S':
			size, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid package size for %q: %q", pkg.Name, value)
			}
			pkg.Size = size
		case 'I':
			installed, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid installed size for %q: %q", pkg.Name, value)
			}
			pkg.Installed = installed
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := addPackage(); err != nil {
		return nil, err
	}
	return index, nil
}

func (index *Index) SelectPackage(name, arch string) (*Package, error) {
	var selected *Package
	for _, pkg := range index.Packages {
		if pkg.Name != name || !packageMatchesArch(pkg.Arch, arch) {
			continue
		}
		if selected == nil || CompareVersions(selected.Version, pkg.Version) < 0 {
			selected = pkg
		}
	}
	if selected == nil {
		return nil, fmt.Errorf("cannot find package %q in archive", name)
	}
	return selected, nil
}

func packageMatchesArch(pkgArch, targetArch string) bool {
	return pkgArch == "" || pkgArch == "noarch" || strings.EqualFold(pkgArch, targetArch)
}
