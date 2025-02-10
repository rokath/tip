package tip

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestX(t *testing.T) {
	table := []byte{2, 0xff, 0xff, 0}
	in := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	exp := []replace{
		{5, 0, 0}, // count = 5, sz & id unused
		{0, 2, 1},
		{2, 2, 1},
		{4, 2, 1},
		{4, 0, byte(len(in))},
	}
	rpl := buildReplaceList(table, in)
	fmt.Println("exp=", exp)
	fmt.Println("act=", rpl)
	assert.Equal(t, exp, rpl)
}
