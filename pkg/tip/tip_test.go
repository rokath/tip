package tip

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestTip(t *testing.T) {
	tt := []struct {
		in  []byte
		exp []byte
	}{
		{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}},
		{[]byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}},
	}

	buf := make([]byte, 100)
	for _, x := range tt {
		n := Pack(buf, x.in)
		act := buf[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.exp, act)
	}
}

func TestTiu(t *testing.T) {
	tt := []struct {
		in  []byte
		exp []byte
	}{
		{[]byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}},
		{[]byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}, []byte{0xA0, 0xbb, 0xcc, 0xaa, 0xbb}},
	}

	buf := make([]byte, 100)
	for _, x := range tt {
		n := Unpack(buf, x.in)
		act := buf[:n]
		assert.Equal(t, x.exp, act)
	}
}

func TestBuffer(t *testing.T) {
	in := [][]byte{
		{0xaa, 0xbb, 0xcc, 0xaa, 0xbb},
		{0xFa, 0xbb, 0xcc, 0xaa, 0xbb},
	}
	buf := make([]byte, 100)
	out := make([]byte, 100)
	var ratio float64
	var i uint
	for _, x := range in {
		n := Pack(buf, x)
		act := buf[:n]
		assertNoZeroes(t, act)

		m := Unpack(out, act)
		res := out[:m]
		assert.Equal(t, x, res)

		ratio += float64(n) / float64(len(x))
		i++
	}
	fmt.Println("ratio ", ratio/float64(i))
}

func assertNoZeroes(t *testing.T, b []byte) {
	for _, x := range b {
		assert.NotEqual(t, x, 0)
	}
}
