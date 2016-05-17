package ed448

//Muti-word (double-length) arithmetic
//This is adapted from https://golang.org/src/math/big/arith.go
type Word uintptr

const (
	// Compute the size _S of a Word in bytes.
	_m    = ^Word(0)
	_logS = _m>>8&1 + _m>>16&1 + _m>>32&1
	_S    = 1 << _logS

	_W = _S << 3 // word size in bits
	_B = 1 << _W // digit base
	_M = _B - 1  // digit mask

	_W2 = _W / 2   // half word size in bits
	_B2 = 1 << _W2 // half digit base
	_M2 = _B2 - 1  // half digit mask
)

// z1<<_W + z0 = x*y
// Adapted from Warren, Hacker's Delight, p. 132.
func mulWW_g(x, y Word) (z1, z0 Word) {
	x0 := x & _M2
	x1 := x >> _W2
	y0 := y & _M2
	y1 := y >> _W2
	w0 := x0 * y0
	t := x1*y0 + w0>>_W2
	w1 := t & _M2
	w2 := t >> _W2
	w1 += x0 * y1
	z1 = x1*y1 + w2 + w1>>_W2
	z0 = x * y
	return
}

// z1<<_W + z0 = x*y + c
func mulAddWWW_g(x, y, c Word) (z1, z0 Word) {
	z1, zz0 := mulWW_g(x, y)
	if z0 = zz0 + c; z0 < zz0 {
		z1++
	}
	return
}

// z1<<_W + z0 = x-y-c, with c == 0 or 1
func subWW_g(x, y, c Word) (z1, z0 Word) {
	yc := y + c
	z0 = x - yc
	if z0 > x || yc < y {
		z1 = 1
	}
	return
}

type DWord struct {
	h Word
	l Word
}

//z = x - y
func subDWord(x, y DWord) (z DWord) {
	z.l = x.l - y.l
	// see "Hacker's Delight", section 2-12 (overflow detection)
	c := (y.l&^x.l | (y.l|^x.l)&z.l) >> (_W - 1)
	z.h = x.h - y.h - c

	return
}

// z = a + b
func addDWord(a, b DWord) (z DWord) {
	c := Word(0)
	z.l = a.l + b.l
	if z.l < a.l {
		c++
	}

	z.h = a.h + b.h + c

	//XXX Should we panic on underflow?
	//if z.h < a.h { panic() }

	return
}

//XXX Everything from here is specific to amd64 architecture
//Should be moved to an architecture-specific file

type limb Word
type bigNumber [Limbs]limb
type serialized [56]byte

func mustDeserialize(in serialized) bigNumber {
	n, ok := deserialize(in)
	if !ok {
		panic("Failed to deserialize")
	}

	return n
}