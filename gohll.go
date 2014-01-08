package gohll

import (
    "encoding/binary"
    "github.com/reusee/mmh3"
)

const (
    SPARSE byte = iota
    NORMAL
)

type HLL struct {
    P uint8

    m1 uint
    m2 uint

    alpha float64
    format byte
    
    tempSet *SparseList
    sparseList *SparseList
    MaxSparseSetSize int

    registers []uint8
}

func NewHLL(p uint8, maxSparseSetSize int) *HLL {
    m1 := uint(1 << p)
    m2 := uint(1 << 25)

    var alpha float64
    switch m1 {
        case 16:
            alpha= 0.673
        case 32:
            alpha= 0.697
        case 64:
            alpha = 0.709
        default:
            alpha = 0.7213/(1 + 1.079/float64(m1))
    }

    format := SPARSE

    tempSet := NewSparseList(p, maxSparseSetSize)
    sparseList := NewSparseList(p, int(m1 * 6))
    maxSparseSetSize = maxSparseSetSize

    return &HLL{
        P: p,
        m1: m1,
        m2: m2,
        alpha: alpha,
        format: format,
        tempSet: tempSet,
        sparseList: sparseList,
        MaxSparseSetSize: maxSparseSetSize,
    }
}

func (h *HLL) Add(value string) {
    x, _ := binary.Uvarint(mmh3.Hash128([]byte(value)))
    switch h.format {
        case NORMAL:
            h.addNormal(x)
        case SPARSE:
            h.addSparse(x)
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
