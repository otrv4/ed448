package ed448

type twExtendedPoint struct {
	x, y, z, t *bigNumber
}

// Based on Hisil's formula 5.1.3: Doubling in E^e
// XXX: Find out if double is always a double of itself
func (p *twExtendedPoint) double(q *twExtendedPoint, beforeDouble bool) {
	a, b, c, d := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	c.square(q.x)
	a.square(q.y)
	d.addRaw(c, a)
	p.t.addRaw(q.y, q.x)
	b.square(p.t)
	exponentBias := uint32(3)
	b.subXBias(b, d, exponentBias)
	p.t.sub(a, c)
	p.x.square(q.z)
	p.z.addRaw(p.x, p.x)
	exponentBias = uint32(4)
	a.subXBias(p.z, p.t, exponentBias)
	p.x.mul(a, b)
	p.z.mul(p.t, a)
	p.y.mul(p.t, d)
	if !beforeDouble {
		p.t.mul(b, d)
	}
}

func (p *twExtendedPoint) decafEncode(dst []byte) {
	t := uint64(0)
	overT := uint64(0)
	serialize(dst, p.deisogenize(t, overT))
}

func (p *twExtendedPoint) deisogenize(t, overT uint64) *bigNumber {
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

func decafDecode(ser serialized, identity uint64) (*twExtendedPoint, uint64) {
	a, b, c, d, e := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	p := &twExtendedPoint{
		x: &bigNumber{},
		y: &bigNumber{},
		z: &bigNumber{},
		t: &bigNumber{},
	}

	n, succ := deserializeReturnMask(ser)
	ok := uint64(succ)

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
	p.y[0] -= uint32(zero)

	return p, ok
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

func (p *twExtendedPoint) add(pn *twPNiels, beforeDouble bool) {
	tmp := &bigNumber{}
	tmp.mul(p.z, pn.z)
	p.z.copyInto(tmp)
	p.addNielsToExtended(pn.n, beforeDouble)
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

func (c *curveT) precomputedScalarMul(scalar Scalar) *twExtendedPoint {

	p := &twExtendedPoint{
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
	}

	scalar2 := scalarAdd(scalar, decafPrecompTable.scalarAdjustment)
	scalar2 = scalarHalve(scalar2, scalarQ)

	var ni *twNiels
	for i := int(decafCombSpacing - 1); i >= 0; i-- {
		if i != int(decafCombSpacing-1) {
			p.double(p, false)
		}

		for j := uintZero; j < decafCombNumber; j++ {
			var tab uint32
			for k := uintZero; k < decafCombTeeth; k++ {
				bit := uint(i) + decafCombSpacing*(k+j*decafCombTeeth)
				if bit < scalarBits {
					tab |= (scalar2[bit/wordBits] >> (bit % wordBits) & 1) << k
				}
			}

			invert := (int32(tab) >> (decafCombTeeth - 1)) - 1
			tab ^= uint32(invert)
			tab &= (1 << (decafCombTeeth - 1)) - 1

			ni = decafPrecompTable.lookup(j, decafCombTeeth, uint(tab))

			ni.conditionalNegate(uint32(invert))

			if i != int(decafCombSpacing-1) || j != 0 {
				p.addNielsToExtended(ni, j == decafCombNumber-1 && i != 0)
			} else {
				p.nielsToExtended(ni)
			}
		}
	}

	return p
}

//TODO implement subfunctions
func pointDoubleScalarMul(
	pointA *twExtendedPoint, scalarX [scalarWords]uint32,
	pointB *twExtendedPoint, scalarY [scalarWords]uint32,
) (out *twExtendedPoint) {
	// 	const decafWindowBits = 5
	// 	const window = decafWindowBits       //5
	// 	const windowMask = (1 << window) - 1 //0x0001f 31
	// 	const windowTMask = windowMask >> 1  //0x0000f 15
	// 	const nTable = 1 << (window - 1)     //0x00010 16

	// 	scalar1x := [scalarWords]uint32{}
	// 	scalar2x := [scalarWords]uint32{}

	// 	scalar1x = scalarAdd(scalarY, decafPrecompTable.scalarAdjustment)
	// 	scalar1x = scalarHalve(scalar1x, scalarQ)
	// 	scalar2x = scalarAdd(scalarY, decafPrecompTable.scalarAdjustment)
	// 	scalar2x = scalarHalve(scalar2x, scalarQ)

	// 	/* Set up a precomputed table with odd multiples of scalarY. */
	// 	point := twPNiels{}
	// 	tmp := twExtendedPoint{}
	// 	multiples1 := [nTable]twPNiels{}
	// 	multiples2 := [nTable]twPNiels{}
	// 	prepareFixedWindow(multiples1, scalarX, nTable) //implement tables
	// 	prepareFixedWindow(multiples2, scalarY, nTable)

	// 	/* Initialize. */
	// 	i := scalarBits - ((scalarBits - 1) % window) - 1

	// 	for (i; i>=0, i-=window) {
	// 		/* Fetch another block of bits */
	// 		bits1 := new(bigNumber)
	// 		bits2 := new(bigNumber)
	// 		bits1 = scalar1x[i/wordBits] >> (i%wordBits)
	// 		bits2 = scalar2x[i/wordBits] >> (i%wordBits)
	// 		if (i%wordBits >= wordBits-window && i/wordBits < scalarWords-1) {
	// 			bits1 ^= scalar1x[i/wordBits+1] << (wordBits - (i%wordBits))
	// 			bits2 ^= scalar2x[i/wordBits+1] << (wordBits - (i%wordBits))
	// 		}
	// 		bits1 &= windowMask
	// 		bits2 &= windowMask
	// 		inv1 := new(bigNumber)
	// 		inv2 := new(bigNumber)
	// 		inv1 = (bits1 >> (window-1)) -1
	// 		inv2 = (bits2 >> (window-1)) -1
	// 		bits1 ^= inv1
	// 		bits2 ^= inv2

	// 		/* Add in from table.  Compute t only on last iteration. */
	// 		decafPrecompTable.lookupxx(point, multiples1, size(point), nTable, bits1 & windowTMask)	//implement lookupxx
	// 		point.n.conditionalNegate(inv1)
	// 		v := 1
	// 		if (v !=0) {
	// 			tmp.twExtendedPoint(point)
	// 			v = 0
	// 		} else {
	// 			/* Using Hisil et al's lookahead method instead of extensible here
	//             * for no particular reason.  Double WINDOW times, but only compute t on
	//             * the last one.
	//             */
	//             for (j=0; j<window-1; j++) {
	//             	pointDoubleInternal(tmp, tmp, -1)
	//             }
	//             pointDoubleInternal(tmp, tmp, 0)
	//             tmp.addNielsToExtended(point, false)
	// 		}

	// 		decafPrecompTable.lookupxx(point, multiples2, size(point), nTable, bits2 & windowTMask)	//implement lookupxx
	// 		point.n.conditionalNegate(inv2)
	// 		tmp.addNielsToExtended(point, i!=0 ? -1 : 0)

	// 	}
	// 	copyInto(out, tmp)
	return out
}
