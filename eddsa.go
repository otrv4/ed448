package ed448

import (
	"golang.org/x/crypto/sha3"
)

func clamp(val []byte) {
	val[0] &= -(Cofactor)
	val[57-1] = 0
	val[57-2] |= 0x80
}

func deriveKeypair(sym [57]byte) (Scalar, Point) {
	digest := [57]byte{}
	sha3.ShakeSum256(digest[:], sym[:])
	clamp(digest[:])

	r := NewScalar(digest[:])
	r.Halve(r)
	r.Halve(r)
	h := PrecomputedScalarMul(r)

	return r, h
}

func hashWithDom(output []byte, input []byte) {
	sha3.ShakeSum256(output, append(append([]byte("SigEd448"), 0x00, 0x00), input...))
}

// DSASign implements EdDSA style signing for Ed448
// - equivalent of goldilocks_ed448_sign
func DSASign(sym [57]byte, pub Point, msg []byte) [114]byte {
	secret := [114]byte{}
	sha3.ShakeSum256(secret[:], sym[:])
	clamp(secret[:])
	sec := NewScalar(secret[0:57])
	seed := secret[57:]

	nonce := make([]byte, 114)
	hashWithDom(nonce, append(seed, msg...))
	nonceScalar := NewScalar(nonce[:])
	nonceScalar2 := NewScalar()
	nonceScalar2.Halve(nonceScalar)
	nonceScalar2.Halve(nonceScalar2)
	noncePoint := PrecomputedScalarMul(nonceScalar2).DSAEncode()

	challenge := make([]byte, 114)
	hashWithDom(challenge, append(append(noncePoint, pub.DSAEncode()...), msg...))

	challengeScalar := NewScalar(challenge)
	challengeScalar.Mul(challengeScalar, sec)
	challengeScalar.Add(challengeScalar, nonceScalar)

	var sig [114]byte
	copy(sig[:], noncePoint)
	copy(sig[57:], challengeScalar.Encode())

	return sig
}

var scalarFour = NewScalar([]byte{0x04})

// DSAVerify implements EdDSA style verifying for Ed448
// equivalent of goldilocks_ed48_verify
func DSAVerify(sig [114]byte, pub Point, msg []byte) bool {
	pub2 := PointScalarMul(pub, scalarFour)
	sig1 := append([]byte{}, sig[:57]...)
	sig2 := append([]byte{}, sig[57:]...)
	rPoint := NewPoint([16]uint32{}, [16]uint32{}, [16]uint32{}, [16]uint32{})
	rPoint.DSADecode(sig1)
	rPoint = PointScalarMul(rPoint, scalarFour)

	challenge := make([]byte, 114)
	hashWithDom(challenge, append(append(sig1, pub.DSAEncode()...), msg...))
	challengeScalar := NewScalar(challenge)
	challengeScalar.Sub(scalarZero, challengeScalar)

	responseScalar := NewScalar(sig2)
	pk := PointDoubleScalarMulNonsecret(pub2, responseScalar, challengeScalar)

	return pk.Equals(rPoint)
}
