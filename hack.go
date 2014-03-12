package gohll

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
