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

//n = x + y
func (n *bigNumber) add(x *bigNumber, y *bigNumber) *bigNumber {
	return n.addRaw(x, y).weakReduce()
}

func (n *bigNumber) addRaw(x *bigNumber, y *bigNumber) *bigNumber {
	for i := 0; i < len(n); i++ {
		n[i] = x[i] + y[i]
	}

	return n
}

//n = x - y
func (n *bigNumber) sub(x *bigNumber, y *bigNumber) *bigNumber {
	return n.subRaw(x, y).bias(2).weakReduce()
}

func (n *bigNumber) subRaw(x *bigNumber, y *bigNumber) *bigNumber {
	for i := 0; i < len(n); i++ {
		n[i] = x[i] - y[i]
	}

	return n
}

//n = x * y
func (n *bigNumber) mul(x *bigNumber, y *bigNumber) *bigNumber {
	//it does not work in place, that why the temporary bigNumber is necessary
	for i, ni := range karatsubaMul(new(bigNumber), x, y) {
		n[i] = ni
	}

	return n
}

//XXX Is there any optimum way of squaring?
func (n *bigNumber) square(x *bigNumber) *bigNumber {
	return n.mul(x, x)
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
	n.negRaw(x)
	n.bias(2)
	n.weakReduce()
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

func (n *bigNumber) zero() (eq bool) {
	x := n.copy().strongReduce()
	r := limb(0)

	for _, ni := range x {
		r |= ni ^ 0
	}

	return r == 0
}

//in is big endian
func (n *bigNumber) setBytes(in []byte) *bigNumber {
	if len(in) != 56 {
		panic("1")
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
