
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

// Mem returns first position of needle in hay or -1.
// Mem calls via CGO the memmem C function and exists to test it.
// See also slice.Index function.
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
