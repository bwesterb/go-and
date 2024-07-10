package and

//go:noescape
func andNEON(dst, a, b *byte, len uint64)

func and(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		andNEON(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	andGeneric(dst[l:], a[l:], b[l:])
}
