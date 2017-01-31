package ed448

const (
	lmask = 0xffffffff

	// For 32 arch
	limbs     = 16
	radix     = 28
	radixMask = word(0xfffffff)

	// The size of the Goldilocks field, in bits.
	fieldBits = 448
	edwardsD  = -39081

	// The size of the Goldilocks field, in bytes.
	fieldBytes = fieldBits / 8 // 56

	// The size of the Goldilocks scalars, in bits.
	scalarBits = fieldBits - 2 // 446
	// The size of the Goldilocks field, in bytes.
	scalarBytes = (scalarBits + 7) / 8 // 56

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
	combNumber  = uint(8)
	combTeeth   = uint(4)
	combSpacing = uint(14)

	// DecafComb configuration
	decafCombNumber  = uint(5)
	decafCombTeeth   = uint(5)
	decafCombSpacing = uint(18)

	uintZero = uint(0)

	montgomeryFactor = word(0xae918bc5)
)

var (
	bigZero = &bigNumber{0}
	bigOne  = &bigNumber{1}
	bigTwo  = &bigNumber{2}

	sqrtDminus1 = mustDeserialize(serialized{
		0x46, 0x9f, 0x74, 0x36, 0x18, 0xe2, 0xd2, 0x79,
		0x01, 0x4f, 0x2b, 0xb4, 0x8d, 0x88, 0x38, 0xea,
		0xde, 0xab, 0x9a, 0x18, 0x5a, 0x06, 0x4c, 0xf1,
		0xa6, 0x5c, 0xe6, 0x51, 0x70, 0x97, 0x4d, 0x42,
		0x7b, 0x9f, 0xa4, 0x56, 0xf6, 0xc5, 0x28, 0x46,
		0xac, 0xdc, 0x4a, 0x73, 0x48, 0x87, 0x3b, 0x44,
		0x49, 0x7a, 0x5b, 0xb2, 0xc0, 0xc0, 0xfe, 0x12,
	})

	scalarQ = &decafScalar{
		0xab5844f3, 0x2378c292,
		0x8dc58f55, 0x216cc272,
		0xaed63690, 0xc44edb49,
		0x7cca23e9, 0xffffffff,
		0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff,
		0xffffffff, 0x3fffffff,
	}

	scalarR2 = &decafScalar{
		0x049b9b60, 0xe3539257,
		0xc1b195d9, 0x7af32c4b,
		0x88ea1859, 0x0d66de23,
		0x5ee4d838, 0xae17cf72,
		0xa3c47c44, 0x1a9cc14b,
		0xe4d070af, 0x2052bcb7,
		0xf823b729, 0x3402a939,
	}
)
