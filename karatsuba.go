package ed448

//c = a * b
func karatsubaMul(c, a, b *bigNumber) *bigNumber {
	var aa, bb [8]limb
	for i := 0; i < 8; i++ {
		aa[i] = a[i] + a[i+8]
		bb[i] = b[i] + b[i+8]
	}

	var z0, z1, z2 uint64

	//j = 0
	z2 = 0

	z2 += uint64(a[0]) * uint64(b[0])
	z1 += uint64(aa[0]) * uint64(bb[0])
	z0 += uint64(a[8+0]) * uint64(b[8])

	z1 -= z2
	z0 += z2
	z2 = 0

	z0 -= uint64(a[8+0-1]) * uint64(b[1])
	z2 += uint64(aa[8+0-1]) * uint64(bb[1])
	z1 += uint64(a[16+0-1]) * uint64(b[8+1])

	z0 -= uint64(a[8+0-2]) * uint64(b[2])
	z2 += uint64(aa[8+0-2]) * uint64(bb[2])
	z1 += uint64(a[16+0-2]) * uint64(b[8+2])

	z0 -= uint64(a[8+0-3]) * uint64(b[3])
	z2 += uint64(aa[8+0-3]) * uint64(bb[3])
	z1 += uint64(a[16+0-3]) * uint64(b[8+3])

	z0 -= uint64(a[8+0-4]) * uint64(b[4])
	z2 += uint64(aa[8+0-4]) * uint64(bb[4])
	z1 += uint64(a[16+0-4]) * uint64(b[8+4])

	z0 -= uint64(a[8+0-5]) * uint64(b[5])
	z2 += uint64(aa[8+0-5]) * uint64(bb[5])
	z1 += uint64(a[16+0-5]) * uint64(b[8+5])

	z0 -= uint64(a[8+0-6]) * uint64(b[6])
	z2 += uint64(aa[8+0-6]) * uint64(bb[6])
	z1 += uint64(a[16+0-6]) * uint64(b[8+6])

	z0 -= uint64(a[8+0-7]) * uint64(b[7])
	z2 += uint64(aa[8+0-7]) * uint64(bb[7])
	z1 += uint64(a[16+0-7]) * uint64(b[8+7])

	z1 += z2
	z0 += z2

	c[0] = limb(z0) & radixMask
	c[0+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 1
	z2 = 0

	z2 += uint64(a[1]) * uint64(b[0])
	z1 += uint64(aa[1]) * uint64(bb[0])
	z0 += uint64(a[8+1]) * uint64(b[8])

	z2 += uint64(a[1-1]) * uint64(b[1])
	z1 += uint64(aa[1-1]) * uint64(bb[1])
	z0 += uint64(a[8+1-1]) * uint64(b[8+1])

	z1 -= z2
	z0 += z2
	z2 = 0

	z0 -= uint64(a[8+1-2]) * uint64(b[2])
	z2 += uint64(aa[8+1-2]) * uint64(bb[2])
	z1 += uint64(a[16+1-2]) * uint64(b[8+2])

	z0 -= uint64(a[8+1-3]) * uint64(b[3])
	z2 += uint64(aa[8+1-3]) * uint64(bb[3])
	z1 += uint64(a[16+1-3]) * uint64(b[8+3])

	z0 -= uint64(a[8+1-4]) * uint64(b[4])
	z2 += uint64(aa[8+1-4]) * uint64(bb[4])
	z1 += uint64(a[16+1-4]) * uint64(b[8+4])

	z0 -= uint64(a[8+1-5]) * uint64(b[5])
	z2 += uint64(aa[8+1-5]) * uint64(bb[5])
	z1 += uint64(a[16+1-5]) * uint64(b[8+5])

	z0 -= uint64(a[8+1-6]) * uint64(b[6])
	z2 += uint64(aa[8+1-6]) * uint64(bb[6])
	z1 += uint64(a[16+1-6]) * uint64(b[8+6])

	z0 -= uint64(a[8+1-7]) * uint64(b[7])
	z2 += uint64(aa[8+1-7]) * uint64(bb[7])
	z1 += uint64(a[16+1-7]) * uint64(b[8+7])

	z1 += z2
	z0 += z2

	c[1] = limb(z0) & radixMask
	c[1+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 2
	z2 = 0

	z2 += uint64(a[2]) * uint64(b[0])
	z1 += uint64(aa[2]) * uint64(bb[0])
	z0 += uint64(a[8+2]) * uint64(b[8])

	z2 += uint64(a[2-1]) * uint64(b[1])
	z1 += uint64(aa[2-1]) * uint64(bb[1])
	z0 += uint64(a[8+2-1]) * uint64(b[8+1])

	z2 += uint64(a[2-2]) * uint64(b[2])
	z1 += uint64(aa[2-2]) * uint64(bb[2])
	z0 += uint64(a[8+2-2]) * uint64(b[8+2])

	z1 -= z2
	z0 += z2
	z2 = 0

	z0 -= uint64(a[8+2-3]) * uint64(b[3])
	z2 += uint64(aa[8+2-3]) * uint64(bb[3])
	z1 += uint64(a[16+2-3]) * uint64(b[8+3])

	z0 -= uint64(a[8+2-4]) * uint64(b[4])
	z2 += uint64(aa[8+2-4]) * uint64(bb[4])
	z1 += uint64(a[16+2-4]) * uint64(b[8+4])

	z0 -= uint64(a[8+2-5]) * uint64(b[5])
	z2 += uint64(aa[8+2-5]) * uint64(bb[5])
	z1 += uint64(a[16+2-5]) * uint64(b[8+5])

	z0 -= uint64(a[8+2-6]) * uint64(b[6])
	z2 += uint64(aa[8+2-6]) * uint64(bb[6])
	z1 += uint64(a[16+2-6]) * uint64(b[8+6])

	z0 -= uint64(a[8+2-7]) * uint64(b[7])
	z2 += uint64(aa[8+2-7]) * uint64(bb[7])
	z1 += uint64(a[16+2-7]) * uint64(b[8+7])

	z1 += z2
	z0 += z2

	c[2] = limb(z0) & radixMask
	c[2+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 3
	z2 = 0

	z2 += uint64(a[3]) * uint64(b[0])
	z1 += uint64(aa[3]) * uint64(bb[0])
	z0 += uint64(a[8+3]) * uint64(b[8])

	z2 += uint64(a[3-1]) * uint64(b[1])
	z1 += uint64(aa[3-1]) * uint64(bb[1])
	z0 += uint64(a[8+3-1]) * uint64(b[8+1])

	z2 += uint64(a[3-2]) * uint64(b[2])
	z1 += uint64(aa[3-2]) * uint64(bb[2])
	z0 += uint64(a[8+3-2]) * uint64(b[8+2])

	z2 += uint64(a[3-3]) * uint64(b[3])
	z1 += uint64(aa[3-3]) * uint64(bb[3])
	z0 += uint64(a[8+3-3]) * uint64(b[8+3])

	z1 -= z2
	z0 += z2
	z2 = 0

	z0 -= uint64(a[8+3-4]) * uint64(b[4])
	z2 += uint64(aa[8+3-4]) * uint64(bb[4])
	z1 += uint64(a[16+3-4]) * uint64(b[8+4])

	z0 -= uint64(a[8+3-5]) * uint64(b[5])
	z2 += uint64(aa[8+3-5]) * uint64(bb[5])
	z1 += uint64(a[16+3-5]) * uint64(b[8+5])

	z0 -= uint64(a[8+3-6]) * uint64(b[6])
	z2 += uint64(aa[8+3-6]) * uint64(bb[6])
	z1 += uint64(a[16+3-6]) * uint64(b[8+6])

	z0 -= uint64(a[8+3-7]) * uint64(b[7])
	z2 += uint64(aa[8+3-7]) * uint64(bb[7])
	z1 += uint64(a[16+3-7]) * uint64(b[8+7])

	z1 += z2
	z0 += z2

	c[3] = limb(z0) & radixMask
	c[3+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 4
	z2 = 0

	z2 += uint64(a[4]) * uint64(b[0])
	z1 += uint64(aa[4]) * uint64(bb[0])
	z0 += uint64(a[8+4]) * uint64(b[8])

	z2 += uint64(a[4-1]) * uint64(b[1])
	z1 += uint64(aa[4-1]) * uint64(bb[1])
	z0 += uint64(a[8+4-1]) * uint64(b[8+1])

	z2 += uint64(a[4-2]) * uint64(b[2])
	z1 += uint64(aa[4-2]) * uint64(bb[2])
	z0 += uint64(a[8+4-2]) * uint64(b[8+2])

	z2 += uint64(a[4-3]) * uint64(b[3])
	z1 += uint64(aa[4-3]) * uint64(bb[3])
	z0 += uint64(a[8+4-3]) * uint64(b[8+3])

	z2 += uint64(a[4-4]) * uint64(b[4])
	z1 += uint64(aa[4-4]) * uint64(bb[4])
	z0 += uint64(a[8+4-4]) * uint64(b[8+4])

	z1 -= z2
	z0 += z2
	z2 = 0

	z0 -= uint64(a[8+4-5]) * uint64(b[5])
	z2 += uint64(aa[8+4-5]) * uint64(bb[5])
	z1 += uint64(a[16+4-5]) * uint64(b[8+5])

	z0 -= uint64(a[8+4-6]) * uint64(b[6])
	z2 += uint64(aa[8+4-6]) * uint64(bb[6])
	z1 += uint64(a[16+4-6]) * uint64(b[8+6])

	z0 -= uint64(a[8+4-7]) * uint64(b[7])
	z2 += uint64(aa[8+4-7]) * uint64(bb[7])
	z1 += uint64(a[16+4-7]) * uint64(b[8+7])

	z1 += z2
	z0 += z2

	c[4] = limb(z0) & radixMask
	c[4+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 5
	z2 = 0

	z2 += uint64(a[5]) * uint64(b[0])
	z1 += uint64(aa[5]) * uint64(bb[0])
	z0 += uint64(a[8+5]) * uint64(b[8])

	z2 += uint64(a[5-1]) * uint64(b[1])
	z1 += uint64(aa[5-1]) * uint64(bb[1])
	z0 += uint64(a[8+5-1]) * uint64(b[8+1])

	z2 += uint64(a[5-2]) * uint64(b[2])
	z1 += uint64(aa[5-2]) * uint64(bb[2])
	z0 += uint64(a[8+5-2]) * uint64(b[8+2])

	z2 += uint64(a[5-3]) * uint64(b[3])
	z1 += uint64(aa[5-3]) * uint64(bb[3])
	z0 += uint64(a[8+5-3]) * uint64(b[8+3])

	z2 += uint64(a[5-4]) * uint64(b[4])
	z1 += uint64(aa[5-4]) * uint64(bb[4])
	z0 += uint64(a[8+5-4]) * uint64(b[8+4])

	z2 += uint64(a[5-5]) * uint64(b[5])
	z1 += uint64(aa[5-5]) * uint64(bb[5])
	z0 += uint64(a[8+5-5]) * uint64(b[8+5])

	z1 -= z2
	z0 += z2
	z2 = 0

	z0 -= uint64(a[8+5-6]) * uint64(b[6])
	z2 += uint64(aa[8+5-6]) * uint64(bb[6])
	z1 += uint64(a[16+5-6]) * uint64(b[8+6])

	z0 -= uint64(a[8+5-7]) * uint64(b[7])
	z2 += uint64(aa[8+5-7]) * uint64(bb[7])
	z1 += uint64(a[16+5-7]) * uint64(b[8+7])

	z1 += z2
	z0 += z2

	c[5] = limb(z0) & radixMask
	c[5+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 6
	z2 = 0

	z2 += uint64(a[6]) * uint64(b[0])
	z1 += uint64(aa[6]) * uint64(bb[0])
	z0 += uint64(a[8+6]) * uint64(b[8])

	z2 += uint64(a[6-1]) * uint64(b[1])
	z1 += uint64(aa[6-1]) * uint64(bb[1])
	z0 += uint64(a[8+6-1]) * uint64(b[8+1])

	z2 += uint64(a[6-2]) * uint64(b[2])
	z1 += uint64(aa[6-2]) * uint64(bb[2])
	z0 += uint64(a[8+6-2]) * uint64(b[8+2])

	z2 += uint64(a[6-3]) * uint64(b[3])
	z1 += uint64(aa[6-3]) * uint64(bb[3])
	z0 += uint64(a[8+6-3]) * uint64(b[8+3])

	z2 += uint64(a[6-4]) * uint64(b[4])
	z1 += uint64(aa[6-4]) * uint64(bb[4])
	z0 += uint64(a[8+6-4]) * uint64(b[8+4])

	z2 += uint64(a[6-5]) * uint64(b[5])
	z1 += uint64(aa[6-5]) * uint64(bb[5])
	z0 += uint64(a[8+6-5]) * uint64(b[8+5])

	z2 += uint64(a[6-6]) * uint64(b[6])
	z1 += uint64(aa[6-6]) * uint64(bb[6])
	z0 += uint64(a[8+6-6]) * uint64(b[8+6])

	z1 -= z2
	z0 += z2
	z2 = 0

	z0 -= uint64(a[8+6-7]) * uint64(b[7])
	z2 += uint64(aa[8+6-7]) * uint64(bb[7])
	z1 += uint64(a[16+6-7]) * uint64(b[8+7])

	z1 += z2
	z0 += z2

	c[6] = limb(z0) & radixMask
	c[6+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	//j = 7
	z2 = 0

	z2 += uint64(a[7]) * uint64(b[0])
	z1 += uint64(aa[7]) * uint64(bb[0])
	z0 += uint64(a[8+7]) * uint64(b[8])

	z2 += uint64(a[7-1]) * uint64(b[1])
	z1 += uint64(aa[7-1]) * uint64(bb[1])
	z0 += uint64(a[8+7-1]) * uint64(b[8+1])

	z2 += uint64(a[7-2]) * uint64(b[2])
	z1 += uint64(aa[7-2]) * uint64(bb[2])
	z0 += uint64(a[8+7-2]) * uint64(b[8+2])

	z2 += uint64(a[7-3]) * uint64(b[3])
	z1 += uint64(aa[7-3]) * uint64(bb[3])
	z0 += uint64(a[8+7-3]) * uint64(b[8+3])

	z2 += uint64(a[7-4]) * uint64(b[4])
	z1 += uint64(aa[7-4]) * uint64(bb[4])
	z0 += uint64(a[8+7-4]) * uint64(b[8+4])

	z2 += uint64(a[7-5]) * uint64(b[5])
	z1 += uint64(aa[7-5]) * uint64(bb[5])
	z0 += uint64(a[8+7-5]) * uint64(b[8+5])

	z2 += uint64(a[7-6]) * uint64(b[6])
	z1 += uint64(aa[7-6]) * uint64(bb[6])
	z0 += uint64(a[8+7-6]) * uint64(b[8+6])

	z2 += uint64(a[7-7]) * uint64(b[7])
	z1 += uint64(aa[7-7]) * uint64(bb[7])
	z0 += uint64(a[8+7-7]) * uint64(b[8+7])

	z1 -= z2
	z0 += z2
	z2 = 0

	z1 += z2
	z0 += z2

	c[7] = limb(z0) & radixMask
	c[7+8] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	// finish

	z0 += z1
	z0 += uint64(c[8])
	z1 += uint64(c[0])

	c[8] = limb(z0) & radixMask
	c[0] = limb(z1) & radixMask

	z0 >>= 28
	z1 >>= 28

	c[9] += limb(z0)
	c[1] += limb(z1)

	return c
}
