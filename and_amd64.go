package and

func and(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		andAVX2(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	andGeneric(dst[l:], a[l:], b[l:])
}
