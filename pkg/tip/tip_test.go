package tip

import (
	"testing"

	"github.com/tj/assert"
)

func TestBuffer(t *testing.T) {
	in := []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}
	buf := Pack(in)
	out := Unpack(buf)
	assert.Equal(t, in, out)
}
