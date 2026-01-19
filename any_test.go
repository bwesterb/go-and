package and

import (
	"math/rand/v2"
	"reflect"
	"runtime"
	"testing"
)

func anyNaive(a []byte) bool {
	for _, v := range a {
		if v != 0 {
			return true
		}
	}
	return false
}

func testAny(t *testing.T, fancy, generic func(a []byte) bool, size int) {
	a := make([]byte, size)
	idx := rand.IntN(size)
	a[idx] = uint8(rand.Int() & 2)
	r1 := fancy(a)
	r2 := generic(a)
	if r1 != r2 {
		t.Fatalf("%s produced a different result from %s at length %d:\n%t\n%t", runtime.FuncForPC(reflect.ValueOf(fancy).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(generic).Pointer()).Name(), size, r1, r2)
	}
}

func TestAny(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAny(t, Any, anyNaive, size)
		testAny(t, anyGeneric, anyNaive, size)
		for j := 0; j < 10; j++ {
			testAny(t, Any, anyNaive, size+rand.IntN(100))
			testAny(t, anyGeneric, anyNaive, size+rand.IntN(100))
		}
	}
}

func BenchmarkAny(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Any(a)
	}
}

func BenchmarkAnyGeneric(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		anyGeneric(a)
	}
}

func BenchmarkAnyNaive(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		anyNaive(a)
	}
}
