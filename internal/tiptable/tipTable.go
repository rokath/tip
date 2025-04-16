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
	flag.IntVar(&ID1Count, "n", 127, "direct ID count ID1Count (count for 2-bytes pattern), 0-127 for u=7 and 0-191 for u=6")
}

func PrintPattern(index int, x pattern.Pattern) {
	s := make([]byte, 32)
	for _, b := range x.Bytes {
		if !(20 <= b && b <= 127) {
			b = ' '
		}
		s = append(s, b)
	}
	//fmt.Printf("cnt:%4d w:%9.1f b:%8.2f rateD:%8.4f rateI:%8.4f hex:%16s, ascci:'%s'\n", len(x.Pos), x.Weight, x.Balance, 1000*x.RateDirect, 1000*x.RateIndirect, hex.EncodeToString(x.Bytes), string(s))
	fmt.Printf("i:%3d, cnt:%6d, ascci:'%s'\n", index, len(x.Pos), string(s))
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
	discardSize := 100 // len(p.Hist)/100
	p.DiscardSeldomPattern(discardSize)
	//p.ComputeValues(maxPatternSize)
	p.ReduceFromSmallerSide()
	//p.ReduceFromLargerSide()
	//p.ComputeValues(maxPatternSize)
	//p.DeleteEmptyKeys()

	//p.PrintInfo("Histogram after Reduce")
	//p.PrintInfo("Histogram after AddWeights")

	// Todo: Reduce bigger keys if smaller keys fit?
	ll := 50
	list := p.ExportAsList()
	pattern.SortByDescCount(list)
	fmt.Println("SortByDescCount")
	for i, x := range list {
		PrintPattern(i, x)
		if i == ll {
			break
		}
	}
	//pattern.SortByDescWeight(list)
	//fmt.Println("SortByDescWeight")
	//for i, x := range list {
	//	PrintPattern(x)
	//	if i == ll {
	//		break
	//	}
	//}
	// pattern.SortByDescBalance(list)
	// fmt.Println("SortByDescBalance")
	// for i, x := range list {
	// 	PrintPattern(x)
	// 	if i == ll {
	// 		break
	// 	}
	// }
	//  pattern.SortByIncrRateDirect(list)
	//  fmt.Println("SortByIncrRateDirect (i len weight balance RateDirect RateIndirect pattern)")
	//  for i, x := range list {
	//  	PrintPattern(x)
	//  	if i == ll {
	//  		break
	//  	}
	//  }
	//  pattern.SortByIncrRateIndirect(list)
	//  fmt.Println("SortByIncrRateIndirect (i len weight balance RateDirect RateIndirect pattern)")
	//  for i, x := range list {
	//  	PrintPattern(x)
	//  	if i == ll {
	//  		break
	//  	}
	//  }

	idList, moreBytes, maxListPatternSize := pattern.Extract2BytesPatterns(list)
	//moreBytes = list
	// lists are sorted by descending count here.
	if len(idList) >= ID1Count {
		idList = idList[:ID1Count]
	}
	// idList contains now up to ID1Count 2-bytes pattern. Remaining 2-bytes patterns are discarded
	// because they give nearly no compression effect when indexed with an indirect ID.

	moreBytesCount := MaxID - ID1Count
	if len(moreBytes) > moreBytesCount {
		moreBytes = moreBytes[:moreBytesCount]
	} else {
		fmt.Printf("warning:more pattern ID space than pattern (LastID %d < MaxID %d)", len(moreBytes), MaxID)
	}
	//for _, x := range moreBytes {
	//	idList = append(idList, x)
	//}
	idList = append(idList, moreBytes...)

	// indirectIndexedCount := min(MaxID, len(moreBytes))
	// idList := pattern.SortByDescLength(list[:idCount])

	//maxListPatternSize := len(idList[ID1Count].Bytes)
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

//! UnreplacableContainerBits is container bit size for unreplacebale bytes.
const unsigned unreplacableContainerBits = %d; // 6 bits or 7 bits
`, UnreplacableContainerBits)

	fmt.Fprintf(oh, `
//! ID1Max is the max possible number of primary IDs. Its value depends on UnreplacableContainerBits.
const unsigned ID1Max = %d; // 7 bits:127 or 6 bits:191
`, ID1Max)

	fmt.Fprintf(oh, `
//! ID1Count is the direct ID count. The resulting indirect ID count is (ID1Max - ID1Count) * 255.
const unsigned ID1Count = %d;
`, ID1Count)

	fmt.Fprintf(oh, `
//! MaxID is a computed value: MaxID = ID1Count + (ID1Max - ID1Count) * 255.
//! It is the max possible amount of pattern in the idTable.
const unsigned MaxID = %d;
`, MaxID)

	fmt.Fprintf(oh, `
//! LastID is pattern count inside the idTable. If it is < MaxID, consider increasing ID1Count.
const unsigned LastID = %d;
`, len(idList))

	fmt.Fprintf(oh, `
//! maxPatternlength is the size of the longest pattern inside idTable.
const uint8_t maxPatternlength = %d;
`, pattern.PatternSizeMax)

	fmt.Fprintln(oh, `
//! idTable is sorted by pattern length and pattern count.
//! The pattern position + 1 is the replace id.
//! The generator pattern max size was`, maxPatternSize, `and the list pattern max size is:`, maxListPatternSize)
	start := fmt.Sprintf("const uint8_t idTable[] = { // from %s\n", loc)
	fill := spaces(len("    xxx, ") + len("0x00, ")*maxListPatternSize)
	fill2 := spaces(maxListPatternSize - len("ASCII"))
	fmt.Fprintf(oh, start+"%s// `ASCII%s`|  count    id (decimal)  id1  id2\n", fill, fill2)
	for i, x := range idList {
		pls := createPatternLineString(x.Bytes, maxListPatternSize) // todo: review and improve code
		sz := len(x.Bytes)
		tipTableSize += 1 + sz
		if i < len(idList) {
			id1, id2 := tipPackageIDs(i + 1)
			if id2 == -1 {
				fmt.Fprintf(oh, "\t%s|%7d  0x%04x (%5d)   %02x   --\n", pls, len(x.Pos), i+1, i+1, id1)
			} else {
				fmt.Fprintf(oh, "\t%s|%7d  0x%04x (%5d)   %02x   %02x\n", pls, len(x.Pos), i+1, i+1, id1, uint8(id2))
			}
		}
	}
	fmt.Fprintln(oh, "\t  0 // table end marker")
	fmt.Fprintln(oh, "};")
	fmt.Fprintln(oh)
	fmt.Fprintf(oh, "// tipTableSize is %d.\n", tipTableSize)
	fmt.Fprintln(oh)

	fmt.Fprint(oh, "// Informal, here are all pattern occuring at least twice:\n")
	for i, x := range list {
		if len(x.Pos) > 1 {
			pls := createPatternLineString(x.Bytes, maxListPatternSize) // todo: review and improve code
			fmt.Fprintf(oh, "//%4d: (%4d) %s\n", i, len(x.Pos), pls)
		}
	}
	fmt.Fprintln(oh)
	return
}

// tipPackageIDs computes id1 and id2 from id = 1...MaxID.
// id1 range is 1...ID1Count
// id2 range is 1...255, when -1, than no id2
func tipPackageIDs(id int) (id1 uint8, id2 int) {
	if id <= ID1Count {
		id1 = uint8(id)
		id2 = -1
		return
	}
	offs := ID1Count + 1
	level := (id - offs) / 255
	id2 = (id-offs)%255 + 1
	id1 = uint8(offs + level)

	// cross check:
	idx := (int(id1)-offs)*255 + id2 - 1 + offs
	if id != idx {
		fmt.Printf("ERROR: id:%5d, id1:0x%02x, id2:0x%02x, id != %5d !!!\n", id, id1, id2, idx)
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
