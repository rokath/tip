package tip

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/tj/assert"
)

var table = []byte{3, 0xaa, 0xaa, 0xaa, 0}

var tipTestTable = []struct {
	unpacked []byte
	packed   []byte
}{
	{[]byte{0xaa, 0xbb}, []byte{0xe0, 0xaa, 0xbb}}, // only unreplacable bytes
	{[]byte{'A', 'B', 'C', 'A', 'B'}, []byte{0x80, 0x80|'A', 0x80|'B', 0x80|'C', 0x80|'A', 0x80|'B'}}, // only unreplacable bytes
	{[]byte{0x41, 0x42, 0x43, 0x41, 0x42}, []byte{0x80, 0xc1, 0xc2, 0xc3, 0xc1, 0xc2}}, // only unreplacable bytes
	{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xfc, 0xaa, 0xbb, 0xcc, 0xaa, 0xbb}}, // only unreplacable bytes
	{[]byte{0xd1, 0xaa, 0xaa, 0xaa, 0xd2}, []byte{0xe0, 0x01, 0xd1, 0xd2}},             // 1 pattern in the middle
	{[]byte{0xd1, 0xd2, 0xaa, 0xaa, 0xaa}, []byte{0xe0, 0xd1, 0x01, 0xd2}},             // 1 pattern in the end
	{[]byte{0xaa, 0xaa, 0xaa, 0xd1, 0xd2}, []byte{0x01, 0xe0, 0xd1, 0xd2}},             // 1 pattern at start
	{[]byte{0xaa, 0xaa, 0xaa}, []byte{0x01}},                                           // Just 1 pattern
	{[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}, []byte{0x01, 0x01}},                   // just 2 pattern
	{[]byte{0xd1, 0xaa, 0xaa, 0xaa, 0xd2, 0xaa, 0xaa, 0xaa, 0xd3}, []byte{0xf0, 0x01, 0xd1, 0x01, 0xd2, 0xd3}}, // 2 pattern with distributed unreplacable bytes
}

func TestTIPack(t *testing.T) {
	packet := make([]byte, 100)
	for _, x := range tipTestTable {
		n := TIPack(packet, table, x.unpacked)
		fmt.Println("Tip pack result:", hex.EncodeToString(packet[:n]))
		assert.Equal(t, len(x.packed), n)
		act := packet[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.packed, act)
	}
}

func TestTIUnpack(t *testing.T) {
	buffer := make([]byte, 100)
	for _, x := range tipTestTable {
		assertNoZeroes(t, x.packed)
		n := TIUnpack(buffer, table, x.packed)
		assert.Equal(t, len(x.unpacked), n)
		act := buffer[:n]
		assert.Equal(t, x.unpacked, act)
	}
}

// uses idTable.c
func _TestPack(t *testing.T) { // uses idTable.c
	packet := make([]byte, 100)
	for _, x := range tipTestTable {
		n := Pack(packet, x.unpacked)
		act := packet[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.packed, act)
	}
}

// uses idTable.c
func _TestUnpack(t *testing.T) { // uses idTable.c
	buffer := make([]byte, 100)
	for _, x := range tipTestTable {
		n := Unpack(buffer, x.packed)
		act := buffer[:n]
		assert.Equal(t, x.packed, act)
	}
}

// TestTIPackTIUnpack packs, checks for no zeroes, unpacks and compares.
func TestTIPackTIUnpack(t *testing.T) {
	buffer := make([]byte, 100)
	packet := make([]byte, 100)
	var ratio float64
	var i uint
	for _, x := range tipTestTable {
		n := TIPack(packet, table, x.unpacked)
		act := packet[:n]

		assertNoZeroes(t, act)

		m := TIUnpack(buffer, table, act)
		res := buffer[:m]
		assert.Equal(t, x.unpacked, res)

		ratio += float64(n) / float64(len(x.unpacked))
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
		{
			"",
			args{
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 0xaa, 2, 0xaa, 0xbb, 0}, // idTable
				[]byte{0xaa, 0xbb, 0xaa, 0xbb, 0xcc},                            // src
			},
			[]IDPos{{3, 0}, {3, 2}},
		},
		{
			"",
			args{
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 0xaa, 2, 0xaa, 0xbb, 0}, // idTable
				[]byte{0xff, 0x00, 0xcc}, // src
			},
			[]IDPos{},
		},
		{
			"",
			args{
				[]byte{5, 0, 0, 0, 0, 0, 3, 0xee, 0xff, 0xaa, 2, 0xaa, 0xaa, 0}, // idTable
				[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa},                            // src
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
