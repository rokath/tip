// Copyright 2025 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package tip_cgo is a helper for testing the target C-code functions.
// Each C function to test gets a Go wrapper which is tested in appropriate test functions.
// For some reason inside the tip_test.go an 'import "C"' is not possible.

// The Go functions defined here are not exported. They are called by the Go test functions in this package.
// This way the test functions are executing the tip C-code compiled with the tipConfig.h here.
// Inside ./testdata this file is named cgoPackage.go where it is maintained.
package tip

// #cgo CFLAGS: -g -Wall -I../../src
// #include <stdlib.h>
// #include "memmem.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// memMem returns position of needle in hay or -1.
func memMem(hay, needle []byte) (pos int) {
	h := (*C.uchar)(unsafe.Pointer(&hay[0]))
	hlen := (C.size_t)(len(hay))

	n := (*C.uchar)(unsafe.Pointer(&needle[0]))
	nlen := (C.size_t)(len(needle))

	// addr := (*C.uchar)(unsafe.Pointer(C.Memmem(h, hlen, n, nlen)))

	fmt.Println(h, hlen, n, nlen)
	//posX := addr - h;

	//pos = int(posX)

	//if addr != 0 {
	//	return addr - n
	//}
	return -1
}
