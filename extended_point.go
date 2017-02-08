package ed448

type twExtendedPoint struct {
	x, y, z, t *bigNumber
}

func (p *twExtendedPoint) isValidPoint() bool {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	a.mul(p.x, p.y)
	b.mul(p.z, p.t)
	valid := a.decafEq(b)
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
	p.x.setUI(0)
	p.y.setUI(1)
	p.z.setUI(1)
	p.t.setUI(0)
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
	p.x.mulW(b, 2-2*edwardsD)
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

// Based on Hisil's formula 5.1.3: Doubling in E^e
func (p *twExtendedPoint) double(beforeDouble bool) *twExtendedPoint {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	c.square(p.x)
	a.square(p.y)
	d.addRaw(c, a)
	p.t.addRaw(p.y, p.x)
	b.square(p.t)
	exponentBias := word(3)
	b.subXBias(b, d, exponentBias)
	p.t.sub(a, c)
	p.x.square(p.z)
	p.z.addRaw(p.x, p.x)
	exponentBias = word(4)
	a.subXBias(p.z, p.t, exponentBias)
	p.x.mul(a, b)
	p.z.mul(p.t, a)
	p.y.mul(p.t, d)
	if !beforeDouble {
		p.t.mul(b, d)
	}
	return p
}

// TODO: this will panic if byte array is not 56
func (p *twExtendedPoint) decafEncode(dst []byte) {
	t := word(0x00)
	overT := word(0x00)
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

func decafDecode(dst *twExtendedPoint, src serialized, useIdentity bool) word {
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

	return succ
}

func (p *twExtendedPoint) addNielsToExtended(p2 *twNiels, beforeDouble bool) {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.sub(p.y, p.x)
	a.mul(p2.a, b)
	b.addRaw(p.x, p.y)
	p.y.mul(p2.b, b)
	p.x.mul(p2.c, p.t)
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

func (p *twExtendedPoint) subNielsFromExtendedPoint(p2 *twNiels, beforeDouble bool) {
	a, b, c := &bigNumber{}, &bigNumber{}, &bigNumber{}
	b.sub(p.y, p.x)
	a.mul(p2.b, b)
	b.addRaw(p.x, p.y)
	p.y.mul(p2.a, b)
	p.x.mul(p2.c, p.t)
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

func (p *twExtendedPoint) addProjectiveNielsToExtended(pn *twPNiels, beforeDouble bool) {
	tmp := &bigNumber{}
	tmp.mul(p.z, pn.z)
	p.z = tmp.copy()
	p.addNielsToExtended(pn.n, beforeDouble)
}

func (p *twExtendedPoint) subProjectiveNielsFromExtendedPoint(p2 *twPNiels, beforeDouble bool) {
	tmp := &bigNumber{}
	tmp.mul(p.z, p2.z)
	p.z = tmp.copy()
	p.subNielsFromExtendedPoint(p2.n, beforeDouble)
}

func (p *twExtendedPoint) nielsToExtended(src *twNiels) {
	p.y.add(src.b, src.a)
	p.x.sub(src.b, src.a)
	p.t.mul(p.y, p.x)
	copy(p.z[:], bigOne[:])
}

func (p *twExtendedPoint) twPNiels() *twPNiels {
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

func pointScalarMul(pointA *twExtendedPoint, scalar *decafScalar) *twExtendedPoint {
	const decafWindowBits = 5            //move this to const file
	const window = decafWindowBits       //5
	const windowMask = (1 << window) - 1 //0x0001f 31
	const windowTMask = windowMask >> 1  //0x0000f 15
	const nTable = 1 << (window - 1)     //0x00010 16

	out := &twExtendedPoint{}

	scalar1x := &decafScalar{}
	scalar1x.scalarAdd(scalar, decafPrecompTable.scalarAdjustment)
	scalar1x.halve(scalar1x, scalarQ)

	multiples := pointA.prepareFixedWindow(nTable)

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
		pNeg := constTimeLookup(multiples, uint32(bits&windowTMask))
		pNeg.n.conditionalNegate(inv)

		if first {
			out = pNeg.twExtendedPoint()
			first = false
		} else {
			//Using Hisil et al's lookahead method instead of
			//extensible here for no particular reason.  Double
			//5 (window) times, but only compute out.t on the last one.
			for j := 0; j < window-1; j++ {
				out.double(true)
			}
			out.double(false)
			out.addProjectiveNielsToExtended(pNeg, false)
		}
	}
	return out
}

func precomputedScalarMul(scalar *decafScalar) *twExtendedPoint {
	p := &twExtendedPoint{
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
	}
	scalar2 := NewDecafScalar([fieldBytes]byte{})
	scalar2.Add(scalar, decafPrecompTable.scalarAdjustment)
	scalar2.halve(scalar2, scalarQ)

	var ni *twNiels
	for i := int(decafCombSpacing - 1); i >= 0; i-- {
		if i != int(decafCombSpacing-1) {
			p.double(false)
		}

		for j := uintZero; j < decafCombNumber; j++ {
			var tab word
			for k := uintZero; k < decafCombTeeth; k++ {
				bit := uint(i) + decafCombSpacing*(k+j*decafCombTeeth)
				if bit < scalarBits {
					tab |= (scalar2.(*decafScalar)[bit/wordBits] >> (bit % wordBits) & 1) << k
				}
			}

			invert := (sword(tab) >> (decafCombTeeth - 1)) - 1
			tab ^= word(invert)
			tab &= (1 << (decafCombTeeth - 1)) - 1

			index := uint32(((j << (decafCombTeeth - 1)) + uint(tab)))
			ni = decafPrecompTable.lookup(index)

			ni.conditionalNegate(word(invert))

			if i != int(decafCombSpacing-1) || j != 0 {
				p.addNielsToExtended(ni, j == decafCombNumber-1 && i != 0)
			} else {
				p.nielsToExtended(ni)
			}
		}
	}

	return p
}

func doubleScalarMul(pointB, pointC *twExtendedPoint, scalarB, scalarC *decafScalar) *twExtendedPoint {
	const decafWindowBits = 5
	const window = decafWindowBits       //5
	const windowMask = (1 << window) - 1 //0x0001f 31
	const windowTMask = windowMask >> 1  //0x0000f 15
	const nTable = 1 << (window - 1)     //0x00010 16

	scalar1x := &decafScalar{}
	scalar1x.scalarAdd(scalarB, decafPrecompTable.scalarAdjustment)
	scalar1x.halve(scalar1x, scalarQ)
	scalar2x := &decafScalar{}
	scalar2x.scalarAdd(scalarC, decafPrecompTable.scalarAdjustment)
	scalar2x.halve(scalar2x, scalarQ)

	multiples1 := pointB.prepareFixedWindow(nTable)
	multiples2 := pointC.prepareFixedWindow(nTable)
	out := &twExtendedPoint{}
	first := true
	for i := scalarBits - ((scalarBits - 1) % window) - 1; i >= 0; i -= window {
		bits1 := scalar1x[i/wordBits] >> uint(i%wordBits)
		bits2 := scalar2x[i/wordBits] >> uint(i%wordBits)
		if i%wordBits >= wordBits-window && i/wordBits < scalarWords-1 {
			bits1 ^= scalar1x[i/wordBits+1] << uint(wordBits-(i%wordBits))
			bits2 ^= scalar2x[i/wordBits+1] << uint(wordBits-(i%wordBits))
		}
		bits1 &= windowMask
		bits2 &= windowMask
		inv1 := (bits1 >> (window - 1)) - 1
		inv2 := (bits2 >> (window - 1)) - 1
		bits1 ^= inv1
		bits2 ^= inv2
		//Add in from table.  Compute t only on last iteration.
		mul1pn := constTimeLookup(multiples1, uint32(bits1&windowTMask))
		mul1pn.n.conditionalNegate(inv1)
		if first {
			out = mul1pn.twExtendedPoint()
			first = false
		} else {
			//Using Hisil et al's lookahead method instead of extensible here
			//for no particular reason.  Double WINDOW times, but only compute t on
			//the last one.
			for j := 0; j < window-1; j++ {
				out.double(true)
			}
			out.double(false)
			out.addProjectiveNielsToExtended(mul1pn, false)
		}
		mul2pn := constTimeLookup(multiples2, uint32(bits2&windowTMask))
		mul2pn.n.conditionalNegate(inv2)
		if i > 0 {
			out.addProjectiveNielsToExtended(mul2pn, true)
		} else {
			out.addProjectiveNielsToExtended(mul2pn, false)
		}
	}
	return out
}

// exposed methods

// NewPoint returns an Ed448 Point from uint32 arrays
func NewPoint(a [limbs]uint32, b [limbs]uint32, c [limbs]uint32, d [limbs]uint32) Point {
	x, y, z, t := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	for i := 0; i < limbs; i++ {
		x[i] = word(a[i])
		y[i] = word(b[i])
		z[i] = word(c[i])
		t[i] = word(d[i])
	}

	return &twExtendedPoint{x, y, z, t}
}

//IsValid tests if a point is valid.
func (p *twExtendedPoint) IsValid() bool {
	return p.isValidPoint()
}

//Equals compares whether two points are equal.
func (p *twExtendedPoint) Equals(q Point) bool {
	valid := p.equals(q.(*twExtendedPoint))
	return valid == decafTrue
}

//Copy copies a point.
func (p *twExtendedPoint) Copy() Point {
	p.copy()
	return Point(p)
}

//Add adds two points to produce a thrid point.
func (p *twExtendedPoint) Add(q, r Point) {
	p.add(q.(*twExtendedPoint), r.(*twExtendedPoint))
}

//Sub subtracts two points to produce a thrid point.
func (p *twExtendedPoint) Sub(q, r Point) {
	p.sub(q.(*twExtendedPoint), r.(*twExtendedPoint))
}

//Encode encodes a point as a sequence of bytes.
func (p *twExtendedPoint) Encode() []byte {
	out := make([]byte, 56)
	p.decafEncode(out)
	return out
}

//Decode decodes a point from a sequence of bytes.
//Every point has a unique encoding, so not every
//sequence of bytes is a valid encoding.  If an invalid
//encoding is given, the output is undefined.
func (p *twExtendedPoint) Decode(src []byte, useIdentity bool) {
	ser := [fieldBytes]byte{}
	copy(ser[:], src[:])
	decafDecode(p, ser, useIdentity)
}

//PointScalarMul multiplies a base point by a scalar.
func PointScalarMul(q Point, a Scalar) Point {
	return pointScalarMul(q.(*twExtendedPoint), a.(*decafScalar))
}

//PrecomputedScalarMul mutiplies a precomputed point by a scalar.
func PrecomputedScalarMul(s Scalar) Point {
	return precomputedScalarMul(s.(*decafScalar))
}

//DoubleScalarMul multiplies two base points by two scalars.
func DoubleScalarMul(q, r Point, a, b Scalar) Point {
	return doubleScalarMul(q.(*twExtendedPoint), r.(*twExtendedPoint), a.(*decafScalar), b.(*decafScalar))
}

//DoubleScalarMulNonsecret multiplies two base points by
//two scalars. It may leak the scalars. Otherwise is
//equivalent to DoubleScalarMul.
func DoubleScalarMulNonsecret(s1, s2 Scalar, b2 Point) Point {
	combo := &twExtendedPoint{
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
	}
	return decafDoubleNonSecretScalarMul(combo, b2.(*twExtendedPoint), s1.(*decafScalar), s2.(*decafScalar))
}
