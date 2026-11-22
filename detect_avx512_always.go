//go:build amd64.v4 && !purego

package and

import "golang.org/x/sys/cpu"

// avx512Worthwhile reports whether the AVX-512 kernels should be used.
//
// GOAMD64=v4 guarantees AVX512F, but it also admits CPUs on which
// 512-bit instructions tend to be slower than AVX2, such as Skylake-SP.
// See detect_avx512_runtime.go for why we check AVX512_VBMI2.
func avx512Worthwhile() bool {
	return cpu.X86.HasAVX512VBMI2
}

// VPOPCNTDQ is not part of GOAMD64=v4, so always check at runtime.
func hasAVX512VPOPCNTDQ() bool {
	return avx512Worthwhile() && cpu.X86.HasAVX512VPOPCNTDQ
}
