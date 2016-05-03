package ed448

import "math/big"

type curve struct {
	p      *big.Int // the order of the underlying field
	n      *big.Int // the order of the base point
	b      *big.Int // the constant of the curve equation
	gx, gy *big.Int // (x,y) of the base point
	size   int      // the size of the underlying field
}

var ed448 curve
var zero, one, two *big.Int

func init() {
	ed448.p, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	ed448.n, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3", 16)
	ed448.b, _ = new(big.Int).SetString("-39081", 10)
	ed448.gx, _ = new(big.Int).SetString("297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f", 16)
	ed448.gy, _ = new(big.Int).SetString("13", 16)
	ed448.size = 448
}

func init() {
	zero = big.NewInt(0)
	one = big.NewInt(1)
	two = big.NewInt(2)
}

func newEd448() curve {
	return ed448
}

// Reports whether the given (x,y) lies on the curve.
func (c *curve) isOnCurve(x, y *big.Int) bool {
	// x² + y² = 1 + bx²y²
	x2 := square(x)
	y2 := square(y)

	x2y2 := mul(x2, y2)
	bx2y2 := mul(c.b, x2y2)

	left := add(x2, y2)
	left = mod(left)
	right := add(one, bx2y2)
	right = mod(right)

	return left.Cmp(right) == 0
}

// Returns the sum of (x1,y1) and (x2,y2)
func (c *curve) add(x1, y1, x2, y2 *big.Int) (x3, y3 *big.Int) {
	// x² + y² = 1 + bx²y²
	// x3 =  x1y2 + y1x2 / 1 + bx1x2y1y2
	// y3 =  y1y2 - x1x2 / 1 - bx1x2y1y2

	bx1x2y1y2 := mul(c.b, mul(x1, mul(x2, mul(y1, y2))))
	bx1x2y1y2 = mod(bx1x2y1y2)

	x3 = mul(x1, y2)
	x3 = add(x3, mul(x2, y1))
	x3 = mod(x3)
	divisor := modInv(mod(add(one, bx1x2y1y2)))
	x3 = mul(x3, divisor)

	y3 = mul(y1, y2)
	y3 = sub(y3, mul(x1, x2))
	y3 = mod(y3)
	divisor = modInv(mod(sub(one, bx1x2y1y2)))
	y3 = mul(y3, divisor)

	return
}

//Returns 2*(x,y)
func (c *curve) double(x1, y1 *big.Int) (x3, y3 *big.Int) {
	// x² + y² = 1 + bx²y²
	// x3 =  2xy / 1 + bx²y² = 2xy / x² + y²
	// y3 =  y² - x² / 1 - bx²y² = y² - x² / 2 - x² - y²

	x2plusy2 := add(mul(x1, x1), mul(y1, y1))
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

//Performs a scalar multiplication and returns k*(Bx,By) where k is a number in big-endian form.
func (c *curve) multiply(x, y *big.Int, k []byte) (kx, ky *big.Int) {
	kx, ky = x, y
	n := new(big.Int).SetBytes(k)

	for n.Cmp(zero) > 0 {
		if new(big.Int).Mod(n, two).Cmp(zero) == 0 {
			kx, ky = c.double(kx, ky)
			n = sub(n, two)
		} else {
			kx, ky = c.add(kx, ky, x, y)
			n = sub(n, one)
		}
	}
	return
}

//Returns k*G, where G is the base point of the group and k is an integer in big-endian form.
func (c *curve) multiplyByBase(k []byte) (kx, ky *big.Int) {
	kx, ky = c.multiply(c.gx, c.gy, k)
	return
}

func add(x, y *big.Int) *big.Int {
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
	return new(big.Int).Mod(x, newEd448().p)
}

func modInv(x *big.Int) *big.Int {
	return new(big.Int).ModInverse(x, newEd448().p)
}
