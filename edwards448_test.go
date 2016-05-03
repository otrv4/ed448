package ed448

import . "gopkg.in/check.v1"

func (s *Ed448Suite) TestBasePointIsOnCurve(c *C) {
	ed448 := newEd448()
	c.Assert(ed448.isOnCurve(ed448.gx, ed448.gy), Equals, true)
}

func (s *Ed448Suite) TestAdd(c *C) {
	ed448 := newEd448()

	x2, y2 := ed448.add(ed448.gx, ed448.gy, ed448.gx, ed448.gy)
	x4, y4 := ed448.add(ed448.gx, ed448.gy, x2, y2)

	c.Assert(ed448.isOnCurve(x2, y2), Equals, true)
	c.Assert(ed448.isOnCurve(x4, y4), Equals, true)
}

func (s *Ed448Suite) TestDouble(c *C) {
	ed448 := newEd448()

	xd2, yd2 := ed448.double(ed448.gx, ed448.gy)
	xd4, yd4 := ed448.double(xd2, yd2)

	c.Assert(ed448.isOnCurve(xd2, yd2), Equals, true)
	c.Assert(ed448.isOnCurve(xd4, yd4), Equals, true)
}

func (s *Ed448Suite) TestMultiplication(c *C) {
	ed448 := newEd448()

	x2, y2 := ed448.multiply(ed448.gx, ed448.gy, []byte{0x05})

	c.Assert(ed448.isOnCurve(x2, y2), Equals, true)
}

func (s *Ed448Suite) TestOperationsAreEquivalent(c *C) {
	ed448 := newEd448()

	addX, addY := ed448.add(ed448.gx, ed448.gy, ed448.gx, ed448.gy)
	doubleX, doubleY := ed448.double(ed448.gx, ed448.gy)
	xBy2, yBy2 := ed448.multiply(ed448.gx, ed448.gy, []byte{2})

	c.Assert(addX, DeepEquals, doubleX)
	c.Assert(addY, DeepEquals, doubleY)
	c.Assert(addX, DeepEquals, xBy2)
	c.Assert(doubleX, DeepEquals, xBy2)
	c.Assert(addY, DeepEquals, yBy2)
	c.Assert(addY, DeepEquals, yBy2)
}

func (s *Ed448Suite) TestBaseMultiplication(c *C) {
	ed448 := newEd448()

	x, y := ed448.multiplyByBase([]byte{0x05})

	c.Assert(ed448.isOnCurve(x, y), Equals, true)
}

func (s *Ed448Suite) BenchmarkAddition(c *C) {
	ed448 := newEd448()
	c.ResetTimer()
	x, y := ed448.gx, ed448.gy
	for i := 0; i < c.N; i++ {
		x, y = ed448.add(x, y, x, y)
	}
}

func (s *Ed448Suite) BenchmarkDoubling(c *C) {
	ed448 := newEd448()
	c.ResetTimer()
	x, y := ed448.gx, ed448.gy
	for i := 0; i < c.N; i++ {
		x, y = ed448.double(x, y)
	}
}

func (s *Ed448Suite) BenchmarkMultiplication(c *C) {
	ed448 := newEd448()
	c.ResetTimer()
	x, y := ed448.gx, ed448.gy
	for i := 0; i < c.N; i++ {
		x, y = ed448.multiply(x, y, []byte{0x03})
	}
}
