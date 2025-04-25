package tip

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/tj/assert"
)

var ppptable = []byte{
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  1: id 1-8
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  2:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  3:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  4:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  5:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  6:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  7:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  8:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', //  9:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', // 10:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', // 11:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', // 12:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', // 13:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', // 14:
	2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', 2, 'x', 'x', // 15: id 113...120
	2, 'p', 'p', // 121
	2, 'x', 'x', // 122
	2, 'x', 'x', // 123
	2, 'x', 'x', // 124
	2, 'x', 'x', // 125
	3, 'p', 'p', 'p', // id 126
	4, 'p', 0xbb, 0xcc, 0xdd, // id 127
	0,
}

var tipTestTable = []struct {
	UnreplacableContainerBitCount int
	unpacked                      []byte
	packed                        []byte
}{         // 0x70 0x70 0x70 0x41 0x70 0x70 0x70         0x7e, 0x80, 0x7e, 0xc1
	{7, []byte{'p', 'p', 'p', 'A', 'p', 'p', 'p'}, []byte{126, 0x80, 126, 0x80 | 'A'}},                             // only unreplacable 1 byte, not optimizable because msb==0
	{7, []byte{'p', 'p', 'p', 0xbb, 'p', 'p', 'p'}, []byte{126, 0xbb, 126}},                                        // only unreplacable 1 byte in the middle, optimizable
	{7, []byte{'p', 'p', 'p', 'p', 'p', 'p', 'A'}, []byte{126, 126, 0x80, 0x80 | 'A'}},                             // only unreplacable 1 byte at the end, not optimizable
	{7, []byte{'p', 'p', 'p', 'p', 'p', 'p', 0xbb}, []byte{126, 126, 0xbb}},                                        // only unreplacable 1 byte at the end, optimizable
	{7, []byte{0x77, 0x77, 'p', 'p', 'p'}, []byte{0x80, 0xf7, 126, 0xf7}},                                          // 1 pattern in the end, not optimizable
	{7, []byte{0xf7, 0xf7, 'p', 'p', 'p'}, []byte{0xf7, 0xf7, 126}},                                                // 1 pattern in the end, optimizable
	{7, []byte{0xf7, 0x77, 'p', 'p', 'p'}, []byte{0x82, 0xf7, 126, 0xf7}},                                          // 1 pattern in the end, not optimizable
	{7, []byte{0x77, 0xf7, 'p', 'p', 'p'}, []byte{0x81, 0xf7, 126, 0xf7}},                                          // 1 pattern in the end, not optimizable
	{7, []byte{0xd1, 'A', 'p', 'p', 'p'}, []byte{0x82, 0xd1, 126, 0x80 | 'A'}},                                     // 1 pattern in the end, not optimizable
	{7, []byte{'A', 0xbb}, []byte{0x81, 0x80 | 'A', 0xbb}},                                                         // only unreplacable bytes, not optimizable
	{7, []byte{'A', 'B', 'C', 'A', 'B'}, []byte{0x80, 0x80 | 'A', 0x80 | 'B', 0x80 | 'C', 0x80 | 'A', 0x80 | 'B'}}, // only unreplacable bytes, not optimizable
	{7, []byte{0x41, 0x42, 0x43, 0x41, 0x42}, []byte{0x80, 0xc1, 0xc2, 0xc3, 0xc1, 0xc2}},                          // only unreplacable bytes, not optimizable
	{7, []byte{0xd1, 'p', 'p', 'p', 0xd2}, []byte{0x83, 126, 0xd1, 0xd2}},                                          // 1 pattern in the middle, not optimizable
	{7, []byte{'p', 'p', 'p', 0xd1, 0xd2}, []byte{126, 0x83, 0xd1, 0xd2}},                                          // 1 pattern at start, not optimizable
	{7, []byte{'p', 'p', 'p'}, []byte{126}},                                                                        // Just 1 pattern
	{7, []byte{'p', 'p', 'p', 'p', 'p', 'p'}, []byte{126, 126}},                                                    // just 2 pattern
	{6, []byte{'p', 'p', 'p'}, []byte{126}},                                                                        // Just 1 pattern
	{6, []byte{'p', 'p', 'p', 'p', 'p', 'p'}, []byte{126, 126}},                                                    // just 2 pattern
	{6, []byte{0xc3, 'p', 'p', 'p'}, []byte{0xc3, 126}},                                                            // 1 pattern in the end
	{6, []byte{'p', 'p', 'p', 0xc3}, []byte{126, 0xc3}},       // a single unreplacable (optimizable)
	{6, []byte{'p', 'p', 'p', 0x33}, []byte{126, 0xc0, 0xf3}}, // a single unreplacable (not optimizable)
	//{6, []byte{0xc1, 'p', 'p', 'p', 0xd2, 'x', 'x', 'p'}, []byte{0xc1, 126, 0xd2, 126}}, // 1 pattern in the end
	//{7, []byte{0xd1, 'p', 'p', 'p', 0x72, 'x', 'x', 'p'}, []byte{0x82, 126, 0xd1, 126, 0xf2}},                      // 2 pattern with distributed unreplacable bytes, not optimizable
	//{7, []byte{0xd1, 'p', 'p', 'p', 0xd2, 'x', 'x', 'p'}, []byte{0xd1, 126, 0xd2, 126}},                            // 2 pattern with distributed unreplacable bytes, optimizable
	//{7, []byte{'A'}, []byte{0x80, 0x80 | 'A'}},                                                                     // only unreplacable 1 byte, not optimizable
	//{7, []byte{'p'}, []byte{'p'}},                                                                                  // only unreplacable 1 byte, optimizable
	//{7, []byte{0xd1, 0xd2, 'x', 'x', 'p'}, []byte{0xd1, 0xd2, 126}},                                                // 1 pattern in the end, optimizable
	//{7, []byte{'p', 0xbb}, []byte{0x83, 'p', 0xbb}},                                                                // only unreplacable bytes, not optimizable
	//{7, []byte{'p', 0xbb, 0xcc, 'p', 0xbb}, []byte{0x9f, 'p', 0xbb, 0xcc, 'p', 0xbb}},                              // only unreplacable bytes, not optimizable
	//{7, []byte{0xd1, 'p', 'p', 'p', 0xd2, 'x', 'x', 'p', 0xd3}, []byte{0x87, 126, 0xd1, 126, 0xd2, 0xd3}},          // 2 pattern with distributed unreplacable bytes, not optimizable
}

func TestTIPack(t *testing.T) {
	packet := make([]byte, 100)
	for i, x := range tipTestTable {
		fmt.Println(i)
		fmt.Println("x.unpacked ", hex.EncodeToString(x.unpacked))
		fmt.Println("  x.packed ", hex.EncodeToString(x.packed))
		n := TIPack(packet, x.unpacked, x.UnreplacableContainerBitCount, 127, ppptable)
		fmt.Println("pack result", hex.EncodeToString(packet[:n]))
		assert.Equal(t, len(x.packed), n)
		act := packet[:n]
		assertNoZeroes(t, act)
		assert.Equal(t, x.packed, act)
	}
}

func TestTIUnpack(t *testing.T) {
	buffer := make([]byte, 100)
	for _, x := range tipTestTable {
		assertNoZeroes(t, x.packed)
		n := TIUnpack(buffer, x.packed, x.UnreplacableContainerBitCount, 127, ppptable)
		//fmt.Println("Tip unpack result:", hex.EncodeToString(buffer[:n]))
		assert.Equal(t, len(x.unpacked), n)
		act := buffer[:n]
		assert.Equal(t, x.unpacked, act)
	}
}

// TestTIPackTIUnpack packs, checks for no zeroes, unpacks and compares.
func TestTIPackTIUnpack(t *testing.T) {
	buffer := make([]byte, 100)
	packet := make([]byte, 100)
	var ratio float64
	var i uint
	for _, x := range tipTestTable {
		n := TIPack(packet, x.unpacked, x.UnreplacableContainerBitCount, 127, ppptable)
		act := packet[:n]

		assertNoZeroes(t, act)

		m := TIUnpack(buffer, act, x.UnreplacableContainerBitCount, 127, ppptable)
		res := buffer[:m]
		assert.Equal(t, x.unpacked, res)

		ratio += float64(n) / float64(len(x.unpacked))
		i++
	}
	fmt.Println("ratio ", ratio/float64(i))
}

// assertNoZeroes checks that b does not contain any zeroes.
func assertNoZeroes(t *testing.T, b []byte) {
	for _, x := range b {
		assert.NotEqual(t, x, 0)
	}
}
