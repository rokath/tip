package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

var (
	version    string // do not initialize, goreleaser will handle that
	commit     string // do not initialize, goreleaser will handle that
	date       string // do not initialize, goreleaser will handle that
	iFn        string // input file name
	oFn        = ".tipTable.c"
	patSizeMax int
)

func init() {
	flag.StringVar(&iFn, "i", "-", "input file name")
	flag.IntVar(&patSizeMax, "z", 8, "max pattern size to find")
	flag.Parse()
	if iFn != "-" {
		oFn = iFn + oFn
	}
}

func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()}
	doit(os.Stdout, fSys)
}

func doit(w io.Writer, fSys *afero.Afero) {
	if len(os.Args) < 3 {
		fmt.Fprintln(w, version, commit, date)
		fmt.Fprintln(w, "Usage: tipTableGen -i inputFileName")
		fmt.Fprintln(w, "Example: `tipTableGen -i fn` creates fn"+oFn)
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}

	// Open and get the input file size.
	ih, err := fSys.Open(iFn)
	if err != nil {
		log.Fatal(err)
	}
	defer ih.Close()
	stat, err := ih.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	iSize := int(stat.Size())
	if iSize > 1024*1024*1014 {
		log.Fatal("input file size", iSize, "is > 1 GB")
	}

	// Read input file into a byte slice.
	iData := make([]byte, iSize)
	_, err = bufio.NewReader(ih).Read(iData)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	tips := generateSortedPatternList(iData, patSizeMax)
	generateTipTable(fSys, oFn, tips, stat, patSizeMax)
}

// patpat_t contains a pattern and its occurances count.
type pat_t struct {
	cnt int    // cnt is the count of occurances.
	pat []byte // pat is the pattern as byte slice.
	key string // key is the pattern as hex string.
}

// ScanForPatternRepetitions searches data for ptLen bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// This pattern search algorithm:
// Start with first ptLen bytes from data as pattern and search data[ptLen:] for a first repetition.
// If a repetition was found at data[n:n+ptLen] continue at data[n+ptLen] and so on.
// The returned map contains all (<=len(data)-ptLen) pattern with their occurances count.
// Example: data = 10 times 1,2,3,4 and ptLen = 2: Result map ["12":10"23":10"34":10"41":100]
// Example: data = 10 times 1,2,3,4 and ptLen = 3: Result map ["123":10"234":10"341":10"412":100]
// Example: data = 10 times 1,2,3,4 and ptLen = 4: Result map ["1234":10"2341":10"3412":10"4123":100]
// Example: data = 10 times 1,2,3,4 and ptLen = 5: Result map ["12341":2 "23412":2 "34123":2 "41232":2]
func ScanForPatternRepetitions(data []byte, ptLen int) map[string]int {
	m := make(map[string]int, 10000)
	limit := len(data) - (ptLen) + 1 // This is the first position in data not to check for repetitions.
	for i := 0; i < limit; i++ {     // Loop over all possible pattern.
		pat := data[i : i+ptLen]
		key := hex.EncodeToString(pat) // We need to convert pat into a key.
		if _, ok := m[key]; !ok {      // On first pattern occurance, add it with count 1 to map.
			m[key] = 1
		} else { // If pattern is in map, for next pattern it could be such case:
			// data =[]byte{1,2,3,4,1,2,3,4,1,2,3,4,1,2,3,4}
			// ptLen=4: At offset 4 we find an existing pattern counted before.
			// ptLen=8: At offset 4 we find an existing pattern not counted before.
			m[key] += 1 // Not sure about that.
			continue    // We would find last repetitions more often, so we ignore the data reminder for this pat.
		}
		for n := i + ptLen; n < limit; n++ { // Start search after pattern.
			chk := data[n : n+ptLen]
			if slices.Equal(pat, chk) {
				m[key] += 1
				n += ptLen // Continue search after pattern.
			}
		}
	}
	return m
}

// buildPatternHistogram searches data for any 2 to maxPatternLength bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
// Pattern of size 1 are skipped, because they give no compression effect when replaced by an id.
func buildPatternHistogram(data []byte, maxPatternLength int) map[string]int {
	subMap := make([]map[string]int, maxPatternLength)
	var wg sync.WaitGroup
	for i := 0; i < maxPatternLength-1; i++ { // loop over pattern sizes
		wg.Add(1)
		go func() {
			defer wg.Done()
			subMap[i] = ScanForPatternRepetitions(data, i+2)
		}()
	}
	wg.Wait()
	m := make(map[string]int, 100000)
	for i := 0; i < maxPatternLength; i++ { // loop over pattern sizes
		maps.Copy(m, subMap[i])
	}
	return m
}

// reduceSubPatternCounts searches for key being a part of an other key.
// If a key A is 3 times in key B, v[A] -= 3 is performed.
// Algorithm:
// * Store
func reduceSubPatternCounts(ps []pat_t) []pat_t {
	return ps
}

func generateSortedPatternList(data []byte, maxPatternSize int) []pat_t {
	m := buildPatternHistogram(data, maxPatternSize)
	list := patternHistToList(m)
	sList := sortPatternByRisingLength(list)     // smallest pattern first
	rList := reduceSubPatternCounts(sList)       // sub pattern are first
	dList := descentingCountAndLengthSort(rList) // biggest cnt first, biggest Length first on equal cnt
	return dList
}

// generateTipTable generates a file oFn containing Go code using list and stat.
// list is assumed to be sorted by list[i].count in decending order.
func generateTipTable(fSys *afero.Afero, oFn string, list []pat_t, stat os.FileInfo, maxPatternSize int) {
	idCount := min(127, len(list))
	oh, err := fSys.Create(oFn)
	if err != nil {
		log.Fatal(err)
	}
	defer oh.Close()
	tipTableSize := 1 // TipTable contains at least table end marker
	fmt.Fprintln(oh, `//! @file tipTable.c`)
	fmt.Fprintln(oh, "//! @brief Generated code - do not edit!")
	fmt.Fprintln(oh)
	fmt.Fprintln(oh, "#include <stdint.h>")
	fmt.Fprintln(oh, "#include <stddef.h>")
	fmt.Fprintln(oh)
	fmt.Fprintln(oh, "//! tipTable is sorted by pattern count and pattern length.")
	fmt.Fprintln(oh, "//! The pattern position + 1 is the replacement id.")
	start := fmt.Sprintf("uint8_t tipTable[] = { // from %s (%s)", stat.Name(), stat.ModTime().String()[:16])
	fill := spaces(9+(6*maxPatternSize - len(start)))
	fill2 := spaces(maxPatternSize-9)
	fmt.Fprintf(oh, start+"%s-- __ASCII__%s|  count  id\n", fill, fill2)
	for i, x := range list {
		pls := createPatternLineString(x.pat, maxPatternSize) // todo: review and improve code
		sz := len(x.pat)
		tipTableSize += 1 + sz
		if i < idCount {
			fmt.Fprintf(oh, "\t%s|%7d  %02x\n", pls, x.cnt, i+1)
		} else {
			if x.cnt > 1 {
				fmt.Fprintf(oh, "//\t%s|%7d  %6d\n", pls, x.cnt, i+1)
			}
		}
	}
	fmt.Fprintln(oh, "\t  0 // table end marker")
	fmt.Fprintln(oh, "};")
	fmt.Fprintln(oh)
	fmt.Fprintf(oh, "const size_t tipTableSize = %d;", tipTableSize)
	fmt.Fprintln(oh)
}

// spaces returns a string consisting of n spaces.
func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	var s strings.Builder
	for range n {
		s.WriteString(" ")
	}
	return s.String()
}

// byteSliceAsASCII returns b as ASCII string size len. Example: "˙Aah˙B˙˙C˙˙     "
// length is used to append spaces until the string has the desired length.
func byteSliceAsASCII(b []byte, length int) string {
	var s strings.Builder
	for _, x := range b {
		if 0x20 <= x && x < 0x7f {
			s.WriteString(fmt.Sprintf("%c", x))
		} else {
			s.WriteString(`˙`)
		}
	}
	s.WriteString(spaces(length - len(b)))
	return s.String()
}

// createPatternLineString writes pattern as "  n, b0, b1, ..., b(n-1), // AAA˙˙AA˙ " string.
func createPatternLineString(pattern []byte, maxPatternSize int) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%3d,", len(pattern))) // start line with pattern size
	for _, x := range pattern {                      // write pattern bytes
		s.WriteString(fmt.Sprintf(" 0x%02x,", x))
	}
	fill := spaces(6*(maxPatternSize - len(pattern)))
	s.WriteString(fmt.Sprintf("%s // ", fill))               // align
	s.WriteString(byteSliceAsASCII(pattern, maxPatternSize)) // write pattern lettes as comment
	return s.String()                                        // no alignment here to keep s length
}

// patternHistToList converts m into list and restores original patterns.
func patternHistToList(m map[string]int) (list []pat_t) {
	list = make([]pat_t, len(m))
	var i int
	for key, cnt := range m {
		list[i].cnt = cnt
		list[i].pat, _ = hex.DecodeString(key)
		list[i].key = key
		i++
	}
	return
}

// descentingCountAndLengthSort returns list ordered for decreasing count and pattern length.
func descentingCountAndLengthSort(list []pat_t) []pat_t {
	compareFn := func(a, b pat_t) int {
		if a.cnt < b.cnt {
			return 1
		}
		if a.cnt > b.cnt {
			return -1
		}
		if len(a.pat) < len(b.pat) {
			return 1
		}
		if len(a.pat) > len(b.pat) {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

// sortPatternByRisingLength returns smallest length pattern first.
// On equal length we do not care about the cnt value.
func sortPatternByRisingLength(list []pat_t) []pat_t {
	compareFn := func(a, b pat_t) int {
		if len(a.pat) < len(b.pat) {
			return 1
		}
		if len(a.pat) > len(b.pat) {
			return -1
		}
		//  if a.cnt < b.cnt {
		//  	return 1
		//  }
		//  if a.cnt > a.cnt {
		//  	return -1
		//  }
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}
