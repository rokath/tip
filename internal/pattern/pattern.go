package pattern

import (
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
)

var (
	SizeMax int
	Verbose bool
)

// Histogram objects hold pattern strings occurences count.
type Histogram struct {
	Hist map[string]int
	mu   *sync.Mutex
}

// NewHistogram returns a new Histogram instance.
func NewHistogram(mu *sync.Mutex) *Histogram {
	h := make(map[string]int, 10000)
	object := Histogram{h, mu}
	return &object
}

// countOverlapping returns sub count in s.
// https://stackoverflow.com/questions/67956996/is-there-a-count-function-in-go-but-for-overlapping-substrings
func countOverlapping(s, sub string) int {
	var c int
	for d := range s {
		if strings.HasPrefix(s[d:], sub) {
			c++
		}
	}
	return c
}

// Reduce searches the keys if they contain sub-keys.
// If a sub-key is found inside a key with count n,
// The sub-key count is reduced by n.
// It uses
func (p *Histogram) Reduce(list []Patt) (rlist []Patt) {
	if Verbose {
		fmt.Println("Reducing histogram with length", len(p.Hist), "...")
	}
	dlist := SortByDescentingCountAndLengthAndAphabetical(rlist)
	for i, x := range dlist {
		key := dlist[i].Key // top entry is longest key

		for k := i; k < len(dlist)-1; k++ {
			n := strings.Count(key, x.Key)
			fmt.Println(n) hier weiter
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			fmt.Print(k)
			//strings.Count(list[k].Key, list[])
			// 	//p.scanForRepetitions(data, k+2)
		}(i)
	}
	wg.Wait()

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}

// Extend searches data for any 2-to-max bytes sequences
// and extends p with them as key strings hex encoded with their increased count as values in hist.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func (p *Histogram) Extend(data []byte, maxPatternSize int) {
	if Verbose {
		fmt.Println("Extending histogram with length", len(p.Hist), "...")
	}
	var wg sync.WaitGroup
	for i := 0; i < maxPatternSize-1; i++ { // loop over pattern sizes
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			p.scanForRepetitions(data, k+2)
		}(i)
	}
	wg.Wait()

	if Verbose {
		fmt.Println("Extending histogram...done. New length is", len(p.Hist))
	}
}

// scanForRepetitions searches data for ptLen bytes sequences
// and adds them as key strings hex encoded with their count as values to p.Hist.
// This pattern search algorithm:
// Start at offset 0 with ptLen bytes from data as pattern and search data for repetitions
// by moving byte by byte.
func (p *Histogram) scanForRepetitions(data []byte, ptLen int) {
	if Verbose {
		fmt.Println("scan for count", ptLen, "repetitions...")
	}
	last := len(data) - (ptLen) // This is the last position in data to check for repetitions.

	var wg sync.WaitGroup
	for i := 0; i <= last; i++ { // Loop over all possible pattern.
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			pat := data[k : k+ptLen]
			key := hex.EncodeToString(pat) // We need to convert pat into a key.
			p.mu.Lock()
			p.Hist[key]++
			p.mu.Unlock()
		}(i)
	}
	wg.Wait()
	if Verbose {
		fmt.Println("scan for count", ptLen, "repetitions...done.")
	}
}

// Patt contains a pattern and its occurances count.
type Patt struct {
	Cnt   int    // cnt is the count of occurances.
	Bytes []byte // Bytes is the pattern as byte slice.
	Key   string // key is the pattern as hex string.
}

// histogramToList converts m into list and restores original patterns.
func (p *Histogram) ExportAsList() (list []Patt) {
	list = make([]Patt, len(p.Hist))
	var i int
	p.mu.Lock()
	for key, cnt := range p.Hist {
		list[i].Cnt = cnt
		list[i].Bytes, _ = hex.DecodeString(key)
		list[i].Key = key
		i++
	}
	p.mu.Unlock()
	return
}
