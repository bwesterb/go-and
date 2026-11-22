//go:generate go run src.go -out ../../and_amd64.s -stubs ../../and_stubs_amd64.go -pkg and
package main

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

func main() {
	ConstraintExpr("!purego")

	gen("and", evex(VPANDQ), AVX512, "Sets dst to the bitwise and of a and b")
	gen("and", VPAND, AVX2, "Sets dst to the bitwise and of a and b")
	gen("and", VPAND, AVX, "Sets dst to the bitwise and of a and b")
	gen("or", evex(VPORQ), AVX512, "Sets dst to the bitwise or of a and b")
	gen("or", VPOR, AVX2, "Sets dst to the bitwise or of a and b")
	gen("or", VPOR, AVX, "Sets dst to the bitwise or of a and b")
	gen("xor", evex(VPXORQ), AVX512, "Sets dst to the bitwise xor of a and b")
	gen("xor", VPXOR, AVX2, "Sets dst to the bitwise xor of a and b")
	gen("xor", VPXOR, AVX, "Sets dst to the bitwise xor of a and b")
	gen("andNot", evex(VPANDNQ), AVX512, "Sets dst to the bitwise and of not(a) and b")
	gen("andNot", VPANDN, AVX2, "Sets dst to the bitwise and of not(a) and b")
	gen("andNot", VPANDN, AVX, "Sets dst to the bitwise and of not(a) and b")
	genNot(AVX512)
	genNot(AVX2)
	genNot(AVX)
	genPopcntAVX512()
	genPopcnt()
	genPopcntMasked(AVX)
	genMemset(AVX512)
	genMemset(AVX2)
	genMemset(AVX)
	genAnyAVX512()
	genAnyAVX()
	genAnyMaskedAVX512()
	genAnyMaskedAVX()
	Generate()
}

type AVXLevel string

var (
	AVX    AVXLevel = "AVX"
	AVX2   AVXLevel = "AVX2"
	AVX512 AVXLevel = "AVX512"
)

func (a AVXLevel) Bits() int {
	switch a {
	case AVX:
		return 128
	case AVX2:
		return 256
	case AVX512:
		return 512
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
	case AVX512:
		return ZMM()
	default:
		panic("invalid level")
	}
}

// Adapts a variadic EVEX instruction to the three-operand form used by gen.
func evex(op func(...Op)) func(Op, Op, Op) {
	return func(a, b, c Op) { op(a, b, c) }
}

// Move instruction for unaligned loads and stores of a full vector register.
func (a AVXLevel) MOVDQU(src, dst Op) {
	if a == AVX512 {
		VMOVDQU64(src, dst)
	} else {
		VMOVDQU(src, dst)
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
		avxLevel.MOVDQU(Mem{Base: a, Disp: avxLevel.Bytes() * i}, as[i])
		avxLevel.MOVDQU(Mem{Base: b, Disp: avxLevel.Bytes() * i}, bs[i])
	}
	for i := 0; i < len(as); i++ {
		op(bs[i], as[i], bs[i])
	}
	for i := 0; i < len(as); i++ {
		avxLevel.MOVDQU(bs[i], Mem{Base: dst, Disp: avxLevel.Bytes() * i})
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

	var allOnes Op
	if avxLevel != AVX512 {
		Comment("Initialize this register to all ones, so we can XOR with it to simulate a NOT")
		allOnes = avxLevel.CreateRegister()
		VPCMPEQB(allOnes, allOnes, allOnes)
	}

	var as []Op
	for i := 0; rounds > i; i++ {
		as = append(as, avxLevel.CreateRegister())
	}

	Label("loop")

	for i := 0; i < len(as); i++ {
		avxLevel.MOVDQU(Mem{Base: a, Disp: avxLevel.Bytes() * i}, as[i])
	}
	for i := 0; i < len(as); i++ {
		if avxLevel == AVX512 {
			// VPTERNLOGQ with truth table 0x55 computes the NOT of its
			// first source operand directly, without an all-ones register.
			VPTERNLOGQ(U8(0x55), as[i], as[i], as[i])
		} else {
			VPXOR(as[i], allOnes, as[i])
		}
	}
	for i := 0; i < len(as); i++ {
		avxLevel.MOVDQU(as[i], Mem{Base: dst, Disp: avxLevel.Bytes() * i})
	}

	ADDQ(U32(avxLevel.Bytes()*rounds), a)
	ADDQ(U32(avxLevel.Bytes()*rounds), dst)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	VZEROALL()
	RET()
}

// Sums the eight quadwords of the ZMM register z into the general purpose
// register ret, clobbering z.
func sumQuadwords(z reg.VecVirtual, ret Op) {
	hi := YMM()
	VEXTRACTI64X4(U8(1), z, hi)
	VPADDQ(hi, z.AsY(), z.AsY())
	hiX := XMM()
	VEXTRACTI128(U8(1), z.AsY(), hiX)
	VPADDQ(hiX, z.AsX(), z.AsX())
	VPSHUFD(U8(0x4e), z.AsX(), hiX)
	VPADDQ(hiX, z.AsX(), z.AsX())
	VMOVQ(z.AsX(), ret)
}

func genPopcntAVX512() {
	TEXT("popcntAVX512", NOSPLIT, "func(a *byte, l uint64) int")

	Pragma("noescape")

	const bytes = 64
	const rounds = 8
	const accs = 4

	Doc(fmt.Sprintf("Counts the number of bits set in a assuming it is %d*l bytes", bytes*rounds))
	a := Load(Param("a"), GP64())
	l := Load(Param("l"), GP64())

	var counts []Op
	for range rounds {
		counts = append(counts, ZMM())
	}
	var acc []reg.VecVirtual
	for range accs {
		z := ZMM()
		VPXORQ(z, z, z)
		acc = append(acc, z)
	}

	Label("loop")

	Comment("Count the bits in each quadword, and accumulate into several",
		"registers to avoid a long dependency chain")
	for i, r := range counts {
		VPOPCNTQ(Mem{Base: a, Disp: bytes * i}, r)
	}
	for i, r := range counts {
		VPADDQ(r, acc[i%accs], acc[i%accs])
	}

	ADDQ(U32(bytes*rounds), a)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	VPADDQ(acc[1], acc[0], acc[0])
	VPADDQ(acc[3], acc[2], acc[2])
	VPADDQ(acc[2], acc[0], acc[0])
	ret := GP64()
	sumQuadwords(acc[0], ret)

	VZEROALL()
	Store(ret, ReturnIndex(0))
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

	Comment("Zero the return register")
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

func genPopcntMasked(avxLevel AVXLevel) {
	TEXT("popcntMaskedAsm"+string(avxLevel), NOSPLIT, "func(a, b *byte, l uint64) int")

	Pragma("noescape")

	const rounds = 2
	gpsPerWide := avxLevel.Bits() / 64

	Doc(fmt.Sprintf("Counts the number of bits set in a&b assuming both are %d*l bytes", rounds*avxLevel.Bytes()))
	a := Load(Param("a"), GP64())
	b := Load(Param("b"), GP64())
	l := Load(Param("l"), GP64())

	ret := GP64()
	Comment("Zero the return register")
	XORQ(ret, ret)

	var as, bs []Op
	for i := 0; rounds > i; i++ {
		as = append(as, avxLevel.CreateRegister())
		bs = append(bs, avxLevel.CreateRegister())
	}
	var split, counts []Op
	for i := 0; len(as)*gpsPerWide > i; i++ {
		split = append(split, GP64())
		counts = append(counts, GP64())
	}

	Label("loop")

	for i := range len(as) {
		VMOVDQU(Mem{Base: a, Disp: avxLevel.Bytes() * i}, as[i])
		VMOVDQU(Mem{Base: b, Disp: avxLevel.Bytes() * i}, bs[i])
	}
	for i := range len(as) {
		VPAND(bs[i], as[i], as[i])
	}
	Commentf("Split the wide registers into general purpose registers on which we can call POPCNT")
	for i := range len(as) {
		for j := range gpsPerWide {
			VPEXTRQ(U8(j), as[i], split[i*gpsPerWide+j])
		}
	}
	for i := range len(split) {
		POPCNTQ(split[i], counts[i])
	}
	for i := range len(split) {
		ADDQ(counts[i], ret)
	}

	ADDQ(U32(avxLevel.Bytes()*len(as)), a)
	ADDQ(U32(avxLevel.Bytes()*len(as)), b)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	VZEROALL()
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
	switch avxLevel {
	case AVX512:
		// VPBROADCASTB to a ZMM register requires AVX512BW, so instead
		// repeat b four times in a general purpose register and use
		// VPBROADCASTD, which only requires AVX512F.
		tmp := GP32()
		MOVBLZX(b.Addr, tmp)
		IMUL3L(U32(0x01010101), tmp, tmp)
		VPBROADCASTD(tmp, bRepeated)
	case AVX2:
		VPBROADCASTB(b.Addr, bRepeated)
	default:
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
		avxLevel.MOVDQU(bRepeated, Mem{Base: dst, Disp: avxLevel.Bytes() * i})
	}

	ADDQ(U32(avxLevel.Bytes()*rounds), dst)
	SUBQ(U32(1), l)
	JNZ(LabelRef("loop"))

	VZEROALL()
	RET()
}

func genAnyAVX512() {
	TEXT("anyAVX512", NOSPLIT, "func(a *byte, l uint64) bool")

	Pragma("noescape")

	const bytes = 64
	const rounds = 8

	Doc(fmt.Sprintf("Returns whether any of the bits of a is set assuming it is %d*l bytes", bytes*rounds))
	a := Load(Param("a"), GP64())
	l := Load(Param("l"), GP64())

	var as []Op
	for range rounds {
		as = append(as, ZMM())
	}
	mask := K()

	Label("loop")

	for i, r := range as {
		VMOVDQU64(Mem{Base: a, Disp: bytes * i}, r)
	}

	Comment("OR all registers together; truth table 0xfe is the OR of three operands")
	VPTERNLOGQ(U8(0xfe), as[0], as[1], as[2])
	VPTERNLOGQ(U8(0xfe), as[3], as[4], as[5])
	VPTERNLOGQ(U8(0xfe), as[2], as[6], as[7])
	VPORQ(as[5], as[7], as[7])

	VPTESTMQ(as[7], as[7], mask)
	KORTESTW(mask, mask)
	JNZ(LabelRef("found"))

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

func genAnyMaskedAVX512() {
	TEXT("anyMaskedAVX512", NOSPLIT, "func(a, b *byte, l uint64) bool")

	Pragma("noescape")

	const bytes = 64
	const rounds = 8

	Doc(fmt.Sprintf("Returns whether any of the bits of the bitwise and of a and b are set assuming all are %d*l bytes", bytes*rounds))
	a := Load(Param("a"), GP64())
	b := Load(Param("b"), GP64())
	l := Load(Param("l"), GP64())

	var as, bs []Op
	for range rounds {
		as = append(as, ZMM())
		bs = append(bs, ZMM())
	}
	mask := K()

	Label("loop")

	for i := range as {
		VMOVDQU64(Mem{Base: a, Disp: bytes * i}, as[i])
		VMOVDQU64(Mem{Base: b, Disp: bytes * i}, bs[i])
	}

	Comment("Compute the OR of a&b over all registers in two independent chains;",
		"truth table 0xf8 computes dst | (src2 & src3)")
	const half = rounds / 2
	acc1, acc2 := bs[0], bs[half]
	VPANDQ(as[0], acc1, acc1)
	VPANDQ(as[half], acc2, acc2)
	for i := 1; i < half; i++ {
		VPTERNLOGQ(U8(0xf8), bs[i], as[i], acc1)
		VPTERNLOGQ(U8(0xf8), bs[half+i], as[half+i], acc2)
	}
	VPORQ(acc1, acc2, acc2)

	VPTESTMQ(acc2, acc2, mask)
	KORTESTW(mask, mask)
	JNZ(LabelRef("found"))

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
