package ed448

import (
	"math/big"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) Test_IsValidMontgomeryPoint(c *C) {
	curve448 := Curve448()
	c.Assert(curve448.IsOnCurve(curve448.Params().Gu, curve448.Params().Gv), Equals, true)

	x, y := new(big.Int).SetInt64(1), new(big.Int).SetInt64(1)
	c.Assert(curve448.IsOnCurve(x, y), Equals, false)
}
