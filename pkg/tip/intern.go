package tip

// #include "tipinternal.h"
import "C"

import "unsafe"

// ! @brief replace_t is a replace type descriptor.
type replace struct {
	bo byte //  offset_t bo; // bo is the buffer offset, where replace bytes starts. // todo: adapt to tipConfig.h automatically
	sz byte //  uint8_t  sz; // sz is the replace size (2-255).
	id byte //  uint8_t  id; // id is the replace byte 0x01 to 0x7f.
}

// buildReplaceList is only for tests 
func buildReplaceList(table, in []byte) (rpl []replace) {
	tbl := (*C.uchar)(unsafe.Pointer(&table[0])) //o := unsafe.Pointer((*C.uchar)(&out[0]))
	src := (*C.uchar)(unsafe.Pointer(&in[0]))    //i := unsafe.Pointer((*C.uchar)(&in[0]))
	slen := (C.size_t)(len(in))
	var rlen int
	rcount := (*C.int)(unsafe.Pointer(&rlen))

	// https://go.dev/wiki/cgo
	// https://stackoverflow.com/questions/11924196/convert-between-slices-of-different-types
	cArray := unsafe.Pointer(C.buildReplaceList(rcount, tbl, src, slen))
	const sizeof_bo = 1                                      // byte
	const sizeof_sz = 1                                      // byte
	const sizeof_id = 1                                      // byte
	const sizeof_replace = sizeof_bo + sizeof_sz + sizeof_id // bytes
	length := int(*rcount) * sizeof_replace // C.getTheArrayLength()
	bytes := C.GoBytes(cArray, C.int(length))
	rpl = make([]replace, *rcount)
	for i := range rpl {
		pos := i*sizeof_replace 
		rpl[i].bo = bytes[pos]
		rpl[i].sz = bytes[pos+1]
		rpl[i].id = bytes[pos+2]
	}

	// This does not work:
	// xArray := C.buildReplaceList(tbl, src, slen)
	// slice := unsafe.Slice((*replace)(xArray), length) // Go 1.17
	// fmt.Printf("slice=%v\n", slice)
	// rpl = &slice[0]
	// That is ok:
	// slice := unsafe.Slice(&rpl[0], rplCount) // Go 1.17
	// fmt.Printf("slice=%v\n", slice)

	return
}
