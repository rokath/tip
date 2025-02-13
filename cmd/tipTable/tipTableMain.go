package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rokath/tip/internal/pattern"
	"github.com/rokath/tip/internal/tiptable"
	"github.com/spf13/afero"
)

var (
	version    string // do not initialize, goreleaser will handle that
	commit     string // do not initialize, goreleaser will handle that
	date       string // do not initialize, goreleaser will handle that
	iFn        string // input file name
	oFn        string // ouput file name
	help       bool
	verbose    bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "v", false, "help")
	flag.StringVar(&iFn, "i", "", "input file name")
	flag.StringVar(&oFn, "o", "idTable.c", "output file name")
	flag.IntVar(&pattern.SizeMax, "z", 8, "max pattern size to find")
}

func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()}
	flag.Parse()
	doit(os.Stdout, fSys)
}

func doit(w io.Writer, fSys *afero.Afero) {
	distributeArgs()
	if help {
		fmt.Fprintln(w, "Usage: tipTable -i inputFileName [-o outputFileName] [-z max pattern size] [-v]")
		fmt.Fprintln(w, "Example: `tipTable -i trice.bin` creates idTable.c")
		flag.PrintDefaults()
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
	tiptable.Generate(fSys, oFn, iFn, pattern.SizeMax)
}

func distributeArgs() {
	tiptable.Verbose = verbose
	pattern.Verbose = verbose
}
