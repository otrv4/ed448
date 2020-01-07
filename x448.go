package ed448

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

	return dst
}
