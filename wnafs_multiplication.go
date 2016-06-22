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
		current = 2*current + int(bit)

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
		control[position].power = j
		control[position].addend = current
		position++
	}

	control[position].power = -1
	control[position].addend = 0

	return
}

func linear_combo_var_fixed_vt(p *twExtensible, x, y []word_t, table []*twNiels) {
	tableBitsVar := 4
	tableBitsPre := 5
	controlVar := make([]smvt_control, scalarBits/(tableBitsVar+1)+3)
	controlPre := make([]smvt_control, scalarBits/(tableBitsPre+1)+3)

	recodeWnaf(controlVar, x, scalarBits, tableBitsVar)
	recodeWnaf(controlPre, y, scalarBits, tableBitsVar)

	precmpVar := make([]*twPNiels, 1<<uint(tableBitsVar))
	prepareWnafTable(precmpVar, p, uint(tableBitsVar))

	contP := 0
	contV := 0
	working := p.copy()

	i := controlVar[0].power
	if i > controlPre[0].power {
		working = precmpVar[controlVar[0].addend>>1].TwistedExtensible()
		contV++
	} else if i == controlPre[0].power && i >= 0 {
		working = precmpVar[controlVar[0].addend>>1].TwistedExtensible()
		working = working.addTwNiels(table[controlPre[0].addend>>1])
		contV++
		contP++
	} else {
		i = controlPre[0].power
		working = table[controlPre[0].addend>>1].TwistedExtensible()
		contP++
	}

	if i < 0 {
		working.setIdentity()
		return
	}

	for i--; i >= 0; i-- {
		working = working.double()

		if i == controlVar[contV].power {
			if controlVar[contV].addend > 0 {
				working = working.addTwPNiels(precmpVar[controlVar[contV].addend>>1])
			} else {
				working = working.subTwPNiels(precmpVar[(-controlVar[contV].addend)>>1])
			}
			contV++
		}

		if i == controlPre[contP].power {
			if controlPre[contP].addend > 0 {
				working = working.addTwNiels(table[controlPre[contP].addend>>1])
			} else {
				working = working.subTwNiels(table[(-controlPre[contP].addend)>>1])
			}
			contP++
		}
	}

	//XXX PERF: should it be in-place?
	p.x = working.x.copy()
	p.y = working.y.copy()
	p.z = working.z.copy()
	p.t = working.t.copy()
	p.u = working.u.copy()

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
