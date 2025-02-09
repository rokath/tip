package tip

// Copyright 2025 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package tip is a wwrapper for executing and a helper for testing the target C-code.
// For some reason inside the tip_test.go an 'import "C"' is not possible.

// #cgo CFLAGS: -g -Wall -I../../src
// #include <stdint.h>
// #include <stddef.h>
// #include "memmem.c"
// #include "shift.c"
// #include "idTable.c"
// #include "tip.c"
// #include "pack.c"
// #include "unpack.c"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

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
// out needs to have a size of at least TIP_PATTERN_SIZE_MAX*len(in)
// for the case if in has max possible compression.
func Unpack(out, in []byte) (ulen int) {
	o := (*C.uchar)(unsafe.Pointer(&out[0]))
	i := (*C.uchar)(unsafe.Pointer(&in[0]))
	ilen := (C.size_t)(len(in))
	olen := C.tiu(o, i, ilen)
	return int(olen)
}

//////////////////////////////////////////////////

// ! @brief replace_t is a replace type descriptor.
type replace struct {
	bo uint16 //  offset_t bo; // bo is the buffer offset, where replace bytes starts. // todo: adapt to tipConfig.h automatically
	sz byte   //  uint8_t  sz; // sz is the replace size (2-255).
	id byte   //  uint8_t  id; // id is the replace byte 0x01 to 0x7f.
}

// Pack compresses in to out with no zeroes in out and returns packed size plen.
// out needs to have a size of at least 8*len(in)/7 + 1 for the case in cannot get compressed.
func buildReplaceList(table, in []byte) (rpl []replace) {
	tbl := (*C.uchar)(unsafe.Pointer(&table[0])) //o := unsafe.Pointer((*C.uchar)(&out[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))    //i := unsafe.Pointer((*C.uchar)(&in[0]))
	slen := (C.size_t)(len(in))

	p := C.buildReplaceList(tbl, src, slen)
	fmt.Printf("p=%v\n", p)
	fmt.Printf("rpl=%v\n", rpl)

	// https://go.dev/wiki/cgo
	var theCArray *C.replace_t = C.buildReplaceList(tbl, src, slen)
	length := 10                             // C.getTheArrayLength()
	slice := unsafe.Slice(theCArray, length) // Go 1.17

	fmt.Printf("slice=%v\n", slice)

	fmt.Println("type of slice", reflect.TypeOf(slice))
	//rpl = &slice[0]
	return
}
