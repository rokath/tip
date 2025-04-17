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
	slices.SortFunc(p.Keys, compareFn)
}

// Extract2BytesPatterns returns all 2-bytes patterns separated from list.
func Extract2BytesPatterns(list []Pattern) (twoBytes, moreBytes []Pattern, max int) {
	twoBytes = make([]Pattern, len(list))
	i := 0
	moreBytes = make([]Pattern, len(list))
	k := 0
	for _, x := range list {
		cnt := len(x.Bytes)
		if cnt > max {
			max = cnt
		}
		if cnt == 2 {
			twoBytes[i] = x
			i++
		} else {
			moreBytes[k] = x
			k++
		}
	}
	twoBytes = twoBytes[:i]
	moreBytes = moreBytes[:k]
	return
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


// SortByDescWeight sorts and returns list ordered for descenting weight, count and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescWeight(list []Pattern) []Pattern {
	compareFn := func(a, b Pattern) int {	
		aWeight := len(a.Pos) * len(a.Bytes)
		bWeight := len(b.Pos) * len(b.Bytes)
		if aWeight < bWeight {
			return 1
		}
		if aWeight > bWeight {
			return -1
		}
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

// SortByDescLength sorts and returns list ordered for descending pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescLength(list []Pattern) []Pattern {
	compareFn := func(a, b Pattern) int {
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
