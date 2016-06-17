package ed448

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) TestPoint(c *C) {
	//Base point
	gx := serialized{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	}
	gy := serialized{0x13}

	basePoint, err := NewPoint(gx, gy)
	c.Assert(err, IsNil)

	curve := newRadixCurve()
	c.Assert(curve.isOnCurve(basePoint), Equals, true)

	p := basePoint.Double()
	c.Assert(curve.isOnCurve(p), Equals, true)

	q := basePoint.Add(basePoint)
	c.Assert(curve.isOnCurve(q), Equals, true)
}

func (s *Ed448Suite) TestMixedAddition(c *C) {
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

func (s *Ed448Suite) TestExtensibleUntwistAndDoubleAndSerialize(c *C) {
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

	c.Assert(ser.equals(exp), Equals, true)
}

func (s *Ed448Suite) TestConditionalNegate(c *C) {
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
	x.conditionalNegate(0xffffffff)
	c.Assert(x, DeepEquals, negN)
}
func (s *Ed448Suite) TestMontgomerySerialize(c *C) {
	bs_in, _ := hex.DecodeString("d03786c1b949c8e1b6046c527542ff55e9acda5c6fe8c7fef9c499ad182e4d84701555454c3ed9d10ff7b95cc4dd94b29c519dc51c29e80e")
	bs_z0, _ := hex.DecodeString("e281b05e4051a52b331430897d9d950529a46637d3ca1f45e1d2dc4fbd164c956f25dd0cf30458b4129e900faa2ba9b8d305dc4ae1e1b343")
	bs_xd, _ := hex.DecodeString("c88f896abf42ca2cbff1edf881d1246ee76abe7385932d7b54fb9d71307fdd8043d8a80c7d0363e7a45443d4e9a03bf3e0aab82fb4714c5f")
	bs_zd, _ := hex.DecodeString("962fa8b019eeedd607eda6b44454e17b76b1536f6b336362257d72c3c1576339514f1f4d2d0ae7b0680469a432a2f54cb7f9dbc14473802d")
	bs_xa, _ := hex.DecodeString("09e41fe2e74667a6676fb0492b496f7d69d45055601ec86839b95e9343407ed592ea357118e5568eea272e9349adf0efbe29307187cfff6e")
	bs_za, _ := hex.DecodeString("b115a615745fc6f453a43d1466e12acd2215ac373cadcd633211235510c6a04c4f041006d07f543f2bd4b050ecdd472be4415ab7a3f79f95")
	in := new(bigNumber).setBytes(bs_in)
	mont := &montgomery{
		new(bigNumber).setBytes(bs_z0),
		new(bigNumber).setBytes(bs_xd),
		new(bigNumber).setBytes(bs_zd),
		new(bigNumber).setBytes(bs_xa),
		new(bigNumber).setBytes(bs_za),
	}

	bs_exp, _ := hex.DecodeString("322d71661943b5e080abed64d9ed331874a975329aaf9b42815e793ac08691e478fe559b29593a5413d5a4475e3ae0735a6d9bc1dc192b7d")
	exp := new(bigNumber).setBytes(bs_exp)

	out, _ := mont.serialize(in)

	c.Assert(out.equals(exp), Equals, true)
}

func (s *Ed448Suite) TestMontgomeryDeserialize(c *C) {
	bs_in, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008")
	in := new(bigNumber).setBytes(bs_in)
	out := new(montgomery)
	out.deserialize(in)
	bs_z0, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040")
	z0 := new(bigNumber).setBytes(bs_z0)
	bs_xd, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	xd := new(bigNumber).setBytes(bs_xd)
	bs_zd, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	zd := new(bigNumber).setBytes(bs_zd)
	bs_xa, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	xa := new(bigNumber).setBytes(bs_xa)
	bs_za, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040")
	za := new(bigNumber).setBytes(bs_za)

	c.Assert(out.z0.equals(z0), Equals, true)
	c.Assert(out.xd.equals(xd), Equals, true)
	c.Assert(out.zd.equals(zd), Equals, true)
	c.Assert(out.xa.equals(xa), Equals, true)
	c.Assert(out.za.equals(za), Equals, true)
}

func (s *Ed448Suite) TestMontgomeryStep(c *C) {
	bs_z0, _ := hex.DecodeString("e281b05e4051a52b331430897d9d950529a46637d3ca1f45e1d2dc4fbd164c956f25dd0cf30458b4129e900faa2ba9b8d305dc4ae1e1b343")
	bs_xd, _ := hex.DecodeString("dc7c2264cf2a3f6178ee7884793f2d0cfe98e602c32adfbec9a5fc225c904f5e1f45c614fc483aec252745e04a38f49e1a4cfc0e8bbf14c5")
	bs_zd, _ := hex.DecodeString("76ad01dbfd7dd72671ad1f827b762fe0c39c808084533b1e22ee18537b7e43c75b995f9e107ec055fbb3df4fb83ad78e69de76a188fb6db6")
	bs_xa, _ := hex.DecodeString("86032c9f990e2680726003f62a1ec5c01f18ad130ce0883b247d2ea9e8d591e6121e6007027d44d94d9659a05fb47e91c3c11b5552cb2185")
	bs_za, _ := hex.DecodeString("5c157ed45b60be2db18a494780b5b7ab79ae1afc3919c9b00c1879495ea079b73990eebf5f0def2897fe8ca78084c07ef89c5bfc336625fd")
	mont := &montgomery{
		new(bigNumber).setBytes(bs_z0),
		new(bigNumber).setBytes(bs_xd),
		new(bigNumber).setBytes(bs_zd),
		new(bigNumber).setBytes(bs_xa),
		new(bigNumber).setBytes(bs_za),
	}
	mont.montgomeryStep()
	bs_z0, _ = hex.DecodeString("e281b05e4051a52b331430897d9d950529a46637d3ca1f45e1d2dc4fbd164c956f25dd0cf30458b4129e900faa2ba9b8d305dc4ae1e1b343")
	bs_xd, _ = hex.DecodeString("c88f896abf42ca2cbff1edf881d1246ee76abe7385932d7b54fb9d71307fdd8043d8a80c7d0363e7a45443d4e9a03bf3e0aab82fb4714c5f")
	bs_zd, _ = hex.DecodeString("962fa8b019eeedd607eda6b44454e17b76b1536f6b336362257d72c3c1576339514f1f4d2d0ae7b0680469a432a2f54cb7f9dbc14473802d")
	bs_xa, _ = hex.DecodeString("09e41fe2e74667a6676fb0492b496f7d69d45055601ec86839b95e9343407ed592ea357118e5568eea272e9349adf0efbe29307187cfff6e")
	bs_za, _ = hex.DecodeString("b115a615745fc6f453a43d1466e12acd2215ac373cadcd633211235510c6a04c4f041006d07f543f2bd4b050ecdd472be4415ab7a3f79f95")
	exp := &montgomery{
		new(bigNumber).setBytes(bs_z0),
		new(bigNumber).setBytes(bs_xd),
		new(bigNumber).setBytes(bs_zd),
		new(bigNumber).setBytes(bs_xa),
		new(bigNumber).setBytes(bs_za),
	}

	c.Assert(mont.z0.equals(exp.z0), Equals, true)
	c.Assert(mont.xd.equals(exp.xd), Equals, true)
	c.Assert(mont.zd.equals(exp.zd), Equals, true)
	c.Assert(mont.xa.equals(exp.xa), Equals, true)
	c.Assert(mont.za.equals(exp.za), Equals, true)
}

func compareNumbers(label string, n *bigNumber, b *big.Int) {
	s := [56]byte{}
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
