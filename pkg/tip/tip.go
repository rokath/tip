// Package tip is a wrapper for testing the target C-code.
// For some reason inside the *_test.go an 'import "C"' is not possible.
package tip

// #cgo CFLAGS: -g -Wall -I../../src -I../../../trice/src -I../../examples/L432_inst/Core/inc
// #include <stdint.h>
// #include <stddef.h>
// #include "memmem.c"
// #include "idTable.c"
// #include "pack.c"
// #include "unpack.c"
import "C"

import (
	"fmt"
	"unsafe"
)

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

// ! NewIDPositionTable is a wrapper for testing C function newIDPosTable and therefore returns posTable.
func NewIDPositionTable(idTable, in []byte) (posTable []IDPos) {
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	idPatTbl := (*C.uchar)(unsafe.Pointer(&idTable[0]))
	C.newIDPosTable(idPatTbl, src, slen)
	n := int(C.IDPosTable.count)
	pt := (*[C.TIP_SRC_BUFFER_SIZE_MAX]C.IDPosition_t)(unsafe.Pointer(&C.IDPosTable.item[0]))
	posTable = make([]IDPos, n)
	for i := range posTable {
		posTable[i].id = byte(pt[i].id)
		posTable[i].start = int(pt[i].start)
	}
	return
}

/* NOT EASY !!!
//  typedef struct {
//      int count; //! count is the actual path count in map.
//      uint8_t path[TIP_MAX_PATH_COUNT][TIP_SRC_BUFFER_SIZE_MAX/2+1];
//  } map_t;

type srcMap struct {
	count int       // paths count
	path  [][]uint8 // slice of paths
}

// ! NewSrcMap is a wrapper for testing C function createSrcMap and therefore returns srcMap.
func NewSrcMap(idTable, in []byte) (sm srcMap) {
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	idPatTbl := (*C.uchar)(unsafe.Pointer(&idTable[0]))
	C.createSrcMap(idPatTbl, src, slen)
	sm.count = int(C.srcMap.count)
	path := make([][]byte, sm.count )
	for i := range sm.count {
		path[i] =  append(path[i], byte(77))

		path[i] =  append(path[i], ([C.TIP_MAX_PATH_COUNT]C.uint8_t)(unsafe.Pointer(&C.srcMap.path[i])))
	}
	paths := (*[C.TIP_MAX_PATH_COUNT]C.srcMap_t)(unsafe.Pointer(&C.srcMap.path[0]))
	fmt.Println(sm)
	fmt.Println(paths)
	sm.path = make([][]byte, C.TIP_MAX_PATH_COUNT)
	for i := range sm.path {
		sm.path[i] = make([]byte, C.TIP_SRC_BUFFER_SIZE_MAX-1)
		sm.path[i] = append( sm.path[i], 0)
	}
	//for i := range sm.path {
	//	posTable[i].id = byte(paths)
	//	posTable[i].start = int(pt[i].start)
	//}
	return
}
*/