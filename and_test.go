package and

import (
	"testing"
)

func BenchmarkAnd(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		And(a, a, bb)
	}
}

func BenchmarkAndGeneric(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andGeneric(a, a, bb)
	}
}
