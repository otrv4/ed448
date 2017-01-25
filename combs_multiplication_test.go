package ed448

import . "gopkg.in/check.v1"

func (s *Ed448Suite) TestRadixScheduleForCombs(c *C) {
	scalar := Scalar{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}

	schedule := make([]uint32, scalarWords)
	scheduleScalarForCombs(schedule, scalar)
	c.Assert(schedule, DeepEquals, []uint32{
		0xfefb9893, 0x4acadc23, 0x2b57a900, 0xcddcdc54,
		0x79beae27, 0x5989b711, 0xc4d0ca21, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x80000000, 0x20000000,
	})

	scalar = Scalar{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0xA, 0xB, 0xC, 0xD, 0xE,
	}
	scheduleScalarForCombs(schedule, scalar)
	c.Assert(schedule, DeepEquals, []uint32{
		0xa94f761a, 0x390e7adb, 0xe474e157, 0x3d267b1c,
		0xa25392e2, 0xf762496f, 0x066bb82f, 0x80000005,
		0x00000004, 0x80000005, 0x00000005, 0x80000006,
		0x00000006, 0x00000007,
	})
}
