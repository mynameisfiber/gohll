package gohll

import (
	"math"
)

// Takes in a 64bit hash and the set precision and outputs a 32bit encoded hash
// for use with the SparseList
func EncodeHash(x uint64, p uint8) uint32 {
	if SliceUint64(x, 63-p, 39) == 0 {
		var result uint32
		result = uint32((x >> 32) &^ 0x7f)
		w := SliceUint64(x, 63-p, 0) << p
		result |= (uint32(LeadingBitUint64(w)) << 1)
		result |= 1
		return result
	}
	return uint32(x>>32) &^ 0x1
}

// Takes a 32bit hash which was encoded for use with the SparseList and
// extracts the meaningful metadata from it using the normal mode precision
// (namely it's index and location of it's leading set bit)
func DecodeHash(x uint32, p uint8) (uint32, uint8) {
	var r uint8
	if x&0x1 == 1 {
		r = uint8(SliceUint32(x, 6, 1))
	} else {
		r = LeadingBitUint32(SliceUint32(x, 31-p, 1) << (1 + p))
	}
	return GetIndex(x, p), r + 1

}

// Returns the normal mode precision (given by p) of an encoded hash
func GetIndex(x uint32, p uint8) uint32 {
	return SliceUint32(x, 31, 32-p)
}

// Returns the sparse mode index of the encoded hash
func GetIndexSparse(x uint32) uint32 {
	return x >> 7
}

// Performs linear counting given the number of registers, m1, and the number
// of empty registers, V
func LinearCounting(m1 uint, V int) float64 {
	return float64(m1) * math.Log(float64(m1)/float64(V))
}

// Estimates the amount of bias in a normal mode cardinality query with an
// estimator value of E and a normal mode precision of p
func EstimateBias(E float64, p uint8) float64 {
	if p > 18 {
		return 0.0
	}
	estimateVector := rawEstimateData[p-4]
	N := len(estimateVector)
	if E < estimateVector[0] || E > estimateVector[N-1] {
		return 0.0
	}

	biasVector := biasData[p-4]

	for i, v := range estimateVector[1:] {
		if v == E {
			return biasVector[i]
		}
		if v > E && estimateVector[i-1] < E {
			return linearInterpolation(estimateVector[i-1:i+1], biasVector[i-1:i+1], E)
		}
	}
	return 0.0
}

func linearInterpolation(x, y []float64, x0 float64) float64 {
	if len(x) != 2 || len(y) != 2 {
		return 0.0
	}
	return y[0] + (y[1]-y[0])*(x0-x[0])/(x[1]-x[0])
}

func threshold(p uint8) float64 {
	switch p {
	case 4:
		return 10
	case 5:
		return 20
	case 6:
		return 40
	case 8:
		return 220
	case 7:
		return 80
	case 9:
		return 400
	case 10:
		return 900
	case 11:
		return 1800
	case 12:
		return 3100
	case 13:
		return 6500
	case 14:
		return 11500
	case 15:
		return 20000
	case 16:
		return 50000
	case 17:
		return 120000
	case 18:
		return 350000
	}
	return 0
}
