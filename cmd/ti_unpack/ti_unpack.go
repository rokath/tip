package main

import (
	"encoding/hex"
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
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "v", false, "help")
	flag.StringVar(&iFn, "i", "", "input file name")
	flag.StringVar(&oFn, "o", "", "output file name")
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
		fmt.Fprintln(w, "Usage: ti_unpack -i inputFileName [-o outputFileName] [-v]")
		fmt.Fprintln(w, "Example: `ti_uunpack -i trice.bin.tip` creates trice.bin.tip.untip")
		flag.PrintDefaults()
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}
	if iFn == "" {
		fmt.Fprintln(w, `"ti_unpack -h" prints help`)
		return
	}
	if oFn == "" {
		oFn = iFn + ".untip"
	}
	if verbose {
		if version == "" && commit == "" && date == "" {
			fmt.Println("experimental version")
		} else {
			fmt.Fprintln(w, version, commit, date)
		}
	}
	fi, err := fSys.Stat(iFn)
	if err != nil {
		return
	}
	packet, err := fSys.ReadFile(iFn)
	if err != nil {
		return err
	}
	for i, x := range packet {
		if x == 0 {
			return fmt.Errorf("%s contains a 0 at offset %d (invalid tip packet)", iFn, i)
		}
	}
	buffer := make([]byte, 24*len(packet)) // assuming 24-bytes pattern matching exactly
	n := tip.Unpack(buffer, packet)
	if verbose {
		fmt.Println("packet:", hex.Dump(packet))
		fmt.Println("buffer:", hex.Dump(buffer[:n]))
		fmt.Fprintln(w, "file size", fi.Size(), "changed to", n, "(rate", 100*int64(n)/fi.Size(), "percent)")
	}
	return fSys.WriteFile(oFn, buffer[:n], 0644)
}
