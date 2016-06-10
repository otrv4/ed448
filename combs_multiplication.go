package ed448

func scheduleScalarForCombs(schedule []word_t, scalar [scalarWords]word_t) {
	table := baseTable
	scalar3 := make([]word_t, len(schedule))

	for i, _ := range scalar3 {
		scalar3[i] = scalar[i]
	}

	scalar3[len(scalar3)-1] &= (word_t(1) << (scalarBits % wordBits)) - 1

	convertToSignedWindowForm(schedule, scalar3, table.adjustments[:])
}

func convertToSignedWindowForm(out []word_t, scalar []word_t, preparedData []word_t) {
	mask := word_t(dword_t(-(scalar[0] & 1)) & 0xffffffff)

	carry := add_nr_ext_packed(out, scalar, scalarWords, preparedData, scalarWords, word_t(^mask))
	carry += add_nr_ext_packed(out, out, scalarWords, preparedData[scalarWords:], scalarWords, word_t(mask))

	for i := 0; i < scalarWords-1; i++ {
		out[i] >>= 1
		out[i] |= out[i+1] << (wordBits - 1)
	}

	out[scalarWords-1] >>= 1
	out[scalarWords-1] |= carry << (wordBits - 1)
}

func add_nr_ext_packed(out []word_t, a []word_t, wordsA uint32, c []word_t, wordsC uint32, mask word_t) word_t {
	i := uint32(0)
	carry := dword_t(0)

	for i = 0; i < wordsC; i++ {
		carry += dword_t(a[i]) + dword_t(c[i]&mask)
		out[i] = word_t(carry)
		carry >>= wordBits
	}

	//XXX Wont execute because in our case words have same size
	for ; i < wordsA; i++ {
		carry += dword_t(a[i])
		out[i] = word_t(carry)
		carry >>= wordBits
	}

	return word_t(carry)
}
