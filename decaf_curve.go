package ed448

import (
	"errors"
	"io"

	"golang.org/x/crypto/sha3"
)

func decafPseudoRandomFunction(sym []byte) []byte {
	hash := sha3.NewShake256()
	hash.Write(sym[:])
	hash.Write([]byte("decaf_448_derive_private_key"))
	var out [64]byte
	hash.Read(out[:])
	return out[:]
}

func (c *decafCurveT) decafDerivePrivateKey(sym [symKeyBytes]byte) (*privateKey, error) {
	k := &privateKey{}
	copy(k.symKey(), sym[:])

	skb := decafPseudoRandomFunction(sym[:])
	secretKey := &scalar{}

	barrettDeserializeAndReduce(secretKey[:], skb, &curvePrimeOrder)

	err := secretKey.encode(k.secretKey())
	if err != nil {
		return nil, err
	}

	publicKey := precomputedScalarMul(secretKey)
	publicKey.decafEncode(k.publicKey())

	return k, nil
}

func (c *decafCurveT) decafGenerateKeys(r io.Reader) (k *privateKey, err error) {
	symKey, err := generateSymmetricKey(r)
	if err != nil {
		return
	}

	return c.decafDerivePrivateKey(symKey)
}

func (c *decafCurveT) decafComputeSecret(myPriv *privateKey, yourPub [fieldBytes]byte) ([]byte, word) {
	var delta, less uint16
	var ser [fieldBytes]byte
	var sk scalar
	invalid := "decaf_448_ss_invalid"

	priv := myPriv.secretKey()
	pub := myPriv.publicKey()
	sym := myPriv.symKey()

	_ = barrettDeserializeReturnMask(sk[:], priv, &curvePrimeOrder)

	// Lexsort keys. Less will be -1 if mine is less, and 0 otherwise.
	for i := uintZero; i < fieldBytes; i++ {
		delta = uint16(pub[i])
		delta -= uint16(yourPub[i])
		// Case:
		// = -> delta = 0 -> hi delta-1 = -1, hi delta = 0
		// > -> delta > 0 -> hi delta-1 = 0, hi delta = 0
		// < -> delta < 0 -> hi delta-1 = (doesnt matter), hi delta = -1
		less &= delta - 1
		less |= delta
	}
	less >>= 8

	// update the lesser
	for j := uintZero; j < fieldBytes; j++ {
		ser[j] = uint8(uint16(pub[j])&less) | uint8(uint16(yourPub[j])&^less) // check
	}

	hash := sha3.NewShake256()
	hash.Write(ser[:])

	// update the greater
	for k := uintZero; k < fieldBytes; k++ {
		ser[k] = uint8(uint16(pub[k])&^less | uint16(yourPub[k])&less)
	}

	hash.Write(ser[:])
	ser, ok := directPointScalarMul(yourPub, &sk, decafFalse)

	// If invalid, replace
	for l := uintZero; l < fieldBytes; l++ {
		ser[l] &= uint8(ok)
		if l < wordBits {
			ser[l] |= sym[l] & ^uint8(ok)
		} else if l-wordBits < uint(len(invalid)) {
			ser[l] |= invalid[l-wordBits] & ^uint8(ok)
		}
	}

	hash.Write(ser[:])
	var shared [fieldBytes]byte
	hash.Read(shared[:])

	//TODO: should we wipe ser bytes?
	return shared[:], ok
}

func decafDeriveNonce(msg []byte, symKey []byte) *scalar {
	h := sha3.NewShake256()
	h.Write(msg)
	h.Write(symKey)
	h.Write([]byte("decaf_448_sign_shake"))
	var out [64]byte
	h.Read(out[:])

	dst := &scalar{}

	barrettDeserializeAndReduce(dst[:], out[:], &curvePrimeOrder)

	return dst
}

func decafDeriveChallenge(pubKey []byte, tmpSignature [fieldBytes]byte, msg []byte) *scalar {
	h := sha3.NewShake256()
	h.Write(msg)
	h.Write(pubKey)
	h.Write(tmpSignature[:])
	var out [64]byte
	h.Read(out[:])

	dst := &scalar{}

	barrettDeserializeAndReduce(dst[:], out[:], &curvePrimeOrder)

	return dst
}

func (c *decafCurveT) decafDeriveTemporarySignature(nonce *scalar) (dst [fieldBytes]byte) {
	point := precomputedScalarMul(nonce)
	point.decafEncode(dst[:])
	return
}

func (c *decafCurveT) decafSign(msg []byte, k *privateKey) (sig [signatureBytes]byte, err error) {
	secretKeyWords := &scalar{}
	//TODO: should secret words be destroyed?

	if ok := barrettDeserialize(secretKeyWords[:], k.secretKey(), &curvePrimeOrder); !ok {
		err = errors.New("Corrupted private key")
		return
	}

	nonce := decafDeriveNonce(msg, k.symKey())
	tmpSignature := c.decafDeriveTemporarySignature(nonce)
	challenge := decafDeriveChallenge(k.publicKey(), tmpSignature, msg)

	challenge.mul(challenge, secretKeyWords)
	nonce.sub(nonce, challenge)

	copy(sig[:fieldBytes], tmpSignature[:])
	nonce.encode(sig[fieldBytes:])

	//TODO: should nonce and challenge be destroyed?
	return
}

func (c *decafCurveT) decafVerify(signature [signatureBytes]byte, msg []byte, k *publicKey) (bool, error) {
	serPubkey := serialized(*k)
	tmpSig := [fieldBytes]byte{}
	copy(tmpSig[:], signature[:])
	challenge := decafDeriveChallenge(serPubkey[:], tmpSig, msg)

	point := &twExtendedPoint{
		x: &bigNumber{},
		y: &bigNumber{},
		z: &bigNumber{},
		t: &bigNumber{},
	}
	pkPoint := &twExtendedPoint{
		x: &bigNumber{},
		y: &bigNumber{},
		z: &bigNumber{},
		t: &bigNumber{},
	}

	ret, err := decafDecodeOld(point, tmpSig, true)
	if err != nil {
		return false, err
	}

	ret1, err := decafDecodeOld(pkPoint, serPubkey, false)
	if err != nil {
		return false, err
	}
	// TODO: hacky. FIX ME.
	ret &= ret1

	response := &scalar{}
	ret &= response.decode(signature[56:])

	pkPoint = decafDoubleNonSecretScalarMul(pkPoint, response, challenge)
	ret &= pkPoint.equals(point)

	if ret != word(lmask) {
		return false, errors.New("unable to verify given signature")
	}
	return true, nil
}
