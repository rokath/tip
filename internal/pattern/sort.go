package pattern

import "slices"

// SortByDescentingCountAndLengthAndAphabetical returns list ordered for decreasing count and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescentingCountAndLengthAndAphabetical(list []Patt) []Patt {
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

// SortByIncreasingLengthAndAlphabetical returns list ordered for increasing pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByIncreasingLengthAndAlphabetical(list []Patt) []Patt {
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

// SortByDescendingLength returns list ordered for descending pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescendingLength(list []Patt) []Patt {
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
