// +build go1.1

package gohll

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkErrorBounds(t *testing.T, c, i, errorRate float64) {
	actualError := math.Abs(c/(i-1) - 1)
	if actualError > 3*errorRate {
		debug.PrintStack()
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
}

func fnv1a(s string) uint64 {
	h := fnv.New64a()
	h.Write(byteSlice(s))
	return h.Sum64()
}

func BenchmarkAddSparse(b *testing.B) {
	h, _ := NewHLL(20)
	h.sparseList.MaxSize = 1e8

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		h.Add(fmt.Sprintf("%d", i))
	}
}

func BenchmarkAddNormal(b *testing.B) {
	h, _ := NewHLL(20)
	h.ToNormal()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		h.Add(fmt.Sprintf("%d", i))
	}
}

func BenchmarkAddSparseFNV1A(b *testing.B) {
	h, _ := NewHLL(20)
	h.Hasher = fnv1a
	h.sparseList.MaxSize = 1e8

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		h.Add(fmt.Sprintf("%d", i))
	}
}

func BenchmarkAddNormalFNV1A(b *testing.B) {
	h, _ := NewHLL(20)
	h.Hasher = fnv1a
	h.ToNormal()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		h.Add(fmt.Sprintf("%d", i))
	}
}

func BenchmarkCardinalityNormal(b *testing.B) {
	h, _ := NewHLL(20)
	h.ToNormal()
	for i := 0; i <= 10000; i++ {
		h.Add(fmt.Sprintf("%d", i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		h.Cardinality()
	}
}

func BenchmarkCardinalitySparse(b *testing.B) {
	h, _ := NewHLL(20)
	h.sparseList.MaxSize = 1e8
	for i := 0; i <= 10000; i++ {
		h.Add(fmt.Sprintf("%d", i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		h.Cardinality()
	}
}

func TestHolistic(t *testing.T) {
	h, err := NewHLL(10)
	assert.Nil(t, err)
	errorRate := 1.04 / math.Sqrt(float64(h.m1))

	var i float64
	for i = 0.0; i < 1000000; i++ {
		h.Add(fmt.Sprintf("%d-%d", i, rand.Uint32()))

		if int(i+1001)%1000 == 0 {
			c := h.Cardinality()
			checkErrorBounds(t, c, i, errorRate)
		}
	}
}

func TestSparse(t *testing.T) {
	h, err := NewHLL(20)
	assert.Nil(t, err)

	var i float64
	for i = 0; i <= 50000; i++ {
		h.Add(fmt.Sprintf("%d-%d", i, rand.Uint32()))
	}

	assert.Equal(t, h.format, SPARSE, "Not using sparse mode")

	c := h.Cardinality()
	errorRate := 1.04 / math.Sqrt(float64(h.m2))
	checkErrorBounds(t, c, i, errorRate)
}

func TestNormal(t *testing.T) {
	h, err := NewHLL(10)
	assert.Nil(t, err)

	h.ToNormal()
	var i float64
	for i = 0; i <= 100000; i++ {
		h.Add(fmt.Sprintf("%d-%d", i, rand.Uint32()))
	}

	assert.Equal(t, h.format, NORMAL, "Not using normal mode")

	c := h.Cardinality()
	errorRate := 1.04 / math.Sqrt(float64(h.m1))
	checkErrorBounds(t, c, i, errorRate)
}

func TestModeChange(t *testing.T) {
	h, err := NewHLL(10)
	assert.Nil(t, err)

	var i float64
	for i = 0; i < 100000; i++ {
		h.Add(fmt.Sprintf("%d-%d", i, rand.Uint32()))
	}

	assert.Equal(t, h.format, NORMAL, "Did not convert to normal mode")
	c := h.Cardinality()
	errorRate := 1.04 / math.Sqrt(float64(h.m1))
	checkErrorBounds(t, c, i, errorRate)
}

func TestUnionNormalNormal(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h1.ToNormal()
	h2.ToNormal()

	testSetOperations(t, h1, h2)
}

func TestUnionNormalSparse(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h1.ToNormal()
	h2.sparseList.MaxSize = 1e8

	testSetOperations(t, h1, h2)
}

func TestUnionSparseNormal(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h1.sparseList.MaxSize = 1e8
	h2.ToNormal()

	testSetOperations(t, h1, h2)
	assert.Equal(t, h1.format, NORMAL, "Did not convert h1 to normal mode")
}

func TestUnionSparseSparse(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h1.sparseList.MaxSize = 1e8
	h2.sparseList.MaxSize = 1e8

	testSetOperations(t, h1, h2)
}

func testSetOperations(t *testing.T, h1, h2 *HLL) {
	var i float64
	h1Format := h1.format
	for i = 0; i <= 75000; i++ {
		h1.Add(fmt.Sprintf("%d", i))
	}
	assert.Equal(t, h1Format, h1.format)

	h2Format := h2.format
	for i = 25000; i <= 100000; i++ {
		h2.Add(fmt.Sprintf("%d", i))
	}
	assert.Equal(t, h2Format, h2.format)

	errorRate := 1.04 / math.Sqrt(float64(h1.m1))

	c, _ := h1.CardinalityUnion(h2)
	checkErrorBounds(t, c, i, errorRate)

	c, _ = h1.CardinalityIntersection(h2)
	checkErrorBounds(t, c, 50000, errorRate)

	h1.Union(h2)
	c = h1.Cardinality()
	checkErrorBounds(t, c, i, errorRate)
}
