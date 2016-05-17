package ed448

func WideMul(x, y Word) DWord {
	z1, z0 := mulWW_g(x, y)
	return DWord{h: z1, l: z0}
}

func multiplyAndAdd(to DWord, x, y Word) (z DWord) {
	var z1 Word
	z1, z.l = mulAddWWW_g(x, y, to.l)
	z.h = z1 + to.h

	//XXX Should we panic on overflow?
	//if z11 < z1 { panic() }

	return
}

func multiplyAndSubtract(from DWord, a, b Word) DWord {
	d := WideMul(a, b)
	return subDWord(from, d)
}

func wideMul(x, y limb) DWord {
	return WideMul(Word(x), Word(y))
}

func mac(x, y limb, acc DWord) DWord {
	return multiplyAndAdd(acc, Word(x), Word(y))
}

func msb(x, y limb, acc DWord) DWord {
	return multiplyAndSubtract(acc, Word(x), Word(y))
}

func karatsubaMul(a, b bigNumber) (c bigNumber) {
	var aa, bb, bbb [4]limb

	mask := Word(0xffffffffffffff)

	for i := 0; i < 4; i++ {
		aa[i] = a[i] + a[i+4]
		bb[i] = b[i] + b[i+4]
		bbb[i] = bb[i] + b[i+4]
	}

	accum2 := wideMul(a[0], b[3])
	accum0 := wideMul(aa[0], bb[3])
	accum1 := wideMul(a[4], b[7])

	accum2 = mac(a[1], b[2], accum2)
	accum0 = mac(aa[1], bb[2], accum0)
	accum1 = mac(a[5], b[6], accum1)

	accum2 = mac(a[2], b[1], accum2)
	accum0 = mac(aa[2], bb[1], accum0)
	accum1 = mac(a[6], b[5], accum1)

	accum2 = mac(a[3], b[0], accum2)
	accum0 = mac(aa[3], bb[0], accum0)
	accum1 = mac(a[7], b[4], accum1)

	//If borrow != 0, we should panic?
	accum0 = subDWord(accum0, accum2)
	accum1 = addDWord(accum1, accum2)

	c[3] = limb(accum1.l & mask)
	c[7] = limb(accum0.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum0 = mac(aa[1], bb[3], accum0)
	accum1 = mac(a[5], b[7], accum1)
	accum0 = mac(aa[2], bb[2], accum0)
	accum1 = mac(a[6], b[6], accum1)
	accum0 = mac(aa[3], bb[1], accum0)

	accum1 = addDWord(accum1, accum0)
	accum2 = wideMul(a[0], b[0])
	accum1 = subDWord(accum1, accum2)
	accum0 = addDWord(accum0, accum2)

	accum0 = msb(a[1], b[3], accum0)
	accum0 = msb(a[2], b[2], accum0)
	accum1 = mac(a[7], b[5], accum1)
	accum0 = msb(a[3], b[1], accum0)
	accum1 = mac(aa[0], bb[0], accum1)
	accum0 = mac(a[4], b[4], accum0)

	//c[3+i], c[3+i mod 7]
	c[0] = limb(accum0.l & mask)
	c[4] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum2 = wideMul(a[2], b[7])
	accum0 = mac(a[6], bb[3], accum0)
	accum1 = mac(aa[2], bbb[3], accum1)

	accum2 = mac(a[3], b[6], accum2)
	accum0 = mac(a[7], bb[2], accum0)
	accum1 = mac(aa[3], bbb[2], accum1)

	accum2 = mac(a[0], b[1], accum2)
	accum1 = mac(aa[0], bb[1], accum1)
	accum0 = mac(a[4], b[5], accum0)

	accum2 = mac(a[1], b[0], accum2)
	accum1 = mac(aa[1], bb[0], accum1)
	accum0 = mac(a[5], b[4], accum0)

	accum1 = subDWord(accum1, accum2)
	accum0 = addDWord(accum0, accum2)

	c[1] = limb(accum0.l & mask)
	c[5] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum2 = wideMul(a[3], b[7])
	accum0 = mac(a[7], bb[3], accum0)
	accum1 = mac(aa[3], bbb[3], accum1)

	accum2 = mac(a[0], b[2], accum2)
	accum1 = mac(aa[0], bb[2], accum1)
	accum0 = mac(a[4], b[6], accum0)

	accum2 = mac(a[1], b[1], accum2)
	accum1 = mac(aa[1], bb[1], accum1)
	accum0 = mac(a[5], b[5], accum0)

	accum2 = mac(a[2], b[0], accum2)
	accum1 = mac(aa[2], bb[0], accum1)
	accum0 = mac(a[6], b[4], accum0)

	accum1 = subDWord(accum1, accum2)
	accum0 = addDWord(accum0, accum2)

	c[2] = limb(accum0.l & mask)
	c[6] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum0 = addDWord(accum0, DWord{0, Word(c[3])})
	accum1 = addDWord(accum1, DWord{0, Word(c[7])})

	c[3] = limb(accum0.l & mask)
	c[7] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	c[0] += limb(accum1.l)
	c[4] += limb(accum0.l + accum1.l)

	return
}
