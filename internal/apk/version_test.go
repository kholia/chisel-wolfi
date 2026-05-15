package apk_test

import (
	. "gopkg.in/check.v1"

	"github.com/canonical/chisel/internal/apk"
)

func (s *S) TestCompareVersions(c *C) {
	tests := []struct {
		a, b string
		cmp  int
	}{
		{"1.9-r0", "1.10-r0", -1},
		{"1.0_rc1-r0", "1.0-r0", -1},
		{"1.0-r2", "1.0-r10", -1},
		{"1.0-r0", "1.0_git1-r0", -1},
		{"1.0_p1-r0", "1.0-r0", 1},
	}
	for _, test := range tests {
		c.Logf("%s vs %s", test.a, test.b)
		c.Assert(apk.CompareVersions(test.a, test.b), Equals, test.cmp)
		c.Assert(apk.CompareVersions(test.b, test.a), Equals, -test.cmp)
	}
}
