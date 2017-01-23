package ed448

import "fmt"

type bigNumber [limbs]word_t
type serialized [56]byte

func mustDeserialize(in serialized) *bigNumber {
	n, ok := deserialize(in)
	if !ok {
		panic("Failed to deserialize")
	}

	return n
}

func isZeroMask(n uint32) uint32 {
	nn := uint64(n)
	nn = nn - 1
	return uint32(nn >> wordBits)
}

func constantTimeGreaterOrEqualP(n *bigNumber) word_t {
	ge := word_t(lmask)

	for i := 0; i < 4; i++ {
		ge &= n[i]
	}

	ge = (ge & (n[4] + 1)) | word_t(isZeroMask(uint32(n[4]^radixMask)))

	for i := 5; i < 8; i++ {
		ge &= n[i]
	}

	return word_t(^isZeroMask(uint32(ge ^ radixMask)))
}

//n = x + y
func (n *bigNumber) add(x *bigNumber, y *bigNumber) *bigNumber {
	return n.addRaw(x, y).weakReduce()
}

func (n *bigNumber) addW(w uint32) *bigNumber {
	n[0] += word_t(w)
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

func (n *bigNumber) setUI(y uint64) *bigNumber {
	n[0] = word_t(y) & radixMask
	n[1] = word_t(y >> radix)
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

//n = x - y
func (n *bigNumber) sub(x *bigNumber, y *bigNumber) *bigNumber {
	return n.subRaw(x, y).bias(2).weakReduce()
}

func (n *bigNumber) subW(w uint32) *bigNumber {
	n[0] -= word_t(w)
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

func (n *bigNumber) subXBias(x *bigNumber, y *bigNumber, amt uint32) *bigNumber {
	return n.subRaw(x, y).bias(amt).weakReduce()
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

func (n *bigNumber) isr(x *bigNumber) *bigNumber {
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

	return n.mul(l2, l0)
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

func (n *bigNumber) weakReduce() *bigNumber {
	tmp := word_t(uint64(n[limbs-1]) >> radix)
	n[limbs/2] += tmp

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

//XXX Security this should be constant time
func (n *bigNumber) mulWSignedCurveConstant(x *bigNumber, c int64) *bigNumber {
	if c >= 0 {
		return n.mulW(x, uint64(c))
	}

	r := n.mulW(x, uint64(-c))
	r.negRaw(r)
	r.bias(2)

	return r
}

func (n *bigNumber) neg(x *bigNumber) *bigNumber {
	return n.negRaw(x).bias(2).weakReduce()
}

func (n *bigNumber) conditionalNegate(neg word_t) *bigNumber {
	return constantTimeSelect(new(bigNumber).neg(n), n, neg)
}

func constantTimeSelect(x, y *bigNumber, first word_t) *bigNumber {
	//XXX this is probably more complicate than it should
	return y.copy().conditionalSwap(x.copy(), first)
}

//if swap == 0xffffffff => n = x, x = n
func (n *bigNumber) conditionalSwap(x *bigNumber, swap word_t) *bigNumber {
	for i, xv := range x {
		s := (xv ^ n[i]) & swap
		x[i] ^= s
		n[i] ^= s
	}

	return n
}

func (n *bigNumber) decafCondNegate(neg dword_t) {
	y := &bigNumber{}
	y.sub(&bigNumber{0}, n)
	n.decafConstTimeSel(n, y, neg)
}

func (n *bigNumber) decafConstTimeSel(x, y *bigNumber, neg dword_t) {
	n[0] = (x[0] & word_t(^neg)) | (y[0] & word_t(neg))
	n[1] = (x[1] & word_t(^neg)) | (y[1] & word_t(neg))
	n[2] = (x[2] & word_t(^neg)) | (y[2] & word_t(neg))
	n[3] = (x[3] & word_t(^neg)) | (y[3] & word_t(neg))
	n[4] = (x[4] & word_t(^neg)) | (y[4] & word_t(neg))
	n[5] = (x[5] & word_t(^neg)) | (y[5] & word_t(neg))
	n[6] = (x[6] & word_t(^neg)) | (y[6] & word_t(neg))
	n[7] = (x[7] & word_t(^neg)) | (y[7] & word_t(neg))
	n[8] = (x[8] & word_t(^neg)) | (y[8] & word_t(neg))
	n[9] = (x[9] & word_t(^neg)) | (y[9] & word_t(neg))
	n[10] = (x[10] & word_t(^neg)) | (y[10] & word_t(neg))
	n[11] = (x[11] & word_t(^neg)) | (y[11] & word_t(neg))
	n[12] = (x[12] & word_t(^neg)) | (y[12] & word_t(neg))
	n[13] = (x[13] & word_t(^neg)) | (y[13] & word_t(neg))
	n[14] = (x[14] & word_t(^neg)) | (y[14] & word_t(neg))
	n[15] = (x[15] & word_t(^neg)) | (y[15] & word_t(neg))
}

func (n *bigNumber) negRaw(x *bigNumber) *bigNumber {
	n[0] = word_t(-x[0])
	n[1] = word_t(-x[1])
	n[2] = word_t(-x[2])
	n[3] = word_t(-x[3])
	n[4] = word_t(-x[4])
	n[5] = word_t(-x[5])
	n[6] = word_t(-x[6])
	n[7] = word_t(-x[7])
	n[8] = word_t(-x[8])
	n[9] = word_t(-x[9])
	n[10] = word_t(-x[10])
	n[11] = word_t(-x[11])
	n[12] = word_t(-x[12])
	n[13] = word_t(-x[13])
	n[14] = word_t(-x[14])
	n[15] = word_t(-x[15])

	return n
}

func (n *bigNumber) copy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

func (n *bigNumber) set(x *bigNumber) *bigNumber {
	copy(n[:], x[:])
	return n
}

func (n *bigNumber) equals(o *bigNumber) (eq bool) {
	r := word_t(0)
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

func decafEq(x, y *bigNumber) dword_t {
	n := &bigNumber{}
	n.sub(x, y)
	n.strongReduce()

	var ret word_t

	for i := 0; i < limbs; i++ {
		ret |= n[i]
	}
	return ((dword_t(ret) - 1) >> 32)
}

func (n *bigNumber) zeroMask() uint32 {
	x := n.copy().strongReduce()
	r := word_t(0)

	r |= x[0] ^ 0
	r |= x[1] ^ 0
	r |= x[2] ^ 0
	r |= x[3] ^ 0
	r |= x[4] ^ 0
	r |= x[5] ^ 0
	r |= x[6] ^ 0
	r |= x[7] ^ 0
	r |= x[8] ^ 0
	r |= x[9] ^ 0
	r |= x[10] ^ 0
	r |= x[11] ^ 0
	r |= x[12] ^ 0
	r |= x[13] ^ 0
	r |= x[14] ^ 0
	r |= x[15] ^ 0

	return isZeroMask(uint32(r))
}

func (n *bigNumber) zero() (eq bool) {
	return n.zeroMask() == lmask
}

//in is big endian
func (n *bigNumber) setBytes(in []byte) *bigNumber {
	if len(in) != 56 {
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

func (n *bigNumber) String() string {
	dst := make([]byte, 56)
	serialize(dst[:], n)
	return fmt.Sprintf("%#v", dst)
	//return fmt.Sprintf("0x%s", new(big.Int).SetBytes(rev(dst)).Text(16))
}

func (n *bigNumber) limbs() []word_t {
	return n[:]
}
