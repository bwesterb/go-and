//go:build !purego

package and

func and(dst, a, b []byte) {
	l := uint64(0)
	if len(a) >= 512 && avx512Worthwhile() {
		l = uint64(len(a)) >> 9
		andAVX512(&dst[0], &a[0], &b[0], l)
		l <<= 9
	}
	if hasAVX2() {
		if n := uint64(len(a)-int(l)) >> 8; n != 0 {
			andAVX2(&dst[l], &a[l], &b[l], n)
			l += n << 8
		}
	} else if hasAVX() {
		l = uint64(len(a)) >> 7
		if l != 0 {
			andAVX(&dst[0], &a[0], &b[0], l)
		}
		l <<= 7
	}
	andGeneric(dst[l:], a[l:], b[l:])
}

func or(dst, a, b []byte) {
	l := uint64(0)
	if len(a) >= 512 && avx512Worthwhile() {
		l = uint64(len(a)) >> 9
		orAVX512(&dst[0], &a[0], &b[0], l)
		l <<= 9
	}
	if hasAVX2() {
		if n := uint64(len(a)-int(l)) >> 8; n != 0 {
			orAVX2(&dst[l], &a[l], &b[l], n)
			l += n << 8
		}
	} else if hasAVX() {
		l = uint64(len(a)) >> 7
		if l != 0 {
			orAVX(&dst[0], &a[0], &b[0], l)
		}
		l <<= 7
	}
	orGeneric(dst[l:], a[l:], b[l:])
}

func xor(dst, a, b []byte) {
	l := uint64(0)
	if len(a) >= 512 && avx512Worthwhile() {
		l = uint64(len(a)) >> 9
		xorAVX512(&dst[0], &a[0], &b[0], l)
		l <<= 9
	}
	if hasAVX2() {
		if n := uint64(len(a)-int(l)) >> 8; n != 0 {
			xorAVX2(&dst[l], &a[l], &b[l], n)
			l += n << 8
		}
	} else if hasAVX() {
		l = uint64(len(a)) >> 7
		if l != 0 {
			xorAVX(&dst[0], &a[0], &b[0], l)
		}
		l <<= 7
	}
	xorGeneric(dst[l:], a[l:], b[l:])
}

func andNot(dst, a, b []byte) {
	l := uint64(0)
	if len(a) >= 512 && avx512Worthwhile() {
		l = uint64(len(a)) >> 9
		andNotAVX512(&dst[0], &a[0], &b[0], l)
		l <<= 9
	}
	if hasAVX2() {
		if n := uint64(len(a)-int(l)) >> 8; n != 0 {
			andNotAVX2(&dst[l], &a[l], &b[l], n)
			l += n << 8
		}
	} else if hasAVX() {
		l = uint64(len(a)) >> 7
		if l != 0 {
			andNotAVX(&dst[0], &a[0], &b[0], l)
		}
		l <<= 7
	}
	andNotGeneric(dst[l:], a[l:], b[l:])
}

func not(dst, a []byte) {
	l := uint64(0)
	if hasAVX2() {
		l = uint64(len(a)) >> 8
		if l != 0 {
			notAVX2(&dst[0], &a[0], l)
		}
		l <<= 8
	} else if hasAVX() {
		l = uint64(len(a)) >> 7
		if l != 0 {
			notAVX(&dst[0], &a[0], l)
		}
		l <<= 7
	}
	notGeneric(dst[l:], a[l:])
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

func popcntMasked(a, b []byte) int {
	l := uint64(0)
	var ret int
	if hasPopcnt() && hasAVX() {
		l = uint64(len(a)) >> 5
		if l != 0 {
			ret = popcntMaskedAsmAVX(&a[0], &b[0], l)
		}
		l <<= 5
	}
	ret += popcntMaskedGeneric(a[l:], b[l:])
	return ret
}

func memset(dst []byte, b byte) {
	l := uint64(0)
	if len(dst) >= 64 && avx512Worthwhile() {
		l = uint64(len(dst)) >> 6
		memsetAVX512(&dst[0], l, b)
		l <<= 6
	}
	if hasAVX2() {
		if n := uint64(len(dst)-int(l)) >> 5; n != 0 {
			memsetAVX2(&dst[l], n, b)
			l += n << 5
		}
	} else if hasAVX() {
		l = uint64(len(dst)) >> 4
		if l != 0 {
			memsetAVX(&dst[0], l, b)
		}
		l <<= 4
	}
	memsetGeneric(dst[l:], b)
}

func any_(a []byte) bool {
	l := uint64(0)
	if hasAVX() {
		l = uint64(len(a)) >> 8
		if l != 0 && anyAVX(&a[0], l) {
			return true
		}
		l <<= 8
	}
	return anyGeneric(a[l:])
}

func anyMasked(a, b []byte) bool {
	l := uint64(0)
	if hasAVX() {
		l = uint64(len(a)) >> 8
		if l != 0 && anyMaskedAVX(&a[0], &b[0], l) {
			return true
		}
		l <<= 8
	}
	return anyMaskedGeneric(a[l:], b[l:])
}
