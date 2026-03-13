package mem

import (
	"testing"

	"github.com/tj/assert"
)

func TestMem(t *testing.T) {
	tt := []struct {
		name   string
		hay    []byte
		needle []byte
		exp    int
	}{
		{"full match", []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, 0},
		{"subsequence match", []byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xcc, 0xaa, 0xbb}, 2},
		{"missing needle", []byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xEE}, -1},
		{"empty needle", []byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{}, -1},
		{"single byte match", []byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa}, 3},
		{"short full match", []byte{0xaa, 0xBB}, []byte{0xaa, 0xBB}, 0},
		{"single byte haystack", []byte{0xBB}, []byte{0xBB}, 0},
		{"nil needle", []byte{0xBB}, nil, -1},
		{"nil haystack", nil, []byte{0xBB}, -1},
		{"both nil", nil, nil, -1},
		{"needle longer than haystack", []byte{0x01, 0x02}, []byte{0x01, 0x02, 0x03}, -1},
	}
	for _, x := range tt {
		t.Run(x.name, func(t *testing.T) {
			pos := Mem(x.hay, x.needle)
			assert.Equal(t, x.exp, pos)
		})
	}
}
