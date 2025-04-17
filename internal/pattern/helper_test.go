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
