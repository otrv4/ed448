package ed448

// Scalar is a interface of Ed448 scalar
type Scalar interface {
	Mul(a, b Scalar)
	Sub(a, b Scalar)
	Add(a, b Scalar)
	Decode(src []byte) error
	Encode(dst []byte)
	Copy() Scalar
	Equals(a Scalar) bool
	// unexposed funcs
	halve(a, b Scalar)
}
