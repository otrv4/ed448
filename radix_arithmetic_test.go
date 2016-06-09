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
	scalar := [ScalarWords]word_t{}
	scalar[ScalarWords-1] = 1000 //big-endian

	p := curve.multiplyByBase(scalar)

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
