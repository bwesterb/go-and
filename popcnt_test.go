package and

import (
	"math/bits"
	"math/rand/v2"
	"testing"
)

func popcntNaive(a []byte) int {
	var ret int
	for i := range a {
		ret += bits.OnesCount8(a[i])
	}
	return ret
}

func testPopcntAgainstGeneric(t *testing.T, size int) {
	a := createRandomBuffer(size)
	got := Popcnt(a)
	want := popcntGeneric(a)
	if got != want {
		t.Fatalf("Popcnt produced a different result from popcntGeneric at length %d: %d; want %d", size, got, want)
	}
}

func TestPopcntAgainstGeneric(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testPopcntAgainstGeneric(t, size)
		for j := 0; j < 10; j++ {
			testPopcntAgainstGeneric(t, size+rand.IntN(100))
		}
	}
}

func BenchmarkPopcnt(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Popcnt(a)
	}
}

func BenchmarkPopcntGeneric(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = popcntGeneric(a)
	}
}

func BenchmarkPopcntNaive(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = popcntNaive(a)
	}
}
