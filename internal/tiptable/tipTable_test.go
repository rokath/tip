package tiptable

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tj/assert"
)

// 
func _TestGenerateAA(t *testing.T) {
	defer Setup(t)() // This executes Setup(t) and puts the returned function into the defer list.
	testTable := []struct {
		data []byte
		exp  string
	}{
		/*{
			[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
			`
	  4, 0xaa, 0xaa, 0xaa, 0xaa, // ˙˙˙˙|      2  01 2 times
	  3, 0xaa, 0xaa, 0xaa,       // ˙˙˙ |      0  02 3 times
	  2, 0xaa, 0xaa,             // ˙˙  |      0  03`,4 times
		},*/
		{
			[]byte{0xaa, 0xaa, 0xaa, 0xaa},
			`
	  4, 0xaa, 0xaa, 0xaa, 0xaa, // ˙˙˙˙|      1  01
	  3, 0xaa, 0xaa, 0xaa,       // ˙˙˙ |      0  02
	  2, 0xaa, 0xaa,             // ˙˙  |      0  03`,

	  //  4, 0xaa, 0xaa, 0xaa, 0xaa, // ˙˙˙˙|      1  01   1 times
	  //  3, 0xaa, 0xaa, 0xaa,       // ˙˙˙ |      0  02   2 times
	  //  2, 0xaa, 0xaa,             // ˙˙  |      0  03   3 times general but 2 times in aa aa aa
      //  -2 aaaa for first aaaaaa and 2n aaaaaa also -2 but one aa is now subtracted twice

	},
		{
			[]byte{0xaa, 0xaa},
			`
	  2, 0xaa, 0xaa, // ˙˙|      1  01`,
		},
		{
			[]byte{0xaa, 0xaa, 0xaa},
			`
	  3, 0xaa, 0xaa, 0xaa, // ˙˙˙|      1  01
	  2, 0xaa, 0xaa,       // ˙˙ |      0  02`,
		},
	}
	patternSizeMax := 4
	iFn := "testData"
	oFn := iFn + ".idTable.c"

	for _, x := range testTable {
		assert.Nil(t, FSys.WriteFile(iFn, x.data, 0777))

		Generate(FSys, oFn, iFn, patternSizeMax)

		result, err := FSys.ReadFile(oFn)
		assert.Nil(t, err)
		xxx := string(result)
		_, after, _ := strings.Cut(xxx, "count  id")
		act, _, _ := strings.Cut(after, "\n\t  0 // table end marker\n};\n\n")

		fmt.Println(act)
		assert.Equal(t, x.exp, act)
	}
}
