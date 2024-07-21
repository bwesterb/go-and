package and

import (
	"math/rand/v2"
	"testing"
)

func testPopcntAgainstGeneric(t *testing.T, size int) {
	a := make([]byte, size)
	rng := rand.New(rand.NewPCG(0, 0))
	for i := range a {
		a[i] = uint8(rng.UintN(256))
	}
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
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Popcnt(a)
	}
}

func BenchmarkPopcntGeneric(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = popcntGeneric(a)
	}
}
