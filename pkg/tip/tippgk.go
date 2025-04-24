// Package tip is a wrapper for testing the target C-code.
// For some reason inside the *_test.go an 'import "C"' is not possible.
package tip

// #cgo CFLAGS: -g -Wall -I../../src.config -I../../src -I../../../trice/src
// #include <stdint.h>
// #include <stddef.h>
// #include "tipInternal.h"
// #include "memmem.c"
// #include "idTable.c"
// #include "ti_pack.c"
// #include "ti_unpack.c"
// #include "shift87bit.c"
// #include "shift78bit.c"
// #include "shift86bit.c"
// #include "shift68bit.c"
// #include "tip.c"
// int optimizeUnreplacablesEnabled(void) {
//     return OPTIMIZE_UNREPLACABLES;
// }
// int unreplacableBitCount(void) {
//     return unreplacableContainerBits;
// }
// int maxPatternSize(void) {
//     return maxPatternlength;
// }
import "C"

import (
	"unsafe"
)

// maxPatternSize returns the length of the longest existing pattern inside ID table.
func MaxPatternSize() int {
	x := C.maxIdPatternLength()
	return int(x)
}

// OptimizeUnreplacablesEnabled returns, if in tipConfig.h OPTIMIZE_UNREPLACABLES was set.
func OptimizeUnreplacablesEnabled() bool {
	x := C.optimizeUnreplacablesEnabled()
	return x > 0
}

// UnreplacableBitCount return the bit count used for unreplacable conversation (6 or 7).
func UnreplacableBitCount() int {
	x := C.unreplacableBitCount()
	return int(x)
}

// Pack compresses in to out with no zeroes in out and returns packed size plen.
// out needs to have a size of at least 8*len(in)/7 + 1 for the case in cannot get compressed.
func TIPack(out, table, in []byte) (plen int) {
	dst := (*C.uchar)(unsafe.Pointer(&out[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	tbl := (*C.uchar)(unsafe.Pointer(&table[0]))
	slen := (C.size_t)(len(in))
	dlen := C.tiPack(dst, tbl, src, slen)
	return int(dlen)
}


// Pack compresses in to out with no zeroes in out and returns packed size plen.
// out needs to have a size of at least 8*len(in)/7 + 1 for the case in cannot get compressed.
func TIPack2(out, in []byte, urc, id1Max int, table []byte) (plen int) {
	dst := (*C.uchar)(unsafe.Pointer(&out[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	tbl := (*C.uchar)(unsafe.Pointer(&table[0]))
	dlen := C.tiPack2(dst, src, slen, (C.uint)(urc), (C.uint)(id1Max), tbl)
	return int(dlen)
}

// TIUnpack decompresses in to out and returns unpacked size ulen.
// for the case if in has max possible compression.
func TIUnpack(out, table, in []byte) (ulen int) {
	dst := (*C.uchar)(unsafe.Pointer(&out[0]))
	tbl := (*C.uchar)(unsafe.Pointer(&table[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	dlen := C.tiUnpack(dst, tbl, src, slen)
	return int(dlen)
}

// Pack compresses in to out with no zeroes in out and returns packed size plen.
// out needs to have a size of at least 8*len(in)/7 + 1 for the case in cannot get compressed.
func Pack(out, in []byte) (plen int) {
	o := (*C.uchar)(unsafe.Pointer(&out[0])) //o := unsafe.Pointer((*C.uchar)(&out[0]))
	i := (*C.uchar)(unsafe.Pointer(&in[0]))  //i := unsafe.Pointer((*C.uchar)(&in[0]))
	ilen := (C.size_t)(len(in))
	olen := C.tip(o, i, ilen)
	return int(olen)
}

// Unpack decompresses in to out and returns unpacked size ulen.
// for the case if in has max possible compression.
func Unpack(out, in []byte) (ulen int) {
	o := (*C.uchar)(unsafe.Pointer(&out[0]))
	i := (*C.uchar)(unsafe.Pointer(&in[0]))
	ilen := (C.size_t)(len(in))
	olen := C.tiu(o, i, ilen)
	return int(olen)
}

type IDPos struct {
	id    byte
	start int
}

// ! NewIDPositionTable is a wrapper for testing C function createIDPosTable and therefore returns posTable.
func NewIDPositionTable(idTable, in []byte) (posTable []IDPos) {
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	idPatTbl := (*C.uchar)(unsafe.Pointer(&idTable[0]))
	C.createIDPosTable(idPatTbl, src, slen)
	n := int(C.IDPosTable.count)
	pt := (*[C.TIP_SRC_BUFFER_SIZE_MAX]C.IDPosition_t)(unsafe.Pointer(&C.IDPosTable.item[0]))
	posTable = make([]IDPos, n)
	for i := range posTable {
		posTable[i].id = byte(pt[i].id)
		posTable[i].start = int(pt[i].start)
	}
	return
}
