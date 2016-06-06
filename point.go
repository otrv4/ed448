package ed448

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
)

var (
	bigNumOne           = mustDeserialize(serialized{1})
	bigNumTwo           = mustDeserialize(serialized{2})
	curveDSigned        = int64(-39081)
	twistedCurveDSigned = int64(-39082)
)

// Point represents a point on the curve in a suitable coordinate system
type Point interface {
	OnCurve() bool
	Add(Point) Point
	Double() Point

	Marshal() []byte
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

type twNiels struct {
	a, b, c *bigNumber
}

func newNielsPoint(a, b, c [56]byte) *twNiels {
	return &twNiels{
		a: mustDeserialize(serialized(a)),
		b: mustDeserialize(serialized(b)),
		c: mustDeserialize(serialized(c)),
	}
}

func (p *twNiels) String() string {
	return fmt.Sprintf("A: %s\nB: %s\nC: %s\n", p.a, p.b, p.c)
}

func (p *twNiels) copy() *twNiels {
	return &twNiels{
		a: p.a.copy(),
		b: p.b.copy(),
		c: p.c.copy(),
	}
}

//XXX SECURITY this should be constant-time
func (nP *twNiels) conditionalNegate(neg bool) {
	if neg {
		tmp := nP.a
		nP.a = nP.b
		nP.b = tmp
		nP.c.neg(nP.c)
	}
}

func (p *twNiels) TwistedExtensible() *twExtensible {
	x := new(bigNumber)
	y := new(bigNumber)
	z := new(bigNumber)
	t := new(bigNumber)
	u := new(bigNumber)

	y = y.add(p.b, p.a)
	x = x.sub(p.b, p.a)
	z = &bigNumber{1}
	t = x.copy()
	u = y.copy()

	//PERF: should it be in-place?
	return &twExtensible{x, y, z, t, u}
}

type twExtensible [5]*bigNumber

func (p *twExtensible) Add(Point) Point {
	return nil
}

func (p *twExtensible) Double() Point {
	return nil
}

func (p *twExtensible) Marshal() []byte {
	return nil
}

func (p *twExtensible) OnCurve() bool {
	x := p[0]
	y := p[1]
	z := p[2]
	t := p[3]
	u := p[4]

	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)

	// Check invariant:
	// 0 = -x*y + z*t*u
	l1 = l1.mul(t, u)
	l2 = l2.mul(z, l1)
	l0 = l0.mul(x, y)
	l1 = l1.neg(l0)
	l0 = l0.add(l1, l2)
	l5 := l0.zero()

	// Check invariant:
	// 0 = d*t^2*u^2 + x^2 - y^2 + z^2 - t^2*u^2

	l2 = l2.square(y)
	l1 = l1.neg(l2)
	l0 = l0.square(x)
	l2 = l2.add(l0, l1)
	l3 = l3.square(u)
	l0 = l0.square(t)
	l1 = l1.mul(l0, l3)
	l3 = l3.mulWSignedCurveConstant(l1, curveDSigned)
	l0 = l0.add(l3, l2)
	l3 = l3.neg(l1)
	l2 = l2.add(l3, l0)
	l1 = l1.square(z)
	l0 = l0.add(l1, l2)
	l4 := l0.zero()

	//XXX SECURITY this might not be constant time (due logical short circuit)
	//zero() should return an mask
	return l4 && l5 && !z.zero()
}

func (p *twExtensible) String() string {
	x := p[0]
	y := p[1]
	z := p[2]
	t := p[3]
	u := p[4]

	ret := fmt.Sprintf("X: %s\n", x)
	ret += fmt.Sprintf("Y: %s\n", y)
	ret += fmt.Sprintf("Z: %s\n", z)
	ret += fmt.Sprintf("T: %s\n", t)
	ret += fmt.Sprintf("U: %s\n", u)

	return ret
}

func (p *twExtensible) equals(p2 *twExtensible) bool {
	ok := true

	for i, pi := range p {
		ok = ok && pi.equals(p2[i])
	}

	return ok
}

func (p *twExtensible) double() *twExtensible {
	x := p[0].copy()
	y := p[1].copy()
	z := p[2].copy()
	t := p[3].copy()
	u := p[4].copy()

	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2 = l2.square(x)
	l0 = l0.square(y)
	u = u.addRaw(l2, l0)
	t = t.addRaw(y, x)
	l1 = l1.square(t)
	t = t.subRaw(l1, u)
	t.bias(3)
	t.weakReduce()
	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	l1 = l1.sub(l0, l2)
	x = x.square(z)
	x.bias(1)
	z = z.addRaw(x, x)
	l0 = l0.subRaw(z, l1)
	l0.weakReduce()
	z = z.mul(l1, l0)
	x = x.mul(l0, t)
	y = y.mul(l1, u)

	//PERF: should it be in-place?
	return &twExtensible{x, y, z, t, u}
}

func (p *twExtensible) addTwNiels(p2 *twNiels) *twExtensible {
	x := p[0].copy()
	y := p[1].copy()
	z := p[2].copy()
	t := p[3].copy()
	u := p[4].copy()

	l0 := new(bigNumber)
	l1 := new(bigNumber)

	l1 = l1.sub(y, x)
	l0 = l0.mul(p2.a, l1)
	l1 = l1.addRaw(x, y)
	y = y.mul(p2.b, l1)
	l1 = l1.mul(u, t)
	x = x.mul(p2.c, l1)

	u = u.addRaw(l0, y)
	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	t = t.sub(y, l0)

	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	y = y.sub(z, x)
	l0 = l0.addRaw(x, z)

	z = z.mul(l0, y)
	x = x.mul(y, t)
	y = y.mul(l0, u)

	//PERF: should it be in-place?
	return &twExtensible{x, y, z, t, u}
}

type twistedHomogeneousProjective [3]*bigNumber

func NewTwistedPoint(x serialized, y serialized) (p Point, e error) {
	xN, ok1 := deserialize(x)
	yN, ok2 := deserialize(y)

	p = newTwistedHomogeneousProjective(xN, yN)

	if !(ok1 && ok2) {
		e = errors.New("invalid coordinates")
	}

	return
}

//Affine to Twisted Homogeneous Projective
func newTwistedHomogeneousProjective(x *bigNumber, y *bigNumber) *twistedHomogeneousProjective {
	x1 := new(bigNumber).mul(x, y)
	x1 = x1.mul(x1, bigNumTwo)

	x2 := new(bigNumber).mul(x, x)
	y2 := new(bigNumber).mul(y, y)
	x2plusy2 := new(bigNumber).add(x2, y2)

	y1 := x2plusy2

	z1 := new(bigNumber).sub(y2, x2)
	z2 := new(bigNumber).sub(bigNumTwo, x2plusy2)

	return &twistedHomogeneousProjective{
		x1.mul(x1, z2).copy(), // X * Z
		y1.mul(y1, z1).copy(), // Y * Z
		z1.mul(z1, z2).copy(), // Z = 1
	}
}

func (hP *twistedHomogeneousProjective) OnCurve() bool {
	// (-x² + y²)z² - z^4 - (d-1)x²y² = 0
	x := hP[0]
	y := hP[1]
	z := hP[2]

	x2 := new(bigNumber).mul(x, x)
	y2 := new(bigNumber).mul(y, y)
	z2 := new(bigNumber).mul(z, z)
	z4 := new(bigNumber).mul(z2, z2)

	x2y2 := new(bigNumber).mul(x2, y2)
	dx2y2 := x2y2.mulWSignedCurveConstant(x2y2, twistedCurveDSigned)
	dx2y2.weakReduce()

	r := new(bigNumber).sub(y2, x2)
	r.mul(r, z2)
	r.sub(r, z4)
	r.sub(r, dx2y2)

	r.strongReduce()
	return r.zero()
}

func (hP *twistedHomogeneousProjective) Add(p Point) Point {
	//a = -1
	//d = -39082
	//A ← Z1*Z2,
	//B ← A^2,
	//C ← X1*X2,
	//D ← Y1*Y2,
	//E ← dC*D,
	//F ← B−E,
	//G ← B+E,
	//U ← C+D,
	//X3 ← A*F*((X1+Y1)*(X2+Y2)−U),
	//Y3 ← A*G*U,
	//Z3 ← F*G.

	x1 := hP[0]
	y1 := hP[1]
	z1 := hP[2]

	hP2 := p.(*twistedHomogeneousProjective)
	x2 := hP2[0]
	y2 := hP2[1]
	z2 := hP2[2]

	a := new(bigNumber).mul(z1, z2)
	b := new(bigNumber).mul(a, a)
	c := new(bigNumber).mul(x1, x2)
	d := new(bigNumber).mul(y1, y2)

	e := new(bigNumber).mulWSignedCurveConstant(c, twistedCurveDSigned)
	e.mul(e, d)
	f := new(bigNumber).sub(b, e)
	g := new(bigNumber).add(b, e)
	u := new(bigNumber).add(c, d)

	//Just reusing e and b (unused) memory
	x3 := e.mul(b.add(x1, y1), e.add(x2, y2))
	x3.sub(x3, u)
	x3.mul(x3, a).mul(x3, f)

	//reuse u
	y3 := u.mul(u, a)
	y3 = y3.mul(y3, g)

	z3 := f.mul(f, g)

	return &twistedHomogeneousProjective{
		x3, y3, z3,
	}
}

func (hP *twistedHomogeneousProjective) Double() Point {
	return nil
}

func (hP *twistedHomogeneousProjective) Marshal() []byte {
	return nil
}

func (hP *twistedHomogeneousProjective) conditionalNegate(neg bool) {
	if neg {
		hP[0].neg(hP[0])
	}
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
	return fmt.Sprintf("X: %s\nY: %s\nZ: %s\n", hP[0], hP[1], hP[2])
}

func (hP *homogeneousProjective) conditionalNegate(neg bool) {
	//XXX this should be constant-time
	if neg {
		hP[0].neg(hP[0])
	}
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
	b.square(b)
	c := new(bigNumber).square(x1)
	d := new(bigNumber).square(y1)
	e := new(bigNumber).add(c, d)
	h := new(bigNumber).square(z1)
	//j := h.mulW(h, 2) // This is slower than adding
	j := h.add(h, h)
	j.sub(e, j)

	xx := b.sub(b, e)
	xx.mul(xx, j)
	yy := c.sub(c, d)
	yy.mul(yy, e)
	zz := e.mul(e, j)

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

	e := new(bigNumber).mulWSignedCurveConstant(c, curveDSigned)
	e.mul(e, d)
	f := new(bigNumber).sub(b, e)
	g := new(bigNumber).add(b, e)

	//Just reusing e and b (unused) memory
	x3 := e.mul(b.add(x1, y1), e.add(x2, y2))
	x3.sub(x3, c).sub(x3, d)
	x3.mul(x3, a).mul(x3, f)

	y3 := d.sub(d, c)
	y3 = y3.mul(y3, a).mul(y3, g)

	z3 := f.mul(f, g)

	return &homogeneousProjective{
		x3, y3, z3,
	}
}

func (hP *homogeneousProjective) Marshal() []byte {
	byteLen := 56

	dst := make([]byte, byteLen)
	serialize(dst, hP[0]) //x little endian
	x := new(big.Int).SetBytes(rev(dst))

	serialize(dst, hP[1]) //y little endian
	y := new(big.Int).SetBytes(rev(dst))

	serialize(dst, hP[2]) //z little endian
	z := new(big.Int).SetBytes(rev(dst))

	//x and y in affine coordinates
	//XXX I'm not sure if I need to covert to affine
	x.Div(x, z)
	y.Div(y, z)

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point

	xBytes := x.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)

	yBytes := y.Bytes()
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)
	return ret
}
