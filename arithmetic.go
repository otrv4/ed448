package ed448

// ModQ produces a byte array mod Q (prime order)
func ModQ(serial []byte) []byte {
	words := [16]word_t{}
	deserializeModQ(words[:], serial)
	out := make([]byte, fieldBytes)
	wordsToBytes(out, words[:])
	return out
}

// Mul multiplies two large values
func Mul(x [fieldBytes]byte, y [fieldBytes]byte) (out [fieldBytes]byte) {
	desX, _ := deserialize(x)
	desY, _ := deserialize(y)
	desX.mulCopy(desX, desY)
	serialize(out[:], desX)
	return out
}

// Add two large values
func Add(x [fieldBytes]byte, y [fieldBytes]byte) (out [fieldBytes]byte) {
	desX, _ := deserialize(x)
	desY, _ := deserialize(y)
	desX.add(desX, desY)
	serialize(out[:], desX)
	return out
}

// Sub subtracts two large values
func Sub(x [fieldBytes]byte, y [fieldBytes]byte) (out [fieldBytes]byte) {
	desX, _ := deserialize(x)
	desY, _ := deserialize(y)
	desX.sub(desX, desY)
	serialize(out[:], desX)
	return out
}
