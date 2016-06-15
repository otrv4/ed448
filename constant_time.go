package ed448

//XXX need a constant time implement
func condSwap(a, b *bigNumber, doswap word_t) {
	if doswap != 0 {
		c := a.copy()
		a = b.copy()
		b = c
	}
	return
}

//XXX need a constant time implement
func mask(a, b *bigNumber, mask uint32) {
	for k := 0; k < len(b); k += 1 {
		a[k] = limb(mask) & b[k]
	}
	return
}
