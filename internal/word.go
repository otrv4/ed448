package galoisfield

// The characteristics of the word
const (
	// WordSizeBits is the size of a word in bits depending on the architecture
	// TODO: for the moment using a 32bits arch
	WordSizeBits = 32
	// WordSizeBytes is the size of a word in bytes depending on the architecture
	WordSizeBytes = 4
	// NLimbs is
	NLimbs = 64 / WordSizeBytes
)

type word uint32
