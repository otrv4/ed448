package ed448

type twExtendedPoint struct {
	x, y, z, t *bigNumber
}

// Based on Hisil's formula 5.1.3: Doubling in E^e
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
	t := dword_t(0)
	overT := dword_t(0)
	serialize(dst, p.deisogenize(t, overT))
}

func (p *twExtendedPoint) deisogenize(t, overT dword_t) *bigNumber {
	a, b, c, d, s := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	a.mulWSignedCurveConstant(p.y, 1-(D))
	c.mul(a, p.t)
	a.mul(p.x, p.z)
	d.sub(c, a)
	a.add(p.z, p.y)
	b.sub(p.z, p.y)
	c.mul(b, a)
	b.mulWSignedCurveConstant(c, (-(D)))
	a.isr(b)
	b.mulWSignedCurveConstant(a, (-(D)))
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

func decafDecode(ser serialized, identity dword_t) (*twExtendedPoint, dword_t) {
	a, b, c, d, e := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}
	p := &twExtendedPoint{
		x: &bigNumber{},
		y: &bigNumber{},
		z: &bigNumber{},
		t: &bigNumber{},
	}

	n, succ := deserializeReturnMask(ser)
	ok := dword_t(succ)

	zero := decafEq(n, bigZero)
	ok &= identity | ^zero
	ok &= ^highBit(n)
	a.square(n)
	p.z.sub(bigOne, a)
	b.square(p.z)
	c.mulWSignedCurveConstant(a, 4-4*(D))
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
	p.y[0] -= word_t(zero)

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

func (p *twExtendedPoint) nielsToExtended(src *twNiels) {
	p.y.add(src.b, src.a)
	p.x.sub(src.b, src.a)
	p.t.mul(p.y, p.x)
	copy(p.z[:], bigOne[:])
}

//func (p *twExtendedPoint) precomputedScalarMul(scalar [scalarWords]word_t) {
//	n := uint(5)
//	t := uint(5)
//	s := uint(18)
//
//	var scalar1 [scalarWords]word_t
//	scalar1 = scalarAdd(scalar, precomputedBaseTable.scalarAdjustment)
//
//	scalar1 = scHalve(scalar1, scP)
//
//	var ni *twNiels
//
//	for i := int(s - 1); i >= 0; i-- {
//		if i != int(s-1) {
//			p.pointDoubleInternal(p, false)
//		}
//
//		for j := uint(0); j < n; j++ {
//			var tab word_t
//			for k := uint(0); k < t; k++ {
//				bit := uint(i) + s*(k+j*t)
//				if bit < 446 { // change 446 to constant
//					tab |= (scalar1[bit/uint(32)] >> (bit % uint(32)) & 1) << k
//					// change uint(32) to constant
//				}
//			}
//
//			invert := (int32(tab) >> (t - 1)) - 1
//			tab ^= word_t(invert)
//			tab &= (1 << (t - 1)) - 1
//
//			ni = precomputedBaseTable.decafLookup(j, t, uint(tab))
//
//			ni.conditionalNegate(word_t(invert))
//
//			if i != int(s-1) || j != 0 {
//				p.addNielsToProjective(ni, j == n-1 && i != 0)
//			} else {
//				convertNielsToPt(p, ni)
//			}
//		}
//	}
//	//pointPrint("x end", p.x)
//	//pointPrint("y end", p.y)
//	//pointPrint("z end", p.z)
//	//pointPrint("t end", p.t)
//}
