package and

import (
	"encoding/binary"
)

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
