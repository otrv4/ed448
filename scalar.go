package ed448

type ScalarI interface {
	Mul(a, b ScalarI)
	Sub(a, b ScalarI)
}
