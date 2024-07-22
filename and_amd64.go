//go:build !purego

package and

func and(dst, a, b []byte) {
	l := uint64(0)
	if hasAVX2() {
		l = uint64(len(a)) >> 8
		if l != 0 {
			andAVX2(&dst[0], &a[0], &b[0], l)
		}
		l <<= 8
	}
	andGeneric(dst[l:], a[l:], b[l:])
}

func or(dst, a, b []byte) {
	l := uint64(0)
	if hasAVX2() {
		l = uint64(len(a)) >> 8
		if l != 0 {
			orAVX2(&dst[0], &a[0], &b[0], l)
		}
		l <<= 8
	}
	orGeneric(dst[l:], a[l:], b[l:])
}

func andNot(dst, a, b []byte) {
	l := uint64(0)
	if hasAVX2() {
		l = uint64(len(a)) >> 8
		if l != 0 {
			andNotAVX2(&dst[0], &a[0], &b[0], l)
		}
		l <<= 8
	}
	andNotGeneric(dst[l:], a[l:], b[l:])
}

func popcnt(a []byte) int {
	l := uint64(0)
	var ret int
	if hasPopcnt() {
		l = uint64(len(a)) >> 6
		if l != 0 {
			ret = popcntAsm(&a[0], l)
		}
		l <<= 6
	}
	ret += popcntGeneric(a[l:])
	return ret
}

func memset(dst []byte, b byte) {
	l := uint64(0)
	if hasAVX2() {
		l = uint64(len(dst)) >> 5
		if l != 0 {
			memsetAVX2(&dst[0], l, b)
		}
		l <<= 5
	}
	memsetGeneric(dst[l:], b)
}
