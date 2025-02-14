package tip

// Copyright 2025 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package tip is a wwrapper for executing and a helper for testing the target C-code.
// For some reason inside the tip_test.go an 'import "C"' is not possible.

// #cgo CFLAGS: -g -Wall -I../../src -I../../../trice/src -I../../examples/L432_inst/Core/inc
// #include <stdint.h>
// #include <stddef.h>
// #include "memmem.c"
// #include "idTable.c"
// #include "pack.c"
// #include "unpack.c"
import "C"

import (
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
