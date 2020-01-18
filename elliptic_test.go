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

func (s *Ed448Suite) Test_AddMontgomeryPoint(c *C) {
	curve448 := Curve448()
	x, y := curve448.Add(curve448.Params().Gu, curve448.Params().Gv, curve448.Params().Gu, curve448.Params().Gv)

	c.Assert(curve448.IsOnCurve(x, y), Equals, false)

	x1, y1 := new(big.Int).SetInt64(0), new(big.Int).SetInt64(0)
	baseX := curve448.Params().Gu
	baseY := curve448.Params().Gv

	x3, y3 := curve448.Add(baseX, baseY, x1, y1)
	c.Assert(x3, DeepEquals, baseX)
	c.Assert(y3, DeepEquals, baseY)
}

func (s *Ed448Suite) Test_DoubleMontgomeryPoint(c *C) {
	curve448 := Curve448()
	x1, y1 := new(big.Int).SetInt64(0), new(big.Int).SetInt64(0)
	x, y := curve448.Double(x1, y1)

	c.Assert(x.Sign(), Equals, 0)
	c.Assert(y.Sign(), Equals, 0)
}

// With RFC7748 test vectors
func (s *Ed448Suite) Test_ScalarMultMontgomeryPoint(c *C) {
	curve448 := Curve448()
	x1 := new(big.Int)
	sc := new(big.Int)
	exp := new(big.Int)

	x1, _ = new(big.Int).SetString("06fce640fa3487bfda5f6cf2d5263f8aad88334cbd07437f020f08f9814dc031ddbdc38c19c6da2583fa5429db94ada18aa7a7fb4ef8a086", 16)
	sc, _ = new(big.Int).SetString("3d262fddf9ec8e88495266fea19a34d28882acef045104d0d1aae121700a779c984c24f8cdd78fbff44943eba368f54b29259a4f1c600ad3", 16)
	y1 := new(big.Int).SetInt64(0)
	exp, _ = new(big.Int).SetString("ce3e4ff95a60dc6697da1db1d85e6afbdf79b50a2412d7546d5f239fe14fbaadeb445fc66a01b0779d98223961111e21766282f73dd96b6f", 16)

	dst := curve448.ScalarMult(x1, y1, sc.Bytes())

	c.Assert(dst, DeepEquals, exp.Bytes())

	x1, _ = new(big.Int).SetString("0fbcc2f993cd56d3305b0b7d9e55d4c1a8fb5dbb52f8e9a1e9b6201b165d015894e56c4d3570bee52fe205e28a78b91cdfbde71ce8d157db", 16)
	sc, _ = new(big.Int).SetString("203d494428b8399352665ddca42f9de8fef600908e0d461cb021f8c538345dd77c3e4806e25f46d3315c44e0a5b4371282dd2c8d5be3095f", 16)
	exp, _ = new(big.Int).SetString("884a02576239ff7a2f2f63b2db6a9ff37047ac13568e1e30fe63c4a7ad1b3ee3a5700df34321d62077e63633c575c1c954514e99da7c179d", 16)

	dst = curve448.ScalarMult(x1, y1, sc.Bytes())

	c.Assert(dst, DeepEquals, exp.Bytes())
}
