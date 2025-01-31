package tiptable

import (
	"fmt"
	"strings"
)

// spaces returns a string consisting of n spaces.
func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	var s strings.Builder
	for range n {
		s.WriteString(" ")
	}
	return s.String()
}

// byteSliceAsASCII returns b as ASCII string size len. Example: "˙Aah˙B˙˙C˙˙     "
// length is used to append spaces until the string has the desired length.
func byteSliceAsASCII(b []byte, length int) string {
	var s strings.Builder
	for _, x := range b {
		if 0x20 <= x && x < 0x7f {
			s.WriteString(fmt.Sprintf("%c", x))
		} else {
			s.WriteString(`˙`)
		}
	}
	s.WriteString(spaces(length - len(b)))
	return s.String()
}
