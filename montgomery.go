package ed448

type montgomery struct {
	z0, xd, zd, xa, za *bigNumber
}

func (a *montgomery) montgomeryStep() {
	L0 := new(bigNumber)
	L1 := new(bigNumber)
	L0.addRaw(a.zd, a.xd)
	L1.subxRaw(a.xd, a.zd)
	a.zd.subxRaw(a.xa, a.za)
	karatsubaMul(a.xd, L0, a.zd)
	a.zd.addRaw(a.za, a.xa)
	karatsubaMul(a.za, L1, a.zd)
	a.xa.addRaw(a.za, a.xd)
	karatsubaSquare(a.zd, a.xa)
	karatsubaMul(a.xa, a.z0, a.zd)
	a.zd.subxRaw(a.xd, a.za)
	karatsubaSquare(a.za, a.zd)
	karatsubaSquare(a.xd, L0)
	karatsubaSquare(L0, L1)
	a.zd.mulWSignedCurveConstant(a.xd, 1-curveDSigned) /* FIXME PERF MULW */
	L1.subxRaw(a.xd, L0)
	karatsubaMul(a.xd, L0, a.zd)
	L0.subRaw(a.zd, L1)
	L0.bias(4 - 2*1 /*is32 ? 2 : 4*/)
	//XXX 64bits don't need this reduce
	L0.weakReduce()
	karatsubaMul(a.zd, L0, L1)
}

func (a *montgomery) serialize(sbz *bigNumber) (b *bigNumber, ok uint32) {
	L0 := new(bigNumber)
	L1 := new(bigNumber)
	L2 := new(bigNumber)
	L3 := new(bigNumber)
	b = new(bigNumber)

	L3.mulCopy(a.z0, a.zd)
	L1.sub(L3, a.xd)
	L3.mulCopy(a.za, L1)
	L2.mulCopy(a.z0, a.xd)
	L1.sub(L2, a.zd)
	L0.mulCopy(a.xa, L1)
	L2.add(L0, L3)
	L1.sub(L3, L0)
	L3.mulCopy(L1, L2)
	L2 = a.z0.copy()
	L2.addW(1)
	L0.squareCopy(L2)
	L1.mulWSignedCurveConstant(L0, curveDSigned-1)
	L2.add(a.z0, a.z0)
	L0.add(L2, L2)
	L2.add(L0, L1)
	L0.mulCopy(a.xd, L2)
	L5 := a.zd.zeroMask()
	L6 := -L5

	// constant_time_mask ( L1, L0, sizeof(L1), L5 );
	mask(L1, L0, L5)
	L2.add(L1, a.zd)
	L4 := ^L5
	L1.mulCopy(sbz, L3)
	L1.addW(L6)
	L3.mulCopy(L2, L1)
	L1.mulCopy(L3, L2)
	L2.mulCopy(L3, a.xd)
	L3.mulCopy(L1, L2)
	L0.isr(L3)
	L2.mulCopy(L1, L0)
	L1.squareCopy(L0)
	L0.mulCopy(L3, L1)

	// constant_time_mask ( b, L2, sizeof(L1), L4 );
	mask(b, L2, L4)
	L0.subW(1)
	L5 = L0.zeroMask()
	L4 = sbz.zeroMask()

	return b, L5 | L4
}

func (a *montgomery) deserialize(sz *bigNumber) {
	a.z0 = new(bigNumber).squareCopy(sz)
	a.xd = new(bigNumber).setUi(1)
	a.zd = new(bigNumber).setUi(0)
	a.xa = new(bigNumber).setUi(1)
	a.za = a.z0.copy()
}
