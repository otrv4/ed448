package ed448

func mask(a, b *bigNumber, mask word) {
	a[0] = word(mask) & b[0]
	a[1] = word(mask) & b[1]
	a[2] = word(mask) & b[2]
	a[3] = word(mask) & b[3]
	a[4] = word(mask) & b[4]
	a[5] = word(mask) & b[5]
	a[6] = word(mask) & b[6]
	a[7] = word(mask) & b[7]
	a[8] = word(mask) & b[8]
	a[9] = word(mask) & b[9]
	a[10] = word(mask) & b[10]
	a[11] = word(mask) & b[11]
	a[12] = word(mask) & b[12]
	a[13] = word(mask) & b[13]
	a[14] = word(mask) & b[14]
	a[15] = word(mask) & b[15]
}

// mask needs to be either decafTrue or decafFalse
func constantTimeSelectBigNumber(bfalse, btrue *bigNumber, mask word) *bigNumber {
	var x bigNumber

	x[0] = (btrue[0] & mask) | (bfalse[0] &^ mask)
	x[1] = (btrue[1] & mask) | (bfalse[1] &^ mask)
	x[2] = (btrue[2] & mask) | (bfalse[2] &^ mask)
	x[3] = (btrue[3] & mask) | (bfalse[3] &^ mask)
	x[4] = (btrue[4] & mask) | (bfalse[4] &^ mask)
	x[5] = (btrue[5] & mask) | (bfalse[5] &^ mask)
	x[6] = (btrue[6] & mask) | (bfalse[6] &^ mask)
	x[7] = (btrue[7] & mask) | (bfalse[7] &^ mask)
	x[8] = (btrue[8] & mask) | (bfalse[8] &^ mask)
	x[9] = (btrue[9] & mask) | (bfalse[9] &^ mask)
	x[10] = (btrue[10] & mask) | (bfalse[10] &^ mask)
	x[11] = (btrue[11] & mask) | (bfalse[11] &^ mask)
	x[12] = (btrue[12] & mask) | (bfalse[12] &^ mask)
	x[13] = (btrue[13] & mask) | (bfalse[13] &^ mask)
	x[14] = (btrue[14] & mask) | (bfalse[14] &^ mask)
	x[15] = (btrue[15] & mask) | (bfalse[15] &^ mask)

	return &x
}

// mask needs to be either decafTrue or decafFalse
func constantTimeSelectPoint(bfalse, btrue *twExtendedPoint, mask word) *twExtendedPoint {
	ret := &twExtendedPoint{}

	ret.x = constantTimeSelectBigNumber(bfalse.x, btrue.x, mask)
	ret.y = constantTimeSelectBigNumber(bfalse.y, btrue.y, mask)
	ret.z = constantTimeSelectBigNumber(bfalse.z, btrue.z, mask)
	ret.t = constantTimeSelectBigNumber(bfalse.t, btrue.t, mask)

	return ret
}

// mask needs to be either decafTrue or decafFalse
func constantTimeSelectScalar(bfalse, btrue *scalar, mask word) *scalar {
	var ret scalar

	ret[0] = (btrue[0] & mask) | (bfalse[0] &^ mask)
	ret[1] = (btrue[1] & mask) | (bfalse[1] &^ mask)
	ret[2] = (btrue[2] & mask) | (bfalse[2] &^ mask)
	ret[3] = (btrue[3] & mask) | (bfalse[3] &^ mask)
	ret[4] = (btrue[4] & mask) | (bfalse[4] &^ mask)
	ret[5] = (btrue[5] & mask) | (bfalse[5] &^ mask)
	ret[6] = (btrue[6] & mask) | (bfalse[6] &^ mask)
	ret[7] = (btrue[7] & mask) | (bfalse[7] &^ mask)
	ret[8] = (btrue[8] & mask) | (bfalse[8] &^ mask)
	ret[9] = (btrue[9] & mask) | (bfalse[9] &^ mask)
	ret[10] = (btrue[10] & mask) | (bfalse[10] &^ mask)
	ret[11] = (btrue[11] & mask) | (bfalse[11] &^ mask)
	ret[12] = (btrue[12] & mask) | (bfalse[12] &^ mask)
	ret[13] = (btrue[13] & mask) | (bfalse[13] &^ mask)

	return &ret
}

// ConstantTimeSelectPoint will use constant time select to choose either the left or right point, depending
// on if the mask is either all zeroes, or all ones
func ConstantTimeSelectPoint(bfalse, btrue Point, mask uint32) Point {
	return constantTimeSelectPoint(bfalse.(*twExtendedPoint), btrue.(*twExtendedPoint), word(mask))

}

// ConstantTimeSelectScalar will use constant time select to choose either the left or right scalar, depending
// on if the mask is either all zeroes, or all ones
func ConstantTimeSelectScalar(bfalse, btrue Scalar, mask uint32) Scalar {
	return constantTimeSelectScalar(bfalse.(*scalar), btrue.(*scalar), word(mask))
}
