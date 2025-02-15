package pattern

import (
	"testing"

	"github.com/tj/assert"
)

func TestSortByDescCountDescLength(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	pat := []Patt{
		{100, []byte{1, 2, 3, 1, 2, 3, 4}, "01020301020304"},
		{100, []byte{1, 2, 3, 4}, "01020304"},
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
		{900, []byte{1, 2}, "0102"},
		{100, []byte{8, 2, 3, 1, 2, 3}, "080203010203"},
		{300, []byte{1, 2, 3}, "010203"},
	}
	exp := []Patt{
		{900, []byte{1, 2}, "0102"},
		{300, []byte{1, 2, 3}, "010203"},
		{100, []byte{1, 2, 3, 1, 2, 3, 4}, "01020301020304"},
		{100, []byte{1, 2, 3, 1, 2, 3}, "010203010203"},
		{100, []byte{8, 2, 3, 1, 2, 3}, "080203010203"},
		{100, []byte{1, 2, 3, 4}, "01020304"},
	}
	act := SortByDescCountDescLength(pat)
	assert.Equal(t, exp, act)
}
