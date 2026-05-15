package apk

import (
	"fmt"
	"io"
)

type InstalledPackage struct {
	Name          string
	Version       string
	Arch          string
	Checksum      string
	Size          int64
	InstalledSize int64
}

func WriteInstalledDatabase(w io.Writer, packages []InstalledPackage) error {
	for _, pkg := range packages {
		if pkg.Name == "" {
			return fmt.Errorf("installed package missing name")
		}
		if pkg.Version == "" {
			return fmt.Errorf("installed package %q missing version", pkg.Name)
		}
		if pkg.Checksum != "" {
			if _, err := fmt.Fprintf(w, "C:%s\n", pkg.Checksum); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "P:%s\nV:%s\n", pkg.Name, pkg.Version); err != nil {
			return err
		}
		if pkg.Arch != "" {
			if _, err := fmt.Fprintf(w, "A:%s\n", pkg.Arch); err != nil {
				return err
			}
		}
		if pkg.Size > 0 {
			if _, err := fmt.Fprintf(w, "S:%d\n", pkg.Size); err != nil {
				return err
			}
		}
		if pkg.InstalledSize > 0 {
			if _, err := fmt.Fprintf(w, "I:%d\n", pkg.InstalledSize); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}
