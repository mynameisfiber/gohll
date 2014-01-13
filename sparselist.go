package gohll

import (
	"sort"
)

type MergableList interface {
	sort.Interface
	Get(int) uint32
}

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
	indexI := GetIndexSparse(sl.Data[i])
	indexJ := GetIndexSparse(sl.Data[j])

	if indexI < indexJ {
		return true
	} else if indexI > indexJ {
		return false
	}
	// If indexI == indexJ we do a reverse sort on the rho values so we can
	// easily find the largest rho value for the same index
	return sl.Data[i] > sl.Data[j]
}

func (sl *SparseList) Add(N uint32) {
	sl.Data = append(sl.Data, N)
}

func (sl *SparseList) Swap(i, j int) {
	sl.Data[i], sl.Data[j] = sl.Data[j], sl.Data[i]
}

func (sl *SparseList) Get(i int) uint32 {
	return sl.Data[i]
}

func (sl *SparseList) Clear() {
	sl.Data = sl.Data[0:0]
}

func (sl *SparseList) Merge(tmpList MergableList) {
	// This function assumes that sl is already sorted!
	if tmpList.Len() == 0 {
		return
	}
	sort.Sort(tmpList)

	var slIndex uint32
	var slStopIteration bool
	sli := int(0)
	if sl.Len() > 0 {
		slIndex = GetIndexSparse(sl.Data[0])
		slStopIteration = false
	} else {
		slStopIteration = true
	}

	slDirty := false
	var lastTmpIndex uint32
	var value uint32
	for i := 0; i < tmpList.Len(); i++ {
		value = tmpList.Get(i)
		tmpIndex := GetIndexSparse(value)
		if tmpIndex == lastTmpIndex && i != 0 {
			continue
		}
		if tmpIndex > slIndex {
			for tmpIndex > slIndex {
				sli += 1
				if sli >= sl.Len() {
					slStopIteration = true
					break
				}
				slIndex = GetIndexSparse(sl.Data[sli])
			}
		}
		if slStopIteration || tmpIndex < slIndex {
			sl.Add(value)
			slDirty = true
		} else if tmpIndex == slIndex {
			if value > sl.Data[sli] {
				slDirty = true
				sl.Data[sli] = value
			}
		}
		lastTmpIndex = tmpIndex
	}
	if slDirty {
		sort.Sort(sl)
	}
}
