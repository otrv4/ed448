package ed448

type barrettPrime struct {
	wordsInP uint32
	pShift   uint32
	lowWords []word_t
}

var curvePrimeOrder = barrettPrime{
	wordsInP: 14,
	pShift:   30,
	lowWords: []word_t{
		0x54a7bb0d,
		0xdc873d6d,
		0x723a70aa,
		0xde933d8d,
		0x5129c96f,
		0x3bb124b6,
		0x8335dc16,
	},
}

func barrettDeserialize(dst []word_t, serial []byte, p *barrettPrime) bool {
	s := p.wordsInP * wordBits / 8
	if p.pShift != 0 {
		s -= (wordBits - p.pShift) / 8
	}

	bytesToWords(dst, serial[:s])

	carry := uint64(0)
	for i, wi := range dst {
		carry >>= wordBits
		carry += uint64(wi)
		if i < len(p.lowWords) {
			carry += uint64(p.lowWords[i])
		}
	}

	if p.pShift != 0 {
		carry >>= p.pShift
	} else {
		carry >>= wordBits
	}

	scarry := int64(carry)
	scarry = -scarry
	scarry >>= wordBits
	scarry >>= wordBits

	return ^scarry != 0
}

func barrettDeserializeAndReduce(dst []word_t, serial [64]byte, p *barrettPrime) {
	tmp := [16]word_t{} //XXX Why is this 16 if dst has len = 14?

	bytesToWords(tmp[:], serial[:])
	barrettReduce(tmp[:], 0, p)

	for i := uint32(0); i < p.wordsInP; i++ {
		dst[i] = tmp[i]
	}
}

func barrettReduce(dst []word_t, carry word_t, p *barrettPrime) {
	for wordsLeft := uint32(len(dst)); wordsLeft >= p.wordsInP; wordsLeft-- {
		//XXX PERF unroll
		for repeat := 0; repeat < 2; repeat++ {
			mand := dst[wordsLeft-1] >> p.pShift
			dst[wordsLeft-1] &= (word_t(1) << p.pShift) - 1

			if p.pShift != 0 && repeat == 0 {
				if wordsLeft < uint32(len(dst)) {
					mand |= dst[wordsLeft] << (wordBits - p.pShift)
					dst[wordsLeft] = 0
				} else {
					mand |= carry << (wordBits - p.pShift)
				}
			}

			carry = widemac(
				dst[wordsLeft-p.wordsInP:wordsLeft],
				p.lowWords, mand, 0)
		}
	}

	cout := addExtPacked(dst, dst[:p.wordsInP], p.lowWords, 0xffffffff)

	if p.pShift != 0 {
		cout = (cout << (wordBits - p.pShift)) + (dst[p.wordsInP-1] >> p.pShift)
		dst[p.wordsInP-1] &= word_t(1)<<p.pShift - 1
	}

	/* mask = carry-1: if no carry then do sub, otherwise don't */
	subExtPacked(dst, dst[:p.wordsInP], p.lowWords, cout-1)
}

func addExtPacked(dst, x, y []word_t, mask word_t) word_t {
	carry := int64(0)
	for i := 0; i < len(y); i++ {
		carry += int64(x[i]) + int64(y[i]&mask)
		dst[i] = word_t(carry)
		carry >>= wordBits
	}

	for i := len(y); i < len(x); i++ {
		carry += int64(x[i])
		dst[i] = word_t(carry)
		carry >>= wordBits
	}

	return word_t(carry)
}

func subExtPacked(dst, x, y []word_t, mask word_t) word_t {
	carry := int64(0)
	for i := 0; i < len(y); i++ {
		carry += int64(x[i]) - (int64(y[i]) & int64(mask))
		dst[i] = word_t(carry)
		carry >>= wordBits
	}

	for i := len(y); i < len(x); i++ {
		carry += int64(x[i])
		dst[i] = word_t(carry)
		carry >>= wordBits
	}

	return word_t(carry)
}

//XXX Is this the same as mulAddVWW_g() ?
func widemac(accum []word_t, mier []word_t, mand, carry word_t) word_t {
	for i := 0; i < len(mier); i++ {
		product := uint64(mand) * uint64(mier[i])
		product += uint64(accum[i])
		product += uint64(carry)

		accum[i] = word_t(product)
		carry = word_t(product >> wordBits)
	}

	for i := len(mier); i < len(accum); i++ {
		sum := uint64(carry) + uint64(accum[i])
		accum[i] = word_t(sum)
		carry = word_t(sum >> wordBits)
	}

	return carry
}
