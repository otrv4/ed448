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

func (n *bigNumber) bias(b uint32) {
	var co1 limb = radixMask * limb(b)
	var co2 limb = co1 - limb(b)
	lo := [4]limb{co1, co1, co1, co1}
	hi := [4]limb{co2, co1, co1, co1}

	n[0] += lo[0]
	n[1] += lo[1]
	n[2] += lo[2]
	n[3] += lo[3]

	n[4] += lo[0]
	n[5] += lo[1]
	n[6] += lo[2]
	n[7] += lo[3]

	n[8] += hi[0]
	n[9] += hi[1]
	n[10] += hi[2]
	n[11] += hi[3]

	n[12] += lo[0]
	n[13] += lo[1]
	n[14] += lo[2]
	n[15] += lo[3]
}

//TODO: double check if this can be used for both 32 and 64 bits
//(at least before unrolling)
func (n *bigNumber) strongReduce() *bigNumber {
	// clear high
	n[8] += n[15] >> 28
	n[0] += n[15] >> 28
	n[15] &= radixMask

	scarry := int64(0)
	for i := 0; i < 16; i++ {
		m := limb(radixMask)
		if i == 8 {
			m = limb(0xffffffe)
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

	return n
}
