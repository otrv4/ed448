package ed448

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type Ed448Suite struct{}

var _ = Suite(&Ed448Suite{})

func (s *Ed448Suite) TestGenerateKeysProducesKeyPair(c *C) {
	curve := NewCurve()
	priv, pub, ok := curve.GenerateKeys()
	c.Assert(ok, Equals, true)
	c.Assert(priv, NotNil)
	c.Assert(pub, NotNil)
}

func (s *Ed448Suite) TestSignAndVerify(c *C) {
	curve := NewCurve()
	priv, pub, ok := curve.GenerateKeys()
	c.Assert(ok, Equals, true)

	message := []byte("sign here.")

	signature, ok := curve.Sign(priv, message)

	c.Assert(ok, Equals, true)
	c.Assert(signature, NotNil)

	valid := curve.Verify(signature, message, pub)

	c.Assert(valid, Equals, true)
}

func (s *Ed448Suite) TestComputeSecret(c *C) {
	curve := NewCurve()
	privA, pubA, _ := curve.GenerateKeys()
	privB, pubB, _ := curve.GenerateKeys()
	secretA := curve.ComputeSecret(privA, pubB)
	secretB := curve.ComputeSecret(privB, pubA)
	c.Assert(secretA, DeepEquals, secretB)
}

func (s *Ed448Suite) TestDeserializeModQ(c *C) {
	serial := []byte{
		0xb3, 0xa4, 0x53, 0x31, 0xb1, 0x2b, 0x41, 0x1a,
		0xda, 0x51, 0xcf, 0xba, 0x0d, 0xea, 0x65, 0xb3,
		0x1b, 0x97, 0x9b, 0x41, 0xfe, 0x18, 0x93, 0x0c,
		0x6e, 0x4c, 0x02, 0x8a, 0x26, 0x24, 0xdf, 0xf0,
		0x24, 0x24, 0x06, 0x01, 0x4a, 0xb6, 0x3c, 0xab,
		0x33, 0x1e, 0xb5, 0xcf, 0x79, 0xc2, 0xc2, 0x6b,
		0xbb, 0x5e, 0xf8, 0xd8, 0x3e, 0x74, 0x26, 0x2c,
	}

	exp := BigNumber{
		827565235, 440478641, 3134149082, 3009800717,
		1100715803, 210966782, 2315406446, 4041155622,
		17179684, 2872882762, 3484753459, 1807925881,
		3640155835, 740717630, 0, 0,
	}

	output := DeserializeModQ(serial)
	c.Assert(output, DeepEquals, exp)
}
