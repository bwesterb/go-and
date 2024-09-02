//go:build !amd64 || purego

package and

func hasAVX() bool {
	return false
}

func hasAVX2() bool {
	return false
}
