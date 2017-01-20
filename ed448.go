package ed448

import (
	"crypto/rand"
	"golang.org/x/crypto/sha3"
)

// Curve is the interface that wraps the basic curve methods.
//XXX It would be better with the use of privateKey and publicKey types.
type Curve interface {
	GenerateKeys() (priv [privKeyBytes]byte, pub [pubKeyBytes]byte, ok bool)
	Sign(priv [privKeyBytes]byte, message []byte) (signature [signatureBytes]byte, ok bool)
	Verify(signature [signatureBytes]byte, message []byte, pub [pubKeyBytes]byte) (valid bool)
	ComputeSecret(private [privKeyBytes]byte, public [pubKeyBytes]byte) (secret [Size512]byte)
}

type curveT struct{}

var (
	curve = &curveT{}
)

// NewCurve returns a Curve.
func NewCurve() Curve {
	return curve
}

// Generates a private key and its correspondent public key.
func (ed *curveT) GenerateKeys() (priv [privKeyBytes]byte, pub [pubKeyBytes]byte, ok bool) {
	var err error
	privKey, err := ed.generateKey(rand.Reader)
	ok = err == nil

	copy(priv[:], privKey[:])
	copy(pub[:], privKey.publicKey())

	return
}

// Signs a message using the provided private key and returns the signature.
func (ed *curveT) Sign(priv [privKeyBytes]byte, message []byte) (signature [signatureBytes]byte, ok bool) {
	pk := privateKey(priv)
	signature, err := ed.sign(message, &pk)
	ok = err == nil
	return
}

// Verify a signature does correspond a message by a public key.
func (ed *curveT) Verify(signature [signatureBytes]byte, message []byte, pub [pubKeyBytes]byte) (valid bool) {
	pk := publicKey(pub)
	valid = ed.verify(signature, message, &pk)
	return
}

// ECDH Compute secret according to private key and peer's public key.
func (ed *curveT) ComputeSecret(private [privKeyBytes]byte, public [pubKeyBytes]byte) (secret [Size512]byte) {
	k := privateKey(private)
	return sha3.Sum512(ed.computeSecret(k.secretKey(), public[:]))
}
