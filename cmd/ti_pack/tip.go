package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rokath/tip/pkg/tip"
	"github.com/spf13/afero"
)

var (
	version string // do not initialize, goreleaser will handle that
	commit  string // do not initialize, goreleaser will handle that
	date    string // do not initialize, goreleaser will handle that
	iFn     string // input file name
	oFn     string // output file name
	help    bool
	verbose bool
	maxSize int
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.StringVar(&iFn, "i", "", "input file name")
	flag.StringVar(&oFn, "o", "", "output file name")
	flag.IntVar(&maxSize, "m", 255, "max file size according to tip.h configuration")
}

func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()}
	flag.Parse()
	err := doit(os.Stdout, fSys)
	if err != nil {
		fmt.Println(err)
	}
}

func doit(w io.Writer, fSys *afero.Afero) (err error) {
	if help {
		fmt.Fprintln(w, "Usage: ti_pack -i inputFileName [-o outputFileName] [-m max file size] [-v]")
		fmt.Fprintln(w, "Example: `ti_pack -i trice.bin` creates trice.bin.tip")
		flag.PrintDefaults()
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}
	if iFn == "" {
		fmt.Fprintln(w, `"ti_pack -h" prints help`)
		return
	}
	if oFn == "" {
		oFn = iFn + ".tip"
	}
	//  if verbose {
	//  	fmt.Fprintln(w, version, commit, date)
	//  }
	fi, err := fSys.Stat(iFn)
	if err != nil {
		return
	}
	if fi.Size() > int64(maxSize) {
		return fmt.Errorf("cannot pack %d bytes. maximum is %d", fi.Size(), maxSize)
	}
	buffer, err := fSys.ReadFile(iFn)
	if err != nil {
		return
	}
	packet := make([]byte, 2*len(buffer))
	n := tip.Pack(packet, buffer)
	if verbose {
		//fmt.Println(hex.Dump(buffer))
		//fmt.Println(hex.Dump(packet[:n]))
		fmt.Fprintln(w, "file size", fi.Size(), "changed to", n, "(rate", 100*n/len(buffer), "percent)")
	}
	return fSys.WriteFile(oFn, packet[:n], 0644)
}
