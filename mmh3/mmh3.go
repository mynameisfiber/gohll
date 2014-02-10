// MurmurHash implementation.
// Cloned from github.com/reusee/mmh3
package mmh3

import (
	"bytes"
	"encoding/binary"
)

const (
	c1_32 uint32 = 0xcc9e2d51
	c2_32 uint32 = 0x1b873593
)

func Hash32(s string) uint32 {
	length := len(s)
	if length == 0 {
		return 0
	}
	key := byteSlice(s)

	nblocks := length >> 2
	var h, k uint32
	buf := bytes.NewBuffer(key)
	for i := 0; i < nblocks; i++ {
		binary.Read(buf, binary.LittleEndian, &k)
		k *= c1_32
		k = (k << 15) | (k >> (32 - 15))
		k *= c2_32
		h ^= k
		h = (h << 13) | (h >> (32 - 13))
		h = (h << 2) + h + 0xe6546b64
	}
	k = 0
	tailIndex := nblocks << 2
	switch length & 3 {
	case 3:
		k ^= uint32(key[tailIndex+2]) << 16
		fallthrough
	case 2:
		k ^= uint32(key[tailIndex+1]) << 8
		fallthrough
	case 1:
		k ^= uint32(key[tailIndex])
		k *= c1_32
		k = (k << 15) | (k >> (32 - 15))
		k *= c2_32
		h ^= k
	}
	h ^= uint32(length)
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

const (
	c1_128 = 0x87c37b91114253d5
	c2_128 = 0x4cf5ad432745937f
)

func Hash128(s string) (uint64, uint64) {
	length := len(s)
	if length == 0 {
		return 0, 0
	}
	key := byteSlice(s)

	nblocks := length >> 4
	var h1, h2, k1, k2 uint64
	buf := bytes.NewBuffer(key)
	for i := 0; i < nblocks; i++ {
		binary.Read(buf, binary.LittleEndian, &k1)
		binary.Read(buf, binary.LittleEndian, &k2)
		k1 *= c1_128
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= c2_128
		h1 ^= k1
		h1 = (h1 << 27) | (h1 >> (64 - 27))
		h1 += h2
		h1 = (h1 << 2) + h1 + 0x52dce729
		k2 *= c2_128
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= c1_128
		h2 ^= k2
		h2 = (h2 << 31) | (h2 >> (64 - 31))
		h2 += h1
		h2 = (h2 << 2) + h2 + 0x38495ab5
	}
	k1, k2 = 0, 0
	tailIndex := nblocks << 4
	switch length & 15 {
	case 15:
		k2 ^= uint64(key[tailIndex+14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(key[tailIndex+13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(key[tailIndex+12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(key[tailIndex+11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(key[tailIndex+10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(key[tailIndex+9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(key[tailIndex+8])
		k2 *= c2_128
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= c1_128
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(key[tailIndex+7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(key[tailIndex+6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(key[tailIndex+5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(key[tailIndex+4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(key[tailIndex+3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(key[tailIndex+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(key[tailIndex+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(key[tailIndex])
		k1 *= c1_128
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= c2_128
		h1 ^= k1
	}
	h1 ^= uint64(length)
	h2 ^= uint64(length)
	h1 += h2
	h2 += h1
	h1 ^= h1 >> 33
	h1 *= 0xff51afd7ed558ccd
	h1 ^= h1 >> 33
	h1 *= 0xc4ceb9fe1a85ec53
	h1 ^= h1 >> 33
	h2 ^= h2 >> 33
	h2 *= 0xff51afd7ed558ccd
	h2 ^= h2 >> 33
	h2 *= 0xc4ceb9fe1a85ec53
	h2 ^= h2 >> 33
	h1 += h2
	h2 += h1

	return h1, h2
}
