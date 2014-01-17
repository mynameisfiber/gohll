package gohll

import ()

type tempSet []uint32

func (ts tempSet) Len() int {
	return len(ts)
}

func (ts tempSet) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts tempSet) Less(i, j int) bool {
	indexI := getIndexSparse(ts[i])
	indexJ := getIndexSparse(ts[j])

	if indexI < indexJ {
		return true
	} else if indexI > indexJ {
		return false
	}
	// If indexI == indexJ we do a reverse sort on the rho values so we can
	// easily find the largest rho value for the same index
	return ts[i] > ts[j]
}

func (ts tempSet) Get(i int) uint32 {
	return ts[i]
}

func (ts tempSet) Clear() {
	ts = ts[0:0]
}

func (ts tempSet) Append(value uint32) *tempSet {
	newTs := append(ts, value)
	return &newTs
}

func (ts tempSet) Full() bool {
	return len(ts) == cap(ts)
}
