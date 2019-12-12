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

func (s *Ed448InternalSuite) Test_WideMul64(c *C) {
	a := uint64(0x04)
	b := uint64(0x05)
	n := widemul64(a, b)
	c.Assert(n, Equals, uint128{0x00, 0x14})

	d := uint64(0xff)
	e := uint64(0xff)
	n = widemul64(d, e)
	c.Assert(n, Equals, uint128{0x00, 0xfe01})
}

func (s *Ed448InternalSuite) Test_IsWord32Zero(c *C) {
	a := uint32(0x00)
	ret := isWord32Zero(a)
	c.Assert(ret, Equals, true)

	b := uint32(0x01)
	ret = isWord32Zero(b)
	c.Assert(ret, Equals, false)
}

func (s *Ed448InternalSuite) Test_IsWord64Zero(c *C) {
	a := uint64(0x00)
	ret := isWord64Zero(a)
	c.Assert(ret, Equals, true)

	b := uint64(0x01)
	ret = isWord64Zero(b)
	c.Assert(ret, Equals, false)

	d := uint64(0x02)
	ret = isWord64Zero(d)
	c.Assert(ret, Equals, false)

	e := uint64(0x64)
	ret = isWord64Zero(e)
	c.Assert(ret, Equals, false)

	f := uint64(0xff)
	ret = isWord64Zero(f)
	c.Assert(ret, Equals, false)
}
