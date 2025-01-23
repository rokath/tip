package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

	// Get the input file size
	stat, err := ih.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	iSize := int(stat.Size())

	if iSize > 1024*1024*1014 {
		log.Fatal("input file size", iSize, "is > 1 GB")
	}

	// Read the file into a byte slice.
	iData := make([]byte, iSize)
	_, err = bufio.NewReader(ih).Read(iData)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	tips := generateSortedPatternHistogram(iData)
	writeGoTipTable(fSys, oFn, tips, stat)

}

// tip contains the count of a pattern and the pattern itself.
type tip struct {
	n int    // n is the count of pattern occurances.
	p []byte // pattern are 1-8 bytes long.
}
