package gohll

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
    "math/rand"
    "math"
)

func TestSparse(t *testing.T) {
	h, err := NewHLL(20, 2500)
	assert.Equal(t, err, nil, "Could not create HLL")

    var i float64
	for i = 0; i <= 5000; i += 1 {
		h.Add(fmt.Sprintf("%d", rand.Float64()))
	}

	assert.Equal(t, h.format, SPARSE, "Not using sparse mode")

	c := h.Cardinality()
    if math.Abs(c - i - 1) > 10 {
        t.Fatalf("Error too high for cardinality estimation: %f should be %f", c, i-1)
    }
}

func TestNormal(t *testing.T) {
    h, err := NewHLL(8, 2500)
	assert.Equal(t, err, nil, "Could not create HLL")

    var i float64
	for i = 0; h.format != NORMAL; i += 1 {
		h.Add(fmt.Sprintf("%d", rand.Float64()))
	}

	assert.Equal(t, h.format, NORMAL, "Not using normal mode")

	c := h.Cardinality()
    if math.Abs(c - i - 1) > 10 {
        t.Fatalf("Error too high for cardinality estimation: %f should be %f", c, i-1)
    }
}
