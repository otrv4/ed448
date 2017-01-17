package ed448

//c = a * b
func karatsubaMul(c, a, b *bigNumber) *bigNumber {
	var aa, bb [8]uint64
	for i := 0; i < 8; i++ {
		aa[i] = uint64(a[i]) + uint64(a[i+8])
		bb[i] = uint64(b[i]) + uint64(b[i+8])
	}

	var z0, z1, z2 uint64

	//j = 0
	z2 = 0
	z2 += uint64(a[0]) * uint64(b[0])
	z1 += aa[0] * bb[0]
	z1 -= z2
	z0 += uint64(a[8]) * uint64(b[8])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[1]
	z2 += aa[6] * bb[2]
	z2 += aa[5] * bb[3]
	z2 += aa[4] * bb[4]
	z2 += aa[3] * bb[5]
	z2 += aa[2] * bb[6]
	z2 += aa[1] * bb[7]

	z1 += uint64(a[15]) * uint64(b[9])
	z1 += uint64(a[14]) * uint64(b[10])
	z1 += uint64(a[13]) * uint64(b[11])
	z1 += uint64(a[12]) * uint64(b[12])
	z1 += uint64(a[11]) * uint64(b[13])
	z1 += uint64(a[10]) * uint64(b[14])
	z1 += uint64(a[9]) * uint64(b[15])
	z1 += z2

	z0 -= uint64(a[7]) * uint64(b[1])
	z0 -= uint64(a[6]) * uint64(b[2])
	z0 -= uint64(a[5]) * uint64(b[3])
	z0 -= uint64(a[4]) * uint64(b[4])
	z0 -= uint64(a[3]) * uint64(b[5])
	z0 -= uint64(a[2]) * uint64(b[6])
	z0 -= uint64(a[1]) * uint64(b[7])
	z0 += z2

	c[0] = word_t(z0) & radixMask
	c[8] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 1
	z2 = 0
	z2 += uint64(a[1]) * uint64(b[0])
	z2 += uint64(a[0]) * uint64(b[1])

	z1 += aa[1] * bb[0]
	z1 += aa[0] * bb[1]
	z1 -= z2

	z0 += uint64(a[9]) * uint64(b[8])
	z0 += uint64(a[8]) * uint64(b[9])
	z0 += z2

	z2 = 0

	z2 += aa[7] * bb[2]
	z2 += aa[6] * bb[3]
	z2 += aa[5] * bb[4]
	z2 += aa[4] * bb[5]
	z2 += aa[3] * bb[6]
	z2 += aa[2] * bb[7]

	z1 += uint64(a[15]) * uint64(b[10])
	z1 += uint64(a[14]) * uint64(b[11])
	z1 += uint64(a[13]) * uint64(b[12])
	z1 += uint64(a[12]) * uint64(b[13])
	z1 += uint64(a[11]) * uint64(b[14])
	z1 += uint64(a[10]) * uint64(b[15])
	z1 += z2

	z0 -= uint64(a[7]) * uint64(b[2])
	z0 -= uint64(a[6]) * uint64(b[3])
	z0 -= uint64(a[5]) * uint64(b[4])
	z0 -= uint64(a[4]) * uint64(b[5])
	z0 -= uint64(a[3]) * uint64(b[6])
	z0 -= uint64(a[2]) * uint64(b[7])
	z0 += z2

	c[1] = word_t(z0) & radixMask
	c[9] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 2
	z2 = 0
	z2 += uint64(a[2]) * uint64(b[0])
	z2 += uint64(a[1]) * uint64(b[1])
	z2 += uint64(a[0]) * uint64(b[2])

	z1 += aa[2] * bb[0]
	z1 += aa[1] * bb[1]
	z1 += aa[0] * bb[2]
	z1 -= z2

	z0 += uint64(a[10]) * uint64(b[8])
	z0 += uint64(a[9]) * uint64(b[9])
	z0 += uint64(a[8]) * uint64(b[10])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[3]
	z2 += aa[6] * bb[4]
	z2 += aa[5] * bb[5]
	z2 += aa[4] * bb[6]
	z2 += aa[3] * bb[7]

	z1 += uint64(a[15]) * uint64(b[11])
	z1 += uint64(a[14]) * uint64(b[12])
	z1 += uint64(a[13]) * uint64(b[13])
	z1 += uint64(a[12]) * uint64(b[14])
	z1 += uint64(a[11]) * uint64(b[15])
	z1 += z2

	z0 -= uint64(a[7]) * uint64(b[3])
	z0 -= uint64(a[6]) * uint64(b[4])
	z0 -= uint64(a[5]) * uint64(b[5])
	z0 -= uint64(a[4]) * uint64(b[6])
	z0 -= uint64(a[3]) * uint64(b[7])
	z0 += z2

	c[2] = word_t(z0) & radixMask
	c[10] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 3
	z2 = 0
	z2 += uint64(a[3]) * uint64(b[0])
	z2 += uint64(a[2]) * uint64(b[1])
	z2 += uint64(a[1]) * uint64(b[2])
	z2 += uint64(a[0]) * uint64(b[3])

	z1 += aa[3] * bb[0]
	z1 += aa[2] * bb[1]
	z1 += aa[1] * bb[2]
	z1 += aa[0] * bb[3]
	z1 -= z2

	z0 += uint64(a[11]) * uint64(b[8])
	z0 += uint64(a[10]) * uint64(b[9])
	z0 += uint64(a[9]) * uint64(b[10])
	z0 += uint64(a[8]) * uint64(b[11])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[4]
	z2 += aa[6] * bb[5]
	z2 += aa[5] * bb[6]
	z2 += aa[4] * bb[7]

	z0 -= uint64(a[7]) * uint64(b[4])
	z0 -= uint64(a[6]) * uint64(b[5])
	z0 -= uint64(a[5]) * uint64(b[6])
	z0 -= uint64(a[4]) * uint64(b[7])
	z0 += z2

	z1 += uint64(a[15]) * uint64(b[12])
	z1 += uint64(a[14]) * uint64(b[13])
	z1 += uint64(a[13]) * uint64(b[14])
	z1 += uint64(a[12]) * uint64(b[15])
	z1 += z2

	c[3] = word_t(z0) & radixMask
	c[11] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 4
	z2 = 0
	z2 += uint64(a[4]) * uint64(b[0])
	z2 += uint64(a[3]) * uint64(b[1])
	z2 += uint64(a[2]) * uint64(b[2])
	z2 += uint64(a[1]) * uint64(b[3])
	z2 += uint64(a[0]) * uint64(b[4])

	z1 += aa[4] * bb[0]
	z1 += aa[3] * bb[1]
	z1 += aa[2] * bb[2]
	z1 += aa[1] * bb[3]
	z1 += aa[0] * bb[4]
	z1 -= z2

	z0 += uint64(a[12]) * uint64(b[8])
	z0 += uint64(a[11]) * uint64(b[9])
	z0 += uint64(a[10]) * uint64(b[10])
	z0 += uint64(a[9]) * uint64(b[11])
	z0 += uint64(a[8]) * uint64(b[12])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[5]
	z2 += aa[6] * bb[6]
	z2 += aa[5] * bb[7]

	z1 += uint64(a[15]) * uint64(b[13])
	z1 += uint64(a[14]) * uint64(b[14])
	z1 += uint64(a[13]) * uint64(b[15])
	z1 += z2

	z0 -= uint64(a[7]) * uint64(b[5])
	z0 -= uint64(a[6]) * uint64(b[6])
	z0 -= uint64(a[5]) * uint64(b[7])
	z0 += z2

	c[4] = word_t(z0) & radixMask
	c[12] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 5
	z2 = 0
	z2 += uint64(a[5]) * uint64(b[0])
	z2 += uint64(a[4]) * uint64(b[1])
	z2 += uint64(a[3]) * uint64(b[2])
	z2 += uint64(a[2]) * uint64(b[3])
	z2 += uint64(a[1]) * uint64(b[4])
	z2 += uint64(a[0]) * uint64(b[5])

	z1 += aa[5] * bb[0]
	z1 += aa[4] * bb[1]
	z1 += aa[3] * bb[2]
	z1 += aa[2] * bb[3]
	z1 += aa[1] * bb[4]
	z1 += aa[0] * bb[5]
	z1 -= z2

	z0 += uint64(a[13]) * uint64(b[8])
	z0 += uint64(a[12]) * uint64(b[9])
	z0 += uint64(a[11]) * uint64(b[10])
	z0 += uint64(a[10]) * uint64(b[11])
	z0 += uint64(a[9]) * uint64(b[12])
	z0 += uint64(a[8]) * uint64(b[13])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[6]
	z2 += aa[6] * bb[7]

	z1 += uint64(a[15]) * uint64(b[14])
	z1 += uint64(a[14]) * uint64(b[15])
	z1 += z2

	z0 -= uint64(a[7]) * uint64(b[6])
	z0 -= uint64(a[6]) * uint64(b[7])
	z0 += z2

	c[5] = word_t(z0) & radixMask
	c[13] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 6
	z2 = 0
	z2 += uint64(a[6]) * uint64(b[0])
	z2 += uint64(a[5]) * uint64(b[1])
	z2 += uint64(a[4]) * uint64(b[2])
	z2 += uint64(a[3]) * uint64(b[3])
	z2 += uint64(a[2]) * uint64(b[4])
	z2 += uint64(a[1]) * uint64(b[5])
	z2 += uint64(a[0]) * uint64(b[6])

	z1 += aa[6] * bb[0]
	z1 += aa[5] * bb[1]
	z1 += aa[4] * bb[2]
	z1 += aa[3] * bb[3]
	z1 += aa[2] * bb[4]
	z1 += aa[1] * bb[5]
	z1 += aa[0] * bb[6]
	z1 -= z2

	z0 += uint64(a[14]) * uint64(b[8])
	z0 += uint64(a[13]) * uint64(b[9])
	z0 += uint64(a[12]) * uint64(b[10])
	z0 += uint64(a[11]) * uint64(b[11])
	z0 += uint64(a[10]) * uint64(b[12])
	z0 += uint64(a[9]) * uint64(b[13])
	z0 += uint64(a[8]) * uint64(b[14])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[7]
	z1 += uint64(a[15]) * uint64(b[15])
	z1 += z2
	z0 -= uint64(a[7]) * uint64(b[7])
	z0 += z2

	c[6] = word_t(z0) & radixMask
	c[14] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 7
	z2 = 0
	z2 += uint64(a[7]) * uint64(b[0])
	z2 += uint64(a[6]) * uint64(b[1])
	z2 += uint64(a[5]) * uint64(b[2])
	z2 += uint64(a[4]) * uint64(b[3])
	z2 += uint64(a[3]) * uint64(b[4])
	z2 += uint64(a[2]) * uint64(b[5])
	z2 += uint64(a[1]) * uint64(b[6])
	z2 += uint64(a[0]) * uint64(b[7])

	z1 += aa[7] * bb[0]
	z1 += aa[6] * bb[1]
	z1 += aa[5] * bb[2]
	z1 += aa[4] * bb[3]
	z1 += aa[3] * bb[4]
	z1 += aa[2] * bb[5]
	z1 += aa[1] * bb[6]
	z1 += aa[0] * bb[7]
	z1 -= z2

	z0 += uint64(a[15]) * uint64(b[8])
	z0 += uint64(a[14]) * uint64(b[9])
	z0 += uint64(a[13]) * uint64(b[10])
	z0 += uint64(a[12]) * uint64(b[11])
	z0 += uint64(a[11]) * uint64(b[12])
	z0 += uint64(a[10]) * uint64(b[13])
	z0 += uint64(a[9]) * uint64(b[14])
	z0 += uint64(a[8]) * uint64(b[15])
	z0 += z2

	z2 = 0
	z1 += z2
	z0 += z2

	c[7] = word_t(z0) & radixMask
	c[15] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	// finish

	z0 += z1
	z0 += uint64(c[8])
	z1 += uint64(c[0])

	c[8] = word_t(z0) & radixMask
	c[0] = word_t(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	c[9] += word_t(z0)
	c[1] += word_t(z1)

	return c
}
