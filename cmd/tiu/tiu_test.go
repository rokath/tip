package main

import (
	"testing"

	"github.com/rokath/tip/pkg/tip"
	"github.com/tj/assert"
)

func Test_main(t *testing.T) {
	in := []byte{0x01, 0x88, 0x88, 0x01}
	out := make([]byte, 1000)
	n := tip.Unpack(out, in)
	assert.Equal(t, 4, n)
}
