package main

import (
	"encoding/hex"
	"fmt"

	"github.com/rokath/tip/pkg/tip"
)

func main() {
	in := []byte{0x01, 0x88, 0x88, 0x01}
	out := make([]byte, 1000)
	n := tip.Pack(out, in)
	fmt.Println(hex.Dump(out[:n]))
}
