package gohll

import (
    "math"
)

func SliceInt(N uint64, start, stop uint8) uint64 {
    // This function will take and integer and indicies into the bits and
    // return the number represended within that range.  The leading bit of an
    // integer is at index 64 while the trailing bit is at index 0.  The
    // numerical value of start is larger than that of stop.
    //
    // Example in 32bit:
    // start, stop = 6, 1
    // N      = 11110000111100001111000011110000
    // mask   = 00000000000000000000000001111110
    // result = 00000000000000000000000000111000
    mask := uint64((1<<(start+1)) - (1<<stop)
    r := N & mask
    if stop > 0 {
        r >>= stop
    }
    return r
}

func LeadingBit(N uint64) uint8 {
	// Returns the index of the leftmost set bit
	return 64 - uint8(math.Log2(float64(N))+1)
}

func EncodeHash(x uint64, p1, p2 uint8) uint64 {
    if SliceInt(x, 63-p1, 64-p2) == 0 {
        result := SliceInt(x, 63, 64-p2) << 7
        result |= uint64(LeadingBit(x) << 1)
        result += 1
        return result
    } else {
        return SliceInt(x, 63, 64-p2) << 1
    }
}

func DecodeHash(x uint64, p1, p2 uint8) (uint64, uint8) {
    var r uint8
    if x & 1 == 1 {
        r = uint8(SliceInt(x, 6, 1)) + (p2 - p1)
    } else {
        r = LeadingBit(SliceInt(x, p2-p1-1, 1))
    }
    return GetIndex(x, p1), r

}

func GetIndex(x uint64, p1 uint8) uint64 {
    if x & 1 == 1 {
        return SliceInt(x, p1+6, 6)
    } else {
        return SliceInt(x, p1+1, 1)
    }
}
