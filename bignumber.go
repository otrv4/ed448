package ed448

import "fmt"

type limb word_t
type bigNumber [Limbs]limb
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
	ge := limb(0xffffffff)

	for i := 0; i < 4; i++ {
		ge &= n[i]
	}

	ge = (ge & (n[4] + 1)) | limb(isZeroMask(uint32(n[4]^radixMask)))

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
	n[0] += limb(w)
	return n
}

func (n *bigNumber) addRaw(x *bigNumber, y *bigNumber) *bigNumber {
	for i := 0; i < len(n); i++ {
		n[i] = x[i] + y[i]
	}

	return n
}

func (n *bigNumber) setUi(y uint64) *bigNumber {
	n[0] = limb(y) & radixMask
	n[1] = limb(y >> Radix)

	for i := 2; i < Limbs; i++ {
		n[i] = 0
	}

	return n
}

//n = x - y
func (n *bigNumber) sub(x *bigNumber, y *bigNumber) *bigNumber {
	return n.subRaw(x, y).bias(2).weakReduce()
}

func (n *bigNumber) subW(w uint32) *bigNumber {
	n[0] -= limb(w)
	return n
}

func (n *bigNumber) subRaw(x *bigNumber, y *bigNumber) *bigNumber {
	for i := 0; i < len(n); i++ {
		n[i] = x[i] - y[i]
	}

	return n
}

func (n *bigNumber) subxRaw(x *bigNumber, y *bigNumber) *bigNumber {
	// XXX Only weakReduce when 32bits
	return n.subRaw(x, y).bias(2).weakReduce()
}

//n = x * y
func (n *bigNumber) mul(x *bigNumber, y *bigNumber) *bigNumber {
	//it does not work in place, that why the temporary bigNumber is necessary
	for i, ni := range karatsubaMul(new(bigNumber), x, y) {
		n[i] = ni
	}

	return n
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

//XXX Is there any optimum way of squaring?
func (n *bigNumber) square(x *bigNumber) *bigNumber {
	return n.mul(x, x)
}

func (n *bigNumber) squareN(x *bigNumber, y uint) *bigNumber {
	if y&1 != 0 {
		n.square(x)
		y -= 1
	} else {
		n.square(x).square(n)
		y -= 2
	}

	for ; y > 0; y -= 2 {
		n.square(n).square(n)
	}

	return n
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
		s := (xv ^ n[i]) & limb(swap)
		x[i] ^= s
		n[i] ^= s
	}

	return n
}

func (n *bigNumber) negRaw(x *bigNumber) *bigNumber {
	for i, xi := range x {
		n[i] = limb(-xi)
	}

	return n
}

func (n *bigNumber) copy() *bigNumber {
	c := &bigNumber{}
	copy(c[:], n[:])
	return c
}

func (n *bigNumber) equals(o *bigNumber) (eq bool) {
	r := limb(0)
	x := n.copy().strongReduce()
	y := o.copy().strongReduce()

	for i, yi := range y {
		r |= x[i] ^ yi
	}

	return r == 0
}

func (n *bigNumber) zeroMask() uint32 {
	x := n.copy().strongReduce()
	r := limb(0)

	for _, ni := range x {
		r |= ni ^ 0
	}

	return isZeroMask(uint32(r))
}

func (n *bigNumber) zero() (eq bool) {
	return n.zeroMask() == 0xffffffff
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

func (n *bigNumber) limbs() []limb {
	return n[:]
}
