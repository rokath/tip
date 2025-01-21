package main

func main() {

}

type replacer struct {
	short byte
	long  []byte
}

var nibbleList = []replacer{
	{7, []byte{0xAA, 0x00, 0xFF}},
	{2, []byte{0xFF, 0xFF}},
	{1, []byte{0xFF, 0xAA}},
	{0, []byte{0xFF, 0x00}},
}

// Translate nibbleList

// 1 1:  0xFF, 0xAA, 0xFF, 0xAA
// 7 2:  0xAA, 0x00, 0xFF 0xFF, 0xFF
// So we have 64 double nibble pattern, we can add to the byte list
// 0bbb0bbb -> 00bbbbbb (bbbbbb==000000 is forbidden to kep 0x00 out)
// In the above pattern 0 0 []byte{ 0xFF, 0x00, 0xFF, 0x00 } we need to forbid,
// so it will be 01000000 then and we use 65...127 = 63 more pattern

// Here we can use 01bbbbbb another 64 byte pattern
var byteList = []replacer{
	{0x02, []byte{0xFF, 0x00, 0xFF}},
	{0x01, []byte{0x00}},
}

// byte list and nibble list result in a combined list sorted by pattern length:
var combinedList = []replacer{
	{0x3F, []byte{0xAA, 0x00, 0xFF, 0xAA, 0x00, 0xFF}},
	{0x02, []byte{0xFF, 0x00, 0xFF}},
	{0x01, []byte{0x00}},
}

// pack converts in to out and returns final lenth.
//
// Algorithm:
// * Start with combind list longest pattern and try to find a match inside in.
// * If a longest possible pattern match was found we have afterwards:
//   - preBytes match postBytes
//   - start over with preBytes and postBytes and so on until we cannot replace any pattern anymore
//   - Then we have: xx xx p7 x p0 p0 xx xx xx for example, where pp are any pattern replace bytes,
//     which all != 0 and all have MSB==0. The xx are the remaining bytes, which can have any values.
//     Of course we need the position information like:
//
// (A) in:  xx xx xx xx xx xx xx xx xx xx xx xx xx xx xx xx
// (B) in:  xx xx P7 P7 P7 P7 xx P0 P0 P0 P0 P0 P0 xx xx xx
// (C) ref:  0  0  1  1  1  1  0  1  1  1  1  1  1  0  0  0
// (D) (in) xx xx      p7     xx    p0    p0       xx xx xx
// * (A) is in and (C) is the result of the first
// Using (C) we collect the remaing bytes: xx xx xx xx xx xx in this example
// We convert them to yy yy yy yy yy yy yy
func pack(in []byte, out []byte) int {

	// in: 0x00, 0xAA, 0xFF, 0xFF
	return 0
}
