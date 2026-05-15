package apk_test

import (
	"crypto/rand"
	"crypto/rsa"

	. "gopkg.in/check.v1"

	"github.com/canonical/chisel/internal/apk"
	"github.com/canonical/chisel/internal/testutil"
)

func (s *S) TestVerifySignature(c *C) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)

	pkg := testutil.MustMakeSignedAPK([]testutil.TarEntry{
		testutil.Reg(0644, "file", "content"),
	}, key, "test.rsa.pub")

	err = apk.VerifySignature(pkg, []*rsa.PublicKey{&key.PublicKey})
	c.Assert(err, IsNil)
}

func (s *S) TestVerifySignatureInvalid(c *C) {
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)
	verifyKey, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)

	pkg := testutil.MustMakeSignedAPK([]testutil.TarEntry{
		testutil.Reg(0644, "file", "content"),
	}, signingKey, "test.rsa.pub")

	err = apk.VerifySignature(pkg, []*rsa.PublicKey{&verifyKey.PublicKey})
	c.Assert(err, ErrorMatches, `cannot verify APK RSA signature ".SIGN.RSA256.test.rsa.pub"`)
}

func (s *S) TestVerifySignatureMissing(c *C) {
	pkg := testutil.MustMakeAPK([]testutil.TarEntry{
		testutil.Reg(0644, "file", "content"),
	}, false)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)

	err = apk.VerifySignature(pkg, []*rsa.PublicKey{&key.PublicKey})
	c.Assert(err, Equals, apk.ErrNoSignature)
}
