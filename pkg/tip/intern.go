package tip

// #include "tip.h"
// replace_t * buildReplaceList(int * rcount, const uint8_t * table, const uint8_t * src, size_t slen);
import "C"

import (
	"encoding/binary"
	"log"
	"unsafe"
)

// ! @brief replace_t is a replace type descriptor.
type replace struct {
	bo offset //  offset_t bo; // bo is the buffer offset, where replace bytes starts. // todo: adapt to tipConfig.h automatically
	sz byte   //  uint8_t  sz; // sz is the replace size (2-255).
	id byte   //  uint8_t  id; // id is the replace byte 0x01 to 0x7f.
}

// buildReplaceList is only for tests.
// To adapt it to different max sizes, change sizeof_bo
func buildReplaceList(table, in []byte) (rpl []replace) {
	tbl := (*C.uchar)(unsafe.Pointer(&table[0])) //o := unsafe.Pointer((*C.uchar)(&out[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))    //i := unsafe.Pointer((*C.uchar)(&in[0]))
	slen := (C.size_t)(len(in))
	var rlen int
	rcount := (*C.int)(unsafe.Pointer(&rlen))

	// https://go.dev/wiki/cgo
	// https://stackoverflow.com/questions/11924196/convert-between-slices-of-different-types
	cArray := unsafe.Pointer(C.buildReplaceList(rcount, tbl, src, slen))

	var x offset
	var sizeof_bo = int(unsafe.Sizeof(x))                  // 1=byte, 2=uint16, 4=uint32
	const sizeof_sz = 1                                    // byte
	const sizeof_id = 1                                    // byte
	var sizeof_replace = sizeof_bo + sizeof_sz + sizeof_id // bytes

	rcnt := int(*rcount)
	length := rcnt * sizeof_replace 
	bytes := C.GoBytes(cArray, C.int(length))
	rpl = make([]replace, rcnt)
	
	for i := range rpl {
		pos := i * sizeof_replace
		rpl[i].bo = readOffset(bytes[pos:])
		rpl[i].sz = bytes[pos+sizeof_bo]
		rpl[i].id = bytes[pos+sizeof_bo+1]
	}
	return
}

// readOffset reds a value of type offset from b.
func readOffset(b []byte) offset {
	ot := offsetType(MaxSize)
	switch v := ot.(type) {
	case byte:
		return offset(b[0])
	case uint16:
		return offset(binary.LittleEndian.Uint16(b))
	case uint32:
		return offset(binary.LittleEndian.Uint32(b))
	default:
		log.Fatalf("I don't know about type %T!\n", v)
		return offset(0)
	}
}

// offsetType returns a value of matching type according to maxSite
func offsetType(maxSize int) interface{} {
	if maxSize < 256 {
		return byte(0)
	}
	if maxSize < 65536 {
		return uint16(0)
	}
	return uint32(0)
}
