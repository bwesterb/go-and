//go:generate go run src.go -out ../../and_amd64.s -stubs ../../and_stubs_amd64.go -pkg and
package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	// Must be called on 32 byte aligned a, b, dst.
	TEXT("andAVX2", NOSPLIT, "func(dst, a, b *byte, l uint64)")

	Pragma("noescape")

	Doc("Sets dst to the bitwise and of a and b assuming all are 256*l bytes")
	a := Load(Param("a"), GP64())
	b := Load(Param("b"), GP64())
	dst := Load(Param("dst"), GP64())
	l := Load(Param("l"), GP64())

	as := []Op{YMM(), YMM(), YMM(), YMM(), YMM(), YMM(), YMM(), YMM()}
	bs := []Op{YMM(), YMM(), YMM(), YMM(), YMM(), YMM(), YMM(), YMM()}

	Label("loop")

	for i := 0; i < len(as); i++ {
		VMOVDQU(Mem{Base: a, Disp: 32 * i}, as[i])
		VMOVDQU(Mem{Base: b, Disp: 32 * i}, bs[i])
	}
	for i := 0; i < len(as); i++ {
		VPAND(as[i], bs[i], bs[i])
	}
	for i := 0; i < len(as); i++ {
		VMOVDQU(Mem{Base: dst, Disp: 32 * i}, bs[i])
	}

	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	RET()
	Generate()
}
