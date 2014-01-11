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
	return float64(m1) * math.Log(float64(m1)/float64(V))
}

func EstimateBias(E float64, p uint8) float64 {
	return 0.0
}

func Threshold(p uint8) float64 {
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
