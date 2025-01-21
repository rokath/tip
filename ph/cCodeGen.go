package main

import (
	"fmt"

	"github.com/spf13/afero"
)

// generateB1CCode writes B1 code into fh.
func writeB1CCode(fh afero.File, m map[byte]int, keys []byte) {
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

// generateB2CCode writes B2 code into fh.
func writeB2CCode(fh afero.File, m map[[2]byte]int, keys [][2]byte) {
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

// generateB3CCode writes B3 code into fh.
func writeB3CCode(fh afero.File, m map[[3]byte]int, keys [][3]byte) {
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

// generateB4CCode writes B4 code into fh.
func writeB4CCode(fh afero.File, m map[[4]byte]int, keys [][4]byte) {
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

// generateB5CCode writes B5 code into fh.
func writeB5CCode(fh afero.File, m map[[5]byte]int, keys [][5]byte) {
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

// generateB6CCode writes B6 code into fh.
func writeB6CCode(fh afero.File, m map[[6]byte]int, keys [][6]byte) {
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

// generateB7CCode writes B7 code into fh.
func writeB7CCode(fh afero.File, m map[[7]byte]int, keys [][7]byte) {
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

// generateB8CCode writes B8 code into fh.
func writeB8CCode(fh afero.File, m map[[8]byte]int, keys [][8]byte) {
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
