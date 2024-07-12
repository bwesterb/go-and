//go:build !amd64

package and

func hasAVX2() bool {
	return false
}
