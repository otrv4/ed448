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

// Copy copies one galoisfield to another.
func (gf *GaloisField448) Copy() *GaloisField448 {
	n := NewGaloisField448()
	*n = *gf

	return n
}

//static INLINE_UNUSED void gf_add_RAW (gf out, const gf a, const gf b);
//static INLINE_UNUSED void gf_sub_RAW (gf out, const gf a, const gf b);
//static INLINE_UNUSED void gf_bias (gf inout, int amount);
//static INLINE_UNUSED void gf_weak_reduce (gf inout);
//
//void gf_strong_reduce (gf inout);
//void gf_add (gf out, const gf a, const gf b);
//void gf_sub (gf out, const gf a, const gf b);
//void gf_mul (gf_s *__restrict__ out, const gf a, const gf b);
//void gf_mulw_unsigned (gf_s *__restrict__ out, const gf a, uint32_t b);
//void gf_sqr (gf_s *__restrict__ out, const gf a);
//mask_t gf_isr(gf a, const gf x); /** a^2 x = 1, QNR, or 0 if x=0.  Return true if successful */
//mask_t gf_eq (const gf x, const gf y);
//mask_t gf_lobit (const gf x);
//
//void gf_serialize (uint8_t *serial, const gf x);
//mask_t gf_deserialize (gf x, const uint8_t serial[SER_BYTES],uint8_t hi_nmask);
