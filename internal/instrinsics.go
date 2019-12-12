package galoisfield

import (
	"math/bits"
)

type uint128 struct {
	hi, lo uint64
}

func widemul32(a, b uint32) uint64 {
	c := uint64(a) * uint64(b)
	return c
}

func widemul64(a, b uint64) uint128 {
	hi, lo := bits.Mul64(a, b)
	return uint128{hi, lo}
}

func isWord32Zero(a uint32) uint32 {
	var ret uint32
	ret = a - ret
	ret, _ = bits.Sub32(ret, ret, ret)
	return ret
}
