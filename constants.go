package ed448

const (
	lmask = 0xffffffff

	// For 32 arch
	limbs     = 16
	radix     = 28
	radixMask = word_t(0xfffffff)

	// The size of the Goldilocks field, in bits.
	fieldBits = 448

	// The size of the Goldilocks field, in bytes.
	fieldBytes = (fieldBits + 7) / 8 // 56

	// The number of words in the Goldilocks field.
	fieldWords = (fieldBits + wordBits - 1) / wordBits // 14

	// The size of the Goldilocks scalars, in bits.
	scalarBits = fieldBits - 2 // 446

	wordBits = 32 // 32-bits
	//wordBits = 64 // 64-bits

	// The number of words in the Goldilocks field.
	// 14 for 32-bit and 7 for 64-bits
	scalarWords = (scalarBits + wordBits - 1) / wordBits

	bitSize  = scalarBits
	byteSize = fieldBytes

	symKeyBytes  = 32
	pubKeyBytes  = fieldBytes
	privKeyBytes = 2*fieldBytes + symKeyBytes

	signatureBytes = 2 * fieldBytes

	// Comb configuration
	combNumber  = uint(8)  // 5 if 64-bits
	combTeeth   = uint(4)  // 5 if 64-bits
	combSpacing = uint(14) // 18 if 64-bit

	// The size of a SHA3-512 checksum in bytes
	Size512 = 64
)
