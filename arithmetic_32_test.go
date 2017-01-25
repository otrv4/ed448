package ed448

import . "gopkg.in/check.v1"

var (
	primeSerial = [fieldBytes]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}

	one  = [fieldBytes]byte{1}
	zero = [fieldBytes]byte{}
)

func (s *Ed448Suite) Test_ModQ_WithPrimeOrder(c *C) {
	primeOrderSerial := []byte{
		0xf3, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
		0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x72, 0xc2,
		0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
		0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
	}

	primeOrderModQ := ModQ(primeOrderSerial)
	c.Assert(primeOrderModQ, DeepEquals, zero[:])

	primeOrderPlusOne := []byte{
		0xf4, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
		0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x72, 0xc2,
		0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
		0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
	}

	primeOrderPlusOneModQ := ModQ(primeOrderPlusOne)
	c.Assert(primeOrderPlusOneModQ, DeepEquals, one[:])
}

func (s *Ed448Suite) Test_PointMul(c *C) {
	resultZero := PointMul(zero, testValue)
	resultValueTimes1 := PointMul(one, testValue)

	c.Assert(resultZero, DeepEquals, zero[:])
	c.Assert(resultValueTimes1, DeepEquals, testValue[:])

	val := PointMul(primeSerial, one)
	c.Assert(val, IsNil)

	val = PointMul(one, primeSerial)
	c.Assert(val, IsNil)

	val = PointMul(primeSerial, primeSerial)
	c.Assert(val, IsNil)
}

func (s *Ed448Suite) Test_Add(c *C) {
	valuePlusOne := [fieldBytes]byte{
		0x04, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
		0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x32, 0xc2,
		0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
		0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
		0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0x2f, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
	}

	resultAddZero := PointAddition(zero, testValue)
	resultAddOne := PointAddition(one, testValue)

	c.Assert(resultAddZero, DeepEquals, testValue[:])
	c.Assert(resultAddOne, DeepEquals, valuePlusOne[:])

	val := PointAddition(primeSerial, primeSerial)
	c.Assert(val, IsNil)

	val = PointAddition(primeSerial, zero)
	c.Assert(val, IsNil)

	val = PointAddition(zero, primeSerial)
	c.Assert(val, IsNil)
}

func (s *Ed448Suite) Test_ScalarSub(c *C) {
	twelve := Scalar{0xc}
	thirteen := Scalar{0xd}
	scalarOne := Scalar{0x1}

	result := ScalarSub(thirteen, twelve)

	c.Assert(result, DeepEquals, scalarOne)
}

func (s *Ed448Suite) Test_ScalarMul(c *C) {
	x := Scalar{
		0xffb823a3, 0xc96a3c35,
		0x7f8ed27d, 0x087b8fb9,
		0x1d9ac30a, 0x74d65764,
		0xc0be082e, 0xa8cb0ae8,
		0xa8fa552b, 0x2aae8688,
		0x2c3dc273, 0x47cf8cac,
		0x3b089f07, 0x1e63e807,
	}

	y := Scalar{
		0xd8bedc42, 0x686eb329,
		0xe416b899, 0x17aa6d9b,
		0x1e30b38b, 0x188c6b1a,
		0xd099595b, 0xbc343bcb,
		0x1adaa0e7, 0x24e8d499,
		0x8e59b308, 0x0a92de2d,
		0xcae1cb68, 0x16c5450a,
	}

	expected := Scalar{
		0xa18d010a, 0x1f5b3197,
		0x994c9c2b, 0x6abd26f5,
		0x08a3a0e4, 0x36a14920,
		0x74e9335f, 0x07bcd931,
		0xf2d89c1e, 0xb9036ff6,
		0x203d424b, 0xfccd61b3,
		0x4ca389ed, 0x31e055c1,
	}

	c.Assert(ScalarMul(x, y), DeepEquals, expected)
}