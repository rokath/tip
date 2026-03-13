package mem

// #cgo CFLAGS: -g -Wall -I../../src
// #include <stdint.h>
// #include <stddef.h>
// #include "memmem.h"
// #include "memmem.c"
import "C"

import (
	"unsafe"
)

// Mem returns the first position of needle in hay, or -1 if needle is not
// found. It calls the C memmem implementation via cgo and exists primarily to
// test that implementation from Go.
func Mem(hay, needle []byte) (pos int) {
	if hay == nil || needle == nil || len(hay) == 0 || len(needle) == 0 {
		return -1
	}
	h := unsafe.Pointer((*C.uchar)(&hay[0]))
	hlen := (C.size_t)(len(hay))

	n := unsafe.Pointer((*C.uchar)(&needle[0]))
	nlen := C.size_t(len(needle))

	pos = int((C.int)(C.MemmemOffset(h, hlen, n, nlen)))
	return
}
