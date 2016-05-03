package ed448

import "math/big"

// Edwards curve domain parameters. See https://safecurves.cr.yp.to
var (
	prime     *big.Int // the order of the underlying field
	rho       *big.Int // the order of the base point
	edCons    *big.Int // the constant of the curve equation
	gx, gy    *big.Int // (x,y) of the base point
	fieldSize int      // the size of the underlying field
)

type curve interface {
	isOnCurve(x, y interface{}) bool
	add(x1, y1, x2, y2 interface{}) (x3, y3 interface{})
	double(x1, y1 interface{}) (x3, y3 interface{})
	multiply(x, y interface{}, k []byte) (kx, ky interface{})
	multiplyByBase(k []byte) (kx, ky interface{})
}

type bigintsCurve struct {
}

var ed448 bigintsCurve
var zero, one, two *big.Int

func init() {
	prime, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	rho, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3", 16)
	edCons, _ = new(big.Int).SetString("-39081", 10)
	gx, _ = new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
	gy, _ = new(big.Int).SetString("13", 16)
	fieldSize = 448
}

func init() {
	zero = big.NewInt(0)
	one = big.NewInt(1)
	two = big.NewInt(2)
}

func newBigintsCurve() curve {
	return &ed448
}

func (c *bigintsCurve) isOnCurve(x, y interface{}) bool {
	return isOnCurve(x.(*big.Int), y.(*big.Int))
}

// Reports whether the given (x,y) lies on the bigintsCurve.
func isOnCurve(x, y *big.Int) bool {
	// x² + y² = 1 + bx²y²
	x2 := square(x)
	y2 := square(y)

	x2y2 := mul(x2, y2)
	bx2y2 := mul(edCons, x2y2)

	left := sum(x2, y2)
	left = mod(left)
	right := sum(one, bx2y2)
	right = mod(right)

	return left.Cmp(right) == 0
}

func (c *bigintsCurve) add(x1, y1, x2, y2 interface{}) (x3, y3 interface{}) {
	return add(x1.(*big.Int), y1.(*big.Int), x2.(*big.Int), y2.(*big.Int))
}

// Returns the sum of (x1,y1) and (x2,y2)
func add(x1, y1, x2, y2 *big.Int) (x3, y3 *big.Int) {
	// x² + y² = 1 + bx²y²
	// x3 =  x1y2 + y1x2 / 1 + bx1x2y1y2
	// y3 =  y1y2 - x1x2 / 1 - bx1x2y1y2

	bx1x2y1y2 := mul(edCons, mul(x1, mul(x2, mul(y1, y2))))
	bx1x2y1y2 = mod(bx1x2y1y2)

	x3 = mul(x1, y2)
	x3 = sum(x3, mul(x2, y1))
	x3 = mod(x3)
	divisor := modInv(mod(sum(one, bx1x2y1y2)))
	x3 = mul(x3, divisor)

	y3 = mul(y1, y2)
	y3 = sub(y3, mul(x1, x2))
	y3 = mod(y3)
	divisor = modInv(mod(sub(one, bx1x2y1y2)))
	y3 = mul(y3, divisor)

	return
}

func (c *bigintsCurve) double(x1, y1 interface{}) (x3, y3 interface{}) {
	return double(x1.(*big.Int), y1.(*big.Int))
}

//Returns 2*(x,y)
func double(x1, y1 *big.Int) (x3, y3 *big.Int) {
	// x² + y² = 1 + bx²y²
	// x3 =  2xy / 1 + bx²y² = 2xy / x² + y²
	// y3 =  y² - x² / 1 - bx²y² = y² - x² / 2 - x² - y²

	x2plusy2 := sum(mul(x1, x1), mul(y1, y1))
	x2plusy2 = mod(x2plusy2)

	x3 = mul(x1, y1)
	x3.Lsh(x3, 1) // x3 = 2xy
	x3 = mod(x3)
	divisor := modInv(x2plusy2)
	x3 = mul(x3, divisor) // x3 = 2xy / x² + y²

	y3 = sub(mul(y1, y1), mul(x1, x1)) // y3 = y² - x²
	y3 = mod(y3)
	divisor = modInv(mod(sub(two, x2plusy2)))
	y3 = mul(y3, divisor) // y3 = y² - x² / 2 - x² - y²

	return
}

func (c *bigintsCurve) multiply(x, y interface{}, k []byte) (kx, ky interface{}) {
	return multiply(x.(*big.Int), y.(*big.Int), k)
}

//Performs a scalar multiplication and returns k*(Bx,By) where k is a number in big-endian form.
func multiply(x, y *big.Int, k []byte) (kx, ky *big.Int) {
	kx, ky = x, y
	n := new(big.Int).SetBytes(k)

	for n.Cmp(zero) > 0 {
		if new(big.Int).Mod(n, two).Cmp(zero) == 0 {
			kx, ky = double(kx, ky)
			n = sub(n, two)
		} else {
			kx, ky = add(kx, ky, x, y)
			n = sub(n, one)
		}
	}
	return
}

func (c *bigintsCurve) multiplyByBase(k []byte) (kx, ky interface{}) {
	return multiplyByBase(k)
}

//Returns k*G, where G is the base point of the group and k is an integer in big-endian form.
func multiplyByBase(k []byte) (kx, ky *big.Int) {
	return multiply(gx, gy, k)
}

func sum(x, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

func sub(x, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}

func mul(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(x, y)
}

func square(v *big.Int) *big.Int {
	return new(big.Int).Mul(v, v)
}

func mod(x *big.Int) *big.Int {
	return new(big.Int).Mod(x, prime)
}

func modInv(x *big.Int) *big.Int {
	return new(big.Int).ModInverse(x, prime)
}
