package galoisfield

import "github.com/awnumar/memguard"

const (
	// WordBits is
	// TODO: for the moment
	WordBits = 32
)

// GaloisField448 is a field with a finite number of elements. The number depends
// on the word bits: 512/word_bits
// This should not be confunsed with the Field Element
type GaloisField448 struct {
	Limb *memguard.LockedBuffer
}

// NewGaloisField448 will return a newly created, empty field element
// TODO: use no escape
func NewGaloisField448() *GaloisField448 {
	var gf GaloisField448
	// TODO: extract to constant. It should be 320 for ed25519
	gf.Limb = memguard.NewBuffer(512 / WordBits)

	gf.Limb.Freeze() // Make it inmmutable

	return &gf
}
