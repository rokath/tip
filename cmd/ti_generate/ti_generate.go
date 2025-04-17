package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/rokath/tip/internal/pattern"
	"github.com/rokath/tip/internal/tiptable"
	"github.com/spf13/afero"

	"gopkg.in/neurosnap/sentences.v1/english"
)

var (
	version string // do not initialize, goreleaser will handle that
	commit  string // do not initialize, goreleaser will handle that
	date    string // do not initialize, goreleaser will handle that
	iFn     string // input file name
	oFn     string // ouput file name
	//tFn     string // token file/folder name
	//rFn     string // random file name
	help    bool
	verbose bool
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.StringVar(&iFn, "i", "", "input file/folder name")
	flag.StringVar(&oFn, "o", "idTable.c", "output file name")
	//flag.StringVar(&tFn, "t", "", "tokenizer file name")
	//flag.StringVar(&rFn, "r", "", "random file name")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()}
	doit(os.Stdout, fSys)
}

func doit(w io.Writer, fSys *afero.Afero) {
	flag.Parse()
	distributeArgs()
	if help {
		fmt.Fprintln(w, "Usage: ti_generate -i inputFileName [options]")
		fmt.Fprintln(w, "Example: `ti_generate -i trice.bin` creates idTable.c")
		flag.PrintDefaults()
		fmt.Fprintln(w, "The TipUserManual explains details.")
		return
	}
	if verbose {
		if version == "" && commit == "" && date == "" {
			fmt.Println("experimental version")
		} else {
			fmt.Fprintln(w, version, commit, date)
		}
	}

	/*
	if rFn != "" {
		fn := rFn+".txt"
		f, err := fSys.Create(fn)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		for range 10000 {
			f.WriteString(randSeq(6))
			f.WriteString("abcd")
			f.WriteString(randSeq(6))
			f.WriteString("abc")
			f.WriteString(randSeq(6))
			f.WriteString("ab")
			f.WriteString(randSeq(6))
			f.WriteString("RSTU")
			f.WriteString(randSeq(6))
			f.WriteString("123")
			f.WriteString(randSeq(6))
			f.WriteString("xy")
			f.WriteString(randSeq(6))
			f.WriteString("AAAAAA")
		}
		fmt.Println( fn, "created")
		return
	}

	if tFn != "" {
		tokenize(w, fSys, tFn)
		return
	}
	*/

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

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// tokenize is a wrapper for tokenizer.Tokenize(text)
func tokenize(w io.Writer, fSys *afero.Afero, tFn string) {
	folder := tFn + ".SAMPLES"
	if ok, _ := fSys.IsDir(folder); ok {
		log.Fatal(folder, " exists")
		return
	}
	err := fSys.Mkdir(folder, 0755)
	if err != nil {
		log.Fatal(err)
	}
	data, err := fSys.ReadFile(tFn)
	if err != nil {
		fmt.Fprintln(w, "oh no!")
		log.Fatal(err)
	}
	text := string(data)
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		fmt.Fprintln(w, "oh no?")
		panic(err)
	}

	sentences := tokenizer.Tokenize(text)

	for _, s := range sentences {
		t := strings.TrimSpace(s.Text)
		f, err := os.CreateTemp(folder+"", "*.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(t)
		//fmt.Printf("%3d:\t'%s'\n", i, t)
	}
}
