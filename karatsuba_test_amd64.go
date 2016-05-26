package ed448

import (
	. "gopkg.in/check.v1"
)

type Ed448Amd64Suite struct{}

var _ = Suite(&Ed448Amd64Suite{})

func (s *Ed448Amd64Suite) TestWideMul(c *C) {
	//No carry on multiplication
	result := WideMul(Word(0x01), Word(0x02))
	c.Assert(result, Equals, DWord{0, 0x02})

	//With carry on multiplication
	result = WideMul(Word(0xffffffffffffffff), Word(0x02))
	c.Assert(result, Equals, DWord{0x01, 0xfffffffffffffffe})
}

func (s *Ed448Amd64Suite) TestGroupedOperations(c *C) {
	//No carry on addition
	result := multiplyAndAdd(DWord{0, 0x01}, Word(0xffffffffffffffff), Word(0x02))
	c.Assert(result, Equals, DWord{0x01, 0xffffffffffffffff})

	//With carry on addition
	result = multiplyAndAdd(DWord{0, 0x02}, Word(0xffffffffffffffff), Word(0x02))
	c.Assert(result, Equals, DWord{0x02, 0x00})

	//No borrow
	//FIXME: This is acc - a * b
	//result = Msb(Word(0xffffffffffffffff), Word(0x02), DWord{0, 0x01})
	//c.Assert(result, Equals, DWord{0x01, 0xfffffffffffffffd})

	//With borrow
	//FIXME: This is acc - a * b
	//result = Msb(Word(0x8000000000000000), Word(0x02), DWord{0, 0x01})
	//c.Assert(result, Equals, DWord{0x00, 0xffffffffffffffff})
}
