//go:build !amd64 || purego

package and

func avx512Worthwhile() bool {
	return false
}

func hasAVX512VPOPCNTDQ() bool {
	return false
}
