package gohll

import (
	"math"
	"math/bits"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeHash(t *testing.T) {
	p1 := uint8(12)
	x := uint64(0xffffffffffffffff)
	result := encodeHash(x, p1)
	ideal := uint32(0xffffffff - 1)
	assert.Equal(t, result, ideal, "Encoded Incorrectly")
}

func TestDecodeHash(t *testing.T) {
	p1 := uint8(12)
	x := uint32(0xffffffff - 1)
	index, rho := decodeHash(x, p1)

	assert.Equal(t, rho, uint8(1), "Did not decode rho properly")
	assert.Equal(t, index, uint32(0xfff), "Did not decode index properly")

	x = uint32(0xffffff00)
	index, rho = decodeHash(x, p1)
	assert.Equal(t, rho, uint8(1), "Did not decode rho properly")
	assert.Equal(t, index, uint32(0xfff), "Did not decode index properly")
}

func TestEncodeDecode1(t *testing.T) {
	p1 := uint8(12)
	// construct number with index = 0f0 and rho = 4
	x := uint64(0x0f00ffffffffffff)

	encoded := encodeHash(x, p1)
	index, rho := decodeHash(encoded, p1)

	assert.Equal(t, index, uint32(0x0f0), "Incorrect index")
	assert.Equal(t, rho, uint8(4)+1, "Incorrect rho")
}

func TestEncodeDecode2(t *testing.T) {
	p1 := uint8(12)
	// construct number with index = 0f0 and rho = 16
	x := uint64(0x0f00000f00000000)

	encoded := encodeHash(x, p1)
	index, rho := decodeHash(encoded, p1)

	assert.Equal(t, index, uint32(0x0f0), "Incorrect index")
	assert.Equal(t, rho, uint8(16)+1, "Incorrect rho")
}

func TestEncodeDecode3(t *testing.T) {
	p := uint8(4)
	var hash uint64
	for i := 0; i < 100; i++ {
		hash = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())

		index := sliceUint64(hash, 63, 64-p)
		w := sliceUint64(hash, 63-p, 0) << p
		rho := bits.LeadingZeros64(w) + 1

		e := encodeHash(hash, p)
		edIndex, edRho := decodeHash(e, p)

		assert.Equal(t, uint64(edIndex), uint64(index), "Incorrect index")
		assert.Equal(t, uint64(edRho), uint64(rho), "Incorrect index")
	}
}

func TestEstimateBias(t *testing.T) {
	bias := estimateBias(27.5, 5)
	actualBias := 17.4134

	if math.Abs(bias/actualBias-1) > 0.01 {
		t.Fatalf("Incorrect bias estimate.  Calculated %f, should be closer to %f", bias, actualBias)
	}
}

func TestEstimateBias2(t *testing.T) {
	bias := estimateBias(11822.412839663843, 14)
	actualBias := 11811.188669

	if math.Abs(bias/actualBias-1) > 0.01 {
		t.Fatalf("Incorrect bias estimate.  Calculated %f, should be closer to %f", bias, actualBias)
	}
}
