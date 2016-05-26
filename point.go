package ed448

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
)

var (
	bigNumOne    = mustDeserialize(serialized{1})
	curveDSigned = int64(-39081)

	//XXX Remove. It's for debugging only
	gCurvePrime, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	gCurveD        = big.NewInt(-39081)
)

// Point represents a point on the curve in a suitable coordinate system
type Point interface {
	OnCurve() bool
	Add(Point) Point
	Double() Point
	//ReAdd(Point) Point //????
	//Affine() *Affine
}

// NewPoint instantiates a new point in a suitable coordinate system.
// The x and y coordinates must be affine coordinates in little-endian
//XXX This should probably receive []byte{}
func NewPoint(x serialized, y serialized) (p Point, e error) {
	xN, ok1 := deserialize(x)
	yN, ok2 := deserialize(y)

	p = newHomogeneousProjective(xN, yN)

	if !(ok1 && ok2) {
		e = errors.New("invalid coordinates")
	}

	return
}

//XXX This should be removed
type Affine [2]*bigNumber

func (aP *Affine) OnCurve() bool {
	// x² + y² - 1 - bx²y² = 0
	a := aP[0]
	b := aP[1]
	x2 := karatsubaMul(a, a)
	y2 := karatsubaMul(b, b)

	x2y2 := karatsubaMul(x2, y2)
	dx2y2 := x2y2.mulWSignedCurveConstant(x2y2, curveDSigned)
	dx2y2.weakReduce()

	r := sumRadix(x2, y2)
	r = subRadix(r, bigNumOne)
	r = subRadix(r, dx2y2)

	r.strongReduce()
	return r.zero()
}

func (aP *Affine) Double() Point {
	return nil
}

func (aP *Affine) Add(Point) Point {
	return nil
}

//HP(X : Y : Z) = Affine(X/Z, Y/Z), Z ≠ 0
type homogeneousProjective [3]*bigNumber

//Affine to Homogeneous Projective
func newHomogeneousProjective(x *bigNumber, y *bigNumber) *homogeneousProjective {
	return &homogeneousProjective{
		x.copy(),         // X * Z
		y.copy(),         // Y * Z
		bigNumOne.copy(), // Z = 1
	}
}

func (hP *homogeneousProjective) String() string {
	return fmt.Sprintf("X: %s, Y: %s, Z: %s", hP[0], hP[1], hP[2])
}

func (hP *homogeneousProjective) OnCurve() bool {
	// (x² + y²)z² - z^4 - dx²y² = 0
	x := hP[0]
	y := hP[1]
	z := hP[2]

	x2 := karatsubaMul(x, x)
	y2 := karatsubaMul(y, y)
	z2 := karatsubaMul(z, z)
	z4 := karatsubaMul(z2, z2)

	x2y2 := karatsubaMul(x2, y2)
	dx2y2 := x2y2.mulWSignedCurveConstant(x2y2, curveDSigned)
	dx2y2.weakReduce()

	r := sumRadix(x2, y2)
	r = karatsubaMul(r, z2)
	r = subRadix(r, z4)
	r = subRadix(r, dx2y2)

	r.strongReduce()
	return r.zero()
}

func rev(in []byte) []byte {
	r := make([]byte, len(in), len(in))

	for i, ni := range in {
		r[len(in)-i-1] = ni
	}

	return r
}

func compareNumbers(label string, n *bigNumber, b *big.Int) {
	s := [56]byte{}
	serialize(s[:], n)

	r := rev(s[:])
	bs := b.Bytes()

	for i := len(r) - len(bs); i > 0; i-- {
		bs = append([]byte{0}, bs...)
	}

	if !bytes.Equal(r, bs) {
		fmt.Printf("%s does not match!\n\t%#v\n\n vs\n\n\t%#v\n", label, r, bs)
	}
}

// See Hisil, formula 5.1
func (hP *homogeneousProjective) Double() Point {
	x1 := hP[0]
	y1 := hP[1]
	z1 := hP[2]

	//gx, _ := new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
	//gy := big.NewInt(0x13)
	//gz := big.NewInt(1)

	//compareNumbers("x", x1, gx)
	//compareNumbers("y", y1, gy)
	//compareNumbers("z", z1, gz)

	//gb := new(big.Int).Add(gx, gy)
	//gb = gb.Mul(gb, gb)
	//gb = gb.Mod(gb, gCurvePrime)

	//gc := new(big.Int).Mul(gx, gx)
	//gc = gc.Mod(gc, gCurvePrime)

	//gd := new(big.Int).Mul(gy, gy)
	//gd = gd.Mod(gd, gCurvePrime)

	b := sumRadix(x1, y1).square().strongReduce()
	c := squareRadix(x1).strongReduce()
	d := squareRadix(y1).strongReduce()

	//compareNumbers("b", b, gb)
	//compareNumbers("c", c, gc)
	//compareNumbers("d", d, gd)

	e := sumRadix(c, d).strongReduce()
	h := squareRadix(z1).strongReduce()
	j := subRadix(e, sumRadix(h, h)).strongReduce() //XXX Is there an optimum double?

	//ge := new(big.Int).Add(gc, gd)
	//gh := new(big.Int).Mul(gz, gz)
	//gh = gh.Mod(gh, gCurvePrime)
	//gj := new(big.Int).Add(gh, gh)
	//gj = gj.Sub(ge, gj)

	//compareNumbers("e", e, ge)
	//compareNumbers("h", h, gh)
	//compareNumbers("j", j, gj)

	bMe := subRadix(b, e)
	xx := karatsubaMul(bMe, j) // a = 1 => F = E + D = C + D
	yy := karatsubaMul(e, subRadix(c, d))
	zz := karatsubaMul(e, j).strongReduce()

	//gbMe := new(big.Int).Sub(gb, ge)
	//gxx := new(big.Int).Mul(gbMe, gj)
	//gxx = gxx.Mod(gxx, gCurvePrime)

	//gyy := new(big.Int).Sub(gc, gd)
	//gyy = gyy.Mul(gyy, ge)
	//gyy = gyy.Mod(gyy, gCurvePrime)

	//gzz := new(big.Int).Mul(ge, gj)
	//gzz = gzz.Mod(gzz, gCurvePrime)

	//compareNumbers("b-e", bMe, gbMe)
	//compareNumbers("x3", xx, gxx)
	//compareNumbers("y3", yy, gyy)
	//compareNumbers("z3", zz, gzz)

	//fmt.Printf("b = %s\n", b)
	//fmt.Printf("e = %s\n", e)

	//fmt.Printf("xx = %s\n", xx)
	//fmt.Printf("gxx = 0x%s\n", gxx.Text(16))

	//fmt.Printf("yy = %s\n", yy)
	//fmt.Printf("gyy = 0x%s\n", gyy.Text(16))

	//fmt.Printf("zz = %s\n", zz)
	//fmt.Printf("gzz = 0x%s\n", gzz.Text(16))

	//XXX Should it change the same instance instead?
	return &homogeneousProjective{
		xx, yy, zz,
	}
}

// See Hisil, formula 5.3
func (hP *homogeneousProjective) Add(p Point) Point {
	//A ← Z1*Z2,
	//B ← A^2,
	//C ← X1*X2,
	//D ← Y1*Y2,
	//E ← dC*D,
	//F ← B−E,
	//G ← B+E,
	//X3 ← A*F*((X1+Y1)*(X2+Y2)−C−D),
	//Y3 ← A*G*(D−aC),
	//Z3 ← F*G.

	//gx, _ := new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
	//gy := big.NewInt(0x13)
	//gz := big.NewInt(1)

	x1 := hP[0]
	y1 := hP[1]
	z1 := hP[2]

	hP2 := p.(*homogeneousProjective)
	x2 := hP2[0]
	y2 := hP2[1]
	z2 := hP2[2]

	//compareNumbers("x1", x1, gx)
	//compareNumbers("y1", y1, gy)
	//compareNumbers("z1", z1, gz)

	//compareNumbers("x2", x2, gx)
	//compareNumbers("y2", y2, gy)
	//compareNumbers("z2", z2, gz)

	a := karatsubaMul(z1, z2)
	b := karatsubaMul(a, a)
	c := karatsubaMul(x1, x2)
	d := karatsubaMul(y1, y2)

	tmp := &bigNumber{}
	tmp.mulWSignedCurveConstant(c, curveDSigned)
	tmp.weakReduce()

	e := karatsubaMul(tmp, d)
	f := subRadix(b, e).strongReduce()
	g := sumRadix(b, e).strongReduce()

	//ga := new(big.Int).Mul(gz, gz)
	//ga = ga.Mod(ga, gCurvePrime)

	//gb := new(big.Int).Mul(ga, ga)
	//gb = gb.Mod(gb, gCurvePrime)

	//gc := new(big.Int).Mul(gx, gx)
	//gc = gc.Mod(gc, gCurvePrime)

	//gd := new(big.Int).Mul(gy, gy)
	//gd = gd.Mod(gd, gCurvePrime)

	//ge := new(big.Int).Mul(gCurveD, gc)
	//ge = ge.Mod(ge, gCurvePrime)
	//ge = ge.Mul(ge, gd)
	//ge = ge.Mod(ge, gCurvePrime)

	//DOES NOT MATCH
	//gf := new(big.Int).Sub(gb, ge)
	//gg := new(big.Int).Add(gb, ge)

	//compareNumbers("a", a, ga)
	//compareNumbers("b", b, gb)
	//compareNumbers("c", c, gc)
	//compareNumbers("d", d, gd)
	//compareNumbers("e", e, ge)
	//compareNumbers("f", f, gf)
	//compareNumbers("g", g, gg)

	x3 := karatsubaMul(sumRadix(x1, y1), sumRadix(x2, y2))
	x3 = subRadix(x3, c)
	x3 = subRadix(x3, d)
	x3 = karatsubaMul(a, karatsubaMul(f, x3))

	//gx3 := new(big.Int).Add(gx, gy)
	//gx3 = gx3.Mul(gx3, gx3)
	//gx3 = gx3.Mod(gx3, gCurvePrime)
	//gx3 = gx3.Sub(gx3, gc)
	//gx3 = gx3.Sub(gx3, gd)
	//gx3 = gx3.Mul(gx3, ga)
	//gx3 = gx3.Mod(gx3, gCurvePrime)
	//gx3 = gx3.Mul(gx3, gf)
	//gx3 = gx3.Mod(gx3, gCurvePrime)

	y3 := karatsubaMul(a, karatsubaMul(g, subRadix(d, c)))

	//gy3 := new(big.Int).Sub(gd, gc)
	//gy3 = gy3.Mul(gy3, ga)
	//gy3 = gy3.Mod(gy3, gCurvePrime)
	//gy3 = gy3.Mul(gy3, gg)
	//gy3 = gy3.Mod(gy3, gCurvePrime)

	z3 := karatsubaMul(f, g)

	//gz3 := new(big.Int).Mul(gf, gg)
	//gz3 = gz3.Mod(gz3, gCurvePrime)

	//compareNumbers("x3", x3, gx3)
	//compareNumbers("y3", y3, gy3)

	//fmt.Printf("x3 = %s\n", x3)
	//fmt.Printf("gx3 = 0x%s\n", gx3.Text(16))

	//fmt.Printf("y3 = %s\n", y3)
	//fmt.Printf("gy3 = 0x%s\n", gy3.Text(16))

	//fmt.Printf("z3 = %s\n", z3)
	//fmt.Printf("gz3 = 0x%s\n", gz3.Text(16))

	return &homogeneousProjective{
		x3, y3, z3,
	}
}
