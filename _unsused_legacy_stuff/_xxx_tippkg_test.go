

// uses idTable.c
func _TestPack(t *testing.T) { // uses idTable.c
	packet := make([]byte, 100)
	for _, x := range testTable() {
		n := Pack(packet, x.unpacked)
		act := packet[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.packed, act)
	}
}

// uses idTable.c
func _TestUnpack(t *testing.T) { // uses idTable.c
	buffer := make([]byte, 100)
	for _, x := range testTable() {
		n := Unpack(buffer, x.packed)
		act := buffer[:n]
		assert.Equal(t, x.packed, act)
	}
}

/*
type IDPos struct {
	id    byte
	start int
}

// ! NewIDPositionTable is a wrapper for testing C function createIDPosTable and therefore returns posTable.
func NewIDPositionTable(idTable, in []byte) (posTable []IDPos) {
	src := (*C.uchar)(unsafe.Pointer(&in[0]))
	slen := (C.size_t)(len(in))
	idPatTbl := (*C.uchar)(unsafe.Pointer(&idTable[0]))
	C.createIDPosTable(idPatTbl, src, slen)
	n := int(C.IDPosTable.count)
	pt := (*[C.TIP_SRC_BUFFER_SIZE_MAX]C.IDPosition_t)(unsafe.Pointer(&C.IDPosTable.item[0]))
	posTable = make([]IDPos, n)
	for i := range posTable {
		posTable[i].id = byte(pt[i].id)
		posTable[i].start = int(pt[i].start)
	}
	return
}
*/


// TODO
func _TestNewIDPositionTable(t *testing.T) {
	type args struct {
		idTable []byte
		in      []byte
	}
	tests := []struct {
		name         string
		args         args
		wantPosTable []IDPos
	}{ // test cases:
		{
			"",
			args{
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 'p', 2, 'p', 0xbb, 0}, // idTable
				[]byte{0xff, 0x00, 'p', 0xbb, 0xee, 0xff, 'p', 0xbb, 0xcc},    // src
			},
			[]IDPos{{3, 2}, {2, 4}, {3, 6}},
		},
		{
			"",
			args{
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 'p', 2, 'p', 0xbb, 0}, // idTable
				[]byte{'p', 0xbb, 'p', 0xbb, 0xcc},                            // src
			},
			[]IDPos{{3, 0}, {3, 2}},
		},
		{
			"",
			args{
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 'p', 2, 'p', 0xbb, 0}, // idTable
				[]byte{0xff, 0x00, 0xcc}, // src
			},
			[]IDPos{},
		},
		{
			"",
			args{
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 'p', 2, 'x', 'x', 0}, // idTable
				[]byte{'p', 'p', 'p', 'p', 'p'},                              // src
			},
			[]IDPos{{3, 0}, {3, 1}, {3, 2}, {3, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPosTable := NewIDPositionTable(tt.args.idTable, tt.args.in); !reflect.DeepEqual(gotPosTable, tt.wantPosTable) {
				t.Errorf("NewIDPositionTable() = %v, want %v", gotPosTable, tt.wantPosTable)
			}
		})
	}
}
