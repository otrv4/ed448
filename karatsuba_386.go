package ed448

//c = a * b
func karatsubaMul(c, a, b *bigNumber) *bigNumber {
	var aa, bb [8]limb
	for i := 0; i < 8; i++ {
		aa[i] = a[i] + a[i+8]
		bb[i] = b[i] + b[i+8]
	}

	var z0, z1, z2 uint64
	for j := 0; j < 8; j++ {
		z2 = 0

		for i := 0; i <= j; i++ {
			z2 += uint64(a[j-i]) * uint64(b[i])
			z1 += uint64(aa[j-i]) * uint64(bb[i])
			z0 += uint64(a[8+j-i]) * uint64(b[8+i])
		}

		z1 -= z2
		z0 += z2
		z2 = 0

		for i := j + 1; i < 8; i++ {
			z0 -= uint64(a[8+j-i]) * uint64(b[i])
			z2 += uint64(aa[8+j-i]) * uint64(bb[i])
			z1 += uint64(a[16+j-i]) * uint64(b[8+i])
		}

		z1 += z2
		z0 += z2

		c[j] = limb(z0) & radixMask
		c[j+8] = limb(z1) & radixMask

		z0 >>= 28
		z1 >>= 28
	}

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
