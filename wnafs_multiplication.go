package ed448

import (
	"math/bits"
)

type smvtControl struct {
	power, addend int
}

const bOver16 = uint(2)

func recodeWNAF2(control []smvtControl, s *scalar, tableBits uint) word {
	tableSize := 446/(tableBits+1) + 3
	position := tableSize - 1

	control[position].power = -1
	control[position].addend = 0
	position--

	current := uint64(s[0] & 0xFFFF)
	mask := uint32((1 << (tableBits + 1)) - 1)

	for w := uint(1); w < 30; w++ {
		if w < 28 {
			// Refill the 16 high bits of current
			current += uint64(uint32((s[w/bOver16] >> uint((16 * (w % bOver16)))) << 16))
		}

		for current&0xFFFF != 0 {
			pos := uint32(bits.TrailingZeros32(uint32(current)))
			odd := uint32(current) >> pos
			delta := int32(odd & mask)
			if odd&(1<<(tableBits+1)) != 0 {
				delta -= (1 << (tableBits + 1))
			}
			current = uint64(int64(current) - int64(delta<<pos))
			control[position].power = int(pos) + 16*(int(w)-1)
			control[position].addend = int(delta)
			position--
		}
		current >>= 16
	}

	position++
	n := uint(tableSize - position)
	for i := uint(0); i < n; i++ {
		control[i] = control[i+position]
	}
	return word(n - 1)
}

func recodeWNAF(control []smvtControl, s *scalar, nBits, tableBits uint) (position word) {
	current := 0
	var i, j int
	position = 0
	for i = int(nBits - 1); i >= 0; i-- {
		bit := (s[i/wordBits] >> uint(i%wordBits)) & 1
		current = (2 * current) + int(bit)

		// Sizing: |current| >= 2^(tableBits+1) -> |current| = 2^0
		// Current loses (tableBits+1) bits every time.  It otherwise gains
		// 1 bit per iteration.  The number of iterations is
		// (nbits + 2 + tableBits), and an additional control word is added at
		// the end.  So the total number of control words is at most
		// ceil((nbits+1) / (tableBits+1)) + 2 = floor((nbits)/(tableBits+1)) + 2.
		// There's also the stopper with power -1, for a total of +3.
		if current >= (2<<word(tableBits)) || current <= -1-(2<<word(tableBits)) {
			delta := (current + 1) >> 1 // |delta| < 2^tablebits
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

func (p *twExtendedPoint) prepareFixedWindow(nTable int) []*twPNiels {
	pOriginal := p.copy()
	pn := p.copy().doubleInternal(false).toPNiels()
	out := make([]*twPNiels, nTable)
	out[0] = pOriginal.toPNiels()
	for i := 1; i < nTable; i++ {
		pOriginal.addProjectiveNielsToExtended(pn, false)
		out[i] = pOriginal.toPNiels()
	}
	return out[:]
}

func prepareWNAFTable(dst []*twPNiels, p *twExtensible, tableSize uint) {
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

func decafPrepareWNAFTable(dst []*twPNiels, p *twExtendedPoint, tableSize uint) {
	dst[0] = p.toPNiels()

	if tableSize == 0 {
		return
	}

	p.doubleInternal(false)

	twOp := p.toPNiels()

	p.addProjectiveNielsToExtended(dst[0], false)
	dst[1] = p.toPNiels()

	for i := 2; i < 1<<tableSize; i++ {
		p.addProjectiveNielsToExtended(twOp, false)
		dst[i] = p.toPNiels()
	}
}

func linearComboVarFixedVt(working *twExtensible, scalarVar, scalarPre *scalar, precmp []*twNiels) {
	tableBitsVar := uint(4) //SCALARMUL_WNAF_COMBO_TABLE_BITS;
	nbitsVar := uint(446)
	nbitsPre := uint(446)
	tableBitsPre := uint(5)

	var controlVar [92]smvtControl // nbitsVar/(tableBitsVar+1)+3
	var controlPre [77]smvtControl // nbitsPre/(tableBitsPre+1)+3

	recodeWNAF(controlVar[:], scalarVar, nbitsVar, tableBitsVar)
	recodeWNAF(controlPre[:], scalarPre, nbitsPre, tableBitsPre)

	var precmpVar [16]*twPNiels // 1 << tableBitsVar
	prepareWNAFTable(precmpVar[:], working, tableBitsVar)

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

func doubleScalarMul(pointB, pointC *twExtendedPoint, scalarB, scalarC *scalar) *twExtendedPoint {
	const decafWindowBits = 5
	const window = decafWindowBits       //5
	const windowMask = (1 << window) - 1 //0x0001f 31
	const windowTMask = windowMask >> 1  //0x0000f 15
	const nTable = 1 << (window - 1)     //0x00010 16

	scalar1x := &scalar{}
	scalar1x.add(scalarB, decafPrecompTable.scalarAdjustment)
	scalar1x.halve(scalar1x)
	scalar2x := &scalar{}
	scalar2x.add(scalarC, decafPrecompTable.scalarAdjustment)
	scalar2x.halve(scalar2x)

	multiples1 := pointB.prepareFixedWindow(nTable)
	multiples2 := pointC.prepareFixedWindow(nTable)
	out := &twExtendedPoint{}
	first := true
	for i := scalarBits - ((scalarBits - 1) % window) - 1; i >= 0; i -= window {
		bits1 := scalar1x[i/wordBits] >> uint(i%wordBits)
		bits2 := scalar2x[i/wordBits] >> uint(i%wordBits)
		if i%wordBits >= wordBits-window && i/wordBits < scalarWords-1 {
			bits1 ^= scalar1x[i/wordBits+1] << uint(wordBits-(i%wordBits))
			bits2 ^= scalar2x[i/wordBits+1] << uint(wordBits-(i%wordBits))
		}
		bits1 &= windowMask
		bits2 &= windowMask
		inv1 := (bits1 >> (window - 1)) - 1
		inv2 := (bits2 >> (window - 1)) - 1
		bits1 ^= inv1
		bits2 ^= inv2
		//Add in from table.  Compute t only on last iteration.
		mul1pn := constTimeLookup(multiples1, word(bits1&windowTMask))
		mul1pn.n.conditionalNegate(inv1)
		if first {
			out = mul1pn.toExtendedPoint()
			first = false
		} else {
			//Using Hisil et al's lookahead method instead of extensible here
			//for no particular reason.  Double WINDOW times, but only compute t on
			//the last one.
			for j := 0; j < window-1; j++ {
				out.doubleInternal(true)
			}
			out.doubleInternal(false)
			out.addProjectiveNielsToExtended(mul1pn, false)
		}
		mul2pn := constTimeLookup(multiples2, word(bits2&windowTMask))
		mul2pn.n.conditionalNegate(inv2)
		if i > 0 {
			out.addProjectiveNielsToExtended(mul2pn, true)
		} else {
			out.addProjectiveNielsToExtended(mul2pn, false)
		}
	}
	return out
}

func decafDoubleNonSecretScalarMul(p *twExtendedPoint, scalarPre, scalarVar *scalar) *twExtendedPoint {
	tableBitsVar := uint(3) // DECAF_WNAF_VAR_TABLE_BITS
	tableBitsPre := uint(5) // DECAF_WNAF_FIXED_TABLE_BITS

	var controlVar [115]smvtControl // nbitsVar/(tableBitsVar+1)+3
	var controlPre [77]smvtControl  // nbitsPre/(tableBitsPre+1)+3

	recodeWNAF2(controlPre[:], scalarPre, tableBitsPre)
	recodeWNAF2(controlVar[:], scalarVar, tableBitsVar)

	var precmpVar [8]*twPNiels
	decafPrepareWNAFTable(precmpVar[:], p, tableBitsVar)

	contp := 0
	contv := 0

	index := controlVar[0].addend >> 1

	i := controlVar[0].power
	out := &twExtendedPoint{
		&bigNumber{0x00},
		&bigNumber{0x00},
		&bigNumber{0x00},
		&bigNumber{0x00},
	}

	if i > controlPre[0].power {
		out = precmpVar[index].toExtendedPoint()
		contv++
	} else if i == controlPre[0].power && i >= 0 {
		out = precmpVar[index].toExtendedPoint()
		out.addNielsToExtended(decafWnafsTable[controlPre[0].addend>>1], i != 0)
		contv++
		contp++
	} else {
		i = controlPre[0].power
		out = decafWnafsTable[controlPre[0].addend>>1].toExtended()
		contp++
	}

	if i < 0 {
		out.setIdentity()
		return out
	}

	for i--; i >= 0; i-- {

		cv := i == controlVar[contv].power
		cp := i == controlPre[contp].power

		out.doubleInternal(i != 0 && !(cv || cp))

		if cv {
			if controlVar[contv].addend > 0 {
				out.addProjectiveNielsToExtended(precmpVar[controlVar[contv].addend>>1], (i != 0 && !cp))
			} else {
				out.subProjectiveNielsFromExtendedPoint(precmpVar[(-controlVar[contv].addend)>>1], (i != 0 && !cp))
			}
			contv++
		}

		if cp {
			if controlPre[contp].addend > 0 {
				out.addNielsToExtended(decafWnafsTable[controlPre[contp].addend>>1], i != 0)
			} else {
				out.subNielsFromExtendedPoint(decafWnafsTable[(-controlPre[contp].addend)>>1], i != 0)
			}
			contp++
		}
	}
	return out
}
