package tip

// Copyright 2025 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package tip is a wwrapper for executing and a helper for testing the target C-code.
// For some reason inside the tip_test.go an 'import "C"' is not possible.

// #cgo CFLAGS: -g -Wall -I../../src
// #include <stdint.h>
// #include <stddef.h>
// #include "tip.h"
// #include "memmem.h"
// #include "memmem.c"
// #include "tipTable.c"
// #include "tip.c"
// #include "pack.c"
// #include "unpack.c"
import "C"

import (
	"unsafe"
)

// Pack compresses in to out with no zeroes in out and returns packed size plen.
// out needs to have a size of at least 8*len(in)/7 + 1 for the case in cannot get compressed.
func Pack(out, in []byte) (plen int) {
	ilen := (C.size_t)(len(in))
	olen := C.tip((*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.uchar)(unsafe.Pointer(&in[0])), ilen)
	return int(olen)
}

// Unpack decompresses in to out and returns unpacked size ulen.
// out needs to have a size of at least TIP_PATTERN_SIZE_MAX*len(in)
// for the case if in has max possible compression.
func Unpack(out, in []byte) (ulen int) {
	ilen := (C.size_t)(len(in))
	olen := C.tiu((*C.uchar)(unsafe.Pointer(&out[0])),
		(*C.uchar)(unsafe.Pointer(&in[0])), ilen)
	return int(olen)
}

/*
// See also tiptable.SliceIndex function.
// MemMem returns position of needle in hay or -1.
func MemMem(hay, needle []byte) (pos int) {
	h := (*C.uchar)(unsafe.Pointer(&hay[0]))
	hlen := (C.size_t)(len(hay))

	n := (*C.uchar)(unsafe.Pointer(&needle[0]))
	nlen := (C.size_t)(len(needle))

	pos = int((C.int)(C.MemmemOffset(h, hlen, n, nlen)))
	return
}
*/
