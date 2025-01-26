// Copyright 2025 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package tip_cgo is a helper for testing the target C-code functions.
// Each C function to test gets a Go wrapper which is tested in appropriate test functions.
// For some reason inside the tip_test.go an 'import "C"' is not possible.

// The Go functions defined here are not exported. They are called by the Go test functions in this package.
// This way the test functions are executing the tip C-code compiled with the tipConfig.h here.
// Inside ./testdata this file is named cgoPackage.go where it is maintained.
package tip_cgo

// #cgo CFLAGS: -g -Wall -I../../src
// #include <stdlib.h>
// #include "memmem.h"
import "C"

import "unsafe"

// memMem returns position of needle in hay or -1.
func memMem(hay, needle []byte) (pos int) {
	h := (*C.uchar)(unsafe.Pointer(&hay[0]))
	n := (*C.uchar)(unsafe.Pointer(&needle[0]))

	addr := (*C.uchar)(unsafe.Pointer(memmem(h, len(hay), n, len(needle))))
	if addr != 0 {
		return addr - n
	}
	return -1
}
