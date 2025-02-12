package shift // package shift is used to test the shift C-functions.

// #cgo CFLAGS: -g -Wall -I../../src
// #include <stdint.h>
// #include <stddef.h>
// #include "pack.c"
// #include "unpack.c"
// #include "memmem.c"  // needed for pack.c code
// #include "idTable.c" // needed for pack.c code
import "C"

import (
	"unsafe"
)

// Shift87bit is shifting all bits in u8 by 1.
// It calls via CGO the shift87bit C function and exists to test it.
// cap is the capacity u8 needs to have. Requirement cap >= 8 * len(u8) / 7 + 1.
// The C-function works in-place to safe memory.
func Shift87bit(u8 []byte) (u7 []byte) {
	src := (*C.uchar)(unsafe.Pointer(&u8[0]))
	sLen := C.size_t(len(u8))
	u7 = make([]byte, 2*len(u8))
	lst := (*C.uchar)(unsafe.Pointer(&u7[len(u7)-1]))
	cnt := int(C.shift87bit(lst, src, sLen))
	u7 = u7[len(u7)-cnt:]
	return
}

// Shift78bit is reverting the Shift87bit operation.
func Shift78bit(u7 []byte) (u8 []byte) {
	src := (*C.uchar)(unsafe.Pointer(&u7[0]))
	sLen := (C.size_t)(len(u7))
	u8 = make([]byte, len(u7))
	dst := (*C.uchar)(unsafe.Pointer(&u8[0]))
	cnt := int(C.shift78bit(dst, src, sLen))
	u8 = u8[:cnt]
	return
}
