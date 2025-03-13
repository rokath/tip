package tip

import (
	"fmt"
	"testing"

	"reflect"

	"github.com/tj/assert"
)

var table = []byte{3, 0xaa, 0xaa, 0xaa, 0}

var tipTestTable = []struct {
	buf []byte
	pkg []byte
}{
	{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xfc, 0xaa, 0xbb, 0xcc, 0xaa, 0xbb}},
	{[]byte{0xd1, 0xaa, 0xaa, 0xaa, 0xd2}, []byte{0xe0, 0x01, 0xd1, 0xd2}},
	{[]byte{0xd1, 0xd2, 0xaa, 0xaa, 0xaa}, []byte{0xe0, 0xd1, 0x01, 0xd2}},
	{[]byte{0xaa, 0xaa, 0xaa, 0xd1, 0xd2}, []byte{0x01, 0xe0, 0xd1, 0xd2}},
	{[]byte{0xd1, 0xaa, 0xaa, 0xaa, 0xd2, 0xaa, 0xaa, 0xaa, 0xd3}, []byte{0xf0, 0x01, 0xd1, 0x01, 0xd2, 0xd3}},
}

func TestTIPack(t *testing.T) {
	packet := make([]byte, 100)
	for _, x := range tipTestTable {
		n := TIPack(packet, table, x.buf)
		assert.Equal(t, len(x.pkg), n)
		act := packet[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.pkg, act)
	}
}

func TestTIUnpack(t *testing.T) {
	buffer := make([]byte, 100)
	for _, x := range tipTestTable {
		assertNoZeroes(t, x.pkg)
		n := TIUnpack(buffer, table, x.pkg)
		assert.Equal(t, len(x.buf), n)
		act := buffer[:n]
		assert.Equal(t, x.buf, act)
	}
}

/*
func TestPack(t *testing.T) { // uses idTable.c
	packet := make([]byte, 100)
	for _, x := range tipTestTable {
		n := Pack(packet, x.buf)
		act := packet[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.pkg, act)
	}
}

func TestUnpack(t *testing.T) { // uses idTable.c
	buffer := make([]byte, 100)
	for _, x := range tipTestTable {
		n := Unpack(buffer, x.pkg)
		act := buffer[:n]
		assert.Equal(t, x.pkg, act)
	}
}
*/

// TestTIPackTIUnpack packs, checks for no zeroes, unpacks and compares.
func TestTIPackTIUnpack(t *testing.T) {
	buffer := make([]byte, 100)
	packet := make([]byte, 100)
	var ratio float64
	var i uint
	for _, x := range tipTestTable {
		n := TIPack(packet, table, x.buf)
		act := packet[:n]

		assertNoZeroes(t, act)

		m := TIUnpack(buffer, table, act)
		res := buffer[:m]
		assert.Equal(t, x.buf, res)

		ratio += float64(n) / float64(len(x.buf))
		i++
	}
	fmt.Println("ratio ", ratio/float64(i))
}

// assertNoZeroes checks that b does not contain any zeroes.
func assertNoZeroes(t *testing.T, b []byte) {
	for _, x := range b {
		assert.NotEqual(t, x, 0)
	}
}

func TestNewIDPositionTable(t *testing.T) {
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
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 0xaa, 2, 0xaa, 0xbb, 0}, // idTable
				[]byte{0xff, 0x00, 0xaa, 0xbb, 0xee, 0xff, 0xaa, 0xbb, 0xcc},    // src
			},
			[]IDPos{{3, 2}, {2, 4}, {3, 6}},
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
