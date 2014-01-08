package gohll

func EncodeHash(x uint64, p uint8) uint32 {
    if SliceUint64(x, 63-p, 39) == 0 {
        var result uint32
        result = uint32(x >> 39) << 7
        result |= uint32(LeadingBitUint64(x) << 1)
        result += 1
        return result
    } else {
        return uint32(x >> 39) << 1
    }
}

func DecodeHash(x uint32, p uint8) (uint32, uint8) {
    var r uint8
    if x & 1 == 1 {
        r = uint8(SliceUint32(x, 6, 1)) + (25 - p)
    } else {
        r = LeadingBitUint32(SliceUint32(x, 24-p, 1))
    }
    return GetIndex(x, p), r

}

func GetIndex(x uint32, p uint8) uint32 {
    if x & 1 == 1 {
        return SliceUint32(x, p+6, 6)
    } else {
        return SliceUint32(x, p+1, 1)
    }
}
