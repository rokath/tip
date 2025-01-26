package main

import (
	"encoding/hex"
	"fmt"

	"github.com/rokath/tip/pkg/tip"
)

func main() {
	in := []byte{0x01, 0x01, 0x01, 0x01}
	out := tip.Pack(in)
	fmt.Println(hex.Dump(out))
}
