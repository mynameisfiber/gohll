package gohll

import (
)

type TempSet []uint32

func (ts TempSet) Len() int {
    return len(ts)
}

func (ts TempSet) Swap(i, j int) {
    ts[i], ts[j] = ts[j], ts[i]
}

func (ts TempSet) Less(i, j int) bool {
	indexI := GetIndexSparse(ts[i])
	indexJ := GetIndexSparse(ts[j])

	if indexI < indexJ {
		return true
	} else if indexI > indexJ {
		return false
	} else {
		// If indexI == indexJ we do a reverse sort on the rho values so we can
		// easily find the largest rho value for the same index
		return ts[i] > ts[j]
	}
}

func (ts TempSet) Get(i int) uint32 {
    return ts[i]
}

func (ts TempSet) Clear() {
    ts = ts[0:0]
}

func (ts TempSet) Append(value uint32) *TempSet {
    newTs := append(ts, value)
    return &newTs
}

func (ts TempSet) Full() bool {
    return len(ts) == cap(ts)
}
