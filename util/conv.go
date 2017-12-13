package util

import (
	"reflect"
	"unsafe"

	"github.com/gramework/runtimer"
)

func init() {
	// debug.SetGCPercent(0)
}

// StringToBytes effectively converts string to bytes
//go:nosplit
func StringToBytes(s string) []byte {
	strstruct := runtimer.StringStructOf(&s)

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(strstruct.Str),
		Len:  strstruct.Len,
		Cap:  strstruct.Len,
	}))
}

// BytesToString effectively converts bytes to string
//go:nosplit
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
