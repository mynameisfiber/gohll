package gohll

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceUint32(t *testing.T) {
	N := uint32(0xf0f0f0f0)
	start := uint8(6)
	stop := uint8(1)
	ideal := uint32(56)
	result := sliceUint32(N, start, stop)
	assert.Equal(t, result, ideal, "Incorrect uint32 slice")

	N = uint32(0xf0f0f0f0)
	start = uint8(31)
	stop = uint8(0)
	ideal = N
	result = sliceUint32(N, start, stop)
	assert.Equal(t, result, ideal, "Incorrect uint32 slice")
}

func TestSliceUint64(t *testing.T) {
	N := uint64(0xf0f0f0f0f0f0f0f0)
	start := uint8(6)
	stop := uint8(1)
	ideal := uint64(56)
	result := sliceUint64(N, start, stop)
	assert.Equal(t, result, ideal, "Incorrect uint64 slice")

	N = uint64(0xf0f0f0f0f0f0f0f0)
	start = uint8(63)
	stop = uint8(0)
	ideal = N
	result = sliceUint64(N, start, stop)
	assert.Equal(t, result, ideal, "Incorrect uint64 slice")
}
