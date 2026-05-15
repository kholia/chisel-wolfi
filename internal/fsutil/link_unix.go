//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package fsutil

import (
	"os"

	"golang.org/x/sys/unix"
)

func linkNoFollow(oldpath, newpath string) error {
	err := unix.Linkat(unix.AT_FDCWD, oldpath, unix.AT_FDCWD, newpath, unix.AT_SYMLINK_NOFOLLOW)
	if err == unix.EINVAL {
		err = unix.Linkat(unix.AT_FDCWD, oldpath, unix.AT_FDCWD, newpath, 0)
	}
	if err == unix.ENOSYS {
		return os.Link(oldpath, newpath)
	}
	if err != nil {
		return &os.LinkError{Op: "link", Old: oldpath, New: newpath, Err: err}
	}
	return err
}
