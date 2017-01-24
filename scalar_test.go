package ed448

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) Test_ScalarAddition(c *C) {
	s1 := [scalarWords]uint32{
		0x529eec33, 0x721cf5b5,
		0xc8e9c2ab, 0x7a4cf635,
		0x44a725bf, 0xeec492d9,
		0x0cd77058, 0x00000002,
	}
	s2 := [scalarWords]uint32{0x00000001}
	expected := [scalarWords]uint32{
		0x529eec34, 0x721cf5b5,
		0xc8e9c2ab, 0x7a4cf635,
		0x44a725bf, 0xeec492d9,
		0x0cd77058, 0x00000002,
	}

	c.Assert(scalarAdd(s1, s2), DeepEquals, expected)
}

func (s *Ed448Suite) Test_ScalarHalve(c *C) {
	expected := [scalarWords]uint32{6}

	c.Assert(scalarHalve([scalarWords]uint32{12}, [scalarWords]uint32{4}),
		DeepEquals,
		expected)
}

//func (s *Ed448Suite) Test_ScalarMul(c *C) {
//	x := &bigNumber{
//		0xffb823a3, 0xc96a3c35,
//		0x7f8ed27d, 0x087b8fb9,
//		0x1d9ac30a, 0x74d65764,
//		0xc0be082e, 0xa8cb0ae8,
//		0xa8fa552b, 0x2aae8688,
//		0x2c3dc273, 0x47cf8cac,
//		0x3b089f07, 0x1e63e807,
//	}
//
//	y := &bigNumber{
//		0xd8bedc42, 0x686eb329,
//		0xe416b899, 0x17aa6d9b,
//		0x1e30b38b, 0x188c6b1a,
//		0xd099595b, 0xbc343bcb,
//		0x1adaa0e7, 0x24e8d499,
//		0x8e59b308, 0x0a92de2d,
//		0xcae1cb68, 0x16c5450a,
//	}
//
//	expected := &bigNumber{
//		0xa18d010a, 0x1f5b3197,
//		0x994c9c2b, 0x6abd26f5,
//		0x08a3a0e4, 0x36a14920,
//		0x74e9335f, 0x07bcd931,
//		0xf2d89c1e, 0xb9036ff6,
//		0x203d424b, 0xfccd61b3,
//		0x4ca389ed, 0x31e055c1,
//	}
//
//	c.Assert(scalarMul(x, y), DeepEquals, expected)
//}
