package ed448

//c = a * b
func karatsubaMul(c, a, b *bigNumber) *bigNumber {
	var aa, bb [8]dword
	for i := 0; i < 8; i++ {
		aa[i] = dword(a[i]) + dword(a[i+8])
		bb[i] = dword(b[i]) + dword(b[i+8])
	}

	var z0, z1, z2 dword

	//j = 0
	z2 = 0
	z2 += dword(a[0]) * dword(b[0])
	z1 += aa[0] * bb[0]
	z1 -= z2
	z0 += dword(a[8]) * dword(b[8])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[1]
	z2 += aa[6] * bb[2]
	z2 += aa[5] * bb[3]
	z2 += aa[4] * bb[4]
	z2 += aa[3] * bb[5]
	z2 += aa[2] * bb[6]
	z2 += aa[1] * bb[7]

	z1 += dword(a[15]) * dword(b[9])
	z1 += dword(a[14]) * dword(b[10])
	z1 += dword(a[13]) * dword(b[11])
	z1 += dword(a[12]) * dword(b[12])
	z1 += dword(a[11]) * dword(b[13])
	z1 += dword(a[10]) * dword(b[14])
	z1 += dword(a[9]) * dword(b[15])
	z1 += z2

	z0 -= dword(a[7]) * dword(b[1])
	z0 -= dword(a[6]) * dword(b[2])
	z0 -= dword(a[5]) * dword(b[3])
	z0 -= dword(a[4]) * dword(b[4])
	z0 -= dword(a[3]) * dword(b[5])
	z0 -= dword(a[2]) * dword(b[6])
	z0 -= dword(a[1]) * dword(b[7])
	z0 += z2

	c[0] = word(z0) & radixMask
	c[8] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 1
	z2 = 0
	z2 += dword(a[1]) * dword(b[0])
	z2 += dword(a[0]) * dword(b[1])

	z1 += aa[1] * bb[0]
	z1 += aa[0] * bb[1]
	z1 -= z2

	z0 += dword(a[9]) * dword(b[8])
	z0 += dword(a[8]) * dword(b[9])
	z0 += z2

	z2 = 0

	z2 += aa[7] * bb[2]
	z2 += aa[6] * bb[3]
	z2 += aa[5] * bb[4]
	z2 += aa[4] * bb[5]
	z2 += aa[3] * bb[6]
	z2 += aa[2] * bb[7]

	z1 += dword(a[15]) * dword(b[10])
	z1 += dword(a[14]) * dword(b[11])
	z1 += dword(a[13]) * dword(b[12])
	z1 += dword(a[12]) * dword(b[13])
	z1 += dword(a[11]) * dword(b[14])
	z1 += dword(a[10]) * dword(b[15])
	z1 += z2

	z0 -= dword(a[7]) * dword(b[2])
	z0 -= dword(a[6]) * dword(b[3])
	z0 -= dword(a[5]) * dword(b[4])
	z0 -= dword(a[4]) * dword(b[5])
	z0 -= dword(a[3]) * dword(b[6])
	z0 -= dword(a[2]) * dword(b[7])
	z0 += z2

	c[1] = word(z0) & radixMask
	c[9] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 2
	z2 = 0
	z2 += dword(a[2]) * dword(b[0])
	z2 += dword(a[1]) * dword(b[1])
	z2 += dword(a[0]) * dword(b[2])

	z1 += aa[2] * bb[0]
	z1 += aa[1] * bb[1]
	z1 += aa[0] * bb[2]
	z1 -= z2

	z0 += dword(a[10]) * dword(b[8])
	z0 += dword(a[9]) * dword(b[9])
	z0 += dword(a[8]) * dword(b[10])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[3]
	z2 += aa[6] * bb[4]
	z2 += aa[5] * bb[5]
	z2 += aa[4] * bb[6]
	z2 += aa[3] * bb[7]

	z1 += dword(a[15]) * dword(b[11])
	z1 += dword(a[14]) * dword(b[12])
	z1 += dword(a[13]) * dword(b[13])
	z1 += dword(a[12]) * dword(b[14])
	z1 += dword(a[11]) * dword(b[15])
	z1 += z2

	z0 -= dword(a[7]) * dword(b[3])
	z0 -= dword(a[6]) * dword(b[4])
	z0 -= dword(a[5]) * dword(b[5])
	z0 -= dword(a[4]) * dword(b[6])
	z0 -= dword(a[3]) * dword(b[7])
	z0 += z2

	c[2] = word(z0) & radixMask
	c[10] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 3
	z2 = 0
	z2 += dword(a[3]) * dword(b[0])
	z2 += dword(a[2]) * dword(b[1])
	z2 += dword(a[1]) * dword(b[2])
	z2 += dword(a[0]) * dword(b[3])

	z1 += aa[3] * bb[0]
	z1 += aa[2] * bb[1]
	z1 += aa[1] * bb[2]
	z1 += aa[0] * bb[3]
	z1 -= z2

	z0 += dword(a[11]) * dword(b[8])
	z0 += dword(a[10]) * dword(b[9])
	z0 += dword(a[9]) * dword(b[10])
	z0 += dword(a[8]) * dword(b[11])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[4]
	z2 += aa[6] * bb[5]
	z2 += aa[5] * bb[6]
	z2 += aa[4] * bb[7]

	z0 -= dword(a[7]) * dword(b[4])
	z0 -= dword(a[6]) * dword(b[5])
	z0 -= dword(a[5]) * dword(b[6])
	z0 -= dword(a[4]) * dword(b[7])
	z0 += z2

	z1 += dword(a[15]) * dword(b[12])
	z1 += dword(a[14]) * dword(b[13])
	z1 += dword(a[13]) * dword(b[14])
	z1 += dword(a[12]) * dword(b[15])
	z1 += z2

	c[3] = word(z0) & radixMask
	c[11] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 4
	z2 = 0
	z2 += dword(a[4]) * dword(b[0])
	z2 += dword(a[3]) * dword(b[1])
	z2 += dword(a[2]) * dword(b[2])
	z2 += dword(a[1]) * dword(b[3])
	z2 += dword(a[0]) * dword(b[4])

	z1 += aa[4] * bb[0]
	z1 += aa[3] * bb[1]
	z1 += aa[2] * bb[2]
	z1 += aa[1] * bb[3]
	z1 += aa[0] * bb[4]
	z1 -= z2

	z0 += dword(a[12]) * dword(b[8])
	z0 += dword(a[11]) * dword(b[9])
	z0 += dword(a[10]) * dword(b[10])
	z0 += dword(a[9]) * dword(b[11])
	z0 += dword(a[8]) * dword(b[12])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[5]
	z2 += aa[6] * bb[6]
	z2 += aa[5] * bb[7]

	z1 += dword(a[15]) * dword(b[13])
	z1 += dword(a[14]) * dword(b[14])
	z1 += dword(a[13]) * dword(b[15])
	z1 += z2

	z0 -= dword(a[7]) * dword(b[5])
	z0 -= dword(a[6]) * dword(b[6])
	z0 -= dword(a[5]) * dword(b[7])
	z0 += z2

	c[4] = word(z0) & radixMask
	c[12] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 5
	z2 = 0
	z2 += dword(a[5]) * dword(b[0])
	z2 += dword(a[4]) * dword(b[1])
	z2 += dword(a[3]) * dword(b[2])
	z2 += dword(a[2]) * dword(b[3])
	z2 += dword(a[1]) * dword(b[4])
	z2 += dword(a[0]) * dword(b[5])

	z1 += aa[5] * bb[0]
	z1 += aa[4] * bb[1]
	z1 += aa[3] * bb[2]
	z1 += aa[2] * bb[3]
	z1 += aa[1] * bb[4]
	z1 += aa[0] * bb[5]
	z1 -= z2

	z0 += dword(a[13]) * dword(b[8])
	z0 += dword(a[12]) * dword(b[9])
	z0 += dword(a[11]) * dword(b[10])
	z0 += dword(a[10]) * dword(b[11])
	z0 += dword(a[9]) * dword(b[12])
	z0 += dword(a[8]) * dword(b[13])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[6]
	z2 += aa[6] * bb[7]

	z1 += dword(a[15]) * dword(b[14])
	z1 += dword(a[14]) * dword(b[15])
	z1 += z2

	z0 -= dword(a[7]) * dword(b[6])
	z0 -= dword(a[6]) * dword(b[7])
	z0 += z2

	c[5] = word(z0) & radixMask
	c[13] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 6
	z2 = 0
	z2 += dword(a[6]) * dword(b[0])
	z2 += dword(a[5]) * dword(b[1])
	z2 += dword(a[4]) * dword(b[2])
	z2 += dword(a[3]) * dword(b[3])
	z2 += dword(a[2]) * dword(b[4])
	z2 += dword(a[1]) * dword(b[5])
	z2 += dword(a[0]) * dword(b[6])

	z1 += aa[6] * bb[0]
	z1 += aa[5] * bb[1]
	z1 += aa[4] * bb[2]
	z1 += aa[3] * bb[3]
	z1 += aa[2] * bb[4]
	z1 += aa[1] * bb[5]
	z1 += aa[0] * bb[6]
	z1 -= z2

	z0 += dword(a[14]) * dword(b[8])
	z0 += dword(a[13]) * dword(b[9])
	z0 += dword(a[12]) * dword(b[10])
	z0 += dword(a[11]) * dword(b[11])
	z0 += dword(a[10]) * dword(b[12])
	z0 += dword(a[9]) * dword(b[13])
	z0 += dword(a[8]) * dword(b[14])
	z0 += z2

	z2 = 0
	z2 += aa[7] * bb[7]
	z1 += dword(a[15]) * dword(b[15])
	z1 += z2
	z0 -= dword(a[7]) * dword(b[7])
	z0 += z2

	c[6] = word(z0) & radixMask
	c[14] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 7
	z2 = 0
	z2 += dword(a[7]) * dword(b[0])
	z2 += dword(a[6]) * dword(b[1])
	z2 += dword(a[5]) * dword(b[2])
	z2 += dword(a[4]) * dword(b[3])
	z2 += dword(a[3]) * dword(b[4])
	z2 += dword(a[2]) * dword(b[5])
	z2 += dword(a[1]) * dword(b[6])
	z2 += dword(a[0]) * dword(b[7])

	z1 += aa[7] * bb[0]
	z1 += aa[6] * bb[1]
	z1 += aa[5] * bb[2]
	z1 += aa[4] * bb[3]
	z1 += aa[3] * bb[4]
	z1 += aa[2] * bb[5]
	z1 += aa[1] * bb[6]
	z1 += aa[0] * bb[7]
	z1 -= z2

	z0 += dword(a[15]) * dword(b[8])
	z0 += dword(a[14]) * dword(b[9])
	z0 += dword(a[13]) * dword(b[10])
	z0 += dword(a[12]) * dword(b[11])
	z0 += dword(a[11]) * dword(b[12])
	z0 += dword(a[10]) * dword(b[13])
	z0 += dword(a[9]) * dword(b[14])
	z0 += dword(a[8]) * dword(b[15])
	z0 += z2

	z2 = 0
	z1 += z2
	z0 += z2

	c[7] = word(z0) & radixMask
	c[15] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	// finish

	z0 += z1
	z0 += dword(c[8])
	z1 += dword(c[0])

	c[8] = word(z0) & radixMask
	c[0] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	c[9] += word(z0)
	c[1] += word(z1)

	return c
}
