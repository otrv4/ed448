package ed448

// Scalar is a interface of Ed448 scalar
type Scalar interface {
	Equals(a Scalar) bool
	Copy() Scalar
	Add(a, b Scalar)
	Sub(a, b Scalar)
	Mul(a, b Scalar)
	Encode() []byte
	Decode(src []byte) error
	// unexposed funcs
	halve(a, b Scalar)
}
