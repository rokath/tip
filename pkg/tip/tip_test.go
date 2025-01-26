package tip

import (
	"testing"

	"github.com/tj/assert"
)

func TestScanBuffer(t *testing.T) {
	in := []byte{0xaa, 0xbb, 0xcc, 0xaa, 0xbb}
	pt := []byte{0xcc, 0xaa}
	assert.Equal(t, in, pt)
}
