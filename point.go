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

func (p *extensibleCoordinates) twist() *twExtensible {
	x := new(bigNumber)
	y := new(bigNumber)
	z := new(bigNumber)
	t := new(bigNumber)
	u := new(bigNumber)

	l0 := new(bigNumber)
	l1 := new(bigNumber)

	u.square(p.z)
	y.square(p.x)
	z.sub(u, y)
	y.add(z, z)
	u.add(y, y)
	y.sub(p.z, p.x)
	x.mul(y, p.y)
	z.sub(p.z, p.y)
	t.mul(z, x)
	l1.mul(t, u)

	x.mul(t, l1)
	l0.isr(x)
	u.mul(t, l0)
	l1.square(l0)
	t.mul(x, l1)
	l1.add(p.x, p.y)
	l0.sub(p.x, p.y)
	x.mul(t, l0)
	l0.add(x, l1)
	t.sub(l1, x)
	x.mul(l0, u)
	x.addW(-y.zeroMask())
	y.mul(t, u)
	y.addW(-z.zeroMask())
	z.setUi(1 + uint64(y.zeroMask()))
	t = x.copy()
	u = y.copy()

	return &twExtensible{x, y, z, t, u}
}

//XXX unused
func (p *extensibleCoordinates) double() *extensibleCoordinates {
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

func (p *twPNiels) TwistedExtensible() *twExtensible {
	u := new(bigNumber).add(p.n.b, p.n.a)
	t := new(bigNumber).sub(p.n.b, p.n.a)

	return &twExtensible{
		x: new(bigNumber).mul(p.z, t),
		y: new(bigNumber).mul(p.z, u),
		z: new(bigNumber).square(p.z),
		t: t,
		u: u,
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

//XXX this may not always work
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

	//XXX PERF: should it be in-place?
	return &twExtensible{x, y, z, t, u}
}

type twExtensible struct {
	x, y, z, t, u *bigNumber
}

func (p *twExtensible) setIdentity() {
	p.x = p.x.setUi(0)
	p.y = p.y.setUi(1)
	p.z = p.z.setUi(1)
	p.t = p.t.setUi(0)
	p.u = p.u.setUi(0)
}

func (p *twExtensible) copy() *twExtensible {
	return &twExtensible{
		x: p.x.copy(),
		y: p.y.copy(),
		z: p.z.copy(),
		t: p.t.copy(),
		u: p.u.copy(),
	}
}

func (p *twExtensible) addTwPNiels(a *twPNiels) *twExtensible {
	// field_mul ( L0, e->z, a->z );
	L0 := new(bigNumber).mul(p.z, a.z)
	// field_copy ( e->z, L0 );
	p.z = L0.copy()
	// add_tw_niels_to_tw_extensible( e, a->n );
	return p.addTwNiels(a.n)
}

func (p *twExtensible) subTwPNiels(a *twPNiels) *twExtensible {
	//XXX PERF: should it be in-place?
	e := p.copy()
	e.z = e.z.mul(e.z, a.z)

	return e.subTwNiels(a.n)
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
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l2 = l2.mul(p2.z, p.x)
	l1 = l1.mul(p.z, p2.x)
	l0 = l0.sub(l2, l1)

	l4 := l0.zeroMask()

	l2 = l2.mul(p2.z, p.y)
	l1 = l1.mul(p.z, p2.y)
	l0 = l0.sub(l2, l1)

	l3 := l0.zeroMask()

	return (l4 & l3) == 0xffffffff
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

	//XXX PERF: should it be in-place?
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

	//XXX PERF: should it be in-place?
	return &twExtensible{x, y, z, t, u}
}

func (p *twExtensible) subTwNiels(e *twNiels) *twExtensible {
	l0 := new(bigNumber)
	l1 := new(bigNumber)

	//XXX PERF: should it be in-place?
	d := p.copy()

	l1 = l1.subxRaw(d.y, d.x)
	l0 = l0.mul(e.b, l1)
	l1 = l1.addRaw(d.x, d.y)
	d.y = d.y.mul(e.a, l1)
	l1 = l1.mul(d.u, d.t)
	d.x = d.x.mul(e.c, l1)
	d.u = d.u.addRaw(l0, d.y)
	d.t = d.t.subxRaw(d.y, l0)
	d.y = d.y.addRaw(d.x, d.z)
	l0 = l0.subxRaw(d.z, d.x)

	d.z = d.z.mul(l0, d.y)
	d.x = d.x.mul(d.y, d.t)
	d.y = d.y.mul(l0, d.u)

	return d
}

func (p *twExtensible) untwistAndDoubleAndSerialize() *bigNumber {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)
	l3 := new(bigNumber)
	b := new(bigNumber)

	l3.mul(p.y, p.x)
	b.add(p.y, p.x)
	l1.square(b)
	l2.add(l3, l3)
	b.sub(l1, l2)
	l2.square(p.z)
	l1.square(l2)
	b.add(b, b)
	l2.mulWSignedCurveConstant(b, curveDSigned-1)
	b.mulWSignedCurveConstant(l2, curveDSigned-1)
	l0.mul(l2, l1)
	l2.mul(b, l0)
	l0.isr(l2)
	l1.mul(b, l0)

	//XXX This is included in the original code, but it seems not to be used
	//b = b.square(l0)
	//l0 = l0.mul(l2, b)

	return b.mul(l1, l3)
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

// See Hisil, formula 5.1
func (p *homogeneousProjective) double() *homogeneousProjective {
	x1 := p.x
	y1 := p.y
	z1 := p.z

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

	a := new(bigNumber).mul(z1, z2)
	b := new(bigNumber).square(a)
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

type montgomery struct {
	z0, xd, zd, xa, za *bigNumber
}

func (a *montgomery) montgomeryStep() {
	L0 := new(bigNumber)
	L1 := new(bigNumber)
	L0.addRaw(a.zd, a.xd)
	L1.subxRaw(a.xd, a.zd)
	a.zd.subxRaw(a.xa, a.za)
	a.xd.mul(L0, a.zd)
	a.zd.addRaw(a.za, a.xa)
	a.za.mul(L1, a.zd)
	a.xa.addRaw(a.za, a.xd)
	a.zd.square(a.xa)
	a.xa.mul(a.z0, a.zd)
	a.zd.subxRaw(a.xd, a.za)
	a.za.square(a.zd)
	a.xd.square(L0)
	L0.square(L1)
	a.zd.mulWSignedCurveConstant(a.xd, 1-curveDSigned) /* FIXME PERF MULW */
	L1.subxRaw(a.xd, L0)
	a.xd.mul(L0, a.zd)
	L0.subRaw(a.zd, L1)
	L0.bias(4 - 2*1 /*is32 ? 2 : 4*/)
	//XXX 64bits don't need this reduce
	L0.weakReduce()
	a.zd.mul(L0, L1)
}

func (a *montgomery) serialize(sbz *bigNumber) (b *bigNumber, ok uint32) {
	L0 := new(bigNumber)
	L1 := new(bigNumber)
	L2 := new(bigNumber)
	L3 := new(bigNumber)
	b = new(bigNumber)

	L3.mul(a.z0, a.zd)
	L1.sub(L3, a.xd)
	L3.mul(a.za, L1)
	L2.mul(a.z0, a.xd)
	L1.sub(L2, a.zd)
	L0.mul(a.xa, L1)
	L2.add(L0, L3)
	L1.sub(L3, L0)
	L3.mul(L1, L2)
	L2 = a.z0.copy()
	L2.addW(1)
	L0.square(L2)
	L1.mulWSignedCurveConstant(L0, curveDSigned-1)
	L2.add(a.z0, a.z0)
	L0.add(L2, L2)
	L2.add(L0, L1)
	L0.mul(a.xd, L2)
	L5 := a.zd.zeroMask()
	L6 := -L5
	// constant_time_mask ( L1, L0, sizeof(L1), L5 );
	mask(L1, L0, L5)
	L2.add(L1, a.zd)
	L4 := ^L5
	L1.mul(sbz, L3)
	L1.addW(L6)
	L3.mul(L2, L1)
	L1.mul(L3, L2)
	L2.mul(L3, a.xd)
	L3.mul(L1, L2)
	L0.isr(L3)
	L2.mul(L1, L0)
	L1.square(L0)
	L0.mul(L3, L1)
	// constant_time_mask ( b, L2, sizeof(L1), L4 );
	mask(b, L2, L4)
	L0.subW(1)
	L5 = L0.zeroMask()
	L4 = sbz.zeroMask()

	return b, L5 | L4
}

func (a *montgomery) deserialize(sz *bigNumber) {
	a.z0 = new(bigNumber).square(sz)
	a.xd = new(bigNumber).setUi(1)
	a.zd = new(bigNumber).setUi(0)
	a.xa = new(bigNumber).setUi(1)
	a.za = a.z0.copy()
}
