package ed448

import "fmt"

type limb Word
type bigNumber [Limbs]limb //XXX Should this type be a pointer to an array?
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

//XXX It may not work on 64-bit
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

func (n *bigNumber) String() string {
	dst := [56]byte{}
	serialize(dst[:], n)
	return fmt.Sprintf("%#v", dst)
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

//XXX Is there any optimum way of squaring?
func (n *bigNumber) square() *bigNumber {
	return n.mul(n)
}

func (n *bigNumber) mulWSignedCurveConstant(x *bigNumber, c int64) *bigNumber {
	if c >= 0 {
		return n.mulW(x, uint64(c))
	}

	r := n.mulW(x, uint64(-c))
	r.negRaw(r)
	r.bias(2)

	return r
}

func (n *bigNumber) negRaw(x *bigNumber) *bigNumber {
	for i, xi := range x {
		n[i] = limb(-xi)
	}

	return n
}
