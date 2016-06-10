package ed448

import (
	"bytes"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestRadixBasePointIsOnCurve(c *C) {
	curve := newRadixCurve()
	p := curve.BasePoint()

	c.Assert(curve.isOnCurve(p), Equals, true)
}

func (s *Ed448Suite) TestRadixMultiplyByBase(c *C) {
	curve := newRadixCurve()
	scalar := [scalarWords]word_t{}
	scalar[scalarWords-1] = 1000 //little-endian

	p := curve.multiplyByBase(scalar)

	c.Assert(curve.isOnCurve(p), Equals, true)
}

//func (s *Ed448Suite) TestRadixMultiplyByBaseAgain(c *C) {
//	curve := newRadixCurve()
//	scalar := [ScalarWords]word_t{}
//	scalar[0] = 2
//
//	//p = basePoint * 2
//	p := curve.multiplyByBase(scalar)
//
//	c.Assert(curve.isOnCurve(p), Equals, true)
//
//	gx := mustDeserialize(serialized{
//		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
//		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
//		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
//		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
//		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
//		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
//		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
//		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
//	})
//
//	gy := mustDeserialize(serialized{0x13})
//
//	q := newExtensible(gx, gy).twist().double()
//	c.Assert(q.OnCurve(), Equals, true)
//	c.Assert(q.equals(p), Equals, true)
//}

func (s *Ed448Suite) TestRadixGenerateKey(c *C) {
	buffer := make([]byte, symKeyBytes)
	buffer[0] = 0x10
	r := bytes.NewReader(buffer[:])

	curve := newRadixCurve()
	privKey, err := curve.generateKey(r)

	expectedSymKey := make([]byte, symKeyBytes)
	expectedSymKey[0] = 0x10

	expectedPriv := []byte{
		0x06, 0x01, 0x3f, 0x3e, 0xb3, 0x3f, 0x9e, 0x10,
		0xde, 0xde, 0x34, 0x23, 0x6a, 0x9a, 0x75, 0x44,
		0x69, 0x41, 0x18, 0x4f, 0x79, 0xb7, 0x52, 0x50,
		0x03, 0xa0, 0x7d, 0xe2, 0x89, 0xee, 0x15, 0x8a,
		0xaf, 0x44, 0xf3, 0x39, 0x78, 0x2c, 0xa6, 0x9b,
		0xbe, 0x5b, 0xb4, 0x1d, 0x25, 0x6a, 0x83, 0x32,
		0x7c, 0xd0, 0xc0, 0x3d, 0xa5, 0x26, 0xf8, 0x37,
	}

	expectedPublic := []byte{
		0x4d, 0xdb, 0xad, 0x93, 0xb8, 0x95, 0x29, 0x61,
		0x67, 0xfc, 0xf4, 0xbd, 0x27, 0x94, 0xb9, 0x0f,
		0x06, 0x09, 0x05, 0xef, 0x8f, 0x32, 0x63, 0x2c,
		0xa6, 0xce, 0x45, 0xfb, 0x1c, 0x83, 0xc5, 0xe7,
		0x0f, 0xf9, 0xf4, 0x43, 0x2a, 0x0c, 0xaf, 0x82,
		0x7a, 0xf5, 0x19, 0xe9, 0x5e, 0x40, 0x17, 0x48,
		0x44, 0xb9, 0xf8, 0x11, 0x88, 0x9a, 0xc3, 0xa5,
	}

	c.Assert(err, IsNil)
	c.Assert(privKey.symKey(), DeepEquals, expectedSymKey)
	c.Assert(privKey.secretKey(), DeepEquals, expectedPriv)
	c.Assert(privKey.publicKey(), DeepEquals, expectedPublic)
}

func (s *Ed448Suite) TestMultiplication(c *C) {
	curve := newRadixCurve()

	scalar := []byte{0x02}
	p1 := curve.multiplyRaw(scalar, curve.BasePoint())
	p2 := curve.multiply(scalar, curve.BasePoint())

	c.Assert(curve.isOnCurve(p1), Equals, true)
	c.Assert(curve.isOnCurve(p2), Equals, true)
	// c.Assert(p2.Marshal(), DeepEquals, p1.Marshal())
}

/*
func (s *Ed448Suite) TestComputeSecret(c *C) {
	curve := newRadixCurve()
	privA, pubA, err := curve.generateKey(rand.Reader)
	fmt.Printf("privA %v\npubA %v\n", privA, pubA)
	c.Assert(err, Equals, nil)
	privB, pubB, err := curve.generateKey(rand.Reader)
	fmt.Printf("privB %v\npubB %v\n", privB, pubB)
	c.Assert(err, Equals, nil)
	out := curve.computeSecret(privA, pubB)
	expected := curve.computeSecret(privB, pubA)
	c.Assert(out, DeepEquals, expected)
}
*/

func (s *Ed448Suite) TestAdd(c *C) {
	curve := newRadixCurve()

	p2 := curve.add(curve.BasePoint(), curve.BasePoint())
	p4 := curve.add(p2, p2)

	c.Assert(curve.isOnCurve(p2), Equals, true)
	c.Assert(curve.isOnCurve(p4), Equals, true)
}

func (s *Ed448Suite) TestDouble(c *C) {
	curve := newRadixCurve()

	p2 := curve.double(curve.BasePoint())
	p4 := curve.double(p2)

	c.Assert(curve.isOnCurve(p2), Equals, true)
	c.Assert(curve.isOnCurve(p4), Equals, true)
}

func (s *Ed448Suite) TestOperationsAreEquivalent(c *C) {
	curve := newRadixCurve()

	//XXX something wrong here
	// addp2 := curve.add(curve.BasePoint(), curve.BasePoint())
	doublep2 := curve.double(curve.BasePoint())
	mulp2 := curve.multiply([]byte{0x02}, curve.BasePoint())

	// c.Assert(addp2, DeepEquals, doublep2)
	c.Assert(doublep2, DeepEquals, mulp2)
}
