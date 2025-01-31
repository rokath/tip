package pattern

import (
	"encoding/hex"
	"maps"
	"slices"
	"sync"

	"github.com/rokath/tip/pkg/slice"
)

// Type contains a pattern and its occurances count.
type Type struct {
	Cnt   int    // cnt is the count of occurances.
	Bytes []byte // Bytes is the pattern as byte slice.
	Key   string // key is the pattern as hex string.
}

// ScanForRepetitions searches data for ptLen bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// This pattern search algorithm:
// Start with first ptLen bytes from data as pattern and search data[ptLen:] for a first repetition.
// If a repetition was found at data[n:n+ptLen] continue at data[n+ptLen] and so on.
// The returned map contains all (<=len(data)-ptLen) pattern with their occurances count.
func ScanForRepetitions(data []byte, ptLen int) map[string]int {
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

// BuildHistogram searches data for any 2 to maxPatternLength bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func BuildHistogram(data []byte, maxPatternLength int) map[string]int {
	subMap := make([]map[string]int, maxPatternLength)
	var wg sync.WaitGroup
	for i := 0; i < maxPatternLength-1; i++ { // loop over pattern sizes
		wg.Add(1)
		go func() {
			defer wg.Done()
			subMap[i] = ScanForRepetitions(data, i+2)
		}()
	}
	wg.Wait()
	m := make(map[string]int, 100000)
	for i := 0; i < maxPatternLength; i++ { // loop over pattern sizes
		maps.Copy(m, subMap[i])
	}
	return m
}

// SortByRisingLength returns smallest length pattern first.
// On equal length we do not care about the cnt value.
func SortByRisingLength(list []Type) []Type {
	compareFn := func(a, b Type) int {
		if len(a.Bytes) < len(b.Bytes) {
			return 1
		}
		if len(a.Bytes) > len(b.Bytes) {
			return -1
		}
		//  if a.cnt < b.cnt {
		//  	return 1
		//  }
		//  if a.cnt > a.cnt {
		//  	return -1
		//  }
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

// SortByDescentingCountAndLength returns list ordered for decreasing count and pattern length.
func SortByDescentingCountAndLength(list []Type) []Type {
	compareFn := func(a, b Type) int {
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
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

// ReduceSubCounts searches for key being a part of an other key.
// ps is assumed to be sortet by rising pattern length.
// If a pattern A is 3 times in pattern B, the pattern A cnt value is decreased by 3.
// Algorithm: Because ps is sorted, we just check from small to big
func ReduceSubCounts(ps []Type) []Type {
	for i, x := range ps {
		if i == len(ps)-1 {
			continue
		}
		sub := x.Bytes               // sub is the next (smaller) pattern we want to check.
		for _, y := range ps[i+1:] { // range over the next patterns
			pat := y.Bytes
			n := slice.Count(pat, sub)
			x.Cnt -= n
		}
	}
	return ps
}

// HistogramToList converts m into list and restores original patterns.
func HistogramToList(m map[string]int) (list []Type) {
	list = make([]Type, len(m))
	var i int
	for key, cnt := range m {
		list[i].Cnt = cnt
		list[i].Bytes, _ = hex.DecodeString(key)
		list[i].Key = key
		i++
	}
	return
}

func GenerateSortedList(data []byte, maxPatternSize int) []Type {
	m := BuildHistogram(data, maxPatternSize)
	list := HistogramToList(m)
	sList := SortByDescentingCountAndLength(list)  // smallest pattern first
	rList := ReduceSubCounts(sList)                // sub pattern are first
	dList := SortByDescentingCountAndLength(rList) // biggest cnt first, biggest Length first on equal cnt
	return dList
}
