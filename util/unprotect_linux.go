//+build linux

package util

import (
	"fmt"
	"syscall"
	"unsafe"
)

func Unprotect(b []byte) {
	if b == nil {
		return
	}

	ptr := unsafe.Pointer(&b)
	off := uintptr(ptr) % 4096
	fmt.Printf("%d - ", ptr)
	ptr = unsafe.Pointer(uintptr(ptr) - off)
	fmt.Printf("%d = %d", off, ptr)

	if err := syscall.Mprotect(*(*[]byte)(ptr), syscall.PROT_READ|syscall.PROT_WRITE); err != nil {
		panic(fmt.Sprintf("can't unprotect: %s", err))
	}
}
