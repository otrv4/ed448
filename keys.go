package ed448

// secret || public || symmetric
type privateKey [privKeyBytes]byte

func (k *privateKey) secretKey() []byte {
	return k[:fieldBytes]
}

func (k *privateKey) publicKey() []byte {
	return k[fieldBytes : 2*fieldBytes]
}

func (k *privateKey) symKey() []byte {
	return k[2*fieldBytes:]
}

type publicKey [pubKeyBytes]byte
