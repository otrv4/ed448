package ed448

import "math/big"

//XXX Why having a class at all and not just exported methods?
type radixCurve struct {
	zero, one, two             bigNumber
	prime, rho, edCons, gx, gy bigNumber
}

var rCurve radixCurve

func init() {
	rCurve = radixCurve{
		//???
		zero: bigNumber{},
		one:  bigNumber{},
		two:  bigNumber{},

		//XXX This is in 64-bits, need to use serialize
		prime: bigNumber{
			0xffffffffffffff,
			0xffffffffffffff,
			0xffffffffffffff,
			0xffffffffffffff,
			0xfffffffffffffe,
			0xffffffffffffff,
			0xffffffffffffff,
			0xffffffffffffff,
		},

		//XXX This is in 64-bits, need to use serialize
		rho: bigNumber{
			0x78c292ab5844f3,
			0xc2728dc58f5523,
			0x49aed63690216c,
			0x7cca23e9c44edb,
			0xffffffffffffff,
			0xffffffffffffff,
			0xffffffffffffff,
			0x3fffffffffffff,
		},

		edCons: bigNumber{},
		gx:     bigNumber{},
		gy:     bigNumber{},
	}

	//prime, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	//rho, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3", 16)
	//edCons, _ = new(big.Int).SetString("-39081", 10)
	//gx, _ = new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
	//gy, _ = new(big.Int).SetString("13", 16)
	//fieldSize = 448
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
