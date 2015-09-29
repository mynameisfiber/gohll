package gohll

import (
	"bytes"
	"encoding/gob"
)

type serializable struct {
	P uint8

	M1 uint
	M2 uint

	Alpha  float64
	Format byte

	TempSet    *tempSet
	SparseList *sparseList

	Registers []uint8
}

// MarshalBinary implements encoding.BinaryMarshaler.
// Does not serialize hasher!
func (h *HLL) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(
		serializable{
			P:          h.P,
			M1:         h.m1,
			M2:         h.m2,
			Alpha:      h.alpha,
			Format:     h.format,
			TempSet:    h.tempSet,
			SparseList: h.sparseList,
			Registers:  h.registers,
		})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// Preserves the hasher.
func (h *HLL) UnmarshalBinary(data []byte) error {
	var s serializable
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&s)
	if err != nil {
		return err
	}
	h.P = s.P
	h.m1 = s.M1
	h.m2 = s.M2
	h.alpha = s.Alpha
	h.format = s.Format
	h.tempSet = s.TempSet
	h.sparseList = s.SparseList
	h.registers = s.Registers

	if h.Hasher == nil {
		h.Hasher = MMH3Hash
	}
	return nil
}
