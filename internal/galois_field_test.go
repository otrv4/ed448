package galoisfield

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448InternalSuite struct{}

var _ = Suite(&Ed448InternalSuite{})

func (s *Ed448InternalSuite) Test_NewGaloisField(c *C) {
	gf := NewGaloisField448(N32Limbs)
	c.Assert(gf, NotNil)

	size := gf.Limb.Size()
	// For an arch of 32 for the moment
	c.Assert(size, Equals, 128)

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

func (s *Ed448InternalSuite) Test_GaloisField_AddRaw(c *C) {
	tmp1 := [128]byte{0x57}
	tmp2 := [128]byte{0x83}
	tmp3 := [128]byte{0xda}

	x := NewGaloisField448FromBytes(tmp1[:])
	y := NewGaloisField448FromBytes(tmp2[:])
	exp := NewGaloisField448FromBytes(tmp3[:])

	v := AddRaw32(x, y)

	c.Assert(v.limbs(), DeepEquals, exp.limbs())

	x.Destroy()
	y.Destroy()
	exp.Destroy()
}
