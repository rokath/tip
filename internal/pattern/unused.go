package pattern

import "slices"

// SortKeysDescSize sorts p.Keys by decending size and alphabetical.
func (p *Histogram) SortKeysByDescSize() {
	compareFn := func(a, b string) int {
		if len(a) < len(b) {
			return 1
		}
		if len(a) > len(b) {
			return -1
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
