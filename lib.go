// and implements bitwise and for two byte-slices.
package and

import (
	"encoding/binary"
)

// Writes bitwise and of a and b to dst.
//
// Panics if len(a) ≠ len(b), or len(dst) ≠ len(a).
func And(dst, a, b []byte) {
	if len(a) != len(b) || len(b) != len(dst) {
		panic("lengths of a, b and dst must be equal")
	}

	if hasAVX2() {
		and(dst, a, b)
		return
	}
	andGeneric(dst, a, b)
}

func andGeneric(dst, a, b []byte) {
	i := 0

	for ; i <= len(a)-8; i += 8 {
		binary.LittleEndian.PutUint64(
			dst[i:],
			binary.LittleEndian.Uint64(a[i:])&binary.LittleEndian.Uint64(b[i:]),
		)
	}

	for ; i < len(a); i++ {
		dst[i] = a[i] & b[i]
	}
}

// Writes bitwise or of a and b to dst.
//
// Panics if len(a) ≠ len(b), or len(dst) ≠ len(a).
func Or(dst, a, b []byte) {
	if len(a) != len(b) || len(b) != len(dst) {
		panic("lengths of a, b and dst must be equal")
	}

	if hasAVX2() {
		or(dst, a, b)
		return
	}
	orGeneric(dst, a, b)
}

func orGeneric(dst, a, b []byte) {
	i := 0

	for ; i <= len(a)-8; i += 8 {
		binary.LittleEndian.PutUint64(
			dst[i:],
			binary.LittleEndian.Uint64(a[i:])|binary.LittleEndian.Uint64(b[i:]),
		)
	}

	for ; i < len(a); i++ {
		dst[i] = a[i] | b[i]
	}
}

// Writes bitwise and of not(a) and b to dst.
//
// Panics if len(a) ≠ len(b), or len(dst) ≠ len(a).
func AndNot(dst, a, b []byte) {
	if len(a) != len(b) || len(b) != len(dst) {
		panic("lengths of a, b and dst must be equal")
	}

	if hasAVX2() {
		andNot(dst, a, b)
		return
	}
	andNotGeneric(dst, a, b)
}

func andNotGeneric(dst, a, b []byte) {
	i := 0

	for ; i <= len(a)-8; i += 8 {
		binary.LittleEndian.PutUint64(
			dst[i:],
			(^binary.LittleEndian.Uint64(a[i:]))&binary.LittleEndian.Uint64(b[i:]),
		)
	}

	for ; i < len(a); i++ {
		dst[i] = (^a[i]) & b[i]
	}
}
