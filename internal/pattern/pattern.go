package pattern

import (
	"encoding/hex"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"
)

var (
	PatternSizeMax int
	Verbose        bool
)

// Histogram objects hold pattern strings occurences count.
type Histogram struct {
	Hist map[string]int
	mu   *sync.Mutex
	Keys []string
}

// NewHistogram returns a new Histogram instance.
func NewHistogram(mu *sync.Mutex) *Histogram {
	h := make(map[string]int, 10000)
	object := Histogram{h, mu, nil}
	return &object
}

// Extend searches data for any 2-to-max bytes sequences and extends p.Hist with them.
// The byte sequences are getting hex encodedand used as keys with their increased count as values in p.Hist.
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
// This pattern search algorithm: Start at offset 0 with ptLen bytes from data as pattern
// and search data for repetitions by moving byte by byte. Extend p.Hist accordingly.
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

// GetKeys extracts all p.Hist keys into p.Keys.
func (p *Histogram) GetKeys() {
	p.mu.Lock()
	for key := range p.Hist {
		p.Keys = append( p.Keys, key )
	}
	p.mu.Unlock()
}

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
	slices.SortFunc(p.Keys, compareFn)
}

/*
// countOverlapping returns sub count in s.
// https://stackoverflow.com/questions/67956996/is-there-a-count-function-in-go-but-for-overlapping-substrings
func countOverlapping(s, sub string) int {
	var c int
	for i := range s {
		if strings.HasPrefix(s[i:], sub) {
			c++
		}
	}
	return c
}
*/
// countOverlapping2 returns sub count in s, assuming s & sub are hex-encoded byte buffers
// https://stackoverflow.com/questions/67956996/is-there-a-count-function-in-go-but-for-overlapping-substrings
func countOverlapping2(s, sub string) int {
	var c int
	for i := 0; i < len(s); i += 2 {
		if strings.HasPrefix(s[i:], sub) {
			c++
		}
	}
	return c
}

// Reduce searches the keys if they contain sub-keys.
// If a sub-key is found inside a key with count n,
// The sub-key count is reduced by n.
// It uses
func (p *Histogram) Reduce() {
	if Verbose {
		fmt.Println("Reducing histogram with length", len(p.Hist), "...")
	}
	if len(p.Keys) < 2 {
		return
	}
	for i := 0; i < len(p.Keys)-1; {
		if len(p.Keys[i]) < len(p.Keys[i+1]) {
			log.Fatal("unsorted keys")
		}

		// Collect 1st group of equal length keys...
		var equalSize1stKey []string
		equal1stLength := len(p.Keys[i])
		for equal1stLength == len(p.Keys[i]) && i < len(p.Keys)-1 {
			equalSize1stKey = append(equalSize1stKey, p.Keys[i])
			i++
		}
		k := i // Keep position
		// Collect 2nd group of equal length keys...
		var equalSize2ndKey []string
		equal2ndLength := len(p.Keys[i])
		for equal2ndLength == len(p.Keys[i]) && i < len(p.Keys)-1 {
			equalSize2ndKey = append(equalSize2ndKey, p.Keys[i])
			i++
		}
		p.ReduceOverlappingKeys(equalSize1stKey, equalSize2ndKey)
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}

func (p *Histogram) ReduceOverlappingKeys(equalSize1stKey, equalSize2ndKey []string) {
	// TODO
	for _, key := range equalSize1stKey {
		for _, sub := range equalSize2ndKey {
			n := countOverlapping2(key, sub) // sub is n-times inside key
			a := p.Hist[key]                 // key count is a
			b := p.Hist[sub]                 // sub count is b
			c := b - a*n                     // new count is c
			p.Hist[sub] = c
		}

	}
	/*

			smallerLength := len(p.Keys[i+1]) // Here is equalLength > smallerLength.
			k := i // Process next group of smaller keys with equal size...
			for smallerLength == len(p.Keys[i+1]) {
				n := countOverlapping2(p.Keys[i], p.Keys[i+1])
				p.Hist[p.Keys[i]] -= n
				i++
			}
			i = k // Process next group of smaller keys with equal size...done.



		//rlist = list
		//  dlist := SortByDescCountDescLength(rlist)
		//  for i, x := range dlist {
		//  	key := dlist[i].Key // top entry is longest key
		//
		//  	for k := i; k < len(dlist)-1; k++ {
		//  		n := strings.Count(key, x.Key)
		//  		fmt.Println(n) hier weiter
		//  	}
		//
		//  	var wg sync.WaitGroup
		//  	wg.Add(1)
		//  	go func(k int) {
		//  		defer wg.Done()
		//  		fmt.Print(k)
		//  		//strings.Count(list[k].Key, list[])
		//  		// 	//p.scanForRepetitions(data, k+2)
		//  	}(i)
		//  }
		//  wg.Wait()

	*/

}

// histogramToList converts p.Hist into list and restores original patterns.
func (p *Histogram) ExportAsList() (list []Patt) {
	list = make([]Patt, len(p.Hist))
	var i int
	p.mu.Lock()
	for key, cnt := range p.Hist {
		list[i].Cnt = cnt
		list[i].Bytes, _ = hex.DecodeString(key) // restore bytes
		list[i].Key = key
		i++
	}
	p.mu.Unlock()
	return
}

// Patt contains a pattern and its occurances count.
type Patt struct {
	Cnt   int    // cnt is the count of occurances.
	Bytes []byte // Bytes is the pattern as byte slice.
	Key   string // key is the pattern as hex string.
}
