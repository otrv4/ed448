package ed448

import "math/big"

//XXX Why having a class at all and not just exported methods?
type radixCurve struct {
	zero, one, two             bigNumber
	prime, rho, edCons, gx, gy bigNumber
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

func init() {
	p, _ := deserialize(primeSerialized)
	rCurve = radixCurve{
		//???
		zero: bigNumber{},
		one:  bigNumber{},
		two:  bigNumber{},

		prime: p,

		//rho: 0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3
		rho: mustDeserialize(serialized{
			0xf3, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
			0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x72, 0xc2,
			0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
			0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
		}),

		//XXX not sure if we should represent this as radix-base because it is negative
		//edCons: -39081
		edCons: bigNumber{},

		//gx: 0x297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f
		gx: mustDeserialize(serialized{
			0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
			0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
			0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
			0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
			0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
			0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
			0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
			0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
		}),

		//gy: 0x13
		gy: bigNumber{0x13},
	}
}

//func init() {
//	zero = big.NewInt(0)
//	one = big.NewInt(1)
//	two = big.NewInt(2)
//}

func newRadixCurve() curve {
	return &rCurve
}

func (c *radixCurve) isOnCurve(x, y interface{}) bool {
	// x² + y² = 1 + bx²y²
	x2 := squareBigints(x.(*big.Int))
	y2 := squareBigints(y.(*big.Int))

	x2y2 := mulBigints(x2, y2)
	bx2y2 := mulBigints(edCons, x2y2)

	left := sumBigints(x2, y2)
	left = modBigints(left)
	right := sumBigints(one, bx2y2)
	right = modBigints(right)

	return left.Cmp(right) == 0
}

// Returns the sum of (x1,y1) and (x2,y2)
func (c *radixCurve) add(x1, y1, x2, y2 interface{}) (x3, y3 interface{}) {
	// x² + y² = 1 + bx²y²
	// x3 =  x1y2 + y1x2 / 1 + bx1x2y1y2
	// y3 =  y1y2 - x1x2 / 1 - bx1x2y1y2

	x1x2 := mulBigints(x1.(*big.Int), x2.(*big.Int))
	y1y2 := mulBigints(y1.(*big.Int), y2.(*big.Int))
	bx1x2y1y2 := mulBigints(edCons, mulBigints(x1x2, y1y2))
	bx1x2y1y2 = modBigints(bx1x2y1y2)

	x1y2 := mulBigints(x1.(*big.Int), y2.(*big.Int))
	x2y1 := mulBigints(x2.(*big.Int), y1.(*big.Int))
	x3 = sumBigints(x1y2, x2y1)
	x3 = modBigints(x3.(*big.Int))
	divisor := modInvBigints(modBigints(sumBigints(one, bx1x2y1y2)))
	x3 = mulBigints(x3.(*big.Int), divisor)

	y3 = subBigints(y1y2, x1x2)
	y3 = modBigints(y3.(*big.Int))
	divisor = modInvBigints(modBigints(subBigints(one, bx1x2y1y2)))
	y3 = mulBigints(y3.(*big.Int), divisor)

	return x3.(*big.Int), y3.(*big.Int)
}

//Returns 2*(x,y)
func (c *radixCurve) double(x, y interface{}) (x2, y2 interface{}) {
	// x² + y² = 1 + bx²y²
	// x3 =  2xy / 1 + bx²y² = 2xy / x² + y²
	// y3 =  y² - x² / 1 - bx²y² = y² - x² / 2 - x² - y²

	x1, y1 := x.(*big.Int), y.(*big.Int)
	x2plusy2 := sumBigints(mulBigints(x1, x1), mulBigints(y1, y1))
	x2plusy2 = modBigints(x2plusy2)

	x3 := mulBigints(x1, y1)
	x3.Lsh(x3, 1) // x3 = 2xy
	x3 = modBigints(x3)
	divisor := modInvBigints(x2plusy2)
	x3 = mulBigints(x3, divisor) // x3 = 2xy / x² + y²

	y3 := subBigints(mulBigints(y1, y1), mulBigints(x1, x1)) // y3 = y² - x²
	y3 = modBigints(y3)
	divisor = modInvBigints(modBigints(subBigints(two, x2plusy2)))
	y3 = mulBigints(y3, divisor) // y3 = y² - x² / 2 - x² - y²

	x2, y2 = x3, y3
	return
}

//Performs a scalar multiplication and returns k*(Bx,By) where k is a number in big-endian form.
func (c *radixCurve) multiply(x, y interface{}, k []byte) (kx, ky interface{}) {
	kx, ky = x.(*big.Int), y.(*big.Int)
	n := new(big.Int).SetBytes(k)

	for n.Cmp(zero) > 0 {
		if new(big.Int).Mod(n, two).Cmp(zero) == 0 {
			kx, ky = c.double(kx, ky)
			n = subBigints(n, two)
		} else {
			kx, ky = c.add(kx, ky, kx, ky)
			n = subBigints(n, one)
		}
	}
	return
}

//Returns k*G, where G is the base point of the group and k is an integer in big-endian form.
func (c *radixCurve) multiplyByBase(k []byte) (kx, ky interface{}) {
	return c.multiply(gx, gy, k)
}
