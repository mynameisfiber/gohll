package mmh3

import (
	"reflect"
	"unsafe"
)

func byteSlice(s string) []byte {
	var b []byte
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	return b
}

func uint32Slice(s string) []uint32 {
	var slice []uint32
	pslice := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pslice.Data = pstring.Data
	pslice.Len = pstring.Len >> 2
	return slice
}

func uint64Slice(s string) []uint64 {
	var slice []uint64
	pslice := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pslice.Data = pstring.Data
	pslice.Len = pstring.Len >> 3
	return slice
}
