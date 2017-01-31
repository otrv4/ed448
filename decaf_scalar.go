package ed448

import "fmt"

type decafScalar [scalarWords]word

// Serializes an array of words into an array of bytes (little-endian)
func (s *decafScalar) serialize(dst []byte) error {
	wordBytes := wordBits / 8
	if len(dst) < fieldBytes {
		return fmt.Errorf("dst length smaller than fieldBytes")
	}

	for i := 0; i*wordBytes < fieldBytes; i++ {
		for j := 0; j < wordBytes; j++ {
			b := s[i] >> uint(8*j)
			dst[wordBytes*i+j] = byte(b)
		}
	}
	return nil
}

func (s *decafScalar) scalarAdd(a, b *decafScalar) {
	out := &decafScalar{}
	var chain dword

	for i := uintZero; i < scalarWords; i++ {
		chain += dword(a[i]) + dword(b[i])
		out[i] = word(chain)
		chain >>= wordBits
	}
	out.scalarSubExtra(out, scalarQ, word(chain))
	copy(s[:], out[:])
}

// unexposed methods for decafSign as an internal method
// XXX: remove these methods when what is needed is exposed
func (s *decafScalar) scalarSub(x, y *decafScalar) {
	noExtra := word(0)
	s.scalarSubExtra(x, y, noExtra)
}

// XXX: remove these methods when what is needed is exposed
func (s *decafScalar) scalarMul(x, y *decafScalar) {
	s.montgomeryMultiply(x, y)
	s.montgomeryMultiply(s, scalarR2)
}

func (s *decafScalar) scalarSubExtra(minuend *decafScalar, subtrahend *decafScalar, carry word) {
	out := &decafScalar{}
	var chain sdword

	for i := uintZero; i < scalarWords; i++ {
		chain += sdword(minuend[i]) - sdword(subtrahend[i])
		out[i] = word(chain)
		chain >>= wordBits
	}

	borrow := chain + sdword(carry)
	chain = 0

	for i := uintZero; i < scalarWords; i++ {
		chain += sdword(out[i]) + (sdword(scalarQ[i]) & borrow)
		out[i] = word(chain)
		chain >>= wordBits
	}
	copy(s[:], out[:])
}

func (s *decafScalar) scalarHalve(a, b *decafScalar) {
	out := &decafScalar{}
	mask := -(a[0] & 1)
	var chain dword
	var i uint

	for i = 0; i < scalarWords; i++ {
		chain += dword(a[i]) + dword(b[i]&mask)
		out[i] = word(chain)
		chain >>= wordBits
	}
	for i = 0; i < scalarWords-1; i++ {
		out[i] = out[i]>>1 | out[i+1]<<(wordBits-1)
	}

	out[i] = out[i]>>1 | word(chain<<(wordBits-1))

	copy(s[:], out[:])
}

func (s *decafScalar) montgomeryMultiply(x, y *decafScalar) {
	out := &decafScalar{}
	carry := word(0)

	for i := 0; i < scalarWords; i++ {
		chain := dword(0)
		for j := 0; j < scalarWords; j++ {
			chain += dword(x[i])*dword(y[j]) + dword(out[j])
			out[j] = word(chain)
			chain >>= wordBits
		}
		saved := word(chain)
		multiplicand := out[0] * montgomeryFactor
		chain = 0
		for j := 0; j < scalarWords; j++ {
			chain += dword(multiplicand)*dword(scalarQ[j]) + dword(out[j])
			if j > 0 {
				out[j-1] = word(chain)
			}
			chain >>= wordBits
		}
		chain += dword(saved) + dword(carry)
		out[scalarWords-1] = word(chain)
		carry = word(chain >> wordBits)
	}
	out.scalarSubExtra(out, scalarQ, carry)
	copy(s[:], out[:])
}

func (s *decafScalar) scalarEquals(x *decafScalar) word {
	diff := word(0)
	for i := uintZero; i < scalarWords; i++ {
		diff |= s[i] ^ x[i]
	}
	return word(((dword(diff)) - 1) >> wordBits)
}

func (s *decafScalar) halve(a, b Scalar) {
	s.scalarHalve(a.(*decafScalar), b.(*decafScalar))
}

func (s *decafScalar) Decode(src []byte) error {
	if len(src) < fieldBytes {
		return fmt.Errorf("src length smaller than fieldBytes")
	}
	barrettDeserializeAndReduce(s[:], src, &curvePrimeOrder)
	return nil
}

func (s *decafScalar) Encode(dst []byte) error {
	return s.serialize(dst)
}

func (s *decafScalar) Copy() Scalar {
	out := &decafScalar{}
	copy(out[:], s[:])
	return out
}

// NewDecafScalar returns a decaf Scalar in Ed448 depending on the arch
func NewDecafScalar(in [fieldBytes]byte) Scalar {
	out := &decafScalar{}
	barrettDeserializeAndReduce(out[:], in[:], &curvePrimeOrder)
	return out
}
