package gohll

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

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
	actualError := math.Abs(c-i-1) / (i - 1)
	if actualError > errorRate {
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
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
	actualError := math.Abs(c-i-1) / (i - 1)
	if actualError > errorRate {
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
}

func TestModeChange(t *testing.T) {
    h, err := NewHLL(10)
	assert.Nil(t, err)
    
    var i float64
    for i = 0; h.format == SPARSE; i += 1 {
		h.Add(fmt.Sprintf("%d-%d", i, rand.Uint32()))
    }

    assert.Equal(t, h.format, NORMAL, "Did not convert to normal mode")
	c := h.Cardinality()
	errorRate := 1.04 / math.Sqrt(float64(h.m1))
	actualError := math.Abs(c-i-1) / (i - 1)
	if actualError > errorRate {
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
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
    testUnion(t, h1, h2)
}

func TestUnionSparseNormal(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

	h2.toNormal()
    testUnion(t, h1, h2)
    assert.Equal(t, h1.format, NORMAL, "Did not convert h1 to normal mode")
}

func TestUnionSparseSparse(t *testing.T) {
	h1, err := NewHLL(10)
	assert.Nil(t, err)

	h2, err := NewHLL(10)
	assert.Nil(t, err)

    testUnion(t, h1, h2)
}

func testUnion(t *testing.T, h1, h2 *HLL) {
	var i float64
	for i = 0; i <= 100000; i += 1 {
		h1.Add(fmt.Sprintf("%d", i))
	}

	h2.toNormal()
	for i = 50000; i <= 150000; i += 1 {
		h2.Add(fmt.Sprintf("%d", i))
	}

    h1.Union(h2)
	c := h1.Cardinality()
	errorRate := 1.04 / math.Sqrt(float64(h1.m1))
	actualError := math.Abs(c-i-1) / (i - 1)
	if actualError > errorRate {
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
}
