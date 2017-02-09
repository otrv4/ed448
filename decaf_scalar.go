package ed448

import (
	"errors"
)

type decafScalar [scalarWords]word

func (s *decafScalar) copy() *decafScalar {
	out := &decafScalar{}
	copy(out[:], s[:])
	return out
}

func (s *decafScalar) scalarEquals(x *decafScalar) bool {
	diff := word(0)
	for i := uintZero; i < scalarWords; i++ {
		diff |= s[i] ^ x[i]
	}
	return word(((dword(diff))-1)>>wordBits) == decafTrue
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
		chain += sdword(out[i]) + (sdword(ScalarQ[i]) & borrow)
		out[i] = word(chain)
		chain >>= wordBits
	}
	copy(s[:], out[:])
}

func (s *decafScalar) scalarAdd(a, b *decafScalar) {
	out := &decafScalar{}
	var chain dword

	for i := uintZero; i < scalarWords; i++ {
		chain += dword(a[i]) + dword(b[i])
		out[i] = word(chain)
		chain >>= wordBits
	}
	out.scalarSubExtra(out, ScalarQ, word(chain))
	copy(s[:], out[:])
}

func (s *decafScalar) scalarSub(x, y *decafScalar) {
	noExtra := word(0)
	s.scalarSubExtra(x, y, noExtra)
}

func (s *decafScalar) scalarMul(x, y *decafScalar) {
	s.montgomeryMultiply(x, y)
	s.montgomeryMultiply(s, scalarR2)
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

// Serializes an array of words into an array of bytes (little-endian)
func (s *decafScalar) serialize(dst []byte) error {
	wordBytes := wordBits / 8
	if len(dst) < fieldBytes {
		return errors.New("dst length smaller than fieldBytes")
	}

	for i := 0; i*wordBytes < fieldBytes; i++ {
		for j := 0; j < wordBytes; j++ {
			b := s[i] >> uint(8*j)
			dst[wordBytes*i+j] = byte(b)
		}
	}
	return nil
}

func (s *decafScalar) decodeShort(b []byte, size uint) {
	k := uint(0)
	for i := uint(0); i < scalarLimbs; i++ {
		out := word(0)
		for j := uint(0); j < 4 && k < size; j, k = j+1, k+1 {
			out |= (word(b[k])) << (8 * j)
		}
		s[i] = out
	}
}

func (s *decafScalar) decode(b []byte) word {
	s.decodeShort(b, scalarBytes)

	accum := int64(0)
	for i := 0; i < 14; i++ {
		accum += int64(s[i]) - int64(ScalarQ[i])
		accum >>= wordBits
	}

	s.scalarMul(s, &decafScalar{0x01})

	return word(accum)
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
			chain += dword(multiplicand)*dword(ScalarQ[j]) + dword(out[j])
			if j > 0 {
				out[j-1] = word(chain)
			}
			chain >>= wordBits
		}
		chain += dword(saved) + dword(carry)
		out[scalarWords-1] = word(chain)
		carry = word(chain >> wordBits)
	}
	out.scalarSubExtra(out, ScalarQ, carry)
	copy(s[:], out[:])
}

//Exported methods

// NewDecafScalar returns a decaf Scalar in Ed448
func NewDecafScalar(in [fieldBytes]byte) Scalar {
	out := &decafScalar{}
	barrettDeserializeAndReduce(out[:], in[:], &curvePrimeOrder)
	return out
}

//Equals compares two scalars. Returns true if they are the same; false, otherwise.
func (s *decafScalar) Equals(x Scalar) bool {
	return s.scalarEquals(x.(*decafScalar))
}

//Copy copies scalars.
func (s *decafScalar) Copy() Scalar {
	out := &decafScalar{}
	copy(out[:], s[:])
	return out
}

//Add adds two scalars. The scalar may use the same memory.
func (s *decafScalar) Add(x, y Scalar) {
	s.scalarAdd(x.(*decafScalar), y.(*decafScalar))
}

//Sub subtracts two scalars. The scalar may use the same memory.
func (s *decafScalar) Sub(x, y Scalar) {
	noExtra := word(0)
	s.scalarSubExtra(x.(*decafScalar), y.(*decafScalar), noExtra)
}

//Mul multiplies two scalars. The scalar may use the same memory.
func (s *decafScalar) Mul(x, y Scalar) {
	s.montgomeryMultiply(x.(*decafScalar), y.(*decafScalar))
	s.montgomeryMultiply(s, scalarR2)
}

func (s *decafScalar) halve(a, b Scalar) {
	s.scalarHalve(a.(*decafScalar), b.(*decafScalar))
}

//Encode serializes a scalar to wire format.
func (s *decafScalar) Encode() []byte {
	dst := make([]byte, fieldBytes)
	s.serialize(dst)
	return dst
}

//Decode reads a scalar from wire format or from bytes and reduces mod scalar prime.
func (s *decafScalar) Decode(src []byte) error {
	if len(src) < fieldBytes {
		return errors.New("ed448: cannot decode a scalar from a byte array with a length unequal to 56")
	}
	barrettDeserializeAndReduce(s[:], src, &curvePrimeOrder)
	return nil
}
