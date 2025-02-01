package slice

import (
	"testing"

	"github.com/tj/assert"
)

func TestIndex(t *testing.T) {
	tt := []struct {
		s   []byte // slice inside to search
		v   []byte // slice to find
		exp int    // expected index
	}{
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2}, 0},
		{[]byte{1, 2, 3, 4, 5}, []byte{3, 4}, 2},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4, 5}, 0},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 4}, -1},
		{[]byte{}, []byte{1, 2, 4}, -1},
		{[]byte{}, []byte{}, 0},
		{[]byte{1, 2, 3, 4, 5}, []byte{}, 0},
		{[]byte{1, 2, 3, 4, 5}, nil, -1},
		{nil, []byte{1, 2, 3, 4, 5}, -1},
		{nil, nil, -1},
	}

	for _, x := range tt {
		idx := Index(x.s, x.v)
		assert.Equal(t, x.exp, idx)
	}
}

func TestCount(t *testing.T) {
	tt := []struct {
		s   []byte // slice inside to search
		v   []byte // slice to count
		exp int    // expected count
	}{
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2}, 1},
		{[]byte{1, 2, 3, 4, 5}, []byte{3, 4}, 1},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4, 5}, 1},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 4}, 0},
		{[]byte{2, 2, 4, 2, 2}, []byte{2, 2}, 2},
		{[]byte{2, 2, 2, 2, 2}, []byte{2, 2}, 2},
		{[]byte{2, 2, 2, 2, 2, 2}, []byte{2, 2}, 3},
		{[]byte{2, 2, 2, 2, 2, 2}, []byte{2}, 6},
		{[]byte{}, []byte{1, 2, 4}, 0},
		{[]byte{}, []byte{}, 0},
		{[]byte{1, 2, 3, 4, 5}, []byte{}, 0},
		{[]byte{1, 2, 3, 4, 5}, nil, 0},
		{nil, []byte{}, 0},
		{nil, nil, 0},
	}
	for _, x := range tt {
		idx := Count(x.s, x.v)
		assert.Equal(t, x.exp, idx)
	}
}
