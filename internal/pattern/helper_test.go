package pattern

import (
	"encoding/hex"
	"math"
	"slices"
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
