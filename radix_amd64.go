package ed448

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

//XXX unroll
func (n *bigNumber) bias(b uint) *bigNumber {
	var co1 limb = radixMask * limb(b)
	var co2 limb = co1 - limb(b)

	for i := 0; i < len(n); i++ {
		if i == 4 {
			n[i] += co2
		} else {
			n[i] += co1
		}
	}

	return n
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

		//overflows
		if c1 == 1 && c2 == 1 {
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
	acc0 := []Word{0, 0}
	acc4 := []Word{0, 0}

	tmp := []Word{0, 0}
	for i := 0; i < 4; i++ {
		mulAddVWW_g(tmp, []Word{Word(x[i]), 0}, Word(w), 0)
		addVV_g(acc0, acc0, tmp) //XXX should we check carry?

		mulAddVWW_g(tmp, []Word{Word(x[i+4]), 0}, Word(w), 0)
		addVV_g(acc4, acc4, tmp) //XXX should we check carry?

		n[i] = limb(acc0[0]) & radixMask
		shrVU_g(acc0, acc0, Radix)

		n[i+4] = limb(acc4[0]) & radixMask
		shrVU_g(acc4, acc4, Radix)
	}

	addVV_g(acc0, acc0, acc4) //XXX should we check carry?
	addVV_g(acc0, acc0, []Word{Word(n[4]), 0, 0})

	n[4] = limb(acc0[0]) & radixMask
	shrVU_g(tmp, acc0, Radix)
	n[5] += limb(tmp[0])

	addVV_g(acc4, acc4, []Word{Word(n[0]), 0, 0})
	n[0] = limb(acc4[0]) & radixMask
	shrVU_g(tmp, acc4, Radix)
	n[1] += limb(tmp[0])

	return n
}
