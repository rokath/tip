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
	version string // do not initialize, goreleaser will handle that
	commit  string // do not initialize, goreleaser will handle that
	date    string // do not initialize, goreleaser will handle that
	iFn     string // input file name
	oFn     string // ouput file name
	help    bool
	verbose bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.StringVar(&iFn, "i", "", "input file/folder name")
	flag.StringVar(&oFn, "o", "idTable.c", "output file name")
}

func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()}
	doit(os.Stdout, fSys)
}

func doit(w io.Writer, fSys *afero.Afero) {
	flag.Parse()
	distributeArgs()
	if help {
		fmt.Fprintln(w, "Usage: ti_generate -i inputFileName [-o outputFileName] [-z max pattern size] [-v]")
		fmt.Fprintln(w, "Example: `ti_generate -i trice.bin` creates idTable.c")
		flag.PrintDefaults()
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}
	if verbose {
		if version == "" && commit == "" && date == "" {
			fmt.Println("experimenal version")
		} else {
			fmt.Fprintln(w, version, commit, date)
		}
	}
	if iFn == "" {
		if !verbose {
			fmt.Fprintln(w, `"ti_generate -h" prints help`)
		}
		return
	}
	err := tiptable.Generate(fSys, oFn, iFn, pattern.PatternSizeMax)
	if err != nil {
		fmt.Println(err)
	}
	if verbose {
		fmt.Println(oFn, "generated")
	}
}

func distributeArgs() {
	tiptable.Verbose = verbose
	pattern.Verbose = verbose
}
