package shift

import (
	"slices"
	"testing"

	"github.com/tj/assert"
)

// tt contains test data. All u7 have a cleared bit 7 for better readability. The code to test set them all.
// The msbits are stored from bit position 6-0 beginning with the last byte, because the data are processed from the back.
// The first msbits-byte carries 1-7 msbits and all further msbits-bytes carry 7 msbits.
var tt = []struct{ u8, u7 []byte }{
	{[]byte{0x00}, []byte{0x00, 0x00}},
	{[]byte{0xff}, []byte{0x40, 0x7f}},
	{[]byte{0x85}, []byte{0x40, 0x05}},
	{[]byte{0x05, 0x06}, []byte{0x00, 0x05, 0x06}},
	{[]byte{0x05, 0x86}, []byte{0x40, 0x05, 0x06}},
	{[]byte{0x85, 0x06}, []byte{0x20, 0x05, 0x06}},
	{[]byte{0x85, 0x86}, []byte{0x60, 0x05, 0x06}},
	{[]byte{0x05, 0x06, 0x07}, []byte{0x00, 0x05, 0x06, 0x07}},
	{[]byte{0x05, 0x06, 0x07, 0x08}, []byte{0x00, 0x05, 0x06, 0x07, 0x08}},
	{[]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}, []byte{0x00, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}},
	{[]byte{0x91, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}, []byte{0x01, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}}, // msb0 is msb of first byte
	{[]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x97}, []byte{0x40, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}}, // msb6 is msb of last byte
	{[]byte{0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97}, []byte{0x7f, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}},
	{[]byte{0x01, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97}, []byte{0x7e, 0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}},
	{[]byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6}, []byte{0x7f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76}},
	{[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, []byte{0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}},
	{[]byte{0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x9a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97}, []byte{0x7f, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x7f, 0x1a, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17}},
	{[]byte{0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89}, []byte{0x60, 0x01, 0x02, 0x7f, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}},
	{[]byte{0x0a, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, []byte{0x00, 0x0a, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}},
	{[]byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}, []byte{0x40, 0x70, 0x7f, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77}},
}

func Test_Shift87bit(t *testing.T) {
	for _, x := range tt {
		act := Shift87bit(x.u8)
		for i, y := range act {
			assert.Equal(t, byte(0x80), byte(0x80)&y) // check bit 7
			act[i] = byte(0x7f) & y                   // remove bit 7
		}
		assert.Equal(t, x.u7, act)
	}
}

func Test_Shift78bit(t *testing.T) {
	for _, x := range tt {
		tmp := slices.Clone(x.u7)
		for i, y := range tmp {
			tmp[i] = byte(0x80) | y // set bit 7
		}
		act := Shift78bit(tmp)
		assert.Equal(t, x.u8, act)
	}
}
