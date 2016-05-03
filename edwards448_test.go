package ed448

import . "gopkg.in/check.v1"

func (s *Ed448Suite) TestBasePointIsOnCurve(c *C) {
	ed448 := newEd448()
	c.Assert(ed448.IsOnCurve(ed448.Gx, ed448.Gy), Equals, true)
}

func (s *Ed448Suite) TestAdd(c *C) {
	ed448 := newEd448()

	x2, y2 := ed448.Add(ed448.Gx, ed448.Gy, ed448.Gx, ed448.Gy)
	x4, y4 := ed448.Add(ed448.Gx, ed448.Gy, x2, y2)

	c.Assert(ed448.IsOnCurve(x2, y2), Equals, true)
	c.Assert(ed448.IsOnCurve(x4, y4), Equals, true)
}

func (s *Ed448Suite) TestDouble(c *C) {
	ed448 := newEd448()

	xd2, yd2 := ed448.Double(ed448.Gx, ed448.Gy)
	xd4, yd4 := ed448.Double(xd2, yd2)

	c.Assert(ed448.IsOnCurve(xd2, yd2), Equals, true)
	c.Assert(ed448.IsOnCurve(xd4, yd4), Equals, true)
}

func (s *Ed448Suite) TestMultiplication(c *C) {
	ed448 := newEd448()

	x2, y2 := ed448.Multiply(ed448.Gx, ed448.Gy, []byte{0x05})

	c.Assert(ed448.IsOnCurve(x2, y2), Equals, true)
}

func (s *Ed448Suite) TestOperationsAreEquivalent(c *C) {
	ed448 := newEd448()

	addX, addY := ed448.Add(ed448.Gx, ed448.Gy, ed448.Gx, ed448.Gy)
	doubleX, doubleY := ed448.Double(ed448.Gx, ed448.Gy)
	xBy2, yBy2 := ed448.Multiply(ed448.Gx, ed448.Gy, []byte{2})

	c.Assert(addX, DeepEquals, doubleX)
	c.Assert(addY, DeepEquals, doubleY)
	c.Assert(addX, DeepEquals, xBy2)
	c.Assert(doubleX, DeepEquals, xBy2)
	c.Assert(addY, DeepEquals, yBy2)
	c.Assert(addY, DeepEquals, yBy2)
}

func (s *Ed448Suite) TestBaseMultiplication(c *C) {
	ed448 := newEd448()

	x, y := ed448.MultiplyByBase([]byte{0x05})

	c.Assert(ed448.IsOnCurve(x, y), Equals, true)
}

func (s *Ed448Suite) BenchmarkAddition(c *C) {
	ed448 := newEd448()
	c.ResetTimer()
	x, y := ed448.Gx, ed448.Gy
	for i := 0; i < c.N; i++ {
		x, y = ed448.Add(x, y, x, y)
	}
}

func (s *Ed448Suite) BenchmarkDoubling(c *C) {
	ed448 := newEd448()
	c.ResetTimer()
	x, y := ed448.Gx, ed448.Gy
	for i := 0; i < c.N; i++ {
		x, y = ed448.Double(x, y)
	}
}

func (s *Ed448Suite) BenchmarkMultiplication(c *C) {
	ed448 := newEd448()
	c.ResetTimer()
	x, y := ed448.Gx, ed448.Gy
	for i := 0; i < c.N; i++ {
		x, y = ed448.Multiply(x, y, []byte{0x03})
	}
}
