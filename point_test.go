package ed448

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestPoint(c *C) {
	//Base point
	gx := serialized{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	}
	gy := serialized{0x13}

	basePoint, err := NewPoint(gx, gy)
	c.Assert(err, IsNil)

	curve := newRadixCurve()
	c.Assert(curve.isOnCurve(basePoint), Equals, true)

	p := basePoint.Double()
	c.Assert(curve.isOnCurve(p), Equals, true)

	q := basePoint.Add(basePoint)
	c.Assert(curve.isOnCurve(q), Equals, true)
}

func (s *Ed448Suite) TestTwistedPoint(c *C) {
	//Base point
	gx := serialized{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	}
	gy := serialized{0x13}

	basePoint, err := NewTwistedPoint(gx, gy)

	c.Assert(err, IsNil)

	c.Assert(basePoint.OnCurve(), Equals, true)

	q := basePoint.Add(basePoint)
	c.Assert(q.OnCurve(), Equals, true)
}

func (s *Ed448Suite) TestMixedAddition(c *C) {
	pa, _ := hex.DecodeString("4b8a632c1feab72769cd96e7aaa577861871b3613945c802b89377e8b85331ecc0ffb1cb20169bfc9c27274d38b0d01e87a1d5d851770bc8")
	pb, _ := hex.DecodeString("81a45f02f41053f8d7d2a1f176a340529b33b7ee4d3fa84de384b750b35a54c315bf36c41d023ade226449916e668396589ea2145da09b95")
	pc, _ := hex.DecodeString("5f5a2b06a2dbf7136f8dc979fd54d631ca7de50397250a196d3be2a721ab7cbaa92c545d9b15b5319e11b64bc031666049d8637e13838b3b")

	n := &twNiels{
		a: new(bigNumber).setBytes(pa),
		b: new(bigNumber).setBytes(pb),
		c: new(bigNumber).setBytes(pc),
	}

	px, _ := hex.DecodeString("e45b0207cf5036bcb75a775cb4eb3e8312a8d2b6c9c309dc6a589d2824427848e1ccc7ddac1a53d028375ff6b329d9f0998ed9bb4c81b4e9")
	py, _ := hex.DecodeString("e7c9798862329c3db188697a564706eade026ad6c773ca35069fd53f5d36c0b9db9fbda22386702aae4694ea2dfbe5e97458dd9040b2b97f")
	pz, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	pt, _ := hex.DecodeString("e45b0207cf5036bcb75a775cb4eb3e8312a8d2b6c9c309dc6a589d2824427848e1ccc7ddac1a53d028375ff6b329d9f0998ed9bb4c81b4e9")
	pu, _ := hex.DecodeString("e7c9798862329c3db188697a564706eade026ad6c773ca35069fd53f5d36c0b9db9fbda22386702aae4694ea2dfbe5e97458dd9040b2b97f")

	tx := &twExtensible{
		new(bigNumber).setBytes(px),
		new(bigNumber).setBytes(py),
		new(bigNumber).setBytes(pz),
		new(bigNumber).setBytes(pt),
		new(bigNumber).setBytes(pu),
	}

	px, _ = hex.DecodeString("ac7cd8b31f6f031e0ee8d606c12dafd6503305fd398f55399e3a677543dd7c8239995872432b428aba728f99d7bef231cb32a125ea9e7a69")
	py, _ = hex.DecodeString("4392419b29d484975fbb27272745408e81ae8ab9d1818c8181c39637437b3e456ec78335e6637c20df95a708da038a42213f29079ff1a457")
	pz, _ = hex.DecodeString("32aebf8745d98e59b0b4ad64b1f5e3ef67413c6ad82993ab24f1d8102f6da3c4206729e3392807daf13acc980447545ab91ce1ec6a7f5728")
	pt, _ = hex.DecodeString("b5c30299cb7bd622abdfc0a47e95b00cfc037561016114052a7bcf8862f7e4c656885e0b1800cc4f7d9046592912e2e0ee12662edc2588e6")
	pu, _ = hex.DecodeString("cf4f0c8cbed27e95dbbef285d2f6d68ba9819fdf44c1e09e0e4fc8be9f5b94aeac0c10886de2fa80f688a45a082861813dcf5afc7cd9f820")

	expected := &twExtensible{
		new(bigNumber).setBytes(px),
		new(bigNumber).setBytes(py),
		new(bigNumber).setBytes(pz),
		new(bigNumber).setBytes(pt),
		new(bigNumber).setBytes(pu),
	}

	ret := tx.addTwNiels(n)
	c.Assert(ret.equals(expected), Equals, true)
}
