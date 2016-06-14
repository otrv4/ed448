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

//XXX dst should have fieldBytes size
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

func (n *bigNumber) bias(b uint32) *bigNumber {
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

	return n
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

func (n *bigNumber) mulW(x *bigNumber, w uint64) *bigNumber {
	whi := uint32(w >> Radix)
	wlo := uint32(w & uint64(radixMask))

	var accum0, accum8 uint64

	accum0 = uint64(wlo) * uint64(x[0])
	accum8 = uint64(wlo) * uint64(x[8])
	accum0 += uint64(whi) * uint64(x[15])
	accum8 += uint64(whi) * uint64(x[15]+x[7])

	n[0] = limb(accum0 & uint64(radixMask))
	accum0 >>= Radix

	n[8] = limb(accum8 & uint64(radixMask))
	accum8 >>= Radix

	for i := 1; i < Limbs/2; i++ {
		accum0 += uint64(wlo) * uint64(x[i])
		accum8 += uint64(wlo) * uint64(x[i+8])
		accum0 += uint64(whi) * uint64(x[i-1])
		accum8 += uint64(whi) * uint64(x[i+7])

		n[i] = limb(accum0 & uint64(radixMask))
		accum0 >>= Radix

		n[i+8] = limb(accum8 & uint64(radixMask))
		accum8 >>= Radix
	}

	accum0 += accum8 + uint64(n[8])
	n[8] = limb(accum0 & uint64(radixMask))
	n[9] += limb(accum0 >> Radix)

	accum8 += uint64(n[0])
	n[0] = limb(accum8 & uint64(radixMask))
	n[1] += limb(accum8 >> Radix)

	return n
}
