//go:generate go run src.go -out ../../and_amd64.s -stubs ../../and_stubs_amd64.go -pkg and
package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	gen("and", VPAND, "Sets dst to the bitwise and of a and b")
	gen("or", VPOR, "Sets dst to the bitwise or of a and b")
	gen("andNot", VPANDN, "Sets dst to the bitwise and of not(a) and b")
	genPopcnt()
	genMemset()
	Generate()
}

func gen(name string, op func(Op, Op, Op), doc string) {
	TEXT(name+"AVX2", NOSPLIT, "func(dst, a, b *byte, l uint64)")

	Pragma("noescape")

	Doc(doc + " assuming all are 256*l bytes")
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
		op(bs[i], as[i], bs[i])
	}
	for i := 0; i < len(as); i++ {
		VMOVDQU(bs[i], Mem{Base: dst, Disp: 32 * i})
	}

	ADDQ(U32(256), a)
	ADDQ(U32(256), b)
	ADDQ(U32(256), dst)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	RET()
}

func genPopcnt() {
	TEXT("popcntAsm", NOSPLIT, "func(a *byte, l uint64) int")

	Pragma("noescape")

	Doc("Counts the number of bits set in a assuming all are 64*l bytes")
	a := Load(Param("a"), GP64())
	l := Load(Param("l"), GP64())

	ret := GP64()

	as := []Op{GP64(), GP64(), GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}
	intermediates := []Op{GP64(), GP64(), GP64(), GP64(), GP64(), GP64(), GP64(), GP64()}

	Doc("Zero the return register")
	XORQ(ret, ret)

	Label("loop")

	for i := 0; i < len(as); i++ {
		MOVQ(Mem{Base: a, Disp: 8 * i}, as[i])
	}
	for i := 0; i < len(as); i++ {
		POPCNTQ(as[i], intermediates[i])
	}
	for i := 0; i < len(as); i++ {
		ADDQ(intermediates[i], ret)
	}

	ADDQ(U32(len(as)*8), a)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	Store(ret, ReturnIndex(0))
	RET()
}

func genMemset() {
	const rounds = 1
	TEXT("memsetAVX2", NOSPLIT, "func(dst *byte, l uint64, b byte)")

	Pragma("noescape")

	Doc("Sets each byte in dst to b")
	dst := Load(Param("dst"), GP64())
	l := Load(Param("l"), GP64())

	bRepeated := YMM()
	b, err := Param("b").Resolve()
	if err != nil {
		panic(err)
	}
	VPBROADCASTB(b.Addr, bRepeated)

	Label("loop")

	for i := 0; i < rounds; i++ {
		VMOVDQU(bRepeated, Mem{Base: dst, Disp: 32 * i})
	}

	ADDQ(U32(32*rounds), dst)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	RET()
}
