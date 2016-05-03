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
	x, y := ed448.Gx, ed448.Gy

	marshaled := Marshal(ed448, x, y)

	uncompressedForm := uint8(4)
	c.Assert(marshaled[0], Equals, uncompressedForm)

	ux, uy := Unmarshal(ed448, marshaled)

	c.Assert(x, DeepEquals, ux)
	c.Assert(y, DeepEquals, uy)
}

func (s *Ed448Suite) TestKeyGeneration(c *C) {
	c.Skip("This is way to slow to run with big.Int arithmetic.")

	ed448 := newEd448()
	random := getReader()

	priv, pub, err := GenerateKey(ed448, random)

	c.Assert(err, Equals, nil)

	px, py := Unmarshal(ed448, pub)
	c.Assert(priv[1] > 0, Equals, true)
	c.Assert(ed448.IsOnCurve(px, py), Equals, true)
}

func getReader() io.Reader {
	len := (newEd448().N.BitLen() + 7)
	zeroes := make([]byte, len)
	return bytes.NewReader(zeroes)
}
