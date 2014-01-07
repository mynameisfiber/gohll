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
    P1 uint8
    P2 uint8

    m1 uint
    m2 uint

    alpha float64
    format byte
    
    tmpSet map[uint64]bool
    sparseList []uint
    maxSparseSetSize uint

    registers []uint8
}

func NewHLL(p1, p2 uint8, maxSparseSetSize uint) *HLL {
    m1 := uint(1 << p1)
    m2 := uint(1 << p2)

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

    tmpSet := make(map[uint]bool)
    sparseList := make([]uint, 0, m1 * 6)
    maxSparseSetSize = maxSparseSetSize

    return &HLL{
        P1: p1,
        P2: p2,
        m1: m1,
        m2: m2,
        alpha: alpha,
        format: format,
        tmpSet: tmpSet,
        sparseList: sparseList,
        maxSparseSetSize: maxSparseSetSize,
    }
}

func (h *HLL) Add(value string) {
    x, _ := Uvarint(mmh3.Hash128(value))
    switch h.format {
        case NORMAL:
            continue // IMPLEMENT
        case SPARSE:
            h.addSparse(x)
    }
}

func (h *HLL) addSparse(hash uint64) {
    k := EncodeHash(hash, h.P1, h.P2)
    h.tmpSet[k] = true
    if len(h.tmpSet) > h.maxSparseSetSize {
        Merge(&h.sparseList, h.tmpSet)
    }
}
