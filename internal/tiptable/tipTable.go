package tiptable

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rokath/tip/internal/pattern"
	"github.com/spf13/afero"
)

// Generategenerates a file oFn containing C code using iFn and max pattern size.
func Generate(fSys *afero.Afero, oFn, iFn string, maxPatternSize int) {

	data, stat := readData(fSys, iFn)
	list := pattern.GenerateSortedList(data, maxPatternSize)
	// list is sorted by list[i].count, len(list[i].Bytes) and alphabetical in decending order.
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
	fill := spaces(9 + (6*maxPatternSize - len(start)))
	fill2 := spaces(maxPatternSize - 9)
	fmt.Fprintf(oh, start+"%s-- __ASCII__%s|  count  id\n", fill, fill2)
	for i, x := range list {
		pls := createPatternLineString(x.Bytes, maxPatternSize) // todo: review and improve code
		sz := len(x.Bytes)
		tipTableSize += 1 + sz
		if i < idCount {
			fmt.Fprintf(oh, "\t%s|%7d  %02x\n", pls, x.Cnt, i+1)
		} else {
			if x.Cnt > 1 {
				fmt.Fprintf(oh, "//\t%s|%7d  %6d\n", pls, x.Cnt, i+1)
			}
		}
	}
	fmt.Fprintln(oh, "\t  0 // table end marker")
	fmt.Fprintln(oh, "};")
	fmt.Fprintln(oh)
	fmt.Fprintf(oh, "const size_t tipTableSize = %d;", tipTableSize)
	fmt.Fprintln(oh)
}

// readData reads file iFn into data and returns also the iFn file info in stat.
func readData(fSys *afero.Afero, iFn string) (data []byte, stat os.FileInfo) {
	// Open and get the input file size.
	ih, err := fSys.Open(iFn)
	if err != nil {
		log.Fatal(err)
	}
	defer ih.Close()
	stat, err = ih.Stat()
	if err != nil {
		log.Fatal(err)
	}
	iSize := int(stat.Size())
	if iSize > 1024*1024*1014 {
		log.Fatal("input file size", iSize, "is > 1 GB")
	}

	// Read input file into a byte slice.
	data = make([]byte, iSize)
	_, err = bufio.NewReader(ih).Read(data)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return
}

// createPatternLineString writes pattern as "  n, b0, b1, ..., b(n-1), // AAA˙˙AA˙ " string.
func createPatternLineString(pattern []byte, maxPatternSize int) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%3d,", len(pattern))) // start line with pattern size
	for _, x := range pattern {                      // write pattern bytes
		s.WriteString(fmt.Sprintf(" 0x%02x,", x))
	}
	fill := spaces(6 * (maxPatternSize - len(pattern)))
	s.WriteString(fmt.Sprintf("%s // ", fill))               // align
	s.WriteString(byteSliceAsASCII(pattern, maxPatternSize)) // write pattern lettes as comment
	return s.String()                                        // no alignment here to keep s length
}
