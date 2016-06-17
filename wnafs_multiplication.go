package ed448

type smvt_control struct {
	power, addend int
}

func recodeWnaf(control []smvt_control, scalar []word_t, nBits, tableBits int) (position uint32) {

	current := 0
	var i, j int
	position = 0
	for i = nBits - 1; i >= 0; i-- {
		bit := (scalar[i/wordBits] >> uint(i%wordBits)) & 1
		current = (2 * current) + int(bit)

		/*
		 * Sizing: |current| >= 2^(tableBits+1) -> |current| = 2^0
		 * So current loses (tableBits+1) bits every time.  It otherwise gains
		 * 1 bit per iteration.  The number of iterations is
		 * (nbits + 2 + tableBits), and an additional control word is added at
		 * the end.  So the total number of control words is at most
		 * ceil((nbits+1) / (tableBits+1)) + 2 = floor((nbits)/(tableBits+1)) + 2.
		 * There's also the stopper with power -1, for a total of +3.
		 */
		if current >= (2<<uint32(tableBits)) || current <= -1-(2<<uint32(tableBits)) {
			delta := (current + 1) >> 1 /* |delta| < 2^tablebits */
			current = -(current & 1)

			for j = i; (delta & 1) == 0; j++ {
				delta >>= 1
			}
			control[position].power = j + 1
			control[position].addend = delta
			position++
		}
	}

	if current != 0 {
		for j = 0; (current & 1) == 0; j++ {
			current >>= 1
		}
		control[position].power = int(j)
		control[position].addend = current
		position++
	}

	control[position].power = -1
	control[position].addend = 0

	return
}

func linear_combo_var_fixed_vt(p *twExtensible, x, y []word_t, table []*twNiels) {
	return
}

func prepareWnafTable(dst []*twPNiels, p *twExtensible, tableSize uint) {
	dst[0] = p.twPNiels()

	if tableSize == 0 {
		return
	}

	p = p.double()
	twOp := p.twPNiels()

	p = p.addTwPNiels(dst[0])
	dst[1] = p.twPNiels()

	for i := 2; i < 1<<tableSize; i++ {
		p = p.addTwPNiels(twOp)
		dst[i] = p.twPNiels()
	}
}
