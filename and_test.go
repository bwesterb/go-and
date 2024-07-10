package and

import (
	"bytes"
	"math/rand/v2"
	"testing"
)

func testAgainstGeneric(t *testing.T, size int) {
	a := make([]byte, size)
	b := make([]byte, size)
	c1 := make([]byte, size)
	c2 := make([]byte, size)
	And(c1, a, b)
	andGeneric(c2, a, b)
	if !bytes.Equal(c1, c2) {
		t.Fatal()
	}
}

func TestAgainstGeneric(t *testing.T) {
	for i := 0; i < 20; i++ {
		size := 1 << i
		testAgainstGeneric(t, size)
		for j := 0; j < 10; j++ {
			testAgainstGeneric(t, size+rand.IntN(100))
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
