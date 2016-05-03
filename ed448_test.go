package ed448

import (
	"bytes"
	"io"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

func (s *Ed448Suite) TestMarshalAndUnmarshal(c *C) {
	ed448 := newEd448()
	x, y := ed448.gx, ed448.gy

	marshaled := marshal(ed448, x, y)

	uncompressedForm := uint8(4)
	c.Assert(marshaled[0], Equals, uncompressedForm)

	ux, uy := unmarshal(ed448, marshaled)

	c.Assert(x, DeepEquals, ux)
	c.Assert(y, DeepEquals, uy)
}

func (s *Ed448Suite) TestKeyGeneration(c *C) {
	c.Skip("This is way to slow to run with big.Int arithmetic.")

	goldilocks := NewGoldilocks()

	priv, pub, err := goldilocks.GenerateKey(getReader())

	c.Assert(err, Equals, nil)
	px, py := unmarshal(ed448, pub)
	c.Assert(priv[1] > 0, Equals, true)
	c.Assert(ed448.isOnCurve(px, py), Equals, true)
}

func getReader() io.Reader {
	len := (newEd448().n.BitLen() + 7)
	zeroes := make([]byte, len)
	return bytes.NewReader(zeroes)
}
