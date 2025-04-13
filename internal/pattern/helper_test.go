package pattern

import (
	"encoding/hex"
	"math"
	"slices"
	"testing"

	"github.com/tj/assert"
)

// SortPositions pouts all positions in ascending order.
func (p *Histogram) SortPositions() {
	for k, _ := range p.Hist {
		slices.Sort(p.Hist[k].Pos)
	}
}

func withinTolerance(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	d := math.Abs(a - b)
	return (d / math.Abs(b)) < epsilon
}

func s2h(s string) string {
	return hex.EncodeToString([]byte(s))
}

func Test_s2h(t *testing.T) {
	s := s2h("XYZ")
	assert.Equal(t, "58595a", s)
}
