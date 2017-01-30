package ed448

// ModQ produces a byte array mod Q (prime order)
func ModQ(serial []byte) []byte {
	words := Scalar{}
	words.deserializeModQ(serial)
	out := make([]byte, fieldBytes)
	words.serialize(out)
	return out
}

// PointMul multiplies a point x with a scalar y
// PointMul automatically reduces the output by P
func PointMul(x [fieldBytes]byte, y [fieldBytes]byte) (out []byte) {
	desX, okX := deserialize(x)
	desY, okY := deserialize(y)
	if !(okX && okY) {
		return nil
	}
	desZ := &bigNumber{}
	desZ.mulCopy(desX, desY)
	out = make([]byte, fieldBytes)
	serialize(out, desZ)
	return out
}

// PointAddition adds two Ed448 points
// Inputs should never be >= prime P. If they are, PointAddition returns nil.
// PointAddition automatically reduces the output by P
func PointAddition(x [fieldBytes]byte, y [fieldBytes]byte) (out []byte) {
	desX, okX := deserialize(x)
	desY, okY := deserialize(y)
	if !(okX && okY) {
		return nil
	}
	desZ := &bigNumber{}
	desZ.add(desX, desY)
	out = make([]byte, fieldBytes)
	serialize(out, desZ)
	return out
}

// ScalarSub subtracts scalar x from scalar y.
// ScalarSub automatically reduces the output by Q
func ScalarSub(x Scalar, y Scalar) (out Scalar) {
	noExtra := uint32(0)
	out.scalarSubExtra(x, y, noExtra)
	return
}

// ScalarMul multiplies scalar x from scalar y.
// ScalarMul automatically reduces the output by Q
func ScalarMul(x Scalar, y Scalar) (out Scalar) {
	out.montgomeryMultiply(x, y)
	out.montgomeryMultiply(out, scalarR2)
	return
}
