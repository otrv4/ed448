package ed448

import (
	"bytes"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestRadixBasePointIsOnCurve(c *C) {
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
	curve := newRadixCurve()
	p, err := NewPoint(gx, gy)

	c.Assert(err, IsNil)
	c.Assert(curve.isOnCurve(p), Equals, true)
}

func (s *Ed448Suite) TestRadixMultiplyByBase(c *C) {
	curve := newRadixCurve()
	scalar := [ScalarWords]word_t{}
	scalar[ScalarWords-1] = 1000 //big-endian

	p := curve.multiplyByBase2(scalar)

	c.Assert(curve.isOnCurve(p), Equals, true)
}

func (s *Ed448Suite) TestRadixGenerateKey(c *C) {
	buffer := make([]byte, FieldBytes)
	buffer[55] = 0x10
	r := bytes.NewReader(buffer[:])

	curve := newRadixCurve()
	priv, _, err := curve.generateKey(r)

	expectedPriv := make([]byte, FieldBytes)
	expectedPriv[55] = 0x11

	c.Assert(err, IsNil)
	c.Assert(priv, DeepEquals, expectedPriv)

	//XXX We need to figure out how to serialize this
	//c.Assert(pub, DeepEquals, []byte{0x4,
	//	//gx
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1,
	//	//gy
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	//	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3,
	//})
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

func (s *Ed448Suite) TestAdd(c *C) {
	curve := newBigintsCurve()

	x2, y2 := curve.add(gx, gy, gx, gy)
	x4, y4 := curve.add(gx, gy, x2, y2)

	c.Assert(curve.isOnCurve(x2, y2), Equals, true)
	c.Assert(curve.isOnCurve(x4, y4), Equals, true)
}

func (s *Ed448Suite) TestDouble(c *C) {
	curve := newBigintsCurve()

	xd2, yd2 := curve.double(gx, gy)
	xd4, yd4 := curve.double(xd2, yd2)

	c.Assert(curve.isOnCurve(xd2, yd2), Equals, true)
	c.Assert(curve.isOnCurve(xd4, yd4), Equals, true)
}

func (s *Ed448Suite) TestMultiplication(c *C) {
	curve := newBigintsCurve()

	x2, y2 := curve.multiply(gx, gy, []byte{0x05})

	c.Assert(curve.isOnCurve(x2, y2), Equals, true)
}

func (s *Ed448Suite) TestOperationsAreEquivalent(c *C) {
	curve := newBigintsCurve()

	addX, addY := curve.add(gx, gy, gx, gy)
	doubleX, doubleY := curve.double(gx, gy)
	xBy2, yBy2 := curve.multiply(gx, gy, []byte{2})

	c.Assert(addX, DeepEquals, doubleX)
	c.Assert(addY, DeepEquals, doubleY)
	c.Assert(addX, DeepEquals, xBy2)
	c.Assert(doubleX, DeepEquals, xBy2)
	c.Assert(addY, DeepEquals, yBy2)
	c.Assert(addY, DeepEquals, yBy2)
}

func (s *Ed448Suite) TestBaseMultiplication(c *C) {
	curve := newBigintsCurve()

	x, y := curve.multiplyByBase([]byte{0x05})

	c.Assert(curve.isOnCurve(x, y), Equals, true)
}
*/
