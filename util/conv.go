package util

import (
	"unsafe"

	"github.com/gramework/runtimer"
)

// StringToBytes effectively converts string to bytes
func StringToBytes(s string) []byte {
	strstruct := runtimer.StringStructOf(&s)
	return *(*[]byte)(unsafe.Pointer(&runtimer.SliceType2{
		Array: strstruct.Str,
		Len:   strstruct.Len,
		Cap:   strstruct.Len,
	}))
}

// BytesToString effectively converts bytes to string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
