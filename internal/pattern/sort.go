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

/*
// SortByDescRateDirect sorts and returns list ordered for descenting count and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByIncrRateDirect(list []Pattern) []Pattern {
	compareFn := func(a, b Pattern) int {	
		if a.RateDirect > b.RateDirect {
			return 1
		}
		if a.RateDirect < b.RateDirect {
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
*/
/*
// SortByDescRateDirect sorts and returns list ordered for descenting count and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByIncrRateIndirect(list []Pattern) []Pattern {
	compareFn := func(a, b Pattern) int {	
		if a.RateIndirect > b.RateIndirect {
			return 1
		}
		if a.RateIndirect < b.RateIndirect {
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

// SortByDescWeight sorts list ordered for descenting weight and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescWeight(list []Pattern){
	compareFn := func(a, b Pattern) int {	
		if a.Weight < b.Weight {
			return 1
		}
		if a.Weight > b.Weight {
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
}

// SortByDescBalance sorts list ordered for descenting weight and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescBalance(list []Pattern){
	compareFn := func(a, b Pattern) int {	
		if a.Balance < b.Balance {
			return 1
		}
		if a.Balance > b.Balance {
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
}
*/
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
