package ed448

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
	// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
	ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int)
	// ScalarBaseMult returns k*G, where G is the base point of the group
	// and k is an integer in big-endian form.
	ScalarBaseMult(k []byte) (x, y *big.Int)
}

type CurveParams struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

type Ed448Curve struct {
	*CurveParams
}

var ed448 Ed448Curve

func init() {
	ed448.CurveParams = &CurveParams{Name: "P-224"}
	ed448.P, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	ed448.N, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3", 16)
	ed448.B, _ = new(big.Int).SetString("-39081", 10)
	ed448.Gx, _ = new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
	ed448.Gy, _ = new(big.Int).SetString("13", 16)
	ed448.BitSize = 448

	// p224FromBig(&p224.gx, p224.Gx)
	// p224FromBig(&p224.gy, p224.Gy)
	// p224FromBig(&p224.b, p224.B)
}

func (curve *CurveParams) Params() *CurveParams {
	return curve
}

func (curve *CurveParams) IsOnCurve(x, y *big.Int) bool {
	// x² + y² = 1 + bx²y²
	x2 := new(big.Int).Mul(x, x)
	x2.Mod(x2, curve.P)

	y2 := new(big.Int).Mul(y, y)
	y2.Mod(y2, curve.P)

	x2y2 := new(big.Int).Mul(x2, y2)
	x2y2.Mod(x2y2, curve.P)

	bx2y2 := new(big.Int).Mul(x2y2, curve.B)
	bx2y2.Mod(bx2y2, curve.P)

	left := new(big.Int).Add(x2, y2)
	right := new(big.Int).Add(big.NewInt(1), bx2y2)

	return left.Cmp(right) == 0
}
