package ed448

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) Test_NewHomogeneousPoint(c *C) {
	//Base point
	gx := []byte{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	}
	gy := []byte{0x13}

	basePoint, err := newPoint(gx, gy)
	c.Assert(err, IsNil)
	c.Assert(basePoint.OnCurve(), Equals, true)

	p := basePoint.double()
	c.Assert(p.OnCurve(), Equals, true)

	q := basePoint.add(basePoint)
	c.Assert(q.OnCurve(), Equals, true)
}

func (s *Ed448Suite) Test_NewExtensiblePoint(c *C) {
	x := &bigNumber{
		0x034365c8, 0x06b2a874, 0x0eb875d7, 0x0ae4c7a7,
		0x0785df04, 0x09929351, 0x01fe8c3b, 0x0f2a0e5f,
		0x0111d39c, 0x07ab52ba, 0x01df4552, 0x01d87566,
		0x0f297be2, 0x027c090f, 0x0a81b155, 0x0d1a562b,
	}
	y := &bigNumber{
		0x00da9708, 0x0e7d583e, 0x0dbcc099, 0x0d2dad89,
		0x05a49ce4, 0x01cb4ddc, 0x0928d395, 0x0098d91d,
		0x0bff16ce, 0x06f02f9a, 0x0ce27cc1, 0x0dab5783,
		0x0b553d94, 0x03251a0c, 0x064d70fb, 0x07fe3a2f,
	}

	exp := &extensibleCoordinates{
		&bigNumber{
			0x034365c8, 0x06b2a874, 0x0eb875d7, 0x0ae4c7a7,
			0x0785df04, 0x09929351, 0x01fe8c3b, 0x0f2a0e5f,
			0x0111d39c, 0x07ab52ba, 0x01df4552, 0x01d87566,
			0x0f297be2, 0x027c090f, 0x0a81b155, 0x0d1a562b,
		},
		&bigNumber{
			0x00da9708, 0x0e7d583e, 0x0dbcc099, 0x0d2dad89,
			0x05a49ce4, 0x01cb4ddc, 0x0928d395, 0x0098d91d,
			0x0bff16ce, 0x06f02f9a, 0x0ce27cc1, 0x0dab5783,
			0x0b553d94, 0x03251a0c, 0x064d70fb, 0x07fe3a2f,
		},
		&bigNumber{0x01},
		&bigNumber{
			0x034365c8, 0x06b2a874, 0x0eb875d7, 0x0ae4c7a7,
			0x0785df04, 0x09929351, 0x01fe8c3b, 0x0f2a0e5f,
			0x0111d39c, 0x07ab52ba, 0x01df4552, 0x01d87566,
			0x0f297be2, 0x027c090f, 0x0a81b155, 0x0d1a562b,
		},
		&bigNumber{
			0x00da9708, 0x0e7d583e, 0x0dbcc099, 0x0d2dad89,
			0x05a49ce4, 0x01cb4ddc, 0x0928d395, 0x0098d91d,
			0x0bff16ce, 0x06f02f9a, 0x0ce27cc1, 0x0dab5783,
			0x0b553d94, 0x03251a0c, 0x064d70fb, 0x07fe3a2f,
		},
	}
	p := newExtensible(x, y)
	c.Assert(p, DeepEquals, exp)
	//c.Assert(basePoint.OnCurve(), Equals, true)
}

func (s *Ed448Suite) Test_MixedAddition(c *C) {
	pa, _ := hex.DecodeString("4b8a632c1feab72769cd96e7aaa577861871b3613945c802b89377e8b85331ecc0ffb1cb20169bfc9c27274d38b0d01e87a1d5d851770bc8")
	pb, _ := hex.DecodeString("81a45f02f41053f8d7d2a1f176a340529b33b7ee4d3fa84de384b750b35a54c315bf36c41d023ade226449916e668396589ea2145da09b95")
	pc, _ := hex.DecodeString("5f5a2b06a2dbf7136f8dc979fd54d631ca7de50397250a196d3be2a721ab7cbaa92c545d9b15b5319e11b64bc031666049d8637e13838b3b")

	n := &twNiels{
		a: new(bigNumber).setBytes(pa),
		b: new(bigNumber).setBytes(pb),
		c: new(bigNumber).setBytes(pc),
	}

	px, _ := hex.DecodeString("e45b0207cf5036bcb75a775cb4eb3e8312a8d2b6c9c309dc6a589d2824427848e1ccc7ddac1a53d028375ff6b329d9f0998ed9bb4c81b4e9")
	py, _ := hex.DecodeString("e7c9798862329c3db188697a564706eade026ad6c773ca35069fd53f5d36c0b9db9fbda22386702aae4694ea2dfbe5e97458dd9040b2b97f")
	pz, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	pt, _ := hex.DecodeString("e45b0207cf5036bcb75a775cb4eb3e8312a8d2b6c9c309dc6a589d2824427848e1ccc7ddac1a53d028375ff6b329d9f0998ed9bb4c81b4e9")
	pu, _ := hex.DecodeString("e7c9798862329c3db188697a564706eade026ad6c773ca35069fd53f5d36c0b9db9fbda22386702aae4694ea2dfbe5e97458dd9040b2b97f")

	tx := &twExtensible{
		new(bigNumber).setBytes(px),
		new(bigNumber).setBytes(py),
		new(bigNumber).setBytes(pz),
		new(bigNumber).setBytes(pt),
		new(bigNumber).setBytes(pu),
	}

	px, _ = hex.DecodeString("ac7cd8b31f6f031e0ee8d606c12dafd6503305fd398f55399e3a677543dd7c8239995872432b428aba728f99d7bef231cb32a125ea9e7a69")
	py, _ = hex.DecodeString("4392419b29d484975fbb27272745408e81ae8ab9d1818c8181c39637437b3e456ec78335e6637c20df95a708da038a42213f29079ff1a457")
	pz, _ = hex.DecodeString("32aebf8745d98e59b0b4ad64b1f5e3ef67413c6ad82993ab24f1d8102f6da3c4206729e3392807daf13acc980447545ab91ce1ec6a7f5728")
	pt, _ = hex.DecodeString("b5c30299cb7bd622abdfc0a47e95b00cfc037561016114052a7bcf8862f7e4c656885e0b1800cc4f7d9046592912e2e0ee12662edc2588e6")
	pu, _ = hex.DecodeString("cf4f0c8cbed27e95dbbef285d2f6d68ba9819fdf44c1e09e0e4fc8be9f5b94aeac0c10886de2fa80f688a45a082861813dcf5afc7cd9f820")

	expected := &twExtensible{
		new(bigNumber).setBytes(px),
		new(bigNumber).setBytes(py),
		new(bigNumber).setBytes(pz),
		new(bigNumber).setBytes(pt),
		new(bigNumber).setBytes(pu),
	}

	ret := tx.addTwNiels(n)
	c.Assert(ret.equals(expected), Equals, true)
}

func (s *Ed448Suite) Test_ExtensibleUntwistAndDoubleAndSerialize(c *C) {
	px, _ := hex.DecodeString("4ed74e709fb89daba40d2aad54b8befa01e3cc2cd9eee3d72f9869a2897e5e44c32990e0366df5da4d36a890f10835a1ff85db9058b346b8")
	py, _ := hex.DecodeString("79c2294410f6371b2074d4ce8c40e366ebcf3770f45867e2280de6cb5e7da2c9e9c53a3ba0e9e38af58ac04092ef2a4d09510502adab1b90")
	pz, _ := hex.DecodeString("0b629561746bb03a5a1806376c6e424d51c704677885fc9947e3ae97d9146726dafa80b16a53f9bf492982b997466bf1c36e0ebaea3c7feb")
	pt, _ := hex.DecodeString("04073f6f22d607005b286fe02183753ffaf9c16d39e4d14b4291e8995cbb638fc123f0276ed08a394605221b0d76b87c80d92e327e49815a")
	pu, _ := hex.DecodeString("1531409e631a1e5f630426b33faf8d7a4f61653b32e4116bbf6cb4e170c143a887c2789a3409bcc5c2bbc3540e5b30a00050b83bfa04ae27")

	p := &twExtensible{
		new(bigNumber).setBytes(px),
		new(bigNumber).setBytes(py),
		new(bigNumber).setBytes(pz),
		new(bigNumber).setBytes(pt),
		new(bigNumber).setBytes(pu),
	}

	b, _ := hex.DecodeString("b690c6bcccee269215e1d7b86728e410ad4f6d1b933acaccf9e3b5b25c81cfe7e3c225e0f24afe060f3160f33cde18df3e6317db48c61aa5")
	exp := new(bigNumber).setBytes(b)

	ser := p.untwistAndDoubleAndSerialize()

	c.Assert(ser.decafEq(exp), Equals, decafTrue)
}

func (s *Ed448Suite) Test_ConditionalNegate(c *C) {
	pa, _ := hex.DecodeString("4b8a632c1feab72769cd96e7aaa577861871b3613945c802b89377e8b85331ecc0ffb1cb20169bfc9c27274d38b0d01e87a1d5d851770bc8")
	pb, _ := hex.DecodeString("81a45f02f41053f8d7d2a1f176a340529b33b7ee4d3fa84de384b750b35a54c315bf36c41d023ade226449916e668396589ea2145da09b95")
	pc, _ := hex.DecodeString("5f5a2b06a2dbf7136f8dc979fd54d631ca7de50397250a196d3be2a721ab7cbaa92c545d9b15b5319e11b64bc031666049d8637e13838b3b")

	n := &twNiels{
		a: new(bigNumber).setBytes(pa),
		b: new(bigNumber).setBytes(pb),
		c: new(bigNumber).setBytes(pc),
	}

	negN := &twNiels{
		a: n.b.copy(),
		b: n.a.copy(),
		c: new(bigNumber).neg(n.c.copy()),
	}

	x := n.copy()
	x.conditionalNegate(lmask)
	c.Assert(x, DeepEquals, negN)
}
func (s *Ed448Suite) Test_MontgomerySerialize(c *C) {
	bsIn, _ := hex.DecodeString("d03786c1b949c8e1b6046c527542ff55e9acda5c6fe8c7fef9c499ad182e4d84701555454c3ed9d10ff7b95cc4dd94b29c519dc51c29e80e")
	bsZ0, _ := hex.DecodeString("e281b05e4051a52b331430897d9d950529a46637d3ca1f45e1d2dc4fbd164c956f25dd0cf30458b4129e900faa2ba9b8d305dc4ae1e1b343")
	bsXd, _ := hex.DecodeString("c88f896abf42ca2cbff1edf881d1246ee76abe7385932d7b54fb9d71307fdd8043d8a80c7d0363e7a45443d4e9a03bf3e0aab82fb4714c5f")
	bsZd, _ := hex.DecodeString("962fa8b019eeedd607eda6b44454e17b76b1536f6b336362257d72c3c1576339514f1f4d2d0ae7b0680469a432a2f54cb7f9dbc14473802d")
	bsXa, _ := hex.DecodeString("09e41fe2e74667a6676fb0492b496f7d69d45055601ec86839b95e9343407ed592ea357118e5568eea272e9349adf0efbe29307187cfff6e")
	bsZa, _ := hex.DecodeString("b115a615745fc6f453a43d1466e12acd2215ac373cadcd633211235510c6a04c4f041006d07f543f2bd4b050ecdd472be4415ab7a3f79f95")
	in := new(bigNumber).setBytes(bsIn)
	mont := &montgomery{
		new(bigNumber).setBytes(bsZ0),
		new(bigNumber).setBytes(bsXd),
		new(bigNumber).setBytes(bsZd),
		new(bigNumber).setBytes(bsXa),
		new(bigNumber).setBytes(bsZa),
	}

	bsExp, _ := hex.DecodeString("322d71661943b5e080abed64d9ed331874a975329aaf9b42815e793ac08691e478fe559b29593a5413d5a4475e3ae0735a6d9bc1dc192b7d")
	exp := new(bigNumber).setBytes(bsExp)

	out, _ := mont.serialize(in)

	c.Assert(out.decafEq(exp), Equals, decafTrue)
}

func (s *Ed448Suite) Test_MontgomeryDeserialize(c *C) {
	bsIn, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008")
	in := new(bigNumber).setBytes(bsIn)
	out := new(montgomery)
	out.deserialize(in)
	bsZ0, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040")
	z0 := new(bigNumber).setBytes(bsZ0)
	bsXd, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	xd := new(bigNumber).setBytes(bsXd)
	bsZd, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	zd := new(bigNumber).setBytes(bsZd)
	bsXa, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	xa := new(bigNumber).setBytes(bsXa)
	bsZa, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040")
	za := new(bigNumber).setBytes(bsZa)

	c.Assert(out.z0.decafEq(z0), Equals, decafTrue)
	c.Assert(out.xd.decafEq(xd), Equals, decafTrue)
	c.Assert(out.zd.decafEq(zd), Equals, decafTrue)
	c.Assert(out.xa.decafEq(xa), Equals, decafTrue)
	c.Assert(out.za.decafEq(za), Equals, decafTrue)
}

func (s *Ed448Suite) Test_MontgomeryStep(c *C) {
	bsZ0, _ := hex.DecodeString("e281b05e4051a52b331430897d9d950529a46637d3ca1f45e1d2dc4fbd164c956f25dd0cf30458b4129e900faa2ba9b8d305dc4ae1e1b343")
	bsXd, _ := hex.DecodeString("dc7c2264cf2a3f6178ee7884793f2d0cfe98e602c32adfbec9a5fc225c904f5e1f45c614fc483aec252745e04a38f49e1a4cfc0e8bbf14c5")
	bsZd, _ := hex.DecodeString("76ad01dbfd7dd72671ad1f827b762fe0c39c808084533b1e22ee18537b7e43c75b995f9e107ec055fbb3df4fb83ad78e69de76a188fb6db6")
	bsXa, _ := hex.DecodeString("86032c9f990e2680726003f62a1ec5c01f18ad130ce0883b247d2ea9e8d591e6121e6007027d44d94d9659a05fb47e91c3c11b5552cb2185")
	bsZa, _ := hex.DecodeString("5c157ed45b60be2db18a494780b5b7ab79ae1afc3919c9b00c1879495ea079b73990eebf5f0def2897fe8ca78084c07ef89c5bfc336625fd")
	mont := &montgomery{
		new(bigNumber).setBytes(bsZ0),
		new(bigNumber).setBytes(bsXd),
		new(bigNumber).setBytes(bsZd),
		new(bigNumber).setBytes(bsXa),
		new(bigNumber).setBytes(bsZa),
	}
	mont.montgomeryStep()
	bsZ0, _ = hex.DecodeString("e281b05e4051a52b331430897d9d950529a46637d3ca1f45e1d2dc4fbd164c956f25dd0cf30458b4129e900faa2ba9b8d305dc4ae1e1b343")
	bsXd, _ = hex.DecodeString("c88f896abf42ca2cbff1edf881d1246ee76abe7385932d7b54fb9d71307fdd8043d8a80c7d0363e7a45443d4e9a03bf3e0aab82fb4714c5f")
	bsZd, _ = hex.DecodeString("962fa8b019eeedd607eda6b44454e17b76b1536f6b336362257d72c3c1576339514f1f4d2d0ae7b0680469a432a2f54cb7f9dbc14473802d")
	bsXa, _ = hex.DecodeString("09e41fe2e74667a6676fb0492b496f7d69d45055601ec86839b95e9343407ed592ea357118e5568eea272e9349adf0efbe29307187cfff6e")
	bsZa, _ = hex.DecodeString("b115a615745fc6f453a43d1466e12acd2215ac373cadcd633211235510c6a04c4f041006d07f543f2bd4b050ecdd472be4415ab7a3f79f95")
	exp := &montgomery{
		new(bigNumber).setBytes(bsZ0),
		new(bigNumber).setBytes(bsXd),
		new(bigNumber).setBytes(bsZd),
		new(bigNumber).setBytes(bsXa),
		new(bigNumber).setBytes(bsZa),
	}

	c.Assert(mont.z0.equals(exp.z0), Equals, true)
	c.Assert(mont.xd.equals(exp.xd), Equals, true)
	c.Assert(mont.zd.equals(exp.zd), Equals, true)
	c.Assert(mont.xa.equals(exp.xa), Equals, true)
	c.Assert(mont.za.equals(exp.za), Equals, true)
}

func (s *Ed448Suite) Test_AddTwNiels(c *C) {
	na, _ := hex.DecodeString("33d7e1341e2291816fa27efbac283c2d0fae711d29b581200d215449fa64ef98a767887486155176a543fc08807a595766b7987e4b4c037f")
	nb, _ := hex.DecodeString("4ad8ff3e6b86b69d349faa7cca6280ed8208997607ed60c842651c0ddac0754664433340bd3e4253dd8565c36713f7ca2c11023891708535")
	nc, _ := hex.DecodeString("8bd294a6cfdee12764081c4e9acaab981fcf3b8bd422f683d37a175081eeaec3a1b5c42dc5b962e5a46a0959b1f725796637306f8723066c")
	n := &twNiels{
		new(bigNumber).setBytes(na),
		new(bigNumber).setBytes(nb),
		new(bigNumber).setBytes(nc),
	}
	ex, _ := hex.DecodeString("779f2f66bbb61f027dc9dcf8207caf539cae89a43499b0840c3fb1d6af841fb7b797642d822d49fb41d55aed37cec81e0aadff753c57ac5d")
	ey, _ := hex.DecodeString("8a38832663ae92cff2a70d6b96af28f1a55460959db9286a5498b88cc153d843369fe02fa4848bc02b9d36dc572b3bf93ee75d50c053bc37")
	ez, _ := hex.DecodeString("eb67f18c89fab163f7ca1396d67a9405f7a91bd6aa9858ac7b784ea82f2b41b255b46ebc4a152dd8eb69ca0250923da7e76cd5fd2360b8cc")
	et, _ := hex.DecodeString("81e01268e0589b800037589f8d1298d039b9ddf57928b3333958ae0fa593e358e8a1df5a0e333a1c4e5976dd95ca3dff2293314e14c498cf")
	eu, _ := hex.DecodeString("8da96259e6d38108ae8c007410371d933fd7209b36e910d47db044a444d61ac4df9649d12ffc3b3a9f6fa79dc7f1fa03e44d3073551d8442")
	e := &twExtensible{
		new(bigNumber).setBytes(ex),
		new(bigNumber).setBytes(ey),
		new(bigNumber).setBytes(ez),
		new(bigNumber).setBytes(et),
		new(bigNumber).setBytes(eu),
	}

	ex, _ = hex.DecodeString("f6e3e13f7662bdf0f468fe98062cf0152a02e35a4f49cb28debc24d2ce9eae08a2ce023c9df521faa06545490e14608a62e59dc5c1c9d3c7")
	ey, _ = hex.DecodeString("a93634cd0770e1846d8280275b9dc0e2a7636eca8c8fbe50edda6fc8966fc26d63d8b2ead7df70e81fb30b2e36c1d0fb0541359bbf2d7b6b")
	ez, _ = hex.DecodeString("567ce9ccf084022de2a1017524ff1dfe9cea601978db5d84e017b473bea82057ed0a58be4567c30fa649126c9dfcb3e083ea6bbe50e1b95c")
	et, _ = hex.DecodeString("b9e37cdbd0622a8d3a0e2dc8bda3e014e7c4e6c159c18a7c25076f9aab26022340c228ff13c5f5be52cabfdeae1bfd4cc4b7be572253d0ad")
	eu, _ = hex.DecodeString("0a789316f892f38685e32b8ee63fea3f5dc90459c6ea3557a88720772c752aba93652c500bee3e9651ed94437a5ba41eb5336e22f5a7a4c5")

	exp := &twExtensible{
		new(bigNumber).setBytes(ex),
		new(bigNumber).setBytes(ey),
		new(bigNumber).setBytes(ez),
		new(bigNumber).setBytes(et),
		new(bigNumber).setBytes(eu),
	}

	e.addTwNiels(n)

	c.Assert(e.x.equals(exp.x), Equals, true)
	c.Assert(e.y.equals(exp.y), Equals, true)
	c.Assert(e.z.equals(exp.z), Equals, true)
	c.Assert(e.t.equals(exp.t), Equals, true)
	c.Assert(e.u.equals(exp.u), Equals, true)
}

func (s *Ed448Suite) Test_DeserializeAndTwistAprox(c *C) {
	b, _ := hex.DecodeString("d03786c1b949c8e1b6046c527542ff55e9acda5c6fe8c7fef9c499ad182e4d84701555454c3ed9d10ff7b95cc4dd94b29c519dc51c29e80e")
	n := new(bigNumber).setBytes(b)

	ex, _ := hex.DecodeString("4d8b77dc973a1f9bcd5358c702ee8159a71cd3e4c1ff95bfb30e7038cffe9f794211dffd758e2a2a693a08a9a454398fde981e5e2669acad")
	ey, _ := hex.DecodeString("27193fda68a08730d1def89d64c7f466d9e3d0ac89d8fdcd17b8cdb446e80404e8cd715d4612c16f70803d50854b66c9b3412e85e2f19b0d")
	ez, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	et, _ := hex.DecodeString("4d8b77dc973a1f9bcd5358c702ee8159a71cd3e4c1ff95bfb30e7038cffe9f794211dffd758e2a2a693a08a9a454398fde981e5e2669acad")
	eu, _ := hex.DecodeString("27193fda68a08730d1def89d64c7f466d9e3d0ac89d8fdcd17b8cdb446e80404e8cd715d4612c16f70803d50854b66c9b3412e85e2f19b0d")
	exp := &twExtensible{
		new(bigNumber).setBytes(ex),
		new(bigNumber).setBytes(ey),
		new(bigNumber).setBytes(ez),
		new(bigNumber).setBytes(et),
		new(bigNumber).setBytes(eu),
	}

	tw, ok := n.deserializeAndTwistApprox()

	c.Assert(tw.x.equals(exp.x), Equals, true)
	c.Assert(tw.y.equals(exp.y), Equals, true)
	c.Assert(tw.z.equals(exp.z), Equals, true)
	c.Assert(tw.t.equals(exp.t), Equals, true)
	c.Assert(tw.u.equals(exp.u), Equals, true)
	c.Assert(ok, Equals, true)
}

func (s *Ed448Suite) Test_UntwistDoubleAndSerialize(c *C) {
	ex, _ := hex.DecodeString("d902fadbeee8dd1ef391dcce59cc75d286c9efc7229dd919a35236a5447384e84617bf94d4129af02d7667fad1df88985132c1ce1b133428")
	ey, _ := hex.DecodeString("ba1d18df944a527ec4ebad9c84cc32643064dcd26bf003a9763dad575104e1a3c9fbb02f971169c2736ed5d8812ad8eeedcfa8226977ddb4")
	ez, _ := hex.DecodeString("2d35e8b251eb6b421291cf3a466597759059e01b7cc89f332f96f801ced244299f4da20b9fcedbaa66c5fd3508dcb61888e2b89bee4fea45")
	et, _ := hex.DecodeString("8713cc3806a247771ae8567b3b73dd874a8261a610de7c34202fab877f15213120e2fd14e5b191663c1e62d404c54b9f63e1e2e3d98eafb2")
	eu, _ := hex.DecodeString("eafb1cd470e2728ee254c7a312092e820656c14a993f2896479aa211b0a1bb515deee36d06acee20a40a1cad5dc5cc38072cdd63447587e9")

	tw := &twExtensible{
		new(bigNumber).setBytes(ex),
		new(bigNumber).setBytes(ey),
		new(bigNumber).setBytes(ez),
		new(bigNumber).setBytes(et),
		new(bigNumber).setBytes(eu),
	}

	b, _ := hex.DecodeString("47e65d055a310a0e094908bb39194bc7d18589a1c93d5a2560ab7d38f4557cabb54d23a0801fad7cdbafc7b6457b0142b9394c71a8048666")
	expected := new(bigNumber).setBytes(b)

	untw := tw.untwistAndDoubleAndSerialize()

	c.Assert(untw.equals(expected), Equals, true)
}

func (s *Ed448Suite) Test_ConvertPNielsToExtended(c *C) {
	pn := &twPNiels{
		&twNiels{
			&bigNumber{
				0x00923994, 0x09e0df3c, 0x05d44abd, 0x0e120043,
				0x0802cd52, 0x03d6dc68, 0x06848872, 0x0b6dc339,
				0x0847737b, 0x0bc369ba, 0x02096051, 0x0bb48580,
				0x07f1feb0, 0x01f3325e, 0x05347dc3, 0x0ac3e06b,
			},
			&bigNumber{
				0x0ebfe2a0, 0x0eb5fbb2, 0x0dc6ec15, 0x0668e6e7,
				0x0a2b46d6, 0x09716cff, 0x0b1f6161, 0x075f2e4b,
				0x092bafec, 0x04a09fe7, 0x0d0693b5, 0x0181e9e6,
				0x0ae2011c, 0x0339acaa, 0x05ab0851, 0x02e7f480,
			},
			&bigNumber{
				0x01034f2a, 0x0eae9078, 0x0f3efa41, 0x039389e4,
				0x04f19d1d, 0x05e8f21c, 0x081e279d, 0x01011eca,
				0x06f90064, 0x068cfe2b, 0x08b95f11, 0x0144de8a,
				0x04022ab1, 0x087d7dca, 0x0129b8ca, 0x0c3a54aa,
			},
		},
		&bigNumber{
			0x0cc630cf, 0x04be9e7e, 0x02806f73, 0x033000c5,
			0x08901863, 0x06198942, 0x057694d2, 0x0c4ea1d8,
			0x0c5c8ca4, 0x0ba40152, 0x0c694dd1, 0x085ae4e6,
			0x05bc4f74, 0x0db70b9f, 0x09da2b77, 0x0a6633ef,
		},
	}

	expected := &twExtendedPoint{
		&bigNumber{
			0x0f28f6b9, 0x031242e7, 0x025075a3, 0x089c6fef,
			0x0c7aad41, 0x08072e8b, 0x00d25064, 0x0fee9931,
			0x039d8756, 0x08c9750e, 0x02a13cdd, 0x05f83a36,
			0x0f15eaf3, 0x08641e57, 0x0283f647, 0x0e176bf0,
		},
		&bigNumber{
			0x0538a2be, 0x05dccef4, 0x0e1f3db4, 0x0987f28b,
			0x025654bb, 0x0015ccbc, 0x0ad1978a, 0x00e352a8,
			0x0b404516, 0x093ebe46, 0x0cef87de, 0x09d51db5,
			0x09936823, 0x075b8fae, 0x0b49e190, 0x029b4497,
		},
		&bigNumber{
			0x026644d3, 0x03978937, 0x02d980ab, 0x0932cef7,
			0x01f40638, 0x07e71e49, 0x0cb9ac94, 0x0db56f31,
			0x008d6c2b, 0x01eaad17, 0x0be10c32, 0x04f7381e,
			0x0fffe02c, 0x0f5d619b, 0x0020e408, 0x02d2ddba,
		},
		&bigNumber{
			0x09892a57, 0x0fe01c67, 0x00431201, 0x040373e2,
			0x0e08643f, 0x01adb24a, 0x0ec4e7ad, 0x0d41c3ed,
			0x02d82f6d, 0x0f010b33, 0x0e0cb94a, 0x08668c81,
			0x0549b386, 0x0efbdfa8, 0x0cf765f1, 0x01332bf3,
		},
	}

	c.Assert(pn.toExtendedPoint(), DeepEquals, expected)
}

func compareNumbers(label string, n *bigNumber, b *big.Int) {
	s := [fieldBytes]byte{}
	serialize(s[:], n)

	r := rev(s[:])
	bs := b.Bytes()

	for i := len(r) - len(bs); i > 0; i-- {
		bs = append([]byte{0}, bs...)
	}

	if !bytes.Equal(r, bs) {
		fmt.Printf("%s does not match!\n\t%#v\n\n vs\n\n\t%#v\n", label, r, bs)
	}

}
