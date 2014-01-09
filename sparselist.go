package gohll

import (
	"sort"
)

type SparseList struct {
	Data    []uint32
	P       uint8
	MaxSize int
}

func NewSparseList(p uint8, capacity int) *SparseList {
	return &SparseList{
		Data:    make([]uint32, 0, capacity),
		P:       p,
		MaxSize: capacity,
	}
}

func (sl *SparseList) Len() int {
	return len(sl.Data)
}

func (sl *SparseList) Full() bool {
	return len(sl.Data) >= sl.MaxSize
}

func (sl *SparseList) Less(i, j int) bool {
	indexI, rhoI := DecodeHash(sl.Data[i], sl.P)
	indexJ, rhoJ := DecodeHash(sl.Data[j], sl.P)

	if indexI < indexJ {
		return true
	} else if indexI > indexJ {
		return false
	} else {
		// If indexI == indexJ we do a reverse sort on the rho values so we can
		// easily find the largest rho value for the same index
		return rhoI > rhoJ
	}
}

func (sl *SparseList) Add(N uint32) {
	sl.Data = append(sl.Data, N)
}

func (sl *SparseList) Swap(i, j int) {
	sl.Data[i], sl.Data[j] = sl.Data[j], sl.Data[i]
}

func (sl *SparseList) Clear() {
	sl.Data = sl.Data[0:0]
}

func (sl *SparseList) Merge(tmpList *SparseList) {
	// This function assumes that sl is already sorted!
	var slIndex uint32
	var slRho uint8

	sort.Sort(tmpList)
	sli := int(0)
	if sl.Len() > 0 {
		slIndex, slRho = DecodeHash(sl.Data[0], sl.P)
	}

	var lastTmpIndex uint32
	slDirty := false
	for _, value := range tmpList.Data {
		tmpIndex, tmpRho := DecodeHash(value, tmpList.P)
		if tmpIndex == lastTmpIndex {
			continue
		} else if tmpIndex > slIndex {
			sl.Add(value)
			slDirty = true
			for slIndex < tmpIndex && sli+1 < sl.Len() {
				sli += 1
				slIndex, slRho = DecodeHash(sl.Data[sli], sl.P)
			}
		} else if tmpIndex < slIndex {
			sl.Add(value)
			slDirty = true
		} else if tmpIndex == slIndex {
			if tmpRho > slRho {
				slDirty = true
				sl.Data[sli] = value
			}
		}
		lastTmpIndex = tmpIndex
	}
	tmpList.Clear()
	if slDirty {
		sort.Sort(sl)
	}
}
