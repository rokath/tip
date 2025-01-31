package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rokath/tip/internal/tiptable"
	"github.com/spf13/afero"
)

var (
	version    string // do not initialize, goreleaser will handle that
	commit     string // do not initialize, goreleaser will handle that
	date       string // do not initialize, goreleaser will handle that
	iFn        string // input file name
	oFn        string // input file name
	patSizeMax int
	help       bool
	verbose    bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "v", false, "help")
	flag.StringVar(&iFn, "i", "", "input file name")
	flag.StringVar(&oFn, "o", "tipTable.c", "output file name")
	flag.IntVar(&patSizeMax, "z", 8, "max pattern size to find")
}

func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()}
	flag.Parse()
	doit(os.Stdout, fSys)
}

func doit(w io.Writer, fSys *afero.Afero) {
	if help {
		fmt.Fprintln(w, "Usage: tipTable -i inputFileName [-o outputFileName] [-z max pattern size] [-v]")
		fmt.Fprintln(w, "Example: `tipTableGen -i trice.bin` creates tipTable.c")
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}
	if iFn == "" {
		fmt.Fprintln(w, `"tipTable -h" prints help`)
		return
	}
	if verbose {
		fmt.Fprintln(w, version, commit, date)
	}
	tiptable.Generate(fSys, oFn, iFn, patSizeMax)
}
