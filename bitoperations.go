package gohll

// SliceUint32 takes a 32bit integer and indicies into the bits and return the
// number represended within that range.  The leading bit of an integer is at
// index 32 while the trailing bit is at index 0.  The numerical value of start
// is larger than that of stop.
//
// Example in 32bit:
// start, stop = 6, 1
// N      = 11110000111100001111000011110000
// mask   = 00000000000000000000000001111110
// result = 00000000000000000000000000111000
func sliceUint32(N uint32, start, stop uint8) uint32 {
	mask := uint32((1 << (start + 1)) - (1 << stop))
	r := N & mask
	if stop > 0 {
		r >>= stop
	}
	return r
}

// SliceUint64 takes a 64bit  integer and indicies into the bits and return the
// number represended within that range.  The leading bit of an integer is at
// index 64 while the trailing bit is at index 0.  The numerical value of start
// is larger than that of stop.
//
// Example in 32bit:
// start, stop = 6, 1
// N      = 11110000111100001111000011110000
// mask   = 00000000000000000000000001111110
// result = 00000000000000000000000000111000
func sliceUint64(N uint64, start, stop uint8) uint64 {
	mask := uint64((1 << (start + 1)) - (1 << stop))
	r := N & mask
	if stop > 0 {
		r >>= stop
	}
	return r
}
