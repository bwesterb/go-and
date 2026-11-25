//go:build amd64 && !amd64.v4 && !purego

package and

import "golang.org/x/sys/cpu"

// avx512Worthwhile reports whether the AVX-512 kernels should be used.
//
// Besides AVX512F we require AVX512_VBMI2, which none of the CPUs on
// which 512-bit instructions tend to be slower than AVX2 have: Intel
// Skylake-SP/X, Cascade Lake and Cooper Lake (frequency throttling) and
// Knights Landing/Mill. Ice Lake and later, and AMD Zen 4 and later,
// do have it.
func avx512Worthwhile() bool {
	return cpu.X86.HasAVX512F && cpu.X86.HasAVX512VBMI2
}

func hasAVX512VPOPCNTDQ() bool {
	return avx512Worthwhile() && cpu.X86.HasAVX512VPOPCNTDQ
}
