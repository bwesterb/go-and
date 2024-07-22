//go:build !purego

package and

//go:noescape
func andNEON(dst, a, b *byte, len uint64)

//go:noescape
func orNEON(dst, a, b *byte, len uint64)

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

func andNot(dst, a, b []byte) {
	// TODO: Write a NEON version for this
	andNotGeneric(dst, a, b)
}

func popcnt(a []byte) int {
	// TODO: Write a NEON version for this
	return popcntGeneric(a)
}
