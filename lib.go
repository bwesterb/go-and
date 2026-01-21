// and implements bitwise and for two byte-slices.
package and

import (
	"encoding/binary"
	"math/bits"
)

// Writes bitwise and of a and b to dst.
//
// Panics if len(a) ≠ len(b), or len(dst) ≠ len(a).
func And(dst, a, b []byte) {
	if len(a) != len(b) || len(b) != len(dst) {
		panic("lengths of a, b and dst must be equal")
	}

	and(dst, a, b)
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

	or(dst, a, b)
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

// Writes bitwise xor of a and b to dst.
//
// Panics if len(a) ≠ len(b), or len(dst) ≠ len(a).
func Xor(dst, a, b []byte) {
	if len(a) != len(b) || len(b) != len(dst) {
		panic("lengths of a, b and dst must be equal")
	}

	xor(dst, a, b)
}

func xorGeneric(dst, a, b []byte) {
	i := 0

	for ; i <= len(a)-8; i += 8 {
		binary.LittleEndian.PutUint64(
			dst[i:],
			binary.LittleEndian.Uint64(a[i:])^binary.LittleEndian.Uint64(b[i:]),
		)
	}

	for ; i < len(a); i++ {
		dst[i] = a[i] ^ b[i]
	}
}

// Writes bitwise and of not(a) and b to dst.
//
// Panics if len(a) ≠ len(b), or len(dst) ≠ len(a).
func AndNot(dst, a, b []byte) {
	if len(a) != len(b) || len(b) != len(dst) {
		panic("lengths of a, b and dst must be equal")
	}

	andNot(dst, a, b)
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

// Writes bitwise not of a to dst.
//
// Panics if len(dst) ≠ len(a).
func Not(dst, a []byte) {
	if len(a) != len(dst) {
		panic("lengths of a and dst must be equal")
	}

	not(dst, a)
}

func notGeneric(dst, a []byte) {
	i := 0

	for ; i <= len(a)-8; i += 8 {
		binary.LittleEndian.PutUint64(dst[i:], ^binary.LittleEndian.Uint64(a[i:]))
	}

	for ; i < len(a); i++ {
		dst[i] = ^a[i]
	}
}

// Popcnt counts the number of bits set in a.
func Popcnt(a []byte) int {
	return popcnt(a)
}

func popcntGeneric(a []byte) int {
	var ret int
	i := 0

	for ; i <= len(a)-8; i += 8 {
		ret += bits.OnesCount64(binary.LittleEndian.Uint64(a[i:]))
	}

	for ; i < len(a); i++ {
		ret += bits.OnesCount8(a[i])
	}
	return ret
}

// PopcntMasked counts the number of bits set in a&b.
//
// Panics if len(a) ≠ len(b).
func PopcntMasked(a, b []byte) int {
	if len(a) != len(b) {
		panic("lengths of a and b must be equal")
	}
	return popcntMasked(a, b)
}

func popcntMaskedGeneric(a, b []byte) int {
	var ret int
	i := 0

	for ; i <= len(a)-8; i += 8 {
		ret += bits.OnesCount64(binary.LittleEndian.Uint64(a[i:]) & binary.LittleEndian.Uint64(b[i:]))
	}

	for ; i < len(a); i++ {
		ret += bits.OnesCount8(a[i] & b[i])
	}
	return ret
}

// Memset sets dst[*] to b.
func Memset(dst []byte, b byte) {
	memset(dst, b)
}

func memsetGeneric(dst []byte, b byte) {
	if b == 0 {
		// Special case that the Go compiler can optimize.
		for i := range dst {
			dst[i] = 0
		}
		return
	}
	eightB := 0x0101010101010101 * uint64(b)
	i := 0
	for ; i <= len(dst)-8; i += 8 {
		binary.LittleEndian.PutUint64(dst[i:], eightB)
	}
	for ; i < len(dst); i++ {
		dst[i] = b
	}
}

// Any returns whether at least one bit is set.
func Any(a []byte) bool {
	return any_(a)
}

func anyGeneric(a []byte) bool {
	i := 0

	for ; i <= len(a)-8; i += 8 {
		if binary.LittleEndian.Uint64(a[i:]) != 0 {
			return true
		}
	}
	for ; i < len(a); i++ {
		if a[i] != 0 {
			return true
		}
	}
	return false
}

// AnyMasked returns whether at least one of the bits of its arguments AND-ed together is set.
func AnyMasked(a, b []byte) bool {
	if len(a) != len(b) {
		panic("lengths of a and b must be equal")
	}

	return anyMasked(a, b)
}

func anyMaskedGeneric(a, b []byte) bool {
	i := 0

	for ; i <= len(a)-8; i += 8 {
		if binary.LittleEndian.Uint64(a[i:])&binary.LittleEndian.Uint64(b[i:]) != 0 {
			return true
		}
	}
	for ; i < len(a); i++ {
		if a[i]&b[i] != 0 {
			return true
		}
	}
	return false
}
