package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

func (s *Ed448Suite) TestGenerateKeysProducesKeyPair(c *C) {
	c.Skip("Public key is not being set yet.")
	ed448 := NewEd448()
	priv, pub, err := ed448.GenerateKeys()
	c.Assert(err, IsNil)
	c.Assert(priv, NotNil)
	c.Assert(pub, NotNil)
}

func (s *Ed448Suite) TestSignAndVerify(c *C) {
	c.Skip("Public key is not being set yet.")
	ed448 := NewEd448()
	priv, pub, _ := ed448.GenerateKeys()
	message := []byte("sign here.")

	signature, err := ed448.Sign(priv, message)

	c.Assert(err, IsNil)
	c.Assert(signature, NotNil)

	valid := ed448.Verify(signature, message, pub)

	c.Assert(valid, Equals, true)
}
