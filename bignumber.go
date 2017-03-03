package ed448

import "fmt"

type bigNumber [limbs]word
type serialized [fieldBytes]byte

func (n *bigNumber) zero() (eq bool) {
	return n.zeroMask() == lmask
}

//n = x + y
func (n *bigNumber) add(x *bigNumber, y *bigNumber) *bigNumber {
	return n.addRaw(x, y).weakReduce()
}

func (n *bigNumber) addW(w word) *bigNumber {
	n[0] += word(w)
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

func mustDeserialize(in serialized) *bigNumber {
	n, ok := deserialize(in)
	if !ok {
		panic("Failed to deserialize")
	}

	return n
}

func (n *bigNumber) invert(x *bigNumber) {
	t1 := &bigNumber{}
	t1.square(x)
	n.isr(t1)
	t1.square(n)
	n.mul(t1, x)
}

func (n *bigNumber) neg(x *bigNumber) *bigNumber {
	return n.negRaw(x).bias(2).weakReduce()
}

func (n *bigNumber) conditionalNegate(neg word) *bigNumber {
	return constantTimeSelect(new(bigNumber).neg(n), n, neg)
}

func constantTimeSelect(x, y *bigNumber, first word) *bigNumber {
	//XXX this is probably more complicate than it should
	return y.copy().conditionalSwap(x.copy(), first)
}

//if swap == 0xffffffff => n = x, x = n
func (n *bigNumber) conditionalSwap(x *bigNumber, swap word) *bigNumber {
	for i, xv := range x {
		s := (xv ^ n[i]) & swap
		x[i] ^= s
		n[i] ^= s
	}

	return n
}

func (n *bigNumber) decafCondNegate(neg word) {
	n.decafConstTimeSel(n, new(bigNumber).sub(bigZero, n), neg)
}

func (n *bigNumber) copy() *bigNumber {
	c := &bigNumber{}
	*c = *n
	return c
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

func (n *bigNumber) String() string {
	dst := make([]byte, fieldBytes)
	serialize(dst[:], n)
	return fmt.Sprintf("%#v", dst)
	//return fmt.Sprintf("0x%s", new(big.Int).SetBytes(rev(dst)).Text(16))
}

func (n *bigNumber) limbs() []word {
	return n[:]
}
