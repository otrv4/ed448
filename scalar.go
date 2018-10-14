package ed448

import (
	"errors"
)

// Scalar is a interface of a Ed448 scalar
type Scalar interface {
	Equals(a Scalar) bool
	EqualsMask(a Scalar) uint32
	Copy() Scalar
	Add(a, b Scalar)
	Sub(a, b Scalar)
	Mul(a, b Scalar)
	Halve(a Scalar)
	Invert() bool
	Encode() []byte
	BarretDecode(src []byte) error
	Decode(src []byte)
}

type scalar [scalarWords]word

func (s *scalar) montgomeryMultiply(x, y *scalar) {
	out := &scalar{}
	carry := word(0x00)

	for i := 0; i < scalarWords; i++ {
		chain := dword(0x00)
		for j := 0; j < scalarWords; j++ {
			chain += dword(x[i])*dword(y[j]) + dword(out[j])
			out[j] = word(chain)
			chain >>= wordBits
		}

		saved := word(chain)
		multiplicand := out[0] * montgomeryFactor
		chain = dword(0x00)

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

	out.subExtra(out, ScalarQ, carry)
	copy(s[:], out[:])
}

func (s *scalar) montgomerySquare(x *scalar) {
	s.montgomeryMultiply(x, x)
}

// Invert a scalar: 1/s
func (s *scalar) invert() bool {
	// Fermat's little theorem, sliding window algorithm
	preComp := make([]*scalar, 8)
	for i := range preComp {
		preComp[i] = new(scalar)
	}

	out := &scalar{}
	scalarWindowBits := uint(3)
	last := (1 << scalarWindowBits) - 1

	// Precompute preCmp = [a^1, a^3, ...]
	preComp[0].montgomeryMultiply(s, scalarR2)

	if last > 0 {
		preComp[last].montgomeryMultiply(preComp[0], preComp[0])
	}

	for i := 1; i <= last; i++ {
		preComp[i].montgomeryMultiply(preComp[i-1], preComp[last])
	}

	// Sliding window
	var residue, trailing, started uint
	for i := scalarBits - 1; i >= int(-scalarWindowBits); i-- {
		if started != 0 {
			out.montgomerySquare(out)
		}

		var w word
		if i >= 0 {
			w = ScalarQ[i/wordBits]
		} else {
			w = 0x00
		}

		if i >= 0 && i < int(wordBits) {
			w -= 2
		}
		residue = uint((word(residue) << 1) | ((w >> (uint(i) % wordBits)) & 1))
		if residue>>scalarWindowBits != 0 {
			trailing = residue
			residue = 0
		}

		if trailing > 0 && (trailing&((1<<scalarWindowBits)-1)) == 0 {
			if started != 0 {
				out.montgomeryMultiply(out, preComp[trailing>>(scalarWindowBits+1)])
			} else {
				out = preComp[trailing>>(scalarWindowBits+1)].copy()
				started = 1
			}

			trailing = 0
		}
		trailing <<= 1
	}

	// demontgomerize
	one := &scalar{0x01}
	out.montgomeryMultiply(out, one)

	copy(s[:], out[:])

	// TODO: memzero
	// True is the output is not zero
	return out.equals(scalarZero) != true
}

func (s *scalar) equals(x *scalar) bool {
	return word(s.equalsMask(x)) == decafTrue
}

func (s *scalar) equalsMask(x *scalar) uint32 {
	diff := word(0x00)
	for i := uintZero; i < scalarWords; i++ {
		diff |= s[i] ^ x[i]
	}
	return uint32(isZeroMask(diff))
}

func (s *scalar) copy() *scalar {
	out := &scalar{}
	copy(out[:], s[:])
	return out
}

func (s *scalar) set(w word) {
	s[0] = w
}

// {minuend , accum} - subtrahend + (one or mor) q
// Must have carry <= 1
func (s *scalar) subExtra(minuend *scalar, subtrahend *scalar, carry word) {
	out := &scalar{}
	var chain sdword

	for i := uintZero; i < scalarWords; i++ {
		chain += sdword(minuend[i]) - sdword(subtrahend[i])
		out[i] = word(chain)
		chain >>= wordBits
	}

	borrow := chain + sdword(carry) // 0 or -1
	chain = 0

	for i := uintZero; i < scalarWords; i++ {
		chain += sdword(out[i]) + (sdword(ScalarQ[i]) & borrow)
		out[i] = word(chain)
		chain >>= wordBits
	}
	copy(s[:], out[:])
}

func (s *scalar) add(a, b *scalar) {
	var chain dword

	for i := uintZero; i < scalarWords; i++ {
		chain += dword(a[i]) + dword(b[i])
		s[i] = word(chain)
		chain >>= wordBits
	}
	s.subExtra(s, ScalarQ, word(chain))
}

func (s *scalar) sub(x, y *scalar) {
	noExtra := word(0x00)
	s.subExtra(x, y, noExtra)
}

func (s *scalar) mul(x, y *scalar) {
	s.montgomeryMultiply(x, y)
	s.montgomeryMultiply(s, scalarR2)
}

func (s *scalar) halve(a *scalar) {
	mask := -(a[0] & 1)
	var chain dword
	var i uint

	for i = uintZero; i < scalarWords; i++ {
		chain += dword(a[i]) + dword(ScalarQ[i]&mask)
		s[i] = word(chain)
		chain >>= wordBits
	}
	for i = uintZero; i < scalarWords-1; i++ {
		s[i] = s[i]>>1 | s[i+1]<<(wordBits-1)
	}

	s[i] = s[i]>>1 | word(chain<<(wordBits-1))
}

// Serializes an array of words into an array of bytes (little-endian)
func (s *scalar) encode(dst []byte) error {
	wordBytes := wordBits / 8
	if len(dst) < fieldBytes {
		return errors.New("dst length smaller than fieldBytes")
	}

	k := uintZero
	for i := uintZero; i < scalarLimbs; i++ {
		for j := uintZero; j < uint(wordBytes); j++ {
			b := s[i] >> (8 * j)
			dst[k] = byte(b)
			k++
		}
	}
	return nil
}

func (s *scalar) decodeShort(b []byte, size uint) {
	k := uintZero
	for i := uintZero; i < scalarLimbs; i++ {
		out := word(0)
		for j := uintZero; j < 4 && k < size; j, k = j+1, k+1 {
			out |= (word(b[k])) << (8 * j)
		}
		s[i] = out
	}
}

func (s *scalar) decode(b []byte) word {
	s.decodeShort(b, scalarBytes)

	accum := sdword(0x00)
	for i := 0; i < 14; i++ {
		accum += sdword(s[i]) - sdword(ScalarQ[i])
		accum >>= wordBits
	}

	s.mul(s, &scalar{0x01})

	return word(accum)
}

// HACKY: either the param or the return
func decodeLong(s *scalar, b []byte) *scalar {
	y := &scalar{}
	bLen := len(b)
	size := bLen - (bLen % scalarSerBytes)

	if bLen == 0 {
		s = scalarZero.copy()
		return s
	}

	if size == bLen {
		size -= scalarSerBytes
	}
	s.decodeShort(b[size:], uint(bLen-size))

	if bLen == scalarBytes {
		s.mul(s, &scalar{0x01})
		return s
	}

	for size > 0 {
		size -= scalarSerBytes
		s.montgomeryMultiply(s, scalarR2)
		y.decode(b[size:])
		s.add(s, y)
	}

	y.destroy()
	return s.copy()
}

func (s *scalar) destroy() {
	copy(s[:], scalarZero[:])
}

//Exported methods

// NewScalar returns a Scalar in Ed448 with decaf
func NewScalar(in ...[]byte) Scalar {
	if len(in) > 1 {
		panic("too many arguments to function call")
	}

	if in == nil {
		return &scalar{}
	}

	out := &scalar{}

	bytes := in[0][:]
	return decodeLong(out, bytes)
}

// Equals compares two scalars. Returns true if they are the same; false, otherwise.
func (s *scalar) Equals(x Scalar) bool {
	return s.equals(x.(*scalar))
}

// EqualsMask compares two scalars.
func (s *scalar) EqualsMask(x Scalar) uint32 {
	return s.equalsMask(x.(*scalar))
}

// Copy copies scalars.
func (s *scalar) Copy() Scalar {
	out := &scalar{}
	copy(out[:], s[:])
	return out
}

// Add adds two scalars. The scalars may use the same memory.
func (s *scalar) Add(x, y Scalar) {
	s.add(x.(*scalar), y.(*scalar))
}

// Sub subtracts two scalars. The scalars may use the same memory.
func (s *scalar) Sub(x, y Scalar) {
	noExtra := word(0)
	s.subExtra(x.(*scalar), y.(*scalar), noExtra)
}

// Mul multiplies two scalars. The scalars may use the same memory.
func (s *scalar) Mul(x, y Scalar) {
	s.montgomeryMultiply(x.(*scalar), y.(*scalar))
	s.montgomeryMultiply(s, scalarR2)
}

// Halve halfs a scalar. The scalars may used the same memory.
func (s *scalar) Halve(x Scalar) {
	s.halve(x.(*scalar))
}

// Invert inverts a scalar. The scalars may used the same memory.
func (s *scalar) Invert() bool {
	return s.invert()
}

// Encode serializes a scalar to wire format.
func (s *scalar) Encode() []byte {
	dst := make([]byte, fieldBytes)
	s.encode(dst)
	return dst
}

// Decode reads a scalar from wire format or from bytes and reduces mod scalar prime.
// TODO: this will reduce with barret, change name and receiver
func (s *scalar) BarretDecode(src []byte) error {
	if len(src) < fieldBytes {
		return errors.New("ed448: cannot decode a scalar from a byte array with a length unequal to 56")
	}
	barrettDeserializeAndReduce(s[:], src, &curvePrimeOrder)
	return nil
}

// Decode reads a scalar from wire format or from bytes and reduces mod scalar prime.
func (s *scalar) Decode(src []byte) {
	decodeLong(s, src)
}
