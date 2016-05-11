package ed448

import (
	"bytes"
	"encoding/binary"
)

type bytesCurve struct {
}

var bsCurve bytesCurve

func newBytesCurve() curve {
	return &bsCurve
}

func (c *bytesCurve) isOnCurve(x, y interface{}) bool {
	// x² + y² = 1 + bx²y²
	//hex representation for 39081
	edsCons := []byte{0x98, 0xA9, 0x0, 0x0}
	x32, y32 := x.([]byte), y.([]byte)
	x2 := mul(x32, x32)
	y2 := mul(y32, y32)
	x2y2 := mul(x2, y2)
	bx2y2 := mul(edsCons, x2y2)
	left := sum(x2, y2)
	right := sub([]byte{0x1, 0x0, 0x0, 0x0}, bx2y2)
	return bytes.Compare(left, right) == 0
}

func (c *bytesCurve) add(x1, y1, x2, y2 interface{}) (x3, y3 interface{}) {
	return make([]byte, 1), make([]byte, 1)
}

func (c *bytesCurve) double(x1, y1 interface{}) (x3, y3 interface{}) {
	return make([]byte, 1), make([]byte, 1)
}

func (c *bytesCurve) multiply(x, y interface{}, k []byte) (kx, ky interface{}) {
	return make([]byte, 1), make([]byte, 1)
}

func (c *bytesCurve) multiplyByBase(k []byte) (kx, ky interface{}) {
	return make([]byte, 1), make([]byte, 1)
}

func sum(a, b []byte) []byte {
	return writeBytes(readUint32(a) + readUint32(b))
}

func sub(a, b []byte) []byte {
	return writeBytes(readUint32(a) - readUint32(b))
}

func mul(a, b []byte) []byte {
	a32, b32, accum := readUint32(a), readUint32(b), uint32(0)
	for a32 > 0 {
		if a32%2 != 0 {
			accum += b32
		}
		a32 = a32 >> 1
		b32 = b32 << 1
	}
	return writeBytes(accum)
}

func readUint32(a []byte) uint32 {
	var a32 uint32
	b := bytes.NewReader(a)
	binary.Read(b, binary.LittleEndian, &a32)
	return a32
}

func writeBytes(a uint32) []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, a)
	return b.Bytes()
}
