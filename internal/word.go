package galoisfield

// The characteristics of the word
const (
	// Word32SizeBits is the size of a word in bits depending on a 32 architecture
	Word32SizeBits = 32
	// Word32SizeBits is the size of a word in bits depending on a 64 architecture
	Word64SizeBits = 64
	// Word32SizeBytes is the size of a word in bytes depending on a 32 architecture
	Word32SizeBytes = 4
	// Word64SizeBytes is the size of a word in bytes depending on a 64 architecture
	Word64SizeBytes = 8
	// N32Limbs is
	N32Limbs = 64 / Word32SizeBytes
	// N64Limbs is
	N64Limbs = 64 / Word64SizeBytes
)

type word32 uint32
type word64 uint64
