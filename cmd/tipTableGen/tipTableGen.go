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
	oFn        = ".tipTable.c"
	patSizeMax int
	help       bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&iFn, "i", "-", "input file name")
	flag.IntVar(&patSizeMax, "z", 8, "max pattern size to find")
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
	if help {
		fmt.Fprintln(w, version, commit, date)
		fmt.Fprintln(w, "Usage: tipTableGen [-h] [-i inputFileName] [-z max pattern size]")
		fmt.Fprintln(w, "Example: `tipTableGen -i fileName -z 12` creates fileName"+oFn+" out of pattern with max size 12")
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}

	tiptable.Generate(fSys, oFn, iFn, patSizeMax)
}
