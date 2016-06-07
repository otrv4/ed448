package ed448

import (
	"io"
	"math/big"
)

const (
	// The size of the Goldilocks field, in bits.
	FieldBits = 448

	// The size of the Goldilocks field, in bytes.
	FieldBytes = (FieldBits + 7) / 8 // 56

	// The size of the Goldilocks scalars, in bits.
	ScalarBits = FieldBits - 2 // 446

	wordBits = 32 // 32-bits
	//wordBits = 64 // 64-bits

	// The number of words in the Goldilocks field.
	// 14 for 32-bit and 7 for 64-bits
	ScalarWords = (ScalarBits + wordBits - 1) / wordBits

	BitSize  = ScalarBits
	ByteSize = FieldBytes

	//Comb configuration
	CombNumber  = uint(8)  // 5 if 64-bits
	CombTeeth   = uint(4)  // 5 in 64-bits
	CombSpacing = uint(14) // 18 in 64-bit
)

type word_t uint32 //32-bits
//type word_t uint64 //64-bits

type dword_t uint64 //32-bits
//type word_t uint128 //64-bits

//XXX Why having a class at all and not just exported methods?
type radixCurve struct {
	zero, one, two            *bigNumber
	prime, primeOrder, edCons *bigNumber
	basePoint                 Point
}

var rCurve radixCurve

//p = 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff
var primeSerialized = serialized{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

func mustNewPoint(x, y serialized) Point {
	p, err := NewPoint(x, y)
	if err != nil {
		panic("failed to create point")
	}

	return p
}

func init() {
	p, _ := deserialize(primeSerialized)
	rCurve = radixCurve{
		//???
		zero: mustDeserialize(serialized{0x0}),
		one:  mustDeserialize(serialized{0x1}),
		two:  mustDeserialize(serialized{0x02}),

		prime: p,

		//primeOrder: 0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3
		primeOrder: mustDeserialize(serialized{
			0xf3, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
			0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x72, 0xc2,
			0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
			0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
		}),

		//edCons: -39081
		edCons: mustDeserialize(serialized{0xa9, 0x98}), // unsigned

		//gx: 0x297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f
		//gy: 0x13
		basePoint: mustNewPoint(serialized{
			0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
			0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
			0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
			0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
			0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
			0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
			0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
			0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
		},
			serialized{0x13},
		),
	}
}

type pointCurve interface {
	isOnCurve(p Point) bool
	add(p1, p2 Point) (p3 Point)
	double(p1 Point) (p2 Point)
	//multiply(p Point, n *bigNumber) (p2 Point)
	//multiplyByBase(n *bigNumber) (p Point)

	multiplyByBase2(scalar [ScalarWords]word_t) Point
	generateKey(rand io.Reader) (priv []byte, pub []byte, err error)
	computeSecret(private []byte, public []byte) Point
}

func newRadixCurve() pointCurve {
	return &rCurve
}

func (c *radixCurve) isOnCurve(p Point) bool {
	return p.OnCurve()
}

// Returns the sum of (x1,y1) and (x2,y2)
func (c *radixCurve) add(p1, p2 Point) Point {
	return p1.Add(p2)
}

//Returns 2*(x,y)
func (c *radixCurve) double(p Point) Point {
	return p.Double()
}

var (
	primeOrder, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab", 16)
)

func (c *radixCurve) multiply(n []byte, p Point) Point {
	m := new(big.Int).SetBytes(n)
	one := big.NewInt(1)

	out := p

	for i := big.NewInt(0); i.Cmp(m) == -1; i.Add(i, one) {
		out.Add(p)
	}

	return out
}

//multiply2 is Montgomery Ladder Exp
func (c *radixCurve) multiply2(n []byte, p Point) Point {
	m := new(big.Int).SetBytes(n)
	R0 := c.basePoint
	R1 := p
	for pos := m.BitLen() - 2; pos >= 0; pos-- {
		if m.Bit(pos) == 0 {
			R1 = c.add(R0, R1)
			R0 = c.double(R0)
		} else {
			R0 = c.add(R0, R1)
			R1 = c.double(R1)
		}
	}
	return R0
}

func (c *radixCurve) multiplyByBase(n []byte) Point {
	m := new(big.Int).SetBytes(n)
	one := big.NewInt(1)

	priv := c.basePoint
	for i := big.NewInt(0); i.Cmp(m) == -1; i.Add(i, one) {
		//XXX could be optimized by using a readdition formula
		priv = priv.Add(c.basePoint)
	}

	return priv
}

func (c *radixCurve) multiplyByBase2(scalar [ScalarWords]word_t) Point {
	out := &twExtensible{
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
	}

	n := CombNumber
	t := CombTeeth
	s := CombSpacing

	schedule := make([]word_t, ScalarWords)
	scheduleScalarForCombs(schedule, scalar)

	var ni *twNiels

	for i := uint(0); i < s; i++ {
		if i != 0 {
			out = out.double()
		}

		for j := uint(0); j < n; j++ {
			tab := word_t(0)

			for k := uint(0); k < t; k++ {
				bit := (s - 1 - i) + k*s + j*(s*t)
				if bit < ScalarWords*wordBits {
					tab |= (schedule[bit/wordBits] >> (bit % wordBits) & 1) << k
				}
			}

			invert := word_t(tab>>(t-1)) - 1
			tab ^= invert
			tab &= (1 << (t - 1)) - 1

			ni = baseTable.lookup(j, t, uint(tab))
			ni.conditionalNegate(invert != 0)

			if i != 0 || j != 0 {
				out = out.addTwNiels(ni)
			} else {
				out = ni.TwistedExtensible()
			}
		}
	}

	//if(!out.OnCurve()){ return nil } //and maybe panic?

	return out
}

func leBytesToWords(dst []word_t, src []byte) {
	wordBytes := uint(wordBits / 8)
	srcLen := uint(len(src))

	dstLen := uint((srcLen + wordBytes - 1) / wordBytes)
	if dstLen < uint(len(dst)) {
		panic("wrong dst size")
	}

	for i := uint(0); i*wordBytes < srcLen; i++ {
		out := word_t(0)
		for j := uint(0); j < wordBytes && wordBytes*i*j < srcLen; j++ {
			out |= word_t(src[wordBytes*i+j]) << (8 * j)
		}

		dst[i] = out
	}
}

func (c *radixCurve) generateKey(read io.Reader) (priv []byte, pub []byte, err error) {
	buff := make([]byte, FieldBytes)
	if _, err = io.ReadFull(read, buff); err != nil {
		return
	}

	m := new(big.Int)
	randomSeed := new(big.Int).SetBytes(buff)
	_, m = m.DivMod(randomSeed, new(big.Int).Sub(primeOrder, one), m)
	m.Add(m, one) //m E [1, primeOrder-1]

	//XXX this does not always have 56bytes
	privBytes := m.Bytes()
	priv = make([]byte, FieldBytes)
	for i := 0; i < len(privBytes); i++ {
		priv[len(priv)-i-1] = privBytes[len(privBytes)-i-1]
	}

	scalar := [14]word_t{}
	leBytesToWords(scalar[:], priv[:])
	publicKey := c.multiplyByBase2(scalar)
	//XXX Hamburg's code makes a untwist_and_double_and_serialize before
	pub = publicKey.Marshal() //I have no idea how to serialize "twisted extensible coordinates"
	return
}

func (c *radixCurve) computeSecret(private []byte, public []byte) Point {
	scalar := [14]word_t{}
	leBytesToWords(scalar[:], private[:])
	ga := c.multiplyByBase2(scalar)
	gab := c.multiply2(public, ga)
	return gab
}
