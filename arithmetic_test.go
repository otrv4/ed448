package ed448

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestConstantTimeGreaterOrEqualP(c *C) {
	//p (little-endian)
	p := [8]int64{
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xfffffffffffffe,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
	}

	greaterThanP := [8]int64{
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
	}

	lesserThanP := [8]int64{
		0xfffffffffffffe,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xfffffffffffffe,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
	}

	c.Assert(constantTimeGreaterOrEqualP(p), Equals, true)
	c.Assert(constantTimeGreaterOrEqualP(greaterThanP), Equals, true)
	c.Assert(constantTimeGreaterOrEqualP(lesserThanP), Equals, false)
}
