package ed448

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestNegate(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	n := new(bigNumber).setBytes(bs)
	out := new(bigNumber).neg(n)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	expected := new(bigNumber).setBytes(bs)

	c.Assert(out, DeepEquals, expected)
}

func (s *Ed448Suite) TestZeroMask(c *C) {
	zero := &bigNumber{}
	one := &bigNumber{1}

	c.Assert(zero.zeroMask(), Equals, uint32(0xffffffff))
	c.Assert(one.zeroMask(), Equals, uint32(0))
}
