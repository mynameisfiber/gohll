package gohll

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeHash(t *testing.T) {
    p1 := uint8(14)
    x := uint64(0xffffffffffffffff)
    result := EncodeHash(x, p1)
    ideal := uint32(((1<<25)-1) << 1)
    assert.Equal(t, result, ideal, "Encoded Incorrectly")
}

func TestDecodeHash(t *testing.T) {
    p1 := uint8(12)
    x := uint32(((1<<25)-1) << 1)
    index, rho := DecodeHash(x, p1)
    assert.Equal(t, rho, uint8(0), "Did not decode rho properly")
    assert.Equal(t, index, uint32(0xfff), "Did not decode index properly")

    x = uint32(0xffffff00)
    index, rho = DecodeHash(x, p1)
    assert.Equal(t, rho, uint8(8-1), "Did not decode rho properly")
    assert.Equal(t, index, uint32((1<<(p1+1)) - 1), "Did not decode index properly")
}
