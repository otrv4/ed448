package ed448

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestDeserialize(c *C) {
	b := serialized{0x1}
	n, ok := deserialize(b)

	c.Assert(n, DeepEquals, bigNumber{1})
	c.Assert(ok, Equals, true)
}

func (s *Ed448Suite) TestConstantTimeGreaterOrEqualP(c *C) {
	//p (little-endian)
	p := bigNumber{
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xfffffffffffffe,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
	}

	greaterThanP := bigNumber{

		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
	}

	lesserThanP := bigNumber{
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

func (s *Ed448Suite) TestSerialize(c *C) {
	dst := [56]byte{}

	one := bigNumber{0x01}
	serialize(dst[:], one)
	c.Assert(dst, DeepEquals, [56]byte{1})

	p := bigNumber{
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
		0xfffffffffffffe,
		0xffffffffffffff,
		0xffffffffffffff,
		0xffffffffffffff,
	}

	serialize(dst[:], p)
	c.Assert(dst, DeepEquals, [56]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})
}
