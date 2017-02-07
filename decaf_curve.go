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

func (c *decafCurveT) decafDerivePrivateKey(sym [symKeyBytes]byte) (privateKey, error) {
	k := privateKey{}
	copy(k.symKey(), sym[:])

	skb := decafPseudoRandomFunction(sym[:])
	secretKey := &decafScalar{}

	barrettDeserializeAndReduce(secretKey[:], skb, &curvePrimeOrder)
	secretKey.serialize(k.secretKey())

	publicKey := precomputedScalarMul(secretKey)

	publicKey.decafEncode(k.publicKey())

	return k, nil
}

func (c *decafCurveT) decafGenerateKeys(r io.Reader) (k privateKey, err error) {
	symKey, err := generateSymmetricKey(r)
	if err != nil {
		return
	}

	return c.decafDerivePrivateKey(symKey)
}

func decafDeriveNonce(msg []byte, symKey []byte) *decafScalar {
	h := sha3.NewShake256()
	h.Write(msg)
	h.Write(symKey)
	h.Write([]byte("decaf_448_sign_shake"))
	var out [64]byte
	h.Read(out[:])

	dst := &decafScalar{}

	barrettDeserializeAndReduce(dst[:], out[:], &curvePrimeOrder)

	return dst
}

func decafDeriveChallenge(pubKey []byte, tmpSignature [fieldBytes]byte, msg []byte) *decafScalar {
	h := sha3.NewShake256()
	h.Write(msg)
	h.Write(pubKey)
	h.Write(tmpSignature[:])
	var out [64]byte
	h.Read(out[:])

	dst := &decafScalar{}

	barrettDeserializeAndReduce(dst[:], out[:], &curvePrimeOrder)

	return dst
}

func (c *decafCurveT) decafDeriveTemporarySignature(nonce *decafScalar) (dst [fieldBytes]byte) {
	point := precomputedScalarMul(nonce)
	point.decafEncode(dst[:])
	return
}

func (c *decafCurveT) decafSign(msg []byte, k *privateKey) (sig [signatureBytes]byte, err error) {
	secretKeyWords := &decafScalar{}
	//XXX: should secret words be destroyed?
	if ok := barrettDeserialize(secretKeyWords[:], k.secretKey(), &curvePrimeOrder); !ok {
		err = errors.New("Corrupted private key")
		return
	}

	nonce := decafDeriveNonce(msg, k.symKey())
	tmpSignature := c.decafDeriveTemporarySignature(nonce)
	challenge := decafDeriveChallenge(k.publicKey(), tmpSignature, msg)

	challenge.scalarMul(challenge, secretKeyWords)
	nonce.scalarSub(nonce, challenge)

	copy(sig[:fieldBytes], tmpSignature[:])
	nonce.serialize(sig[fieldBytes:])

	//XXX: should nonce and challenge be destroyed?
	return
}

func (c *decafCurveT) decafVerify(signature [signatureBytes]byte, msg []byte, k *publicKey) bool {

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

	ret := decafDecode(point, tmpSig, word(lmask))
	ret &= decafDecode(pkPoint, serPubkey, word(0x00))
	response := &decafScalar{}
	ret &= response.decode(signature[56:])

	pkPoint = decafDoubleNonSecretScalarMul(pkPoint, pkPoint, response, challenge)

	ret &= pkPoint.equals(point)

	return ret == word(lmask)
}
