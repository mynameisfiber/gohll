package gohll

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRhoCondition(t *testing.T) {
	// Create two numbers with the
	// same index but one with a larger rho and merge them

	s1 := newSparseList(12, 10)
	s2 := newSparseList(12, 10)

	n1 := encodeHash(0x0f00000f00000000, 12)
	n2 := encodeHash(0x0f000000f0000000, 12)

	s1.Add(n1)
	s2.Add(n2)

	s1.Merge(s2)

	assert.Equal(t, s1.Len(), 1, "Did not merge properly")
	assert.Equal(t, s1.Data[0], n2, "Did not pick correct number")
}

func TestIndexSkipping(t *testing.T) {
	// Create two numbers with the
	// same index but one with a larger rho and merge them

	s1 := newSparseList(12, 10)
	s2 := newSparseList(12, 10)

	n1 := encodeHash(0x0f00000f00000000, 12)
	n2 := encodeHash(0x00f00000f0000000, 12)

	s1.Add(n1)
	s2.Add(n2)

	s1.Merge(s2)

	assert.Equal(t, s1.Len(), 2, "Did not merge properly")
	assert.Equal(t, s1.Data[0], n2, "Did not pick correct number")
	assert.Equal(t, s1.Data[1], n1, "Did not pick correct number")
}
