package ed448

import "math/big"

	// N is the order of the base point
	// N, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3", 16)
	// BitSize is the size of the underlying field
	// BitSize = 448

type Point struct {
	x, y *big.Int
}

type Edwards448Curve struct {
	p         *big.Int
	b         *big.Int
	basePoint *Point
}

func (p *Point) SetX(x *big.Int) {
	p.x = x
}

func (p *Point) SetY(y *big.Int) {
	p.y = y
}

func (p Point) GetX() *big.Int {
	return p.x
}

func (p Point) GetY() *big.Int {
	return p.y
}

func (c *Edwards448Curve) setFieldSize(p *big.Int) {
	c.p = p
}

func (c *Edwards448Curve) setConstant(b *big.Int) {
	c.b = b
}

func (c *Edwards448Curve) GetFieldSize() *big.Int {
	if nil == c.p {
		p, _ := new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
		c.setFieldSize(p)
	}
	return c.p
}

func (c *Edwards448Curve) GetConstant() *big.Int {
	if nil == c.b {
		b, _ := new(big.Int).SetString("-39081", 10)
		c.setConstant(b)
	}
	return c.b
}

func (c *Edwards448Curve) GetBasePoint() *Point {
	if nil == c.basePoint {
		c.basePoint = new(Point)
		x, _ := new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
		y, _ := new(big.Int).SetString("13", 16)
		c.basePoint.SetX(x)
		c.basePoint.SetY(y)
	}
	return c.basePoint
}

func (c Edwards448Curve) CheckMembershipOf(p *Point) bool {
	// x² + y² = 1 + bx²y²
	bp := c.GetBasePoint()
	x := bp.GetX()
	y := bp.GetY()

	x2 := new(big.Int).Mul(x, x)
	x2.Mod(x2, c.GetFieldSize())

	y2 := new(big.Int).Mul(y, y)
	y2.Mod(y2, c.GetFieldSize())

	x2y2 := new(big.Int).Mul(x2, y2)
	x2y2.Mod(x2y2, c.GetFieldSize())

	// TODO: we may use shifting to multiply
	bx2y2 := new(big.Int).Mul(x2y2, c.GetConstant())
	bx2y2.Mod(bx2y2, c.GetFieldSize())

	left := new(big.Int).Add(x2, y2)
	right := new(big.Int).Add(big.NewInt(1), bx2y2)

	return left.Cmp(right) == 0
}

/*
func (curve *CurveParams) Add(x1, y1, x2, y2 *big.Int) (x3, y3 *big.Int) {
	x3 = new(big.Int).Mul(x1, y2)
	x3.Add(x3, new(big.Int).Mul(x2, y1))

	y3 = new(big.Int).Mul(y1, y2)
	y3.Sub(x3, new(big.Int).Mul(x1, x2))

	// TODO: Consider mod after each mul
	bx1x2x2y2 := new(big.Int).Mul(
		curve.B, new(big.Int).Mul(x1, new(big.Int).Mul(x2, new(big.Int).Mul(y1, y2))))
	bx1x2x2y2.Mod(bx1x2x2y2, curve.P)

	x3.Div(x3, new(big.Int).Add(big.NewInt(1), bx1x2x2y2))
	x3.Mod(x3, curve.P)

	y3.Div(y3, new(big.Int).Sub(big.NewInt(1), bx1x2x2y2))
	y3.Mod(y3, curve.P)

	return
}

func (curve *CurveParams) Double(x1, y1 *big.Int) (x3, y3 *big.Int) {
	// x3 =  2xy / 1 + bx²y² = 2xy / x² + y²
	x3 = new(big.Int).Mul(x1, y1)
	x3 = x3.Lsh(x3, 1) // x3 = 2xy
	x2plusy2 := new(big.Int).Add(new(big.Int).Mul(x1, x1), new(big.Int).Mul(y1, y1))
	x2plusy2 = x2plusy2.Mod(x2plusy2, curve.P)
	x3 = x3.Div(x3, x2plusy2) // x3 = 2xy / x² + y²
	x3 = x3.Mod(x3, curve.P)

	// y3 =  y² - x² / 1 + bx²y² = y² - x² / 1 - x² - y²
	y3 = new(big.Int).Sub(new(big.Int).Mul(y1, y1), new(big.Int).Mul(x1, x1))
	y3 = y3.Mod(y3, curve.P)
	y3 = y3.Div(y3, new(big.Int).Sub(big.NewInt(1), x2plusy2)) // y3 = y² - x² / 1 - x² - y²
	y3 = y3.Mod(y3, curve.P)

	return curve.Add(x1, y1, x1, y1)
}
*/
