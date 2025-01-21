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
	oFn     = ".ph.c"
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

	oh, err := fSys.Create(oFn)
	if err != nil {
		log.Fatal(err)
	}
	defer oh.Close()

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

	fmt.Fprintln(oh, `//! \file`, oh.Name())
	fmt.Fprintln(oh, "//!")
	fmt.Fprintln(oh, "//! Generated code - do not edit!")
	fmt.Fprintln(oh)
	/*
		writeB1Histogram(iData, oh)
		writeB2Histogram(iData, oh)
		writeB3Histogram(iData, oh)
		writeB4Histogram(iData, oh)
		writeB5Histogram(iData, oh)
		writeB6Histogram(iData, oh)
		writeB7Histogram(iData, oh)
		writeB8Histogram(iData, oh)
	*/
	m1, k1 := createB1Histogram(iData)
	m2, k2 := createB2Histogram(iData)
	m3, k3 := createB3Histogram(iData)
	m4, k4 := createB4Histogram(iData)
	m5, k5 := createB5Histogram(iData)
	m6, k6 := createB6Histogram(iData)
	m7, k7 := createB7Histogram(iData)
	m8, k8 := createB8Histogram(iData)

	writeB1CCode(oh, m1, k1)
	writeB2CCode(oh, m2, k2)
	writeB3CCode(oh, m3, k3)
	writeB4CCode(oh, m4, k4)
	writeB5CCode(oh, m5, k5)
	writeB6CCode(oh, m6, k6)
	writeB7CCode(oh, m7, k7)
	writeB8CCode(oh, m8, k8)

}

/*
func writeB1Histogram(data []byte, fh afero.File) {
	m, keys := createB1Histogram( data )
	writeB1CCode(fh, m, keys)
}

func writeB2Histogram(data []byte, fh afero.File) {
	m, keys := createB2Histogram( data )
	writeB2CCode(fh, m, keys)
}


func writeB3Histogram(data []byte, fh afero.File) {
	m, keys := createB3Histogram( data )
	writeB3CCode(fh, m, keys)
}

func writeB4Histogram(data []byte, fh afero.File) {
	m, keys := createB4Histogram( data )
	writeB4CCode(fh, m, keys)
}

func writeB5Histogram(data []byte, fh afero.File) {
	m, keys := createB5Histogram( data )
	writeB5CCode(fh, m, keys)
}

func writeB6Histogram(data []byte, fh afero.File) {
	m, keys := createB6Histogram( data )
	writeB6CCode(fh, m, keys)
}

func writeB7Histogram(data []byte, fh afero.File) {
}

func writeB8Histogram(data []byte, fh afero.File) {
}
*/
