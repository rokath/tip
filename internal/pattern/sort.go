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

// Extract2BytesPatterns returns all 2-bytes patterns separated from list.
func Extract2BytesPatterns(list []Patt)( twoBytes, moreBytes []Patt , max int){
	twoBytes = make([]Patt, len(list))
	i := 0 
	moreBytes = make([]Patt, len(list))
	k := 0
	for _, x := range list {
		cnt := len(x.Bytes)
		if cnt > max {
			max = cnt
		}
		if cnt == 2 {
			twoBytes[i] = x
			i++
		}else{
			moreBytes[k] = x
			k++
		}
	}
	twoBytes = twoBytes[:i]
	moreBytes = moreBytes[:k]
	return 
}

// SortByDescCountDescLength returns list ordered for descenting count and pattern length.
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
