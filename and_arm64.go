//go:build !purego

package and

//go:noescape
func andNEON(dst, a, b *byte, l uint64)

//go:noescape
func orNEON(dst, a, b *byte, l uint64)

//go:noescape
func xorNEON(dst, a, b *byte, l uint64)

//go:noescape
func andNotNEON(dst, a, b *byte, l uint64)

//go:noescape
func popcntNEON(a *byte, l uint64) uint64

//go:noescape
func memsetNEON(dst *byte, l uint64, b byte)

func and(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		andNEON(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	andGeneric(dst[l:], a[l:], b[l:])
}

func or(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		orNEON(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	orGeneric(dst[l:], a[l:], b[l:])
}

func xor(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		xorNEON(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	xorGeneric(dst[l:], a[l:], b[l:])
}

func andNot(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		andNotNEON(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	andNotGeneric(dst[l:], a[l:], b[l:])
}

func not(dst, a []byte) {
	notGeneric(dst, a)
}

func popcnt(a []byte) int {
	ret := 0
	l := uint64(len(a)) >> 8
	if l != 0 {
		ret = int(popcntNEON(&a[0], l))
		l <<= 8
	}
	return ret + popcntGeneric(a[l:])
}

func memset(dst []byte, b byte) {
	l := uint64(len(dst)) >> 8
	if l != 0 {
		memsetNEON(&dst[0], l, b)
	}
	l <<= 8
	memsetGeneric(dst[l:], b)
}

func anyMasked(a, b []byte) bool {
	// TODO: Write NEON implementation
	return anyMaskedGeneric(a, b)
}
