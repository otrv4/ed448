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
	return barrettDeserializeReturnMask(dst, serial, p) != 0
}

func barrettDeserializeReturnMask(dst []word_t, serial []byte, p *barrettPrime) word_t {
	s := p.wordsInP * wordBits / 8
	if p.pShift != 0 {
		s -= (wordBits - p.pShift) / 8
	}

	bytesToWords(dst, serial[:s])

	carry := dword_t(0)
	for i, wi := range dst {
		carry >>= wordBits
		carry += dword_t(wi)
		if i < len(p.lowWords) {
			carry += dword_t(p.lowWords[i])
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

	return word_t(^scarry)
}

func barrettDeserializeAndReduce(dst []word_t, serial []byte, p *barrettPrime) {
	wordLen := wordBits / 8
	size := (len(serial) + wordLen - 1) / wordLen
	if size < int(p.wordsInP) {
		size = int(p.wordsInP)
	}

	tmp := make([]word_t, size)
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

	cout := addExtPacked(dst, dst[:p.wordsInP], p.lowWords, lmask)

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
		product := dword_t(mand) * dword_t(mier[i])
		product += dword_t(accum[i])
		product += dword_t(carry)

		accum[i] = word_t(product)
		carry = word_t(product >> wordBits)
	}

	for i := len(mier); i < len(accum); i++ {
		sum := dword_t(carry) + dword_t(accum[i])
		accum[i] = word_t(sum)
		carry = word_t(sum >> wordBits)
	}

	return carry
}

func barrettNegate(dst []word_t, p *barrettPrime) {
	barrettReduce(dst, 0, p)

	carry := int64(0)
	for i := 0; i < len(p.lowWords); i++ {
		carry = carry - int64(p.lowWords[i]) - int64(dst[i])
		dst[i] = word_t(carry)
		carry >>= wordBits
	}

	for i := len(p.lowWords); i < int(p.wordsInP); i++ {
		carry = carry - int64(dst[i])
		dst[i] = word_t(carry)
		if i < int(p.wordsInP-1) {
			carry >>= wordBits
		}
	}

	carry = carry + int64(word_t(1)<<p.pShift)
	dst[p.wordsInP-1] = word_t(carry)
}

func barrettMac(dst, x, y []word_t, p *barrettPrime) {
	nWords := int(p.wordsInP)
	if nWords < len(x) {
		nWords = len(x)
	}
	nWords++

	if nWords < len(dst) {
		nWords = len(dst)
	}

	tmp := make([]word_t, nWords)

	for bpos := len(y) - 1; bpos >= 0; bpos-- {
		for idown := nWords - 2; idown >= 0; idown-- {
			tmp[idown+1] = tmp[idown]
		}

		tmp[0] = 0

		carry := widemac(tmp, x, y[bpos], 0)
		barrettReduce(tmp, carry, p)
	}

	cout := addPacked(tmp, dst)
	barrettReduce(tmp, cout, p)

	for i := 0; i < nWords && i < len(dst); i++ {
		dst[i] = tmp[i]
	}

	for i := nWords; i < len(dst); i++ {
		dst[i] = 0
	}
}

func addPacked(dst, x []word_t) word_t {
	carry := dword_t(0)

	//dst can be longer than x
	for i := 0; i < len(x); i++ {
		carry = carry + dword_t(dst[i]) + dword_t(x[i])
		dst[i] = word_t(carry)
		carry >>= wordBits
	}

	return word_t(carry)
}
