package pattern

import (
	"encoding/hex"
	"fmt"
	"sync"
)

var (
	SizeMax int
	Verbose bool
)

// PatternHistogram objects hold pattern strings occurences count.
type PatternHistogram struct {
	Hist map[string]int
	mu   *sync.Mutex
}

// Extend searches data for any 2-to-max bytes sequences
// and extends p with them as key strings hex encoded with their increased count as values in hist.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func (p *PatternHistogram) Extend(data []byte, maxPatternSize int) {
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
func (p *PatternHistogram) scanForRepetitions(data []byte, ptLen int) {
	if Verbose {
		fmt.Println("scan for count", ptLen, "repetitions...")
	}
	last := len(data) - (ptLen) // This is the last position in data to check for repetitions.

	var wg sync.WaitGroup
	for i := 0; i <= last; i++ { // Loop over all possible pattern.
		wg.Add(1)
		go
		func(k int) {
			defer wg.Done()
			pat := data[k : k+ptLen]
			key := hex.EncodeToString(pat) // We need to convert pat into a key.
			if k == 3 {
				fmt.Println(k, key)
			}
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
