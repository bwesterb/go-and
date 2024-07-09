package and

func andGeneric(dst, a, b []byte) {
	for i := 0; i < len(a); i++ {
		dst[i] = a[i] & b[i]
	}
}
