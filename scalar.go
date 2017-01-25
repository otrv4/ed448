package ed448

func scalarAdd(a, b [scalarWords]uint32) (out [scalarWords]uint32) {
	var chain uint64

	for i := uintZero; i < scalarWords; i++ {
		chain += uint64(a[i]) + uint64(b[i])
		out[i] = uint32(chain)
		chain >>= wordBits
	}

	return scalarSubExtra(out[:], scalarQ, uint32(chain))
}

func scalarSubExtra(minuend []uint32, subtrahend [scalarWords]uint32, carry uint32) (out [scalarWords]uint32) {
	var chain int64

	for i := uintZero; i < scalarWords; i++ {
		chain += int64(minuend[i]) - int64(subtrahend[i])
		out[i] = uint32(chain)
		chain >>= wordBits
	}

	borrow := chain + int64(carry)
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

func scalarMul(x, y [scalarWords]uint32) [scalarWords]uint32 {
	out := [scalarWords + 1]uint32{0}
	carry := uint32(0)

	for i := 0; i < scalarWords; i++ {
		chain := uint64(0)
		for j := 0; j < scalarWords; j++ {
			chain += uint64(x[i])*uint64(y[j]) + uint64(out[j])
			out[j] = uint32(chain)
			chain >>= wordBits
		}
		out[scalarWords] = uint32(chain)
		multiplicand := out[0] * montgomeryFactor
		chain = 0
		for j := 0; j < scalarWords; j++ {
			chain += uint64(multiplicand)*uint64(scalarQ[j]) + uint64(out[j])
			if j > 0 {
				out[j-1] = uint32(chain)
			}
			chain >>= wordBits
		}
		chain += uint64(out[scalarWords]) + uint64(carry)
		out[scalarWords-1] = uint32(chain)
		carry = uint32(chain >> wordBits)
	}
	return scalarSubExtra(out[:], scalarQ, carry)
}
