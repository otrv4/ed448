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
	var chain dword_t
	var i uint

	for i = 0; i < scalarWords; i++ {
		chain += dword_t(a[i]) + dword_t(b[i]&mask)
		out[i] = word_t(chain)
		chain >>= wordBits
	}
	for i = 0; i < scalarWords-1; i++ {
		out[i] = out[i]>>1 | out[i+1]<<(wordBits-1)
	}

	out[i] = out[i]>>1 | word_t(chain<<(wordBits-1))

	return
}
