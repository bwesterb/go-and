//go:build amd64 && go1.26 && goexperiment.simd

package and

import "simd/archsimd"

// TODO: Optimize the and() implementation further to match performance, and then use it for all operations (including AVX1)

func and(dst, a, b []byte) {
	var l int
	if hasAVX2() {
		l = len(dst) & ^31
		dst := dst[:l]
		a := a[:l]
		b := b[:l]
		for len(dst) >= 256 {
			_ = a[255]
			_ = b[255]
			a0 := archsimd.LoadUint8x32Slice(a[0:])
			b0 := archsimd.LoadUint8x32Slice(b[0:])
			a1 := archsimd.LoadUint8x32Slice(a[32:])
			b1 := archsimd.LoadUint8x32Slice(b[32:])
			a2 := archsimd.LoadUint8x32Slice(a[64:])
			b2 := archsimd.LoadUint8x32Slice(b[64:])
			a3 := archsimd.LoadUint8x32Slice(a[96:])
			b3 := archsimd.LoadUint8x32Slice(b[96:])
			a4 := archsimd.LoadUint8x32Slice(a[128:])
			b4 := archsimd.LoadUint8x32Slice(b[128:])
			a5 := archsimd.LoadUint8x32Slice(a[160:])
			b5 := archsimd.LoadUint8x32Slice(b[160:])
			a6 := archsimd.LoadUint8x32Slice(a[192:])
			b6 := archsimd.LoadUint8x32Slice(b[192:])
			a7 := archsimd.LoadUint8x32Slice(a[224:])
			b7 := archsimd.LoadUint8x32Slice(b[224:])
			a0 = a0.And(b0)
			a1 = a1.And(b1)
			a2 = a2.And(b2)
			a3 = a3.And(b3)
			a4 = a4.And(b4)
			a5 = a5.And(b5)
			a6 = a6.And(b6)
			a7 = a7.And(b7)
			a0.StoreSlice(dst[0:])
			a1.StoreSlice(dst[32:])
			a2.StoreSlice(dst[64:])
			a3.StoreSlice(dst[96:])
			a4.StoreSlice(dst[128:])
			a5.StoreSlice(dst[160:])
			a6.StoreSlice(dst[192:])
			a7.StoreSlice(dst[224:])

			dst = dst[256:]
			a = a[256:]
			b = b[256:]
		}
		archsimd.ClearAVXUpperBits()
	} else if hasAVX() {
		for ; len(dst)-16 >= l; l += 16 {
			left := archsimd.LoadUint8x16Slice(a[l:])
			right := archsimd.LoadUint8x16Slice(b[l:])
			d := left.And(right)
			d.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	}
	andGeneric(dst[l:], a[l:], b[l:])
}

func or(dst, a, b []byte) {
	var l int
	if hasAVX2() {
		for ; len(dst)-32 >= l; l += 32 {
			left := archsimd.LoadUint8x32Slice(a[l:])
			right := archsimd.LoadUint8x32Slice(b[l:])
			d := left.Or(right)
			d.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	} else if hasAVX() {
		for ; len(dst)-16 >= l; l += 16 {
			left := archsimd.LoadUint8x16Slice(a[l:])
			right := archsimd.LoadUint8x16Slice(b[l:])
			d := left.Or(right)
			d.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	}
	orGeneric(dst[l:], a[l:], b[l:])
}

func xor(dst, a, b []byte) {
	var l int
	if hasAVX2() {
		for ; len(dst)-32 >= l; l += 32 {
			left := archsimd.LoadUint8x32Slice(a[l:])
			right := archsimd.LoadUint8x32Slice(b[l:])
			d := left.Xor(right)
			d.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	} else if hasAVX() {
		for ; len(dst)-16 >= l; l += 16 {
			left := archsimd.LoadUint8x16Slice(a[l:])
			right := archsimd.LoadUint8x16Slice(b[l:])
			d := left.Xor(right)
			d.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	}
	xorGeneric(dst[l:], a[l:], b[l:])
}

func andNot(dst, a, b []byte) {
	var l int
	if hasAVX2() {
		for ; len(dst)-32 >= l; l += 32 {
			left := archsimd.LoadUint8x32Slice(a[l:])
			right := archsimd.LoadUint8x32Slice(b[l:])
			d := right.AndNot(left)
			d.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	} else if hasAVX() {
		for ; len(dst)-16 >= l; l += 16 {
			left := archsimd.LoadUint8x16Slice(a[l:])
			right := archsimd.LoadUint8x16Slice(b[l:])
			d := right.AndNot(left)
			d.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	}
	andNotGeneric(dst[l:], a[l:], b[l:])
}

func not(dst, a []byte) {
	var l int
	if hasAVX2() {
		for ; len(dst)-32 >= l; l += 32 {
			v := archsimd.LoadUint8x32Slice(a[l:])
			v = v.Not()
			v.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	} else if hasAVX() {
		for ; len(dst)-16 >= l; l += 16 {
			v := archsimd.LoadUint8x16Slice(a[l:])
			v = v.Not()
			v.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	}
	notGeneric(dst[l:], a[l:])
}

func popcnt(a []byte) int {
	// TODO: Use the ASM implementation here. archsimd doesn't offer features for this.
	//       This is non-trivial because we've excluded all of and_arm64.s with build tags.
	return popcntGeneric(a)
}

func memset(dst []byte, b byte) {
	var l int
	if hasAVX2() {
		v := archsimd.BroadcastUint8x32(b)
		for ; len(dst)-32 >= l; l += 32 {
			v.StoreSlice(dst[l:])
		}
		archsimd.ClearAVXUpperBits()
	}
	memsetGeneric(dst[l:], b)
}
