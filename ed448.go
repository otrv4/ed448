package ed448

import "crypto/rand"

type Ed448 interface {
	GenerateKeys() (priv, pub []byte, ok bool)
	Sign(priv, message []byte) (signature []byte, ok bool)
	Verify(signature, message, pub []byte) (valid bool)
	ComputeSecret(private, public []byte) (secret []byte)
}

type ed448 struct{}

func NewEd448() Ed448 {
	return &ed448{}
}

// Generates a private key and its correspondent public key.
// XXX This is missing the symmetricKey
func (ed *ed448) GenerateKeys() (priv, pub []byte, ok bool) {
	var err error
	privKey, err := newRadixCurve().generateKey(rand.Reader)
	ok = err == nil

	return privKey.secretKey(), privKey.publicKey(), ok
}

// Signs a message using the provided private key and returns the signature.
func (ed *ed448) Sign(priv, message []byte) (signature []byte, ok bool) {
	return
}

// Verify a signature does correspond a message by a public key.
func (ed *ed448) Verify(signature, message, pub []byte) (valid bool) {
	return
}

// ECDH Compute secret according to private key and peer's public key.
func (ed *ed448) ComputeSecret(private, public []byte) (secret []byte) {
	return newRadixCurve().computeSecret(private, public)
}
