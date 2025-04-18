package pattern

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/tj/assert"
)

// SortPositions puts all positions in ascending order.
func (p *Histogram) SortPositions() {
	for k, _ := range p.Hist {
		slices.Sort(p.Hist[k].Pos)
	}
}

func s2h(s string) string {
	return hex.EncodeToString([]byte(s))
}

func Test_s2h(t *testing.T) {
	s := s2h("XYZ")
	assert.Equal(t, "58595a", s)
}

// SortByDescCount sorts and returns list ordered for descenting count and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescCount(list []Pattern) []Pattern {
	compareFn := func(a, b Pattern) int {
		if len(a.Pos) < len(b.Pos) {
			return 1
		}
		if len(a.Pos) > len(b.Pos) {
			return -1
		}
		if len(a.Bytes) < len(b.Bytes) {
			return 1
		}
		if len(a.Bytes) > len(b.Bytes) {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

func TestSortByDescCountDescLength(t *testing.T) {
	//defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	pat := []Pattern{
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10}},
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10, 20}},
	}
	exp := []Pattern{
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10, 20}},
		{Bytes: []byte{1, 2, 3, 1, 2, 3, 4}, Pos: []int{0, 10}},
	}
	act := SortByDescCount(pat)
	assert.Equal(t, exp, act)
}
