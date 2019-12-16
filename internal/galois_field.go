package galoisfield

import (
	"github.com/awnumar/memguard"
)

// GaloisField448 is a field with a finite number of elements. The number depends
// on the word bits: 512/word_bits
// This should not be confunsed with the Field Element
type GaloisField448 struct {
	Limb *memguard.LockedBuffer
}

// NewGaloisField448 will return a newly created, empty field element
// TODO: use no escape
func NewGaloisField448(nlimbs int) *GaloisField448 {
	var gf GaloisField448
	// TODO: extract to constant. It should be 320 for ed25519
	// equivalent to: gf.Limb = memguard.NewBuffer(512 / WordBits)
	// TODO: check this
	gf.Limb = memguard.NewBuffer(nlimbs * 8)

	gf.Limb.Freeze() // Make it inmmutable

	return &gf
}

// NewGaloisField448FromBytes will return a newly created field element with the
// provided byte array
// TODO: use no escape
func NewGaloisField448FromBytes(src []byte) *GaloisField448 {
	var gf GaloisField448

	if len(src) > N32Limbs*8 || len(src) < N64Limbs*8 {
		panic("Wrong Len for buffer initialization")
	}

	gf.Limb = memguard.NewBufferFromBytes(src)

	gf.Limb.Freeze() // Make it inmmutable

	return &gf
}

// limbs will return a uint32 slice for the limb used.
func (gf *GaloisField448) limbs32() []uint32 {
	return gf.Limb.Uint32()
}

// limbs will return a uint64 slice for the limb used.
func (gf *GaloisField448) limbs64() []uint64 {
	return gf.Limb.Uint64()
}

// Destroy securely wipes and frees the underlying memory of the gf.Limb
func (gf *GaloisField448) Destroy() {
	gf.Limb.Destroy()
}

// Copy copies one galoisfield to another.
func (gf *GaloisField448) Copy() *GaloisField448 {
	n := NewGaloisField448(N32Limbs)
	*n = *gf

	return n
}

// AddRaw32 adds one galoisfield to another. For a 32 arch
// TODO: how should be error?
func AddRaw32(x *GaloisField448, y *GaloisField448) *GaloisField448 {
	if x.Limb.Size() != 128 || y.Limb.Size() != 128 {
		return nil
	}

	gf := NewGaloisField448(N32Limbs)

	gf.Limb.Melt()
	defer gf.Limb.Freeze()

	n := gf.limbs32()
	t := x.limbs32()
	z := y.limbs32()

	n[0] = t[0] + z[0]
	n[1] = t[1] + z[1]
	n[2] = t[2] + z[2]
	n[3] = t[3] + z[3]
	n[4] = t[4] + z[4]
	n[5] = t[5] + z[5]
	n[6] = t[6] + z[6]
	n[7] = t[7] + z[7]
	n[8] = t[8] + z[8]
	n[9] = t[9] + z[9]
	n[10] = t[10] + z[10]
	n[11] = t[11] + z[11]
	n[12] = t[12] + z[12]
	n[13] = t[13] + z[13]
	n[14] = t[14] + z[14]
	n[15] = t[15] + z[15]

	return gf
}

// AddRaw64 adds one galoisfield to another. For a 32 arch
func AddRaw64(x *GaloisField448, y *GaloisField448) *GaloisField448 {
	if x.Limb.Size() != 64 || y.Limb.Size() != 64 {
		return nil
	}

	gf := NewGaloisField448(N64Limbs)

	gf.Limb.Melt()
	defer gf.Limb.Freeze()

	n := gf.limbs64()
	t := x.limbs64()
	z := y.limbs64()

	n[0] = t[0] + z[0]
	n[1] = t[1] + z[1]
	n[2] = t[2] + z[2]
	n[3] = t[3] + z[3]
	n[4] = t[4] + z[4]
	n[5] = t[5] + z[5]
	n[6] = t[6] + z[6]
	n[7] = t[7] + z[7]

	return gf
}

// SubRaw32 subtracts one galoisfield to another. For a 32 arch
// TODO: how should be error?
func SubRaw32(x *GaloisField448, y *GaloisField448) *GaloisField448 {
	if x.Limb.Size() != 128 || y.Limb.Size() != 128 {
		return nil
	}

	gf := NewGaloisField448(N32Limbs)

	gf.Limb.Melt()
	defer gf.Limb.Freeze()

	n := gf.limbs32()
	t := x.limbs32()
	z := y.limbs32()

	n[0] = t[0] - z[0]
	n[1] = t[1] - z[1]
	n[2] = t[2] - z[2]
	n[3] = t[3] - z[3]
	n[4] = t[4] - z[4]
	n[5] = t[5] - z[5]
	n[6] = t[6] - z[6]
	n[7] = t[7] - z[7]
	n[8] = t[8] - z[8]
	n[9] = t[9] - z[9]
	n[10] = t[10] - z[10]
	n[11] = t[11] - z[11]
	n[12] = t[12] - z[12]
	n[13] = t[13] - z[13]
	n[14] = t[14] - z[14]
	n[15] = t[15] - z[15]

	return gf
}

// SubRaw64 subtracts one galoisfield to another. For a 32 arch
// TODO: how should be error?
func SubRaw64(x *GaloisField448, y *GaloisField448) *GaloisField448 {
	if x.Limb.Size() != 64 || y.Limb.Size() != 64 {
		return nil
	}

	gf := NewGaloisField448(N64Limbs)

	gf.Limb.Melt()
	defer gf.Limb.Freeze()

	n := gf.limbs64()
	t := x.limbs64()
	z := y.limbs64()

	n[0] = t[0] - z[0]
	n[1] = t[1] - z[1]
	n[2] = t[2] - z[2]
	n[3] = t[3] - z[3]
	n[4] = t[4] - z[4]
	n[5] = t[5] - z[5]
	n[6] = t[6] - z[6]
	n[7] = t[7] - z[7]

	return gf
}

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
