package gohll

//**
// HLL++ Implemintation by Micha Gorelick
// paper -- http://im.micha.gd/1dc0z0S
//**

import (
	"errors"
	"github.com/mynameisfiber/gohll/mmh3"
	"math"
)

const (
	SPARSE byte = iota
	NORMAL
)

var (
	ErrInvalidP             = errors.New("invalid value of P, must be 4<=p<=25")
	ErrSameP                = errors.New("both HLL instances must have the same value of P")
	ErrErrorRateOutOfBounds = errors.New("error rate must be 0.26>=errorRate>=0.00025390625")
)

// MMH3Hash is the default hasher and uses murmurhash to return a uint64
func MMH3Hash(value string) uint64 {
	h1, _ := mmh3.Hash128(value)
	return h1
}

type HLL struct {
	P uint8

	Hasher func(string) uint64

	m1 uint
	m2 uint

	alpha  float64
	format byte

	tempSet    *tempSet
	sparseList *sparseList

	registers []uint8
}

// NewHLLByError creates a new HLL object with error rate given by `errorRate`.
// The error must be between 26% and 0.0253%
func NewHLLByError(errorRate float64) (*HLL, error) {
	if errorRate < 0.00025390625 || errorRate > 0.26 {
		return nil, ErrErrorRateOutOfBounds
	}
	p := uint8(math.Ceil(math.Log2(math.Pow(1.04/errorRate, 2))))
	return NewHLL(p)
}

// NewHLL creates a new HLL object given a normal mode precision between 4 and
// 25
func NewHLL(p uint8) (*HLL, error) {
	if p < 4 || p > 25 {
		return nil, ErrInvalidP
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
	sparseList := newSparseList(p, int(m1/4))
	tempSet := make(tempSet, 0, int(m1/16))

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

// Add will add the given string value to the HLL using the currently set
// Hasher function
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
	index := sliceUint64(hash, 63, 64-h.P)
	w := sliceUint64(hash, 63-h.P, 0) << h.P
	rho := leadingBitUint64(w) + 1
	if h.registers[index] < rho {
		h.registers[index] = rho
	}
}

func (h *HLL) addSparse(hash uint64) {
	k := encodeHash(hash, h.P)
	h.tempSet = h.tempSet.Append(k)
	if h.tempSet.Full() {
		h.mergeSparse()
		h.checkModeChange()
	}
}

func (h *HLL) mergeSparse() {
	h.sparseList.Merge(h.tempSet)
	h.tempSet.Clear()
}

func (h *HLL) checkModeChange() {
	if h.sparseList.Full() {
		h.ToNormal()
	}
}

// ToNormal will convert the current HLL to normal mode, maintaining any data
// already inserted into the structure, if it is in sparse mode
func (h *HLL) ToNormal() {
	if h.format != SPARSE {
		return
	}
	h.format = NORMAL
	h.registers = make([]uint8, h.m1)
	for _, value := range h.sparseList.Data {
		index, rho := decodeHash(value, h.P)
		if h.registers[index] < rho {
			h.registers[index] = rho
		}
	}
	for _, value := range *(h.tempSet) {
		index, rho := decodeHash(value, h.P)
		if h.registers[index] < rho {
			h.registers[index] = rho
		}
	}
	h.tempSet.Clear()
	h.sparseList.Clear()
}

// Cardinality returns the estimated cardinality of the current HLL object
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
		Eprime = E - estimateBias(E, h.P)
	} else {
		Eprime = E
	}

	var H float64
	if V != 0 {
		H = linearCounting(h.m1, V)
	} else {
		H = Eprime
	}

	if H <= threshold(h.P) {
		return H
	}
	return Eprime
}

func (h *HLL) cardinalitySparse() float64 {
	h.mergeSparse()
	return linearCounting(h.m2, int(h.m2)-h.sparseList.Len())
}

// Union will merge all data in another HLL object into this one.
func (h *HLL) Union(other *HLL) error {
	if h.P != other.P {
		return ErrSameP
	}
	if other.format == NORMAL {
		if h.format == SPARSE {
			h.ToNormal()
		}
		for i := uint(0); i < h.m1; i++ {
			if other.registers[i] > h.registers[i] {
				h.registers[i] = other.registers[i]
			}
		}
	} else if h.format == NORMAL && other.format == SPARSE {
		other.mergeSparse()
		for _, value := range other.sparseList.Data {
			index, rho := decodeHash(value, h.P)
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

// CardinalityIntersection returns the estimated cardinality of the
// intersection between this HLL object and another one.  That is, it returns
// an estimate of the number of unique items that occur in both this and the
// other HLL object.  This is done with the Inclusion–exclusion principle and
// does not satisfy the error guarantee.
func (h *HLL) CardinalityIntersection(other *HLL) (float64, error) {
	if h.P != other.P {
		return 0.0, ErrSameP
	}
	A := h.Cardinality()
	B := other.Cardinality()
	AuB, _ := h.CardinalityUnion(other)
	return A + B - AuB, nil
}

// CardinalityUnion returns the estimated cardinality of the union between this
// and another HLL object.  This result would be the same as first taking the
// union between this and the other object and then calling Cardinality.
// However, by calling this function we are not making any changes to the HLL
// object.
func (h *HLL) CardinalityUnion(other *HLL) (float64, error) {
	if h.P != other.P {
		return 0.0, ErrSameP
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
		index, rho := decodeHash(value, other.P)
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
			idxH = getIndexSparse(h.sparseList.Get(i))
		}
		if j < other.sparseList.Len() {
			idxOther = getIndexSparse(other.sparseList.Get(j))
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
	return linearCounting(h.m2, int(h.m2)-V)
}
