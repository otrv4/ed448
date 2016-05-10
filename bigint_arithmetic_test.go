package ed448

import (
	"math/big"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestBigintsBasePointIsOnCurve(c *C) {
	curve := newBigintsCurve()
	c.Assert(curve.isOnCurve(gx, gy), Equals, true)
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

func (s *Ed448Suite) BenchmarkAddition(c *C) {
	curve := newBigintsCurve()
	c.ResetTimer()
	x, y := gx, gy
	for i := 0; i < c.N; i++ {
		rx, ry := curve.add(x, y, x, y)
		x, y = rx.(*big.Int), ry.(*big.Int)
	}
}

func (s *Ed448Suite) BenchmarkDoubling(c *C) {
	curve := newBigintsCurve()
	c.ResetTimer()
	x, y := gx, gy
	for i := 0; i < c.N; i++ {
		rx, ry := curve.double(x, y)
		x, y = rx.(*big.Int), ry.(*big.Int)
	}
}

func (s *Ed448Suite) BenchmarkMultiplication(c *C) {
	curve := newBigintsCurve()
	c.ResetTimer()
	x, y := gx, gy
	for i := 0; i < c.N; i++ {
		rx, ry := curve.multiply(x, y, []byte{0x03})
		x, y = rx.(*big.Int), ry.(*big.Int)
	}
}
