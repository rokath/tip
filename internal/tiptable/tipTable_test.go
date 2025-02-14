package tiptable

import (
	"strings"
	"testing"

	"github.com/tj/assert"
)

// todo: Sort results also alphabetically to ensure equal test results.
func _TestGenerate(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	data := []byte{0x01, 0x88, 0xaa, 0xaa, 0x01, 0x88, 0x01, 0x88, 0xaa, 0xbb}
	patternSizeMax := 4
	iFn := "testData"
	oFn := iFn + ".idTable.c"
	assert.Nil(t, FSys.WriteFile(iFn, data, 0777))
	in, err := FSys.ReadFile(iFn)
	assert.Nil(t, err)
	assert.Equal(t, data, in)

	Generate(FSys, oFn, iFn, patternSizeMax)

	tt, err := FSys.ReadFile(oFn)
	assert.Nil(t, err)

	act := string(tt)

	exp := `//! @file idTable.c
//! @brief Generated code - do not edit!
#include <stdint.h>
#include <stddef.h>

//! idTable is sorted by pattern length and pattern count.
//! The pattern position + 1 is the replace id.
//! The pattern max size is 4\nuint8_t tipTable[] = { // from testData ()
                                 // ASCII|  count  id
	  4, 0xaa, 0xaa, 0x01, 0x88, // ˙˙˙˙|      1  01
	  4, 0xaa, 0x01, 0x88, 0x01, // ˙˙˙˙|      1  02
	  4, 0x88, 0xaa, 0xaa, 0x01, // ˙˙˙˙|      1  03
	  4, 0x88, 0x01, 0x88, 0xaa, // ˙˙˙˙|      1  04
	  4, 0x01, 0x88, 0xaa, 0xbb, // ˙˙˙˙|      1  05
	  4, 0x01, 0x88, 0xaa, 0xaa, // ˙˙˙˙|      1  06
	  4, 0x01, 0x88, 0x01, 0x88, // ˙˙˙˙|      1  07
	  3, 0x88, 0xaa, 0xbb,       // ˙˙˙ |      0  08
	  3, 0xaa, 0xaa, 0x01,       // ˙˙˙ |     -1  09
	  3, 0xaa, 0x01, 0x88,       // ˙˙˙ |     -1  0a
	  3, 0x88, 0xaa, 0xaa,       // ˙˙˙ |     -1  0b
	  3, 0x88, 0x01, 0x88,       // ˙˙˙ |     -1  0c
	  3, 0x01, 0x88, 0xaa,       // ˙˙˙ |     -1  0d
	  3, 0x01, 0x88, 0x01,       // ˙˙˙ |     -1  0e
	  2, 0xaa, 0xbb,             // ˙˙  |     -1  0f
	  2, 0xaa, 0xaa,             // ˙˙  |     -4  10
	  2, 0xaa, 0x01,             // ˙˙  |     -4  11
	  2, 0x88, 0x01,             // ˙˙  |     -4  12
	  2, 0x88, 0xaa,             // ˙˙  |     -6  13
	  2, 0x01, 0x88,             // ˙˙  |     -9  14
	  0 // table end marker
};

const size_t tipTableSize = 82;
`
	// remove date and check
	before, _, _ := strings.Cut(act, "(")
	_, after, _ := strings.Cut(act, ")")
	assert.Equal(t, exp, before+"()"+after)
}
