package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

func (s *Ed448Suite) TestGenerateKeysProducesKeyPair(c *C) {
	ed448 := NewCurve()
	priv, pub, ok := ed448.GenerateKeys()
	c.Assert(ok, Equals, true)
	c.Assert(priv, NotNil)
	c.Assert(pub, NotNil)
}

func (s *Ed448Suite) TestSignAndVerify(c *C) {
	ed448 := NewCurve()
	priv, pub, _ := ed448.GenerateKeys()
	message := []byte("sign here.")

	signature, ok := ed448.Sign(priv, message)

	c.Assert(ok, Equals, true)
	c.Assert(signature, NotNil)

	valid := ed448.Verify(signature, message, pub)

	c.Assert(valid, Equals, true)
}

func (s *Ed448Suite) TestComputeSecret(c *C) {
	ed448 := NewCurve()
	privA, pubA, _ := ed448.GenerateKeys()
	privB, pubB, _ := ed448.GenerateKeys()
	secretA := ed448.ComputeSecret(privA, pubB)
	secretB := ed448.ComputeSecret(privB, pubA)
	c.Assert(secretA, DeepEquals, secretB)
}
