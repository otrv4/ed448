package ed448

import (
	"io"
	"math/big"
)

type Goldilocks interface {
	GenerateKey(prg io.Reader) (priv []byte, pub []byte, err error)
}

type goldilocks struct {
	curve curve
}

var mask = []byte{0xff, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f}

func NewBigintsGoldilocks() Goldilocks {
	return &goldilocks{curve: newBigintsCurve()}
}

/*
func NewGoldilocks() Goldilocks {
	return &goldilocks{curve: newField()}
}
*/
// GenerateKey returns a public/private key pair. The private key is
// generated using the given reader, which must return random data.
func (g goldilocks) GenerateKey(rand io.Reader) (priv []byte, pub []byte, err error) {
	n := ed448.n
	bitSize := n.BitLen()
	byteLen := (bitSize + 7) >> 3
	priv = make([]byte, byteLen)

	var x, y *big.Int

	for x == nil {
		_, err = io.ReadFull(rand, priv)
		if err != nil {
			return
		}
		// We have to mask off any excess bits in the case that the size of the
		// underlying field is not a whole number of bytes.
		priv[0] &= mask[bitSize%8]
		// This is because, in tests, rand will return all zeros and we don't
		// want to get the point at infinity and loop forever.
		priv[1] ^= 0x42

		// If the scalar is out of range, sample another random number.
		if new(big.Int).SetBytes(priv).Cmp(n) >= 0 {
			continue
		}

		bx, by := g.curve.multiplyByBase(priv)
		x, y = bx.(*big.Int), by.(*big.Int)
	}

	pub = marshal(g.curve, x, y)
	return
}

// Marshal converts a point into the form specified in section 4.3.6 of ANSI X9.62.
func marshal(curve curve, x, y *big.Int) []byte {
	byteLen := (ed448.size + 7) >> 3

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point

	xBytes := x.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	yBytes := y.Bytes()
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)
	return ret
}

// Unmarshal converts a point, serialized by Marshal, into an x, y pair.
// It is an error if the point is not on the curve. On error, x = nil.
func unmarshal(curve curve, data []byte) (x, y *big.Int) {
	byteLen := (ed448.size + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}
	if data[0] != 4 { // uncompressed form
		return
	}
	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
	if !curve.isOnCurve(x, y) {
		x, y = nil, nil
	}
	return
}
