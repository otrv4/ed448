package ed448

// TODO: this has changed
//c = a * a
func karatsubaSquare(c, a *bigNumber) *bigNumber {
	aa := [8]dword{
		dword(a[0]) + dword(a[8]), // 0 - 8
		dword(a[1]) + dword(a[9]), // 1 - 9
		dword(a[2]) + dword(a[10]),
		dword(a[3]) + dword(a[11]),
		dword(a[4]) + dword(a[12]),
		dword(a[5]) + dword(a[13]),
		dword(a[6]) + dword(a[14]),
		dword(a[7]) + dword(a[15]), //7 - 15
	}

	var z0, z1, z2 dword

	//j = 0
	z2 = 0
	z2 += dword(a[0]) * dword(a[0])
	z1 += aa[0] * aa[0]
	z1 -= z2
	z0 += dword(a[8]) * dword(a[8])
	z0 += z2

	z2 = 0
	z2 += (aa[7] * aa[1]) << 1 // (a7+a15) * (a1+a9)
	z2 += (aa[6] * aa[2]) << 1 // (a6+a14) * (a2+a10)
	z2 += (aa[5] * aa[3]) << 1 // (a5+a13) * (a3+a11)
	z2 += aa[4] * aa[4]

	z1 += (dword(a[15]) * dword(a[9])) << 1
	z1 += (dword(a[14]) * dword(a[10])) << 1
	z1 += (dword(a[13]) * dword(a[11])) << 1
	z1 += dword(a[12]) * dword(a[12])
	z1 += z2

	z0 -= (dword(a[7]) * dword(a[1])) << 1
	z0 -= (dword(a[6]) * dword(a[2])) << 1
	z0 -= (dword(a[5]) * dword(a[3])) << 1
	z0 -= dword(a[4]) * dword(a[4])
	z0 += z2

	c[0] = word(z0) & radixMask
	c[8] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 1
	z2 = (dword(a[1]) * dword(a[0])) << 1

	z1 += (aa[1] * aa[0]) << 1
	z1 -= z2

	z0 += (dword(a[9]) * dword(a[8])) << 1
	z0 += z2

	z2 = 0
	z2 += aa[7] * aa[2]
	z2 += aa[6] * aa[3]
	z2 += aa[5] * aa[4]
	z2 <<= 1

	z1 += (dword(a[15]) * dword(a[10])) << 1
	z1 += (dword(a[14]) * dword(a[11])) << 1
	z1 += (dword(a[13]) * dword(a[12])) << 1
	z1 += z2

	z0 -= (dword(a[7]) * dword(a[2])) << 1
	z0 -= (dword(a[6]) * dword(a[3])) << 1
	z0 -= (dword(a[5]) * dword(a[4])) << 1
	z0 += z2

	c[1] = word(z0) & radixMask
	c[9] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 2
	z2 = 0
	z2 += (dword(a[2]) * dword(a[0])) << 1
	z2 += dword(a[1]) * dword(a[1])

	z1 += (aa[2] * aa[0]) << 1
	z1 += aa[1] * aa[1]
	z1 -= z2

	z0 += (dword(a[10]) * dword(a[8])) << 1
	z0 += dword(a[9]) * dword(a[9])
	z0 += z2

	z2 = 0
	z2 += aa[7] * aa[3]
	z2 += aa[6] * aa[4]
	z2 <<= 1
	z2 += aa[5] * aa[5]

	z1 += (dword(a[15]) * dword(a[11])) << 1
	z1 += (dword(a[14]) * dword(a[12])) << 1
	z1 += dword(a[13]) * dword(a[13])
	z1 += z2

	z0 -= (dword(a[7]) * dword(a[3])) << 1
	z0 -= (dword(a[6]) * dword(a[4])) << 1
	z0 -= dword(a[5]) * dword(a[5])
	z0 += z2

	c[2] = word(z0) & radixMask
	c[10] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 3
	z2 = 0
	z2 += dword(a[3]) * dword(a[0])
	z2 += dword(a[2]) * dword(a[1])
	z2 <<= 1

	z1 += (aa[3] * aa[0]) << 1
	z1 += (aa[2] * aa[1]) << 1
	z1 -= z2

	z0 += (dword(a[11]) * dword(a[8])) << 1
	z0 += (dword(a[10]) * dword(a[9])) << 1
	z0 += z2

	z2 = 0
	z2 += (aa[7] * aa[4]) << 1
	z2 += (aa[6] * aa[5]) << 1

	z0 -= (dword(a[7]) * dword(a[4])) << 1
	z0 -= (dword(a[6]) * dword(a[5])) << 1
	z0 += z2

	z1 += (dword(a[15]) * dword(a[12])) << 1
	z1 += (dword(a[14]) * dword(a[13])) << 1
	z1 += z2

	c[3] = word(z0) & radixMask
	c[11] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 4
	z2 = 0
	z2 += (dword(a[4]) * dword(a[0])) << 1
	z2 += (dword(a[3]) * dword(a[1])) << 1
	z2 += dword(a[2]) * dword(a[2])

	z1 += (aa[4] * aa[0]) << 1
	z1 += (aa[3] * aa[1]) << 1
	z1 += aa[2] * aa[2]
	z1 -= z2

	z0 += (dword(a[12]) * dword(a[8])) << 1
	z0 += (dword(a[11]) * dword(a[9])) << 1
	z0 += dword(a[10]) * dword(a[10])
	z0 += z2

	z2 = 0
	z2 += (aa[7] * aa[5]) << 1
	z2 += aa[6] * aa[6]

	z1 += (dword(a[15]) * dword(a[13])) << 1
	z1 += dword(a[14]) * dword(a[14])
	z1 += z2

	z0 -= (dword(a[7]) * dword(a[5])) << 1
	z0 -= dword(a[6]) * dword(a[6])
	z0 += z2

	c[4] = word(z0) & radixMask
	c[12] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 5
	z2 = 0
	z2 += (dword(a[5]) * dword(a[0])) << 1
	z2 += (dword(a[4]) * dword(a[1])) << 1
	z2 += (dword(a[3]) * dword(a[2])) << 1

	z1 += (aa[5] * aa[0]) << 1
	z1 += (aa[4] * aa[1]) << 1
	z1 += (aa[3] * aa[2]) << 1
	z1 -= z2

	z0 += (dword(a[13]) * dword(a[8])) << 1
	z0 += (dword(a[12]) * dword(a[9])) << 1
	z0 += (dword(a[11]) * dword(a[10])) << 1
	z0 += z2

	z2 = 0
	z2 += (aa[7] * aa[6]) << 1

	z1 += (dword(a[15]) * dword(a[14])) << 1
	z1 += z2

	z0 -= (dword(a[7]) * dword(a[6])) << 1
	z0 += z2

	c[5] = word(z0) & radixMask
	c[13] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 6
	z2 = 0
	z2 += (dword(a[6]) * dword(a[0])) << 1
	z2 += (dword(a[5]) * dword(a[1])) << 1
	z2 += (dword(a[4]) * dword(a[2])) << 1
	z2 += dword(a[3]) * dword(a[3])

	z1 += (aa[6] * aa[0]) << 1
	z1 += (aa[5] * aa[1]) << 1
	z1 += (aa[4] * aa[2]) << 1
	z1 += aa[3] * aa[3]
	z1 -= z2

	z0 += (dword(a[14]) * dword(a[8])) << 1
	z0 += (dword(a[13]) * dword(a[9])) << 1
	z0 += (dword(a[12]) * dword(a[10])) << 1
	z0 += dword(a[11]) * dword(a[11])
	z0 += z2

	z2 = 0
	z2 += aa[7] * aa[7]
	z1 += dword(a[15]) * dword(a[15])
	z1 += z2
	z0 -= dword(a[7]) * dword(a[7])
	z0 += z2

	c[6] = word(z0) & radixMask
	c[14] = word(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 7
	z2 = 0
	z2 += (dword(a[7]) * dword(a[0])) << 1
	z2 += (dword(a[6]) * dword(a[1])) << 1
	z2 += (dword(a[5]) * dword(a[2])) << 1
	z2 += (dword(a[4]) * dword(a[3])) << 1

	z1 += (aa[7] * aa[0]) << 1
	z1 += (aa[6] * aa[1]) << 1
	z1 += (aa[5] * aa[2]) << 1
	z1 += (aa[4] * aa[3]) << 1
	z1 -= z2

	z0 += (dword(a[15]) * dword(a[8])) << 1
	z0 += (dword(a[14]) * dword(a[9])) << 1
	z0 += (dword(a[13]) * dword(a[10])) << 1
	z0 += (dword(a[12]) * dword(a[11])) << 1
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
