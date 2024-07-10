//go:build !amd64 && !arm64

package and

func and(dst, a, b []byte) {
	andGeneric(dst, a, b)
}
