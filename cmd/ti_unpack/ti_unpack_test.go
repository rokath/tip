package main

import (
	"testing"

	"github.com/rokath/tip/pkg/tip"
	"github.com/tj/assert"
)

var idTable = []byte{3, 'A', 'B', 'C', 0}

func Test_main(t *testing.T) {
	var pkg []byte
	if tip.OptimizeUnreplacablesEnabled(){
		if tip.UnreplacableBitCount() == 7 {
			pkg = []byte{0xd1, 0xd2, 0xd3, 0x01}
		} else { // == 6
			pkg = []byte{0xd1, 0xd2, 0xd3, 0x01}
		}
	}else{
		if tip.UnreplacableBitCount() == 7 {
			pkg = []byte{0x87, 0xd1, 0xd2, 0xd3, 0x01}
		} else { // == 6
			pkg = []byte{0xff, 0xd1, 0xd2, 0xd3, 0x01}
		}
	}
	exp:=[]byte{0xd1, 0xd2, 0xd3, 'A', 'B', 'C'}
	out := make([]byte, 1000)
	n := tip.TIUnpack(out, idTable, pkg)
	assert.Equal(t, 6, n)
	act := out[:n]
	assert.Equal(t, exp, act)
}
