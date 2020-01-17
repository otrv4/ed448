package ed448

import "math/big"

// CurveParams contains the parameters of an elliptic curve and also provides
// a generic, non-constant time implementation of Curve.
// This are the Montgomery params
type CurveParams struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	A       *big.Int // the constant of the curve equation
	Gu, Gv  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

// A GoldilocksCurve represents Goldilocks curve or edwards448.
// See https://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
type GoldilocksCurve interface {
	// Params returns the parameters for the curve.
	Params() *CurveParams
	// Params returns the edwards parameters for the curve.
	EdwardsParams() *CurveParams
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
