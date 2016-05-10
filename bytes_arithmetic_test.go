package ed448

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestBytesBasePointIsOnCurve(c *C) {
	c.Skip("This is not yet implemented.")
	curve := newBytesCurve()
	c.Assert(curve.isOnCurve(gx, gy), Equals, true)
}

func (s *Ed448Suite) TestSum(c *C) {
	c.Assert(
		sum([]byte{0x57, 0x0, 0x0, 0x0}, []byte{0x83, 0x0, 0x0, 0x0}),
		DeepEquals,
		[]byte{0xda, 0x0, 0x0, 0x0})

	c.Assert(
		sum([]byte{0x75, 0xbc, 0xd1, 0x5}, []byte{0x1, 0x0, 0x0, 0x0}),
		DeepEquals,
		[]byte{0x76, 0xbc, 0xd1, 0x5})
}
