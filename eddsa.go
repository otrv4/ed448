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

// DSAVerify implements EdDSA style verifying for Ed448
// equivalent of goldilocks_ed48_verify
func DSAVerify(sig [114]byte, pub Point, msg []byte) bool {
	rPoint := NewPoint([16]uint32{}, [16]uint32{}, [16]uint32{}, [16]uint32{})
	// should fail if this fails
	rPoint.DSADecode(sig[:])

	// API_NS(scalar_p) challenge_scalar;
	// {
	//     /* Compute the challenge */
	//     hash_ctx_p hash;
	//     hash_init_with_dom(hash,prehashed,0,context,context_len);
	//     hash_update(hash,signature,GOLDILOCKS_EDDSA_448_PUBLIC_BYTES);
	//     hash_update(hash,pubkey,GOLDILOCKS_EDDSA_448_PUBLIC_BYTES);
	//     hash_update(hash,message,message_len);
	//     uint8_t challenge[2*GOLDILOCKS_EDDSA_448_PRIVATE_BYTES];
	//     hash_final(hash,challenge,sizeof(challenge));
	//     hash_destroy(hash);
	//     API_NS(scalar_decode_long)(challenge_scalar,challenge,sizeof(challenge));
	//     goldilocks_bzero(challenge,sizeof(challenge));
	// }
	// API_NS(scalar_sub)(challenge_scalar, API_NS(scalar_zero), challenge_scalar);

	// API_NS(scalar_p) response_scalar;
	// API_NS(scalar_decode_long)(
	//     response_scalar,
	//     &signature[GOLDILOCKS_EDDSA_448_PUBLIC_BYTES],
	//     GOLDILOCKS_EDDSA_448_PRIVATE_BYTES
	// );

	// for (unsigned c=1; c<GOLDILOCKS_448_EDDSA_DECODE_RATIO; c<<=1) {
	//     API_NS(scalar_add)(response_scalar,response_scalar,response_scalar);
	// }

	// /* pk_point = -c(x(P)) + (cx + k)G = kG */
	// API_NS(base_double_scalarmul_non_secret)(
	//     pk_point,
	//     response_scalar,
	//     pk_point,
	//     challenge_scalar
	// );
	// return goldilocks_succeed_if(API_NS(point_eq(pk_point,r_point)));

	return false
}
