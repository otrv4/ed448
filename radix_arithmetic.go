package ed448

import (
	"fmt"
	"io"
	"math/big"
	"math/rand"
)

const (
	// The size of the Goldilocks field, in bits.
	FieldBits = 448

	// The size of the Goldilocks field, in bytes.
	FieldBytes = (FieldBits + 7) / 8 // 56

	// The size of the Goldilocks scalars, in bits.
	ScalarBits = FieldBits - 2 // 446

	//define FIELD_BYTES          (1+(FIELD_BITS-1)/8)
	//define FIELD_WORDS          (1+(FIELD_BITS-1)/sizeof(word_t))

	BitSize  = ScalarBits
	ByteSize = FieldBytes
)

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

	generateKey(rand io.Reader) (priv []byte, pub []byte, err error)
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

func (c *radixCurve) generateKey(read io.Reader) (priv []byte, pub []byte, err error) {
	priv = make([]byte, ByteSize)

	if _, err = io.ReadFull(read, priv); err != nil {
		return
	}

	//XXX FIXME
	//This is just to not break the API ;)
	//We use the array of random bytes to generate a seed
	seed := new(big.Int).SetBytes(priv)
	seed.Mod(seed, big.NewInt(int64(9223372036854775807)))
	r := rand.New(rand.NewSource(seed.Int64()))

	one := big.NewInt(1)
	m := new(big.Int).Rand(r, new(big.Int).Sub(primeOrder, one)) //[0, primeOrder-1)
	m.Add(m, one)                                                //[1, primeOrder-1]

	priv = m.Bytes()

	fmt.Printf("Private is: %#v\n", m.Bytes())

	//XXX This is sooooooo slow. We need to use an algorithm with some pre-computation
	publicKey := c.multiplyByBase(m.Bytes())
	pub = publicKey.Marshal()
	return
}
