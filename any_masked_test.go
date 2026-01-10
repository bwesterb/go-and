package and

import (
	"math/rand/v2"
	"reflect"
	"runtime"
	"testing"
)

func anyMaskedNaive(a, b []byte) bool {
	for i := range a {
		if a[i]&b[i] != 0 {
			return true
		}
	}
	return false
}

func testAgainstBool(t *testing.T, fancy, generic func(a, b []byte) bool, size int) {
	a := make([]byte, size)
	b := make([]byte, size)
	idx := rand.IntN(size)
	a[idx] = uint8(rand.Int())
	b[idx] = uint8(rand.Int())
	r1 := fancy(a, b)
	r2 := generic(a, b)
	if r1 != r2 {
		t.Fatalf("%s produced a different result from %s at length %d:\n%t\n%t", runtime.FuncForPC(reflect.ValueOf(fancy).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(generic).Pointer()).Name(), size, r1, r2)
	}
}

func TestAnyMasked(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainstBool(t, AnyMasked, anyMaskedNaive, size)
		testAgainstBool(t, anyMaskedGeneric, anyMaskedNaive, size)
		for j := 0; j < 10; j++ {
			testAgainstBool(t, AnyMasked, anyMaskedNaive, size+rand.IntN(100))
			testAgainstBool(t, anyMaskedGeneric, anyMaskedNaive, size+rand.IntN(100))
		}
	}
}

func BenchmarkAnyMasked(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AnyMasked(a, bb)
	}
}

func BenchmarkAnyMaskedGeneric(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		anyMaskedGeneric(a, bb)
	}
}

func BenchmarkAnyMaskedNaive(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		anyMaskedNaive(a, bb)
	}
}
