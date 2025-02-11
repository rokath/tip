package tip

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

var table = []byte{3, 0xaa, 0xaa, 0xaa, 0}

func TestTIPack(t *testing.T) {
	tt := []struct {
		in  []byte
		exp []byte
	}{
		{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xfc, 0xaa, 0xbb, 0xcc, 0xaa, 0xbb}},
		{[]byte{0xd1, 0xaa, 0xaa, 0xaa, 0xd2}, []byte{0xe0, 0x01, 0xd1, 0xd2}},
		{[]byte{0xd1, 0xd2, 0xaa, 0xaa, 0xaa}, []byte{0xe0, 0xd1, 0x01, 0xd2}},
		{[]byte{0xaa, 0xaa, 0xaa, 0xd1, 0xd2}, []byte{0x01, 0xe0, 0xd1, 0xd2}},
		{[]byte{0xd1, 0xaa, 0xaa, 0xaa, 0xd2, 0xaa, 0xaa, 0xaa, 0xd3 }, []byte{0xf0, 0x01, 0xd1, 0x01, 0xd2, 0xd3 }},
	}

	buf := make([]byte, 100)
	for _, x := range tt {
		n := TIPack(buf, table, x.in)
		act := buf[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.exp, act)
	}
}

func TestTIUnpack(t *testing.T) {
	tt := []struct {
		in  []byte
		exp []byte
	}{
		{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}},
	}

	buf := make([]byte, 100)
	for _, x := range tt {
		n := TIUnpack(buf, table, x.in)
		act := buf[:n]
		assert.Equal(t, x.exp, act)
	}
}

func TestPack(t *testing.T) { // uses idTable.c
	tt := []struct {
		in  []byte
		exp []byte
	}{
		{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xfc, 0xaa, 0xbb, 0xcc, 0xaa, 0xbb}},
	}

	buf := make([]byte, 100)
	for _, x := range tt {
		n := Pack(buf, x.in)
		act := buf[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.exp, act)
	}
}

func TestUnpack(t *testing.T) { // uses idTable.c
	tt := []struct {
		in  []byte
		exp []byte
	}{
		{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}},
	}

	buf := make([]byte, 100)
	for _, x := range tt {
		n := Unpack(buf, x.in)
		act := buf[:n]
		assert.Equal(t, x.exp, act)
	}
}

// TestTIPackTIUnpack packs, checks for no zeroes, unpacks and compares.
func _TestTIPackTIUnpack(t *testing.T) {
	in := [][]byte{
		{0xaa, 0xbb, 0xcc, 0xaa, 0xbb},
	}
	buf := make([]byte, 100)
	out := make([]byte, 100)
	var ratio float64
	var i uint
	for _, x := range in {
		n := TIPack(buf, table, x)
		act := buf[:n]
		assertNoZeroes(t, act)

		m := TIUnpack(out, table, act)
		res := out[:m]
		assert.Equal(t, x, res)

		ratio += float64(n) / float64(len(x))
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
