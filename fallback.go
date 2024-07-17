//go:build !amd64 && !arm64

package and

func and(dst, a, b []byte) {
	andGeneric(dst, a, b)
}

func or(dst, a, b []byte) {
	orGeneric(dst, a, b)
}

func andNot(dst, a, b []byte) {
	andNotGeneric(dst, a, b)
}
