//go:build !aix && !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !solaris

package fsutil

import "os"

func linkNoFollow(oldpath, newpath string) error {
	return os.Link(oldpath, newpath)
}
