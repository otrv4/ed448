package ed448

import "fmt"

type word uint32
type sword int32
type dword uint64
type sdword int64

type bigNumber [nLimbs]word
type serialized [fieldBytes]byte

func isZeroMask(n word) word {
	nn := dword(n)
	nn = nn - 1
	return word(nn >> wordBits)
}

func (n *bigNumber) isZero() (eq bool) {
	return n.zeroMask() == lmask
}

func (n *bigNumber) copy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

// TODO: delete this
func (n *bigNumber) equals(o *bigNumber) (eq bool) {
	r := word(0)
	x := n.copy().strongReduce()
	y := o.copy().strongReduce()

	r |= x[0] ^ y[0]
	r |= x[1] ^ y[1]
	r |= x[2] ^ y[2]
	r |= x[3] ^ y[3]
	r |= x[4] ^ y[4]
	r |= x[5] ^ y[5]
	r |= x[6] ^ y[6]
	r |= x[7] ^ y[7]
	r |= x[8] ^ y[8]
	r |= x[9] ^ y[9]
	r |= x[10] ^ y[10]
	r |= x[11] ^ y[11]
	r |= x[12] ^ y[12]
	r |= x[13] ^ y[13]
	r |= x[14] ^ y[14]
	r |= x[15] ^ y[15]

	return r == 0
}

// TODO: make this the canonical equals
// Compare n == x
// If it is equal, it will return 0. Otherwise the lmask.
func (n *bigNumber) decafEq(x *bigNumber) word {
	y := &bigNumber{}
	y.sub(n, x)
	y.strongReduce()

	var ret word

	ret |= y[0]
	ret |= y[1]
	ret |= y[2]
	ret |= y[3]
	ret |= y[4]
	ret |= y[5]
	ret |= y[6]
	ret |= y[7]
	ret |= y[8]
	ret |= y[9]
	ret |= y[10]
	ret |= y[11]
	ret |= y[12]
	ret |= y[13]
	ret |= y[14]

	return isZeroMask(ret)
}

func (n *bigNumber) set(x *bigNumber) *bigNumber {
	copy(n[:], x[:])
	return n
}

//in is big endian
func (n *bigNumber) setBytes(in []byte) *bigNumber {
	if len(in) != fieldBytes {
		return nil
	}

	s := serialized{}
	for i, si := range in {
		s[len(s)-i-1] = si
	}

	d, ok := deserialize(s)
	if !ok {
		return nil
	}

	for i, di := range d {
		n[i] = di
	}

	return n
}

func (n *bigNumber) setUI(y dword) *bigNumber {
	n[0] = word(y) & radixMask
	n[1] = word(y >> radix)
	n[2] = 0
	n[3] = 0
	n[4] = 0
	n[5] = 0
	n[6] = 0
	n[7] = 0
	n[8] = 0
	n[9] = 0
	n[10] = 0
	n[11] = 0
	n[12] = 0
	n[13] = 0
	n[14] = 0
	n[15] = 0

	return n
}

func (n *bigNumber) zeroMask() word {
	x := n.copy().strongReduce()
	r := word(0)

	r |= x[0]
	r |= x[1]
	r |= x[2]
	r |= x[3]
	r |= x[4]
	r |= x[5]
	r |= x[6]
	r |= x[7]
	r |= x[8]
	r |= x[9]
	r |= x[10]
	r |= x[11]
	r |= x[12]
	r |= x[13]
	r |= x[14]
	r |= x[15]

	return isZeroMask(word(r))
}

// Return high bit of x = low bit of 2x mod p
func highBit(x *bigNumber) word {
	y := &bigNumber{}
	y.add(x, x)
	y.strongReduce()
	return word(-(y[0] & 1))
}

func lowBit(x *bigNumber) word {
	x.strongReduce()
	return -(x[0] & 1)
}

func (n *bigNumber) bias(b word) *bigNumber {
	co1 := radixMask * b
	co2 := co1 - b
	lo := [4]word{co1, co1, co1, co1}
	hi := [4]word{co2, co1, co1, co1}

	n[0] += lo[0]
	n[1] += lo[1]
	n[2] += lo[2]
	n[3] += lo[3]

	n[4] += lo[0]
	n[5] += lo[1]
	n[6] += lo[2]
	n[7] += lo[3]

	n[8] += hi[0]
	n[9] += hi[1]
	n[10] += hi[2]
	n[11] += hi[3]

	n[12] += lo[0]
	n[13] += lo[1]
	n[14] += lo[2]
	n[15] += lo[3]

	return n
}

func (n *bigNumber) addRaw(x *bigNumber, y *bigNumber) *bigNumber {
	n[0] = x[0] + y[0]
	n[1] = x[1] + y[1]
	n[2] = x[2] + y[2]
	n[3] = x[3] + y[3]
	n[4] = x[4] + y[4]
	n[5] = x[5] + y[5]
	n[6] = x[6] + y[6]
	n[7] = x[7] + y[7]
	n[8] = x[8] + y[8]
	n[9] = x[9] + y[9]
	n[10] = x[10] + y[10]
	n[11] = x[11] + y[11]
	n[12] = x[12] + y[12]
	n[13] = x[13] + y[13]
	n[14] = x[14] + y[14]
	n[15] = x[15] + y[15]

	return n
}

//n = x + y
func (n *bigNumber) add(x *bigNumber, y *bigNumber) *bigNumber {
	return n.addRaw(x, y).weakReduce()
}

func (n *bigNumber) addW(w word) *bigNumber {
	n[0] += word(w)
	return n
}

func (n *bigNumber) subRaw(x *bigNumber, y *bigNumber) *bigNumber {
	n[0] = x[0] - y[0]
	n[1] = x[1] - y[1]
	n[2] = x[2] - y[2]
	n[3] = x[3] - y[3]
	n[4] = x[4] - y[4]
	n[5] = x[5] - y[5]
	n[6] = x[6] - y[6]
	n[7] = x[7] - y[7]
	n[8] = x[8] - y[8]
	n[9] = x[9] - y[9]
	n[10] = x[10] - y[10]
	n[11] = x[11] - y[11]
	n[12] = x[12] - y[12]
	n[13] = x[13] - y[13]
	n[14] = x[14] - y[14]
	n[15] = x[15] - y[15]

	return n
}

//n = x - y
func (n *bigNumber) sub(x *bigNumber, y *bigNumber) *bigNumber {
	return n.subRaw(x, y).bias(2).weakReduce()
}

func (n *bigNumber) subW(w word) *bigNumber {
	n[0] -= word(w)
	return n
}

func (n *bigNumber) subXBias(x *bigNumber, y *bigNumber, amt word) *bigNumber {
	return n.subRaw(x, y).bias(amt).weakReduce()
}

func (n *bigNumber) mulW(x *bigNumber, w dword) *bigNumber {
	whi := word(w >> radix)
	wlo := word(w & dword(radixMask))

	var accum0, accum8 dword

	accum0 = dword(wlo) * dword(x[0])
	accum8 = dword(wlo) * dword(x[8])
	accum0 += dword(whi) * dword(x[15])
	accum8 += dword(whi) * dword(x[15]+x[7])

	n[0] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[8] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// 1
	accum0 += dword(wlo) * dword(x[1])
	accum8 += dword(wlo) * dword(x[9])
	accum0 += dword(whi) * dword(x[0])
	accum8 += dword(whi) * dword(x[8])

	n[1] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[9] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// 2
	accum0 += dword(wlo) * dword(x[2])
	accum8 += dword(wlo) * dword(x[10])
	accum0 += dword(whi) * dword(x[1])
	accum8 += dword(whi) * dword(x[9])

	n[2] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[10] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// 3
	accum0 += dword(wlo) * dword(x[3])
	accum8 += dword(wlo) * dword(x[11])
	accum0 += dword(whi) * dword(x[2])
	accum8 += dword(whi) * dword(x[10])

	n[3] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[11] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// 4
	accum0 += dword(wlo) * dword(x[4])
	accum8 += dword(wlo) * dword(x[12])
	accum0 += dword(whi) * dword(x[3])
	accum8 += dword(whi) * dword(x[11])

	n[4] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[12] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// 5
	accum0 += dword(wlo) * dword(x[5])
	accum8 += dword(wlo) * dword(x[13])
	accum0 += dword(whi) * dword(x[4])
	accum8 += dword(whi) * dword(x[12])

	n[5] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[13] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// 6
	accum0 += dword(wlo) * dword(x[6])
	accum8 += dword(wlo) * dword(x[14])
	accum0 += dword(whi) * dword(x[5])
	accum8 += dword(whi) * dword(x[13])

	n[6] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[14] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// 7
	accum0 += dword(wlo) * dword(x[7])
	accum8 += dword(wlo) * dword(x[15])
	accum0 += dword(whi) * dword(x[6])
	accum8 += dword(whi) * dword(x[14])

	n[7] = word(accum0 & dword(radixMask))
	accum0 >>= radix

	n[15] = word(accum8 & dword(radixMask))
	accum8 >>= radix

	// finish
	accum0 += accum8 + dword(n[8])
	n[8] = word(accum0 & dword(radixMask))
	n[9] += word(accum0 >> radix)

	accum8 += dword(n[0])
	n[0] = word(accum8 & dword(radixMask))
	n[1] += word(accum8 >> radix)

	return n
}

//n = x * y
func (n *bigNumber) mulCopy(x *bigNumber, y *bigNumber) *bigNumber {
	//it does not work in place, that why the temporary bigNumber is necessary
	return n.set(new(bigNumber).mul(x, y))
}

//n = x * y
func (n *bigNumber) mul(x *bigNumber, y *bigNumber) *bigNumber {
	//it does not work in place, that why the temporary bigNumber is necessary
	return karatsubaMul(n, x, y)
}

//TODO Security this is not constant time
func (n *bigNumber) mulWSignedCurveConstant(x *bigNumber, c sdword) *bigNumber {
	if c >= 0 {
		return n.mulW(x, dword(c))
	}

	r := n.mulW(x, dword(-c))
	return r.sub(bigZero, r)
}

func (n *bigNumber) square(x *bigNumber) *bigNumber {
	return karatsubaSquare(n, x)
}

func (n *bigNumber) squareN(x *bigNumber, y uint) *bigNumber {
	if y&1 != 0 {
		n.square(x)
		y--
	} else {
		n.square(new(bigNumber).square(x))
		y -= 2
	}

	for ; y > 0; y -= 2 {
		n.square(new(bigNumber).square(n))
	}

	return n
}

// n^2 x = 1, QNR; or 0 if x = 0.  Return true if successful
func (n *bigNumber) isr(x *bigNumber) word {
	l0 := new(bigNumber)
	l1 := new(bigNumber)
	l2 := new(bigNumber)

	l1.square(x)
	l2.mul(x, l1)
	l1.square(l2)
	l2.mul(x, l1)
	l1.squareN(l2, 3)
	l0.mul(l2, l1)
	l1.squareN(l0, 3)
	l0.mul(l2, l1)
	l2.squareN(l0, 9)
	l1.mul(l0, l2)
	l0.square(l1)
	l2.mul(x, l0)
	l0.squareN(l2, 18)
	l2.mul(l1, l0)
	l0.squareN(l2, 37)
	l1.mul(l2, l0)
	l0.squareN(l1, 37)
	l1.mul(l2, l0)
	l0.squareN(l1, 111)
	l2.mul(l1, l0)
	l0.square(l2)
	l1.mul(x, l0)
	l0.squareN(l1, 223)
	l1.mul(l2, l0)
	l2.square(l1)
	l0.mul(l2, x)
	n.set(l1)

	return l0.decafEq(bigOne)
}

func invert(x *bigNumber) *bigNumber {
	t1, t2 := &bigNumber{}, &bigNumber{}
	t1.square(x)
	t2.isr(t1)
	t1.square(t2)
	t2.mul(t1, x)
	return t2.copy()
}

func (n *bigNumber) negRaw(x *bigNumber) *bigNumber {
	n[0] = word(-x[0])
	n[1] = word(-x[1])
	n[2] = word(-x[2])
	n[3] = word(-x[3])
	n[4] = word(-x[4])
	n[5] = word(-x[5])
	n[6] = word(-x[6])
	n[7] = word(-x[7])
	n[8] = word(-x[8])
	n[9] = word(-x[9])
	n[10] = word(-x[10])
	n[11] = word(-x[11])
	n[12] = word(-x[12])
	n[13] = word(-x[13])
	n[14] = word(-x[14])
	n[15] = word(-x[15])

	return n
}

func (n *bigNumber) neg(x *bigNumber) *bigNumber {
	return n.negRaw(x).bias(2).weakReduce()
}

func (n *bigNumber) conditionalNegate(neg word) *bigNumber {
	return constantTimeSelect(new(bigNumber).neg(n), n, neg)
}

func (n *bigNumber) decafCondNegate(neg word) {
	n.decafConstTimeSel(n, new(bigNumber).sub(bigZero, n), neg)
}

//if swap == 0xffffffff => n = x, x = n
// This is constant time
func (n *bigNumber) conditionalSwap(x *bigNumber, swap word) *bigNumber {
	for i, xv := range x {
		s := (xv ^ n[i]) & swap
		x[i] ^= s
		n[i] ^= s
	}

	return n
}

func constantTimeSelect(x, y *bigNumber, first word) *bigNumber {
	//TODO this is probably more complicate than it should
	return y.copy().conditionalSwap(x.copy(), first)
}

func (n *bigNumber) decafConstTimeSel(x, y *bigNumber, neg word) {
	n[0] = (x[0] & (^neg)) | (y[0] & (neg))
	n[1] = (x[1] & (^neg)) | (y[1] & (neg))
	n[2] = (x[2] & (^neg)) | (y[2] & (neg))
	n[3] = (x[3] & (^neg)) | (y[3] & (neg))
	n[4] = (x[4] & (^neg)) | (y[4] & (neg))
	n[5] = (x[5] & (^neg)) | (y[5] & (neg))
	n[6] = (x[6] & (^neg)) | (y[6] & (neg))
	n[7] = (x[7] & (^neg)) | (y[7] & (neg))
	n[8] = (x[8] & (^neg)) | (y[8] & (neg))
	n[9] = (x[9] & (^neg)) | (y[9] & (neg))
	n[10] = (x[10] & (^neg)) | (y[10] & (neg))
	n[11] = (x[11] & (^neg)) | (y[11] & (neg))
	n[12] = (x[12] & (^neg)) | (y[12] & (neg))
	n[13] = (x[13] & (^neg)) | (y[13] & (neg))
	n[14] = (x[14] & (^neg)) | (y[14] & (neg))
	n[15] = (x[15] & (^neg)) | (y[15] & (neg))
}

func constantTimeGreaterOrEqualP(n *bigNumber) word {
	ge := word(lmask)

	for i := 0; i < 4; i++ {
		ge &= n[i]
	}

	ge = (ge & (n[4] + 1)) | isZeroMask(word(n[4]^radixMask))

	for i := 5; i < 8; i++ {
		ge &= n[i]
	}

	return ^isZeroMask(word(ge ^ radixMask))
}

func (n *bigNumber) weakReduce() *bigNumber {
	tmp := word(dword(n[nLimbs-1]) >> radix)
	n[nLimbs/2] += tmp

	n[15] = (n[15] & radixMask) + (n[14] >> radix)
	n[14] = (n[14] & radixMask) + (n[13] >> radix)
	n[13] = (n[13] & radixMask) + (n[12] >> radix)
	n[12] = (n[12] & radixMask) + (n[11] >> radix)
	n[11] = (n[11] & radixMask) + (n[10] >> radix)
	n[10] = (n[10] & radixMask) + (n[9] >> radix)
	n[9] = (n[9] & radixMask) + (n[8] >> radix)
	n[8] = (n[8] & radixMask) + (n[7] >> radix)
	n[7] = (n[7] & radixMask) + (n[6] >> radix)
	n[6] = (n[6] & radixMask) + (n[5] >> radix)
	n[5] = (n[5] & radixMask) + (n[4] >> radix)
	n[4] = (n[4] & radixMask) + (n[3] >> radix)
	n[3] = (n[3] & radixMask) + (n[2] >> radix)
	n[2] = (n[2] & radixMask) + (n[1] >> radix)
	n[1] = (n[1] & radixMask) + (n[0] >> radix)

	n[0] = (n[0] & radixMask) + tmp

	return n
}

// Reduce to canonical form
func (n *bigNumber) strongReduce() *bigNumber {
	// clear high
	n.weakReduce()

	// total is less than 2p
	// compute total_value - p.  No need to reduce mod p.

	scarry := sdword(0)
	scarry += sdword(n[0]) - 0xfffffff
	n[0] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[1]) - 0xfffffff
	n[1] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[2]) - 0xfffffff
	n[2] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[3]) - 0xfffffff
	n[3] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[4]) - 0xfffffff
	n[4] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[5]) - 0xfffffff
	n[5] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[6]) - 0xfffffff
	n[6] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[7]) - 0xfffffff
	n[7] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[8]) - 0xffffffe
	n[8] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[9]) - 0xfffffff
	n[9] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[10]) - 0xfffffff
	n[10] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[11]) - 0xfffffff
	n[11] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[12]) - 0xfffffff
	n[12] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[13]) - 0xfffffff
	n[13] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[14]) - 0xfffffff
	n[14] = word(scarry) & radixMask
	scarry >>= 28

	scarry += sdword(n[15]) - 0xfffffff
	n[15] = word(scarry) & radixMask
	scarry >>= 28

	// uncommon case: it was >= p, so now scarry = 0 and this = x
	// common case: it was < p, so now scarry = -1 and this = x - p + 2^255
	// so let's add back in p.  will carry back off the top for 2^255.
	// it can be asserted:
	//assert(isZero(scarry) | isZero(scarry+1));

	scarryMask := word(scarry) & word(radixMask)
	carry := dword(0)
	m := dword(scarryMask)

	carry += dword(n[0]) + m
	n[0] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[1]) + m
	n[1] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[2]) + m
	n[2] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[3]) + m
	n[3] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[4]) + m
	n[4] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[5]) + m
	n[5] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[6]) + m
	n[6] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[7]) + m
	n[7] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[8]) + m&dword(0xfffffffffffffffe)
	n[8] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[9]) + m
	n[9] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[10]) + m
	n[10] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[11]) + m
	n[11] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[12]) + m
	n[12] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[13]) + m
	n[13] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[14]) + m
	n[14] = word(carry) & radixMask
	carry >>= 28

	carry += dword(n[15]) + m
	n[15] = word(carry) & radixMask
	carry >>= 28

	return n
}

func serialize(dst []byte, n *bigNumber) {
	if len(dst) != 56 {
		panic("Failed to serialize")
	}
	src := n.copy()
	src.strongReduce()

	for i := 0; i < 8; i++ {
		l := dword(src[2*i]) + dword(src[2*i+1])<<28
		for j := 0; j < 7; j++ {
			dst[7*i+j] = byte(l)
			l >>= 8
		}
	}
}

func deserializeReturnMask(in serialized) (*bigNumber, word) {
	n := &bigNumber{}

	for i := uint(0); i < 8; i++ {
		out := dword(0)
		for j := uint(0); j < 7; j++ {
			out |= dword(in[7*i+j]) << (8 * j)
		}

		n[2*i] = word(out) & radixMask
		n[2*i+1] = word(out >> 28)
	}

	return n, constantTimeGreaterOrEqualP(n)
}

func deserialize(in serialized) (n *bigNumber, ok bool) {
	n, mask := deserializeReturnMask(in)
	ok = mask == lmask
	return
}

func mustDeserialize(in serialized) *bigNumber {
	n, ok := deserialize(in)
	if !ok {
		panic("Failed to deserialize")
	}

	return n
}

func dsaLikeSerialize(dst []byte, n *bigNumber) {
	n.strongReduce()
	x := n.copy()

	j, fill := uint(0), uint(0)
	buffer := dword(0)

	// TODO: unroll my power!!
	for i := uint(0); i < fieldBytes; i++ {
		if fill < uint(8) && j < nLimbs {
			buffer |= dword(x[j]) << fill
			fill += radix
			j++
		}
		dst[i] = byte(buffer)
		fill -= 8
		buffer >>= 8
	}
}

// TODO: make in type serialized?
func dsaLikeDeserialize(n *bigNumber, in []byte) word {
	j, fill := uint(0), uint(0)
	buffer := dword(0x00)
	scarry := sdword(0x00)

	// TODO: unroll my power!!
	for i := uint(0); i < nLimbs; i++ {
		for fill < radix && j < fieldBytes {
			buffer |= dword(in[j]) << fill
			fill += 8
			j++
		}

		if !(i < nLimbs-1) {

			n[i] = word(buffer)
		}
		n[i] = word(buffer & ((dword(1 << radix)) - 1))

		fill -= radix
		buffer >>= radix
		scarry = sdword((word(scarry) + n[i] - modulus[i]) >> 8 * 4)
	}

	// TODO: check me, and add case when hibit is one
	var high word = 0x01
	succ := -(high)
	succ &= isZeroMask(word(buffer))
	succ &= ^(isZeroMask(word(scarry)))

	return succ
}

func (n *bigNumber) String() string {
	dst := make([]byte, fieldBytes)
	serialize(dst[:], n)
	return fmt.Sprintf("%#v", dst)
	//return fmt.Sprintf("0x%s", new(big.Int).SetBytes(rev(dst)).Text(16))
}

func (n *bigNumber) limbs() []word {
	return n[:]
}
