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
	for i = 0; i <= 5000; i += 1 {
		h.Add(fmt.Sprintf("%d", rand.Float64()))
	}

	assert.Equal(t, h.format, SPARSE, "Not using sparse mode")

	c := h.Cardinality()
	if math.Abs(c-i-1) > 10 {
		t.Fatalf("Error too high for cardinality estimation: %f should be %d", c, i-1)
	}
}
