package galoisfield

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448InternalSuite struct{}

var _ = Suite(&Ed448InternalSuite{})

func (s *Ed448InternalSuite) Test_NewGaloisField32(c *C) {
	gf := NewGaloisField448(N32Limbs)
	c.Assert(gf, NotNil)

	size := gf.Limb.Size()
	// For an arch of 32 for the moment
	c.Assert(size, Equals, 128)

	gf.Destroy()
}

func (s *Ed448InternalSuite) Test_NewGaloisField64(c *C) {
	gf := NewGaloisField448(N64Limbs)
	c.Assert(gf, NotNil)

	size := gf.Limb.Size()
	// For an arch of 64 for the moment
	c.Assert(size, Equals, 64)

	gf.Destroy()
}

func (s *Ed448InternalSuite) Test_GaloisField_Copy(c *C) {
	gf := NewGaloisField448(N32Limbs)
	c.Assert(gf, NotNil)

	n := gf.Copy()
	c.Assert(n, NotNil)

	c.Assert(n, DeepEquals, gf)

	gf.Destroy()
	n.Destroy()
}

func (s *Ed448InternalSuite) Test_GaloisField_AddRaw32(c *C) {
	tmp1 := [128]byte{0x57}
	tmp2 := [128]byte{0x83}
	tmp3 := [128]byte{0xda}

	x := NewGaloisField448FromBytes(tmp1[:])
	y := NewGaloisField448FromBytes(tmp2[:])
	exp := NewGaloisField448FromBytes(tmp3[:])

	v := AddRaw32(x, y)

	c.Assert(v.limbs32(), DeepEquals, exp.limbs32())
	c.Assert(v.Limb.Size(), Equals, 128)

	x.Destroy()
	y.Destroy()
	exp.Destroy()
}

func (s *Ed448InternalSuite) Test_GaloisField_AddRaw64(c *C) {
	tmp1 := [64]byte{0x01}
	tmp2 := [64]byte{0x02}
	tmp3 := [64]byte{0x03}

	x := NewGaloisField448FromBytes(tmp1[:])
	y := NewGaloisField448FromBytes(tmp2[:])
	exp := NewGaloisField448FromBytes(tmp3[:])

	v := AddRaw64(x, y)

	c.Assert(v.limbs64(), DeepEquals, exp.limbs64())
	c.Assert(v.Limb.Size(), Equals, 64)

	x.Destroy()
	y.Destroy()
	exp.Destroy()
}

func (s *Ed448InternalSuite) Test_GaloisField_SubRaw32(c *C) {
	tmp1 := [128]byte{0x10}
	tmp2 := [128]byte{0x05}
	tmp3 := [128]byte{0x0b}

	x := NewGaloisField448FromBytes(tmp1[:])
	y := NewGaloisField448FromBytes(tmp2[:])
	exp := NewGaloisField448FromBytes(tmp3[:])

	v := SubRaw32(x, y)

	c.Assert(v.limbs32(), DeepEquals, exp.limbs32())
	c.Assert(v.Limb.Size(), Equals, 128)

	x.Destroy()
	y.Destroy()
	exp.Destroy()
}

func (s *Ed448InternalSuite) Test_GaloisField_SubRaw64(c *C) {
	tmp1 := [64]byte{0x02}
	tmp2 := [64]byte{0x01}
	tmp3 := [64]byte{0x01}

	x := NewGaloisField448FromBytes(tmp1[:])
	y := NewGaloisField448FromBytes(tmp2[:])
	exp := NewGaloisField448FromBytes(tmp3[:])

	v := SubRaw64(x, y)

	c.Assert(v.limbs64(), DeepEquals, exp.limbs64())
	c.Assert(v.Limb.Size(), Equals, 64)

	x.Destroy()
	y.Destroy()
	exp.Destroy()
}
