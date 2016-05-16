package ed448

const (
	Limbs = 8
	Radix = 56
)

type word uint64
type limb word
type bigNumber [Limbs]limb
type serialized [Radix]byte

func deserialize(in serialized) (n bigNumber, ok bool) {
	const (
		columns = Limbs
		rows    = Limbs - 1
	)

	for i := uint(0); i < columns; i++ {
		for j := uint(0); j < rows; j++ {
			n[i] |= limb(in[rows*i+j]) << (columns * j)
		}
	}

	ok = !constantTimeGreaterOrEqualP(n)

	return
}

//TODO: Make this work with a word parameter
func isZero(n int64) int64 {
	return ^n
}

func constantTimeGreaterOrEqualP(n bigNumber) bool {
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

func serialize(dst []byte, src bigNumber) {
	const (
		rows    = Limbs
		columns = Radix / Limbs
	)

	var n bigNumber
	copy(n[:], src[:])

	for i := uint(0); i < rows; i++ {
		for j := uint(0); j < columns; j++ {
			dst[columns*i+j] = byte(n[i])
			n[i] >>= 8
		}
	}
}

/*This is adapted from https://golang.org/src/math/big/arith.go */
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

func WideMul(x, y Word) DWord {
	z1, z0 := mulWW_g(x, y)
	return DWord{h: z1, l: z0}
}

func Mac(x, y Word, acc DWord) DWord {
	z1, z0 := mulAddWWW_g(x, y, acc.l)

	z11 := z1 + acc.h
	if z11 < z1 {
		panic("high word overflow")
	}

	return DWord{h: z11, l: z0}
}

//z = a - b
func SubDWord(a, b DWord) (z DWord, c Word) {
	xi := a.l
	yi := b.l
	zi := xi - yi

	z.l = zi
	// see "Hacker's Delight", section 2-12 (overflow detection)
	c = (yi&^xi | (yi|^xi)&zi) >> (_W - 1)

	xi = a.h
	yi = b.h
	zi = xi - yi - c

	z.h = zi
	// see "Hacker's Delight", section 2-12 (overflow detection)
	c = (yi&^xi | (yi|^xi)&zi) >> (_W - 1)

	return
}

// z = a + b
func AddDWord(a, b DWord) (z DWord) {
	c := Word(0)
	z.l = a.l + b.l
	if z.l < a.l {
		c++
	}

	z.h = a.h + b.h + c

	if z.h < a.h {
		panic("high word overflow")
	}

	return
}

func Msb(a, b Word, acc DWord) (DWord, Word) {
	d := WideMul(a, b)
	return SubDWord(d, acc)
}

func wideMul(x, y limb) DWord {
	return WideMul(Word(x), Word(y))
}

func mac(x, y limb, acc DWord) DWord {
	return Mac(Word(x), Word(y), acc)
}

func msb(x, y limb, acc DWord) (DWord, Word) {
	return Msb(Word(x), Word(y), acc)
}

func karatsubaMul(a, b bigNumber) (c bigNumber) {
	var aa, bb, bbb [4]limb

	mask := Word(0xffffffffffffff)

	for i := 0; i < 4; i++ {
		aa[i] = a[i] + a[i+4]
		bb[i] = b[i] + b[i+4]
		bbb[i] = bb[i] + b[i+4]
	}

	accum2 := wideMul(a[0], b[3])
	accum0 := wideMul(aa[0], bb[3])
	accum1 := wideMul(a[4], b[7])

	accum2 = mac(a[1], b[2], accum2)
	accum0 = mac(aa[1], bb[2], accum0)
	accum1 = mac(a[5], b[6], accum1)

	accum2 = mac(a[2], b[1], accum2)
	accum0 = mac(aa[2], bb[1], accum0)
	accum1 = mac(a[6], b[5], accum1)

	accum2 = mac(a[3], b[0], accum2)
	accum0 = mac(aa[3], bb[0], accum0)
	accum1 = mac(a[7], b[4], accum1)

	//If borrow != 0, we should panic?
	accum0, _ = SubDWord(accum0, accum2)
	accum1 = AddDWord(accum1, accum2)

	c[3] = limb(accum1.l & mask)
	c[7] = limb(accum0.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum0 = mac(aa[1], bb[3], accum0)
	accum1 = mac(a[5], b[7], accum1)
	accum0 = mac(aa[2], bb[2], accum0)
	accum1 = mac(a[6], b[6], accum1)
	accum0 = mac(aa[3], bb[1], accum0)

	accum1 = AddDWord(accum1, accum0)

	accum2 = wideMul(a[0], b[0])
	accum1, _ = SubDWord(accum1, accum2)
	accum0 = AddDWord(accum0, accum2)

	accum0, _ = msb(a[1], b[3], accum0)
	accum0, _ = msb(a[2], b[2], accum0)
	accum1 = mac(a[7], b[5], accum1)
	accum0, _ = msb(a[3], b[1], accum0)
	accum1 = mac(aa[0], bb[0], accum1)
	accum0 = mac(a[4], b[4], accum0)

	//c[3+i], c[3+i mod 7]
	c[0] = limb(accum0.l & mask) // THIS IS WRONG
	c[4] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum2 = wideMul(a[2], b[7])
	accum0 = mac(a[6], bb[3], accum0)
	accum1 = mac(aa[2], bbb[3], accum1)

	accum2 = mac(a[3], b[6], accum2)
	accum0 = mac(a[7], bb[2], accum0)
	accum1 = mac(aa[3], bbb[2], accum1)

	accum2 = mac(a[0], b[1], accum2)
	accum1 = mac(aa[0], bb[1], accum1)
	accum0 = mac(a[4], b[5], accum0)

	accum2 = mac(a[1], b[0], accum2)
	accum1 = mac(aa[1], bb[0], accum1)
	accum0 = mac(a[5], b[4], accum0)

	accum1, _ = SubDWord(accum1, accum2)
	accum0 = AddDWord(accum0, accum2)

	c[1] = limb(accum0.l & mask) // THIS IS WRONG
	c[5] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum2 = wideMul(a[3], b[7])
	accum0 = mac(a[7], bb[3], accum0)
	accum1 = mac(aa[3], bbb[3], accum1)

	accum2 = mac(a[0], b[2], accum2)
	accum1 = mac(aa[0], bb[2], accum1)
	accum0 = mac(a[4], b[6], accum0)

	accum2 = mac(a[1], b[1], accum2)
	accum1 = mac(aa[1], bb[1], accum1)
	accum0 = mac(a[5], b[5], accum0)

	accum2 = mac(a[2], b[0], accum2)
	accum1 = mac(aa[2], bb[0], accum1)
	accum0 = mac(a[6], b[4], accum0)

	accum1, _ = SubDWord(accum1, accum2)
	accum0 = AddDWord(accum0, accum2)

	c[2] = limb(accum0.l & mask) // THIS IS WRONG
	c[6] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	accum0 = AddDWord(accum0, DWord{0, Word(c[3])})
	accum1 = AddDWord(accum1, DWord{0, Word(c[7])})

	c[3] = limb(accum0.l & mask)
	c[7] = limb(accum1.l & mask)

	accum0.l >>= 56
	accum0.l |= (accum0.h << 8)
	accum0.h >>= 56

	accum1.l >>= 56
	accum1.l |= (accum1.h << 8)
	accum1.h >>= 56

	c[0] += limb(accum1.l)
	c[4] += limb(accum0.l + accum1.l)

	return
}
