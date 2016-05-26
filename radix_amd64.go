package ed448

import "fmt"

const (
	Limbs     = 8
	Radix     = 56
	radixMask = limb(0xffffffffffffff)
)

func deserialize(in serialized) (n *bigNumber, ok bool) {
	n = &bigNumber{}

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

func serialize(dst []byte, src *bigNumber) {
	const (
		rows    = Limbs
		columns = Radix / Limbs
	)

	n := bigNumber{}
	copy(n[:], src[:])

	for i := uint(0); i < rows; i++ {
		for j := uint(0); j < columns; j++ {
			dst[columns*i+j] = byte(n[i])
			n[i] >>= 8
		}
	}
}

func (n *bigNumber) bias(b uint) {
	//noop
}

func (n *bigNumber) strongReduce() *bigNumber {
	mask := radixMask - 1

	//TODO
	n[4] += n[7] >> 56
	n[0] += n[7] >> 56
	n[7] &= radixMask

	acc := []Word{0, 0, 0}
	scarry := []Word{0, 0, 0}
	for i := 0; i < 8; i++ {
		m := radixMask
		if i == 4 {
			m = mask
		}

		c1 := subVV_g(acc, []Word{Word(n[i]), 0, 0}, []Word{Word(m), 0, 0})
		c2 := addVV_g(scarry, scarry, acc)

		fmt.Println(c1, c2)

		//overflows
		if c1 == 1 && c2 == 1 {
			fmt.Printf("%#v\n", scarry)
			scarry[1] = 0xffffffffffffffff
		}

		n[i] = limb(scarry[0]) & radixMask

		shrVU_g(scarry, scarry, 56)
	}

	scarryMask := scarry[0] & Word(radixMask)

	carry := []Word{0, 0, 0}
	for i := 0; i < 8; i++ {
		m := []Word{scarryMask, 0, 0}
		if i == 4 {
			m[0] &= 0xfffffffffffffffe
		}

		addVV_g(acc, []Word{Word(n[i]), 0, 0}, m)
		addVV_g(carry, carry, acc)

		n[i] = limb(carry[0]) & radixMask

		shrVU_g(carry, carry, 56)
	}

	return n
}

func (n *bigNumber) mulW(x *bigNumber, w uint64) *bigNumber {
	//TODO
	return n
}
