package ed448

import (
	"math/big"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) BenchmarkBigintsAddition(c *C) {
	curve := newBigintsCurve()
	c.ResetTimer()
	x, y := gx, gy
	for i := 0; i < c.N; i++ {
		rx, ry := curve.add(x, y, x, y)
		x, y = rx.(*big.Int), ry.(*big.Int)
	}
}

func (s *Ed448Suite) BenchmarkRadixAddition(c *C) {
	curve := newRadixCurve()
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
