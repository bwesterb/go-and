//go:build !amd64

package and

func and(dst, a, b []byte) {
	andGeneric(dst, a, b)
}
