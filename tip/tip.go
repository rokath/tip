package main

func main() {

}

// tiPack converts in to out and returns final lenth.
//
// Algorithm:
// * Start with tip list longest pattern and try to find a match inside in.
// * If a longest possible pattern match was found we have afterwards:
//   - preBytes match postBytes
//   - start over with preBytes and postBytes and so on until we cannot replace any pattern anymore
//   - Then we have: xx xx p7 x p0 p0 xx xx xx for example, where pp are any pattern replace bytes,
//     which all != 0 and all have MSB==0. The xx are the remaining bytes, which can have any values.
//     Of course we need the position information like:
//
// (A) in:  xx xx xx xx xx xx xx xx xx xx xx xx xx xx xx xx
// (B) in:  xx xx P7 P7 P7 P7 xx P0 P0 P0 P0 P0 P0 xx xx xx
// (C) ref:  0  0  1  1  1  1  0  1  1  1  1  1  1  0  0  0
// (D) (in) xx xx      p7     xx    p0    p0       xx xx xx
// * (A) is in and (C) is the result of the first
// Using (C) we collect the remaing bytes: xx xx xx xx xx xx in this example
// We convert them to yy yy yy yy yy yy yy

// Worst case length, when no compression is possible
// in |bits| 7-bits    |out
// -- | -- | --------- | --
// 0 |  0 | 0 * 7 + 0 |  0
// 1 |  8 | 1 * 7 + 1 |  2
// 2 | 16 | 2 * 7 + 2 |  3
// 3 | 24 | 3 * 7 + 3 |  4
// 4 | 32 | 4 * 7 + 4 |  5
// 5 | 40 | 5 * 7 + 5 |  6
// 6 | 48 | 6 * 7 + 6 |  7
// 7 | 56 | 8 * 7 + 0 |  8 (reserving 9)
// 8 | 64 | 9 * 7 + 1 | 10
// 9 | 72 |10 * 7 + 2 | 11
// ...
func tiPack(in []byte) (out []byte) {
	maxLen := 8*len(in)/7 + 1 // if no compression is possible, for each byte 1 more bit is needed
	out = make([]byte, maxLen)
	/*
		ref := make([]bool, maxLen) // ref gets 1s where matching patterns are found

		f := string(in)

		for _, x := range tipTable {

			s := string(x)

			n := strings.Index(f, s)
			ref[n] = true

			for k, y := range in {

			}

			comp := func(a, b []byte) bool {
				return true
			}
			ref[0] = slices.IndexFunc(in, comp, x)

		}
	*/
	return out
}

// ScanBuffer returns the offset of the first occurence of pattern in buf,
// or -1 if search was not found in buf.
func ScanBuffer(buf, pattern []byte) int {
	offset := 0
	ix := 0
	for ix < len(pattern) {
		b := buf[offset]
		if pattern[ix] == b {
			ix++
		} else {
			ix = 0
		}
		offset++
	}
	return offset
}
