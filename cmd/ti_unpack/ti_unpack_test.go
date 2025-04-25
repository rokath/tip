package main

import (
	"testing"

	"github.com/rokath/tip/pkg/tip"
	"github.com/tj/assert"
)

var idTable = []byte{3, 'A', 'B', 'C', 0}

func Test_main(t *testing.T) {
	var tt = []struct {
		ubc int
		pkg []byte
		exp []byte
	}{
		{7, []byte{0xd1, 0xd2, 0xd3, 0x01}, []byte{0xd1, 0xd2, 0xd3, 'A', 'B', 'C'} },
		{6, []byte{0xd1, 0xd2, 0xd3, 0x01}, []byte{0xd1, 0xd2, 0xd3, 'A', 'B', 'C'} },
	}
	out := make([]byte, 1000)
	for _, x := range tt {
		n := tip.TIUnpack(out, x.pkg, x.ubc, 127, idTable)
		assert.Equal(t, len(x.exp), n)
		act := out[:n]
		assert.Equal(t, x.exp, act)
	}
}
