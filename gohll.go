package gohll

//**
// HLL++ Implemintation by Micha Gorelick
// paper -- http://im.micha.gd/1dc0z0S
//**

import (
	"errors"
	"github.com/reusee/mmh3"
	"math"
)

const (
	SPARSE byte = iota
	NORMAL
)

var (
	InvalidPError        = errors.New("Invalid value of P, must be 4<=p<=25")
	SamePError           = errors.New("Both HLL instances must have the same value of P")
	ErrorRateOutOfBounds = errors.New("Error rate must be 0.26>=errorRate>=0.00025390625")
)

func MMH3Hash(value string) uint64 {
	hashBytes := mmh3.Hash128([]byte(value))
	var hash uint64
	for i, value := range hashBytes {
		hash |= uint64(value) << uint(i*8)
	}
	return hash
}

type HLL struct {
	P uint8

	Hasher func(string) uint64

	m1 uint
	m2 uint

	alpha  float64
	format byte

	tempSet    *TempSet
	sparseList *SparseList

	registers []uint8
}

func NewHLLByError(errorRate float64) (*HLL, error) {
	if errorRate < 0.00025390625 || errorRate > 0.26 {
		return nil, ErrorRateOutOfBounds
	}
	p := uint8(math.Ceil(math.Log2(math.Pow(1.04/errorRate, 2))))
	return NewHLL(p)
}

func NewHLL(p uint8) (*HLL, error) {
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

	// Since HLL.registers is a uint8 slice and the SparseList is a uint32
	// slice, we switch from sparse to normal with the sparse list is |m1/4| in
	// size (ie: the same size as the registers would be.
	sparseList := NewSparseList(p, int(m1/4))
	tempSet := make(TempSet, 0, int(m1/8))

	return &HLL{
		P:          p,
		Hasher:     MMH3Hash,
		m1:         m1,
		m2:         m2,
		alpha:      alpha,
		format:     format,
		tempSet:    &tempSet,
		sparseList: sparseList,
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
	w := SliceUint64(hash, 63-h.P, 0) << h.P
	rho := LeadingBitUint64(w) + 1
	if h.registers[index] < rho {
		h.registers[index] = rho
	}
}

func (h *HLL) addSparse(hash uint64) {
	k := EncodeHash(hash, h.P)
	h.tempSet = h.tempSet.Append(k)
	if h.tempSet.Full() {
		h.mergeSparse()
	}
	h.checkModeChange()
}

func (h *HLL) mergeSparse() {
	h.sparseList.Merge(h.tempSet)
	h.tempSet.Clear()
}

func (h *HLL) checkModeChange() {
	if h.sparseList.Full() {
		h.toNormal()
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
	for _, value := range *(h.tempSet) {
		index, rho := DecodeHash(value, h.P)
		if h.registers[index] < rho {
			h.registers[index] = rho
		}
	}
	h.tempSet.Clear()
	h.sparseList.Clear()
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
	var V int
	Ebottom := 0.0
	for _, value := range h.registers {
		Ebottom += math.Pow(2, -1.0*float64(value))
		if value == 0 {
			V += 1
		}
	}

	return h.cardinalityNormalCorrected(Ebottom, V)
}

func (h *HLL) cardinalityNormalCorrected(Ebottom float64, V int) float64 {
	E := h.alpha * float64(h.m1*h.m1) / Ebottom
	var Eprime float64
	if E < 5*float64(h.m1) {
		Eprime = E - EstimateBias(E, h.P)
	} else {
		Eprime = E
	}

	var H float64
	if V != 0 {
		H = LinearCounting(h.m1, V)
	} else {
		H = Eprime
	}

	if H <= Threshold(h.P) {
		return H
	}
	return Eprime
}

func (h *HLL) cardinalitySparse() float64 {
	h.mergeSparse()
	return LinearCounting(h.m2, int(h.m2)-h.sparseList.Len())
}

func (h *HLL) Union(other *HLL) error {
	if h.P != other.P {
		return SamePError
	}
	if other.format == NORMAL {
		if h.format == SPARSE {
			h.toNormal()
		}
		for i := uint(0); i < h.m1; i++ {
			if other.registers[i] > h.registers[i] {
				h.registers[i] = other.registers[i]
			}
		}
	} else if h.format == NORMAL && other.format == SPARSE {
		other.mergeSparse()
		for _, value := range other.sparseList.Data {
			index, rho := DecodeHash(value, h.P)
			if h.registers[index] < rho {
				h.registers[index] = rho
			}
		}
	} else if h.format == SPARSE && other.format == SPARSE {
		h.mergeSparse()
		other.mergeSparse()
		h.sparseList.Merge(other.sparseList)
		h.checkModeChange()
	}
	return nil
}

func (h *HLL) CardinalityIntersection(other *HLL) (float64, error) {
	if h.P != other.P {
		return 0.0, SamePError
	}
	A := h.Cardinality()
	B := other.Cardinality()
	AuB, _ := h.CardinalityUnion(other)
	return A + B - AuB, nil
}

func (h *HLL) CardinalityUnion(other *HLL) (float64, error) {
	if h.P != other.P {
		return 0.0, SamePError
	}
	cardinality := 0.0
	if h.format == NORMAL && other.format == NORMAL {
		cardinality = h.cardinalityUnionNN(other)
	} else if h.format == NORMAL && other.format == SPARSE {
		cardinality = h.cardinalityUnionNS(other)
	} else if h.format == SPARSE && other.format == NORMAL {
		cardinality, _ = other.CardinalityUnion(h)
	} else if h.format == SPARSE && other.format == SPARSE {
		cardinality = h.cardinalityUnionSS(other)
	}
	return cardinality, nil
}

func (h *HLL) cardinalityUnionNN(other *HLL) float64 {
	var V int
	Ebottom := 0.0
	for i, value := range h.registers {
		if other.registers[i] > value {
			value = other.registers[i]
		}
		Ebottom += math.Pow(2, -1.0*float64(value))
		if value == 0 {
			V += 1
		}
	}
	return h.cardinalityNormalCorrected(Ebottom, V)
}

func (h *HLL) cardinalityUnionNS(other *HLL) float64 {
	var V int
	other.mergeSparse()
	registerOther := make([]uint8, h.m1)
	for _, value := range other.sparseList.Data {
		index, rho := DecodeHash(value, other.P)
		if registerOther[index] < rho {
			registerOther[index] = rho
		}
	}
	Ebottom := 0.0
	for i, value := range h.registers {
		if registerOther[i] > value {
			value = registerOther[i]
		}
		Ebottom += math.Pow(2, -1.0*float64(value))
		if value == 0 {
			V += 1
		}
	}
	registerOther = registerOther[:0]
	return h.cardinalityNormalCorrected(Ebottom, V)
}

func (h *HLL) cardinalityUnionSS(other *HLL) float64 {
	h.mergeSparse()
	other.mergeSparse()
	if h.sparseList.Len() == 0 {
		return other.Cardinality()
	} else if other.sparseList.Len() == 0 {
		return h.Cardinality()
	}
	var i, j, V int
	var idxH, idxOther uint32
	for i < h.sparseList.Len()-1 || j < other.sparseList.Len()-1 {
		if i < h.sparseList.Len() {
			idxH = GetIndexSparse(h.sparseList.Get(i))
		}
		if j < other.sparseList.Len() {
			idxOther = GetIndexSparse(other.sparseList.Get(j))
		}
		V += 1
		if idxH < idxOther {
			i += 1
		} else if idxH > idxOther {
			j += 1
		} else {
			i += 1
			j += 1
		}
	}
	return LinearCounting(h.m2, int(h.m2)-V)
}
