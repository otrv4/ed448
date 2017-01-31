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

func (c *curveT) decafDerivePrivateKey(sym [symKeyBytes]byte) (privateKey, error) {

	k := privateKey{}
	copy(k.symKey(), sym[:])

	skb := decafPseudoRandomFunction(sym[:])
	secretKey := &decafScalar{}

	barrettDeserializeAndReduce(secretKey[:], skb, &curvePrimeOrder)
	secretKey.serialize(k.secretKey())

	publicKey := c.precomputedScalarMul(secretKey)

	publicKey.decafEncode(k.publicKey())

	return k, nil
}

func (c *curveT) decafGenerateKeys(r io.Reader) (k privateKey, err error) {
	symKey, err := generateSymmetricKey(r)
	if err != nil {
		return
	}

	return c.decafDerivePrivateKey(symKey)
}

func decafDeriveNonce(msg []byte, symKey []byte) (dst decafScalar) {
	h := sha3.NewShake256()
	h.Write(msg)
	h.Write(symKey)
	h.Write([]byte("decaf_448_sign_shake"))
	var out [64]byte
	h.Read(out[:])

	barrettDeserializeAndReduce(dst[:], out[:], &curvePrimeOrder)

	return
}

func decafDeriveChallenge(pubKey []byte, tmpSignature [fieldBytes]byte, msg []byte) (dst decafScalar) {
	h := sha3.NewShake256()
	h.Write(msg)
	h.Write(pubKey)
	h.Write(tmpSignature[:])
	var out [64]byte
	h.Read(out[:])

	barrettDeserializeAndReduce(dst[:], out[:], &curvePrimeOrder)

	return
}

func (c *curveT) decafDeriveTemporarySignature(nonce *decafScalar) (dst [fieldBytes]byte) {
	point := c.precomputedScalarMul(nonce)
	point.decafEncode(dst[:])
	return
}

func (c *curveT) decafSign(msg []byte, k *privateKey) (sig [signatureBytes]byte, err error) {
	secretKeyWords := decafScalar{}
	//XXX: should secret words be destroyed?
	if ok := barrettDeserialize(secretKeyWords[:], k.secretKey(), &curvePrimeOrder); !ok {
		err = errors.New("Corrupted private key")
		return
	}

	nonce := decafDeriveNonce(msg, k.symKey())
	tmpSignature := c.decafDeriveTemporarySignature(&nonce)
	challenge := decafDeriveChallenge(k.publicKey(), tmpSignature, msg)

	challenge.scalarMul(&challenge, &secretKeyWords)
	nonce.scalarSub(&nonce, &challenge)

	copy(sig[:fieldBytes], tmpSignature[:])
	nonce.serialize(sig[fieldBytes:])

	//XXX: should nonce and challenge be destroyed?
	return
}
