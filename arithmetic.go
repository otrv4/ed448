package ed448

const (
	Limbs = 8
	Radix = 56
)

type word uint64
type limb word
type bigNumber [Limbs]limb
type serialized [Radix]byte

func deserialize(in serialized) (n bigNumber, ok bool) {
	const (
		columns = Limbs
		rows    = Limbs - 1
	)

	for i := uint(0); i < columns; i++ {
		for j := uint(0); j < rows; j++ {
			n[i] |= limb(in[rows*i+j]) << (columns * j)
		}
	}

	ok = !constantTimeGreaterOrEqualP(n)

	return
}

//TODO: Make this work with a word parameter
func isZero(n int64) int64 {
	return ^n
}

func constantTimeGreaterOrEqualP(n bigNumber) bool {
	var (
		ge   = int64(-1)
		mask = int64(1)<<Radix - 1
	)

	for i := 0; i < 4; i++ {
		ge &= int64(n[i])
	}

	ge = (ge & (int64(n[4]) + 1)) | isZero(int64(n[4])^mask)

	for i := 5; i < 8; i++ {
		ge &= int64(n[i])
	}

	return ge == mask
}

func serialize(dst []byte, src bigNumber) {
	const (
		rows    = Limbs
		columns = Radix / Limbs
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
