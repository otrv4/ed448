package ed448

type smvt_control struct {
	power, addend int
}

func recodeWnaf(control []smvt_control, scalar []word_t, nBits, tableBits uint) (position uint32) {
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

func linear_combo_var_fixed_vt(
	working *twExtensible, scalar_var, scalar_pre []word_t, precmp []*twNiels) {
	table_bits_var := uint(4) //SCALARMUL_WNAF_COMBO_TABLE_BITS;
	nbits_var := uint(446)
	nbits_pre := uint(446)
	table_bits_pre := uint(5)

	var control_var [92]smvt_control // nbits_var/(table_bits_var+1)+3
	var control_pre [77]smvt_control // nbits_pre/(table_bits_pre+1)+3

	recodeWnaf(control_var[:], scalar_var, nbits_var, table_bits_var)
	recodeWnaf(control_pre[:], scalar_pre, nbits_pre, table_bits_pre)

	var precmp_var [16]*twPNiels // 1 << table_bits_var
	prepareWnafTable(precmp_var[:], working, uint(table_bits_var))

	contp := 0
	contv := 0

	i := control_var[0].power
	if i > control_pre[0].power {
		convertTwPnielsToTwExtensible(working, precmp_var[control_var[0].addend>>1])
		contv++
	} else if i == control_pre[0].power && i >= 0 {
		convertTwPnielsToTwExtensible(working, precmp_var[control_var[0].addend>>1])
		working.addTwNiels(precmp[control_pre[0].addend>>1])
		contv++
		contp++
	} else {
		i = control_pre[0].power
		convertTwNielsToTwExtensible(working, precmp[control_pre[0].addend>>1])
		contp++
	}

	if i < 0 {
		working.setIdentity()
		return
	}

	for i--; i >= 0; i-- {
		working.double()

		if i == control_var[contv].power {
			if control_var[contv].addend > 0 {
				working.addTwPNiels(precmp_var[control_var[contv].addend>>1])
			} else {
				working.subTwPNiels(precmp_var[(-control_var[contv].addend)>>1])
			}
			contv++
		}

		if i == control_pre[contp].power {
			if control_pre[contp].addend > 0 {
				working.addTwNiels(precmp[control_pre[contp].addend>>1])
			} else {
				working.subTwNiels(precmp[(-control_pre[contp].addend)>>1])
			}
			contp++
		}
	}
}
