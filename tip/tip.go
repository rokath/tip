package main

func main() {

}

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
