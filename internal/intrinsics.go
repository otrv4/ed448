package galoisfield

import (
	"math/bits"
)

const (
	lmask = 0xffffffff
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

// for assembly fun:
// __asm__("subs %0, %1, #1;\n\tsbc %0, %0, %0" : "=r"(ret) : "r"(a) : "cc");
func isWord32Zero(a uint32) bool {
	nn := uint64(a)
	nn = nn - 1
	tmp := uint32(nn >> 32)
	return tmp == lmask
}

// for assembly fun:
// __asm__ volatile("neg %0; sbb %0, %0;" : "+r"(x));
func isWord64Zero(a uint64) bool {
	var hi uint64
	a = a - 1
	hi = hi - 1
	a = a >> 32
	hi = hi >> 32
	tmp := a & hi

	return tmp == lmask
}
