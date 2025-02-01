package pattern

import (
	"encoding/hex"
	"maps"
	"slices"
	"sync"

	"github.com/rokath/tip/pkg/slice"
)

// patt contains a pattern and its occurances count.
type patt struct {
	Cnt   int    // cnt is the count of occurances.
	Bytes []byte // Bytes is the pattern as byte slice.
	key   string // key is the pattern as hex string.
}

// scanForRepetitions searches data for ptLen bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// This pattern search algorithm:
// Start with first ptLen bytes from data as pattern and search data[ptLen:] for a first repetition.
// If a repetition was found at data[n:n+ptLen] continue at data[n+ptLen] and so on.
// The returned map contains all (<=len(data)-ptLen) pattern with their occurances count.
func scanForRepetitions(data []byte, ptLen int) map[string]int {
	m := make(map[string]int, 10000)
	last := len(data) - (ptLen)  // This is the last position in data to check for repetitions.
	for i := 0; i <= last; i++ { // Loop over all possible pattern.
		pat := data[i : i+ptLen]
		key := hex.EncodeToString(pat) // We need to convert pat into a key.
		if _, ok := m[key]; !ok {      // On first pattern occurance, add it with count 1 to map.
			m[key] = 1
		} else {
			continue // pat was already counted
		}
		var n int
		for n = i + ptLen; n <= last; { // Search data after pattern.
			chk := data[n : n+ptLen]
			if slices.Equal(pat, chk) { // found
				m[key] += 1
				n += ptLen // Continue search after pattern.
			} else {
				n++
			}
		}
	}
	return m
}

// buildHistogram searches data for any 2 to maxPatternLength bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func buildHistogram(data []byte, maxPatternLength int) map[string]int {
	subMap := make([]map[string]int, maxPatternLength)
	var wg sync.WaitGroup
	for i := 0; i < maxPatternLength-1; i++ { // loop over pattern sizes
		wg.Add(1)
		go func() {
			defer wg.Done()
			subMap[i] = scanForRepetitions(data, i+2)
		}()
	}
	wg.Wait()
	m := make(map[string]int, 100000)
	for i := 0; i < maxPatternLength; i++ { // loop over pattern sizes
		maps.Copy(m, subMap[i])
	}
	return m
}

// sortByDescentingCountAndLengthAndAphabetical returns list ordered for decreasing count and pattern length.
// It also sorts alphabetical to get reproducable results.
func sortByDescentingCountAndLengthAndAphabetical(list []patt) []patt {
	compareFn := func(a, b patt) int {
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
		if a.key < b.key {
			return 1
		}
		if a.key > b.key {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

// sortByIncreasingLength returns list ordered for increasing pattern length.
// It also sorts alphabetical to get reproducable results.
func sortByIncreasingLength(list []patt) []patt {
	compareFn := func(a, b patt) int {
		if len(a.Bytes) > len(b.Bytes) {
			return 1
		}
		if len(a.Bytes) < len(b.Bytes) {
			return -1
		}
		if a.key > b.key {
			return 1
		}
		if a.key < b.key {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

// reduceSubCounts searches for p[i].Bytes being a part of an other p[k].Bytes with i < k.
// If a pattern A is 3 times in pattern B, the pattern A.Cnt value is decreased by 3.
// Algorithm: check from small to big
func reduceSubCounts(p []patt) []patt {
	if len(p) <= 1 {
		return p // nothing to do
	}
	list := sortByIncreasingLength(p) // smallest pattern first
	for i, x := range list[:len(list)-1] {
		sub := x.Bytes                 // sub is the next (smaller) pattern we want to check.
		for _, y := range list[i+1:] { // range over the next patterns
			n := slice.Count(y.Bytes, sub)
			list[i].Cnt -= n * y.Cnt
		}
	}
	return list
}

// histogramToList converts m into list and restores original patterns.
func histogramToList(m map[string]int) (list []patt) {
	list = make([]patt, len(m))
	var i int
	for key, cnt := range m {
		list[i].Cnt = cnt
		list[i].Bytes, _ = hex.DecodeString(key)
		list[i].key = key
		i++
	}
	return
}

func GenerateSortedList(data []byte, maxPatternSize int) []patt {
	m := buildHistogram(data, maxPatternSize)
	list := histogramToList(m)
	rList := reduceSubCounts(list) // sub pattern are first
	sList := sortByDescentingCountAndLengthAndAphabetical(rList)
	return sList // biggest cnt first, biggest length first on equal cnt
}
