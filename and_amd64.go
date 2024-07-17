package and

func and(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		andAVX2(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	andGeneric(dst[l:], a[l:], b[l:])
}

func or(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		orAVX2(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	orGeneric(dst[l:], a[l:], b[l:])
}

func andNot(dst, a, b []byte) {
	l := uint64(len(a)) >> 8
	if l != 0 {
		andNotAVX2(&dst[0], &a[0], &b[0], l)
	}
	l <<= 8
	andNotGeneric(dst[l:], a[l:], b[l:])
}
