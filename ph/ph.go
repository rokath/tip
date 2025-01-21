package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

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

	writeB1Histogram(iData, oh)
	writeB2Histogram(iData, oh)
	writeB3Histogram(iData, oh)
	writeB4Histogram(iData, oh)
	writeB5Histogram(iData, oh)
	writeB6Histogram(iData, oh)
	writeB7Histogram(iData, oh)
	writeB8Histogram(iData, oh)

}

func writeB1Histogram(data []byte, fh afero.File) {
	// Create a histogram for 1-byte sequences.
	// The keys are the single bytes and the values are their occurance count.
	m := make(map[byte]int)
	for _, x := range data {
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort keys according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b1_t b1_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, byte} // char, idx")
	count := 0
	for i, x := range keys {
		count++
		y := x
		if x <= 0x20 || x >= 0x80 {
			y = ' '
		}
		if i < len(keys)-1 && count < 127 {
			if m[x] >= 100 {
				fmt.Fprintf(fh, "\t{%5d, 0x%02x }, // '%c',\t%3d\n", m[x], x, y, i)
			}
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x }  // '%c',\t%3d\n", m[x], x, y, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}

// func createB2Histogram(data []byte)(h map[[2]byte]int, k [][2]byte{
//
// }

func createKeyFromByteSlice(b []byte) string {
	return fmt.Sprintf("%q", b)
}

func createByteSliceFromKey(s string) []byte {
	// ...
}

func writeB2Histogram(data []byte, fh afero.File) {
	// Create a histogram for 2-byte sequences.
	// The keys are the 2-bytes sequences and the values are their occurance count.
	m := make(map[[2]byte]int)
	for i := 0; i < len(data)-1; i++ {
		x := [2]byte{data[i], data[i+1]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([][2]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b2_t b2_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, word} // chars, idx")
	count := 0
	for i, x := range keys {
		count++
		x0 := x[0]
		x1 := x[1]
		if x0 <= 0x20 || x0 >= 0x80 {
			x0 = ' '
		}
		if x1 <= 0x20 || x1 >= 0x80 {
			x1 = ' '
		}
		if i < len(keys)-1 && count < 127 {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x }, // '%c%c',\t%3d\n", m[x], x[0], x[1], x0, x1, i)
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x }  // '%c%c',\t%3d\n", m[x], x[0], x[1], x0, x1, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}

func writeB3Histogram(data []byte, fh afero.File) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m := make(map[[3]byte]int)
	for i := 0; i < len(data)-2; i++ {
		x := [3]byte{data[i], data[i+1], data[i+2]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([][3]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b3_t b3_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, bytes} // chars, idx")
	count := 0
	for i, x := range keys {
		count++
		x0 := x[0]
		x1 := x[1]
		x2 := x[2]
		if x0 <= 0x20 || x0 >= 0x80 {
			x0 = ' '
		}
		if x1 <= 0x20 || x1 >= 0x80 {
			x1 = ' '
		}
		if x2 <= 0x20 || x2 >= 0x80 {
			x2 = ' '
		}
		if i < len(keys)-1 && count < 127 {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x }, // '%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x0, x1, x2, i)
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x }  // '%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x0, x1, x2, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}

func writeB4Histogram(data []byte, fh afero.File) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m := make(map[[4]byte]int)
	for i := 0; i < len(data)-3; i++ {
		x := [4]byte{data[i], data[i+1], data[i+2], data[i+3]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([][4]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b4_t b4_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, bytes} // chars, idx")
	count := 0
	for i, x := range keys {
		count++
		x0 := x[0]
		x1 := x[1]
		x2 := x[2]
		x3 := x[3]
		if x0 <= 0x20 || x0 >= 0x80 {
			x0 = ' '
		}
		if x1 <= 0x20 || x1 >= 0x80 {
			x1 = ' '
		}
		if x2 <= 0x20 || x2 >= 0x80 {
			x2 = ' '
		}
		if x3 <= 0x20 || x3 >= 0x80 {
			x3 = ' '
		}
		if i < len(keys)-1 && count < 127 {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x }, // '%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x0, x1, x2, x3, i)
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x }  // '%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x0, x1, x2, x3, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}

func writeB5Histogram(data []byte, fh afero.File) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m := make(map[[5]byte]int)
	for i := 0; i < len(data)-4; i++ {
		x := [5]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([][5]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b5_t b5_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, bytes} // chars, idx")
	count := 0
	for i, x := range keys {
		count++
		x0 := x[0]
		x1 := x[1]
		x2 := x[2]
		x3 := x[3]
		x4 := x[4]
		if x0 <= 0x20 || x0 >= 0x80 {
			x0 = ' '
		}
		if x1 <= 0x20 || x1 >= 0x80 {
			x1 = ' '
		}
		if x2 <= 0x20 || x2 >= 0x80 {
			x2 = ' '
		}
		if x3 <= 0x20 || x3 >= 0x80 {
			x3 = ' '
		}
		if x4 <= 0x20 || x4 >= 0x80 {
			x4 = ' '
		}
		if i < len(keys)-1 && count < 127 {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x }, // '%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x0, x1, x2, x3, x4, i)
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x }  // '%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x0, x1, x2, x3, x4, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}

func writeB6Histogram(data []byte, fh afero.File) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m := make(map[[6]byte]int)
	for i := 0; i < len(data)-5; i++ {
		x := [6]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([][6]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b6_t b6_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, bytes} // chars, idx")
	count := 0
	for i, x := range keys {
		count++
		x0 := x[0]
		x1 := x[1]
		x2 := x[2]
		x3 := x[3]
		x4 := x[4]
		x5 := x[5]
		if x0 <= 0x20 || x0 >= 0x80 {
			x0 = ' '
		}
		if x1 <= 0x20 || x1 >= 0x80 {
			x1 = ' '
		}
		if x2 <= 0x20 || x2 >= 0x80 {
			x2 = ' '
		}
		if x3 <= 0x20 || x3 >= 0x80 {
			x3 = ' '
		}
		if x4 <= 0x20 || x4 >= 0x80 {
			x4 = ' '
		}
		if x5 <= 0x20 || x5 >= 0x80 {
			x5 = ' '
		}
		if i < len(keys)-1 && count < 127 {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x,0x%02x, 0x%02x }, // '%c%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x[5], x0, x1, x2, x3, x4, x5, i)
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x,0x%02x, 0x%02x }  // '%c%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x[5], x0, x1, x2, x3, x4, x5, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}

func writeB7Histogram(data []byte, fh afero.File) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m := make(map[[7]byte]int)
	for i := 0; i < len(data)-6; i++ {
		x := [7]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5], data[i+6]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([][7]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b7_t b7_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, bytes} // chars, idx")
	count := 0
	for i, x := range keys {
		count++
		x0 := x[0]
		x1 := x[1]
		x2 := x[2]
		x3 := x[3]
		x4 := x[4]
		x5 := x[5]
		x6 := x[6]
		if x0 <= 0x20 || x0 >= 0x80 {
			x0 = ' '
		}
		if x1 <= 0x20 || x1 >= 0x80 {
			x1 = ' '
		}
		if x2 <= 0x20 || x2 >= 0x80 {
			x2 = ' '
		}
		if x3 <= 0x20 || x3 >= 0x80 {
			x3 = ' '
		}
		if x4 <= 0x20 || x4 >= 0x80 {
			x4 = ' '
		}
		if x5 <= 0x20 || x5 >= 0x80 {
			x5 = ' '
		}
		if x6 <= 0x20 || x6 >= 0x80 {
			x6 = ' '
		}
		if i < len(keys)-1 && count < 127 {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x }, // '%c%c%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x[5], x[6], x0, x1, x2, x3, x4, x5, x6, i)
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x }  // '%c%c%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x[5], x[6], x0, x1, x2, x3, x4, x5, x6, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}

func writeB8Histogram(data []byte, fh afero.File) {
	// Create a histogram for 3-byte sequences.
	// The keys are the 3-bytes sequences and the values are their occurance count.
	m := make(map[[8]byte]int)
	for i := 0; i < len(data)-7; i++ {
		x := [8]byte{data[i], data[i+1], data[i+2], data[i+3], data[i+4], data[i+5], data[i+6], data[i+7]}
		if n, ok := m[x]; ok {
			m[x] = n + 1
		} else {
			m[x] = 1
		}
	}

	// Get a list of all m keys.
	keys := make([][8]byte, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort m according to their values.
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	fmt.Fprintln(fh, `b8_t b8_ph = {`)
	fmt.Fprintln(fh, "\t// {cnt, bytes} // chars, idx")
	count := 0
	for i, x := range keys {
		count++
		x0 := x[0]
		x1 := x[1]
		x2 := x[2]
		x3 := x[3]
		x4 := x[4]
		x5 := x[5]
		x6 := x[6]
		x7 := x[7]
		if x0 <= 0x20 || x0 >= 0x80 {
			x0 = ' '
		}
		if x1 <= 0x20 || x1 >= 0x80 {
			x1 = ' '
		}
		if x2 <= 0x20 || x2 >= 0x80 {
			x2 = ' '
		}
		if x3 <= 0x20 || x3 >= 0x80 {
			x3 = ' '
		}
		if x4 <= 0x20 || x4 >= 0x80 {
			x4 = ' '
		}
		if x5 <= 0x20 || x5 >= 0x80 {
			x5 = ' '
		}
		if x6 <= 0x20 || x6 >= 0x80 {
			x6 = ' '
		}
		if x7 <= 0x20 || x7 >= 0x80 {
			x7 = ' '
		}
		if i < len(keys)-1 && count < 127 {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x }, // '%c%c%c%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7], x0, x1, x2, x3, x4, x5, x6, x7, i)
		} else {
			fmt.Fprintf(fh, "\t{%5d, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x, 0x%02x }  // '%c%c%c%c%c%c%c%c',\t%3d\n", m[x], x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7], x0, x1, x2, x3, x4, x5, x6, x7, i)
			fmt.Fprintf(fh, "};\n\n")
			break
		}
	}
}
