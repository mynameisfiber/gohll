package gohll

// This function will take a 32bit integer and indicies into the bits and
// return the number represended within that range.  The leading bit of an
// integer is at index 32 while the trailing bit is at index 0.  The numerical
// value of start is larger than that of stop.
//
// Example in 32bit:
// start, stop = 6, 1
// N      = 11110000111100001111000011110000
// mask   = 00000000000000000000000001111110
// result = 00000000000000000000000000111000
func SliceUint32(N uint32, start, stop uint8) uint32 {
	mask := uint32((1 << (start + 1)) - (1 << stop))
	r := N & mask
	if stop > 0 {
		r >>= stop
	}
	return r
}

// This function will take a 64bit  integer and indicies into the bits and
// return the number represended within that range.  The leading bit of an
// integer is at index 64 while the trailing bit is at index 0.  The numerical
// value of start is larger than that of stop.
//
// Example in 32bit:
// start, stop = 6, 1
// N      = 11110000111100001111000011110000
// mask   = 00000000000000000000000001111110
// result = 00000000000000000000000000111000
func SliceUint64(N uint64, start, stop uint8) uint64 {
	mask := uint64((1 << (start + 1)) - (1 << stop))
	r := N & mask
	if stop > 0 {
		r >>= stop
	}
	return r
}

// Given a 32bit integer, this function returns what index has the first
// non-zero bit with the leading bit being indexed at 0 and the last bit
// indexed at 32
//
// LeadingBitUint32(0xffffffff) = 0
// LeadingBitUint32(0x0fffffff) = 3
// LeadingBitUint32(0x00ffffff) = 7
// LeadingBitUint32(0x00000000) = 32
func LeadingBitUint32(N uint32) uint8 {
	if N == 0 {
		return 32
	}
	t := uint32(1 << 31)
	r := uint8(0)
	for (N & t) == 0 {
		t >>= 1
		r += 1
	}
	return r
}

// Given a 64bit integer, this function returns what index has the first
// non-zero bit with the leading bit being indexed at 0 and the last bit
// indexed at 64
func LeadingBitUint64(N uint64) uint8 {
	if N == 0 {
		return 64
	}
	t := uint64(1 << 63)
	r := uint8(0)
	for (N & t) == 0 {
		t >>= 1
		r += 1
	}
	return r
}
