package tip

import (
	"testing"

	"github.com/tj/assert"
)

func TestBuffer(t *testing.T) {
	in := []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}

	buf := make([]byte, 100)
	n := Pack(buf, in)

	assertNoZeroes(t, buf)
	
	out := make([]byte, 100)
	m := Unpack(out, buf[:n])
	
	assert.Equal(t, in, out[:m])
}

func assertNoZeroes(t *testing.T, b []byte) {
	for x := range b{
		assert.NotEqual(t,x,0)
	}
}
	
