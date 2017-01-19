package ed448

import (
	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) Test_ScalarAddition(c *C) {
	s1 := [scalarWords]word_t{
		0x529eec33, 0x721cf5b5,
		0xc8e9c2ab, 0x7a4cf635,
		0x44a725bf, 0xeec492d9,
		0x0cd77058, 0x00000002,
	}
	s2 := [scalarWords]word_t{0x00000001}
	expected := [scalarWords]word_t{
		0x529eec34, 0x721cf5b5,
		0xc8e9c2ab, 0x7a4cf635,
		0x44a725bf, 0xeec492d9,
		0x0cd77058, 0x00000002,
	}

	c.Assert(scalarAdd(s1, s2), DeepEquals, expected)
}

func (s *Ed448Suite) Test_ScalarHalve(c *C) {
	expected := [scalarWords]word_t{6}

	c.Assert(scalarHalve([scalarWords]word_t{12}, [scalarWords]word_t{4}),
		DeepEquals,
		expected)
}
