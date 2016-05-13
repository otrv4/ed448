package ed448

const (
	radix = 56
	limbs = 8
)

type word uint64
type bigNumber [limbs]word

func isZero(n int64) int64 {
	return ^n
}

func constantTimeGreaterOrEqualP(n [limbs]int64) bool {
	var (
		ge   = int64(-1)
		mask = int64(1)<<radix - 1
	)

	for i := 0; i < 4; i++ {
		ge &= n[i]
	}

	ge = (ge & (n[4] + 1)) | isZero(n[4]^mask)

	for i := 5; i < 8; i++ {
		ge &= n[i]
	}

	return ge == mask
}

func serialize(dst []byte, src bigNumber) {
	const (
		rows    = limbs
		columns = radix / limbs
	)

	var n bigNumber
	copy(n[:], src[:])

	for i := uint(0); i < rows; i++ {
		for j := uint(0); j < columns; j++ {
			dst[columns*i+j] = byte(n[i])
			n[i] >>= 8
		}
	}
}
