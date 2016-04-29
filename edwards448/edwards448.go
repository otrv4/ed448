package edwards448

import "math/big"

type Curve interface {
	// Params returns the parameters for the curve.
	Params() *CurveParams
	// IsOnCurve reports whether the given (x,y) lies on the curve.
	IsOnCurve(x, y *big.Int) bool
	// Add returns the sum of (x1,y1) and (x2,y2)
	Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int)
	// Double returns 2*(x,y)
	Double(x1, y1 *big.Int) (x, y *big.Int)
	// Multiply performs a scalar multiplication and returns k*(Bx,By) where k is a number in big-endian form.
	Multiply(x1, y1 *big.Int, k int) (x, y *big.Int)
	// // ScalarBaseMult returns k*G, where G is the base point of the group
	// // and k is an integer in big-endian form.
	// ScalarBaseMult(k []byte) (x, y *big.Int)
}

type CurveParams struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

type ed448Curve struct {
	*CurveParams
}

var ed448 ed448Curve

func init() {
	ed448.CurveParams = &CurveParams{Name: "Ed-448"}
	ed448.P, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	ed448.N, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3", 16)
	ed448.B, _ = new(big.Int).SetString("-39081", 10)
	ed448.Gx, _ = new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
	ed448.Gy, _ = new(big.Int).SetString("13", 16)
	ed448.BitSize = 448
}

func Ed448() Curve {
	return ed448
}

func (c *CurveParams) Params() *CurveParams {
	return c
}

func (c *CurveParams) IsOnCurve(x, y *big.Int) bool {
	// x² + y² = 1 + bx²y²
	x2 := new(big.Int).Mul(x, x)
	y2 := new(big.Int).Mul(y, y)

	x2y2 := new(big.Int).Mul(x2, y2)
	bx2y2 := new(big.Int).Mul(c.B, x2y2)

	left := new(big.Int).Add(x2, y2)
	left.Mod(left, c.P)
	right := new(big.Int).Add(big.NewInt(1), bx2y2)
	right.Mod(right, c.P)

	return left.Cmp(right) == 0
}

func (c *CurveParams) Add(x1, y1, x2, y2 *big.Int) (x3, y3 *big.Int) {
	// x² + y² = 1 + bx²y²
	// x3 =  x1y2 + y1x2 / 1 + bx1x2y1y2
	// y3 =  y1y2 - x1x2 / 1 - bx1x2y1y2

	bx1x2y1y2 := new(big.Int).Mul(
		c.B, new(big.Int).Mul(x1, new(big.Int).Mul(x2, new(big.Int).Mul(y1, y2))))
	bx1x2y1y2.Mod(bx1x2y1y2, c.P)

	x3 = new(big.Int).Mul(x1, y2)
	x3.Add(x3, new(big.Int).Mul(x2, y1))
	x3.Mod(x3, c.P)
	divisor := new(big.Int).ModInverse(new(big.Int).Mod(new(big.Int).Add(big.NewInt(1), bx1x2y1y2), c.P), c.P)
	x3.Mul(x3, divisor)

	y3 = new(big.Int).Mul(y1, y2)
	y3.Sub(y3, new(big.Int).Mul(x1, x2))
	y3.Mod(y3, c.P)
	divisor = new(big.Int).ModInverse(new(big.Int).Mod(new(big.Int).Sub(big.NewInt(1), bx1x2y1y2), c.P), c.P)
	y3.Mul(y3, divisor)

	return
}

func (c *CurveParams) Double(x1, y1 *big.Int) (x3, y3 *big.Int) {
	// x² + y² = 1 + bx²y²
	// x3 =  2xy / 1 + bx²y² = 2xy / x² + y²
	// y3 =  y² - x² / 1 - bx²y² = y² - x² / 2 - x² - y²

	x2plusy2 := new(big.Int).Add(new(big.Int).Mul(x1, x1), new(big.Int).Mul(y1, y1))
	x2plusy2.Mod(x2plusy2, c.P)

	x3 = new(big.Int).Mul(x1, y1)
	x3.Lsh(x3, 1) // x3 = 2xy
	x3.Mod(x3, c.P)
	divisor := new(big.Int).ModInverse(x2plusy2, c.P)
	x3.Mul(x3, divisor) // x3 = 2xy / x² + y²

	y3 = new(big.Int).Sub(new(big.Int).Mul(y1, y1), new(big.Int).Mul(x1, x1)) // y3 = y² - x²
	y3.Mod(y3, c.P)
	divisor = new(big.Int).ModInverse(new(big.Int).Mod(new(big.Int).Sub(big.NewInt(2), x2plusy2), c.P), c.P)
	y3.Mul(y3, divisor) // y3 = y² - x² / 2 - x² - y²

	return
}

func (c *CurveParams) Multiply(x, y *big.Int, k int) (kx, ky *big.Int) {
	kx, ky = x, y
	for k > 0 {
		if k % 2 == 0 {
			kx, ky = c.Double(kx, ky)
			k = k - 2
		} else {
			kx, ky = c.Add(kx, ky, kx, ky)
			k = k - 1
		}
	}
	return
}
