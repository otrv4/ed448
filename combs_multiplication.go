package ed448

func scheduleScalarForCombs(schedule []word_t, scalar [scalarWords]word_t) {
	table := baseTable
	tmp := make([]word_t, len(schedule))
	copy(tmp, scalar[:])

	tmp[len(tmp)-1] &= (word_t(1) << (scalarBits % wordBits)) - 1

	convertToSignedWindowForm(schedule, tmp, table.adjustments[:])
}

func convertToSignedWindowForm(out []word_t, scalar []word_t, preparedData []word_t) {
	mask := word_t(dword_t(-(scalar[0] & 1)) & lmask)

	carry := addExtPacked(out, scalar, preparedData[:scalarWords], word_t(^mask))
	carry += addExtPacked(out, out, preparedData[scalarWords:], word_t(mask))

	for i := 0; i < scalarWords-1; i++ {
		out[i] >>= 1
		out[i] |= out[i+1] << (wordBits - 1)
	}

	out[scalarWords-1] >>= 1
	out[scalarWords-1] |= carry << (wordBits - 1)
}
