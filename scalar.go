package ed448

type Scalar interface {
	Mul(a, b Scalar)
	Sub(a, b Scalar)
}
