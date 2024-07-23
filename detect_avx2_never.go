//go:build !amd64 || purego

package and

func hasAVX2() bool {
	return false
}
