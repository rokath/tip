package mem

import (
	"testing"

	"github.com/tj/assert"
)

func TestMem(t *testing.T) {
	tt := []struct {
		hay    []byte
		needle []byte
		exp    int
	}{
		{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, 0},
		{[]byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xcc, 0xaa, 0xbb}, 2},
		{[]byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xEE}, -1},
		{[]byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{}, -1},
		{[]byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa}, 3},
		{[]byte{0xaa, 0xBB}, []byte{0xaa, 0xBB}, 0},
		{[]byte{0xBB}, []byte{0xBB}, 0},
		{[]byte{0xBB}, nil, -1},
		{nil, []byte{0xBB}, -1},
		{nil, nil, -1},
	}
	for _, x := range tt {
		pos := Mem(x.hay, x.needle)
		assert.Equal(t, x.exp, pos)
	}
}
