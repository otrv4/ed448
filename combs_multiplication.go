package ed448

func scheduleScalarForCombs(schedule []uint32, scalar [scalarWords]uint32) {
	table := baseTable
	tmp := make([]uint32, len(schedule))
	copy(tmp, scalar[:])

	tmp[len(tmp)-1] &= (uint32(1) << (scalarBits % wordBits)) - 1

	convertToSignedWindowForm(schedule, tmp, table.adjustments[:])
}

func convertToSignedWindowForm(out []uint32, scalar []uint32, preparedData []uint32) {
	mask := uint32(uint64(-(scalar[0] & 1)) & lmask)

	carry := addExtPacked(out, scalar, preparedData[:scalarWords], uint32(^mask))
	carry += addExtPacked(out, out, preparedData[scalarWords:], uint32(mask))

	for i := 0; i < scalarWords-1; i++ {
		out[i] >>= 1
		out[i] |= out[i+1] << (wordBits - 1)
	}

	out[scalarWords-1] >>= 1
	out[scalarWords-1] |= carry << (wordBits - 1)
}
