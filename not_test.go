package and

import (
	"math/rand/v2"
	"testing"
)

func notNaive(dst, a []byte) {
	for i, b := range a {
		dst[i] = ^b
	}
}

func testNot(t *testing.T, size int) {
	a := make([]byte, size)
	got := make([]byte, size)
	Not(got, a)
	for i, v := range got {
		if v != ^a[i] {
			t.Errorf("Not failed: %08b != %08b", v, ^a[i])
		}
	}
}

func TestNot(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testNot(t, size)
		for j := 0; j < 10; j++ {
			testNot(t, size+rand.IntN(100))
		}
	}
}

func BenchmarkNot(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Not(a, a)
	}
}

func BenchmarkNotGeneric(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		notGeneric(a, a)
	}
}

func BenchmarkNotNaive(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		notNaive(a, a)
	}
}
