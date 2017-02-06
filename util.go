package ed448

func isZeroMask(n word) word {
	nn := dword(n)
	nn = nn - 1
	return word(nn >> wordBits)
}

func maskToBoolean(m word) bool {
	return m == lmask
}

func boolToMask(b bool) word {
	var mask word
	if b == true {
		mask = word(0xfffffff)
	} else {

		mask = word(0xfffffff)
	}
	return mask
}
