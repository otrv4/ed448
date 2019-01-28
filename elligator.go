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
	r, a, b, c, n, e := &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}, &bigNumber{}

	p := &twExtendedPoint{
		x: new(bigNumber),
		y: new(bigNumber),
		z: new(bigNumber),
		t: new(bigNumber),
	}

	// probable nonresidue
	r0, _ := deserialize(ser)
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

	// TODO: do this
	// assert(API_NS(point_valid)(p));

	return p
}

func pointFromUniformHash(ser [112]byte) *twExtendedPoint {
	var ser1 [56]byte
	var ser2 [56]byte

	copy(ser1[:], ser[:56])
	copy(ser2[:], ser[:56])

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
