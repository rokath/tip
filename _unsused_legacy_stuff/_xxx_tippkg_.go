package tip

/* NOT EASY !!!
//  typedef struct {
//      int count; //! count is the actual path count in map.
//      uint8_t path[TIP_MAX_PATH_COUNT][TIP_SRC_BUFFER_SIZE_MAX/2+1];
//  } map_t;

type srcMap struct {
	count int       // paths count
	path  [][]uint8 // slice of paths
}

// ! NewSrcMap is a wrapper for testing C function createSrcMap and therefore returns srcMap.
func NewSrcMap(idTable, in []byte) (sm srcMap) {
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	idPatTbl := (*C.uchar)(unsafe.Pointer(&idTable[0]))
	C.createSrcMap(idPatTbl, src, slen)
	sm.count = int(C.srcMap.count)
	path := make([][]byte, sm.count )
	for i := range sm.count {
		path[i] =  append(path[i], byte(77))

		path[i] =  append(path[i], ([C.TIP_MAX_PATH_COUNT]C.uint8_t)(unsafe.Pointer(&C.srcMap.path[i])))
	}
	paths := (*[C.TIP_MAX_PATH_COUNT]C.srcMap_t)(unsafe.Pointer(&C.srcMap.path[0]))
	fmt.Println(sm)
	fmt.Println(paths)
	sm.path = make([][]byte, C.TIP_MAX_PATH_COUNT)
	for i := range sm.path {
		sm.path[i] = make([]byte, C.TIP_SRC_BUFFER_SIZE_MAX-1)
		sm.path[i] = append( sm.path[i], 0)
	}
	//for i := range sm.path {
	//	posTable[i].id = byte(paths)
	//	posTable[i].start = int(pt[i].start)
	//}
	return
}
*/
