package ed448

const (
	Limbs     = 16
	Radix     = 28
	radixMask = limb(0xfffffff)
)

func deserialize(in serialized) (n *bigNumber, ok bool) {
	n = &bigNumber{}

	for i := uint(0); i < 8; i++ {
		out := uint64(0)
		for j := uint(0); j < 7; j++ {
			out |= uint64(in[7*i+j]) << (8 * j)
		}

		n[2*i] = limb(out) & radixMask
		n[2*i+1] = limb(out >> 28)
	}

	ok = !constantTimeGreaterOrEqualP(n)
	return
}

func serialize(dst []byte, n *bigNumber) {
	src := n.copy()
	src.strongReduce()

	for i := 0; i < 8; i++ {
		l := uint64(src[2*i]) + uint64(src[2*i+1])<<28
		for j := 0; j < 7; j++ {
			dst[7*i+j] = byte(l)
			l >>= 8
		}
	}
}

func (n *bigNumber) copy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

//TODO: double check if this can be used for both 32 and 64 bits
//(at least before unrolling)
func (n *bigNumber) strongReduce() {
	// clear high
	n[8] += n[15] >> 28
	n[0] += n[15] >> 28
	n[15] &= radixMask

	scarry := int64(0)
	for i := 0; i < 16; i++ {
		m := uint32(radixMask)
		if i == 8 {
			m = uint32(0xffffffe)
		}

		scarry += int64(n[i]) - int64(m)

		n[i] = limb(scarry) & radixMask
		scarry >>= 28
	}

	scarryMask := Word(scarry) & Word(radixMask)
	carry := uint64(0)
	for i := 0; i < 16; i++ {
		m := uint64(scarryMask)
		if i == 8 {
			m &= uint64(0xfffffffffffffffe)
		}

		carry += uint64(n[i]) + m
		n[i] = limb(carry) & radixMask
		carry >>= 28
	}
}
