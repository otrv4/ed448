package ed448

func linear_combo_var_fixed_vt(p *twExtensible, x, y []word_t, table []*twNiels) {
	//TODO
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
