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
	// PatternSizeMax is the longest pattern to search for when building
	// `idTable.c`.
	PatternSizeMax int
	// Verbose enables progress output while scanning sample data.
	Verbose bool
)

func init() {
	flag.IntVar(&PatternSizeMax, "z", 8, "max pattern size to find")
}

// Pattern contains a byte sequence and the positions where it occurs.
type Pattern struct {
	Bytes      []byte // Bytes is the pattern as a byte slice.
	Pos        []int  // Pos holds all start positions of Bytes. Its length is the occurrence count.
	DeletedPos []int  // DeletedPos holds all removed elements from Pos.
}

// Histogram holds hex-encoded patterns and their occurrences.
type Histogram struct {
	Hist map[string]Pattern // Hist stores all collected patterns. len(Pat.Pos) is the occurrence count.
	mu   *sync.Mutex        // mu guarantees mutually exclusive access to the histogram.
	Keys []string           // Keys caches all histogram keys for later processing.
}

// NewHistogram returns a new Histogram instance.
func NewHistogram(mu *sync.Mutex) *Histogram {
	h := make(map[string]Pattern, 10000)
	object := Histogram{h, mu, nil}
	return &object
}

// ScanData searches data for all 2-byte to maxPatternSize-byte sequences and
// adds them to p.Hist. The byte sequences are hex-encoded and used as keys.
// Patterns of size 1 are skipped because replacing them with an ID would not
// improve compression. When ring is true, the data is treated as a ring buffer.
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

// DeleteEmptyKeys removes all histogram entries that have no positions left.
func (p *Histogram) DeleteEmptyKeys() {
	for k, v := range p.Hist {
		if len(v.Pos) == 0 {
			delete(p.Hist, k)
		}
	}
}

// DiscardSeldomPattern removes all keys occurring discardSize times or fewer.
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

// scanForRepetitions searches data for ptLen-byte sequences
// and adds them as hex-encoded keys to p.Hist.
// Also the pattern positions are recorded.
// The search starts at offset 0 with ptLen bytes from data as the pattern and
// advances byte by byte. When ring is true, the data is treated as a ring
// buffer.
func (p *Histogram) scanForRepetitions(data []byte, ptLen int, ring bool) {
	if ring {
		data = append(data, data[:ptLen-1]...)
	}
	last := len(data) - ptLen // last is the final position in data to check for repetitions.
	var wg sync.WaitGroup
	for i := 0; i <= last; i++ { // Loop over all possible patterns.
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			pat := data[k : k+ptLen]
			key := hex.EncodeToString(pat) // pat must be converted into a map key.
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

// ExportAsList converts p.Hist into a slice of patterns.
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

// ScanFile reads iFn and adds its data to the histogram.
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

// SortByDescWeight sorts and returns list ordered by descending weight, count,
// and pattern length.
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
