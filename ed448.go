package ed448

import (
	"crypto/rand"
	"crypto/sha512"
)

type Ed448 interface {
	GenerateKeys() (priv [privKeyBytes]byte, pub [pubKeyBytes]byte, ok bool)
	Sign(priv [privKeyBytes]byte, message []byte) (signature [signatureBytes]byte, ok bool)
	Verify(signature [signatureBytes]byte, message []byte, pub [pubKeyBytes]byte) (valid bool)
	ComputeSecret(private [privKeyBytes]byte, public [pubKeyBytes]byte) (secret [sha512.Size]byte)
}

type ed448 struct{}

func NewEd448() Ed448 {
	return &ed448{}
}

// Generates a private key and its correspondent public key.
// XXX This is missing the symmetricKey
func (ed *ed448) GenerateKeys() (priv [privKeyBytes]byte, pub [pubKeyBytes]byte, ok bool) {
	var err error
	privKey, err := newRadixCurve().generateKey(rand.Reader)
	ok = err == nil

	copy(priv[:], privKey.secretKey())
	copy(pub[:], privKey.publicKey())

	return
}

// Signs a message using the provided private key and returns the signature.
func (ed *ed448) Sign(priv [privKeyBytes]byte, message []byte) (signature [signatureBytes]byte, ok bool) {
	pk := privateKey(priv)
	signature, err := newRadixCurve().sign(message, &pk)
	ok = err == nil
	return
}

// Verify a signature does correspond a message by a public key.
func (ed *ed448) Verify(signature [signatureBytes]byte, message []byte, pub [pubKeyBytes]byte) (valid bool) {
	pk := publicKey(pub)
	valid = newRadixCurve().verify(signature, message, &pk)
	return
}

// ECDH Compute secret according to private key and peer's public key.
func (ed *ed448) ComputeSecret(private [privKeyBytes]byte, public [pubKeyBytes]byte) (secret [sha512.Size]byte) {
	return //sha512.Sum512(newRadixCurve().computeSecret(private, public))
}
