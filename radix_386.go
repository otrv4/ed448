package ed448

const (
	Limbs = 16
	Radix = 28
)

func deserialize(in serialized) (n bigNumber, ok bool) {
	mask := limb(0xfffffff)

	for i := uint(0); i < 8; i++ {
		out := uint64(0)
		for j := uint(0); j < 7; j++ {
			out |= uint64(in[7*i+j]) << (8 * j)
		}

		n[2*i] = limb(out) & mask
		n[2*i+1] = limb(out >> 28)
	}

	ok = !constantTimeGreaterOrEqualP(n)
	return
}

func serialize(dst []byte, src bigNumber) {
	//TODO strong reduce

	for i := 0; i < 8; i++ {
		l := uint64(src[2*i]) + uint64(src[2*i+1])<<28
		for j := 0; j < 7; j++ {
			dst[7*i+j] = byte(l)
			l >>= 8
		}
	}

	//TODO
}
