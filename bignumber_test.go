package ed448

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestZero(c *C) {
	notZero := mustDeserialize(serialized{0x01})
	c.Assert(notZero.zero(), Equals, false)

	zero := mustDeserialize(serialized{0x00})
	c.Assert(zero.zero(), Equals, true)
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

func (s *Ed448Suite) TestSubRadix(c *C) {
	x := mustDeserialize(serialized{0x57})
	y := mustDeserialize(serialized{0x83})
	z := mustDeserialize(serialized{0xda})
	c.Assert(new(bigNumber).sub(z, y).strongReduce(), DeepEquals, x)

	x = mustDeserialize(serialized{0xff, 0xff, 0xff, 0xf0})
	y = mustDeserialize(serialized{0x01})
	z = mustDeserialize(serialized{0x00, 0x00, 0x00, 0xf1})
	c.Assert(new(bigNumber).sub(z, y).strongReduce(), DeepEquals, x)
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

func (s *Ed448Suite) TestNegate(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	n := new(bigNumber).setBytes(bs)
	out := new(bigNumber).neg(n)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	expected := new(bigNumber).setBytes(bs)

	c.Assert(out, DeepEquals, expected)
}

func (s *Ed448Suite) TestConditionalNegateNumber(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	n := new(bigNumber).setBytes(bs)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	negated := new(bigNumber).setBytes(bs)

	c.Assert(n.copy().conditionalNegate(lmask), DeepEquals, negated)
	c.Assert(n.copy().conditionalNegate(0), DeepEquals, n)
}

func (s *Ed448Suite) TestConditionalSelect(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	x := new(bigNumber).setBytes(bs)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	y := new(bigNumber).setBytes(bs)

	c.Assert(constantTimeSelect(x, y, lmask), DeepEquals, x)
	c.Assert(constantTimeSelect(x, y, 0), DeepEquals, y)

}

func (s *Ed448Suite) TestConditionalSwap(c *C) {
	bs, _ := hex.DecodeString("e6f5b8ae49cef779e577dc29824eff453f1c4106030088115ea49b4ee84a7b7cdfe06e0d622fc55c7c559ab1f6c3ea3257c07979809026de")
	x := new(bigNumber).setBytes(bs)

	bs, _ = hex.DecodeString("190a4751b63108861a8823d67db100bac0e3bef9fcff77eea15b64b017b58483201f91f29dd03aa383aa654e093c15cda83f86867f6fd921")
	y := new(bigNumber).setBytes(bs)

	a := x.copy()
	b := y.copy()
	a.conditionalSwap(b, lmask)

	c.Assert(a, DeepEquals, y)
	c.Assert(b, DeepEquals, x)

	a.conditionalSwap(b, 0)
	c.Assert(a, DeepEquals, y)
	c.Assert(b, DeepEquals, x)
}

func (s *Ed448Suite) Test_DecafConditionalNegateNumber(c *C) {
	n := &bigNumber{
		0x08db85c2, 0x0fd2361e,
		0x0ce2105d, 0x06a17729,
		0x0a137aa5, 0x0e3ca84d,
		0x0985ee61, 0x05a26d64,
		0x0734c5f3, 0x0da853af,
		0x01d955b7, 0x03160ecd,
		0x0a59046d, 0x0c32cf71,
		0x98dce72d, 0x00007fff}

	expected := &bigNumber{
		0x07247a3d, 0x002dc9e1,
		0x031defa2, 0x095e88d6,
		0x05ec855a, 0x01c357b2,
		0x067a119e, 0x0a5d929b,
		0x08cb3a0b, 0x0257ac50,
		0x0e26aa48, 0x0ce9f132,
		0x05a6fb92, 0x03cd308e,
		0x072318d2, 0x0fff8007,
	}

	n.decafCondNegate(lmask)

	c.Assert(n, DeepEquals, expected)

	n1 := &bigNumber{}

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

	n1.decafCondNegate(lmask)

	// 0 mod p = n1
	c.Assert(n1, DeepEquals, p)
}
