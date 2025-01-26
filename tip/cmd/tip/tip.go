package main

func main() {

}

func tiPack(in []byte) (out []byte) {
	maxLen := 8*len(in)/7 + 1 // if no compression is possible, for each byte 1 more bit is needed
	out = make([]byte, 0, maxLen)
	return out
}
