package tiptable

import (
	"flag"
	"fmt"
	"log"
	"slices"
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
	fmt.Printf("i:%3d, weight:%8d, cnt:%6d, ascci:'%s'\n", index, len(x.Pos)*len(x.Bytes), len(x.Pos), string(s))
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
	if ID1Count > ID1Max {
		log.Fatalf("Invalid value %d for ID1Count (must not be bigger than ID1Max=%d)", ID1Count, ID1Max)
	}
	MaxID = ID1Count + (ID1Max-ID1Count)*255

	var m sync.Mutex
	p := pattern.NewHistogram(&m)

	if ok, _ := fSys.IsDir(loc); ok {
		err = p.ScanAllFiles(fSys, loc, maxPatternSize)
	} else {
		err = p.ScanFile(fSys, loc, maxPatternSize)
	}

	list := p.ExportAsList()
	pattern.SortByDescWeight(list)
	idList1 := list[:ID1Count]                                                   // direct IDs
	idList2 := slices.DeleteFunc(list[ID1Count:], func(x pattern.Pattern) bool { // indirect IDs
		return len(x.Bytes) <= 2 // 2-bytes patteren make no sense for indirect IDs
	})
	cList := append(idList1, idList2...) // combined list
	var idList []pattern.Pattern         // used ID
	if len(cList) > MaxID {
		idList = cList[:MaxID] // used part
	} else {
		fmt.Printf("idList len %d is smaller than MaxID %d - consider parameter change\n", len(cList), MaxID)
		idList = cList
	}
	var maxListPatternSize int
	for _, x := range idList {
		if len(x.Bytes) > maxListPatternSize {
			maxListPatternSize = len(x.Bytes)
		}
	}
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
	switch maxListPatternSize {
	case 2:
		fill = fill[:len(fill)-3]
	case 3:
		fill = fill[:len(fill)-2]
	case 4:
		fill = fill[:len(fill)-1]
	}
	fill2 := spaces(maxListPatternSize - len("ASCII"))
	fmt.Fprintf(oh, start+"%s// `ASCII%s`|  count   weight   id (decimal)  id1  id2\n", fill, fill2)
	for i, x := range idList {
		pls := createPatternLineString(x.Bytes, maxListPatternSize) // todo: review and improve code
		sz := len(x.Bytes)
		tipTableSize += 1 + sz
		if i < len(idList) {
			id1, id2 := tipPackageIDs(i + 1)
			w := len(x.Pos) * len(x.Bytes)
			if id2 == -1 {
				fmt.Fprintf(oh, "\t%s|%7d %8d 0x%04x (%5d)   %02x   --\n", pls, len(x.Pos), w, i+1, i+1, id1)
			} else {
				fmt.Fprintf(oh, "\t%s|%7d %8d 0x%04x (%5d)   %02x   %02x\n", pls, len(x.Pos), w, i+1, i+1, id1, uint8(id2))
			}
		}
	}
	fmt.Fprintln(oh, "\t  0 // table end marker")
	fmt.Fprintln(oh, "};")
	fmt.Fprintln(oh)
	fmt.Fprintf(oh, "// tipTableSize is %d.\n", tipTableSize)
	fmt.Fprintln(oh)

	if Verbose {
		fmt.Fprint(oh, "// Informal, here are all by weight sorted pattern occuring at least twice:\n")
		fmt.Fprint(oh, "//   index: ( count) len, bytes... // ASCII\n")
		for i, x := range list {
			if len(x.Pos) > 1 {
				pls := createPatternLineString(x.Bytes, maxListPatternSize) // todo: review and improve code
				fmt.Fprintf(oh, "//%8d: (%6d) %s\n", i, len(x.Pos), pls)
			}
		}
		fmt.Fprintln(oh)
	}
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
	s.WriteString("`")
	for _, x := range b {
		if 0x20 <= x && x < 0x7f {
			s.WriteString(fmt.Sprintf("%c", x))
		} else {
			s.WriteString(`˙`)
		}
	}
	s.WriteString("`")
	s.WriteString(spaces(length - len(b)))
	return s.String()
}
