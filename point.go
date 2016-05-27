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
	x := aP[0]
	y := aP[1]

	x2 := new(bigNumber).mul(x, x)
	y2 := new(bigNumber).mul(y, y)

	x2y2 := new(bigNumber).mul(x2, y2)
	dx2y2 := x2y2.mulWSignedCurveConstant(x2y2, curveDSigned)
	dx2y2.weakReduce()

	r := new(bigNumber).add(x2, y2)
	r.sub(r, bigNumOne)
	r.sub(r, dx2y2)

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

	x2 := new(bigNumber).mul(x, x)
	y2 := new(bigNumber).mul(y, y)
	z2 := new(bigNumber).mul(z, z)
	z4 := new(bigNumber).mul(z2, z2)

	x2y2 := new(bigNumber).mul(x2, y2)
	dx2y2 := x2y2.mulWSignedCurveConstant(x2y2, curveDSigned)
	dx2y2.weakReduce()

	r := new(bigNumber).add(x2, y2)
	r.mul(r, z2)
	r.sub(r, z4)
	r.sub(r, dx2y2)

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

	b := new(bigNumber).add(x1, y1)
	b = b.square(b).strongReduce()
	c := new(bigNumber).square(x1).strongReduce()
	d := new(bigNumber).square(y1).strongReduce()

	e := new(bigNumber).add(c, d).strongReduce()
	h := new(bigNumber).square(z1).strongReduce()
	j := new(bigNumber).add(h, h) //XXX Is there an optimum double?
	j.sub(e, j).strongReduce()

	xx := new(bigNumber).sub(b, e)
	xx.mul(xx, j) // a = 1 => F = E + D = C + D
	yy := new(bigNumber).sub(c, d)
	yy.mul(yy, e)
	zz := new(bigNumber).mul(e, j).strongReduce()

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

	a := new(bigNumber).mul(z1, z2)
	b := new(bigNumber).mul(a, a)
	c := new(bigNumber).mul(x1, x2)
	d := new(bigNumber).mul(y1, y2)

	tmp := &bigNumber{}
	tmp.mulWSignedCurveConstant(c, curveDSigned)
	tmp.weakReduce()

	e := new(bigNumber).mul(tmp, d)
	f := new(bigNumber).sub(b, e).strongReduce()
	g := new(bigNumber).add(b, e).strongReduce()

	x3 := new(bigNumber).mul(new(bigNumber).add(x1, y1), new(bigNumber).add(x2, y2))
	x3.sub(x3, c)
	x3.sub(x3, d)
	x3.mul(a, x3.mul(x3, f))

	y3 := new(bigNumber).mul(a, new(bigNumber).mul(g, new(bigNumber).sub(d, c)))

	z3 := new(bigNumber).mul(f, g)

	return &homogeneousProjective{
		x3, y3, z3,
	}
}
