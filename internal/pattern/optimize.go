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
		if Verbose {
			fmt.Println("Collect 1st group of equal length smaller keys...")
		}
		var smallerKeys []string
		smallerLength := len(p.Key[i]) // is multiple of 2
		for i < len(p.Key)-1 && smallerLength == len(p.Key[i]) {
			smallerKeys = append(smallerKeys, p.Key[i])
			i++
		}
		k := i // Keep position
		if Verbose {
			fmt.Println("Collect 2nd group of equal bigger length keys...")
		}
		var biggerKeys []string
		biggerLength := len(p.Key[i]) // is multiple of 2
		for i < len(p.Key) &&  biggerLength == len(p.Key[i]) {
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

func (p *Histogram) ReduceOverlappingKeys(biggerKeys, smallerKeys []string) {
	var wg sync.WaitGroup
	for _, bkey := range biggerKeys {
		wg.Add(1)
		go func(bigKey string) {
			defer wg.Done()
			for _, subKey := range smallerKeys {
				n := countOverlapping2(bigKey, subKey) // sub is n-times inside key
				p.mu.Lock()
				a := p.Hist[bigKey].Weight // bkey has a count
				v := p.Hist[subKey]
				b := v.Weight      // sub has b count
				v.Weight = b - a*n // new sub count
				p.Hist[subKey] = v
				if Verbose && n > 0 {
					fmt.Printf("%s(%d) is %d inside %s(%d). -> %s(%d)\n", subKey, b, n, bigKey, a, subKey, v.Weight)
				}
				p.mu.Unlock()
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

