package ed448

import "encoding/hex"

func bytesFromHex(s string) []byte {
	val, _ := hex.DecodeString(s)
	return val
}
