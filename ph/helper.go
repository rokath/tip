package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/afero"
)

// writeTipTable generates a file oFn containing Go code using pn and stat
func writeGoTipTable(fSys *afero.Afero, oFn string, pn []nPatt, stat os.FileInfo) {
	oh, err := fSys.Create(oFn)
	if err != nil {
		log.Fatal(err)
	}
	defer oh.Close()

	fmt.Fprintln(oh, `package main`)
	fmt.Fprintln(oh, "// Generated code - do not edit!")
	fmt.Fprintln(oh)
	fmt.Fprintf(oh, "var tipTable = [127]tipReplace{ // from %s (%s)\n", stat.Name(), stat.ModTime().String())
	for i, x := range pn[len(pn)-127:] {
		s := byteSliceAsGoCode(x.pattern)
		y := 52 - len(s)
		fmt.Fprintf(oh, "\t{ 0x%02x, %s }, %s// %s cnt = %5d\n", 0x7f-i, s, spaces(y), string(x.pattern), x.n)
	}
	fmt.Fprintf(oh, "}\n")
}

// https://gist.github.com/chmike/05da938833328a9a94e02506922f2e7b
func dumpByteSlice(b []byte) {
	var a [16]byte
	n := (len(b) + 15) &^ 15
	for i := 0; i < n; i++ {
		// if i%16 == 0 {
		// 	fmt.Printf("%4d", i)
		// }
		if i%8 == 0 {
			fmt.Print(" ")
		}
		if i < len(b) {
			fmt.Printf(" %02X", b[i])
		} else {
			fmt.Print("   ")
		}
		if i >= len(b) {
			a[i%16] = ' '
		} else if b[i] < 32 || b[i] > 126 {
			a[i%16] = '.'
		} else {
			a[i%16] = b[i]
		}
		if i%16 == 15 {
			fmt.Printf("  %s", string(a[:]))
		}
	}
	fmt.Println()
}

// byteSliceAsGoCode returns b as a Go code string. Evample:  []byte{ 0x5a, 0xf8, 0xbb}
func byteSliceAsGoCode(b []byte) string {
	var s strings.Builder
	s.WriteString("[]byte{")
	for i, x := range b {
		s.WriteString(fmt.Sprintf(" 0x%02x", x))
		if i < len(b)-1 {
			s.WriteString(",")
		}
	}
	s.WriteString("}")
	return s.String()
}

/*
// byteSliceAsASCII returns b as ASCII string. Example:  .Aah.B..C
func byteSliceAsASCII(b []byte) string {
	var s strings.Builder
	//s.WriteString("[]byte{")
	for i, x := range b {
		s.WriteString(fmt.Sprintf(" 0x%02x", x))
		if i < len(b)-1 {
			s.WriteString(",")
		}
	}
	s.WriteString("}")
	return s.String()
}
*/
// spaces returns a string consisting of n spaces.
func spaces(n int) (s string) {
	switch n {
	case 0:
		s = ""
	case 1:
		s = " "
	case 2:
		s = "  "
	case 3:
		s = "   "
	case 4:
		s = "    "
	case 5:
		s = "     "
	case 6:
		s = "      "
	case 7:
		s = "       "
	case 8:
		s = "        "
	case 9:
		s = "         "
	case 10:
		s = "          "
	case 11:
		s = "           "
	case 12:
		s = "            "
	case 13:
		s = "             "
	case 14:
		s = "              "
	case 15:
		s = "               "
	case 16:
		s = "                "
	case 17:
		s = "                 "
	case 18:
		s = "                  "
	case 19:
		s = "                   "
	case 20:
		s = "                    "
	case 21:
		s = "                     "
	case 22:
		s = "                      "
	case 23:
		s = "                       "
	case 24:
		s = "                        "
	case 25:
		s = "                         "
	case 26:
		s = "                          "
	case 27:
		s = "                           "
	case 28:
		s = "                            "
	case 29:
		s = "                             "
	case 30:
		s = "                              "
	case 31:
		s = "                               "
	case 32:
		s = "                                "
	case 33:
		s = "                                 "
	case 34:
		s = "                                  "
	case 35:
		s = "                                   "
	case 36:
		s = "                                    "
	case 37:
		s = "                                     "
	case 38:
		s = "                                      "
	case 39:
		s = "                                       "
	case 40:
		s = "                                        "
	case 41:
		s = "                                         "
	case 42:
		s = "                                          "
	case 43:
		s = "                                           "
	case 44:
		s = "                                            "
	case 45:
		s = "                                             "
	case 46:
		s = "                                              "
	case 47:
		s = "                                               "
	case 48:
		s = "                                                "
	case 49:
		s = "                                                 "
	case 50:
		s = "                                                  "
	case 51:
		s = "                                                   "
	case 52:
		s = "                                                    "
	case 53:
		s = "                                                     "
	case 54:
		s = "                                                      "
	case 55:
		s = "                                                       "
	case 56:
		s = "                                                        "
	case 57:
		s = "                                                         "
	case 58:
		s = "                                                          "
	case 59:
		s = "                                                           "
	case 60:
		s = "                                                            "
	case 61:
		s = "                                                             "
	case 62:
		s = "                                                              "
	case 63:
		s = "                                                               "
	case 64:
		s = "                                                                "
	case 65:
		s = "                                                                 "
	case 66:
		s = "                                                                  "
	case 67:
		s = "                                                                   "
	case 68:
		s = "                                                                    "
	case 69:
		s = "                                                                     "
	case 70:
		s = "                                                                      "
	case 71:
		s = "                                                                       "
	case 72:
		s = "                                                                        "
	case 73:
		s = "                                                                         "
	case 74:
		s = "                                                                          "
	case 75:
		s = "                                                                           "
	case 76:
		s = "                                                                            "
	case 77:
		s = "                                                                             "
	case 78:
		s = "                                                                              "
	case 79:
		s = "                                                                               "
	case 80:
		s = "                                                                                "
	}
	return s
}
