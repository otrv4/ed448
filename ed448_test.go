package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

func (s *Ed448Suite) TestMarshalAndUnmarshal(c *C) {
	ed448 := Ed448()
	x, y := ed448.Params().Gx, ed448.Params().Gy

	marshaled := Marshal(ed448, x, y)

	uncompressedForm := uint8(4)
	c.Assert(marshaled[0], Equals, uncompressedForm)

	ux, uy := Unmarshal(ed448, marshaled)

	c.Assert(x, DeepEquals, ux)
	c.Assert(y, DeepEquals, uy)
}
