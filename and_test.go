package and

import (
	"bytes"
	"math/rand/v2"
	"reflect"
	"runtime"
	"testing"
)

func xorNaive(dst, a, b []byte) {
	for i := range dst {
		dst[i] = a[i] ^ b[i]
	}
}

func andNaive(dst, a, b []byte) {
	for i := range dst {
		dst[i] = a[i] & b[i]
	}
}

func orNaive(dst, a, b []byte) {
	for i := range dst {
		dst[i] = a[i] | b[i]
	}
}

func andNotNaive(dst, a, b []byte) {
	for i := range dst {
		dst[i] = (^a[i]) & b[i]
	}
}

func createRandomBuffer(size int) []byte {
	ret := make([]byte, size)
	rng := rand.New(rand.NewPCG(rand.Uint64(), 0))
	for i := range ret {
		ret[i] = uint8(rng.UintN(256))
	}
	return ret
}

func testAgainst(t *testing.T, fancy, generic func(dst, a, b []byte), size int) {
	a := createRandomBuffer(size)
	b := createRandomBuffer(size)
	c1 := make([]byte, size)
	c2 := make([]byte, size)
	fancy(c1, a, b)
	generic(c2, a, b)
	if !bytes.Equal(c1, c2) {
		t.Fatalf("%s produced a different result from %s at length %d:\n%x\n%x", runtime.FuncForPC(reflect.ValueOf(fancy).Pointer()).Name(), runtime.FuncForPC(reflect.ValueOf(generic).Pointer()).Name(), size, c1, c2)
	}
}

func TestAnd(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainst(t, And, andNaive, size)
		testAgainst(t, andGeneric, andNaive, size)
		for j := 0; j < 10; j++ {
			testAgainst(t, And, andNaive, size+rand.IntN(100))
			testAgainst(t, andGeneric, andNaive, size+rand.IntN(100))
		}
	}
}

func TestXor(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainst(t, Xor, xorNaive, size)
		testAgainst(t, xorGeneric, xorNaive, size)
		for j := 0; j < 10; j++ {
			testAgainst(t, Xor, xorNaive, size+rand.IntN(100))
			testAgainst(t, xorGeneric, xorNaive, size+rand.IntN(100))
		}
	}
}

func TestOr(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainst(t, Or, orNaive, size)
		testAgainst(t, orGeneric, orNaive, size)
		for j := 0; j < 10; j++ {
			testAgainst(t, Or, orNaive, size+rand.IntN(100))
			testAgainst(t, orGeneric, orNaive, size+rand.IntN(100))
		}
	}
}

func TestAndNot(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainst(t, AndNot, andNotNaive, size)
		testAgainst(t, andNotGeneric, andNotNaive, size)
		for j := 0; j < 10; j++ {
			testAgainst(t, AndNot, andNotNaive, size+rand.IntN(100))
			testAgainst(t, andNotGeneric, andNotNaive, size+rand.IntN(100))
		}
	}
}

func BenchmarkAnd(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		And(a, a, bb)
	}
}

func BenchmarkAndGeneric(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andGeneric(a, a, bb)
	}
}

func BenchmarkAndNaive(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andNaive(a, a, bb)
	}
}

func BenchmarkOr(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Or(a, a, bb)
	}
}

func BenchmarkOrGeneric(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		orGeneric(a, a, bb)
	}
}

func BenchmarkOrNaive(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		orNaive(a, a, bb)
	}
}

func BenchmarkXor(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Xor(a, a, bb)
	}
}

func BenchmarkXorGeneric(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		xorGeneric(a, a, bb)
	}
}

func BenchmarkXorNaive(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		xorNaive(a, a, bb)
	}
}

func BenchmarkAndNot(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AndNot(a, a, bb)
	}
}

func BenchmarkAndNotGeneric(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andNotGeneric(a, a, bb)
	}
}

func BenchmarkAndNotNaive(b *testing.B) {
	b.StopTimer()
	size := 32000
	a := createRandomBuffer(size)
	bb := createRandomBuffer(size)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andNotNaive(a, a, bb)
	}
}
