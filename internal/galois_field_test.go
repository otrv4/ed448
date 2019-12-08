package galoisfield

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448InternalSuite struct{}

var _ = Suite(&Ed448InternalSuite{})

func (s *Ed448InternalSuite) Test_NewGaloisField(c *C) {
	gf := NewGaloisField448()
	c.Assert(gf, NotNil)

	size := gf.Limb.Size()
	// For an arch of 32 for the moment
	c.Assert(size, Equals, 16)

	gf.Destroy()
}

func (s *Ed448InternalSuite) Test_GaloisField_Copy(c *C) {
	gf := NewGaloisField448()
	c.Assert(gf, NotNil)

	n := gf.Copy()
	c.Assert(n, NotNil)

	c.Assert(n, DeepEquals, gf)

	gf.Destroy()
	n.Destroy()
}
