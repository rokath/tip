package main

import (
	"testing"

	"github.com/rokath/tip/pkg/tip"
	"github.com/tj/assert"
)

func Test_main(t *testing.T) {
	var pkg []byte
	if tip.UnreplacableBitCount() == 7 {
		pkg = []byte{0x87, 0xd1, 0xd2, 0xd3}
	} else { // == 6
		pkg = []byte{0xff, 0xd1, 0xd2, 0xd3}
	}
	exp := []byte{0xd1, 0xd2, 0xd3}
	out := make([]byte, 1000)
	n := tip.Unpack(out, pkg)
	assert.Equal(t, 3, n)
	act := out[:n]
	assert.Equal(t, exp, act)
}
