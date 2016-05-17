package ed448

const (
	Limbs = 8
	Radix = 56
)

func deserialize(in serialized) (n bigNumber, ok bool) {
	const (
		columns = Limbs
		rows    = Limbs - 1
	)

	for i := uint(0); i < columns; i++ {
		for j := uint(0); j < rows; j++ {
			n[i] |= limb(in[rows*i+j]) << (columns * j)
		}
	}

	ok = !constantTimeGreaterOrEqualP(n)

	return
}

func serialize(dst []byte, src bigNumber) {
	const (
		rows    = Limbs
		columns = Radix / Limbs
	)

	var n bigNumber
	copy(n[:], src[:])

	for i := uint(0); i < rows; i++ {
		for j := uint(0); j < columns; j++ {
			dst[columns*i+j] = byte(n[i])
			n[i] >>= 8
		}
	}
}
