package ed448

// y^2/x^2
// needed: 4*v*(u^2 - 1)/(u^4 - 2*u^2 + 4*v^2 + 1)
func (p *twExtendedPoint) x448LikeEncode(dst []byte) {
	if len(dst) != x448FieldBytes {
		panic("Attempted to encode with a destination that is not 56 bytes")
	}

	q := p.copy()
	q.t = invert(p.x) // 1/x
	q.z.mul(q.t, q.y) // y/x
	q.y.square(q.z)   // (y/x)^2

	dsaLikeSerialize(dst[:], q.y)

	// wipe out
	q.x.set(bigZero)
	q.y.set(bigZero)
	q.z.set(bigZero)
	q.t.set(bigZero)
}

func fromEdDSATox448(ed []byte) [x448FieldBytes]byte {
	if len(ed) != dsaFieldBytes {
		panic("Attempted to convert an array that is not 57 bytes")
	}

	y, n, d := &bigNumber{}, &bigNumber{}, &bigNumber{}
	mask := uint(0xfe << 7)

	dsaLikeDeserialize(y, ed[:], mask)

	// u = y^2 * (1-dy^2) / (1-y^2)
	n.square(y)                            // y^2
	d.sub(bigOne, n)                       // (1-y^2)
	d = invert(d)                          // 1 / (1-y^2)
	y.mul(n, d)                            // y^2 / (1-y^2)
	d.mulWSignedCurveConstant(n, edwardsD) // dy^2
	d.sub(bigOne, d)                       // 1 - dy^2
	n.mul(y, d)                            // y^2 * (1-dy^2) / (q-y^2)

	var dst [x448FieldBytes]byte
	dsaLikeSerialize(dst[:], n)

	// wipe out
	y.set(bigZero)
	n.set(bigZero)
	d.set(bigZero)

	return dst
}

func x448BasePointScalarMul(s []byte) [x448FieldBytes]byte {
	if len(s) != x448FieldBytes {
		panic("Wrong scalar length: should be 56 bytes")
	}

	scalar2 := append([]byte{}, s...)
	// Scalar conditioning
	scalar2[0] &= -(byte(Cofactor))

	theScalar := &scalar{}

	scalar2[x448FieldBytes-1] &= ^(-1 << ((x448FieldBytes + 7) % 8))
	scalar2[x448FieldBytes-1] |= 1 << ((x448FieldBytes + 7) % 8)

	theScalar.decode(scalar2)

	for i := uint(1); i < 2; i <<= 1 {
		theScalar.halve(theScalar)
	}

	p := precomputedScalarMul(theScalar)

	var out [x448FieldBytes]byte
	p.x448LikeEncode(out[:])

	// wipe out
	p.y.set(bigZero)
	p.x.set(bigZero)
	p.z.set(bigZero)
	p.t.set(bigZero)

	return out
}

func x448ScalarMul(point []byte, s []byte) ([x448FieldBytes]byte, bool) {
	if len(s) != x448FieldBytes || len(point) != x448FieldBytes {
		panic("Wrong scalar or base length: should be 56 bytes")
	}
	x1, t1, t2 := &bigNumber{}, &bigNumber{}, &bigNumber{}

	swap := word(0)

	dsaLikeDeserialize(x1, point, uint(0))
	x2 := bigOne.copy()
	z2 := bigZero.copy()
	x3 := x1.copy()
	z3 := bigOne.copy()

	for t := int(x448FieldBits - 1); t >= 0; t-- {
		sb := byte(s[t/8])
		var kT word

		// Scalar conditioning
		if t/8 == 0 {
			sb &= -byte(Cofactor)
		} else if t == (x448FieldBits - 1) {
			sb = -byte(byteOne)
		}

		kT = word((sb >> byte(t%8)) & 1)
		kT = -kT // set to all 0s or all 1s

		swap ^= kT
		x2.conditionalSwap(x3, swap)
		z2.conditionalSwap(z3, swap)
		swap = kT

		t1.addRaw(x2, z2) // A = x2 + z2 // 2+e
		t2.sub(x2, z2)    // B = x2 - z2 // 3+e
		z2.sub(x3, z3)    // D = x3 - z3 // 3+e
		x2.mul(t1, z2)    //DA
		z2.addRaw(z3, x3) // C = x3 + z3 // 2+e
		x3.mul(t2, z2)    // CB
		z3.sub(x2, x3)    // DA - CB
		z2.square(z3)     // (DA - CB)^2
		z3.mul(x1, z2)    // z3 = x1(DA-CB)^2
		z2.addRaw(x2, x3) // (DA + CB) // 2+e
		x3.square(z2)     // x3 = (DA+CB)^2

		z2.square(t1)  // AA = A^2
		t1.square(t2)  // BB = B^2
		x2.mul(z2, t1) // x2 = AA*BB
		t2.sub(z2, t1) //  E = AA-BB // 3+e

		t1.mulW(t2, -edwardsD) // E*-d = a24*E
		t1.addRaw(t1, z2)      // AA + a24*E // 2+e
		z2.mul(t2, t1)         // z2 = E(AA+a24*E)
	}

	x2.conditionalSwap(x3, swap)
	z2.conditionalSwap(z3, swap)
	z2 = invert(z2)
	x1.mul(x2, z2)

	var out [x448FieldBytes]byte
	dsaLikeSerialize(out[:], x1)

	ok := !(x1.equals(bigZero))

	// wipe out
	x1.set(bigZero)
	x2.set(bigZero)
	z2.set(bigZero)
	x3.set(bigZero)
	z3.set(bigZero)
	t1.set(bigZero)
	t2.set(bigZero)

	return out, ok
}
