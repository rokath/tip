package pattern

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

var (
	PatternSizeMax int // PatternSizeMax is the longest pattern to search for to build idTable.c
	Verbose        bool
)

// Pat is the pattern descriptor of Key.
type Pat struct {
	Weight int // Weight is first len(Pos) but gets modifikated later
	Pos []int // Pos holds all start occurances of Key
}

// Histogram objects hold pattern strings occurences count.
type Histogram struct {
	Hist map[string]Pat // Hist is the created histogram. len(Pat.Pos) is the occurrances count.
	mu   *sync.Mutex    // mu guaranties mutual exclusion access to the histogram.
	Key  []string       // Key holds all histogram keys separately for faster processing.
}

// NewHistogram returns a new Histogram instance.
func NewHistogram(mu *sync.Mutex) *Histogram {
	h := make(map[string]Pat, 10000)
	object := Histogram{h, mu, nil}
	return &object
}

// Extend searches data for any 2-to-max bytes sequences
// and extends p with them as key strings hex encoded with their increased count as values in hist.
// Extend searches data for any 2-to-max bytes sequences and extends p.Hist with them.
// The byte sequences are getting hex encodedand used as keys with their increased count as values in p.Hist.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func (p *Histogram) Extend(data []byte, maxPatternSize int) {
	var wg sync.WaitGroup
	for i := range maxPatternSize - 1 { // loop over pattern sizes
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			p.scanForRepetitions(data, k+2)
		}(i)
	}
	wg.Wait()
}

// scanForRepetitions searches data for ptLen bytes sequences
// and adds them as key strings hex encoded with their count as values to p.Hist.
// Also the pattern positions are recorded.
// This pattern search algorithm: Start at offset 0 with ptLen bytes from data as pattern
// and search data for repetitions by moving byte by byte. Extend p.Hist accordingly.
func (p *Histogram) scanForRepetitions(data []byte, ptLen int) {
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
			pt.Pos = append(pt.Pos, k)
			pt.Weight++
			p.Hist[key] = pt
			p.mu.Unlock()
		}(i)
	}
	wg.Wait()
}

// GetKeys extracts all p.Hist keys into p.Keys.
func (p *Histogram) GetKeys() {
	p.mu.Lock()
	for key := range p.Hist {
		p.Key = append(p.Key, key)
	}
	p.mu.Unlock()
}

// ExportAsList converts p.Hist into list and restores original patterns.
func (p *Histogram) ExportAsList() (list []Patt) {
	list = make([]Patt, len(p.Hist))
	var i int
	p.mu.Lock()
	for key, cnt := range p.Hist {
		list[i].Cnt = cnt.Weight
		list[i].Bytes, _ = hex.DecodeString(key) // restore bytes
		list[i].Key = key
		i++
	}
	p.mu.Unlock()
	return
}

// Patt contains a pattern and its occurrances count.
type Patt struct {
	Cnt   int    // Cnt is the count of occurrances.
	Bytes []byte // Bytes is the pattern as byte slice.
	Key   string // key is the pattern as hex string.
}

// ScanFile reads iFn ands its data to the histogram.
func (p *Histogram) ScanFile(fSys *afero.Afero, iFn string, maxPatternSize int) (err error) {
	data, err := fSys.ReadFile(iFn)
	if err != nil {
		return err
	}

	//ss := strings.Split(string(data), ". ") // split ASCII text into sentences (TODO)
        ss := []string{string(data)}

	var wg sync.WaitGroup
	for i, sent := range ss {
		wg.Add(1)
		go func(k int, sentence string) {
			defer wg.Done()
			p.Extend([]byte(sentence), maxPatternSize)
		}(i, sent)
	}
	wg.Wait()
	return
}

// ScanAllFiles traverses location and adds all files as sample data.
func (p *Histogram) ScanAllFiles(fSys *afero.Afero, location string, maxPatternSize int) (err error) {
	numScanned := 0
	err = filepath.Walk(location, func(path string, _ os.FileInfo, _ error) error {
		numScanned++
		fmt.Println(path)
		if dir, e := fSys.IsDir(path); dir {
			return e
		}
		return p.ScanFile(fSys, path, maxPatternSize)
	})
	if Verbose {
		fmt.Println(numScanned, "files scanned")
	}
	return
}
