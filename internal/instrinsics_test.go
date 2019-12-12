package galoisfield

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448InternalSuite) Test_IsWordZero(c *C) {
	a := uint32(0x0)
	ret := isWord32Zero(a)
	c.Assert(ret, Equals, uint32(0x0))

	b := uint32(0x01)
	ret = isWord32Zero(b)
	c.Assert(ret, Equals, uint32(0xffffffff))
}
