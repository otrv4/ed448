package ed448

import (
	"errors"
)

// Point is a interface of a Ed448 point
type Point interface {
	IsOnCurve() bool
	Equals(q Point) bool
	EqualsMask(q Point) uint32
	Copy() Point
	Add(q, r Point)
	Sub(q, r Point)
	Double() Point
	Encode() []byte
	Decode(src []byte, identity bool) (bool, error)
	DSAEncode() []byte
	DSADecode(src []byte) bool
}

// Extended Homogenous Projective coordinates: (X : Y : T : Z), which
// correspond to the affine point (X/Z, Y/Z) with Z ≠ 0
type twExtendedPoint struct {
	x, y, z, t *bigNumber
}

func (p *twExtendedPoint) isOnCurve() bool {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	// x * y == z * t
	a.mul(p.x, p.y)
	b.mul(p.z, p.t)
	valid := a.decafEq(b)

	// y^2 - x^2 == z^2 - t^2 * (1 - D)
	a.square(p.x)
	b.square(p.y)
	a.sub(b, a)
	b.square(p.t)
	c.mulW(b, 1-edwardsD)
	b.square(p.z)
	b.sub(b, c)
	valid &= a.decafEq(b)
	valid &= ^(p.z.decafEq(bigZero))

	return valid == decafTrue
}

func (p *twExtendedPoint) copy() *twExtendedPoint {
	n := &twExtendedPoint{}
	n.x = p.x.copy()
	n.y = p.y.copy()
	n.z = p.z.copy()
	n.t = p.t.copy()
	return n
}

func (p *twExtendedPoint) setIdentity() {
	p.x.setUI(0x00)
	p.y.setUI(0x01)
	p.z.setUI(0x01)
	p.t.setUI(0x00)
}

func (p *twExtendedPoint) equals(q *twExtendedPoint) word {
	a, b := &bigNumber{}, &bigNumber{}
	a.mul(p.y, q.x)
	b.mul(q.y, p.x)
	return a.decafEq(b)
}

func (p *twExtendedPoint) add(q, r *twExtendedPoint) {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.sub(q.y, q.x)
	c.sub(r.y, r.x)
	d.addRaw(r.y, r.x)
	a.mul(c, b)
	b.addRaw(q.y, q.x)
	p.y.mul(d, b)
	b.mul(r.t, q.t)
	p.x.mulW(b, 2*effD)
	b.addRaw(a, p.y)
	c.sub(p.y, a)
	a.mul(q.z, r.z)
	a.addRaw(a, a)
	p.y.addRaw(a, p.x)
	a.sub(a, p.x)
	p.z.mul(a, p.y)
	p.x.mul(p.y, c)
	p.y.mul(a, b)
	p.t.mul(b, c)
}

func (p *twExtendedPoint) sub(q *twExtendedPoint, r *twExtendedPoint) {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.sub(q.y, q.x)
	d.sub(r.y, r.x)
	c.addRaw(r.y, r.x)
	a.mul(c, b)
	b.addRaw(q.y, q.x)
	p.y.mul(d, b)
	b.mul(r.t, q.t)
	p.x.mulW(b, 2-2*edwardsD)
	b.addRaw(a, p.y)
	c.sub(p.y, a)
	a.mul(q.z, r.z)
	a.addRaw(a, a)
	p.y.sub(a, p.x)
	a.addRaw(a, p.x)
	p.z.mul(a, p.y)
	p.x.mul(p.y, c)
	p.y.mul(a, b)
	p.t.mul(b, c)
}

func (p *twExtendedPoint) negate(q *twExtendedPoint) {
	p.x.sub(bigZero, q.x)
	p.y = q.y.copy()
	p.z = q.z.copy()
	p.t.sub(bigZero, q.t)
}

func (p *twExtendedPoint) torque(q *twExtendedPoint) {
	p.x.sub(bigZero, q.x)
	p.y.sub(bigZero, q.y)
	p.z = q.z.copy()
	p.t = q.t.copy()
}

// Based on Hisil's formula 5.1.3: Doubling in E^e
func (p *twExtendedPoint) doubleInternal(beforeDouble bool) *twExtendedPoint {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	c.square(p.x)
	a.square(p.y)
	d.addRaw(c, a)
	p.t.addRaw(p.y, p.x)
	b.square(p.t)
	exponentBias := word(0x03)
	b.subXBias(b, d, exponentBias)
	p.t.sub(a, c)
	p.x.square(p.z)
	p.z.addRaw(p.x, p.x)
	exponentBias = word(0x04)
	a.subXBias(p.z, p.t, exponentBias)
	p.x.mul(a, b)
	p.z.mul(p.t, a)
	p.y.mul(p.t, d)
	if !beforeDouble {
		p.t.mul(b, d)
	}
	return p
}

func (p *twExtendedPoint) double() *twExtendedPoint {
	return p.doubleInternal(false)
}

func (p *twExtendedPoint) decafEncode(dst []byte) {
	if len(dst) != fieldBytes {
		panic("Attempted an encode with a destination that is not 56 bytes")
	}
	t, overT := allZeros, allZeros
	serialize(dst, p.deisogenize(t, overT))
}

func (p *twExtendedPoint) deisogenize(t, overT word) *bigNumber {
	a, b, c, d, s := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	a.mulWSignedCurveConstant(p.y, 1-(edwardsD))
	c.mul(a, p.t)
	a.mul(p.x, p.z)
	d.sub(c, a)
	a.add(p.z, p.y)
	b.sub(p.z, p.y)
	c.mul(b, a)
	b.mulWSignedCurveConstant(c, (-(edwardsD)))
	a.isr(b)
	b.mulWSignedCurveConstant(a, (-(edwardsD)))
	c.mul(b, a)
	a.mul(c, d)
	d.add(b, b)
	c.mul(d, p.z)
	b.decafCondNegate(overT ^ (^(highBit(c))))
	c.decafCondNegate(overT ^ (^(highBit(c))))
	d.mul(b, p.y)
	s.add(a, d)
	s.decafCondNegate(overT ^ highBit(s))

	return s
}

// TODO: should this return a bool and an error?
func decafDecodeOld(dst *twExtendedPoint, src serialized, useIdentity bool) (word, error) {
	a, b, c, d, e := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	n, succ := deserializeReturnMask(src)
	zero := n.decafEq(bigZero)
	if useIdentity {
		succ &= decafTrue | ^zero
	} else {
		succ &= decafFalse | ^zero
	}
	succ &= ^highBit(n)

	a.square(n)
	dst.z.sub(bigOne, a)
	b.square(dst.z)
	c.mulWSignedCurveConstant(a, 4-4*(edwardsD))
	c.add(c, b)
	b.mul(c, a)
	d.isr(b)
	e.square(d)
	a.mul(e, b)
	a.add(a, bigOne)
	succ &= ^(a.decafEq(bigZero))
	b.mul(c, d)
	d.decafCondNegate(highBit(b))
	dst.x.add(n, n)
	c.mul(d, n)
	b.sub(bigTwo, dst.z)
	a.mul(b, c)
	dst.y.mul(a, dst.z)
	dst.t.mul(dst.x, a)
	dst.y[0] -= zero

	var err error
	if succ != decafTrue {
		err = errors.New("unable to decode given point")
		return succ, err
	}
	return succ, err
}

func decafDecode(dst *twExtendedPoint, src serialized, useIdentity bool) (word, error) {
	s2 := &bigNumber{}
	num := &bigNumber{}
	tmp := &bigNumber{}
	tmp2 := &bigNumber{}
	isr := dst.x
	den := dst.t
	ynum := dst.z

	s, succ := deserializeReturnMask(src)
	zero := s.decafEq(bigZero)
	if useIdentity {
		succ &= decafTrue | ^zero
	} else {
		succ &= decafFalse | ^zero
	}
	succ &= ^lowBit(s)

	s2.square(s)
	den.sub(bigOne, s2)
	ynum.add(bigOne, s2)
	num.mulWSignedCurveConstant(s2, -4*(edwardsD-1))
	tmp.square(den)
	num.add(tmp, num)
	tmp2.mul(num, tmp)
	succ &= isr.isr(tmp2)
	tmp.mul(isr, den)
	dst.y.mul(tmp, ynum)
	tmp2.mul(tmp, s)
	tmp2.add(tmp2, tmp2)
	tmp.mul(tmp2, isr)
	dst.x.mul(tmp, num)
	tmp.mul(tmp2, factor)

	dst.x.decafCondNegate(lowBit(tmp))
	dst.z.set(bigOne)
	dst.t.mul(dst.x, dst.y)

	var err error
	if succ != decafTrue {
		err = errors.New("unable to decode given point")
		return succ, err
	}

	return succ, err
}

func (p *twExtendedPoint) dsaLikeEncode(dst []byte) {
	if len(dst) != dsaFieldBytes {
		panic("Attempted to encode with a destination that is not 57 bytes")
	}

	x, y, z, t, u := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	// untwist by 4-isogeny: 2xy/(y^+x^2), (y^2-x^2)/(2z^2-y^2+x^2)
	x.square(p.x)
	t.square(p.y)
	u.add(x, t)
	z.add(p.y, p.x)
	y.square(z)
	y.sub(u, y)
	z.sub(t, x)
	x.square(p.z)
	t.add(x, x)
	t.sub(t, z)
	x.mul(t, y)
	y.mul(z, u)
	z.mul(u, t)

	u.set(bigZero)

	// convert to affine
	z = invert(z)
	t.mul(x, z)
	x.mul(y, z)

	dst[fieldBytes] = byte(allZeros)
	dsaLikeSerialize(dst[:], x)
	dst[fieldBytes] |= byte(zeroMask & lowBit(t))

	// wipe out
	x.set(bigZero)
	y.set(bigZero)
	z.set(bigZero)
	t.set(bigZero)
}

func dsaLikeDecode(p *twExtendedPoint, srcOrg []byte) word {
	if len(srcOrg) != dsaFieldBytes {
		panic("Attempted to decode with a source that is not 57 bytes")
	}
	src := append([]byte{}, srcOrg...)

	succ := decafTrue
	var cofactorMask uint = zeroMask

	low := ^isZeroMask(word(src[fieldBytes] & zeroMask))
	src[fieldBytes] &= byte(^(cofactorMask))

	succ = isZeroMask(word(src[fieldBytes]))
	succ &= dsaLikeDeserialize(p.y, src[:])

	p.x.square(p.y)
	p.z.sub(bigOne, p.x)                       // num = 1- (y^2)
	p.t.mulWSignedCurveConstant(p.x, edwardsD) // d * (y^2)
	p.t.sub(bigOne, p.t)                       // denom = 1 - d * (y^2)
	p.x.mul(p.z, p.t)
	p.t.isr(p.x)      // 1/sqrt(num * denom) // implement it with check
	p.x.mul(p.t, p.z) // sqrt(num / denom)
	p.x.decafCondNegate(^lowBit(p.x) ^ low)
	p.z = bigOne.copy()

	// 4-isogeny 2xy/(y^2-ax^2), (y^2+ax^2)/(2-y^2-ax^2)
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	c.square(p.x)
	a.square(p.y)
	d.add(c, a)
	p.t.add(p.y, p.x)
	b.square(p.t)
	b.sub(b, d)
	p.t.sub(a, c)
	p.x.square(p.z)
	p.z.add(p.x, p.x)
	a.sub(p.z, d)
	p.x.mul(a, b)
	p.z.mul(p.t, a)
	p.y.mul(p.t, d)
	p.t.mul(b, d)

	// wipe out
	a.set(bigZero)
	b.set(bigZero)
	c.set(bigZero)
	d.set(bigZero)
	src = make([]byte, 57)

	ok := p.isOnCurve()
	if !ok {
		return decafFalse
	}

	res := pointScalarMul(p, scalarOneFourth.(*scalar))
	p.x = res.x
	p.y = res.y
	p.z = res.z
	p.t = res.t

	return succ
}

func (p *twExtendedPoint) addNielsToExtended(np *twNiels, beforeDouble bool) {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.sub(p.y, p.x)
	a.mul(np.a, b)
	b.addRaw(p.x, p.y)
	p.y.mul(np.b, b)
	p.x.mul(np.c, p.t)
	c.addRaw(a, p.y)
	b.sub(p.y, a)
	p.y.sub(p.z, p.x)
	a.addRaw(p.x, p.z)
	p.z.mul(a, p.y)
	p.x.mul(p.y, b)
	p.y.mul(a, c)
	if !beforeDouble {
		p.t.mul(b, c)
	}
}

func (p *twExtendedPoint) subNielsFromExtendedPoint(np *twNiels, beforeDouble bool) {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.sub(p.y, p.x)
	a.mul(np.b, b)
	b.addRaw(p.x, p.y)
	p.y.mul(np.a, b)
	p.x.mul(np.c, p.t)
	c.addRaw(a, p.y)
	b.sub(p.y, a)
	p.y.addRaw(p.z, p.x)
	a.sub(p.z, p.x)
	p.z.mul(a, p.y)
	p.x.mul(p.y, b)
	p.y.mul(a, c)
	if !beforeDouble {
		p.t.mul(b, c)
	}
}

func (p *twExtendedPoint) addProjectiveNielsToExtended(np *twPNiels, beforeDouble bool) {
	tmp := &bigNumber{}
	tmp.mul(p.z, np.z)
	p.z = tmp.copy()
	p.addNielsToExtended(np.n, beforeDouble)
}

func (p *twExtendedPoint) subProjectiveNielsFromExtendedPoint(p2 *twPNiels, beforeDouble bool) {
	tmp := &bigNumber{}
	tmp.mul(p.z, p2.z)
	p.z = tmp.copy()
	p.subNielsFromExtendedPoint(p2.n, beforeDouble)
}

// Convert from the extended twisted Edwards representation of a point to affine
// Given (X : Y : Z : T), compute X/Z^2, Y/Z^3 and ignore T.
// If the point is ∞ it returns 0, 0.
// TODO: check me
func (p *twExtendedPoint) toAffine() *affineCoordinates {
	out := &affineCoordinates{
		&bigNumber{},
		&bigNumber{},
	}

	if p.equals(identity) == decafTrue || p.z.decafEq(bigZero) == decafTrue {
		return out
	}

	s, t, r := &bigNumber{}, &bigNumber{}, &bigNumber{}

	r = invert(p.z)
	s.square(r)

	out.x.mul(p.x, s).strongReduce()
	t.mul(p.y, s)
	out.y.mul(t, r).strongReduce()

	return out
}

//TODO: extendedPoint should not know about twNiels
func (np *twNiels) toExtended() *twExtendedPoint {
	p := &twExtendedPoint{
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
	}

	p.y.add(np.b, np.a)
	p.x.sub(np.b, np.a)
	p.t.mul(p.y, p.x)
	copy(p.z[:], bigOne[:])
	return p
}

func (p *twExtendedPoint) toPNiels() *twPNiels {
	a := &bigNumber{}
	a.sub(p.y, p.x)

	b := &bigNumber{}
	b.add(p.x, p.y)

	c := &bigNumber{}
	c.mulWSignedCurveConstant(p.t, 2*edwardsD-2)

	z := &bigNumber{}
	z.add(p.z, p.z)

	return &twPNiels{
		&twNiels{a, b, c},
		z,
	}
}

func (p *twExtendedPoint) isogenizeToMontgomery() []byte {
	var out [56]byte

	p.t = invert(p.x) // 1/x
	p.z.mul(p.t, p.y) // y/x
	p.y.square(p.z)   // (y/x)^2

	dsaLikeSerialize(out[:], p.y)
	return out[:]
}

func pointScalarMul(p *twExtendedPoint, s *scalar) *twExtendedPoint {
	const decafWindowBits = 5            //move this to const file
	const window = decafWindowBits       //5
	const windowMask = (1 << window) - 1 //0x0001f 31
	const windowTMask = windowMask >> 1  //0x0000f 15
	const nTable = 1 << (window - 1)     //0x00010 16

	out := &twExtendedPoint{}

	scalar1x := &scalar{}
	scalar1x.add(s, decafPrecompTable.scalarAdjustment)
	scalar1x.halve(scalar1x)

	multiples := p.prepareFixedWindow(nTable)

	first := true
	for i := scalarBits - ((scalarBits - 1) % window) - 1; i >= 0; i -= window {
		bits := scalar1x[i/wordBits] >> uint(i%wordBits)
		if i%wordBits >= wordBits-window && i/wordBits < scalarWords-1 {
			bits ^= scalar1x[i/wordBits+1] << uint(wordBits-(i%wordBits))
		}
		bits &= windowMask
		inv := (bits >> (window - 1)) - 1
		bits ^= inv

		//Add in from table.  Compute out.t (point) only on last iteration.
		pNeg := constTimeLookup(multiples, word(bits&windowTMask))
		pNeg.n.conditionalNegate(inv)

		if first {
			out = pNeg.toExtendedPoint()
			first = false
		} else {
			//Using Hisil et al's lookahead method instead of
			//extensible here for no particular reason.  Double
			//5 (window) times, but only compute out.t on the last one.
			for j := 0; j < window-1; j++ {
				out.doubleInternal(true)
			}
			out.doubleInternal(false)
			out.addProjectiveNielsToExtended(pNeg, false)
		}
	}
	return out
}

func precomputedScalarMul(s *scalar) *twExtendedPoint {
	p := &twExtendedPoint{
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
	}
	scalar2 := &scalar{}
	scalar2.add(s, decafPrecompTable.scalarAdjustment)
	scalar2.halve(scalar2)

	var np *twNiels
	for i := int(decafCombSpacing - 1); i >= 0; i-- {
		if i != int(decafCombSpacing-1) {
			p.doubleInternal(false)
		}

		for j := uintZero; j < decafCombNumber; j++ {
			var tab word
			for k := uintZero; k < decafCombTeeth; k++ {
				bit := uint(i) + decafCombSpacing*(k+j*decafCombTeeth)
				if bit < scalarBits {
					tab |= (scalar2[bit/wordBits] >> (bit % wordBits) & 1) << k
				}
			}

			invert := (sword(tab) >> (decafCombTeeth - 1)) - 1
			tab ^= word(invert)
			tab &= (1 << (decafCombTeeth - 1)) - 1

			index := word(((j << (decafCombTeeth - 1)) + uint(tab)))
			np = decafPrecompTable.lookup(index)

			np.conditionalNegate(word(invert))

			if i != int(decafCombSpacing-1) || j != 0 {
				p.addNielsToExtended(np, j == decafCombNumber-1 && i != 0)
			} else {
				p = np.toExtended()
			}
		}
	}

	return p
}

// using the montgomery ladder
// TODO: implement the one not using montgomery?
func directPointScalarMul(p [fieldBytes]byte, s *scalar, useIdentity word) ([fieldBytes]byte, word) {
	var out [56]byte
	xa, xs, zs, l0, l1 := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	s0, succ := deserializeReturnMask(p)
	succ &= useIdentity | ^(s0.decafEq(bigZero))

	// Prepare the Montgomery ladder: Q = 1:0, P+Q = P
	xa.square(s0)
	x0 := xa.copy()
	za := bigOne.copy()
	xd := bigOne.copy()
	zd := bigZero.copy()

	pflip := word(0x00)
	for i := scalarBits - 1; i >= 0; i-- {
		// Augmented Montgomery ladder
		flip := -(s[i/wordBits] >> word(i%wordBits) & 1)
		// Differential add
		xs.addRaw(xa, za)
		zs.sub(xa, za)
		xa.addRaw(xd, zd)
		za.sub(xd, zd)

		l0.decafConstTimeSel(xa, xs, flip^pflip)
		l1.decafConstTimeSel(za, zs, flip^pflip)

		xd.mul(xa, zs)
		zd.mul(xs, za)
		xs.addRaw(xd, zd)
		zd.sub(xd, zd)
		zs.mul(zd, s0)
		xa.square(xs)
		za.square(zs)

		// double
		zd.square(l0)
		l0.square(l1)
		l1.sub(zd, l0)
		xd.mul(l0, zd)
		zd.mulWSignedCurveConstant(l1, 1-(edwardsD))
		l0.addRaw(l0, zd)
		zd.mul(l0, l1)

		pflip = flip
	}
	xa.conditionalSwap(xd, pflip)
	za.conditionalSwap(zd, pflip)

	// TODO: should be constant time
	// reserialize TODO: simplify this reserialization
	xzD, xzA, xzS, den, l2, l3 := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	xzS.mul(xs, zs)
	xzD.mul(xd, zd)
	xzA.mul(xa, za)
	zeroOut := xzD.decafEq(bigZero)
	xzD[0] -= zeroOut // make xzD always nonzero
	zCase := zeroOut | xzA.decafEq(bigZero)
	zeroZA := za.decafEq(bigZero)

	// Curve test in zcase, compute x0^2 + (2d-4)x0 + 1
	l0.add(x0, bigOne)
	l1.square(l0)
	l0.mulWSignedCurveConstant(x0, -4*edwardsD)
	l1.add(l1, l0)
	xzA.decafConstTimeSel(xzA, l1, zCase)

	// Compute denominator = x0 xa za xd zd
	l0.mul(x0, xzA)
	l1.mul(l0, xzD)
	den.isr(l1)

	// Check that the square root is valid.
	l2.square(den)
	l3.mul(l0, l2) // x0 xa za den^2 = 1/xzD
	l0.mul(l1, l2)
	l0.add(l0, bigOne)
	succ &= ^highBit(s0) & ^(l0.decafEq(bigZero))

	// Compute y/x for input and output point.
	l1.mul(x0, xd)
	l1.sub(zd, l1)
	l0.mul(za, l1) // L0 = "opq"
	l1.mul(x0, zd)
	l1.sub(l1, xd)
	l2.mul(xa, l1) // L2 = "pqr"
	l1.sub(l0, l2)
	l0.add(l0, l2)
	l2.mul(l1, den) // L2 = y0 / x0
	l1.mul(l0, den) // L1 = yO / xO
	sflip := lowBit(l1) ^ lowBit(l2) | zeroZA

	// If xa==0 or za ==0: return 0
	// Else if za == 0: return s0 * (sflip ? zd : xd)^2 * L3
	// Else if zd == 0: return s0 * (sflip ? zd : xd)^2 * L3
	// Else if pflip: return xs * zs * (sflip ? zd : xd) * L3
	// Else: return s0 * xs * zs * (sflip ? zd : xd) * den
	xd.decafConstTimeSel(xd, zd, sflip)
	den.decafConstTimeSel(den, l3, pflip|zCase)
	xzS.decafConstTimeSel(xzS, xd, zCase)
	s0.decafConstTimeSel(s0, bigOne, pflip & ^zCase)
	s0.decafConstTimeSel(s0, bigZero, zeroOut)

	l0.mul(xd, den)
	l1.mul(l0, s0)
	l0.mul(l1, xzS)

	l0.conditionalNegate(highBit(l0))
	serialize(out[:], l0)

	return out, succ
}

// exposed methods

// NewPoint returns an Ed448 point from 4 arrays of 16 uint32.
func NewPoint(a [nLimbs]uint32, b [nLimbs]uint32, c [nLimbs]uint32, d [nLimbs]uint32) Point {
	x, y, z, t := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	for i := 0; i < nLimbs; i++ {
		x[i] = word(a[i])
		y[i] = word(b[i])
		z[i] = word(c[i])
		t[i] = word(d[i])
	}

	return &twExtendedPoint{x, y, z, t}
}

// NewPointFromBytes returns an Ed448 point from a byte slice.
func NewPointFromBytes(in ...[]byte) Point {
	if len(in) > 1 {
		panic("too many arguments to function call")
	}

	out := &twExtendedPoint{
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
	}

	if in == nil {
		return out
	}

	bytes := in[0][:]
	if len(bytes) != 56 {
		panic("byte input needs to be size 56")
	}
	tmpIn := [fieldBytes]byte{}
	copy(tmpIn[:], bytes[:])
	decafDecodeOld(out, tmpIn, false)

	return out
}

// IsOnCurve reports whether the given point (p) lies on the curve.
func (p *twExtendedPoint) IsOnCurve() bool {
	return p.isOnCurve()
}

// Equals compares whether two points (p, q) are equal.
func (p *twExtendedPoint) Equals(q Point) bool {
	valid := p.equals(q.(*twExtendedPoint))
	return valid == decafTrue
}

// EqualsMask compares whether two points (p, q) are equal.
func (p *twExtendedPoint) EqualsMask(q Point) uint32 {
	return uint32(p.equals(q.(*twExtendedPoint)))
}

// Copy returns a copy of a given point (p).
func (p *twExtendedPoint) Copy() Point {
	p.copy()
	return Point(p)
}

// Add gives the sum of two points (q, r) and produces a thrid point (p).
func (p *twExtendedPoint) Add(q, r Point) {
	p.add(q.(*twExtendedPoint), r.(*twExtendedPoint))
}

// Sub gives the subtraction of two points (q, r) and produces a thrid point (p).
func (p *twExtendedPoint) Sub(q, r Point) {
	p.sub(q.(*twExtendedPoint), r.(*twExtendedPoint))
}

// Double gives the doubling of a point (p).
func (p *twExtendedPoint) Double() Point {
	return p.double()
}

// Encode returns the encoding of a point (p) as a sequence of bytes.
// This uses the 'decaf' technique. See `Decaf: Eliminating cofactors through
// point compression``, Mike Hamburg, Advances in Cryptology (Crypto 2015).
// This technique removes the cofactor through quotients and isogenies.
// The internal representation of points is as "even" elements of a twisted
// Edwards curve with a=-1. Using this subgroup removes a factor of 2 from the
// cofactor. The remaining factor of 2 is removed with a quotient group: any two
// points which differ by an element of the 2- or 4-torsion subgroup are
// considered equal to each other.
// When a point is written out to wire format, it is converted (by isogeny)
// to a Jacobi quartic curve. The x-coordinate of this point is written out.
func (p *twExtendedPoint) Encode() []byte {
	out := make([]byte, fieldBytes)
	p.decafEncode(out)
	return out
}

// Decode gives the decoding a point from a sequence of bytes (src).
// Every point has a unique encoding, so not every sequence of bytes is a valid
// encoding.  If an invalid encoding is given, the output is undefined.
// Set 'useIdentity' true  if the identity is a legal input.
func (p *twExtendedPoint) Decode(src []byte, useIdentity bool) (bool, error) {
	ser := [fieldBytes]byte{}
	copy(ser[:], src[:])

	valid, err := decafDecodeOld(p, ser, useIdentity)
	if err != nil {
		return false, err
	}
	return valid == decafTrue, nil
}

// DSAEncode returns the encoding of a point (p) as a sequence of bytes.
// This uses the eddsa techinique. See ``Edwards-Curve Digital Signature
// Algorithm (EdDSA)``, S. Josefsson and I. Liusvaara, Internet Research Task
// Force (IRTF).
// Multiplies the point to the cofactor first.
func (p *twExtendedPoint) DSAEncode() []byte {
	out := make([]byte, dsaFieldBytes)
	p.dsaLikeEncode(out)
	return out
}

// DSADecode gives the decoding of a point (p) as a sequence of bytes (src).
func (p *twExtendedPoint) DSADecode(src []byte) bool {
	ok := dsaLikeDecode(p, src)

	return ok == decafTrue
}

// PointScalarMul returns the multiplication of a given point (p) by a given
// scalar (a): q * a.
func PointScalarMul(q Point, a Scalar) Point {
	return pointScalarMul(q.(*twExtendedPoint), a.(*scalar))
}

// PrecomputedScalarMul returns the multiplication of a given scalar (a) by the
// precomputed base point of the curve: basePoint * a.
func PrecomputedScalarMul(a Scalar) Point {
	return precomputedScalarMul(a.(*scalar))
}

// PointDoubleScalarMul returns the addition of two multiplications: a given
// point (q) by a given scalar (a) and a given point (r) by a given scalar (b):
// q * a + r * b.
func PointDoubleScalarMul(q, r Point, a, b Scalar) Point {
	return doubleScalarMul(q.(*twExtendedPoint), r.(*twExtendedPoint), a.(*scalar), b.(*scalar))
}

// PointDoubleScalarMulNonsecret returns the addition of two multiplications:
// a given point (q) by a given scalar (b) and the base point of the curve by a
// given scalar (a): q * b + basePoint * a.
// @warning: This function takes variable time, and may leak the scalars used.
// It is designed for signature verification.
func PointDoubleScalarMulNonsecret(q Point, a, b Scalar) Point {
	return decafDoubleNonSecretScalarMul(q.(*twExtendedPoint), a.(*scalar), b.(*scalar))
}
