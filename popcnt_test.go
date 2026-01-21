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

func popcntMaskedNaive(a, b []byte) int {
	var ret int
	for i := range a {
		ret += bits.OnesCount8(a[i] & b[i])
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
	size := 1000000
	a := createRandomBuffer(size)
	b.SetBytes(int64(size))
	for b.Loop() {
		_ = Popcnt(a)
	}
}

func BenchmarkPopcntGeneric(b *testing.B) {
	size := 1000000
	a := createRandomBuffer(size)
	b.SetBytes(int64(size))
	for b.Loop() {
		_ = popcntGeneric(a)
	}
}

func BenchmarkPopcntNaive(b *testing.B) {
	size := 1000000
	a := createRandomBuffer(size)
	b.SetBytes(int64(size))
	for b.Loop() {
		_ = popcntNaive(a)
	}
}

func testPopcntMaskedAgainstGeneric(t *testing.T, size int) {
	a := createRandomBuffer(size)
	b := createRandomBuffer(size)
	got := PopcntMasked(a, b)
	want := popcntMaskedGeneric(a, b)
	if got != want {
		t.Fatalf("PopcntMasked produced a different result from popcntMaskedGeneric at length %d: %d; want %d", size, got, want)
	}
}

func TestPopcntMaskedAgainstGeneric(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testPopcntMaskedAgainstGeneric(t, size)
		for j := 0; j < 10; j++ {
			testPopcntMaskedAgainstGeneric(t, size+rand.IntN(100))
		}
	}
}

func BenchmarkPopcntMasked(b *testing.B) {
	size := 1000000
	a := createRandomBuffer(size)
	c := createRandomBuffer(size)
	b.SetBytes(int64(size))
	for b.Loop() {
		_ = PopcntMasked(a, c)
	}
}

func BenchmarkPopcntMaskedGeneric(b *testing.B) {
	size := 1000000
	a := createRandomBuffer(size)
	c := createRandomBuffer(size)
	b.SetBytes(int64(size))
	for b.Loop() {
		_ = popcntMaskedGeneric(a, c)
	}
}

func BenchmarkPopcntMaskedNaive(b *testing.B) {
	size := 1000000
	a := createRandomBuffer(size)
	c := createRandomBuffer(size)
	b.SetBytes(int64(size))
	for b.Loop() {
		_ = popcntMaskedNaive(a, c)
	}
}
