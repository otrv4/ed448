package ed448

import . "gopkg.in/check.v1"

func (s *Ed448Suite) Test_DecafLookup(c *C) {
	expA := &bigNumber{
		0x0ad825f1, 0x0d37716c, 0x0ba9552a, 0x0883870c,
		0x05c762e3, 0x08ef785f, 0x00469242, 0x06cb253e,
		0x0ee9d967, 0x07b8f17f, 0x032b52b6, 0x0a43de69,
		0x02af783c, 0x01aca9fe, 0x0ff0b680, 0x08967778,
	}

	expB := &bigNumber{
		0x0dc6c9c3, 0x06400c4c, 0x0691083f, 0x01e8c978,
		0x0f68e0c5, 0x0ad74f01, 0x072b5f6a, 0x0f7feb03,
		0x05ade13a, 0x02f60d17, 0x0221a678, 0x098ec54a,
		0x071f244e, 0x0fcfea8a, 0x0e45ded2, 0x0dea6660,
	}

	expC := &bigNumber{
		0x0a8d6752, 0x02585b4a, 0x015a2089, 0x0e62da76,
		0x01f39b68, 0x010c1c74, 0x0ced9f65, 0x0569bb1e,
		0x04daa724, 0x0ba6d09e, 0x0ef281b9, 0x07d3e20a,
		0x0ca3ffdc, 0x0bd7f65a, 0x050288a8, 0x0dea434a,
	}

	point := decafPrecompTable.lookup(word(0x09))

	c.Assert(point.a, DeepEquals, expA)
	c.Assert(point.b, DeepEquals, expB)
	c.Assert(point.c, DeepEquals, expC)

}

func (s *Ed448Suite) Test_SelectMask(c *C) {
	m := selectMask(1, 1)
	c.Assert(m, Equals, allOnes)

	m = selectMask(1, 0)
	c.Assert(m, Equals, allZeros)
}
