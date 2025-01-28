package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

var (
	version string // do not initialize, goreleaser will handle that
	commit  string // do not initialize, goreleaser will handle that
	date    string // do not initialize, goreleaser will handle that
	iFn     string // input file name
	oFn     = ".tipTable.c"
)

func init() {
	flag.StringVar(&iFn, "i", "-", "input file name")
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
	if len(os.Args) != 3 {
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
	maxPatternSize := 20
	tips := generateSortedPatternList(iData, maxPatternSize)
	tips = tips[:127]
	generateTipTable(fSys, oFn, tips, stat, maxPatternSize)
}

// pc contains a pattern and its count.
type pc struct {
	s string // s is the as string encoded pattern already as print line.
	n int    // n is the count of pattern occurances.
}

// byteSliceAsASCII returns b as ASCII string. Example:  .Aah.B..C
func byteSliceAsASCII(b []byte) string {
	var s strings.Builder
	for _, x := range b {
		if 0x20 <= x && x <= 0x7f {
			s.WriteString(fmt.Sprintf("%c", x))
		} else {
			s.WriteString(" ")
		}
	}
	return s.String()
}

// createPatternLineString writes pattern as "  n, b0, b1, ..., b(n-1), // AAA  AA" string.
func createPatternLineString(pattern []byte) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%3d,", len(pattern))) // start line with pattern size
	for _, x := range pattern {                      // write pattern bytes
		s.WriteString(fmt.Sprintf(" 0x%02x,", x))
	}
	s.WriteString(fmt.Sprintf("%s // ", spaces(120-6*len(pattern)))) // align
	s.WriteString(byteSliceAsASCII(pattern))                         // write pattern lettes as comment
	return s.String()                                                // no alignment here to keep s length
}

// buildPatternHistogram searches data for any 2 to maxPatternLength bytes sequences
// and returns them as key strings hex encoded with their count as values in m.
func buildPatternHistogram(data []byte, maxPatternLength int) map[string]int {
	m := make(map[string]int, 100000)
	// This pattern search is very simple and may get improved.
	for ptLen := 2; ptLen <= maxPatternLength; ptLen++ { // loop over pattern sizes
		for i := 0; i < len(data)-(ptLen-1); i++ { // loop over data
			key := createPatternLineString(data[i : i+(ptLen)])
			if value, ok := m[key]; ok { // pattern already fount before
				m[key] = value + 1 // increment count
			} else { // first time pattern occurance
				m[key] = 1 // count is 1
			}
		}
	}
	return m
}

func patternHistToList(m map[string]int) (list []pc) {
	list = make([]pc, len(m))
	i := 0
	for k, v := range m {
		list[i].s = k
		list[i].n = v
		i++
	}
	return
}

// descentingCountSort returns list ordered for decreasing count and pattern length.
func descentingCountSort(list []pc) []pc {
	compareFn := func(a, b pc) int {
		if a.n < b.n {
			return 1
		}
		if a.n > b.n {
			return -1
		}
		if len(a.s) < len(b.s) {
			return 1
		}
		if len(a.s) > len(b.s) {
			return -1
		}
		return 0
	}
	slices.SortFunc(list, compareFn)
	return list
}

func generateSortedPatternList(data []byte, maxPatternSize int) []pc {
	m := buildPatternHistogram(data, maxPatternSize)
	list := patternHistToList(m)
	return descentingCountSort(list)
}

// generateTipTable generates a file oFn containing Go code using t[:127] and stat
func generateTipTable(fSys *afero.Afero, oFn string, t []pc, stat os.FileInfo, maxPatternSize int) {
	oh, err := fSys.Create(oFn)
	if err != nil {
		log.Fatal(err)
	}
	defer oh.Close()
	var tipTableSize int
	fmt.Fprintln(oh, `//! @file tipTable.c`)
	fmt.Fprintln(oh, "//! @brief Generated code - do not edit!")
	fmt.Fprintln(oh)
	fmt.Fprintln(oh, "#include <stdint.h>")
	fmt.Fprintln(oh, "#include <stddef.h>")
	fmt.Fprintln(oh)
	fmt.Fprintln(oh, "//! tipTable is sorted by pattern count and pattern length.")
	fmt.Fprintln(oh, "//! The pattern position + 1 is the replacement id.")
	fmt.Fprintf(oh, "uint8_t tipTable[] = { // from %s (%s)%s-- __ASCII__          |  count  id\n", stat.Name(), stat.ModTime().String()[:16], spaces(7*maxPatternSize-59-len(stat.Name())))
	for i, x := range t {
		before, _, found := strings.Cut(x.s, ",")
		if !found {
			fmt.Println("could not split", x.s, "after first ,")
		}
		sz, err := strconv.Atoi(strings.TrimSpace(before))
		if err != nil {
			fmt.Println(before, "delivered", err)
		}
		tipTableSize += 1 + sz
		sp := spaces(129 + maxPatternSize - len(x.s))
		if i < 127 {
			fmt.Fprintf(oh, "\t%s%s|%7d  %02x\n", x.s, sp, x.n, i+1)
		} else {
			if x.n > 1 {
				fmt.Fprintf(oh, " // %s%s|%7d  %6d\n", x.s, sp, x.n, i+1)
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
	var s strings.Builder
	for range n {
		s.WriteString(" ")
	}
	return s.String()
}
