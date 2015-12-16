package gohll

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGobEmpty(t *testing.T) {
	h := &HLL{}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(h)
	assert.Nil(t, err)
	var h2 HLL
	err = gob.NewDecoder(&buf).Decode(&h2)
	assert.Nil(t, err)
	if h2.tempSet == nil {
		t.Fatal("Gob decode failed for h2.tempSet")
	}
	if h2.sparseList == nil {
		t.Fatal("Gob decode failed for h2.sparseList")
	}
}
func TestGobSparse(t *testing.T) {
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
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(h)
	assert.Nil(t, err)
	var h2 HLL
	err = gob.NewDecoder(&buf).Decode(&h2)
	assert.Nil(t, err)

	assert.Equal(t, h2.format, SPARSE, "Not using sparse mode")
	assert.Equal(t, c, h2.Cardinality())

	for i = 0; i <= 40000; i++ {
		v := rand.Uint32()
		h.Add(fmt.Sprintf("%d-%d", i, v))
		h2.Add(fmt.Sprintf("%d-%d", i, v))
	}

	assert.Equal(t, h.Cardinality(), h2.Cardinality())
}

func TestGobNormal(t *testing.T) {
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
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(h)
	assert.Nil(t, err)
	var h2 HLL
	err = gob.NewDecoder(&buf).Decode(&h2)
	assert.Nil(t, err)

	assert.Equal(t, h2.format, NORMAL, "Not using normal mode")
	assert.Equal(t, c, h2.Cardinality())

	for i = 0; i <= 40000; i++ {
		v := rand.Uint32()
		h.Add(fmt.Sprintf("%d-%d", i, v))
		h2.Add(fmt.Sprintf("%d-%d", i, v))
	}

	assert.Equal(t, h.Cardinality(), h2.Cardinality())
}
