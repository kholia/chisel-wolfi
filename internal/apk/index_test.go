package apk_test

import (
	"bytes"
	"strings"

	. "gopkg.in/check.v1"

	"github.com/canonical/chisel/internal/apk"
	"github.com/canonical/chisel/internal/testutil"
)

func (s *S) TestIndexReader(c *C) {
	data := testutil.MustMakeAPKIndex("P:mypkg\nV:1.0-r0\nA:x86_64\n\n", true)
	reader, err := apk.IndexReader(bytes.NewReader(data))
	c.Assert(err, IsNil)
	defer reader.Close()

	index, err := apk.ParseIndex(reader)
	c.Assert(err, IsNil)
	c.Assert(index.Packages, DeepEquals, []*apk.Package{{
		Name:    "mypkg",
		Version: "1.0-r0",
		Arch:    "x86_64",
	}})
}

func (s *S) TestParseIndexAndSelect(c *C) {
	index, err := apk.ParseIndex(strings.NewReader(`
C:Q1aaaa
P:mypkg
V:1.9-r0
A:x86_64
S:1
I:2

C:Q1bbbb
P:mypkg
V:1.10-r0
A:x86_64
S:3
I:4

C:Q1cccc
P:mypkg
V:2.0-r0
A:aarch64
S:5
I:6

`[1:]))
	c.Assert(err, IsNil)

	pkg, err := index.SelectPackage("mypkg", "x86_64")
	c.Assert(err, IsNil)
	c.Assert(pkg, DeepEquals, &apk.Package{
		Name:      "mypkg",
		Version:   "1.10-r0",
		Arch:      "x86_64",
		Checksum:  "Q1bbbb",
		Size:      3,
		Installed: 4,
	})
}
