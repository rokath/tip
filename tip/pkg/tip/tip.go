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
// size_t TiP( uint8_t* dst, uint8_t const * src, size_t len );
// #include "tipTable.c"
// #include "tiPack.c"
// #include "tiUnpack.c"
// #include "memmem.c"
import "C"

import "unsafe"

func Pack(in []byte) (out []byte) {
	limit := 2 * len(in) // 8*len(in)/7 + 1 is what we need if no compression is possible.
	out = make([]byte, limit)
	olen := C.TiP((*C.uchar)(unsafe.Pointer(&out[0])),
		/*     */ (*C.uchar)(unsafe.Pointer(&in[0])), len(in))
	return out[:olen]
}

func Unack(in []byte) (out []byte) {
	limit := 20 * len(in) // 8*len(in) is what we need if max compression is possible.
	out = make([]byte, limit)
	olen := C.TiU((*C.uchar)(unsafe.Pointer(&out[0])),
		/*     */ (*C.uchar)(unsafe.Pointer(&in[0])), len(in))
	return out[:olen]
}

/*
// setSrcBuffer tells the underlying C code where the tip src buffer is.
func setSrcBuffer(in []byte) {
	C.SetSrc((*C.uchar)(unsafe.Pointer(&in[0]), len(in) )
}

// setDstBuffer tells the underlying C code where the tip dst buffer is.
func setDstBuffer(out []byte) {
	C.SetDst((*C.uchar)(unsafe.Pointer(&out[0]), len(out) )
}

// setOutBuffer tells the underlying C code where to output the tip byte stream.
func setOutBuffer(o []byte) {
	Cout := (*C.uchar)(unsafe.Pointer(&o[0]))
	C.CgoSetOutBuffer(Cout)
}

// triceCheck performs triceCheck C-code sequence n.
func triceCheck(n int) {
	C.TriceCheck(C.int(n))
}

// triceTransfer performs the deferred trice output.
func triceTransfer() {
	C.TriceTransfer()
}
*/
