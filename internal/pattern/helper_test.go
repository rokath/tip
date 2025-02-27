package pattern

import (
	"slices"
)

// SortPositions pouts all positions in ascending order.
func (p *Histogram) SortPositions() {
	for k, _ := range p.Hist {
		slices.Sort(p.Hist[k].Pos)
	}
}
