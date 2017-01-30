package ed448

type smvtControl struct {
	power, addend int
}

func recodeWnaf(control []smvtControl, scalar scalar32, nBits, tableBits uint) (position uint32) {
	current := 0
	var i, j int
	position = 0
	for i = int(nBits - 1); i >= 0; i-- {
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

func prepareWnafTable(dst []*twPNiels, p *twExtensible, tableSize uint) {
	dst[0] = p.twPNiels()

	if tableSize == 0 {
		return
	}

	p.double()
	twOp := p.twPNiels()

	p.addTwPNiels(dst[0])
	dst[1] = p.twPNiels()

	for i := 2; i < 1<<tableSize; i++ {
		p.addTwPNiels(twOp)
		dst[i] = p.twPNiels()
	}
}

func decafPrepareWnafTable(dst []*twPNiels, p *twExtendedPoint, tableSize uint) {
	dst[0] = p.twPNiels()

	if tableSize == 0 {
		return
	}

	p.double(false)

	twOp := p.twPNiels()

	p.addProjectiveNielsToExtended(dst[0], false)
	dst[1] = p.twPNiels()

	for i := 2; i < 1<<tableSize; i++ {
		p.addProjectiveNielsToExtended(twOp, false)
		dst[i] = p.twPNiels()
	}
}

func linearComboVarFixedVt(
	working *twExtensible, scalarVar, scalarPre scalar32, precmp []*twNiels) {
	tableBitsVar := uint(4) //SCALARMUL_WNAF_COMBO_TABLE_BITS;
	nbitsVar := uint(446)
	nbitsPre := uint(446)
	tableBitsPre := uint(5)

	var controlVar [92]smvtControl // nbitsVar/(tableBitsVar+1)+3
	var controlPre [77]smvtControl // nbitsPre/(tableBitsPre+1)+3

	recodeWnaf(controlVar[:], scalarVar, nbitsVar, tableBitsVar)
	recodeWnaf(controlPre[:], scalarPre, nbitsPre, tableBitsPre)

	var precmpVar [16]*twPNiels // 1 << tableBitsVar
	prepareWnafTable(precmpVar[:], working, uint(tableBitsVar))

	contp := 0
	contv := 0

	i := controlVar[0].power
	if i > controlPre[0].power {
		convertTwPnielsToTwExtensible(working, precmpVar[controlVar[0].addend>>1])
		contv++
	} else if i == controlPre[0].power && i >= 0 {
		convertTwPnielsToTwExtensible(working, precmpVar[controlVar[0].addend>>1])
		working.addTwNiels(precmp[controlPre[0].addend>>1])
		contv++
		contp++
	} else {
		i = controlPre[0].power
		convertTwNielsToTwExtensible(working, precmp[controlPre[0].addend>>1])
		contp++
	}

	if i < 0 {
		working.setIdentity()
		return
	}

	for i--; i >= 0; i-- {
		working.double()

		if i == controlVar[contv].power {
			if controlVar[contv].addend > 0 {
				working.addTwPNiels(precmpVar[controlVar[contv].addend>>1])
			} else {
				working.subTwPNiels(precmpVar[(-controlVar[contv].addend)>>1])
			}
			contv++
		}

		if i == controlPre[contp].power {
			if controlPre[contp].addend > 0 {
				working.addTwNiels(precmp[controlPre[contp].addend>>1])
			} else {
				working.subTwNiels(precmp[(-controlPre[contp].addend)>>1])
			}
			contp++
		}
	}
}
