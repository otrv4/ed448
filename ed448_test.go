package ed448

import	(
	"testing"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ed448Suite struct{}

var _ = Suite(&ed448Suite{})

func (s *ed448Suite) Test_IsOnCurve_findsPoint(c *C) {
	ed448 := Ed448()
	isOnCurve := ed448.IsOnCurve(ed448.Params().Gx, ed448.Params().Gy)
	c.Assert(isOnCurve, Equals, true)
}

func TestDouble(t *testing.T) {
	ed448 := Ed448()
	x2, y2 := ed448.Double(ed448.Params().Gx, ed448.Params().Gy)
	x3, y3 := ed448.Add(ed448.Params().Gx, ed448.Params().Gy, ed448.Params().Gx, ed448.Params().Gy)
	if x2.Cmp(x3) != 0 || y2.Cmp(y3) != 0 {
		t.Errorf("FAIL")
	}
}
