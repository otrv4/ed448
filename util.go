package ed448

func isZeroMask(n word) word {
	nn := dword(n)
	nn = nn - 1
	return word(nn >> wordBits)
}

func maskToBoolean(m word) bool {
	return m == lmask
}
