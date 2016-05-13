package ed448

type word uint64

func isZero(n int64) int64 {
	return ^n
}

func constantTimeGreaterOrEqualP(n [8]int64) bool {
	var (
		ge   = int64(-1)
		mask = int64(1)<<56 - 1
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
