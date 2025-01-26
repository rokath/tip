package main

import(
    "encoding/hex"
	"github.com/rokath/tip/pkg/tip"
	)

func main() {
	in := []byte{0x01, 0x01, 0x01, 0x01}
	out := tip.Unpack(in)
	fmt.Println(hex.Dump(out))
}

/*
func tip.Unpack(in []byte) (out []byte) {
	maxLen := 8*len(in)/7 + 1 // if no compression is possible, for each byte 1 more bit is needed
	out = make([]byte, 0, maxLen)
	return out
}
*/
