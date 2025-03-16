package tiptable

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/rokath/tip/internal/pattern"
	"github.com/spf13/afero"
)

var (
	Verbose bool
)

// Generate writes a file oFn containing C code using loc file(s) and max pattern size.
// https://en.wikipedia.org/wiki/Dictionary_coder
// https://cs.stackexchange.com/questions/112901/algorithm-to-find-repeated-patterns-in-a-large-string
func Generate(fSys *afero.Afero, oFn, loc string, maxPatternSize int) (err error) {
	var m sync.Mutex
	p := pattern.NewHistogram(&m)

	if ok, _ := fSys.IsDir(loc); ok {
		err = p.ScanAllFiles(fSys, loc, maxPatternSize)
	} else {
		err = p.ScanFile(fSys, loc, maxPatternSize)
	}
	if Verbose {
		p.PrintInfo("Histogram after Scan")
	}
	// All these trials did not result in significantly improved
	//p.DiscardSeldomPattern(10)
	//p.PrintInfo("Histogram after DiscardSeldomPattern")
	//p.BalanceByteUsage(maxPatternSize)
	//p.PrintInfo("Histogram after Balance")
	//p.Reduce()
	//p.DeleteEmptyKeys()
	//p.PrintInfo("Histogram after Reduce")
	//p.AddWeigths()
	//p.PrintInfo("Histogram after AddWeights")

	// Todo: Reduce bigger keys if smaller keys fit?

	rlist := p.ExportAsList()

	list := pattern.SortByDescCountDescLength(rlist)

	idCount := min(127, len(list))
	idList := pattern.SortByDescLength(list[:idCount])
	maxListPatternSize := len(idList[0].Bytes)
	oh, err := fSys.Create(oFn)
	if err != nil {
		log.Fatal(err)
	}
	defer oh.Close()
	tipTableSize := 1 // TipTable contains at least table end marker
	fmt.Fprintln(oh, `//! @file idTable.c
//! @brief Generated code - do not edit!

#include <stdint.h>
#include <stddef.h>

//! idTable is sorted by pattern length and pattern count.
//! The pattern position + 1 is the replace id.
//! The generator pattern max size was`, maxPatternSize, `and the list pattern max size is:`, maxListPatternSize)
	start := fmt.Sprintf("const uint8_t idTable[] = { // from %s\n", loc)
	fill := spaces(len("    xxx, ") + len("0x00, ")*maxListPatternSize)
	fill2 := spaces(maxListPatternSize - len("ASCII"))
	fmt.Fprintf(oh, start+"%s// `ASCII%s`|  count  id\n", fill, fill2)
	for i, x := range idList {
		pls := createPatternLineString(x.Bytes, maxListPatternSize) // todo: review and improve code
		sz := len(x.Bytes)
		tipTableSize += 1 + sz
		if i < idCount {
			fmt.Fprintf(oh, "\t%s|%7d  %02x\n", pls, x.Cnt, i+1)
		}
	}
	fmt.Fprintln(oh, "\t  0 // table end marker")
	fmt.Fprintln(oh, "};")
	fmt.Fprintln(oh)
	fmt.Fprintf(oh, "// tipTableSize is %d.\n", tipTableSize)
	fmt.Fprintln(oh)
	for i, x := range list {
		if i == 127 {
			fmt.Fprintln(oh, "// --------------------------------")
		}
		if x.Cnt > 1 {
			pls := createPatternLineString(x.Bytes, maxListPatternSize) // todo: review and improve code
			fmt.Fprintf(oh, "//%4d: (%4d) %s\n", i, x.Cnt, pls)
		}
	}
	fmt.Fprintln(oh)
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
