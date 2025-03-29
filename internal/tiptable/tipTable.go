package tiptable

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/rokath/tip/internal/pattern"
	"github.com/spf13/afero"
)

var (
	Verbose                   bool
	ID1Max                    int // ID1Max results immediately from UnreplacableContainerBits.
	ID1Count                  int // ID1Count its the direct IDs count.
	MaxID                     int // MaxID is the max possible amount of pattern in the idTable.
	UnreplacableContainerBits = 6 // UnreplacableContainerBits is container bit size for unreplacebale bytes.
)

func init() {
	flag.IntVar(&UnreplacableContainerBits, "u", UnreplacableContainerBits, "unreplacable bytes container bit size (6 or 7)")
	flag.IntVar(&ID1Count, "n", 127, "direct ID count ID1Count, 0-127 for u=7 and 0-191 for u=7")
}

// Generate writes a file oFn containing C code using loc file(s) and max pattern size.
// https://en.wikipedia.org/wiki/Dictionary_coder
// https://cs.stackexchange.com/questions/112901/algorithm-to-find-repeated-patterns-in-a-large-string
func Generate(fSys *afero.Afero, oFn, loc string, maxPatternSize int) (err error) {

	if UnreplacableContainerBits == 7 {
		ID1Max = 127
	} else if UnreplacableContainerBits == 6 {
		ID1Max = 191
	} else {
		log.Fatalf("Invalid value %d for UnreplacableContainerBits", UnreplacableContainerBits)
	}
	MaxID = ID1Count + (ID1Max-ID1Count)*255

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
	idList, moreBytes := pattern.Extract2BytesPatterns(rlist)
	// lists are sorted by descending count here.
	if len(idList) >= ID1Count {
		idList = idList[:ID1Count]
	}
	// idList contains now up to ID1Count 2-bytes pattern. Remaining 2-bytes patterns are discarded
	// because they give nearly no compression effect when indexed with an indirect ID.
    
	moreBytesCount := MaxID-ID1Count
	moreBytes = moreBytes[:moreBytesCount]
	for _, x := range moreBytes {
		idList = append( idList, x)
	}

	// indirectIndexedCount := min(MaxID, len(moreBytes))
	// idList := pattern.SortByDescLength(list[:idCount])

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
#include <stddef.h>`)

	fmt.Fprintf(oh, `

// UnreplacableContainerBits is container bit size for unreplacebale bytes.
const unsigned unreplacableContainerBits = %d; // 6 bits or 7 bits
`, UnreplacableContainerBits)

	fmt.Fprintf(oh, `
// ID1Max is the max possible number of primary IDs. Its value depends on UnreplacableContainerBits.
const unsigned ID1Max = %d; // 7 bits:127 or 6 bits:191
`, ID1Max)

	fmt.Fprintf(oh, `
// ID1Count is the direct ID count. The resulting indirect ID count is (ID1Max - ID1Count) * 255.
const unsigned ID1Count = %d;
`, ID1Count)

	fmt.Fprintf(oh, `
// MaxID is a computed value: MaxID = ID1Count + (ID1Max - ID1Count) * 255.
// It is the max possible amount of pattern in the idTable.
const unsigned MaxID = %d;
`, MaxID)

	fmt.Fprintf(oh, `
// LastID is pattern count inside the idTable. If it is < MaxID, consider increasing ID1Count.
const unsigned LastID = %d;
`, len(idList))

	fmt.Fprintln(oh, `
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
		if i < len(idList) {
			fmt.Fprintf(oh, "\t%s|%7d  %02x\n", pls, x.Cnt, i+1)
		}
	}
	fmt.Fprintln(oh, "\t  0 // table end marker")
	fmt.Fprintln(oh, "};")
	fmt.Fprintln(oh)
	fmt.Fprintf(oh, "// tipTableSize is %d.\n", tipTableSize)
	fmt.Fprintln(oh)

	fmt.Fprint(oh, "// Informal, here are all pattern occuring at least twice:\n")
	for i, x := range list {
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
