package pattern

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

// 1 2 3 4 -> 12:1 23:1 34:1 123:1 234:1 1234:1 -> weighted: 12:2 23:2 34:2 123:3 234:3 1234:4
//         -> 12:0 23:- 34:0 123:0 234:0 1234:1 -> weighted: 12:0 23:- 34:0 123:0 234:0 1234:4
// 1 1 1 1 -> 11:3           111:2       1111:1 -> weighted: 11:6           111:6       1111:4
//         -> 11:2           111:1       1111:1 -> weighted: 11:4           111:3       1111:4

func (p *Histogram) AddWeigths() {
	for k, v := range p.Hist {
		p.Hist[k] = v * len(k)
	}
}

// Reduce searches the keys if they contain sub-keys.
// If a sub-key is found inside a key with count n,
// The sub-key count is reduced by n.
// It uses
func (p *Histogram) Reduce() {
	if Verbose {
		fmt.Println("Reducing histogram with length", len(p.Hist), "...")
	}
	if len(p.Key) < 2 { // less than 2 keys
		return
	}
	for i := 0; i < len(p.Key)-1; { // iterate over by increasing length sorted keys
		if len(p.Key[i]) > len(p.Key[i+1]) {
			log.Fatal("unsorted keys")
		}

		if Verbose {
			fmt.Println("Collect 1st group of equal length smaller keys...")
		}
		var smallerKeys []string
		smallerLength := len(p.Key[i]) // is multiple of 2
		for smallerLength == len(p.Key[i]) && i < len(p.Key)-1 {
			smallerKeys = append(smallerKeys, p.Key[i])
			i++
		}
		k := i // Keep position
		if Verbose {
			fmt.Println("Collect 2nd group of equal bigger length keys...")
		}
		var biggerKeys []string
		biggerLength := len(p.Key[i]) // is multiple of 2
		for biggerLength == len(p.Key[i]) && i < len(p.Key)-1 {
			biggerKeys = append(biggerKeys, p.Key[i])
			i++
		}
		if smallerLength == biggerLength {
			fmt.Printf("WARNING: smallerLength == biggerLength == %d\n", biggerLength )
		}
		if smallerLength > biggerLength {
			log.Fatalf("ERROR: smallerLength %d > biggerLength %d", smallerLength, biggerLength )
		}
		p.ReduceOverlappingKeys(biggerKeys, smallerKeys)
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}

func (p *Histogram) ReduceOverlappingKeys(biggerKeys, smallerKeys []string) {
	var wg sync.WaitGroup
	for _, key1st := range biggerKeys {
		wg.Add(1)
		go func(bkey string) {
			defer wg.Done()
			for _, sub := range smallerKeys {
				n := countOverlapping2(bkey, sub) // sub is n-times inside key
				p.mu.Lock()
				a := p.Hist[bkey] // bkey has a count
				b := p.Hist[sub]  // sub has b count
				c := b - 1        // a*n      // new sub count is c
				p.Hist[sub] = c
				if Verbose && n > 0 {
					fmt.Printf("%s(%d) is %d inside %s(%d). -> %s(%d)\n", sub, b, n, bkey, a, sub, c)
				}
				p.mu.Unlock()
			}
		}(key1st)
	}
	wg.Wait()
}

// countOverlapping2 returns sub count in s, assuming s & sub are hex-encoded byte buffers (even length).
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

/*


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


*/

/*
// SortByIncrLength returns list ordered for increasing pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByIncrLength(list []Patt) []Patt {
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
*/
