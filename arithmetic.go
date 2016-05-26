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

	//XXX Should we panic on overflow?
	//if z.h < a.h { panic() }

	return
}

// The resulting carry c is either 0 or 1.
func addVV_g(z, x, y []Word) (c Word) {
	for i, xi := range x[:len(z)] {
		yi := y[i]
		zi := xi + yi + c
		z[i] = zi
		// see "Hacker's Delight", section 2-12 (overflow detection)
		c = (xi&yi | (xi|yi)&^zi) >> (_W - 1)
	}
	return
}

// The resulting carry c is either 0 or 1.
func subVV_g(z, x, y []Word) (c Word) {
	for i, xi := range x[:len(z)] {
		yi := y[i]
		zi := xi - yi - c
		z[i] = zi
		// see "Hacker's Delight", section 2-12 (overflow detection)
		c = (yi&^xi | (yi|^xi)&zi) >> (_W - 1)
	}
	return
}

func shlVU_g(z, x []Word, s uint) (c Word) {
	if n := len(z); n > 0 {
		ŝ := _W - s
		w1 := x[n-1]
		c = w1 >> ŝ
		for i := n - 1; i > 0; i-- {
			w := w1
			w1 = x[i-1]
			z[i] = w<<s | w1>>ŝ
		}
		z[0] = w1 << s
	}
	return
}

func shrVU_g(z, x []Word, s uint) (c Word) {
	if n := len(z); n > 0 {
		ŝ := _W - s
		w1 := x[0]
		c = w1 << ŝ
		for i := 0; i < n-1; i++ {
			w := w1
			w1 = x[i+1]
			z[i] = w>>s | w1<<ŝ
		}
		z[n-1] = w1 >> s
	}
	return
}

//XXX Everything from here is specific to amd64 architecture
//Should be moved to an architecture-specific file

type limb Word
type bigNumber [Limbs]limb
type serialized [56]byte

func mustDeserialize(in serialized) *bigNumber {
	n, ok := deserialize(in)
	if !ok {
		panic("Failed to deserialize")
	}

	return n
}

//TODO: Make this work with a word parameter
func isZero(n int64) int64 {
	return ^n
}

func constantTimeGreaterOrEqualP(n *bigNumber) bool {
	var (
		ge   = int64(-1)
		mask = int64(1)<<Radix - 1
	)

	for i := 0; i < 4; i++ {
		ge &= int64(n[i])
	}

	ge = (ge & (int64(n[4]) + 1)) | isZero(int64(n[4])^mask)

	for i := 5; i < 8; i++ {
		ge &= int64(n[i])
	}

	return ge == mask
}

//TODO: should not create a new bigNumber to save memory
func sumRadix(a, b *bigNumber) (c *bigNumber) {
	return a.copy().add(b)
}

//XXX Is there an optimum way of squaring with karatsuba?
func squareRadix(a *bigNumber) (c *bigNumber) {
	return karatsubaMul(a, a)
}

func (n *bigNumber) weakReduce() *bigNumber {
	tmp := limb(uint64(n[Limbs-1]) >> Radix)

	n[Limbs/2] += tmp

	for i := Limbs - 1; i > 0; i-- {
		n[i] = (n[i] & radixMask) + (n[i-1] >> Radix)
	}

	n[0] = (n[0] & radixMask) + tmp

	return n
}

func subRadix(a, b *bigNumber) (c *bigNumber) {
	c = subRadixRaw(a, b)
	c.bias(2)      //???
	c.weakReduce() //???
	return c
}

func subRadixRaw(a, b *bigNumber) (c *bigNumber) {
	c = &bigNumber{}
	for i := 0; i < len(c); i++ {
		c[i] = a[i] - b[i]
	}

	return
}

func (n *bigNumber) copy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

func (n *bigNumber) equals(o *bigNumber) (eq bool) {
	r := limb(0)

	for i, oi := range o {
		r |= n[i] ^ oi
	}

	return r == 0
}

func (n *bigNumber) zero() (eq bool) {
	r := limb(0)

	for _, ni := range n {
		r |= ni ^ 0
	}

	return r == 0
}

func (n *bigNumber) add(x *bigNumber) *bigNumber {
	for i, xi := range x {
		n[i] += xi
	}

	return n
}

func (n *bigNumber) mul(x *bigNumber) *bigNumber {
	for i, mi := range karatsubaMul(n, x) {
		n[i] = mi
	}

	return n
}

//XXX Is there any optimum way of squaring
func (n *bigNumber) square() *bigNumber {
	return n.mul(n)
}
