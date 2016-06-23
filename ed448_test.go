package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

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
