package main

import (
	"testing"

	"github.com/rokath/tip/pkg/tip"
	"github.com/tj/assert"
)

var idTable = []byte{0}

func Test_main(t *testing.T) {
	in := []byte{0x01, 0x88, 0x88, 0x01}
	out := make([]byte, 1000)
	n := tip.TIPack(out, idTable, in)
	assert.Equal(t, 5, n)
}
