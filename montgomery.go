package ed448

type montgomery struct {
	z0, xd, zd, xa, za *bigNumber
}

func (a *montgomery) montgomeryStep() {
	L0 := new(bigNumber)
	L1 := new(bigNumber)
	L0.addRaw(a.zd, a.xd)
	L1.sub(a.xd, a.zd)
	a.zd.sub(a.xa, a.za)
	a.xd.mul(L0, a.zd)
	a.zd.addRaw(a.za, a.xa)
	a.za.mul(L1, a.zd)
	a.xa.addRaw(a.za, a.xd)
	a.zd.square(a.xa)
	a.xa.mul(a.z0, a.zd)
	a.zd.sub(a.xd, a.za)
	a.za.square(a.zd)
	a.xd.square(L0)
	L0.square(L1)
	a.zd.mulWSignedCurveConstant(a.xd, 1-edwardsD) /* TODO: PERF MULW */
	L1.sub(a.xd, L0)
	a.xd.mul(L0, a.zd)
	L0.subRaw(a.zd, L1)
	L0.bias(4 - 2*1 /*is32 ? 2 : 4*/)
	//TODO 64bits don't need this reduce
	L0.weakReduce()
	a.zd.mul(L0, L1)
}

func (a *montgomery) serialize(sbz *bigNumber) (b *bigNumber, ok word) {
	L0 := new(bigNumber)
	L1 := new(bigNumber)
	L2 := new(bigNumber)
	L3 := new(bigNumber)
	b = new(bigNumber)

	L3.mul(a.z0, a.zd)
	L1.sub(L3, a.xd)
	L3.mul(a.za, L1)
	L2.mul(a.z0, a.xd)
	L1.sub(L2, a.zd)
	L0.mul(a.xa, L1)
	L2.add(L0, L3)
	L1.sub(L3, L0)
	L3.mul(L1, L2)
	L2 = a.z0.copy()
	L2.addW(1)
	L0.square(L2)
	L1.mulWSignedCurveConstant(L0, edwardsD-1)
	L2.add(a.z0, a.z0)
	L0.add(L2, L2)
	L2.add(L0, L1)
	L0.mul(a.xd, L2)
	L5 := a.zd.zeroMask()
	L6 := -L5

	// constant_time_mask ( L1, L0, sizeof(L1), L5 );
	mask(L1, L0, L5)
	L2.add(L1, a.zd)
	L4 := ^L5
	L1.mul(sbz, L3)
	L1.addW(L6)
	L3.mul(L2, L1)
	L1.mul(L3, L2)
	L2.mul(L3, a.xd)
	L3.mul(L1, L2)
	L0.isr(L3)
	L2.mul(L1, L0)
	L1.square(L0)
	L0.mul(L3, L1)

	// constant_time_mask ( b, L2, sizeof(L1), L4 );
	mask(b, L2, L4)
	L0.subW(1)
	L5 = L0.zeroMask()
	L4 = sbz.zeroMask()

	return b, L5 | L4
}

func (a *montgomery) deserialize(sz *bigNumber) {
	a.z0 = new(bigNumber).square(sz)
	a.xd = new(bigNumber).setUI(1)
	a.zd = new(bigNumber).setUI(0)
	a.xa = new(bigNumber).setUI(1)
	a.za = a.z0.copy()
}
