package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

var testValue = [fieldBytes]byte{
	0x03, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
	0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x32, 0xc2,
	0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
	0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
	0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0x2f, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
}

func (s *Ed448Suite) TestGenerateKeysProducesKeyPair(c *C) {
	curve := NewCurve()
	priv, pub, ok := curve.GenerateKeys()
	c.Assert(ok, Equals, true)
	c.Assert(priv, NotNil)
	c.Assert(pub, NotNil)
}

func (s *Ed448Suite) TestSignAndVerify(c *C) {
	curve := NewCurve()
	priv, pub, ok := curve.GenerateKeys()
	c.Assert(ok, Equals, true)

	message := []byte("sign here.")

	signature, ok := curve.Sign(priv, message)

	c.Assert(ok, Equals, true)
	c.Assert(signature, NotNil)

	valid := curve.Verify(signature, message, pub)

	c.Assert(valid, Equals, true)
}

func (s *Ed448Suite) TestComputeSecret(c *C) {
	curve := NewCurve()
	privA, pubA, _ := curve.GenerateKeys()
	privB, pubB, _ := curve.GenerateKeys()
	secretA := curve.ComputeSecret(privA, pubB)
	secretB := curve.ComputeSecret(privB, pubA)
	c.Assert(secretA, DeepEquals, secretB)
}
