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

	m1, k1 := createB1Histogram(iData)
	m2, k2 := createB2Histogram(iData)
	m3, k3 := createB3Histogram(iData)
	m4, k4 := createB4Histogram(iData)
	m5, k5 := createB5Histogram(iData)
	m6, k6 := createB6Histogram(iData)
	m7, k7 := createB7Histogram(iData)
	m8, k8 := createB8Histogram(iData)
	/*
		type struct{
			count int
			pattern []byte
		}


		for i:= 0; i < len(k1); i++ {
			n1, _ := m1[k1[i]]
			n2, _ := m2[k2[i]]
			n3, _ := m3[k3[i]]
			n4, _ := m4[k4[i]]
			n5, _ := m5[k5[i]]
			n6, _ := m6[k6[i]]
			n7, _ := m7[k7[i]]
			n8, _ := m8[k8[i]]
		}
	*/

	writeB1CCode(oh, m1, k1)
	writeB2CCode(oh, m2, k2)
	writeB3CCode(oh, m3, k3)
	writeB4CCode(oh, m4, k4)
	writeB5CCode(oh, m5, k5)
	writeB6CCode(oh, m6, k6)
	writeB7CCode(oh, m7, k7)
	writeB8CCode(oh, m8, k8)

}
