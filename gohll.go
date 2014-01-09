package gohll

import (
	"errors"
	"fmt"
	"github.com/reusee/mmh3"
)

const (
	SPARSE byte = iota
	NORMAL
)

var (
	InvalidPError = errors.New("Invalid value of P, must be 4<=p<=25")
)

func MMH3Hash(value string) uint64 {
	hashBytes := mmh3.Hash128([]byte(value))
	var hash uint64
	for i, value := range hashBytes {
		hash |= uint64(value) << uint(i*8)
	}
	fmt.Printf("%0.64b\n", hash)
	return hash
}

type HLL struct {
	P uint8

	Hasher func(string) uint64

	m1 uint
	m2 uint

	alpha  float64
	format byte

	tempSet          *SparseList
	sparseList       *SparseList
	MaxSparseSetSize int

	registers []uint8
}

func NewHLL(p uint8, maxSparseSetSize int) (*HLL, error) {
	if p < 4 || p > 25 {
		return nil, InvalidPError
	}

	m1 := uint(1 << p)
	m2 := uint(1 << 25)

	var alpha float64
	switch m1 {
	case 16:
		alpha = 0.673
	case 32:
		alpha = 0.697
	case 64:
		alpha = 0.709
	default:
		alpha = 0.7213 / (1 + 1.079/float64(m1))
	}

	format := SPARSE

	tempSet := NewSparseList(p, maxSparseSetSize)
	sparseList := NewSparseList(p, int(m1*6))
	maxSparseSetSize = maxSparseSetSize

	return &HLL{
		P:                p,
		MaxSparseSetSize: maxSparseSetSize,
		Hasher:           MMH3Hash,
		m1:               m1,
		m2:               m2,
		alpha:            alpha,
		format:           format,
		tempSet:          tempSet,
		sparseList:       sparseList,
	}, nil
}

func (h *HLL) Add(value string) {
	hash := h.Hasher(value)
	switch h.format {
	case NORMAL:
		h.addNormal(hash)
	case SPARSE:
		h.addSparse(hash)
	}
}

func (h *HLL) addNormal(hash uint64) {
	index := SliceUint64(hash, 63, 64-h.P)
	w := SliceUint64(hash, 63-h.P, 0)
	rho := LeadingBitUint64(w)
	if h.registers[index] < rho {
		h.registers[index] = rho
	}
}

func (h *HLL) addSparse(hash uint64) {
	k := EncodeHash(hash, h.P)
	h.tempSet.Add(k)
	if h.tempSet.Full() {
		h.sparseList.Merge(h.tempSet)
		if h.sparseList.Full() {
			h.toNormal()
		}
	}
}

func (h *HLL) toNormal() {
	h.format = NORMAL
	h.registers = make([]uint8, h.m1)
	for _, value := range h.sparseList.Data {
		index, rho := DecodeHash(value, h.P)
		if h.registers[index] < rho {
			h.registers[index] = rho
		}
	}
}

func (h *HLL) Cardinality() float64 {
	var cardinality float64
	switch h.format {
	case NORMAL:
		cardinality = h.cardinalityNormal()
	case SPARSE:
		cardinality = h.cardinalitySparse()
	}
	return cardinality
}

func (h *HLL) cardinalityNormal() float64 {
	return 0.0
}

func (h *HLL) cardinalitySparse() float64 {
	h.sparseList.Merge(h.tempSet)
	return LinearCounting(h.m2, int(h.m2)-h.sparseList.Len())
}
