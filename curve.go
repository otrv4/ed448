package ed448

import (
	"crypto/sha512"
	"errors"
	"io"
	"math/big"
)

const (
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

	//Comb configuration
	combNumber  = uint(8)  // 5 if 64-bits
	combTeeth   = uint(4)  // 5 if 64-bits
	combSpacing = uint(14) // 18 if 64-bit
)

type word_t uint32 //32-bits
//type word_t uint64 //64-bits

type dword_t uint64 //32-bits
//type word_t uint128 //64-bits

//XXX Why having a class at all and not just exported methods?
type radixCurve struct {
	zero, one, two            *bigNumber
	prime, primeOrder, edCons *bigNumber
	basePoint                 *homogeneousProjective
}

var rCurve radixCurve

//p = 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff
var primeSerialized = serialized{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

func mustNewPoint(x, y serialized) *homogeneousProjective {
	p, err := NewPoint(x, y)
	if err != nil {
		panic("failed to create point")
	}

	return p
}

func init() {
	p, _ := deserialize(primeSerialized)
	rCurve = radixCurve{
		//???
		zero: mustDeserialize(serialized{0x0}),
		one:  mustDeserialize(serialized{0x1}),
		two:  mustDeserialize(serialized{0x02}),

		prime: p,

		//primeOrder: 0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3
		primeOrder: mustDeserialize(serialized{
			0xf3, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
			0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x72, 0xc2,
			0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
			0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
		}),

		//edCons: -39081
		edCons: mustDeserialize(serialized{0xa9, 0x98}), // unsigned

		//gx: 0x297ea0ea2692ff1b4faff46098453a6a26adf733245f065c3c59d0709cecfa96147eaaf3932d94c63d96c170033f4ba0c7f0de840aed939f
		//gy: 0x13
		basePoint: mustNewPoint(serialized{
			0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
			0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
			0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
			0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
			0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
			0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
			0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
			0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
		},
			serialized{0x13},
		),
	}
}

//XXX We dont need an interface anymore
type pointCurve interface {
	BasePoint() *homogeneousProjective

	multiplyByBase(scalar [scalarWords]word_t) *twExtensible
	generateKey(rand io.Reader) (k privateKey, err error)
	sign(msg []byte, k *privateKey) ([signatureBytes]byte, error)
	verify(signature [signatureBytes]byte, msg []byte, k *publicKey) bool
	computeSecret(private []byte, public []byte) []byte
}

func newRadixCurve() pointCurve {
	return &rCurve
}

func (c *radixCurve) BasePoint() *homogeneousProjective {
	return c.basePoint
}

var (
	primeOrder, _ = new(big.Int).SetString("3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab", 16)
)

func (c *radixCurve) multiplyMontgomery(in *bigNumber, scalar [fieldWords]word_t, nbits, n_extra_doubles int) (*bigNumber, word_t) {
	mont := new(montgomery)
	mont.deserialize(in)
	var i, j, n int
	n = (nbits - 1) % wordBits
	pflip := word_t(0)

	for j = (nbits+wordBits-1)/wordBits - 1; j >= 0; j-- {
		w := scalar[j]
		for i = n; i >= 0; i-- {
			flip := -((w >> uint(i)) & 1)

			swap := flip ^ pflip
			mont.xa.conditionalSwap(mont.xd, swap)
			mont.za.conditionalSwap(mont.zd, swap)
			mont.montgomeryStep()
			pflip = flip
		}
		n = wordBits - 1
	}

	mont.xa.conditionalSwap(mont.xd, pflip)
	mont.za.conditionalSwap(mont.zd, pflip)
	//assert(n_extra_doubles < INT_MAX);
	for j = 0; j < n_extra_doubles; j++ {
		mont.montgomeryStep()
	}

	out, ok := mont.serialize(in)
	return out, word_t(ok)
}

func (c *radixCurve) multiplyByBase(scalar [scalarWords]word_t) *twExtensible {
	out := &twExtensible{
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
		new(bigNumber),
	}

	n := combNumber
	t := combTeeth
	s := combSpacing

	schedule := make([]word_t, scalarWords)
	scheduleScalarForCombs(schedule, scalar)

	var ni *twNiels

	for i := uint(0); i < s; i++ {
		if i != 0 {
			out = out.double()
		}

		for j := uint(0); j < n; j++ {
			tab := word_t(0)

			for k := uint(0); k < t; k++ {
				bit := (s - 1 - i) + k*s + j*(s*t)
				if bit < scalarWords*wordBits {
					tab |= (schedule[bit/wordBits] >> (bit % wordBits) & 1) << k
				}
			}

			invert := word_t(tab>>(t-1)) - 1
			tab ^= invert
			tab &= (1 << (t - 1)) - 1

			ni = baseTable.lookup(j, t, uint(tab))
			ni.conditionalNegate(invert)

			if (i | j) != 0 {
				out = out.addTwNiels(ni)
			} else {
				out = ni.TwistedExtensible()
			}
		}
	}

	//if(!out.OnCurve()){ return nil } //and maybe panic?

	return out
}

// Deserializes an array of bytes (little-endian) into an array of words
func bytesToWords(dst []word_t, src []byte) {
	wordBytes := uint(wordBits / 8)
	srcLen := uint(len(src))

	dstLen := uint((srcLen + wordBytes - 1) / wordBytes)
	if dstLen < uint(len(dst)) {
		panic("wrong dst size")
	}

	for i := uint(0); i*wordBytes < srcLen; i++ {
		out := word_t(0)
		for j := uint(0); j < wordBytes && wordBytes*i+j < srcLen; j++ {
			out |= word_t(src[wordBytes*i+j]) << (8 * j)
		}

		dst[i] = out
	}
}

// Serializes an array of words into an array of bytes (little-endian)
func wordsToBytes(dst []byte, src []word_t) {
	wordBytes := wordBits / 8

	for i := 0; i*wordBytes < len(dst); i++ {
		for j := 0; j < wordBytes; j++ {
			b := src[i] >> uint(8*j)
			dst[wordBytes*i+j] = byte(b)
		}
	}
}

//See Goldilocks spec, "Public and private keys" section.
//This is equivalent to PRF(k)
func pseudoRandomFunction(k [symKeyBytes]byte) []byte {
	h := sha512.New()
	h.Write([]byte("derivepk"))
	h.Write(k[:])
	return h.Sum(nil)
}

//See Goldilocks spec, "Public and private keys" section.
//This is equivalent to DESERMODq()
func deserializeModQ(dst []word_t, serial []byte) {
	barrettDeserializeAndReduce(dst, serial, &curvePrimeOrder)
}

func generateSymmetricKey(read io.Reader) (symKey [symKeyBytes]byte, err error) {
	_, err = io.ReadFull(read, symKey[:])
	return
}

func (c *radixCurve) derivePrivateKey(symmetricKey [symKeyBytes]byte) (privateKey, error) {
	k := privateKey{}
	copy(k.symKey(), symmetricKey[:])

	skb := pseudoRandomFunction(symmetricKey)
	secretKey := [fieldWords]word_t{}
	deserializeModQ(secretKey[:], skb)
	wordsToBytes(k.secretKey(), secretKey[:])

	publicKey := c.multiplyByBase(secretKey)
	serializedPublicKey := publicKey.untwistAndDoubleAndSerialize()
	serialize(k.publicKey(), serializedPublicKey)

	return k, nil
}

func (c *radixCurve) generateKey(read io.Reader) (k privateKey, err error) {
	symKey, err := generateSymmetricKey(read)
	if err != nil {
		return
	}

	return c.derivePrivateKey(symKey)
}

func (c *radixCurve) computeSecret(private, public []byte) []byte {
	var sk [fieldWords]word_t
	var pub serialized
	copy(pub[:], public)

	msucc := word_t(0xffffffff)
	pk, succ := deserializeReturnMask(pub)

	msucc &= barrettDeserializeReturnMask(sk[:], private, &curvePrimeOrder)

	ok := word_t(0)
	pk, ok = c.multiplyMontgomery(pk, sk, scalarBits, 1)
	succ &= ok

	gxy := make([]byte, fieldBytes)
	serialize(gxy, pk)

	//XXX SECURITY should we wipe the temporary variables?

	//XXX add error conditions based on succ and msucc
	return gxy
}

func (c *radixCurve) sign(msg []byte, k *privateKey) (s [signatureBytes]byte, e error) {
	secretKeyWords := [fieldWords]word_t{}
	if ok := barrettDeserialize(secretKeyWords[:], k.secretKey(), &curvePrimeOrder); !ok {
		//XXX SECURITY should we wipe secretKeyWords?
		e = errors.New("corrupted private key")
		return
	}

	nonce := [fieldWords]word_t{}
	deriveNonce(nonce[:], msg, k.symKey())

	// tmpSig = 4 * nonce * basePoint
	tmpSig := [fieldBytes]byte{}
	gsk := c.multiplyByBase(nonce).double().untwistAndDoubleAndSerialize()
	serialize(tmpSig[:], gsk)

	challenge := [fieldWords]word_t{}
	deriveChallenge(challenge[:], k.publicKey(), tmpSig, msg)

	barrettNegate(challenge[:], &curvePrimeOrder)
	barrettMac(nonce[:], challenge[:], secretKeyWords[:], &curvePrimeOrder)

	carry := addExtPacked(nonce[:], nonce[:], nonce[:], 0xffffffff)
	barrettReduce(nonce[:], carry, &curvePrimeOrder)

	copy(s[:], tmpSig[:fieldBytes])
	wordsToBytes(s[fieldBytes:], nonce[:])

	//XXX SECURITY Should we wipe nonce, gsk, secretKeyWords, tmpSig, challenge?

	/* response = 2(nonce_secret - sk*challenge)
	 * Nonce = 8[nonce_secret]*G
	 * PK = 2[sk]*G, except doubled (TODO)
	 * so [2] ( [response]G + 2[challenge]PK ) = Nonce
	 */

	return
}

//XXX Should pubKey have a fixed size here?
func deriveChallenge(dst []word_t, pubKey []byte, tmpSignature [fieldBytes]byte, msg []byte) {
	h := sha512.New()
	h.Write(pubKey)
	h.Write(tmpSignature[:])
	h.Write(msg)

	barrettDeserializeAndReduce(dst, h.Sum(nil), &curvePrimeOrder)
}

func deriveNonce(dst []word_t, msg []byte, symKey []byte) {
	h := sha512.New()
	h.Write([]byte("signonce"))
	h.Write(symKey)
	h.Write(msg)
	h.Write(symKey)

	barrettDeserializeAndReduce(dst, h.Sum(nil), &curvePrimeOrder)

	//XXX SECURITY should we wipe r?
}

func (c *radixCurve) verify(signature [signatureBytes]byte, msg []byte, k *publicKey) bool {
	//TODO
	return false
}
