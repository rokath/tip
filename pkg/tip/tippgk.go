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
/*
// tipInit is usable to change settings. Invalid values are ignored silently.
static void tipInit( 
	unsigned urcb,         //!< @param unreplacableContainerBits, only 6 and 7 are accepted
	unsigned id1Count,     //!< @param ID1Count, used ID1 values of ID1Max. The remaining space is for indirect indexing.
	uint8_t const * idt ){ //!< @param idTable
 
	if( urcb == 6 || urcb == 7 ){
		unreplacableContainerBits = urcb;
	}
	if( unreplacableContainerBits == 6 ){
		ID1Max = 191; 
	} else if (unreplacableContainerBits == 7){
		ID1Max = 127;
	} 
	if (id1Count <= ID1Max){
		ID1Count = id1Count;
	}else if (ID1Count > ID1Max) {
		ID1Count = ID1Max;
	}
	unsigned ID2Count =  (ID1Max - ID1Count) * 255;
	MaxID = ID1Count + ID2Count;
	if (idt) {
		IDTable = idt;
	}
	uint8_t const * p = IDTable;
	maxPatternlength = 0;
	LastID = 0;
	while( *p ){
		LastID++;
		maxPatternlength = *p < maxPatternlength ? maxPatternlength : *p;
		p += *p + 1;
	}
}

static unsigned urcb;        //!< stored unreplacableContainerBits, only 6 and 7 are accepted
static unsigned id1Count;    //!< stored ID1Count, used ID1 values. The remaining space is for indirect indexing.
static uint8_t const * idt;  //!< stored idTable

//! @brief tipStoreGlobals reads global varibles (to restore them later).
static void tipStoreGlobals( void ) {
	urcb = unreplacableContainerBits;
	id1Count = ID1Count;    
	idt = IDTable;
}

//! @brief tipRestoreGlobals restores global varibles (after storing them).
static void tipRestoreGlobals( void ) {
	unreplacableContainerBits = urcb;
	ID1Count = id1Count;    
	IDTable = idt;
}

//! @brief tiPack encodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip encoding it uses the passed ubc, id1Count and ID table.
//! @param dst is the destination buffer.
//! @param src is the source buffer.
//! @param slen is the valid data length inside source buffer.
//! @param ubc is the unreplacable container bits count, only 6 and 7 are accepted.
//! @param id1Count is the used ID1 values from ID1Max. The remaining space is for indirect indexing.
//! @param idt is the used id table.
//! @retval length of tip packet
size_t tiPack( uint8_t * dst, uint8_t const * src,size_t slen, unsigned ubc, unsigned id1Count, uint8_t const * idt ){
	tipStoreGlobals();
    tipInit(ubc, id1Count, idt);
    size_t plen = tip( dst, src, slen );
	tipRestoreGlobals();
	return plen;
}

//! @brief tiUnpack decodes src buffer with size len into dst buffer and returns encoded len.
//! @details For the tip decoding it uses the passed ubc, id1Count and ID table.
//! @param dst is the destination buffer.
//! @param src is the source buffer.
//! @param slen is the valid data length inside source buffer.
//! @param ubc is the unreplacable container bits count, only 6 and 7 are accepted.
//! @param id1Count is the used ID1 values from ID1Max. The remaining space is for indirect indexing.
//! @param idt is the used id table.
//! @retval length of untip packet
size_t tiUnpack( uint8_t* dst, const uint8_t * src, size_t slen, unsigned ubc, unsigned id1Count, const uint8_t * idt ){
	tipStoreGlobals();
    tipInit(ubc, id1Count, idt);
    size_t plen = tiu( dst, src, slen );
	tipRestoreGlobals();
	return plen;
}
*/
import "C"

import (
	"unsafe"
)

// Pack compresses in to out with no zeroes in out and returns packed size plen.
// out needs to have a size of at least 8*len(in)/7 + 1 for the case in cannot get compressed.
func TIPack(out, in []byte, urc, id1Max int, table []byte) (plen int) {
	dst := (*C.uchar)(unsafe.Pointer(&out[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	tbl := (*C.uchar)(unsafe.Pointer(&table[0]))
	dlen := C.tiPack(dst, src, slen, (C.uint)(urc), (C.uint)(id1Max), tbl)
	return int(dlen)
}

// TIUnpack decompresses in to out and returns unpacked size ulen.
// for the case if in has max possible compression.
func TIUnpack(out, in []byte, urbc, id1Count int, table []byte) (ulen int) {
	dst := (*C.uchar)(unsafe.Pointer(&out[0]))
	tbl := (*C.uchar)(unsafe.Pointer(&table[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	dlen := C.tiUnpack(dst, src, slen, (C.uint)(urbc), (C.uint)(id1Count), tbl)
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
