package ed448

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestStrongReduce(c *C) {
	p, _ := deserialize(serialized{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})

	//p = p mod p = 0
	p.strongReduce()

	c.Assert(p, DeepEquals, &bigNumber{})

	n := mustDeserialize(serialized{
		0xf5, 0x81, 0x74, 0xd5, 0x7a, 0x33, 0x72,
		0x36, 0x3c, 0x0d, 0x9f, 0xcf, 0xaa, 0x3d,
		0xc1, 0x8b, 0x1e, 0xff, 0x7e, 0x89, 0xbf,
		0x76, 0x78, 0x63, 0x65, 0x80, 0xd1, 0x7d,
		0xd8, 0x4a, 0x87, 0x3b, 0x14, 0xb9, 0xc0,
		0xe1, 0x68, 0x0b, 0xbd, 0xc8, 0x76, 0x47,
		0xf3, 0xc3, 0x82, 0x90, 0x2d, 0x2f, 0x58,
		0xd2, 0x75, 0x4b, 0x39, 0xbc, 0xa8, 0x74,
	})

	n.strongReduce()

	c.Assert(n, DeepEquals, mustDeserialize(serialized{
		0xf5, 0x81, 0x74, 0xd5, 0x7a, 0x33, 0x72,
		0x36, 0x3c, 0x0d, 0x9f, 0xcf, 0xaa, 0x3d,
		0xc1, 0x8b, 0x1e, 0xff, 0x7e, 0x89, 0xbf,
		0x76, 0x78, 0x63, 0x65, 0x80, 0xd1, 0x7d,
		0xd8, 0x4a, 0x87, 0x3b, 0x14, 0xb9, 0xc0,
		0xe1, 0x68, 0x0b, 0xbd, 0xc8, 0x76, 0x47,
		0xf3, 0xc3, 0x82, 0x90, 0x2d, 0x2f, 0x58,
		0xd2, 0x75, 0x4b, 0x39, 0xbc, 0xa8, 0x74,
	}))
}

func (s *Ed448Suite) TestSumRadix(c *C) {
	x := mustDeserialize(serialized{0x57})
	y := mustDeserialize(serialized{0x83})
	z := mustDeserialize(serialized{0xda})
	c.Assert(new(bigNumber).add(x, y), DeepEquals, z)

	x = mustDeserialize(serialized{0xff, 0xff, 0xff, 0xf0})
	y = mustDeserialize(serialized{0x01})
	z = mustDeserialize(serialized{0x00, 0x00, 0x00, 0xf1})
	c.Assert(new(bigNumber).add(x, y), DeepEquals, z)
}

//XXX This is broken in 64-bits, but everything else works
//func (s *Ed448Suite) TestSubRadix(c *C) {
//	x := mustDeserialize(serialized{0x57})
//	y := mustDeserialize(serialized{0x83})
//	z := mustDeserialize(serialized{0xda})
//	c.Assert(subRadix(z, y).strongReduce(), DeepEquals, x)
//
//	x = mustDeserialize(serialized{0xff, 0xff, 0xff, 0xf0})
//	y = mustDeserialize(serialized{0x01})
//	z = mustDeserialize(serialized{0x00, 0x00, 0x00, 0xf1})
//	c.Assert(subRadix(z, y).strongReduce(), DeepEquals, x)
//}

func (s *Ed448Suite) TestEquals(c *C) {
	p, _ := deserialize(serialized{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})

	c.Assert(p.equals(p), Equals, true)

	x := mustDeserialize(serialized{0x01, 0x01})
	y := mustDeserialize(serialized{0x01, 0x02})
	c.Assert(x.equals(y), Equals, false)
}

func (s *Ed448Suite) TestZero(c *C) {
	notZero := mustDeserialize(serialized{0x01})
	c.Assert(notZero.zero(), Equals, false)

	zero := mustDeserialize(serialized{0x00})
	c.Assert(zero.zero(), Equals, true)
}

func (s *Ed448Suite) TestNegate(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	n := new(bigNumber).setBytes(bs)
	out := new(bigNumber).neg(n)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	expected := new(bigNumber).setBytes(bs)

	c.Assert(out, DeepEquals, expected)
}

func (s *Ed448Suite) TestConditionalSelect(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	x := new(bigNumber).setBytes(bs)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	y := new(bigNumber).setBytes(bs)

	c.Assert(constantTimeSelect(x, y, 0xffffffff), DeepEquals, x)
	c.Assert(constantTimeSelect(x, y, 0), DeepEquals, y)
}

func (s *Ed448Suite) TestConditionalSwap(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	x := new(bigNumber).setBytes(bs)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	y := new(bigNumber).setBytes(bs)

	a := x.copy()
	b := y.copy()
	a.conditionalSwap(b, 0xffffffff)

	c.Assert(a, DeepEquals, y)
	c.Assert(b, DeepEquals, x)

	a.conditionalSwap(b, 0)
	c.Assert(a, DeepEquals, y)
	c.Assert(b, DeepEquals, x)
}

func (s *Ed448Suite) TestConditionalNegateNumber(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	n := new(bigNumber).setBytes(bs)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	negated := new(bigNumber).setBytes(bs)

	c.Assert(n.copy().conditionalNegate(0xffffffff), DeepEquals, negated)
	c.Assert(n.copy().conditionalNegate(0), DeepEquals, n)
}

func (s *Ed448Suite) TestZeroMask(c *C) {
	zero := &bigNumber{}
	one := &bigNumber{1}

	c.Assert(zero.zeroMask(), Equals, uint32(0xffffffff))
	c.Assert(one.zeroMask(), Equals, uint32(0))
}

func (s *Ed448Suite) TestSquareN(c *C) {
	gx := mustDeserialize(serialized{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	})

	exp := gx.copy()
	for i := 0; i < 5; i++ {
		exp = new(bigNumber).square(exp)
	}

	n := new(bigNumber).squareN(gx, 5)

	c.Assert(n.equals(exp), Equals, true)

	exp = gx.copy()
	for i := 0; i < 6; i++ {
		exp = new(bigNumber).square(exp)
	}

	n = n.squareN(gx, 6)

	c.Assert(n.equals(exp), Equals, true)
}

func (s *Ed448Suite) TestISR(c *C) {
	gx := mustDeserialize(serialized{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	})

	gx.isr(gx)

	bs, _ := hex.DecodeString("04027d13a34bbe052fdf4247b02a4a3406268203a09076e56dee9dc2b699c4abc66f2832a677dfd0bf7e70ee72f01db170839717d1c64f02")
	exp := new(bigNumber).setBytes(bs)

	c.Assert(gx.equals(exp), Equals, true)
}
