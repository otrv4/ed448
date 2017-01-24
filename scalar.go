package ed448

func scalarAdd(a, b [scalarWords]uint32) (out [scalarWords]uint32) {
	var chain uint64

	for i := uintZero; i < scalarWords; i++ {
		chain += uint64(a[i]) + uint64(b[i])
		out[i] = uint32(chain)
		chain >>= wordBits
	}

	return scalarSubExtra(out, scalarQ, uint32(chain))
}

func scalarSubExtra(accum, sub [scalarWords]uint32, extra uint32) (out [scalarWords]uint32) {
	var chain int64

	for i := uintZero; i < scalarWords; i++ {
		chain += int64(accum[i]) - int64(sub[i])
		out[i] = uint32(chain)
		chain >>= wordBits
	}

	borrow := chain + int64(extra)
	chain = 0

	for i := uintZero; i < scalarWords; i++ {
		chain += int64(out[i]) + (int64(scalarQ[i]) & borrow)
		out[i] = uint32(chain)
		chain >>= wordBits
	}
	return
}

func scalarHalve(a, b [scalarWords]uint32) (out [scalarWords]uint32) {
	mask := -(a[0] & 1)
	var chain uint64
	var i uint

	for i = 0; i < scalarWords; i++ {
		chain += uint64(a[i]) + uint64(b[i]&mask)
		out[i] = uint32(chain)
		chain >>= wordBits
	}
	for i = 0; i < scalarWords-1; i++ {
		out[i] = out[i]>>1 | out[i+1]<<(wordBits-1)
	}

	out[i] = out[i]>>1 | uint32(chain<<(wordBits-1))

	return
}
