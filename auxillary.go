package gohll

import (
	"math"
)

func EncodeHash(x uint64, p uint8) uint32 {
	if SliceUint64(x, 63-p, 39) == 0 {
		var result uint32
		result = uint32(x >> 32 &^ 0x7f)
		w := SliceUint64(x, 63-p, 0) << p
		result |= uint32(LeadingBitUint64(w) << 1)
		result += 1
		return result
	} else {
		return uint32(x>>32) ^ 1
	}
}

func DecodeHash(x uint32, p uint8) (uint32, uint8) {
	var r uint8
	if x&1 == 1 {
		r = uint8(SliceUint32(x, 6, 1))
	} else {
		r = LeadingBitUint32(SliceUint32(x, 31-p, 1) << (1 + p))
	}
	return GetIndex(x, p), r

}

func GetIndex(x uint32, p uint8) uint32 {
	return SliceUint32(x, 31, 32-p)
}

func LinearCounting(m1 uint, V int) float64 {
	return float64(m1) * math.Log2(float64(m1)/float64(V))
}
