package ed448

import (
	"crypto/rand"
	"errors"
)

type Ed448 interface {
	GenerateKeys() (priv, pub []byte, err error)
	Sign(priv, message []byte) (signature []byte, err error)
	Verify(signature, message, pub []byte) (valid bool)
}

type ed448 struct{}

func NewEd448() Ed448 {
	return &ed448{}
}

// Generates a private key and its correspondent public key.
func (ed *ed448) GenerateKeys() (priv, pub []byte, err error) {
	priv, pub, err = newRadixCurve().generateKey(rand.Reader)

	if err != nil {
		errors.New("Generation of keys has failed.")
	}

	return
}

// Signs a message using the provided private key and returns the signature.
func (ed *ed448) Sign(priv, message []byte) (signature []byte, err error) {
	return
}

// Verify a signature does correspond a message by a public key.
func (ed *ed448) Verify(signature, message, pub []byte) (valid bool) {
	return
}
