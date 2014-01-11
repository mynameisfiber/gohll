package gohll

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestSparse(t *testing.T) {
	h, err := NewHLL(20, 2500)
	assert.Equal(t, err, nil, "Could not create HLL")

	var i float64
	for i = 0; i <= 50000; i += 1 {
		h.Add(fmt.Sprintf("%d", rand.Float64()))
	}

	assert.Equal(t, h.format, SPARSE, "Not using sparse mode")

	c := h.Cardinality()
	errorRate := 1.04 / math.Sqrt(float64(h.m1))
	actualError := math.Abs(c-i-1) / (i - 1)
	if actualError > errorRate {
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
}

func TestNormal(t *testing.T) {
	h, err := NewHLL(10, 2500)
	assert.Equal(t, err, nil, "Could not create HLL")

	h.toNormal()
	var i float64
	for i = 0; i <= 100000; i += 1 {
		h.Add(fmt.Sprintf("%d", rand.Float64()))
	}

	assert.Equal(t, h.format, NORMAL, "Not using normal mode")

	c := h.Cardinality()
	errorRate := 1.04 / math.Sqrt(float64(h.m1))
	actualError := math.Abs(c-i-1) / (i - 1)
	if actualError > errorRate {
		t.Fatalf("Error too high for cardinality estimation: %f should be %f (rate of %f instead of %f)", c, i-1, actualError, errorRate)
	}
}
