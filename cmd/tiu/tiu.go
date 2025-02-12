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
		fmt.Fprintln(w, "Usage: tiu -i inputFileName [-o outputFileName] [-v]")
		fmt.Fprintln(w, "Example: `tiu -i trice.bin` creates trice.bin.untip")
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}
	if iFn == "" {
		fmt.Fprintln(w, `"tiu -h" prints help`)
		return
	}
	if oFn == "" {
		oFn = iFn + ".untip"
	}
	if verbose {
		fmt.Fprintln(w, version, commit, date)
	}

	packet, err := fSys.ReadFile(iFn)
	if err != nil {
		return err
	}
	buffer := make([]byte, 24*len(packet)) // assuming 24-bytes pattern matching exactly
	n := tip.Unpack(buffer, packet)
	if verbose {
		fmt.Println(hex.Dump(packet))
		fmt.Println(hex.Dump(buffer[:n]))
		fmt.Println("Unpack rate is", 100*n/len(packet), "percent.")
	}

	return fSys.WriteFile(oFn, buffer[:n], 0644)
}
