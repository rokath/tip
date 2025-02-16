package pattern

import "slices"


// SortKeysByIncrSize sorts p.Keys by decending size and alphabetical.
func (p *Histogram) SortKeysByIncrSize() {
	compareFn := func(a, b string) int {
		if len(a) < len(b) {
			return -1
		}
		if len(a) > len(b) {
			return +1
		}
		if a > b {
			return 1
		}
		if a < b {
			return -1
		}
		return 0
	}
	slices.SortFunc(p.Key, compareFn)
}



// SortByDescCountDescLength returns list ordered for decreasing count and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescCountDescLength(list []Patt) []Patt {
	compareFn := func(a, b Patt) int {
		if a.Cnt < b.Cnt {
			return 1
		}
		if a.Cnt > b.Cnt {
			return -1
		}
		if len(a.Bytes) < len(b.Bytes) {
			return 1
		}
		if len(a.Bytes) > len(b.Bytes) {
			return -1
		}
		if a.Key > b.Key {
			return 1
		}
		if a.Key < b.Key {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}
/*
// SortByIncrLength returns list ordered for increasing pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByIncrLength(list []Patt) []Patt {
	compareFn := func(a, b Patt) int {
		if len(a.Bytes) > len(b.Bytes) {
			return 1
		}
		if len(a.Bytes) < len(b.Bytes) {
			return -1
		}
		if a.Key > b.Key {
			return 1
		}
		if a.Key < b.Key {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}
*/
// SortByDescLength returns list ordered for descending pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescLength(list []Patt) []Patt {
	compareFn := func(a, b Patt) int {
		if len(a.Bytes) < len(b.Bytes) {
			return 1
		}
		if len(a.Bytes) > len(b.Bytes) {
			return -1
		}
		if a.Key > b.Key {
			return 1
		}
		if a.Key < b.Key {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}
