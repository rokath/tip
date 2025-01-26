
package tip

// Copyright 2025 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package tip is a helper for testing the target C-code, but also usable in Go code.
// Each C function gets a Go wrapper which is tested in appropriate test functions.
// For some reason inside the tip_test.go an 'import "C"' is not possible.

// The Go functions defined here are not exported. They are called by the Go test functions in this package.
// This way the test functions are executing the trice C-code compiled with the triceConfig.h here.
// Inside ./testdata this file is named cgoPackage.go where it is maintained.
// The ../renewIDs_in_examples_and_test_folder.sh script copies this file under the name generated_cgoPackage.go into various
// package folders, where it is used separately.package tip

// #cgo CFLAGS: -g -Wall -I../../src
// #include "tip.h"
// #include "memmem.c"
// #include "tipTable.c"
// #include "tip.c"
// #include "pack.c"
// #include "unpack.c"
import "C"

import "unsafe"

// Pack compresses in to out with no zeroes in out and returns packed size plen.
// out needs to have a size of at least 8*len(in)/7 + 1 for the case in cannot get compressed.
func Pack(out, in []byte) (plen int) {
	ilen := (C.size_t)(len(in))
	olen := C.tip((*C.uchar)(unsafe.Pointer(&out[0])),
		/*     */ (*C.uchar)(unsafe.Pointer(&in[0])), ilen)
	return int(olen)
}

// Unpack decompresses in to out and returns unpacked size ulen.
// out needs to have a size of at least TIP_PATTERN_SIZE_MAX*len(in)
// for the case if in has max possible compression was.
func Unpack(out, in []byte) (ulen int) {
	ilen := (C.size_t)(len(in))
	olen := C.tiu((*C.uchar)(unsafe.Pointer(&out[0])),
		/*     */ (*C.uchar)(unsafe.Pointer(&in[0])), ilen)
	return int(olen)
}
