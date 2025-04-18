package pattern

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/spf13/afero"
)

var (
	PatternSizeMax int // PatternSizeMax is the longest pattern to search for to build idTable.c
	Verbose        bool
)

func init() {
	flag.IntVar(&PatternSizeMax, "z", 8, "max pattern size to find")
}

// Pattern contains a pattern and its occurrances count.
type Pattern struct {
	Bytes      []byte // Bytes is the pattern as byte slice.
	Pos        []int  // Pos holds all start occurances of Bytes. Its len is the occurances count.
	DeletedPos []int  // DeletedPos holds all deleted Pos elements.
}

// Histogram objects hold pattern strings occurences count.
type Histogram struct {
	Hist map[string]Pattern // Hist is the created histogram. len(Pat.Pos) is the occurrances count.
	mu   *sync.Mutex        // mu guaranties mutual exclusion access to the histogram.
	Keys []string           // Key holds all histogram keys separately for faster processing.
}

// NewHistogram returns a new Histogram instance.
func NewHistogram(mu *sync.Mutex) *Histogram {
	h := make(map[string]Pattern, 10000)
	object := Histogram{h, mu, nil}
	return &object
}

// ScanData searches data for any 2-to-max bytes sequences and extends p
// with them as key strings hex encoded with their increased count as values in hist.
// ScanData( searches data for any 2-to-max bytes sequences and extends p.Hist with them.
// The byte sequences are getting hex encodedand used as keys with their increased count as values in p.Hist.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
// When ring is true, the data are considered as ring.
func (p *Histogram) ScanData(data []byte, maxPatternSize int, ring bool) {
	var wg sync.WaitGroup
	for i := range maxPatternSize - 1 { // loop over pattern sizes
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			p.scanForRepetitions(data, k+2, ring)
		}(i)
	}
	wg.Wait()
}

func (p *Histogram) DeleteEmptyKeys() {
	for k, v := range p.Hist {
		if len(v.Pos) == 0 {
			delete(p.Hist, k)
		}
	}
}

// DiscardSeldomPattern removes all keys occuring only discardSize or less often.
func (p *Histogram) DiscardSeldomPattern(discardSize int) {
	hlen := len(p.Hist)
	counts := 0 // make([]int, discardSize)
	for k, v := range p.Hist {
		if len(v.Pos) <= discardSize {
			delete(p.Hist, k)
			counts++
		}
	}
	if Verbose {
		fmt.Println(counts, "of", hlen, "patterns removed;", len(p.Hist), "remaining,")
	}
}

// scanForRepetitions searches data for ptLen bytes sequences
// and adds them as key strings hex encoded with their count as values to p.Hist.
// Also the pattern positions are recorded.
// This pattern search algorithm: Start at offset 0 with ptLen bytes from data as pattern
// and search data for repetitions by moving byte by byte.
// When ring is true, the data are considered as ring.
func (p *Histogram) scanForRepetitions(data []byte, ptLen int, ring bool) {
	if ring {
		data = append(data, data[:ptLen-1]...)
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

			pt := p.Hist[key]
			pt.Bytes = pat
			pt.Pos = append(pt.Pos, k)
			p.Hist[key] = pt

			p.mu.Unlock()
		}(i)
	}
	wg.Wait()
}

// ExportAsList converts p.Hist into list and restores original patterns.
func (p *Histogram) ExportAsList() (list []Pattern) {
	list = make([]Pattern, len(p.Hist))
	var i int
	p.mu.Lock()
	for _, value := range p.Hist {
		list[i] = value
		i++
	}
	p.mu.Unlock()
	return
}

// ScanFile reads iFn ands its data to the histogram.
func (p *Histogram) ScanFile(fSys *afero.Afero, iFn string, maxPatternSize int) (err error) {
	data, err := fSys.ReadFile(iFn)
	if err != nil {
		return err
	}
	p.ScanData(data, maxPatternSize, false)
	return nil
}

// ScanAllFiles traverses location and adds all files as sample data.
func (p *Histogram) ScanAllFiles(fSys *afero.Afero, location string, maxPatternSize int) (err error) {
	numScanned := 0
	var wg sync.WaitGroup
	err = filepath.Walk(location, func(path string, _ os.FileInfo, _ error) error {
		numScanned++
		fmt.Println(path)
		if dir, e := fSys.IsDir(path); dir {
			return e
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := p.ScanFile(fSys, path, maxPatternSize)
			if err != nil {
				log.Fatal(err)
			}
		}()
		return nil
	})
	wg.Wait()
	if Verbose {
		fmt.Println(numScanned, "files scanned")
	}
	return
}

// SortByDescWeight sorts and returns list ordered for descenting weight, count and pattern length.
// It also sorts alphabetical to get reproducable results.
func SortByDescWeight(list []Pattern) []Pattern {
	compareFn := func(a, b Pattern) int {
		aWeight := len(a.Pos) * len(a.Bytes)
		bWeight := len(b.Pos) * len(b.Bytes)
		if aWeight < bWeight {
			return 1
		}
		if aWeight > bWeight {
			return -1
		}
		if len(a.Pos) < len(b.Pos) {
			return 1
		}
		if len(a.Pos) > len(b.Pos) {
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
