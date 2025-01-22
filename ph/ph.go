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
	// oFn     = ".ph.c.txt"
)

func init() {
	flag.StringVar(&iFn, "i", "-", "input file name")
	flag.Parse()
	//if iFn != "-" {
	//	oFn = iFn + oFn
	//}
}
func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()}
	doit(os.Stdout, fSys)
}

func doit(w io.Writer, fSys *afero.Afero) {

	if len(os.Args) != 3 {
		fmt.Fprintln(w, version, commit, date)
		fmt.Fprintln(w, "Usage: ph -i inputFileName")
		//fmt.Fprintln(w, "Example: `ph -i fn` creates fn"+oFn)
		return
	}

	ih, err := fSys.Open(iFn)
	if err != nil {
		log.Fatal(err)
	}
	defer ih.Close()

	// oh, err := fSys.Create(oFn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer oh.Close()

	// Get the file size
	stat, err := ih.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	iSize := int(stat.Size())

	// Read the file into a byte slice
	iData := make([]byte, iSize)
	_, err = bufio.NewReader(ih).Read(iData)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	// fmt.Fprintln(oh, `//! \file`, oh.Name())
	// fmt.Fprintln(oh, "//!")
	// fmt.Fprintln(oh, "//! Generated code - do not edit!")
	// fmt.Fprintln(oh)

	m1, k1 := createB1Histogram(iData)
	m2, k2 := createB2Histogram(iData)
	m3, k3 := createB3Histogram(iData)
	m4, k4 := createB4Histogram(iData)
	m5, k5 := createB5Histogram(iData)
	m6, k6 := createB6Histogram(iData)
	m7, k7 := createB7Histogram(iData)
	m8, k8 := createB8Histogram(iData)
	// Maps m1...m8 contain the pattern histograms.
	// Key slices k1...k8 are the appropriate keys ordered by size downwards.

	type nPatt struct {
		n       int    // n is the count of pattern occurances.
		pattern []byte // pattern are 1-8 bytes long.
	}

	pn := make([]nPatt, 0, 1024)

	for i := 0; i < len(k1); i++ {
		n1 := m1[k1[i]]
		pn = append(pn, nPatt{n1, []byte{k1[i]}})
	}
	for i := 0; i < len(k2); i++ {
		n := m2[k2[i]]
		pn = append(pn, nPatt{n, k2[i][:]})
	}
	for i := 0; i < len(k3); i++ {
		n := m3[k3[i]]
		pn = append(pn, nPatt{n, k3[i][:]})
	}
	for i := 0; i < len(k4); i++ {
		n := m4[k4[i]]
		pn = append(pn, nPatt{n, k4[i][:]})
	}
	for i := 0; i < len(k5); i++ {
		n := m5[k5[i]]
		pn = append(pn, nPatt{n, k5[i][:]})
	}
	for i := 0; i < len(k6); i++ {
		n := m6[k6[i]]
		pn = append(pn, nPatt{n, k6[i][:]})
	}
	for i := 0; i < len(k7); i++ {
		n := m7[k7[i]]
		pn = append(pn, nPatt{n, k7[i][:]})
	}
	for i := 0; i < len(k8); i++ {
		n := m8[k8[i]]
		pn = append(pn, nPatt{n, k8[i][:]})
	}

	/* sort pn for count */

	compareFn := func(a, b nPatt) int {
		if a.n > b.n {
			return 1
		}
		if a.n < b.n {
			return -1
		}
		return 0
	}

	slices.SortFunc(pn, compareFn)

	offset := 0
	if len(pn) >= 128 {
		offset = len(pn) - 128
	}
	for i, x := range pn[offset:] {
		if i < 120 {
			fmt.Printf("%02x - %5d: ", 184-i, x.n)
		}
		if i >= 120 {
			fmt.Printf("%02x - %5d: ", 127-i, x.n)
		}
		dumpByteSlice(x.pattern)
	}
	fmt.Println(len(pn))

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

}

// https://gist.github.com/chmike/05da938833328a9a94e02506922f2e7b
func dumpByteSlice(b []byte) {
	var a [16]byte
	n := (len(b) + 15) &^ 15
	for i := 0; i < n; i++ {
		// if i%16 == 0 {
		// 	fmt.Printf("%4d", i)
		// }
		if i%8 == 0 {
			fmt.Print(" ")
		}
		if i < len(b) {
			fmt.Printf(" %02X", b[i])
		} else {
			fmt.Print("   ")
		}
		if i >= len(b) {
			a[i%16] = ' '
		} else if b[i] < 32 || b[i] > 126 {
			a[i%16] = '.'
		} else {
			a[i%16] = b[i]
		}
		if i%16 == 15 {
			fmt.Printf("  %s\n", string(a[:]))
		}
	}
}

//writeB1CCode(oh, m1, k1)
//writeB2CCode(oh, m2, k2)
//writeB3CCode(oh, m3, k3)
//writeB4CCode(oh, m4, k4)
//writeB5CCode(oh, m5, k5)
//writeB6CCode(oh, m6, k6)
//writeB7CCode(oh, m7, k7)
//writeB8CCode(oh, m8, k8)
