package ed448

//XXX need a constant time implement
func mask(a, b *bigNumber, mask uint32) {
	for k := 0; k < len(b); k++ {
		a[k] = uint32(mask) & b[k]
	}
	return
}
