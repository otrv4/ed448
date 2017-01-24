package ed448

// ModQ produces a byte array mod Q (prime order)
func ModQ(serial []byte) []byte {
	words := [16]uint32{}
	deserializeModQ(words[:], serial)
	out := make([]byte, fieldBytes)
	wordsToBytes(out, words[:])
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
	desX.mulCopy(desX, desY)
	out = make([]byte, fieldBytes)
	serialize(out, desX)
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
	desX.add(desX, desY)
	out = make([]byte, fieldBytes)
	serialize(out, desX)
	return out
}

// ScalarSub subtracts scalar x from scalar y.
// ScalarSub automatically reduces the output by Q
func ScalarSub(x [scalarWords]uint32, y [scalarWords]uint32) (out [scalarWords]uint32) {
	noExtra := uint32(0)
	return scalarSubExtra(x, y, noExtra)
}
