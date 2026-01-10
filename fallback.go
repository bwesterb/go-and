//go:build (!amd64 && !arm64) || purego

package and

func and(dst, a, b []byte) {
	andGeneric(dst, a, b)
}

func or(dst, a, b []byte) {
	orGeneric(dst, a, b)
}

func xor(dst, a, b []byte) {
	xorGeneric(dst, a, b)
}

func andNot(dst, a, b []byte) {
	andNotGeneric(dst, a, b)
}

func not(dst, a []byte) {
	notGeneric(dst, a)
}

func popcnt(a []byte) int {
	return popcntGeneric(a)
}

func memset(dst []byte, b byte) {
	memsetGeneric(dst, b)
}

func anyMasked(a, b []byte) bool {
	return anyMaskedGeneric(a, b)
}
