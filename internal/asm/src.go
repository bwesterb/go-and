//go:generate go run src.go -out ../../and_amd64.s -stubs ../../and_stubs_amd64.go -pkg and
package main

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	ConstraintExpr("!purego")

	gen("and", VPAND, AVX2, "Sets dst to the bitwise and of a and b")
	gen("and", VPAND, AVX, "Sets dst to the bitwise and of a and b")
	gen("or", VPOR, AVX2, "Sets dst to the bitwise or of a and b")
	gen("or", VPOR, AVX, "Sets dst to the bitwise or of a and b")
	gen("xor", VPXOR, AVX2, "Sets dst to the bitwise xor of a and b")
	gen("xor", VPXOR, AVX, "Sets dst to the bitwise xor of a and b")
	gen("andNot", VPANDN, AVX2, "Sets dst to the bitwise and of not(a) and b")
	gen("andNot", VPANDN, AVX, "Sets dst to the bitwise and of not(a) and b")
	genNot(AVX2)
	genNot(AVX)
	genPopcnt()
	genMemset(AVX2)
	genMemset(AVX)
	genAnyAVX()
	genAnyMaskedAVX()
	Generate()
}

type AVXLevel string

var (
	AVX  AVXLevel = "AVX"
	AVX2 AVXLevel = "AVX2"
)

func (a AVXLevel) Bits() int {
	switch a {
	case AVX:
		return 128
	case AVX2:
		return 256
	default:
		panic("invalid level")
	}
}

func (a AVXLevel) Bytes() int {
	return a.Bits() / 8
}

func (a AVXLevel) CreateRegister() Op {
	switch a {
	case AVX:
		return XMM()
	case AVX2:
		return YMM()
	default:
		panic("invalid level")
	}
}

func gen(name string, op func(Op, Op, Op), avxLevel AVXLevel, doc string) {
	TEXT(name+string(avxLevel), NOSPLIT, "func(dst, a, b *byte, l uint64)")

	Pragma("noescape")

	const rounds = 8

	Doc(fmt.Sprintf("%s assuming all are %d*l bytes", doc, avxLevel.Bits()))
	a := Load(Param("a"), GP64())
	b := Load(Param("b"), GP64())
	dst := Load(Param("dst"), GP64())
	l := Load(Param("l"), GP64())

	var as, bs []Op
	for i := 0; rounds > i; i++ {
		as = append(as, avxLevel.CreateRegister())
		bs = append(bs, avxLevel.CreateRegister())
	}

	Label("loop")

	for i := 0; i < len(as); i++ {
		VMOVDQU(Mem{Base: a, Disp: avxLevel.Bytes() * i}, as[i])
		VMOVDQU(Mem{Base: b, Disp: avxLevel.Bytes() * i}, bs[i])
	}
	for i := 0; i < len(as); i++ {
		op(bs[i], as[i], bs[i])
	}
	for i := 0; i < len(as); i++ {
		VMOVDQU(bs[i], Mem{Base: dst, Disp: avxLevel.Bytes() * i})
	}

	ADDQ(U32(avxLevel.Bytes()*rounds), a)
	ADDQ(U32(avxLevel.Bytes()*rounds), b)
	ADDQ(U32(avxLevel.Bytes()*rounds), dst)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	VZEROALL()
	RET()
}

func genNot(avxLevel AVXLevel) {
	const rounds = 8
	TEXT("not"+string(avxLevel), NOSPLIT, "func(dst, a *byte, l uint64)")

	Pragma("noescape")

	Doc("Bitwise inverts each byte of a into dst")
	dst := Load(Param("dst"), GP64())
	a := Load(Param("a"), GP64())
	l := Load(Param("l"), GP64())

	Comment("Initialize this register to all ones, so we can XOR with it to simulate a NOT")
	allOnes := avxLevel.CreateRegister()
	VPCMPEQB(allOnes, allOnes, allOnes)

	var as []Op
	for i := 0; rounds > i; i++ {
		as = append(as, avxLevel.CreateRegister())
	}

	Label("loop")

	for i := 0; i < len(as); i++ {
		VMOVDQU(Mem{Base: a, Disp: avxLevel.Bytes() * i}, as[i])
	}
	for i := 0; i < len(as); i++ {
		VPXOR(as[i], allOnes, as[i])
	}
	for i := 0; i < len(as); i++ {
		VMOVDQU(as[i], Mem{Base: dst, Disp: avxLevel.Bytes() * i})
	}

	ADDQ(U32(avxLevel.Bytes()*rounds), a)
	ADDQ(U32(avxLevel.Bytes()*rounds), dst)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	VZEROALL()
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

func genMemset(avxLevel AVXLevel) {
	const rounds = 1
	TEXT("memset"+string(avxLevel), NOSPLIT, "func(dst *byte, l uint64, b byte)")

	Pragma("noescape")

	Doc("Sets each byte in dst to b")
	dst := Load(Param("dst"), GP64())
	l := Load(Param("l"), GP64())

	bRepeated := avxLevel.CreateRegister()
	b, err := Param("b").Resolve()
	if err != nil {
		panic(err)
	}
	if avxLevel == AVX2 {
		VPBROADCASTB(b.Addr, bRepeated)
	} else {
		zeroes := GLOBL("zeroes", RODATA|NOPTR)
		DATA(0, U32(0))
		DATA(4, U32(0))
		DATA(8, U32(0))
		DATA(12, U32(0))
		tmp := GP64()
		MOVB(b.Addr, tmp.As8())
		MOVQ(tmp, bRepeated)
		VPSHUFB(zeroes.Offset(0), bRepeated, bRepeated)
	}

	Label("loop")

	for i := 0; i < rounds; i++ {
		VMOVDQU(bRepeated, Mem{Base: dst, Disp: avxLevel.Bytes() * i})
	}

	ADDQ(U32(avxLevel.Bytes()*rounds), dst)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	VZEROALL()
	RET()
}

func genAnyAVX() {
	TEXT("anyAVX", NOSPLIT, "func(a *byte, l uint64) bool")

	Pragma("noescape")

	const bits = 256 // Even though this only needs AVX and not AVX2
	const bytes = bits / 8
	const rounds = 8

	Doc("Returns whether any of the bits of of a is set")
	a := Load(Param("a"), GP64())
	l := Load(Param("l"), GP64())

	var as []Op
	for range rounds {
		as = append(as, YMM())
	}

	Label("loop")

	for i, r := range as {
		VMOVDQU(Mem{Base: a, Disp: bytes * i}, r)
	}
	for _, r := range as {
		VPTEST(r, r)
		JNZ(LabelRef("found"))
	}

	ADDQ(U32(bytes*rounds), a)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	ret := GP8()

	XORB(ret, ret) // return false
	Store(ret, ReturnIndex(0))
	VZEROALL()
	RET()

	Label("found")
	MOVB(U8(1), ret) // return true
	Store(ret, ReturnIndex(0))
	VZEROALL()
	RET()
}

func genAnyMaskedAVX() {
	TEXT("anyMaskedAVX", NOSPLIT, "func(a, b *byte, l uint64) bool")

	Pragma("noescape")

	const bits = 256 // Even though this only needs AVX and not AVX2
	const bytes = bits / 8
	const rounds = 8

	Doc(fmt.Sprintf("Returns whether any of the bits of the bitwise and of a and b are set assuming all are %d*l bytes", bits))
	a := Load(Param("a"), GP64())
	b := Load(Param("b"), GP64())
	l := Load(Param("l"), GP64())

	var as, bs []Op
	for range rounds {
		as = append(as, YMM())
		bs = append(bs, YMM())
	}

	Label("loop")

	for i := range as {
		VMOVDQU(Mem{Base: a, Disp: bytes * i}, as[i])
		VMOVDQU(Mem{Base: b, Disp: bytes * i}, bs[i])
	}
	for i := range as {
		VPTEST(as[i], bs[i])
		JNZ(LabelRef("found"))
	}

	ADDQ(U32(bytes*rounds), a)
	ADDQ(U32(bytes*rounds), b)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	ret := GP8()

	XORB(ret, ret) // return false
	Store(ret, ReturnIndex(0))
	VZEROALL()
	RET()

	Label("found")
	MOVB(U8(1), ret) // return true
	Store(ret, ReturnIndex(0))
	VZEROALL()
	RET()
}
