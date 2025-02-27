package pattern

import (
	"fmt"
	"math"
	"strings"
	"sync"
)

func (p *Histogram) PrintInfo(message string) {
	var (
		smallest = math.MaxInt
		biggest  = math.MinInt
		sum      int
		count    int
	)
	for _, v := range p.Hist {
		smallest = min(smallest, v.Weight)
		biggest = max(biggest, v.Weight)
		sum += v.Weight
		count++
	}
	fmt.Println(message, "-> count:", count, "sum:", sum, "average:", sum/count, "smallest:", smallest, "biggest:", biggest)
}

// BalanceByteUsage multiplies each key value with maxPatternSize / len(key) to achieve a balance
// in byte usage for pattern of different length. To avoid floats, we use a 1000 times bigger value
func (p *Histogram) BalanceByteUsage(maxPatternSize int) {
	for k, v := range p.Hist {
		v.Weight = v.Weight * 2000 / len(k)
		p.Hist[k] = v
	}
}

// AddWeigths multiplies weight values with key len.
func (p *Histogram) AddWeigths() {
	for k, v := range p.Hist {
		v.Weight *= len(k)
		p.Hist[k] = v
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
	p.SortKeysByIncrSize()
	for i := 0; i < len(p.Key)-1; { // iterate over by increasing length sorted keys
		var smallerKeys []string
		smallerLength := len(p.Key[i]) // is multiple of 2
		for i < len(p.Key)-1 && smallerLength == len(p.Key[i]) {
			smallerKeys = append(smallerKeys, p.Key[i])
			i++
		}
		k := i // Keep position
		var biggerKeys []string
		biggerLength := len(p.Key[i]) // is multiple of 2
		for i < len(p.Key) && biggerLength == len(p.Key[i]) {
			biggerKeys = append(biggerKeys, p.Key[i])
			i++
		}
		if smallerLength < biggerLength { // == on last item
			p.ReduceOverlappingKeys(biggerKeys, smallerKeys)
		}
		i = k // restore position
	}

	if Verbose {
		fmt.Println("Reducinging histogram...done. New length is", len(p.Hist))
	}
}

// positionMatch return pos if a and b have one value common or -1.
func positionMatch(a, b []int) int {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return x
			}
		}
	}
	return -1
}

func (p *Histogram) ReduceOverlappingKeys(biggerKeys, smallerKeys []string) {
	if Verbose {
		fmt.Println("ReduceOverlappingKeys")
		fmt.Println(len(biggerKeys), "\tbiggerKeys:", biggerKeys[:3], "...")
		fmt.Println(len(smallerKeys), "\tsmallerKeys:", smallerKeys[:3], "...")
	}
	var wg sync.WaitGroup
	for _, bkey := range biggerKeys {
		wg.Add(1)
		go func(bigKey string) {
			defer wg.Done()
			for _, subKey := range smallerKeys {
				idx := strings.Index(bkey, subKey)
				if idx == -1 { // subKey not inside bkey
					continue
				}
				p.mu.Lock()
				bkeyPos := p.Hist[bkey].Pos
				subKeyPos := p.Hist[subKey].Pos
				p.mu.Unlock()
				if pos := positionMatch(bkeyPos, subKeyPos); pos >= 0 {
					p.mu.Lock()
					v := p.Hist[subKey]
					v.Weight -= 1
					p.Hist[subKey] = v
					//  if Verbose {
					//  	fmt.Printf("%s(%d) found inside %s(%d) at pos %d.\n", subKey, v.Weight, bigKey, p.Hist[bkey].Weight, pos)
					//  }
					p.mu.Unlock()
				}

				// n := countOverlapping2(bigKey, subKey) // sub is n-times inside key
				// p.mu.Lock()
				// a := p.Hist[bigKey].Weight // bkey has a count
				// v := p.Hist[subKey]
				// b := v.Weight      // sub has b count
				// v.Weight = b - a*n // new sub count
				// p.Hist[subKey] = v
				// if Verbose && n > 0 {
				// 	// fmt.Printf("%s(%d) is %d inside %s(%d). -> %s(%d)\n", subKey, b, n, bigKey, a, subKey, v.Weight)
				// }
				// p.mu.Unlock()
			}
		}(bkey)
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
