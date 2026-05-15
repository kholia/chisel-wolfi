package apk_test

import (
	"bytes"

	. "gopkg.in/check.v1"

	"github.com/canonical/chisel/internal/apk"
)

func (s *S) TestWriteInstalledDatabase(c *C) {
	var buf bytes.Buffer
	err := apk.WriteInstalledDatabase(&buf, []apk.InstalledPackage{{
		Name:          "busybox",
		Version:       "1.37.0-r58",
		Arch:          "aarch64",
		Checksum:      "Q1abcd",
		Size:          123,
		InstalledSize: 456,
	}, {
		Name:    "noarch-data",
		Version: "1.0-r0",
	}})
	c.Assert(err, IsNil)
	c.Assert(buf.String(), Equals, ""+
		"C:Q1abcd\n"+
		"P:busybox\n"+
		"V:1.37.0-r58\n"+
		"A:aarch64\n"+
		"S:123\n"+
		"I:456\n"+
		"\n"+
		"P:noarch-data\n"+
		"V:1.0-r0\n"+
		"\n")
}
