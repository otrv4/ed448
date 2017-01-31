package ed448

type scalar32 [scalarWords]uint32

// Serializes an array of words into an array of bytes (little-endian)
func (s *scalar32) serialize(dst []byte) {
	wordBytes := wordBits / 8

	for i := 0; i*wordBytes < len(dst); i++ {
		for j := 0; j < wordBytes; j++ {
			b := s[i] >> uint(8*j)
			dst[wordBytes*i+j] = byte(b)
		}
	}
}

func (s *scalar32) scalarAdd(a, b *scalar32) {
	out := &scalar32{}
	var chain uint64

	for i := uintZero; i < scalarWords; i++ {
		chain += uint64(a[i]) + uint64(b[i])
		out[i] = uint32(chain)
		chain >>= wordBits
	}
	out.scalarSubExtra(out, scalarQ, uint32(chain))
	copy(s[:], out[:])
}

func (s *scalar32) scalarSubExtra(minuend *scalar32, subtrahend *scalar32, carry uint32) {
	out := &scalar32{}
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
	copy(s[:], out[:])
}

func (s *scalar32) scalarHalve(a, b *scalar32) {
	out := &scalar32{}
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

	copy(s[:], out[:])
}

func (s *scalar32) montgomeryMultiply(x, y *scalar32) {
	out := &scalar32{}
	carry := uint32(0)

	for i := 0; i < scalarWords; i++ {
		chain := uint64(0)
		for j := 0; j < scalarWords; j++ {
			chain += uint64(x[i])*uint64(y[j]) + uint64(out[j])
			out[j] = uint32(chain)
			chain >>= wordBits
		}
		saved := uint32(chain)
		multiplicand := out[0] * montgomeryFactor
		chain = 0
		for j := 0; j < scalarWords; j++ {
			chain += uint64(multiplicand)*uint64(scalarQ[j]) + uint64(out[j])
			if j > 0 {
				out[j-1] = uint32(chain)
			}
			chain >>= wordBits
		}
		chain += uint64(saved) + uint64(carry)
		out[scalarWords-1] = uint32(chain)
		carry = uint32(chain >> wordBits)
	}
	out.scalarSubExtra(out, scalarQ, carry)
	copy(s[:], out[:])
}

func (s *scalar32) Decode(serial []byte) {
	barrettDeserializeAndReduce(s[:], serial, &curvePrimeOrder)
}

func (s *scalar32) Encode(dst []byte) {
	s.serialize(dst)
}

func (s *scalar32) Copy() Scalar {
	out := &scalar32{}
	copy(out[:], s[:])
	return out
}
