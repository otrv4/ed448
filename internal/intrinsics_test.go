package galoisfield

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448InternalSuite) Test_WideMul32(c *C) {
	a := uint32(0x01)
	b := uint32(0x02)
	n := widemul32(a, b)
	c.Assert(n, Equals, uint64(0x02))

	d := uint32(0xff)
	n = widemul32(a, d)
	c.Assert(n, Equals, uint64(0xff))
}

func (s *Ed448InternalSuite) Test_IsWordZero(c *C) {
	a := uint32(0x00)
	ret := isWord32Zero(a)
	c.Assert(ret, Equals, uint32(0x00))

	b := uint32(0x01)
	ret = isWord32Zero(b)
	c.Assert(ret, Equals, uint32(0xffffffff))
}
