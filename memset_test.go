package and

import (
	"math/rand/v2"
	"testing"
)

func memsetNaive(dst []byte, b byte) {
	for i := range dst {
		dst[i] = b
	}
}

func testMemset(t *testing.T, size int) {
	a := make([]byte, size)
	Memset(a, 0xff)
	for i, v := range a {
		if v != 0xff {
			t.Errorf("Memset failed to set a[%d] to 0xff", i)
		}
	}
}

func TestMemsetAgainstGeneric(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testMemset(t, size)
		for j := 0; j < 10; j++ {
			testMemset(t, size+rand.IntN(100))
		}
	}
}

func BenchmarkMemset(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Memset(a, 0xff)
	}
}

func BenchmarkMemsetGeneric(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		memsetGeneric(a, 0xff)
	}
}

func BenchmarkMemsetNaive(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		memsetNaive(a, 0xff)
	}
}
