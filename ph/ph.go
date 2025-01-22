package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/spf13/afero"
)

var (
	version string // do not initialize, goreleaser will handle that
	commit  string // do not initialize, goreleaser will handle that
	date    string // do not initialize, goreleaser will handle that
	iFn     string // input file name
	oFn     = ".tipTable.go.txt"
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
		fmt.Fprintln(w, "Usage: ph -i inputFileName")
		fmt.Fprintln(w, "Example: `ph -i fn` creates fn"+oFn)
		return
	}

	ih, err := fSys.Open(iFn)
	if err != nil {
		log.Fatal(err)
	}
	defer ih.Close()

	// Get the file size
	stat, err := ih.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	iSize := int(stat.Size())

	stat.Name()

	// Read the file into a byte slice
	iData := make([]byte, iSize)
	_, err = bufio.NewReader(ih).Read(iData)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	pn := generateSortedPatternHistogram(iData)
	writeGoTipTable(fSys, oFn, pn, stat)

	//	offset := 0
	//	if len(pn) >= 128 {
	//		offset = len(pn) - 128
	//	}
	//	for i, x := range pn[offset:] {
	//		if i < 120 {
	//			fmt.Printf("%02x - %5d: ", 184-i, x.n)
	//		}
	//		if i >= 120 {
	//			fmt.Printf("%02x - %5d: ", 127-i, x.n)
	//		}
	//		dumpByteSlice(x.pattern)
	//	}
	//	fmt.Println(len(pn))

	// Now we have:
	// 52 -  1557:   FF FE FF                                          ...
	// 51 -  1557:   FF FF FE FF                                       ....
	// 50 -  1566:   FF FF FE                                          ...
	// 4f -  1618:   FB                                                .
	// 4e -  1632:   FC                                                .
	// 4d -  1678:   FE FF                                             ..
	// 4c -  1835:   FD                                                .
	// 4b -  1986:   FF FE                                             ..
	// 4a -  2216:   FE                                                .
	// 49 -  3039:   FF FF FF FF FF FF FF                              .......
	// 48 -  3127:   00 00                                             ..
	// 47 -  4288:   41 41 41 41 41 41 41 41                           AAAAAAAA
	// 46 -  4412:   41 41 41 41 41 41 41                              AAAAAAA
	// 45 -  4540:   41 41 41 41 41 41                                 AAAAAA
	// 44 -  4672:   41 41 41 41 41                                    AAAAA
	// 43 -  4808:   41 41 41 41                                       AAAA
	// 42 -  4948:   41 41 41                                          AAA
	// 41 -  5092:   41 41                                             AA
	// -----------------------------------------------------------------------------
	// 07 -  5392:   41                                                A
	// 06 -  5649:   FF FF FF FF FF FF                                 ......
	// 05 -  7567:   00                                                .
	// 04 -  8259:   FF FF FF FF FF                                    .....
	// 03 - 11687:   FF FF FF FF                                       ....
	// 02 - 19578:   FF FF FF                                          ...
	// 01 - 28018:   FF FF                                             ..
	// 00 - 39778:   FF

	// Generate all combinations of 00-07

	cn := CombineP0P7(pn[len(pn)-8:])

	for i, x := range cn {
		fmt.Printf("%02x - %5d: ", i, x.n)
		dumpByteSlice(x.pattern)
	}
	fmt.Println(len(pn))

	// Now we have:
	// * 8 single nibble pattern which we need to keep.
	// * + 120 byte pattern
	// * + 64 double nibble pattern which may occur also in the byte pattern.

}

func CombineP0P7(pn []nPatt) []nPatt {
	pt := pn[:8]
	comb := make([]nPatt, 64)
	idx := 0
	for _, x := range pt {
		for _, y := range pt {
			comb[idx].pattern = slices.Concat(x.pattern, y.pattern)
			comb[idx].n = x.n + y.n
			idx++
		}
	}
	return comb
}

// nPatt contains the count of a pattern and the pattern itself.
type nPatt struct {
	n       int    // n is the count of pattern occurances.
	pattern []byte // pattern are 1-8 bytes long.
}

// tipReplace contains the replacement byte of a pattern and the pattern itself.
type tipReplace struct {
	n       byte   // n is the count of pattern occurances.
	pattern []byte // pattern are 1-8 bytes long.
}

//writeB1CCode(oh, m1, k1)
//writeB2CCode(oh, m2, k2)
//writeB3CCode(oh, m3, k3)
//writeB4CCode(oh, m4, k4)
//writeB5CCode(oh, m5, k5)
//writeB6CCode(oh, m6, k6)
//writeB7CCode(oh, m7, k7)
//writeB8CCode(oh, m8, k8)
