package pattern

import (
	"encoding/hex"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/rokath/tip/pkg/slice"
)

var (
	SizeMax int
	Verbose bool
)

// Patt contains a pattern and its occurances count.
type Patt struct {
	Cnt   int    // cnt is the count of occurances.
	Bytes []byte // Bytes is the pattern as byte slice.
	Key   string // key is the pattern as hex string.
}

// buildHistogram searches data for any 2-to-max bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func buildHistogram(data []byte, max int) map[string]int {
	if Verbose {
		fmt.Println("Building histogram...")
	}
	subMap := make([]map[string]int, max) // maps slice
	var wg sync.WaitGroup
	for i := 0; i < max-1; i++ { // loop over pattern sizes
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			subMap[k] = scanForRepetitions(data, k+2)
		}(i)
	}
	wg.Wait()
	m := make(map[string]int, 100000)
	for i := 0; i < max; i++ { // loop over pattern sizes
		maps.Copy(m, subMap[i])
	}
	if Verbose {
		fmt.Println("Building histogram...done. Length is", len(m))
	}
	return m
}

// scanForRepetitions searches data for ptLen bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// This pattern search algorithm:
// Start with first ptLen bytes from data as pattern and search data[ptLen:] for a first repetition.
// If a repetition was found at data[n:n+ptLen] continue at data[n+ptLen] and so on.
// The returned map contains all (<=len(data)-ptLen) pattern with their occurances count.
func scanForRepetitions(data []byte, ptLen int) map[string]int {
	if Verbose {
		fmt.Println("scan for count", ptLen, "repetitions...")
	}
	m := make(map[string]int, 1000)
	last := len(data) - (ptLen) // This is the last position in data to check for repetitions.

	// It seems to be not possible to run go routines here, because "pat was already counted" is not clearly defined.
	//var wg sync.WaitGroup
	//var mu sync.Mutex
	for i := 0; i <= last; i++ { // Loop over all possible pattern.
		//wg.Add(1)
		//go
		func(k int) {
			//defer wg.Done()
			pat := data[k : k+ptLen]
			key := hex.EncodeToString(pat) // We need to convert pat into a key.
			//mu.Lock()
			_, ok := m[key]
			//mu.Unlock()
			if ok {
				return // pat was already counted
			} else {
				//mu.Lock()
				m[key] = 1 // On first pattern occurance, add it with count 1 to map.
				//mu.Unlock()
			}
			var n int
			for n = k + ptLen; n <= last; { // Search data after pattern.
				chk := data[n : n+ptLen]
				if slices.Equal(pat, chk) { // found
					//mu.Lock()
					m[key] += 1
					//mu.Unlock()
					n += ptLen // Continue search after pattern.
				} else {
					n++
				}
			}
		}(i)
	}
	//wg.Wait()
	if Verbose {
		fmt.Println("scan for count", ptLen, "repetitions...done. Count is", len(m))
	}
	return m
}

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

// reduceSubCounts searches for p[i].Bytes being a part of an other p[k].Bytes with i < k.
// Example: If a pattern A is 3 times in pattern B, the pattern A.Cnt value is decreased by 3.
// Algorithm: check from small to big
func reduceSubCounts(p []Patt) []Patt {
	if Verbose {
		fmt.Println("Reducing sub pattern counts...")
	}
	if len(p) <= 1 {
		return p // nothing to do
	}
	list := SortByIncreasingLengthAndAlphabetical(p) // smallest pattern first

	count := getCounts(list) // get a copy to work on
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, x := range list[:len(list)-1] { // last list element is longest pattern
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			if Verbose {
				fmt.Println(k, "...")
			}
			sub := x.Bytes                 // sub is the next (smaller) pattern we want to check.
			for _, y := range list[k+1:] { // range over the next patterns
				n := slice.Count(y.Bytes, sub)
				if n > 0 {
					mu.Lock()
					count[k] -= n * y.Cnt
					mu.Unlock()
				}
			}
			if Verbose {
				fmt.Println(k, "...done")
			}
		}(i)
	}
	wg.Wait()
	setCounts(list, count)
	if Verbose {
		fmt.Println("Reducing sub pattern counts...done")
	}
	return list
}

func getCounts(list []Patt) []int {
	count := make([]int, len(list))
	for i, x := range list {
		count[i] = x.Cnt
	}
	return count
}

func setCounts(list []Patt, count []int) {
	for i := range list {
		list[i].Cnt = count[i]
	}
}

// histogramToList converts m into list and restores original patterns.
func histogramToList(m map[string]int) (list []Patt) {
	list = make([]Patt, len(m))
	var i int
	for key, cnt := range m {
		list[i].Cnt = cnt
		list[i].Bytes, _ = hex.DecodeString(key)
		list[i].Key = key
		i++
	}
	return
}

func GenerateDescendingCountSortedList(data []byte, maxPatternSize int) []Patt {
	m := buildHistogram(data, maxPatternSize)
	list := histogramToList(m)
	//rList := list // reduceSubCounts(list)
	//sList := SortByDescentingCountAndLengthAndAphabetical(rList)
	return list // biggest cnt first, biggest length first on equal cnt
}
