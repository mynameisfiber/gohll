package gohll

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
    "math/rand"
)

func TestSparse(t *testing.T) {
	h, err := NewHLL(20, 10000)
	assert.Equal(t, err, nil, "Could not create HLL")

	for i := 0; i <= 5000; i += 1 {
		h.Add(fmt.Sprintf("%d", rand.Float64()))
        fmt.Printf("(%d, %0.0f) ", i, h.Cardinality())
	}

	assert.Equal(t, h.format, SPARSE, "Not using sparse mode")

	c := h.Cardinality()

	fmt.Println("m2: ", h.m2)
	fmt.Println("Sizeof sparseList: ", h.sparseList.Len())
	fmt.Println("Sizeof tempSet: ", h.tempSet.Len())
	assert.Equal(t, c, 0, "blah")
}
