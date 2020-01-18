package ed448

// This function runs Elligator2 on the decaf Jacobi quartic model.  It then
// uses the isogeny to put the result in twisted Edwards form.  As a result,
// it is safe (cannot produce points of order 4), and would be compatible with
// hypothetical other implementations of Decaf using a Montgomery or untwisted
// Edwards model.
// This function isn't quite indifferentiable from a random oracle.
// However, it is suitable for many protocols, including SPEKE and SPAKE2 EE.
// Furthermore, calling it twice with independent seeds and adding the results
// is indifferentiable from a random oracle.
func pointFromNonUniformHash(ser [56]byte) *twExtendedPoint {
	r0, r, a, b, c, n, e := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	p := &twExtendedPoint{
		x: new(bigNumber),
		y: new(bigNumber),
		z: new(bigNumber),
		t: new(bigNumber),
	}

	mask := uint(0xfe << 7)
	dsaLikeDeserialize(r0, ser[:], mask)
	r0.strongReduce()
	a.square(r0) //r^2
	r.sub(bigZero, a)

	// From Decaf paper
	// Compute D2 := (dr+a-d)(dr-ar-d) with a=1
	a.sub(r, bigOne)
	b.mulWSignedCurveConstant(a, edwardsD) // (d * r) - d
	a.add(b, bigOne)
	b.sub(b, r)
	c.mul(a, b)

	// compute N := (r+1)(a-2d)
	a.add(r, bigOne)
	n.mulWSignedCurveConstant(a, 1-2*edwardsD)

	// e = +-sqrt(1/ND) or +-r0 * sqrt(qnr/ND)
	a.mul(c, n)
	square := b.isr(a)
	c = constantTimeSelect(bigOne, r0, square) // r? = isSquare ? 1 : r0
	e.mul(b, c)

	// s2 = +- |N . e|
	a.mul(n, e)
	a.decafCondNegate(lowBit(a) ^ ^square) // NB

	// t2 = -+ cN(r - 1)((a - (2 * d))e)^ 2 - 1
	c.mulWSignedCurveConstant(e, 1-2*edwardsD) // ( a - (2 * d))e
	b.square(c)
	e.sub(r, bigOne)
	c.mul(b, e)
	b.mul(c, n)
	b.decafCondNegate(square)
	b.sub(b, bigOne)

	// isogenize
	c.square(a) // s^2
	a.add(a, a) // 2s
	e.add(c, bigOne)
	p.t.mul(a, e) // 2s(1+s^2)
	p.x.mul(a, b) // 2st
	a.sub(bigOne, c)
	p.y.mul(e, a) // (1+s^2)(1-s^2)
	p.z.mul(a, b) // (1-s^2)t

	if !p.isOnCurve() {
		return nil
	}

	return p
}

func pointFromUniformHash(ser [112]byte) *twExtendedPoint {
	var ser1 [56]byte
	var ser2 [56]byte

	copy(ser1[:], ser[:56])
	copy(ser2[:], ser[56:])

	p := pointFromNonUniformHash(ser1)
	q := pointFromNonUniformHash(ser2)

	r := &twExtendedPoint{
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
		&bigNumber{},
	}

	r.add(p, q)
	return r
}

func invertElligatorNonUniform(p *twExtendedPoint, hint word) ([56]byte, bool) {
	sgnS := word(-(hint & 1))
	sgnAltX := word(-((hint >> 1) & 1))
	sgnR0 := word(-((hint >> 2) & 1))

	a, b := &bigNumber{}, &bigNumber{}

	c := p.deisogenizeNew(a, b, sgnS, sgnAltX)

	isIdentity := p.t.decafEq(bigZero)
	a.decafConstTimeSel(a, bigOne, isIdentity&sgnAltX)
	b.decafConstTimeSel(b, bigOne, isIdentity&sgnS&(^sgnAltX))

	c.mulWSignedCurveConstant(a, edwardsD-1)
	a.add(c, a)
	c.sub(c, b)
	a.add(a, b)
	c.conditionalSwap(a, sgnS)
	b.sub(bigZero, a)
	a.mul(b, c)
	succ := b.isr(a)
	succ |= a.decafEq(bigZero)
	a.mul(b, c)
	a.decafCondNegate(sgnR0 ^ lowBit(a))
	// Eliminate duplicate values for identity
	succ &= ^(a.decafEq(bigZero)&sgnR0 | sgnS)

	var dst [56]byte

	dsaLikeSerialize(dst[:], a)
	return dst, true

	// TODO: check: recovered_hash[SER_BYTES-1] ^= (hint>>3)<<0;
	// return goldilocks_succeed_if(mask_to_bool(succ));
	//return dst, false
}

func invertElligatorUniform(src [112]byte, p *twExtendedPoint, hint word) ([112]byte, bool) {
	p2 := &twExtendedPoint{}
	var partialHash, partialHash2 [56]byte
	var ok bool
	var dst [112]byte

	copy(partialHash[:], src[56:])
	p2 = pointFromNonUniformHash(partialHash)

	p2.sub(p, p2)

	partialHash2, ok = invertElligatorNonUniform(p2, hint)
	if !ok {
		return dst, false
	}

	copy(dst[56:], partialHash[:])
	copy(dst[:56], partialHash2[:])

	return dst, true
}
