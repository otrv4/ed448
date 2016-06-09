package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

func (s *Ed448Suite) TestKeysGenerationProducesKeyPair(c *C) {
	c.Skip("Public key is not being set yet.")
	ed448 := NewEd448()
	priv, pub, err := ed448.GenerateKeys()
	c.Assert(err, IsNil)
	c.Assert(priv, NotNil)
	c.Assert(pub, NotNil)
}
