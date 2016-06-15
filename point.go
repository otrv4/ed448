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
	sqrtDminus1         = mustDeserialize(serialized{
		0xd2, 0xe2, 0x18, 0x36, 0x74, 0x9f, 0x46,
		0x88, 0x8d, 0xb4, 0x2b, 0x4f, 0x01, 0x79,
		0x5a, 0x18, 0x9a, 0xab, 0xde, 0xea, 0x38,
		0x51, 0xe6, 0x5c, 0xa6, 0xf1, 0x4c, 0x06,
		0xa4, 0x9f, 0x7b, 0x42, 0x4d, 0x97, 0x70,
		0xdc, 0xac, 0x46, 0x28, 0xc5, 0xf6, 0x56,
		0x49, 0x44, 0x3b, 0x87, 0x48, 0x73, 0x4a,
		0x12, 0xfe, 0xc0, 0xc0, 0xb2, 0x5b, 0x7a,
	})
)

func maskToBoolean(m uint32) bool {
	return m == 0xffffffff
}

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

type extensibleCoordinates struct {
	x, y, z, t, u *bigNumber
}

//Affina(x,y) => extensible(X, Y, Z, T, U)
func newExtensible(px, py *bigNumber) *extensibleCoordinates {
	x := px.copy()
	y := py.copy()
	z := &bigNumber{1}
	t := x.copy()
	u := y.copy()

	return &extensibleCoordinates{
		x: x,
		y: y,
		z: z,
		t: t,
		u: u,
	}
}

func (p *extensibleCoordinates) twist() *twExtensible {
	x := new(bigNumber)
	y := new(bigNumber)
	z := new(bigNumber)
	t := new(bigNumber)
	u := new(bigNumber)

	l0 := new(bigNumber)
	l1 := new(bigNumber)

	u = u.square(p.z)
	y = y.square(p.x)
	z = z.sub(u, y)
	y = y.add(z, z)
	u = u.add(y, y)
	y = y.sub(p.z, p.x)
	x = x.mul(y, p.y)
	z = z.sub(p.z, p.y)
	t = t.mul(z, x)
	l1 = l1.mul(t, u)

	x = x.mul(t, l1)
	l0 = l0.isr(x)
	u = u.mul(t, l0)
	l1 = l1.square(l0)
	t = t.mul(x, l1)
	l1 = l1.add(p.x, p.y)
	l0 = l0.sub(p.x, p.y)
	x = x.mul(t, l0)
	l0 = l0.add(x, l1)
	t = t.sub(l1, x)
	x = x.mul(l0, u)
	x = x.addW(-y.zeroMask())
	y = y.mul(t, u)
	y = y.addW(-z.zeroMask())
	z = z.setUi(1 + uint64(y.zeroMask()))
	t = x.copy()
	u = y.copy()

	return &twExtensible{x, y, z, t, u}
}

func (p *extensibleCoordinates) Double() *extensibleCoordinates {
	x := p.x.copy()
	y := p.y.copy()
	z := p.z.copy()
	t := p.t.copy()
	u := p.u.copy()

	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2 = l2.square(x)
	l0 = l0.square(y)
	l1 = l1.addRaw(l2, l0)
	t = t.addRaw(y, x)
	u = u.square(t)
	t = t.subRaw(u, l1).bias(3).weakReduce()
	u = u.sub(l0, l2) // equivalent to subx in 32-bits
	x = x.square(z).bias(2)
	z = z.addRaw(x, x)
	l0 = l0.subRaw(z, l1).weakReduce()
	z = z.mul(l1, l0)
	x = x.mul(l0, t)
	y = y.mul(l1, u)

	return &extensibleCoordinates{
		x: x,
		y: y,
		z: z,
		t: t,
		u: u,
	}
}

func (p *extensibleCoordinates) OnCurve() bool {
	x := p.x
	y := p.y
	z := p.z
	t := p.t
	u := p.u

	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)

	// Check invariant:
	// 0 = d*t^2*u^2 - x^2 - y^2 + z^2
	l2 = l2.square(y)
	l1 = l1.neg(l2)
	l0 = l0.square(z)
	l2 = l2.add(l0, l1)
	l3 = l3.square(u)
	l0 = l0.square(t)
	l1 = l1.mul(l0, l3)
	l0 = l0.mulWSignedCurveConstant(l1, curveDSigned)
	l1 = l1.add(l0, l2)
	l0 = l0.square(x)
	l2 = l2.neg(l0)
	l0 = l0.add(l2, l1)
	l5 := l0.zeroMask()

	// Check invariant:
	// 0 = -x*y + z*t*u
	l1 = l1.mul(t, u)
	l2 = l2.mul(z, l1)
	l0 = l0.mul(x, y)
	l1 = l1.neg(l0)
	l0 = l0.add(l1, l2)

	l4 := l0.zeroMask()

	ret := l4 & l5 & (^z.zeroMask())
	return maskToBoolean(ret)
}

func (p *extensibleCoordinates) equals(q *extensibleCoordinates) bool {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2 = l2.mul(q.z, p.x)
	l1 = l1.mul(p.z, q.x)
	l0 = l0.sub(l2, l1)
	l4 := l0.zeroMask()

	l2 = l2.mul(q.z, p.y)
	l1 = l1.mul(p.z, q.y)
	l0 = l0.sub(l2, l1)
	l3 := l0.zeroMask()

	return maskToBoolean(l4 & l3)
}

type twPNiels struct {
	n *twNiels
	z *bigNumber
}

func newTwistedPNiels(a, b, c, z [56]byte) *twPNiels {
	return &twPNiels{
		&twNiels{
			a: mustDeserialize(serialized(a)),
			b: mustDeserialize(serialized(b)),
			c: mustDeserialize(serialized(c)),
		},
		mustDeserialize(serialized(z)),
	}
}

func (p *twPNiels) String() string {
	return fmt.Sprintf("A: %s\nB: %s\nC: %s\nZ: %s\n", p.n.a, p.n.b, p.n.c, p.z)
}

func (p *twPNiels) equals(p2 *twPNiels) bool {
	ok := true

	ok = ok && p.n.equals(p2.n)
	ok = ok && p.z.equals(p2.z)

	return ok
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

func (p *twNiels) equals(p2 *twNiels) bool {
	ok := true

	ok = ok && p.a.equals(p2.a)
	ok = ok && p.b.equals(p2.b)
	ok = ok && p.c.equals(p2.c)

	return ok
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

type twExtensible struct {
	x, y, z, t, u *bigNumber
}

func (p *twExtensible) Add(p1 Point) Point {
	p.addTwPNiels(p1.(*twExtensible).twPNiels())
	return p
}

func (p *twExtensible) addTwPNiels(a *twPNiels) *twExtensible {
	// field_mul ( L0, e->z, a->z );
	L0 := new(bigNumber).mul(p.z, a.z)
	// field_copy ( e->z, L0 );
	p.z = L0.copy()
	// add_tw_niels_to_tw_extensible( e, a->n );
	return p.addTwNiels(a.n)
}

func (p *twExtensible) Double() Point {
	p = p.double()
	return p
}

func (p *twExtensible) Marshal() []byte {
	return nil
}

func (a *twExtensible) twPNiels() *twPNiels {
	// field_sub ( b->n->a, a->y, a->x );
	na := new(bigNumber).sub(a.y, a.x)
	// field_add ( b->n->b, a->x, a->y );
	nb := new(bigNumber).add(a.x, a.y)
	// field_mul ( b->z, a->u, a->t );
	z := new(bigNumber).mul(a.u, a.t)
	// field_mulw_scc_wr ( b->n->c, b->z, 2*EDWARDS_D-2 );
	nc := new(bigNumber).mulWSignedCurveConstant(z, curveDSigned*2-2)
	// field_add ( b->z, a->z, a->z );
	z.add(a.z, a.z)
	return &twPNiels{
		n: &twNiels{
			a: na,
			b: nb,
			c: nc,
		},
		z: z,
	}
}

func (p *twExtensible) OnCurve() bool {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)

	// Check invariant:
	// 0 = -x*y + z*t*u
	l1 = l1.mul(p.t, p.u)
	l2 = l2.mul(p.z, l1)
	l0 = l0.mul(p.x, p.y)
	l1 = l1.neg(l0)
	l0 = l0.add(l1, l2)
	l5 := l0.zeroMask()

	// Check invariant:
	// 0 = d*t^2*u^2 + x^2 - y^2 + z^2 - t^2*u^2

	l2 = l2.square(p.y)
	l1 = l1.neg(l2)
	l0 = l0.square(p.x)
	l2 = l2.add(l0, l1)
	l3 = l3.square(p.u)
	l0 = l0.square(p.t)
	l1 = l1.mul(l0, l3)
	l3 = l3.mulWSignedCurveConstant(l1, curveDSigned)
	l0 = l0.add(l3, l2)
	l3 = l3.neg(l1)
	l2 = l2.add(l3, l0)
	l1 = l1.square(p.z)
	l0 = l0.add(l1, l2)
	l4 := l0.zeroMask()

	ret := l4 & l5 & (^p.z.zeroMask())
	return maskToBoolean(ret)
}

func (p *twExtensible) String() string {
	ret := fmt.Sprintf("X: %s\n", p.x)
	ret += fmt.Sprintf("Y: %s\n", p.y)
	ret += fmt.Sprintf("Z: %s\n", p.z)
	ret += fmt.Sprintf("T: %s\n", p.t)
	ret += fmt.Sprintf("U: %s\n", p.u)

	return ret
}

func (p *twExtensible) equals(p2 *twExtensible) bool {
	ok := true

	ok = ok && p.x.equals(p2.x)
	ok = ok && p.y.equals(p2.y)
	ok = ok && p.z.equals(p2.z)
	ok = ok && p.t.equals(p2.t)
	ok = ok && p.u.equals(p2.u)

	return ok
}

func (p *twExtensible) double() *twExtensible {
	x := p.x.copy()
	y := p.y.copy()
	z := p.z.copy()
	t := p.t.copy()
	u := p.u.copy()

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
	x := p.x.copy()
	y := p.y.copy()
	z := p.z.copy()
	t := p.t.copy()
	u := p.u.copy()

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

func (p *twExtensible) untwistAndDoubleAndSerialize() *bigNumber {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)
	b := new(bigNumber)

	l3 = l3.mul(p.y, p.x)
	b = b.add(p.y, p.x)
	l1 = l1.square(b)
	l2 = l2.add(l3, l3)
	b = b.sub(l1, l2)
	l2 = l2.square(p.z)
	l1 = l1.square(l2)
	b = b.add(b, b)
	l2 = l2.mulWSignedCurveConstant(b, curveDSigned-1)
	b = b.mulWSignedCurveConstant(l2, curveDSigned-1)
	l0 = l0.mul(l2, l1)
	l2 = l2.mul(b, l0)
	l0 = l0.isr(l2)
	l1 = l1.mul(b, l0)

	//XXX This is included in the original code, but it seems not to be used
	//b = b.square(l0)
	//l0 = l0.mul(l2, b)

	return b.mul(l1, l3)
}

//HP(X : Y : Z) = Affine(X/Z, Y/Z), Z ≠ 0
//XXX This can be replaced by extensible for simplicity
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

	//XXX PERF Should it change the same instance instead?
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
