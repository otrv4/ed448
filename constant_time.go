package ed448

//XXX need a constant time implement
func mask(a, b *bigNumber, mask uint32) {
	for k := 0; k < len(b); k += 1 {
		a[k] = word_t(mask) & b[k]
	}
	return
}
