package tiptable

import (
	"testing"

	"github.com/tj/assert"
)

func TestSliceIndex(t *testing.T) {
	type e struct {
		s []byte // slice inside to seache
		v []byte // slice to find
		exp int // expected index
	}
	tt := []e{
		{ []byte{1, 2, 3, 4, 5}, []byte{1, 2}, 0 },
		{ []byte{1, 2, 3, 4, 5}, []byte{3, 4}, 2 },
	}

	for _, x := range tt {
		idx := sliceIndex(x.s, x.v)
		assert.Equal(t, x.exp, idx)
	}
}
