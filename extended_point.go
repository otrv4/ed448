package ed448

type twExtendedPoint struct {
	x, y, z, t *bigNumber
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
	return decafEq(a, b)
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

/**
 * Decode a point from a sequence of bytes.
 *
 * Every point has a unique encoding, so not every
 * sequence of bytes is a valid encoding.  If an invalid
 * encoding is given, the output is undefined.
 *
 * @param [in] ser The serialized version of the point.
 * @param [in] allow_identity DECAF_TRUE (-1) if the identity is a legal input.
 *
 * Returns:
 * out - The decoded point.
 * ok  - Whether the decoding succeeded
 *       Success: uint64(-1)
 *       Failure: uint64( 0) if the base does not represent a point
 */
// XXX: name this dst?
func decafDecode(p *twExtendedPoint, ser serialized, identity word) word {
	a, b, c, d, e := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	n, succ := deserializeReturnMask(ser)
	ok := succ // ommit this
	zero := decafEq(n, bigZero)
	ok &= identity | ^zero
	ok &= ^highBit(n)
	a.square(n)
	p.z.sub(bigOne, a)
	b.square(p.z)
	c.mulWSignedCurveConstant(a, 4-4*(edwardsD))
	c.add(c, b)
	b.mul(c, a)
	d.isr(b)
	e.square(d)
	a.mul(e, b)
	a.add(a, bigOne)
	ok &= ^decafEq(a, bigZero)
	b.mul(c, d)
	d.decafCondNegate(highBit(b))
	p.x.add(n, n)
	c.mul(d, n)
	b.sub(bigTwo, p.z)
	a.mul(b, c)
	p.y.mul(a, p.z)
	p.t.mul(p.x, a)
	p.y[0] -= zero

	return ok
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

func (c *curveT) precomputedScalarMul(scalar *decafScalar) *twExtendedPoint {
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

			ni = decafPrecompTable.lookup(j, decafCombTeeth, uint(tab))

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

func (p *twExtendedPoint) DoubleScalarMul(s1 Scalar, p1 Point, s2 Scalar, p2 Point) {
	p = doubleScalarMul(p1.(*twExtendedPoint), s1, p2.(*twExtendedPoint), s2).copy()
}

func doubleScalarMul(
	pointB *twExtendedPoint, scalarB Scalar,
	pointC *twExtendedPoint, scalarC Scalar,
) *twExtendedPoint {
	const decafWindowBits = 5
	const window = decafWindowBits       //5
	const windowMask = (1 << window) - 1 //0x0001f 31
	const windowTMask = windowMask >> 1  //0x0000f 15
	const nTable = 1 << (window - 1)     //0x00010 16

	scalar1x := NewDecafScalar([fieldBytes]byte{})
	scalar1x.Add(scalarB, decafPrecompTable.scalarAdjustment)
	scalar1x.halve(scalar1x, scalarQ)
	scalar2x := NewDecafScalar([fieldBytes]byte{})
	scalar2x.Add(scalarC, decafPrecompTable.scalarAdjustment)
	scalar2x.halve(scalar2x, scalarQ)

	multiples1 := pointB.prepareFixedWindow(nTable)
	multiples2 := pointC.prepareFixedWindow(nTable)
	out := &twExtendedPoint{}
	first := true
	for i := scalarBits - ((scalarBits - 1) % window) - 1; i >= 0; i -= window {
		bits1 := scalar1x.(*decafScalar)[i/wordBits] >> uint(i%wordBits)
		bits2 := scalar2x.(*decafScalar)[i/wordBits] >> uint(i%wordBits)
		if i%wordBits >= wordBits-window && i/wordBits < scalarWords-1 {
			bits1 ^= scalar1x.(*decafScalar)[i/wordBits+1] << uint(wordBits-(i%wordBits))
			bits2 ^= scalar2x.(*decafScalar)[i/wordBits+1] << uint(wordBits-(i%wordBits))
		}
		bits1 &= windowMask
		bits2 &= windowMask
		inv1 := (bits1 >> (window - 1)) - 1
		inv2 := (bits2 >> (window - 1)) - 1
		bits1 ^= inv1
		bits2 ^= inv2
		/* Add in from table.  Compute t only on last iteration. */
		mul1pn := constTimeLookup(multiples1, bits1&windowTMask).copy()
		mul1pn.n.conditionalNegate(inv1)
		if first {
			out = mul1pn.twExtendedPoint()
			first = false
		} else {
			/* Using Hisil et al's lookahead method instead of extensible here
			 * for no particular reason.  Double WINDOW times, but only compute t on
			 * the last one.
			 */
			for j := 0; j < window-1; j++ {
				out.double(true)
			}
			out.double(false)
			out.addProjectiveNielsToExtended(mul1pn, false)
		}
		mul2pn := constTimeLookup(multiples2, bits2&windowTMask).copy()
		mul2pn.n.conditionalNegate(inv2)
		if i > 0 {
			out.addProjectiveNielsToExtended(mul2pn, true)
		} else {
			out.addProjectiveNielsToExtended(mul2pn, false)
		}
	}
	return out
}
