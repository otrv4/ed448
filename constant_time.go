package ed448

//XXX need a constant time implement
func mask(a, b *bigNumber, mask word) {
	for k := 0; k < len(b); k++ {
		a[k] = word(mask) & b[k]
	}
	return
}
