package gohll

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func checkErrorBounds(t *testing.T, c, i, errorRate float64) {
	actualError := math.Abs(c/(i-1) - 1)
	if actualError > 3*errorRate {
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
}

func TestHolistic(t *testing.T) {
	h, err := NewHLL(6)
	assert.Nil(t, err)
	errorRate := 1.04 / math.Sqrt(float64(h.m1))

	var i float64
	for i = 0.0; i < 1000000; i += 1 {
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
	for i = 0; i <= 50000; i += 1 {
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

	h.toNormal()
	var i float64
	for i = 0; i <= 100000; i += 1 {
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
	for i = 0; i < 100000; i += 1 {
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

	h1.toNormal()
	h2.toNormal()

	testUnion(t, h1, h2)
}

func TestUnionNormalSparse(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h1.toNormal()
	h2.sparseList.MaxSize = 1e10

	testUnion(t, h1, h2)
}

func TestUnionSparseNormal(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h1.sparseList.MaxSize = 1e10
	h2.toNormal()

	testUnion(t, h1, h2)
	assert.Equal(t, h1.format, NORMAL, "Did not convert h1 to normal mode")
}

func TestUnionSparseSparse(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h1.sparseList.MaxSize = 1e10
	h2.sparseList.MaxSize = 1e10

	testUnion(t, h1, h2)
}

func testUnion(t *testing.T, h1, h2 *HLL) {
	var i float64
	h1Format := h1.format
	for i = 0; i <= 50000; i += 1 {
		h1.Add(fmt.Sprintf("%d", i))
	}
	assert.Equal(t, h1Format, h1.format)

	h2Format := h2.format
	for i = 50000; i <= 100000; i += 1 {
		h2.Add(fmt.Sprintf("%d", i))
	}
	assert.Equal(t, h2Format, h2.format)

	errorRate := 1.04 / math.Sqrt(float64(h1.m1))

	c, _ := h1.UnionCardinality(h2)
	checkErrorBounds(t, c, i, errorRate)

	h1.Union(h2)
	c = h1.Cardinality()
	checkErrorBounds(t, c, i, errorRate)
}
