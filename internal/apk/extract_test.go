package apk_test

import (
	"archive/tar"
	"bytes"
	"io"

	. "gopkg.in/check.v1"

	"github.com/canonical/chisel/internal/apk"
	"github.com/canonical/chisel/internal/testutil"
)

func (s *S) TestDataReader(c *C) {
	for _, signed := range []bool{false, true} {
		c.Logf("signed: %v", signed)
		pkg := testutil.MustMakeAPK([]testutil.TarEntry{
			testutil.Dir(0755, "dir/"),
			testutil.Reg(0644, "dir/file", "apk data"),
		}, signed)
		dataReader, err := apk.DataReader(bytes.NewReader(pkg))
		c.Assert(err, IsNil)
		defer dataReader.Close()

		tarReader := tar.NewReader(dataReader)
		header, err := tarReader.Next()
		c.Assert(err, IsNil)
		c.Assert(header.Name, Equals, "dir/")
		header, err = tarReader.Next()
		c.Assert(err, IsNil)
		c.Assert(header.Name, Equals, "dir/file")
		data, err := io.ReadAll(tarReader)
		c.Assert(err, IsNil)
		c.Assert(string(data), Equals, "apk data")
	}
}

func (s *S) TestExtract(c *C) {
	pkg := testutil.MustMakeAPK([]testutil.TarEntry{
		testutil.Dir(0755, "dir/"),
		testutil.Reg(0644, "dir/file", "apk data"),
	}, true)
	targetDir := c.MkDir()
	err := apk.Extract(bytes.NewReader(pkg), &apk.ExtractOptions{
		Package:   "test-package",
		TargetDir: targetDir,
		Extract: map[string][]apk.ExtractInfo{
			"/dir/file": {{
				Path: "/dir/file",
			}},
		},
	})
	c.Assert(err, IsNil)
	c.Assert(testutil.TreeDump(targetDir), DeepEquals, map[string]string{
		"/dir/":     "dir 0755",
		"/dir/file": "file 0644 62721614",
	})
}

func (s *S) TestControlChecksum(c *C) {
	pkg := testutil.MustMakeAPK([]testutil.TarEntry{
		testutil.Reg(0644, "file", "content"),
	}, true)
	checksum, err := apk.ControlChecksum(bytes.NewReader(pkg))
	c.Assert(err, IsNil)
	c.Assert(checksum, Matches, `Q1[A-Za-z0-9+/]+=*`)
}
