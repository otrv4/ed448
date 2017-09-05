package ed448

func scheduleScalarForCombs(schedule []word, sc scalar) {
	table := baseTable
	tmp := make([]word, len(schedule))
	copy(tmp, sc[:])

	tmp[len(tmp)-1] &= (word(1) << (scalarBits % wordBits)) - 1

	convertToSignedWindowForm(schedule, tmp, table.adjustments[:])
}

func convertToSignedWindowForm(out []word, scalar []word, preparedData []word) {
	mask := word(dword(-(scalar[0] & 1)) & lmask)

	carry := addExtPacked(out, scalar, preparedData[:scalarWords], word(^mask))
	carry += addExtPacked(out, out, preparedData[scalarWords:], word(mask))

	for i := 0; i < scalarWords-1; i++ {
		out[i] >>= 1
		out[i] |= out[i+1] << (wordBits - 1)
	}

	out[scalarWords-1] >>= 1
	out[scalarWords-1] |= carry << (wordBits - 1)
}
