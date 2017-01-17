package ed448

const (
	Limbs     = 16
	Radix     = 28
	radixMask = word_t(0xfffffff)
)

func deserializeReturnMask(in serialized) (*bigNumber, word_t) {
	n := &bigNumber{}

	for i := uint(0); i < 8; i++ {
		out := uint64(0)
		for j := uint(0); j < 7; j++ {
			out |= uint64(in[7*i+j]) << (8 * j)
		}

		n[2*i] = word_t(out) & radixMask
		n[2*i+1] = word_t(out >> 28)
	}

	return n, constantTimeGreaterOrEqualP(n)
}

func deserialize(in serialized) (n *bigNumber, ok bool) {
	n, mask := deserializeReturnMask(in)
	ok = mask == 0xffffffff
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
	var co1 word_t = radixMask * word_t(b)
	var co2 word_t = co1 - word_t(b)
	lo := [4]word_t{co1, co1, co1, co1}
	hi := [4]word_t{co2, co1, co1, co1}

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

func (n *bigNumber) strongReduce() *bigNumber {
	// clear high
	n[8] += n[15] >> 28
	n[0] += n[15] >> 28
	n[15] &= radixMask

	//first for

	scarry := int64(0)
	scarry += int64(n[0]) - 0xfffffff
	n[0] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[1]) - 0xfffffff
	n[1] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[2]) - 0xfffffff
	n[2] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[3]) - 0xfffffff
	n[3] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[4]) - 0xfffffff
	n[4] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[5]) - 0xfffffff
	n[5] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[6]) - 0xfffffff
	n[6] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[7]) - 0xfffffff
	n[7] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[8]) - 0xffffffe
	n[8] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[9]) - 0xfffffff
	n[9] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[10]) - 0xfffffff
	n[10] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[11]) - 0xfffffff
	n[11] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[12]) - 0xfffffff
	n[12] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[13]) - 0xfffffff
	n[13] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[14]) - 0xfffffff
	n[14] = word_t(scarry) & radixMask
	scarry >>= 28

	scarry += int64(n[15]) - 0xfffffff
	n[15] = word_t(scarry) & radixMask
	scarry >>= 28

	// second for

	scarryMask := word_t(scarry) & word_t(radixMask)
	carry := uint64(0)
	m := uint64(scarryMask)

	carry += uint64(n[0]) + m
	n[0] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[1]) + m
	n[1] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[2]) + m
	n[2] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[3]) + m
	n[3] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[4]) + m
	n[4] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[5]) + m
	n[5] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[6]) + m
	n[6] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[7]) + m
	n[7] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[8]) + m&uint64(0xfffffffffffffffe)
	n[8] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[9]) + m
	n[9] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[10]) + m
	n[10] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[11]) + m
	n[11] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[12]) + m
	n[12] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[13]) + m
	n[13] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[14]) + m
	n[14] = word_t(carry) & radixMask
	carry >>= 28

	carry += uint64(n[15]) + m
	n[15] = word_t(carry) & radixMask
	carry >>= 28

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

	n[0] = word_t(accum0 & uint64(radixMask))
	accum0 >>= Radix

	n[8] = word_t(accum8 & uint64(radixMask))
	accum8 >>= Radix

	for i := 1; i < Limbs/2; i++ {
		accum0 += uint64(wlo) * uint64(x[i])
		accum8 += uint64(wlo) * uint64(x[i+8])
		accum0 += uint64(whi) * uint64(x[i-1])
		accum8 += uint64(whi) * uint64(x[i+7])

		n[i] = word_t(accum0 & uint64(radixMask))
		accum0 >>= Radix

		n[i+8] = word_t(accum8 & uint64(radixMask))
		accum8 >>= Radix
	}

	accum0 += accum8 + uint64(n[8])
	n[8] = word_t(accum0 & uint64(radixMask))
	n[9] += word_t(accum0 >> Radix)

	accum8 += uint64(n[0])
	n[0] = word_t(accum8 & uint64(radixMask))
	n[1] += word_t(accum8 >> Radix)

	return n
}
