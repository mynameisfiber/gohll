package gohll

import (
	"sort"
)

// Interface defining what objects are mergable with the sparseList object.
// Note: it is assumed that this list is sorted in the same way as the
// sparseList
type mergableList interface {
	sort.Interface
	Get(int) uint32
}

type sparseList struct {
	Data    []uint32
	P       uint8
	MaxSize int
}

func newSparseList(p uint8, capacity int) *sparseList {
	return &sparseList{
		Data:    make([]uint32, 0),
		P:       p,
		MaxSize: capacity,
	}
}

func (sl *sparseList) Len() int {
	return len(sl.Data)
}

func (sl *sparseList) Full() bool {
	return len(sl.Data) >= sl.MaxSize
}

func (sl *sparseList) Less(i, j int) bool {
	indexI := getIndexSparse(sl.Data[i])
	indexJ := getIndexSparse(sl.Data[j])

	if indexI < indexJ {
		return true
	} else if indexI > indexJ {
		return false
	}
	// If indexI == indexJ we do a reverse sort on the rho values so we can
	// easily find the largest rho value for the same index
	return sl.Data[i] > sl.Data[j]
}

func (sl *sparseList) Add(N uint32) {
	sl.Data = append(sl.Data, N)
}

func (sl *sparseList) Swap(i, j int) {
	sl.Data[i], sl.Data[j] = sl.Data[j], sl.Data[i]
}

func (sl *sparseList) Get(i int) uint32 {
	return sl.Data[i]
}

func (sl *sparseList) Clear() {
	sl.Data = sl.Data[0:0]
}

// Merge will merge this sparse list with another mergable list.  This is done
// by having the 32bit integers within the list sorted by it's encoded index
// and, if another item with the same index exists, only keeping the one with
// the largest number of leading zero bits (as given by leadingBitUint32).
//
// NOTE: This function assumes that this list is already sorted with the given
// Less() function
func (sl *sparseList) Merge(tmpList mergableList) {
	// This function assumes that sl is already sorted!
	if tmpList.Len() == 0 {
		return
	}
	sort.Sort(tmpList)

	var slIndex uint32
	var slStopIteration bool
	sli := int(0)
	if sl.Len() > 0 {
		slIndex = getIndexSparse(sl.Data[0])
		slStopIteration = false
	} else {
		slStopIteration = true
	}

	slDirty := false
	var lastTmpIndex uint32
	var value uint32
	for i := 0; i < tmpList.Len(); i++ {
		value = tmpList.Get(i)
		tmpIndex := getIndexSparse(value)
		if tmpIndex == lastTmpIndex && i != 0 {
			continue
		}
		if !slStopIteration && tmpIndex > slIndex {
			for tmpIndex > slIndex {
				sli++
				if sli >= sl.Len() {
					slStopIteration = true
					break
				}
				slIndex = getIndexSparse(sl.Data[sli])
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
