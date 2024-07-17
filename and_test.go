package and

import (
	"bytes"
	"math/rand/v2"
	"reflect"
	"runtime"
	"testing"
)

func testAgainstGeneric(t *testing.T, fancy, generic func(dst, a, b []byte), size int) {
	a := make([]byte, size)
	b := make([]byte, size)
	c1 := make([]byte, size)
	c2 := make([]byte, size)
	rng := rand.New(rand.NewPCG(0, 0))
	for i := range a {
		a[i] = uint8(rng.UintN(256))
		b[i] = uint8(rng.UintN(256))
	}
	fancy(c1, a, b)
	generic(c2, a, b)
	if !bytes.Equal(c1, c2) {
		t.Fatalf("%s produced a different result from %s at length %d:\n%x\n%x", runtime.FuncForPC(reflect.ValueOf(fancy).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(generic).Pointer()).Name(), size, c1, c2)
	}
}

func TestAndAgainstGeneric(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainstGeneric(t, And, andGeneric, size)
		for j := 0; j < 10; j++ {
			testAgainstGeneric(t, And, andGeneric, size+rand.IntN(100))
		}
	}
}

func TestOrAgainstGeneric(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainstGeneric(t, Or, orGeneric, size)
		for j := 0; j < 10; j++ {
			testAgainstGeneric(t, Or, orGeneric, size+rand.IntN(100))
		}
	}
}

func TestAndNotAgainstGeneric(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainstGeneric(t, AndNot, andNotGeneric, size)
		for j := 0; j < 10; j++ {
			testAgainstGeneric(t, AndNot, andNotGeneric, size+rand.IntN(100))
		}
	}
}

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

func BenchmarkOr(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Or(a, a, bb)
	}
}

func BenchmarkOrGeneric(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		orGeneric(a, a, bb)
	}
}

func BenchmarkAndNot(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AndNot(a, a, bb)
	}
}

func BenchmarkAndNotGeneric(b *testing.B) {
	b.StopTimer()
	size := 1000000
	a := make([]byte, size)
	bb := make([]byte, size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andNotGeneric(a, a, bb)
	}
}
