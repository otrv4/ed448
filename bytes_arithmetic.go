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
	//x2 := x.([]byte) << 1
	//y2 := y.([]byte) << 1
	return false
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
	var a16, b16, s uint32
	r := bytes.NewReader(a)
	binary.Read(r, binary.LittleEndian, &a16)
	r = bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &b16)

	s = a16 + b16

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s)
	return buf.Bytes()
}

func mul(a, b []byte) []byte {
	/*
		c := 0
		for bytes.Compare(a, 1) > 0 {
			if a[len(a)-1] == 1 {
				c = sum(a, c)
			}
			a = a >> 1
			b = b << 1
		}
		return c
	*/
	return make([]byte, 1)
}
