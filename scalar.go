package ed448

func scalarAdd(a, b [scalarWords]word_t) (out [scalarWords]word_t) {
	var chain word_t

	for i := uintZero; i < scalarWords; i++ {
		chain += a[i] + b[i]
		out[i] = chain
		chain >>= wordBits
	}

	return scalarSubExtra(out, scalarP, chain)
}

func scalarSubExtra(accum, sub [scalarWords]word_t, extra word_t) (out [scalarWords]word_t) {
	var chain int64

	for i := uintZero; i < scalarWords; i++ {
		chain += int64(accum[i]) - int64(sub[i])
		out[i] = word_t(chain)
		chain >>= wordBits
	}

	borrow := chain + int64(extra)
	chain = 0

	for i := uintZero; i < scalarWords; i++ {
		chain += int64(out[i]) + (int64(scalarP[i]) & borrow)
		out[i] = word_t(chain)
		chain >>= wordBits
	}
	return
}

func scalarHalve(a, b [scalarWords]word_t) (out [scalarWords]word_t) {
	mask := -(a[0] & 1)
	var chain word_t

	for i := uintZero; i < scalarWords; i++ {
		chain += a[i] + b[i]&mask
		out[i] = chain
		chain >>= wordBits
	}

	out[0] = out[0]>>1 | out[1]<<(wordBits-1)
	out[1] = out[1]>>1 | out[2]<<(wordBits-1)
	out[2] = out[2]>>1 | out[3]<<(wordBits-1)
	out[3] = out[3]>>1 | out[4]<<(wordBits-1)
	out[4] = out[4]>>1 | out[5]<<(wordBits-1)
	out[5] = out[5]>>1 | out[6]<<(wordBits-1)
	out[6] = out[6]>>1 | out[7]<<(wordBits-1)
	out[7] = out[7]>>1 | out[8]<<(wordBits-1)
	out[8] = out[8]>>1 | out[9]<<(wordBits-1)
	out[9] = out[9]>>1 | out[10]<<(wordBits-1)
	out[10] = out[10]>>1 | out[11]<<(wordBits-1)
	out[11] = out[11]>>1 | out[12]<<(wordBits-1)
	out[12] = out[12]>>1 | out[13]<<(wordBits-1)
	out[13] = out[13]>>1 | chain<<(wordBits-1)

	return
}
