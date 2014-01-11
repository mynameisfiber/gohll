package gohll

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestEncodeHash(t *testing.T) {
	p1 := uint8(12)
	x := uint64(0xffffffffffffffff)
	result := EncodeHash(x, p1)
	ideal := uint32(0xffffffff - 1)
	assert.Equal(t, result, ideal, "Encoded Incorrectly")
}

func TestDecodeHash(t *testing.T) {
	p1 := uint8(12)
	x := uint32(0xffffffff - 1)
	index, rho := DecodeHash(x, p1)

	assert.Equal(t, rho, uint8(0), "Did not decode rho properly")
	assert.Equal(t, index, uint32(0xfff), "Did not decode index properly")

	x = uint32(0xffffff00)
	index, rho = DecodeHash(x, p1)
	assert.Equal(t, rho, uint8(0), "Did not decode rho properly")
	assert.Equal(t, index, uint32(0xfff), "Did not decode index properly")
}

func TestEncodeDecode1(t *testing.T) {
	p1 := uint8(12)
	// construct number with index = 0f0 and rho = 4
	x := uint64(0x0f00ffffffffffff)

	encoded := EncodeHash(x, p1)
	index, rho := DecodeHash(encoded, p1)

	assert.Equal(t, index, uint32(0x0f0), "Incorrect index")
	assert.Equal(t, rho, uint8(4), "Incorrect rho")
}

func TestEncodeDecode2(t *testing.T) {
	p1 := uint8(12)
	// construct number with index = 0f0 and rho = 16
	x := uint64(0x0f00000f00000000)

	encoded := EncodeHash(x, p1)
	index, rho := DecodeHash(encoded, p1)

	assert.Equal(t, index, uint32(0x0f0), "Incorrect index")
	assert.Equal(t, rho, uint8(16), "Incorrect rho")
}

func TestEstimateBias(t *testing.T) {
	bias := EstimateBias(27.5, 5)
	actualBias := 17.4134

	if math.Abs(bias/actualBias-1) > 0.01 {
		t.Fatalf("Incorrect bias estimate.  Calculated %f, should be closer to %f", bias, actualBias)
	}
}
