package ed448

import (
	"crypto/elliptic"
	"crypto/subtle"
	"math/big"
	"sync"

	curve25519p "golang.org/x/crypto/curve25519"
)

// CurveParams contains the parameters of an elliptic curve and also provides
// a generic, non-constant time implementation of Curve.
// These are the Montgomery params.
type CurveParams struct {
	ep *elliptic.CurveParams
}

// Curve25519Params contains the parameters of an elliptic curve and also provides
// a generic, non-constant time implementation of Curve.
// These are the Montgomery params.
type Curve25519Params struct {
	ep *elliptic.CurveParams
}

// EdwardsCurveParams contains the parameters of an elliptic curve and also provides
// a generic, non-constant time implementation of Curve.
// These are the Edwards params.
type EdwardsCurveParams struct {
	P       *big.Int // the order of the underlying finite field
	N       *big.Int // the prime order of the base point
	D       *big.Int // the non-zero element
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

// A GoldilocksCurve represents the curve448.
type GoldilocksCurve interface {
	elliptic.Curve
}

// A Curve25519 represents the curve25519.
type Curve25519 interface {
	elliptic.Curve
}

// A GoldilocksEdCurve represents Goldilocks edwards448.
// This uses the decaf technique
type GoldilocksEdCurve interface {
	// Params returns the parameters for the curve.
	Params() *EdwardsCurveParams
	// IsOnCurveEdwards reports whether the given p lies on the curve.
	IsOnCurve(p Point) bool
	// AddEdwards returns the sum of p and q
	Add(p, q Point) Point
	// DoubleEdwards returns 2*p
	Double(p Point) Point
	// ScalarMultEdwards returns k*(p) where k is an scalar.
	ScalarMult(p Point, k Scalar) Point
	// ScalarBaseMultEdwards returns k*G, where G is the base point of the group
	// and k is an scalar
	ScalarBaseMult(k Scalar) Point
}

// Params returns the parameters for the curve.
func (curve *CurveParams) Params() *elliptic.CurveParams {
	return curve.ep
}

// Params returns the parameters for the curve.
func (curve *Curve25519Params) Params() *elliptic.CurveParams {
	return curve.ep
}

// IsOnCurve verifies if a given point in montgomery is valid
// v^2 = u^3 + A*u^2 + u
func (curve *CurveParams) IsOnCurve(x, y *big.Int) bool {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)

	t0.Mul(x, x)
	t0.Mul(t0, curve.ep.B)

	t2.Mul(x, x)
	t2.Mul(t2, x)

	t0.Add(t0, t2)
	t0.Add(t0, x)
	t0.Mod(t0, curve.ep.P)

	t1.Mul(y, y)
	t1.Mod(t1, curve.ep.P)

	return t0.Cmp(t1) == 0
}

// IsOnCurve verifies if a given point in montgomery is valid
// v^2 = u^3 + A*u^2 + u
func (curve *Curve25519Params) IsOnCurve(x, y *big.Int) bool {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)

	t0.Mul(x, x)
	t0.Mul(t0, curve.ep.B)

	t2.Mul(x, x)
	t2.Mul(t2, x)

	t0.Add(t0, t2)
	t0.Add(t0, x)
	t0.Mod(t0, curve.ep.P)

	t1.Mul(y, y)
	t1.Mod(t1, curve.ep.P)

	return t0.Cmp(t1) == 0
}

func inv(curve *CurveParams, x *big.Int) *big.Int {
	pMinus2 := big.NewInt(2)
	pMinus2.Sub(curve.ep.P, pMinus2)

	return x.Exp(x, pMinus2, curve.ep.P)
}

func isZero(a *big.Int) bool {
	return a.Sign() == 0
}

func isEqual(x, y *big.Int) bool {
	return isZero(new(big.Int).Sub(x, y))
}

func cMov(x, y *big.Int, b bool) *big.Int {
	z := new(big.Int)

	if b {
		z.Set(y)
	} else {
		z.Set(x)
	}

	return z
}

func isSquare(curve *CurveParams, x *big.Int) bool {
	pMinus1div2 := big.NewInt(1)
	pMinus1div2.Sub(curve.ep.P, pMinus1div2)
	pMinus1div2.Rsh(pMinus1div2, 1)

	return isEqual(new(big.Int).Exp(x, pMinus1div2, curve.ep.P), new(big.Int).SetInt64(1))
}

func sqrt(curve *CurveParams, x *big.Int) *big.Int {
	e := big.NewInt(1)
	e.Add(curve.ep.P, e)
	e.Rsh(e, 2)

	return new(big.Int).Exp(x, e, curve.ep.P)
}

func sgn0LE(x *big.Int) int {
	return 1 - 2*int(x.Bit(0))
}

// Add adds two points in montgomery
// x3 = ((y2-y1)^2/(x2-x1)^2)-A-x1-x2
// y3 = (2*x1+x2+a)*(y2-y1)/(x2-x1)-b*(y2-y1)3/(x2-x1)3-y1
// See: https://www.hyperelliptic.org/EFD/g1p/auto-montgom.html
func (curve *CurveParams) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	if x1.Sign() == 0 || y1.Sign() == 0 {
		return x2, y2
	}

	if x2.Sign() == 0 || y2.Sign() == 0 {
		return x1, y1
	}

	if x1.Cmp(x2) == 0 {
		tmp := new(big.Int)
		tmp.Neg(y1).Mod(tmp, curve.ep.P)

		if y2.Cmp(tmp) == 0 {
			y.SetInt64(1)
			return x, y
		}
	}

	t0.Sub(y2, y1)
	t1.Sub(x2, x1)
	t1 = inv(curve, t1)
	t2.Mul(t0, t1)

	t0.Mul(t2, t2)
	t0.Mul(t0, new(big.Int).SetInt64(1))
	t0.Sub(t0, curve.ep.B)
	t0.Sub(t0, x1)
	x.Sub(t0, x2)

	t0.Sub(x1, x)
	t0.Mul(t0, t2)
	y.Sub(t0, y1)

	x.Mod(x, curve.ep.P)
	y.Mod(y, curve.ep.P)

	return x, y
}

// Add adds two points in montgomery
// x3 = ((y2-y1)^2/(x2-x1)^2)-A-x1-x2
// y3 = (2*x1+x2+a)*(y2-y1)/(x2-x1)-b*(y2-y1)3/(x2-x1)3-y1
// See: https://www.hyperelliptic.org/EFD/g1p/auto-montgom.html
func (curve *Curve25519Params) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	if x1.Sign() == 0 || y1.Sign() == 0 {
		return x2, y2
	}

	if x2.Sign() == 0 || y2.Sign() == 0 {
		return x1, y1
	}

	t0.Sub(y2, y1)
	t1.Sub(x2, x1)
	t1 = new(big.Int).ModInverse(t1, curve.ep.P)
	t2.Mul(t0, t1)

	t0.Mul(t2, t2)
	t0.Mul(t0, new(big.Int).SetInt64(1))
	t0.Sub(t0, curve.ep.B)
	t0.Sub(t0, x1)
	x.Sub(t0, x2)

	t0.Sub(x1, x)
	t0.Mul(t0, t2)
	y.Sub(t0, y1)

	x.Mod(x, curve.ep.P)
	y.Mod(y, curve.ep.P)

	return x, y
}

// Double doubles two points in montgomery
// x3 = b*(3*x12+2*a*x1+1)2/(2*b*y1)2-a-x1-x1
// y3 = (2*x1+x1+a)*(3*x12+2*a*x1+1)/(2*b*y1)-b*(3*x12+2*a*x1+1)3/(2*b*y1)3-y1
// See: https://www.hyperelliptic.org/EFD/g1p/auto-montgom.html
func (curve *CurveParams) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	if x1.Sign() == 0 && y1.Sign() == 0 {
		return x1, y1
	}

	x := new(big.Int)
	y := new(big.Int)

	if y1.Sign() == 0 {
		y.SetInt64(1)
		return x, y
	}

	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)

	t0.Mul(new(big.Int).SetInt64(3), x1)
	t1.Mul(new(big.Int).SetInt64(2), curve.ep.B)
	t0.Add(t0, t1)
	t0.Mul(t0, x1)
	t1.Add(t0, new(big.Int).SetInt64(1))

	t0.Mul(new(big.Int).SetInt64(2), new(big.Int).SetInt64(1))
	t0.Mul(t0, y1)
	t0 = inv(curve, t0)
	t2.Mul(t1, t0)

	t0.Mul(t2, t2)
	t0.Mul(t0, new(big.Int).SetInt64(1))
	t0.Sub(t0, curve.ep.B)
	t0.Sub(t0, x1)
	x.Sub(t0, x1)

	t0.Sub(x1, x)
	t0.Mul(t0, t2)
	y.Sub(t0, y1)

	x.Mod(x, curve.ep.P)
	y.Mod(y, curve.ep.P)

	return x, y
}

// Double doubles two points in montgomery
// x3 = b*(3*x12+2*a*x1+1)2/(2*b*y1)2-a-x1-x1
// y3 = (2*x1+x1+a)*(3*x12+2*a*x1+1)/(2*b*y1)-b*(3*x12+2*a*x1+1)3/(2*b*y1)3-y1
// See: https://www.hyperelliptic.org/EFD/g1p/auto-montgom.html
func (curve *Curve25519Params) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	if x1.Sign() == 0 && y1.Sign() == 0 {
		return x1, y1
	}

	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	t0.Mul(new(big.Int).SetInt64(3), x1)
	t1.Mul(new(big.Int).SetInt64(2), curve.ep.B)
	t0.Add(t0, t1)
	t0.Mul(t0, x1)
	t1.Add(t0, new(big.Int).SetInt64(1))

	t0.Mul(new(big.Int).SetInt64(2), new(big.Int).SetInt64(1))
	t0.Mul(t0, y1)
	t0 = new(big.Int).ModInverse(t0, curve.ep.P)
	t2.Mul(t1, t0)

	t0.Mul(t2, t2)
	t0.Mul(t0, new(big.Int).SetInt64(1))
	t0.Sub(t0, curve.ep.B)
	t0.Sub(t0, x1)
	x.Sub(t0, x1)

	t0.Sub(x1, x)
	t0.Mul(t0, t2)
	y.Sub(t0, y1)

	x.Mod(x, curve.ep.P)
	y.Mod(y, curve.ep.P)

	return x, y
}

// Basepoint is the generator for curve448 in montgomery form
var Basepoint []byte

var basePoint = [32]byte{5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func init() { Basepoint = basePoint[:] }

func checkBasepoint() {
	if subtle.ConstantTimeCompare(Basepoint, []byte{
		0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}) != 1 {
		panic("curve448: global Basepoint value was modified")
	}
}

// LadderScalarMult returns k*(Bx,By) where k is a number in little-endian form.
// This uses the montgomery ladder
func LadderScalarMult(curve GoldilocksCurve, x1, y1 *big.Int, k []byte) (*big.Int, *big.Int) {
	var dst [x448FieldBytes]byte
	var ok bool
	s := [x448FieldBytes]byte{}
	uB := [x448FieldBytes]byte{}

	b := x1.Bytes()
	copy(s[:], k)
	copy(uB[:], b)

	if &uB[0] == &Basepoint[0] {
		checkBasepoint()
		dst = x448BasePointScalarMul(s[:])
	} else {
		var zero [x448FieldBytes]byte

		dst, ok = x448ScalarMul(uB[:], s[:])
		if !ok {
			return nil, nil
		}
		if subtle.ConstantTimeCompare(dst[:], zero[:]) == 1 {
			return nil, nil // low order point
		}
	}

	u := new(big.Int).SetBytes(dst[:])
	v := new(big.Int)

	return u, v
}

// ScalarMult returns k*(Bx,By) where k is a number in little-endian form.
// This uses the double and add method
func (curve *CurveParams) ScalarMult(x1, y1 *big.Int, k []byte) (*big.Int, *big.Int) {
	x2, y2, _ := ToWeierstrassCurve(curve.Params().P, x1, y1)
	x, y := new(big.Int), new(big.Int)

	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			x, y = curve.Double(x, y)
			if byte&0x80 == 0x80 {
				x, y = curve.Add(x2, y2, x, y)
			}
			byte <<= 1
		}
	}

	x3, y3 := ToMontgomeryCurve(x, y)

	return x3, y3
}

// ScalarMult returns k*(Bx,By) where k is a number in little-endian form.
func (curve *Curve25519Params) ScalarMult(x1, y1 *big.Int, k []byte) (*big.Int, *big.Int) {
	var dst [32]byte
	s := [32]byte{}
	uB := [32]byte{}

	b := x1.Bytes()
	copy(s[:], k)
	copy(uB[:], b)

	curve25519p.ScalarMult(&dst, &s, &uB)

	u := new(big.Int).SetBytes(dst[:])
	v := new(big.Int)

	return u, v
}

// LadderScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func LadderScalarBaseMult(curve GoldilocksCurve, k []byte) (*big.Int, *big.Int) {
	var dst [x448FieldBytes]byte
	s := [x448FieldBytes]byte{}

	copy(s[:], k)
	dst = x448BasePointScalarMul(s[:])

	u := new(big.Int).SetBytes(dst[:])
	v := new(big.Int)

	return u, v
}

// ScalarBaseMult returns k*G, where G is the base point of the group
// and k is an integer in big-endian form.
func (curve *CurveParams) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	return curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, k)
}

// ScalarBaseMult returns k*(Bx,By) where k is a number in little-endian form.
func (curve *Curve25519Params) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	var dst [32]byte
	s := [32]byte{}

	copy(s[:], k)

	curve25519p.ScalarBaseMult(&dst, &s)

	u := new(big.Int).SetBytes(dst[:])
	v := new(big.Int)

	return u, v
}

// ToWeierstrassCurve converts from Montgomery form to Weierstrass
func ToWeierstrassCurve(p, u, v *big.Int) (*big.Int, *big.Int, *big.Int) {
	invB := new(big.Int)
	a := new(big.Int)
	b := new(big.Int)
	c := new(big.Int).SetInt64(1)

	invB.ModInverse(new(big.Int).SetInt64(1), p)
	a.Mul(u, invB)
	b.Mul(v, invB)

	return a, b, c
}

// ToMontgomeryCurve converts from Weierstrass to Montgomery
func ToMontgomeryCurve(x, y *big.Int) (*big.Int, *big.Int) {
	u := new(big.Int)
	v := new(big.Int)

	b := new(big.Int).SetInt64(1)
	u.Mul(x, b)
	v.Mul(y, b)

	return u, v
}

// MapToCurve calculates a point on the elliptic curve from an element of the finite field F. This implements Elligator2,
// according to https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-05, section 6.7.1.1.
func (curve *CurveParams) MapToCurve(u *big.Int) (*big.Int, *big.Int) {
	t1, x1, x2, gx1, gx2, y2, x, y := new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	var e1, e2, e3 bool
	z := new(big.Int).SetInt64(-1)

	t1.Mul(u, u)                                // t1 = u^2
	t1.Mul(z, t1)                               // Z * u^2
	e1 = isEqual(t1, new(big.Int).SetInt64(-1)) // Z * u^2 == -1
	t1 = cMov(t1, new(big.Int).SetInt64(0), e1) // if t1 == -1, set t1 = 0
	x1.Add(t1, new(big.Int).SetInt64(1))        // x1 = t1 + 1
	x1 = inv(curve, x1)                         // x1 = inv0(x1)
	x1.Mul(new(big.Int).Neg(curve.ep.B), x1)    // x1 = -A / (1 + Z * u^2)
	gx1.Add(x1, curve.ep.B)                     // gx1 = x1 + A
	gx1.Mul(gx1, x1)                            // gx1 = gx1 * x1
	gx1.Add(gx1, new(big.Int).SetInt64(1))      // gx1 = gx1 + B
	gx1.Mul(gx1, x1)                            // gx1 = x1^3 + A * x1^2 + B * x1

	x2.Sub(new(big.Int).Neg(x1), curve.ep.B) //x2 = -x1 - A
	gx2.Mul(t1, gx1)                         // gx2 = t1 * gx1
	e2 = isSquare(curve, gx1)                // e2 = is_square(gx1)
	x = cMov(x2, x1, e2)                     // If is_square(gx1), x = x1, else x = x2
	y2 = cMov(gx2, gx1, e2)                  // If is_square(gx1), y2 = gx1, else y2 = gx2
	y = sqrt(curve, y2)                      // y = sqrt(y2)
	e3 = sgn0LE(u) == sgn0LE(y)              // Fix sign of y: e3 = sgn0(u) == sgn0(y)
	y = cMov(new(big.Int).Neg(y), y, e3)     // y = CMOV(-y, y, e3)

	x.Mod(x, curve.ep.P)
	y.Mod(y, curve.ep.P)

	return x, y
}

// Params returns the parameters for the curve.
func (curve *EdwardsCurveParams) Params() *EdwardsCurveParams {
	return curve
}

// IsOnCurve reports whether the given point (p) lies on the curve.
func (curve *EdwardsCurveParams) IsOnCurve(p Point) bool {
	return p.(*twExtendedPoint).isOnCurve()
}

// Add gives the sum of two points (p, q) and produces a third point (p).
func (curve *EdwardsCurveParams) Add(p, q Point) Point {
	r := &twExtendedPoint{}
	r.add(p.(*twExtendedPoint), q.(*twExtendedPoint))

	return r
}

// Double gives the doubling of a point (p).
func (curve *EdwardsCurveParams) Double(p Point) Point {
	p.(*twExtendedPoint).double()

	return p
}

// ScalarMult returns the multiplication of a given point (p) by a given
// scalar (a): p * k.
func ScalarMult(p Point, k Scalar) Point {
	return pointScalarMul(p.(*twExtendedPoint), k.(*scalar))
}

// ScalarBaseMult returns the multiplication of a given scalar (k) by the
// precomputed base point of the curve: basePoint * k.
func ScalarBaseMult(k Scalar) Point {
	return precomputedScalarMul(k.(*scalar))
}

var initonce sync.Once
var curve448 *CurveParams
var curve25519 *Curve25519Params
var ed448 *EdwardsCurveParams

func initAll() {
	initCurve25519()
	initCurve448()
	initEd448()
}

func initCurve448() {
	// See https://safecurves.cr.yp.to/field.html and https://tools.ietf.org/html/rfc7748#section-4.2
	P, _ := new(big.Int).SetString("726838724295606890549323807888004534353641360687318060281490199180612328166730772686396383698676545930088884461843637361053498018365439", 10)
	N, _ := new(big.Int).SetString("181709681073901722637330951972001133588410340171829515070372549795146003961539585716195755291692375963310293709091662304773755859649779", 10)
	A, _ := new(big.Int).SetString("156326", 10)
	Gu, _ := new(big.Int).SetString("5", 10)
	Gv, _ := new(big.Int).SetString("355293926785568175264127502063783334808976399387714271831880898435169088786967410002932673765864550910142774147268105838985595290606362", 10)
	curve448 = &CurveParams{&elliptic.CurveParams{Name: "curve-448",
		P:       P,
		N:       N,
		B:       A,
		Gx:      Gu,
		Gy:      Gv,
		BitSize: 448,
	}}
}

func initCurve25519() {
	// See https://safecurves.cr.yp.to/field.html and https://tools.ietf.org/html/rfc7748#section-4.2
	P, _ := new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)
	N, _ := new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)
	A, _ := new(big.Int).SetString("486662", 10)
	Gu, _ := new(big.Int).SetString("9", 10)
	Gv, _ := new(big.Int).SetString("14781619447589544791020593568409986887264606134616475288964881837755586237401", 10)
	curve25519 = &Curve25519Params{&elliptic.CurveParams{Name: "curve-25519",
		P:       P,
		N:       N,
		B:       A,
		Gx:      Gu,
		Gy:      Gv,
		BitSize: 256,
	}}
}

func initEd448() {
	// See https://safecurves.cr.yp.to/field.html and https://tools.ietf.org/html/rfc7748#section-4.2
	ed448 = &EdwardsCurveParams{Name: "ed-448"}
	ed448.P, _ = new(big.Int).SetString("726838724295606890549323807888004534353641360687318060281490199180612328166730772686396383698676545930088884461843637361053498018365439", 10)
	ed448.N, _ = new(big.Int).SetString("181709681073901722637330951972001133588410340171829515070372549795146003961539585716195755291692375963310293709091662304773755859649779", 10)
	ed448.D, _ = new(big.Int).SetString("-39081", 10)
	ed448.Gx, _ = new(big.Int).SetString("224580040295924300187604334099896036246789641632564134246125461686950415467406032909029192869357953282578032075146446173674602635247710", 10)
	ed448.Gy, _ = new(big.Int).SetString("298819210078481492676017930443930673437544040154080242095928241372331506189835876003536878655418784733982303233503462500531545062832660", 10)
	ed448.BitSize = 448
}

// Curve448 returns a Curve which implements curve448
func Curve448() GoldilocksCurve {
	initonce.Do(initAll)
	return curve448
}

// CurveP25519 returns a Curve which implements curve448
func CurveP25519() Curve25519 {
	initonce.Do(initAll)
	return curve25519
}
