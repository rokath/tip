package slice

import (
	"testing"

	"github.com/tj/assert"
)

func TestIndex(t *testing.T) {
	type e struct {
		s   []byte // slice inside to search
		v   []byte // slice to find
		exp int    // expected index
	}
	tt := []e{
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2}, 0},
		{[]byte{1, 2, 3, 4, 5}, []byte{3, 4}, 2},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4, 5}, 0},
		{[]byte{1, 2, 3, 4, 5}, []byte{1, 2, 4}, -1},
		{[]byte{}, []byte{1, 2, 4}, -1},
		{[]byte{}, []byte{}, 0},
		{[]byte{1, 2, 3, 4, 5}, []byte{}, 0},
	}

	for _, x := range tt {
		idx := Index(x.s, x.v)
		assert.Equal(t, x.exp, idx)
	}
}

func TestCount(t *testing.T) {
	type e struct {
		s   []byte // slice inside to search
		v   []byte // slice to count
		exp int    // expected count
	}
	tt := []e{
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
	}
	for _, x := range tt {
		idx := Count(x.s, x.v)
		assert.Equal(t, x.exp, idx)
	}
}
