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

func (s *Ed448Suite) TestWordMult(c *C) {
	//No carry on multiplication
	result := WideMul(Word(0x01), Word(0x02))
	c.Assert(result, Equals, DWord{0, 0x02})

	//With carry on multiplication
	result = WideMul(Word(0xffffffffffffffff), Word(0x02))
	c.Assert(result, Equals, DWord{0x01, 0xfffffffffffffffe})

	//No carry on addition
	result = Mac(Word(0xffffffffffffffff), Word(0x02), DWord{0, 0x01})
	c.Assert(result, Equals, DWord{0x01, 0xffffffffffffffff})

	//With carry on addition
	result = Mac(Word(0xffffffffffffffff), Word(0x02), DWord{0, 0x02})
	c.Assert(result, Equals, DWord{0x02, 0x00})

	//No borrow
	result, _ = Msb(Word(0xffffffffffffffff), Word(0x02), DWord{0, 0x01})
	c.Assert(result, Equals, DWord{0x01, 0xfffffffffffffffd})

	//With borrow
	result, _ = Msb(Word(0x8000000000000000), Word(0x02), DWord{0, 0x01})
	c.Assert(result, Equals, DWord{0x00, 0xffffffffffffffff})
}
