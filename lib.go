package and

func And(dst, a, b []byte) {
	if len(a) != len(b) || len(b) != len(dst) {
		panic("lengths of a, b and dst must be equal")
	}

	and(dst, a, b)
}
