package ed448

import (
	"encoding/hex"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestWNAFSMultiplication(c *C) {
	px, _ := hex.DecodeString("4d8b77dc973a1f9bcd5358c702ee8159a71cd3e4c1ff95bfb30e7038cffe9f794211dffd758e2a2a693a08a9a454398fde981e5e2669acad")
	py, _ := hex.DecodeString("27193fda68a08730d1def89d64c7f466d9e3d0ac89d8fdcd17b8cdb446e80404e8cd715d4612c16f70803d50854b66c9b3412e85e2f19b0d")
	pz, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	pt, _ := hex.DecodeString("4d8b77dc973a1f9bcd5358c702ee8159a71cd3e4c1ff95bfb30e7038cffe9f794211dffd758e2a2a693a08a9a454398fde981e5e2669acad")
	pu, _ := hex.DecodeString("27193fda68a08730d1def89d64c7f466d9e3d0ac89d8fdcd17b8cdb446e80404e8cd715d4612c16f70803d50854b66c9b3412e85e2f19b0d")

	p := &twExtensible{
		new(bigNumber).setBytes(px),
		new(bigNumber).setBytes(py),
		new(bigNumber).setBytes(pz),
		new(bigNumber).setBytes(pt),
		new(bigNumber).setBytes(pu),
	}

	x := [fieldWords]word_t{
		0x6c226d73, 0x70edcfc3,
		0x44156c47, 0x084f4695,
		0xe72606ac, 0x9d0ce5e5,
		0xed96d3ba, 0x9ff3fa11,
		0x4a15c383, 0xca38a0af,
		0xead789b3, 0xb96613ba,
		0x48ba4461, 0x34eb2031,
	}

	y := [fieldWords]word_t{
		0x2118b8c6, 0x4356acd5,
		0x26d7e73c, 0x459174b7,
		0xf10bea31, 0x83e528bb,
		0xb960d695, 0xd0da7e28,
		0xbad7f9a1, 0xe9f5ba01,
		0x94ea1518, 0x12c58cca,
		0x302c76eb, 0x3bd0363e,
	}

	px, _ = hex.DecodeString("d902fadbeee8dd1ef391dcce59cc75d286c9efc7229dd919a35236a5447384e84617bf94d4129af02d7667fad1df88985132c1ce1b133428")
	py, _ = hex.DecodeString("ba1d18df944a527ec4ebad9c84cc32643064dcd26bf003a9763dad575104e1a3c9fbb02f971169c2736ed5d8812ad8eeedcfa8226977ddb4")
	pz, _ = hex.DecodeString("2d35e8b251eb6b421291cf3a466597759059e01b7cc89f332f96f801ced244299f4da20b9fcedbaa66c5fd3508dcb61888e2b89bee4fea45")
	pt, _ = hex.DecodeString("8713cc3806a247771ae8567b3b73dd874a8261a610de7c34202fab877f15213120e2fd14e5b191663c1e62d404c54b9f63e1e2e3d98eafb2")
	pu, _ = hex.DecodeString("eafb1cd470e2728ee254c7a312092e820656c14a993f2896479aa211b0a1bb515deee36d06acee20a40a1cad5dc5cc38072cdd63447587e9")
	expectedP := &twExtensible{
		new(bigNumber).setBytes(px),
		new(bigNumber).setBytes(py),
		new(bigNumber).setBytes(pz),
		new(bigNumber).setBytes(pt),
		new(bigNumber).setBytes(pu),
	}

	linear_combo_var_fixed_vt(p, x[:], y[:], wnfsTable[:])

	c.Assert(p.equals(expectedP), Equals, true)
}
