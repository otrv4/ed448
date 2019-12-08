package galoisfield

import (
	"github.com/awnumar/memguard"
)

const (
	// WordBits is
	// TODO: for the moment using a 32bits arch
	WordBits = 32
	// WordBytes is
	WordBytes = 4
	// NLimbs is
	NLimbs = 64 / WordBytes
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
	// equivalent to: gf.Limb = memguard.NewBuffer(512 / WordBits)
	gf.Limb = memguard.NewBuffer(NLimbs)

	gf.Limb.Freeze() // Make it inmmutable

	return &gf
}

// Destroy securely wipes and frees the underlying memory of the gf.Limb
func (gf *GaloisField448) Destroy() {
	gf.Limb.Destroy()
}
