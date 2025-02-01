package tiptable

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/tj/assert"
)

var FSys *afero.Afero // ram file system for the tests

func init() {
	// All id tests should be executed only on a memory mapped file system.
	FSys = &afero.Afero{Fs: afero.NewMemMapFs()}
}

// todo: Sort results also alphabetically to ensure equal test results.
func TestGenerate(t *testing.T) {
	data := []byte{0x01, 0x88, 0x88, 0x01}
	patSizeMax := 4
	iFn := "testData"
	oFn := iFn + ".tipTable.c"
	assert.Nil(t, FSys.WriteFile(iFn, data, 0777))
	in, err := FSys.ReadFile(iFn)
	assert.Nil(t, err)
	assert.Equal(t, data, in)

	Generate(FSys, oFn, iFn, patSizeMax)

	tt, err := FSys.ReadFile(oFn)
	assert.Nil(t, err)

	act := string(tt)

	exp := `//! @file tipTable.c
//! @brief Generated code - do not edit!

#include <stdint.h>
#include <stddef.h>

//! tipTable is sorted by pattern count and pattern length.
//! The pattern position + 1 is the replacement id.
uint8_t tipTable[] = { // from testData ()-- __ASCII__|  count  id
	  4, 0x01, 0x88, 0x88, 0x01, // ˙˙˙˙|      1  01
	  3, 0x88, 0x88, 0x01,       // ˙˙˙ |      1  02
	  3, 0x01, 0x88, 0x88,       // ˙˙˙ |      1  03
	  2, 0x88, 0x88,             // ˙˙  |      1  04
	  2, 0x88, 0x01,             // ˙˙  |      1  05
	  2, 0x01, 0x88,             // ˙˙  |      1  06
	  0 // table end marker
};

const size_t tipTableSize = 23;
`

	// remove date and check
	before, _, _ := strings.Cut(act, "(")
	_, after, _ := strings.Cut(act, ")")
	assert.Equal(t, exp, before+"()"+after)
}
