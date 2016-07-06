package ed448

import (
	"errors"
	"fmt"
)

var (
	bigNumOne           = mustDeserialize(serialized{1})
	bigNumTwo           = mustDeserialize(serialized{2})
	curveDSigned        = int64(-39081)
	twistedCurveDSigned = int64(-39082)
	sqrtDminus1         = mustDeserialize(serialized{
		0x46, 0x9f, 0x74, 0x36, 0x18, 0xe2, 0xd2, 0x79,
		0x01, 0x4f, 0x2b, 0xb4, 0x8d, 0x88, 0x38, 0xea,
		0xde, 0xab, 0x9a, 0x18, 0x5a, 0x06, 0x4c, 0xf1,
		0xa6, 0x5c, 0xe6, 0x51, 0x70, 0x97, 0x4d, 0x42,
		0x7b, 0x9f, 0xa4, 0x56, 0xf6, 0xc5, 0x28, 0x46,
		0xac, 0xdc, 0x4a, 0x73, 0x48, 0x87, 0x3b, 0x44,
		0x49, 0x7a, 0x5b, 0xb2, 0xc0, 0xc0, 0xfe, 0x12,
	})
)

func maskToBoolean(m uint32) bool {
	return m == 0xffffffff
}

// NewPoint instantiates a new point in a suitable coordinate system.
// The x and y coordinates must be affine coordinates in little-endian
//XXX This should probably receive []byte{}
func NewPoint(x serialized, y serialized) (p *homogeneousProjective, e error) {
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
	l2 = l2.squareCopy(y)
	l1 = l1.neg(l2)
	l0 = l0.squareCopy(z)
	l2 = l2.add(l0, l1)
	l3 = l3.squareCopy(u)
	l0 = l0.squareCopy(t)
	l1 = l1.mulCopy(l0, l3)
	l0 = l0.mulWSignedCurveConstant(l1, curveDSigned)
	l1 = l1.add(l0, l2)
	l0 = l0.squareCopy(x)
	l2 = l2.neg(l0)
	l0 = l0.add(l2, l1)
	l5 := l0.zeroMask()

	// Check invariant:
	// 0 = -x*y + z*t*u
	l1 = l1.mulCopy(t, u)
	l2 = l2.mulCopy(z, l1)
	l0 = l0.mulCopy(x, y)
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

	l2 = l2.mulCopy(q.z, p.x)
	l1 = l1.mulCopy(p.z, q.x)
	l0 = l0.sub(l2, l1)
	l4 := l0.zeroMask()

	l2 = l2.mulCopy(q.z, p.y)
	l1 = l1.mulCopy(p.z, q.y)
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

func (nP *twNiels) conditionalNegate(neg word_t) {
	nP.a.conditionalSwap(nP.b, neg)
	nP.c = nP.c.conditionalNegate(neg)
}

func convertTwNielsToTwExtensible(dst *twExtensible, src *twNiels) {
	dst.y = dst.y.add(src.b, src.a)
	dst.x = dst.x.sub(src.b, src.a)
	dst.z = dst.z.setUi(1)
	dst.t = dst.x.copy()
	dst.u = dst.y.copy()
}

type twExtensible struct {
	x, y, z, t, u *bigNumber
}

func (p *twExtensible) copy(e *twExtensible) *twExtensible {
	p.x = e.x.copy()
	p.y = e.y.copy()
	p.z = e.z.copy()
	p.t = e.t.copy()
	p.u = e.u.copy()

	return p
}

func (p *twExtensible) addTwPNiels(a *twPNiels) *twExtensible {
	p.z = p.z.mulCopy(p.z, a.z)
	return p.addTwNiels(a.n)
}

func (e *twExtensible) subTwPNiels(a *twPNiels) {
	e.z = e.z.mulCopy(e.z, a.z)
	e.subTwNiels(a.n)
}

func convertTwExtensibleToTwPNiels(dst *twPNiels, src *twExtensible) {
	dst.n.a.sub(src.y, src.x)
	dst.n.b.add(src.x, src.y)
	karatsubaMul(dst.z, src.u, src.t)
	dst.n.c.mulWSignedCurveConstant(dst.z, curveDSigned*2-2)
	dst.z.add(src.z, src.z)
}

func (a *twExtensible) twPNiels() *twPNiels {
	ret := &twPNiels{
		n: &twNiels{
			a: new(bigNumber),
			b: new(bigNumber),
			c: new(bigNumber),
		},
		z: new(bigNumber),
	}

	convertTwExtensibleToTwPNiels(ret, a)
	return ret
}

func convertTwPnielsToTwExtensible(dst *twExtensible, src *twPNiels) {
	dst.u.add(src.n.b, src.n.a)
	dst.t.sub(src.n.b, src.n.a)
	karatsubaMul(dst.x, src.z, dst.t)
	karatsubaMul(dst.y, src.z, dst.u)
	karatsubaSquare(dst.z, src.z)
}

func (p *twExtensible) OnCurve() bool {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)

	// Check invariant:
	// 0 = -x*y + z*t*u
	l1 = l1.mulCopy(p.t, p.u)
	l2 = l2.mulCopy(p.z, l1)
	l0 = l0.mulCopy(p.x, p.y)
	l1 = l1.neg(l0)
	l0 = l0.add(l1, l2)
	l5 := l0.zeroMask()

	// Check invariant:
	// 0 = d*t^2*u^2 + x^2 - y^2 + z^2 - t^2*u^2

	l2 = l2.squareCopy(p.y)
	l1 = l1.neg(l2)
	l0 = l0.squareCopy(p.x)
	l2 = l2.add(l0, l1)
	l3 = l3.squareCopy(p.u)
	l0 = l0.squareCopy(p.t)
	l1 = l1.mulCopy(l0, l3)
	l3 = l3.mulWSignedCurveConstant(l1, curveDSigned)
	l0 = l0.add(l3, l2)
	l3 = l3.neg(l1)
	l2 = l2.add(l3, l0)
	l1 = l1.squareCopy(p.z)
	l0 = l0.add(l1, l2)
	l4 := l0.zeroMask()

	ret := l4 & l5 & (^p.z.zeroMask())
	return maskToBoolean(ret)
}

func (p *twExtensible) setIdentity() {
	p.x.setUi(0)
	p.y.setUi(1)
	p.z.setUi(1)
	p.t.setUi(0)
	p.u.setUi(0)
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
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2 = l2.mulCopy(p2.z, p.x)
	l1 = l1.mulCopy(p.z, p2.x)
	l0 = l0.sub(l2, l1)

	l4 := l0.zeroMask()

	l2 = l2.mulCopy(p2.z, p.y)
	l1 = l1.mulCopy(p.z, p2.y)
	l0 = l0.sub(l2, l1)

	l3 := l0.zeroMask()

	return (l4 & l3) == 0xffffffff
}

func (p *twExtensible) double() *twExtensible {
	x := p.x
	y := p.y
	z := p.z
	t := p.t
	u := p.u

	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	//We use karatsubaSquare and karatsubaMul directly because we know it is safe
	//to use them (and it's faster - it avoids creating one intermediate object)
	karatsubaSquare(l2, x)
	karatsubaSquare(l0, y)
	u = u.addRaw(l2, l0)
	t = t.addRaw(y, x)
	karatsubaSquare(l1, t)
	t = t.subRaw(l1, u)
	t.bias(3)
	t.weakReduce()
	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	l1 = l1.sub(l0, l2)
	karatsubaSquare(x, z)
	x.bias(1)
	z = z.addRaw(x, x)
	l0 = l0.subRaw(z, l1)
	l0.weakReduce()
	karatsubaMul(z, l1, l0)
	karatsubaMul(x, l0, t)
	karatsubaMul(y, l1, u)

	return p
}

func (p *twExtensible) addTwNiels(p2 *twNiels) *twExtensible {
	x := p.x
	y := p.y
	z := p.z
	t := p.t
	u := p.u

	l0 := new(bigNumber)
	l1 := new(bigNumber)

	l1 = l1.sub(y, x)
	karatsubaMul(l0, p2.a, l1)
	l1 = l1.addRaw(x, y)
	karatsubaMul(y, p2.b, l1)
	karatsubaMul(l1, u, t)
	karatsubaMul(x, p2.c, l1)

	u = u.addRaw(l0, y)
	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	t = t.sub(y, l0)

	// This is equivalent do subx_nr in 32 bits. Change if using 64-bits
	y = y.sub(z, x)
	l0 = l0.addRaw(x, z)

	karatsubaMul(z, l0, y)
	karatsubaMul(x, y, t)
	karatsubaMul(y, l0, u)

	return p
}

func (d *twExtensible) subTwNiels(e *twNiels) {
	L1 := new(bigNumber).subxRaw(d.y, d.x)
	L0 := karatsubaMul(new(bigNumber), e.b, L1)
	L1.addRaw(d.x, d.y)
	karatsubaMul(d.y, e.a, L1)
	karatsubaMul(L1, d.u, d.t)
	karatsubaMul(d.x, e.c, L1)
	d.u.addRaw(L0, d.y)
	d.t.subxRaw(d.y, L0)
	d.y.addRaw(d.x, d.z)
	L0.subxRaw(d.z, d.x)
	karatsubaMul(d.z, L0, d.y)
	karatsubaMul(d.x, d.y, d.t)
	karatsubaMul(d.y, L0, d.u)
}

func (p *twExtensible) untwistAndDoubleAndSerialize() *bigNumber {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)
	b := new(bigNumber)

	karatsubaMul(l3, p.y, p.x)
	b.add(p.y, p.x)
	karatsubaSquare(l1, b)
	l2.add(l3, l3)
	b.sub(l1, l2)
	karatsubaSquare(l2, p.z)
	karatsubaSquare(l1, l2)
	b.add(b, b)
	l2.mulWSignedCurveConstant(b, curveDSigned-1)
	b.mulWSignedCurveConstant(l2, curveDSigned-1)
	karatsubaMul(l0, l2, l1)
	karatsubaMul(l2, b, l0)
	l0.isr(l2)
	karatsubaMul(l1, b, l0)

	//XXX This is included in the original code, but it seems not to be used
	//b = b.square(l0)
	//l0 = l0.mul(l2, b)

	return karatsubaMul(b, l1, l3)
}

//HP(X : Y : Z) = Affine(X/Z, Y/Z), Z ≠ 0
//XXX This can be replaced by extensible for simplicity if we neither use ADD
//on the basePoint in test and benchmark (it is not used elsewhere)
type homogeneousProjective struct {
	x, y, z *bigNumber
}

//Affine to Homogeneous Projective
func newHomogeneousProjective(x *bigNumber, y *bigNumber) *homogeneousProjective {
	return &homogeneousProjective{
		x: x.copy(),
		y: y.copy(),
		z: bigNumOne.copy(),
	}
}

func (p *homogeneousProjective) String() string {
	return fmt.Sprintf("X: %s\nY: %s\nZ: %s\n", p.x, p.y, p.z)
}

func (p *homogeneousProjective) OnCurve() bool {
	// (x² + y²)z² - z^4 - dx²y² = 0
	x := p.x
	y := p.y
	z := p.z

	x2 := new(bigNumber).mulCopy(x, x)
	y2 := new(bigNumber).mulCopy(y, y)
	z2 := new(bigNumber).mulCopy(z, z)
	z4 := new(bigNumber).mulCopy(z2, z2)

	x2y2 := new(bigNumber).mulCopy(x2, y2)
	dx2y2 := x2y2.mulWSignedCurveConstant(x2y2, curveDSigned)
	dx2y2.weakReduce()

	r := new(bigNumber).add(x2, y2)
	r.mulCopy(r, z2)
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

// See Hisil, formula 5.1
//XXX Used only for testing
func (p *homogeneousProjective) double() *homogeneousProjective {
	x1 := p.x
	y1 := p.y
	z1 := p.z

	b := new(bigNumber).add(x1, y1)
	b.squareCopy(b)
	c := new(bigNumber).squareCopy(x1)
	d := new(bigNumber).squareCopy(y1)
	e := new(bigNumber).add(c, d)
	h := new(bigNumber).squareCopy(z1)
	//j := h.mulW(h, 2) // This is slower than adding
	j := h.add(h, h)
	j.sub(e, j)

	xx := b.sub(b, e)
	xx.mulCopy(xx, j)
	yy := c.sub(c, d)
	yy.mulCopy(yy, e)
	zz := e.mulCopy(e, j)

	//XXX PERF Should it change the same instance instead?
	return &homogeneousProjective{
		xx, yy, zz,
	}
}

// See Hisil, formula 5.3
func (p *homogeneousProjective) add(p2 *homogeneousProjective) *homogeneousProjective {
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

	x1 := p.x
	y1 := p.y
	z1 := p.z

	x2 := p2.x
	y2 := p2.y
	z2 := p2.z

	a := new(bigNumber).mulCopy(z1, z2)
	b := new(bigNumber).squareCopy(a)
	c := new(bigNumber).mulCopy(x1, x2)
	d := new(bigNumber).mulCopy(y1, y2)

	e := new(bigNumber).mulWSignedCurveConstant(c, curveDSigned)
	e.mulCopy(e, d)
	f := new(bigNumber).sub(b, e)
	g := new(bigNumber).add(b, e)

	//Just reusing e and b (unused) memory
	x3 := e.mulCopy(b.add(x1, y1), e.add(x2, y2))
	x3.sub(x3, c).sub(x3, d)
	x3.mulCopy(x3, a).mulCopy(x3, f)

	y3 := d.sub(d, c)
	y3 = y3.mulCopy(y3, a).mulCopy(y3, g)

	z3 := f.mulCopy(f, g)

	return &homogeneousProjective{
		x3, y3, z3,
	}
}

//XXX Move: bigNumber should not know about points
func (sz *bigNumber) deserializeAndTwistApprox() (*twExtensible, bool) {
	a := &twExtensible{
		x: new(bigNumber),
		y: new(bigNumber),
		z: new(bigNumber),
		u: new(bigNumber),
		t: new(bigNumber),
	}

	var L0, L1 *bigNumber
	L0 = new(bigNumber)
	L1 = new(bigNumber)
	// field_sqr ( a->z, sz );
	a.z.squareCopy(sz)
	// field_copy ( a->y, a->z );
	a.y = a.z.copy()
	// field_addw ( a->y, 1 );
	a.y.addW(1)
	// field_sqr ( L0, a->y );
	L0.squareCopy(a.y)
	// field_mulw_scc ( a->x, L0, EDWARDS_D-1 );
	a.x.mulWSignedCurveConstant(L0, curveDSigned-1)
	// field_add ( a->y, a->z, a->z );
	a.y.add(a.z, a.z)
	// field_add ( a->u, a->y, a->y );
	a.u.add(a.y, a.y)
	// field_add ( a->y, a->u, a->x );
	a.y.add(a.u, a.x)
	// field_sqr ( a->x, a->z );
	a.x.squareCopy(a.z)
	// field_neg ( a->u, a->x );
	a.u.neg(a.x)
	// field_addw ( a->u, 1 );
	a.u.addW(1)
	// field_mul ( a->x, sqrt_d_minus_1, a->u );
	a.x.mulCopy(sqrtDminus1, a.u)
	// field_mul ( L0, a->x, a->y );
	L0.mulCopy(a.x, a.y)
	// field_mul ( a->t, L0, a->y );
	a.t.mulCopy(L0, a.y)
	// field_mul ( a->u, a->x, a->t );
	a.u.mulCopy(a.x, a.t)
	// field_mul ( a->t, a->u, L0 );
	a.t.mulCopy(a.u, L0)
	// field_mul ( a->y, a->x, a->t );
	a.y.mulCopy(a.x, a.t)
	// field_isr ( L0, a->y );
	L0.isr(a.y)
	// field_mul ( a->y, a->u, L0 );
	a.y.mulCopy(a.u, L0)
	// field_sqr ( L1, L0 );
	L1.squareCopy(L0)
	// field_mul ( a->u, a->t, L1 );
	a.u.mulCopy(a.t, L1)
	// field_mul ( a->t, a->x, a->u );
	a.t.mulCopy(a.x, a.u)
	// field_add ( a->x, sz, sz );
	a.x.add(sz, sz)
	// field_mul ( L0, a->u, a->x );
	L0.mulCopy(a.u, a.x)
	// field_copy ( a->x, a->z );
	a.x = a.z.copy()
	// field_neg ( L1, a->x );
	L1.neg(a.x)
	// field_addw ( L1, 1 );
	L1.addW(1)
	// field_mul ( a->x, L1, L0 );
	a.x.mulCopy(L1, L0)
	// field_mul ( L0, a->u, a->y );
	L0.mulCopy(a.u, a.y)
	// field_addw ( a->z, 1 );
	a.z.addW(1)
	// field_mul ( a->y, a->z, L0 );
	a.y.mulCopy(a.z, L0)
	// field_subw( a->t, 1 );
	a.t.subW(1)
	// mask_t ret = field_is_zero( a->t );
	// XXX maybe related with constant time
	ret := a.t.zero()
	// field_set_ui( a->z, 1 );
	a.z.setUi(1)
	// field_copy ( a->t, a->x );
	a.t = a.x.copy()
	// field_copy ( a->u, a->y );
	a.u = a.y.copy()
	// return ret;

	return a, !ret
}
