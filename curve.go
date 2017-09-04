package ed448

import (
	"crypto/sha512"
	"errors"
	"io"
)

var (
	//p = 0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffffffffffff
	prime, _ = deserialize(serialized{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})

	//This is the prime order of the curve E
	//primeOrder: 0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3
	primeOrder = mustDeserialize(serialized{
		0xf3, 0x44, 0x58, 0xab, 0x92, 0xc2, 0x78,
		0x23, 0x55, 0x8f, 0xc5, 0x8d, 0x72, 0xc2,
		0x6c, 0x21, 0x90, 0x36, 0xd6, 0xae, 0x49,
		0xdb, 0x4e, 0xc4, 0xe9, 0x23, 0xca, 0x7c,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f,
	})

	//edCons: -39081
	edCons = mustDeserialize(serialized{0xa9, 0x98}) // unsigned

	//This is the base point in the curve 4E
	// gx: 22458004029592430018760433409989603624678964163256413424612546168695
	// 0415467406032909029192869357953282578032075146446173674602635247710,
	// gy=29881921007848149267601793044393067343754404015408024209592824137233
	// 1506189835876003536878655418784733982303233503462500531545062832660
	// XXX: change to correct in hex
	basePoint = mustNewPoint([]byte{
		0x9f, 0x93, 0xed, 0x0a, 0x84, 0xde, 0xf0,
		0xc7, 0xa0, 0x4b, 0x3f, 0x03, 0x70, 0xc1,
		0x96, 0x3d, 0xc6, 0x94, 0x2d, 0x93, 0xf3,
		0xaa, 0x7e, 0x14, 0x96, 0xfa, 0xec, 0x9c,
		0x70, 0xd0, 0x59, 0x3c, 0x5c, 0x06, 0x5f,
		0x24, 0x33, 0xf7, 0xad, 0x26, 0x6a, 0x3a,
		0x45, 0x98, 0x60, 0xf4, 0xaf, 0x4f, 0x1b,
		0xff, 0x92, 0x26, 0xea, 0xa0, 0x7e, 0x29,
	},
		[]byte{0x13},
	)
)

func mustNewPoint(x, y []byte) *homogeneousProjective {
	p, err := newPoint(x, y)
	if err != nil {
		panic("failed to create point")
	}

	return p
}

func (c *curveT) multiplyMontgomery(in *bigNumber, scalar decafScalar, nbits, extraDoubles int) (*bigNumber, word) {
	mont := new(montgomery)
	mont.deserialize(in)
	var i, j, n int
	n = (nbits - 1) % wordBits
	pflip := word(0)

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
	//assert(extraDoubles < INT_MAX);
	for j = 0; j < extraDoubles; j++ {
		mont.montgomeryStep()
	}

	out, ok := mont.serialize(in)
	return out, word(ok)
}

func (c *curveT) multiplyByBase(scalar decafScalar) *twExtensible {
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

	schedule := make([]word, scalarWords)
	scheduleScalarForCombs(schedule, scalar)

	var ni *twNiels

	for i := uint(0); i < s; i++ {
		if i != 0 {
			out.double()
		}

		for j := uint(0); j < n; j++ {
			tab := word(0)

			for k := uint(0); k < t; k++ {
				bit := (s - 1 - i) + k*s + j*(s*t)
				if bit < scalarWords*wordBits {
					tab |= (schedule[bit/wordBits] >> (bit % wordBits) & 1) << k
				}
			}

			invert := word(tab>>(t-1)) - 1
			tab ^= invert
			tab &= (1 << (t - 1)) - 1

			ni = baseTable.lookup(j, t, uint(tab))
			ni.conditionalNegate(invert)

			if i != 0 || j != 0 {
				out.addTwNiels(ni)
			} else {
				convertTwNielsToTwExtensible(out, ni)
			}
		}
	}

	//if(!out.OnCurve()){ return nil } //and maybe panic?

	return out
}

// Deserializes an array of bytes (little-endian) into an array of words
func bytesToWords(dst []word, src []byte) {
	wordBytes := uint(wordBits / 8)
	srcLen := uint(len(src))

	dstLen := uint((srcLen + wordBytes - 1) / wordBytes)
	if dstLen < uint(len(dst)) {
		panic("wrong dst size")
	}

	for i := uint(0); i*wordBytes < srcLen; i++ {
		out := word(0)
		for j := uint(0); j < wordBytes && wordBytes*i+j < srcLen; j++ {
			out |= word(src[wordBytes*i+j]) << (8 * j)
		}

		dst[i] = out
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

func generateSymmetricKey(read io.Reader) (symKey [symKeyBytes]byte, err error) {
	_, err = io.ReadFull(read, symKey[:])
	return
}

func (c *curveT) derivePrivateKey(symmetricKey [symKeyBytes]byte) (privateKey, error) {
	k := privateKey{}
	copy(k.symKey(), symmetricKey[:])

	skb := pseudoRandomFunction(symmetricKey)
	secretKey := decafScalar{}
	secretKey.BarretDecode(skb)
	secretKey.serialize(k.secretKey())

	publicKey := c.multiplyByBase(secretKey)
	serializedPublicKey := publicKey.untwistAndDoubleAndSerialize()
	serialize(k.publicKey(), serializedPublicKey)

	return k, nil
}

func (c *curveT) generateKey(read io.Reader) (k privateKey, err error) {
	symKey, err := generateSymmetricKey(read)
	if err != nil {
		return
	}

	return c.derivePrivateKey(symKey)
}

//XXX Is private only the secret part of the privateKey?
func (c *curveT) computeSecret(private, public []byte) []byte {
	var sk decafScalar
	var pub serialized
	copy(pub[:], public)

	msucc := word(lmask)
	pk, succ := deserializeReturnMask(pub)

	msucc &= barrettDeserializeReturnMask(sk[:], private, &curvePrimeOrder)

	ok := word(0)
	pk, ok = c.multiplyMontgomery(pk, sk, scalarBits, 1)
	succ &= ok

	gxy := make([]byte, fieldBytes)
	serialize(gxy, pk)

	//XXX SECURITY should we wipe the temporary variables?

	//XXX add error conditions based on succ and msucc
	return gxy
}

func (c *curveT) sign(msg []byte, k *privateKey) (s [signatureBytes]byte, err error) {
	secretKeyWords := decafScalar{}
	if ok := barrettDeserialize(secretKeyWords[:], k.secretKey(), &curvePrimeOrder); !ok {
		//XXX SECURITY should we wipe secretKeyWords?
		err = errors.New("corrupted private key")
		return
	}

	nonce := deriveNonce(msg, k.symKey())
	tmpSig := c.deriveTemporarySignature(nonce) // 4 * nonce * G
	challenge := deriveChallenge(k.publicKey(), tmpSig, msg)

	//response = 2(nonce - sk*challenge)
	barrettNegate(challenge[:], &curvePrimeOrder)
	barrettMac(nonce[:], challenge[:], secretKeyWords[:], &curvePrimeOrder)
	carry := addExtPacked(nonce[:], nonce[:], nonce[:], lmask)
	barrettReduce(nonce[:], carry, &curvePrimeOrder)

	// signature = tmpSignature || nonce
	copy(s[:fieldBytes], tmpSig[:])
	nonce.serialize(s[fieldBytes:])

	//XXX SECURITY Should we wipe nonce, gsk, secretKeyWords, tmpSig, challenge?

	/* response = 2(nonce_secret - sk*challenge)
	 * Nonce = 8[nonce_secret]*G
	 * PK = 2[sk]*G, except doubled (TODO)
	 * so [2] ( [response]G + 2[challenge]PK ) = Nonce
	 */

	return
}

func (c *curveT) deriveTemporarySignature(nonce decafScalar) (dst [fieldBytes]byte) {
	// tmpSig = 4 * nonce * basePoint
	fourTimesGTimesNonce := c.multiplyByBase(nonce).double().untwistAndDoubleAndSerialize()
	serialize(dst[:], fourTimesGTimesNonce)
	return
}

//XXX Should pubKey have a fixed size here?
func deriveChallenge(pubKey []byte, tmpSignature [fieldBytes]byte, msg []byte) *decafScalar {
	h := sha512.New()
	h.Write(pubKey)
	h.Write(tmpSignature[:])
	h.Write(msg)

	dst := &decafScalar{}

	barrettDeserializeAndReduce(dst[:], h.Sum(nil), &curvePrimeOrder)

	return dst
}

func deriveNonce(msg []byte, symKey []byte) (dst decafScalar) {
	h := sha512.New()
	h.Write([]byte("signonce"))
	h.Write(symKey)
	h.Write(msg)
	h.Write(symKey)

	barrettDeserializeAndReduce(dst[:], h.Sum(nil), &curvePrimeOrder)

	//XXX SECURITY should we wipe r?
	return
}

func (c *curveT) verify(signature [signatureBytes]byte, msg []byte, k *publicKey) bool {
	serPubkey := serialized(*k)
	pk, ok := deserialize(serPubkey)
	if !ok {
		return false
	}

	nonce := &decafScalar{}
	ok = barrettDeserialize(nonce[:], signature[fieldBytes:2*fieldBytes], &curvePrimeOrder)
	if !ok {
		return false
	}

	tmpSig := [fieldBytes]byte{}
	copy(tmpSig[:], signature[:])
	challenge := deriveChallenge(serPubkey[:], tmpSig, msg)

	eph, ok := deserialize(tmpSig)
	if !ok {
		return false
	}

	//pubKeyBytes -> pubKeyWireFormat -> (DESERPT & twist) -> PK(X, y)
	pkPoint, ok := pk.deserializeAndTwistApprox()
	if !ok {
		return false
	}

	//magic
	//PK_2(X, Y) = PK(X,Y) * ????
	linearComboVarFixedVt(pkPoint, challenge, nonce, wnfsTable[:])

	//PK_2(X,Y) -> (untwist & double & SERPT) -> 2*pubKeyWireFormat
	//In the end, this should be = 4 * nonce * G
	pk = pkPoint.untwistAndDoubleAndSerialize()

	return eph.equals(pk)
}
