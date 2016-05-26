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

	b := sumRadix(x1, y1).square().strongReduce()
	c := squareRadix(x1).strongReduce()
	d := squareRadix(y1).strongReduce()

	e := sumRadix(c, d).strongReduce()
	h := squareRadix(z1).strongReduce()
	j := subRadix(e, sumRadix(h, h)).strongReduce() //XXX Is there an optimum double?

	bMe := subRadix(b, e)
	xx := karatsubaMul(bMe, j) // a = 1 => F = E + D = C + D
	yy := karatsubaMul(e, subRadix(c, d))
	zz := karatsubaMul(e, j).strongReduce()

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

	x1 := hP[0]
	y1 := hP[1]
	z1 := hP[2]

	hP2 := p.(*homogeneousProjective)
	x2 := hP2[0]
	y2 := hP2[1]
	z2 := hP2[2]

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

	x3 := karatsubaMul(sumRadix(x1, y1), sumRadix(x2, y2))
	x3 = subRadix(x3, c)
	x3 = subRadix(x3, d)
	x3 = karatsubaMul(a, karatsubaMul(f, x3))

	y3 := karatsubaMul(a, karatsubaMul(g, subRadix(d, c)))

	z3 := karatsubaMul(f, g)

	return &homogeneousProjective{
		x3, y3, z3,
	}
}
